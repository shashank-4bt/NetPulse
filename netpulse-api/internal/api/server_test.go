package api_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/accounts"
	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/api"
	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/config"
	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/contract"
	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/diagnostics"
	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/storage"
	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/storage/memory"
	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/validation"
	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/worker"
)

type stubRunner struct {
	measurements []contract.Measurement
	err          error
}

func (s stubRunner) Run(context.Context, validation.Target) ([]contract.Measurement, error) {
	return s.measurements, s.err
}

func setup(t *testing.T) (*httptest.Server, *memory.Store, *worker.Worker) {
	t.Helper()
	store := memory.New()
	svc := &diagnostics.Service{Store: store, Queue: store}
	w := &worker.Worker{Store: store, Measurements: store, Queue: store, Runner: stubRunner{
		measurements: []contract.Measurement{{
			ID: "measurement-dns", Key: "dns", Label: "DNS", Value: "1", Measured: true,
		}},
	}}
	server := &api.Server{
		Cfg:         config.Config{CORSOrigin: "http://localhost:3000", RateLimitPerMin: 40, EngineVersion: "0.9.0", AuthDevTokens: true, SessionTTLHours: 168},
		Diagnostics: svc,
		Accounts:    &accounts.Service{Accounts: store, Diagnoses: store, DevTokens: true, SessionTTL: 7 * 24 * time.Hour},
		Limiter:     store,
		StorageInfo: map[string]string{"postgres": "memory", "clickhouse": "memory", "redis": "memory"},
	}
	ts := httptest.NewServer(server.Handler())
	t.Cleanup(ts.Close)
	return ts, store, w
}

func TestHealthAndServices(t *testing.T) {
	ts, _, _ := setup(t)
	res, err := http.Get(ts.URL + "/v1/health")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		t.Fatalf("health %d", res.StatusCode)
	}
	var health contract.Envelope
	if err := json.NewDecoder(res.Body).Decode(&health); err != nil {
		t.Fatal(err)
	}
	if health.Incidents == nil {
		t.Fatal("health must encode incidents as a list, never null")
	}
	res, err = http.Get(ts.URL + "/v1/services")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		t.Fatalf("services %d", res.StatusCode)
	}
	var listed contract.Envelope
	if err := json.NewDecoder(res.Body).Decode(&listed); err != nil {
		t.Fatal(err)
	}
	if !listed.OK || len(listed.Services) == 0 {
		t.Fatal("services catalog must be editorial entries, not an invented health feed")
	}
	res, err = http.Get(ts.URL + "/v1/services/youtube")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		t.Fatalf("service %d", res.StatusCode)
	}
	res, err = http.Get(ts.URL + "/v1/incidents")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	var env contract.Envelope
	if err := json.NewDecoder(res.Body).Decode(&env); err != nil {
		t.Fatal(err)
	}
	if !env.OK || env.Incidents == nil {
		t.Fatal("incidents must be an empty list, not invented items")
	}
	if len(env.Incidents) != 0 {
		t.Fatal("empty store must not invent incidents")
	}
}

func TestCreateRejectsLocalhost(t *testing.T) {
	ts, _, _ := setup(t)
	res, err := http.Post(ts.URL+"/v1/diagnoses", "application/json", bytes.NewBufferString(`{"target":"localhost"}`))
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 403 && res.StatusCode != 400 {
		t.Fatalf("localhost should be rejected, got %d", res.StatusCode)
	}
}

func TestDiagnosisJobProducesHonestReport(t *testing.T) {
	ts, store, w := setup(t)
	res, err := http.Post(ts.URL+"/v1/diagnoses", "application/json", bytes.NewBufferString(`{"target":"example.com"}`))
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	var created contract.Envelope
	if err := json.NewDecoder(res.Body).Decode(&created); err != nil {
		t.Fatal(err)
	}
	if !created.OK || created.Diagnosis == nil {
		t.Fatal("expected queued diagnosis")
	}
	job, ok := store.Dequeue(context.Background())
	if !ok {
		t.Fatal("expected queued job")
	}
	w.ProcessOne(context.Background(), job)

	res, err = http.Get(ts.URL + "/v1/diagnoses/" + created.Diagnosis.ID)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	var got contract.Envelope
	if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
		t.Fatal(err)
	}
	if got.Diagnosis == nil || got.Diagnosis.Report == nil {
		t.Fatal("expected analyzed report")
	}
	if got.Diagnosis.Report.LikelyCause != nil {
		t.Fatal("must not invent a likely cause")
	}
	if !got.Diagnosis.Report.InsufficientEvidence.Determined {
		t.Fatal("insufficient evidence must be first-class")
	}
}

func TestUnknownServiceIsNotInvented(t *testing.T) {
	ts, _, _ := setup(t)
	res, err := http.Get(ts.URL + "/v1/services/not-a-catalog-slug")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 404 {
		t.Fatalf("unknown service %d", res.StatusCode)
	}
}

func TestUnknownDiagnosisIsNotAFailedMeasurement(t *testing.T) {
	ts, _, _ := setup(t)
	res, err := http.Get(ts.URL + "/v1/diagnoses/00000000-0000-4000-8000-000000000000")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 404 {
		t.Fatalf("got %d", res.StatusCode)
	}
}

