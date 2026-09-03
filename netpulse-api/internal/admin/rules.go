package admin

import (
	"context"
	"fmt"
	"strconv"

	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/contract"
	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/id"
	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/incidents"
	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/opsconfig"
)

func (s *Service) Rules(ctx context.Context) ([]contract.DiagnosticRule, *contract.RuleOutcomes, *contract.APIError, int) {
	versions := contract.CurrentVersions()
	minRecoveries := s.ConfigInt(ctx, opsconfig.IncidentMinRecoveries, incidents.MinRecoveriesToResolve)
	minLatency := s.ConfigInt(ctx, opsconfig.SLAMinLatencySamples, 5)
	rules := []contract.DiagnosticRule{
		{
			ID: "ssrf.host", Name: "SSRF host deny", Version: versions.RuleVersion, Layer: "safety",
			Thresholds: map[string]string{"localhost": "denied", "private": "denied", "link_local": "denied", "cloud_metadata": "denied"},
			Summary:    "Diagnose targets cannot be loopback, private, link-local, or cloud metadata.",
		},
		{
			ID: "ssrf.webhook", Name: "Webhook HTTPS + SSRF", Version: versions.RuleVersion, Layer: "safety",
			Thresholds: map[string]string{"scheme": "https"},
			Summary:    "User webhook URLs must be HTTPS and pass the same SSRF host policy.",
		},
		{
			ID: "incident.resolve", Name: "Incident resolution", Version: versions.RuleVersion, Layer: "incident",
			Thresholds: map[string]string{"minRecoveries": strconv.Itoa(minRecoveries), "identifiedCause": "required"},
			Summary:    "Resolution requires independent recoveries and an identified cause unless an audited operator override is stored.",
		},
		{
			ID: "sla.latency", Name: "Latency percentiles", Version: versions.RuleVersion, Layer: "metrics",
			Thresholds: map[string]string{"minLatencySamples": strconv.Itoa(minLatency)},
			Summary:    "Percentiles are not estimated from fewer stored latencies than the configured minimum.",
		},
		{
			ID: "diagnosis.insufficient", Name: "Insufficient evidence default", Version: versions.DiagnosticEngineVersion, Layer: "analysis",
			Thresholds: map[string]string{"rootCause": "insufficient_evidence"},
			Summary:    "Worker facts are classified. Root cause stays insufficient without isolation evidence.",
		},
		{
			ID: "auth.rate_limit", Name: "Sliding-window rate limits", Version: versions.RuleVersion, Layer: "abuse",
			Thresholds: map[string]string{"window": "1m", "diagnose.rateLimitPerMin": strconv.Itoa(s.DiagnoseLimit(ctx))},
			Summary:    "Auth, diagnose, and API keys use a one-minute window. Blocks are stored as abuse events.",
		},
	}

	recs, err := s.Store.ListAllDiagnoses(ctx)
	if err != nil {
		return nil, nil, apiErr("unavailable", "Diagnosis store is unavailable."), 503
	}
	labels, err := s.Store.ListRuleLabels(ctx)
	if err != nil {
		return nil, nil, apiErr("unavailable", "Rule labels are unavailable."), 503
	}
	statuses := map[string]int{}
	for _, rec := range recs {
		statuses[rec.Diagnosis.Status]++
	}
	fp, fn := 0, 0
	for _, label := range labels {
		switch label.Kind {
		case "false_positive":
			fp++
		case "false_negative":
			fn++
		}
	}
	outcomes := &contract.RuleOutcomes{
		Statuses: statuses, FalsePositives: fp, FalseNegatives: fn, SampleCount: len(recs),
		Summary: fmt.Sprintf("Diagnostic statuses are stored worker outcomes. False positives: %d. False negatives: %d. Both stay 0 until an operator stores a label.", fp, fn),
	}
	if statuses == nil {
		outcomes.Statuses = map[string]int{}
	}
	return rules, outcomes, nil, 200
}

func (s *Service) LabelDiagnosis(ctx context.Context, actorID, diagnosisID, kind string) (*contract.RuleOutcomes, *contract.APIError, int) {
	switch kind {
	case "false_positive", "false_negative":
	default:
		return nil, apiErr("validation_error", "Kind must be false_positive or false_negative."), 400
	}
	recs, err := s.Store.ListAllDiagnoses(ctx)
	if err != nil {
		return nil, apiErr("unavailable", "Diagnosis store is unavailable."), 503
	}
	found := false
	for _, rec := range recs {
		if rec.Diagnosis.ID == diagnosisID {
			found = true
			break
		}
	}
	if !found {
		return nil, apiErr("not_found", "not found"), 404
	}
	if err := s.Store.AddRuleLabel(ctx, contract.RuleLabel{
		ID: id.New(), DiagnosisID: diagnosisID, Kind: kind, ActorID: actorID, At: s.stamp(),
	}); err != nil {
		return nil, apiErr("unavailable", "Rule labels are unavailable."), 503
	}
	s.Audit(ctx, actorID, "rule.label", "diagnosis:"+diagnosisID, kind, "An operator labeled a stored diagnosis.")
	_, outcomes, errResp, status := s.Rules(ctx)
	return outcomes, errResp, status
}
