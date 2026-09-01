package diagnostics

import (
	"testing"
	"time"

	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/contract"
	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/validation"
)

func TestAnalyzeNeverNamesAUserPathCause(t *testing.T) {
	summary := "status=200 redirects=0 bytes_read=12 truncated=false final_host=example.com"
	at := "2026-09-01T00:00:00Z"
	report := Analyze("rid", validation.Target{Raw: "example.com", Hostname: "example.com", Kind: "domain"}, []contract.Measurement{{
		ID: "measurement-http", Key: "http", Label: "HTTP", Value: "completed", Measured: true, MeasuredAt: &at, Summary: &summary,
	}}, time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC))

	if report.LikelyCause != nil {
		t.Fatal("worker facts must not become a likely cause")
	}
	if report.Confidence.Percent != nil || report.Confidence.Level != nil {
		t.Fatal("confidence must stay unset without isolation evidence")
	}
	if !report.InsufficientEvidence.Determined {
		t.Fatal("insufficient evidence must be first-class")
	}
	if len(report.Evidence) == 0 || report.Evidence[0].EvidenceClass != "measured_fact" {
		t.Fatal("measured HTTP should produce a measured fact")
	}
	for _, node := range report.Graph {
		if node.Status == "failed" {
			t.Fatalf("unknown/unmeasured must not be failed: %s", node.ID)
		}
	}
	if report.Recommendations[0].AutoExecute {
		t.Fatal("recommendations must not auto-execute")
	}
}
