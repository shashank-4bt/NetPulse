package diagnostics

import (
	"time"

	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/contract"
	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/validation"
)

var stepOrder = []struct {
	ID    string
	Label string
}{
	{"initializing", "Initializing"},
	{"device", "Device"},
	{"wifi", "Wi-Fi"},
	{"dns", "DNS"},
	{"connectivity", "Connectivity"},
	{"isp", "ISP"},
	{"routing", "Routing"},
	{"tls", "TLS"},
	{"http", "HTTP"},
	{"service", "Service"},
	{"regional_comparison", "Regional comparison"},
	{"network_comparison", "Network comparison"},
	{"analysis", "Analysis"},
	{"complete", "Complete"},
}

var graphOrder = []struct {
	ID    string
	Label string
}{
	{"device", "Device"},
	{"wifi", "Wi-Fi"},
	{"router", "Router"},
	{"isp", "ISP"},
	{"route", "Route"},
	{"cdn", "CDN"},
	{"service", "Service"},
}

func Analyze(id string, target validation.Target, measurements []contract.Measurement, generatedAt time.Time) contract.Report {
	ts := generatedAt.UTC().Format(time.RFC3339)
	measured := map[string]contract.Measurement{}
	evidence := []contract.Evidence{}
	for _, item := range measurements {
		measured[item.Key] = item
		if !item.Measured || item.Summary == nil {
			continue
		}
		evidence = append(evidence, contract.Evidence{
			ID:             "evidence-" + item.Key,
			EvidenceClass:  "measured_fact",
			Title:          "Worker " + item.Label + " observation",
			Body:           "Observed from a NetPulse worker vantage, not the user's device. " + *item.Summary,
			MeasurementIDs: []string{item.ID},
			Layer:          item.Layer,
			ObservedAt:     item.MeasuredAt,
		})
	}

	tests := make([]contract.Step, 0, len(stepOrder))
	for _, step := range stepOrder {
		state := "unavailable"
		note := "Not measured from this worker vantage."
		switch step.ID {
		case "initializing":
			state, note = "complete", "Run accepted and enqueued."
		case "complete":
			state, note = "complete", "Analysis finished without isolating a user-path cause."
		case "analysis":
			state, note = "complete", "Worker facts were classified; root cause remains insufficient."
		case "dns", "tls", "http":
			if item, ok := measured[step.ID]; ok && item.Measured {
				state, note = "complete", "Worker probe recorded a measured observation."
			} else if item, ok := measured[step.ID]; ok {
				state, note = "unavailable", "Worker probe did not produce a measured value."
				_ = item
			}
		case "connectivity":
			if item, ok := measured["tcp"]; ok && item.Measured {
				state, note = "complete", "Worker TCP connect was recorded."
			}
		case "service":
			if item, ok := measured["http"]; ok && item.Measured {
				state, note = "complete", "Worker HTTP observation was recorded."
			}
		}
		noteCopy := note
		tests = append(tests, contract.Step{ID: step.ID, Label: step.Label, State: state, Note: &noteCopy})
	}

	graph := make([]contract.GraphNode, 0, len(graphOrder))
	for _, node := range graphOrder {
		status := "not_measured"
		ids := []string{}
		mids := []string{}
		if node.ID == "service" {
			if item, ok := measured["http"]; ok && item.Measured {
				status = "insufficient_evidence"
				mids = []string{item.ID}
				ids = []string{"evidence-http"}
			}
		}
		graph = append(graph, contract.GraphNode{
			ID:             node.ID,
			Label:          node.Label,
			Status:         status,
			Confidence:     contract.EmptyConfidence(),
			EvidenceIDs:    ids,
			MeasurementIDs: mids,
		})
	}

	filled := fillTechnical(measurements, ts)
	return contract.Report{
		ReportID:              id,
		Target:                contract.Target{Raw: target.Raw, Hostname: target.Hostname, Kind: target.Kind, ServiceSlug: target.ServiceSlug},
		Timestamp:             ts,
		Outcome:               "insufficient_evidence",
		Tests:                 tests,
		Measurements:          filled,
		Evidence:              evidence,
		Hypotheses:            []contract.Hypothesis{},
		AlternativeHypotheses: []contract.Hypothesis{},
		LikelyCause:           nil,
		Confidence:            contract.EmptyConfidence(),
		Recommendations: []contract.Recommendation{{
			ID:             "do-not-treat-worker-as-user-path",
			Action:         "Do not change device, DNS, or ISP settings based only on this worker vantage.",
			Reason:         "DNS/TCP/TLS/HTTP here were observed by NetPulse workers, not the reporter's network.",
			Risk:           "Acting as if this isolated the user path can mis-escalate an incident.",
			ExpectedResult: "A later user-path or multi-vantage run can attach more layers.",
			Verification:   "Accept a cause only when device/network layers are measured or independent evidence exists.",
			SafetyClass:    "safe",
			AutoExecute:    false,
			EvidenceIDs:    evidenceIDs(evidence),
		}},
		VerificationSteps: []contract.VerificationStep{{
			ID:     "compare-user-path",
			Label:  "Compare against a user-path measurement",
			Status: "not_run",
			Note:   "Verification requires a second measured run from another vantage.",
		}},
		EscalationConditions: []contract.EscalationCondition{{
			ID:          "no-vendor-escalation-on-worker-only",
			When:        "Only worker-vantage probes exist.",
			Action:      "Do not escalate to an ISP or vendor using this report as proof of a user-path failure.",
			SafetyClass: "advisory",
		}},
		Graph:    graph,
		Versions: contract.CurrentVersions(),
		InsufficientEvidence: contract.InsufficientEvidence{
			Determined: true,
			Message:    contract.InsufficientCause,
			NextCheck:  contract.InsufficientNext,
		},
		EngineVersion: contract.EngineVersion,
	}
}

func fillTechnical(measured []contract.Measurement, ts string) []contract.Measurement {
	byKey := map[string]contract.Measurement{}
	for _, item := range measured {
		byKey[item.Key] = item
	}
	keys := []struct {
		Key   string
		Label string
	}{
		{"dns", "DNS"},
		{"tcp", "TCP"},
		{"tls", "TLS"},
		{"http", "HTTP"},
		{"latency", "Latency"},
		{"packet_loss", "Packet loss"},
		{"routing", "Route metadata"},
		{"network", "Network/ASN"},
		{"region", "Region"},
		{"timestamp", "Timestamp"},
	}
	out := make([]contract.Measurement, 0, len(keys))
	for _, spec := range keys {
		if item, ok := byKey[spec.Key]; ok {
			out = append(out, item)
			continue
		}
		summary := "No worker value for " + spec.Label + "."
		item := contract.Measurement{ID: "measurement-" + spec.Key, Key: spec.Key, Label: spec.Label, Measured: false, Summary: &summary}
		if spec.Key == "timestamp" {
			item.Value = ts
			item.Summary = &ts
		}
		out = append(out, item)
	}
	return out
}

func evidenceIDs(items []contract.Evidence) []string {
	ids := make([]string, 0, len(items))
	for _, item := range items {
		ids = append(ids, item.ID)
	}
	return ids
}
