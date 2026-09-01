package incidents

import (
	"testing"
	"time"

	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/contract"
)

func fixtureIncidents() []contract.Incident {
	return []contract.Incident{
		{
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
		},
		{
			ID:               "22222222-2222-4222-8222-222222222222",
			Title:            "Elevated connectivity failures observed toward the service",
			Severity:         "moderate",
			Status:           "detected",
			Scope:            "github",
			StartedAt:        "2026-08-01T00:00:00Z",
			LastUpdatedAt:    "2026-08-01T01:00:00Z",
			AffectedServices: []string{"github"},
			Regions:          []string{"us-east"},
			Networks:         []string{"AS64501"},
		},
	}
}

func TestFilterDoesNotInventRows(t *testing.T) {
	got := Filter(nil, Query{Page: 1, PageSize: 20})
	if got.Page.Total != 0 || len(got.Items) != 0 {
		t.Fatal("empty store must stay empty")
	}
}

func TestFilterAndPaginate(t *testing.T) {
	now := time.Date(2026, 9, 1, 12, 0, 0, 0, time.UTC)
	got := Filter(fixtureIncidents(), Query{
		Service: "youtube", Region: "eu-west", Network: "AS64500",
		Severity: "high", Status: "investigating", Search: "elevated",
		Time: "30d", Sort: "started_desc", Page: 1, PageSize: 10, Now: now,
	})
	if got.Page.Total != 1 || got.Items[0].ID != "11111111-1111-4111-8111-111111111111" {
		t.Fatalf("unexpected filter result: %+v", got)
	}
}

func TestTimeWindowDropsOldRows(t *testing.T) {
	now := time.Date(2026, 9, 1, 12, 0, 0, 0, time.UTC)
	got := Filter(fixtureIncidents(), Query{Time: "24h", Page: 1, PageSize: 20, Now: now})
	if got.Page.Total != 1 {
		t.Fatalf("expected 1 recent incident, got %d", got.Page.Total)
	}
}

func TestSearchMissIsEmptyNotInvented(t *testing.T) {
	got := Filter(fixtureIncidents(), Query{Search: "no-such-incident", Page: 1, PageSize: 20})
	if got.Page.Total != 0 {
		t.Fatal("unmatched search must not invent a row")
	}
}