func TestMapAggregatesAreEmptyWithoutGeography(t *testing.T) {
	ts, _, _ := setup(t)
	res, err := http.Get(ts.URL + "/v1/map/aggregates")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		t.Fatalf("map %d", res.StatusCode)
	}
	var env contract.Envelope
	if err := json.NewDecoder(res.Body).Decode(&env); err != nil {
		t.Fatal(err)
	}
	if !env.OK || env.Map == nil {
		t.Fatal("map aggregates must be present")
	}
	if env.Map.Cells == nil || len(env.Map.Cells) != 0 {
		t.Fatal("empty store must not invent map cells")
	}
	if env.Map.HasCoordinates {
		t.Fatal("must not invent coordinates")
	}
	raw, _ := json.Marshal(env.Map)
	var payload map[string]any
	if err := json.Unmarshal(raw, &payload); err != nil {
		t.Fatal(err)
	}
	if _, ok := payload["measurements"]; ok {
		t.Fatal("map payload must not include raw measurements")
	}
}

func TestMapAggregatesDrillFromStoredIncidentsWithoutGeocoding(t *testing.T) {
	ts, store, _ := setup(t)
	store.ReplaceIncidents([]contract.Incident{{
		ID:               "11111111-1111-4111-8111-111111111111",
		Title:            "Elevated connectivity failures observed",
		Severity:         "high",
		Status:           "investigating",
		Scope:            "youtube",
		StartedAt:        "2026-09-01T04:00:00Z",
		LastUpdatedAt:    "2026-09-01T05:00:00Z",
		AffectedServices: []string{"youtube"},
		Regions:          []string{"eu-west"},
		Networks:         []string{"AS64500"},
		SampleCount:      3,
	}})
	res, err := http.Get(ts.URL + "/v1/map/aggregates?parent=region:eu-west&level=network&layers=network")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	var env contract.Envelope
	if err := json.NewDecoder(res.Body).Decode(&env); err != nil {
		t.Fatal(err)
	}
	if env.Map == nil || len(env.Map.Cells) != 1 || env.Map.Cells[0].ID != "network:AS64500" {
		t.Fatal("region click must yield a network summary from stored labels")
	}
	if env.Map.Cells[0].Lon != nil || env.Map.HasCoordinates {
		t.Fatal("must not geocode region labels")
	}
}

func TestUnknownIncidentIsNotInvented(t *testing.T) {
	ts, _, _ := setup(t)
	res, err := http.Get(ts.URL + "/v1/incidents/00000000-0000-4000-8000-000000000000")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 404 {
		t.Fatalf("got %d", res.StatusCode)
	}
}

func TestServiceIntelligenceIsUnmeasuredWithoutSeries(t *testing.T) {
	ts, _, _ := setup(t)
	res, err := http.Get(ts.URL + "/v1/services/youtube")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	var env contract.Envelope
	if err := json.NewDecoder(res.Body).Decode(&env); err != nil {
		t.Fatal(err)
	}
	if env.Intelligence == nil || env.Intelligence.CurrentState != "not_measured" {
		t.Fatal("catalog health must stay not_measured without samples")
	}
	if env.Intelligence.Health != nil || env.Intelligence.Availability.Measured {
		t.Fatal("must not invent availability or a health score")
	}
	if env.Intelligence.Availability.SampleCount != 0 {
		t.Fatal("sample count must be observed zero, not a population rate")
	}
}

func TestIncidentFiltersDoNotInventMatches(t *testing.T) {
	ts, store, _ := setup(t)
	store.ReplaceIncidents([]contract.Incident{{
		ID:               "11111111-1111-4111-8111-111111111111",
		Title:            "Elevated connectivity failures observed",
		Severity:         "high",
		Status:           "investigating",
		Scope:            "youtube",
		StartedAt:        "2026-09-01T04:00:00Z",
		LastUpdatedAt:    "2026-09-01T05:00:00Z",
		AffectedServices: []string{"youtube"},
		Regions:          []string{"eu-west"},
		Networks:         []string{"AS64500"},
	}})
	res, err := http.Get(ts.URL + "/v1/incidents?service=github")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	var env contract.Envelope
	if err := json.NewDecoder(res.Body).Decode(&env); err != nil {
		t.Fatal(err)
	}
	if len(env.Incidents) != 0 || env.Page == nil || env.Page.Total != 0 {
		t.Fatal("unmatched filters must return an empty list")
	}
	res, err = http.Get(ts.URL + "/v1/incidents/11111111-1111-4111-8111-111111111111")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	var got contract.Envelope
	if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
		t.Fatal(err)
	}
	if got.Incident == nil || got.Incident.AffectedUserCount != nil {
		t.Fatal("stored incidents must not invent affected-user counts")
	}
}

func TestRateLimit(t *testing.T) {
	store := memory.New()
	if !store.Allow("test", 1) {
		t.Fatal("first should pass")
	}
	if store.Allow("test", 1) {
		t.Fatal("second in the same minute should fail")
	}
	time.Sleep(time.Millisecond)
	_ = storage.Job{}
}
