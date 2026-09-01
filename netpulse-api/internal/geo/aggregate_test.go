package geo

import (
	"encoding/json"
	"testing"

	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/contract"
)

func TestAggregateEmptyStoreDoesNotInventGeography(t *testing.T) {
	agg := Aggregate(nil, ParseQuery(map[string]string{}))
	if len(agg.Cells) != 0 || len(agg.IncidentRefs) != 0 {
		t.Fatal("empty store must not invent cells")
	}
	if agg.HasCoordinates || agg.TotalSamples != 0 {
		t.Fatal("empty store must not invent coordinates or samples")
	}
	if agg.Precision != "none" {
		t.Fatalf("precision %s", agg.Precision)
	}
}

func TestAggregateUsesIncidentLabelsWithoutCoordinates(t *testing.T) {
	agg := Aggregate([]contract.Incident{fixtureIncident()}, ParseQuery(map[string]string{
		"level":  "world",
		"layers": "global,regional,network,service,incidents",
	}))
	if agg.HasCoordinates {
		t.Fatal("region labels must not be geocoded")
	}
	if got := ids(agg.Cells); !containsAll(got, "world", "region:eu-west", "network:AS64500", "service:youtube") {
		t.Fatalf("cells %v", got)
	}
	for _, cell := range agg.Cells {
		if cell.Lon != nil || cell.Lat != nil {
			t.Fatalf("cell %s must not include coordinates", cell.ID)
		}
		if cell.Status == "operational" || cell.Status == "degraded" {
			t.Fatal("must not invent a health color from incident labels")
		}
	}
	if len(agg.IncidentRefs) != 1 {
		t.Fatalf("incident refs %d", len(agg.IncidentRefs))
	}
}

func TestCountryLevelDoesNotInventCountriesFromRegionLabels(t *testing.T) {
	agg := Aggregate([]contract.Incident{fixtureIncident()}, ParseQuery(map[string]string{
		"level":  "country",
		"parent": "world",
		"layers": "regional,global",
	}))
	if len(agg.Cells) != 0 {
		t.Fatalf("country view must stay empty without country codes, got %v", ids(agg.Cells))
	}
}

func TestCountryParentHasNoInventedMapping(t *testing.T) {
	agg := Aggregate([]contract.Incident{fixtureIncident()}, ParseQuery(map[string]string{
		"level":  "region",
		"parent": "country:us",
	}))
	if len(agg.Cells) != 0 {
		t.Fatalf("country codes are not stored, got %v", ids(agg.Cells))
	}
}

func TestRegionParentYieldsNetworkSummary(t *testing.T) {
	agg := Aggregate([]contract.Incident{fixtureIncident()}, ParseQuery(map[string]string{
		"level":  "network",
		"parent": "region:eu-west",
		"layers": "network",
	}))
	if got := ids(agg.Cells); len(got) != 1 || got[0] != "network:AS64500" {
		t.Fatalf("network summary %v", got)
	}
}

func TestNetworkParentYieldsServiceSummary(t *testing.T) {
	agg := Aggregate([]contract.Incident{fixtureIncident()}, ParseQuery(map[string]string{
		"level":  "service",
		"parent": "network:AS64500",
		"layers": "service",
	}))
	if got := ids(agg.Cells); len(got) != 1 || got[0] != "service:youtube" {
		t.Fatalf("service summary %v", got)
	}
}

func TestLayerFilterOmitsOtherLayers(t *testing.T) {
	agg := Aggregate([]contract.Incident{fixtureIncident()}, ParseQuery(map[string]string{
		"layers": "regional",
	}))
	for _, cell := range agg.Cells {
		if cell.Layer != LayerRegional {
			t.Fatalf("unexpected layer %s", cell.Layer)
		}
	}
	if len(agg.IncidentRefs) != 0 {
		t.Fatal("incidents layer was off")
	}
}

func TestServiceFilterDoesNotInventMatches(t *testing.T) {
	agg := Aggregate([]contract.Incident{fixtureIncident()}, ParseQuery(map[string]string{
		"service": "github",
	}))
	if len(agg.Cells) != 0 || agg.TotalSamples != 0 {
		t.Fatal("unmatched service must yield empty aggregates")
	}
}

func TestViewportAndCapUseAggregatesNotRawMeasurements(t *testing.T) {
	cells := make([]contract.MapCell, 0, 300)
	for i := 0; i < 300; i++ {
		lon := float64(i%40) + 0.44
		lat := float64(i%20) + 0.44
		cells = append(cells, contract.MapCell{
			ID:          "fixture:" + itoa(i),
			Level:       LevelRegion,
			Label:       "cell-" + itoa(i),
			Lon:         &lon,
			Lat:         &lat,
			SampleCount: 1,
			Layer:       LayerRegional,
		})
	}
	west, south, east, north := -1.0, -1.0, 5.0, 5.0
	agg := Finalize(cells, nil, Query{
		Level:  LevelRegion,
		Limit:  MaxCells,
		Layers: DefaultLayers,
		West:   &west,
		South:  &south,
		East:   &east,
		North:  &north,
	})
	if len(agg.Cells) > MaxCells {
		t.Fatalf("cap exceeded: %d", len(agg.Cells))
	}
	raw, err := json.Marshal(agg)
	if err != nil {
		t.Fatal(err)
	}
	var payload map[string]any
	if err := json.Unmarshal(raw, &payload); err != nil {
		t.Fatal(err)
	}
	if _, ok := payload["measurements"]; ok {
		t.Fatal("map payload must not include raw measurements")
	}
	listed, _ := payload["cells"].([]any)
	for _, item := range listed {
		cell, _ := item.(map[string]any)
		if _, ok := cell["measurements"]; ok {
			t.Fatal("cells must not embed measurements")
		}
		if lon, ok := cell["lon"].(float64); ok {
			if lon != float64(int(lon)) {
				t.Fatalf("sub-degree lon %v", lon)
			}
		}
	}
}

func TestClusterMergesSameDegree(t *testing.T) {
	lonA, latA := 12.3, 48.1
	lonB, latB := 12.4, 48.4
	agg := Finalize([]contract.MapCell{
		{ID: "a", Level: LevelRegion, Label: "a", Lon: &lonA, Lat: &latA, SampleCount: 2, Layer: LayerRegional},
		{ID: "b", Level: LevelRegion, Label: "b", Lon: &lonB, Lat: &latB, SampleCount: 3, Layer: LayerRegional},
	}, nil, Query{Level: LevelRegion, Limit: MaxCells, Layers: DefaultLayers})
	if len(agg.Cells) != 1 {
		t.Fatalf("expected one clustered cell, got %d", len(agg.Cells))
	}
	if agg.Cells[0].SampleCount != 5 {
		t.Fatalf("sample count %d", agg.Cells[0].SampleCount)
	}
	if *agg.Cells[0].Lon != 12 || *agg.Cells[0].Lat != 48 {
		t.Fatalf("cluster centroid %v,%v", *agg.Cells[0].Lon, *agg.Cells[0].Lat)
	}
}

func fixtureIncident() contract.Incident {
	return contract.Incident{
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
	}
}

func ids(cells []contract.MapCell) []string {
	out := make([]string, 0, len(cells))
	for _, cell := range cells {
		out = append(out, cell.ID)
	}
	return out
}

func containsAll(got []string, want ...string) bool {
	set := map[string]bool{}
	for _, id := range got {
		set[id] = true
	}
	for _, id := range want {
		if !set[id] {
			return false
		}
	}
	return true
}
