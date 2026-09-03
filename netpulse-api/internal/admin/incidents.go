package admin

import (
	"context"
	"strings"

	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/contract"
	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/id"
	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/incidents"
	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/opsconfig"
)

func (s *Service) ListIncidents(ctx context.Context) ([]contract.AdminIncident, *contract.APIError, int) {
	if s.Diagnoses == nil {
		return []contract.AdminIncident{}, nil, 200
	}
	items, err := s.Diagnoses.ListIncidents(ctx)
	if err != nil {
		return nil, apiErr("unavailable", "Incident store is unavailable."), 503
	}
	out := make([]contract.AdminIncident, 0, len(items))
	for _, item := range items {
		adminItem, errResp, status := s.decorateIncident(ctx, item)
		if errResp != nil {
			return nil, errResp, status
		}
		out = append(out, adminItem)
	}
	return out, nil, 200
}

func (s *Service) GetIncident(ctx context.Context, incidentID string) (*contract.AdminIncident, *contract.APIError, int) {
	if s.Diagnoses == nil {
		return nil, apiErr("not_found", "not found"), 404
	}
	item, err := s.Diagnoses.GetIncident(ctx, incidentID)
	if err != nil {
		return nil, apiErr("unavailable", "Incident store is unavailable."), 503
	}
	if item == nil {
		return nil, apiErr("not_found", "not found"), 404
	}
	adminItem, errResp, status := s.decorateIncident(ctx, *item)
	if errResp != nil {
		return nil, errResp, status
	}
	return &adminItem, nil, 200
}

func (s *Service) decorateIncident(ctx context.Context, item contract.Incident) (contract.AdminIncident, *contract.APIError, int) {
	notes, err := s.Store.ListIncidentNotes(ctx, item.ID)
	if err != nil {
		return contract.AdminIncident{}, apiErr("unavailable", "Incident notes are unavailable."), 503
	}
	if notes == nil {
		notes = []contract.IncidentNote{}
	}
	ov, err := s.Store.GetIncidentOverride(ctx, item.ID)
	if err != nil {
		return contract.AdminIncident{}, apiErr("unavailable", "Incident overlay is unavailable."), 503
	}
	return contract.AdminIncident{Incident: contract.NormalizeIncident(item), Notes: notes, Override: ov}, nil, 200
}

func (s *Service) Annotate(ctx context.Context, actorID, incidentID, body string) (*contract.AdminIncident, *contract.APIError, int) {
	return s.addNote(ctx, actorID, incidentID, "annotate", strings.TrimSpace(body), "annotated", "An operator annotation was stored.")
}

func (s *Service) Investigate(ctx context.Context, actorID, incidentID, body string) (*contract.AdminIncident, *contract.APIError, int) {
	item, errResp, status := s.requireIncident(ctx, incidentID)
	if errResp != nil {
		return nil, errResp, status
	}
	item.Status = "investigating"
	item.LastUpdatedAt = s.stamp()
	if err := s.Store.UpdateIncident(ctx, *item); err != nil {
		return nil, apiErr("unavailable", "Incident store is unavailable."), 503
	}
	return s.addNote(ctx, actorID, incidentID, "investigate", fallback(body, "Investigation opened."), "investigating", "An operator opened an investigation.")
}

func (s *Service) Escalate(ctx context.Context, actorID, incidentID, body string) (*contract.AdminIncident, *contract.APIError, int) {
	item, errResp, status := s.requireIncident(ctx, incidentID)
	if errResp != nil {
		return nil, errResp, status
	}
	item.Severity = nextSeverity(item.Severity)
	item.LastUpdatedAt = s.stamp()
	if err := s.Store.UpdateIncident(ctx, *item); err != nil {
		return nil, apiErr("unavailable", "Incident store is unavailable."), 503
	}
	return s.addNote(ctx, actorID, incidentID, "escalate", fallback(body, "Escalated to "+item.Severity+"."), "escalated", "An operator escalated the stored incident severity.")
}

func (s *Service) Resolve(ctx context.Context, actorID, incidentID, reason string, override bool, recoveries int, identifiedCause bool) (*contract.AdminIncident, *contract.APIError, int) {
	item, errResp, status := s.requireIncident(ctx, incidentID)
	if errResp != nil {
		return nil, errResp, status
	}
	minRecoveries := s.ConfigInt(ctx, opsconfig.IncidentMinRecoveries, incidents.MinRecoveriesToResolve)
	decision := incidents.CanMarkResolved(incidents.ResolutionInput{
		RecoverySampleCount: recoveries,
		IdentifiedCause:     identifiedCause,
	})
	if minRecoveries != incidents.MinRecoveriesToResolve && recoveries < minRecoveries && !override {
		decision = incidents.ResolutionDecision{OK: false, Reason: "Independent recoveries required by remote config were not met."}
	}
	result := "resolved"
	summary := "An operator marked the stored incident resolved."
	if !decision.OK {
		if !override {
			return nil, apiErr("validation_error", decision.Reason+" Set override to store an audited resolution anyway."), 400
		}
		if strings.TrimSpace(reason) == "" {
			return nil, apiErr("validation_error", "An override resolution requires a reason."), 400
		}
		result = "resolved_override"
		summary = "An operator overrode automated resolution gates. " + decision.Reason
		s.Audit(ctx, actorID, "incident.resolve.override", "incident:"+incidentID, result, summary+" Reason: "+strings.TrimSpace(reason))
	} else {
		s.Audit(ctx, actorID, "incident.resolve", "incident:"+incidentID, result, summary)
	}
	item.Status = "resolved"
	item.LastUpdatedAt = s.stamp()
	if err := s.Store.UpdateIncident(ctx, *item); err != nil {
		return nil, apiErr("unavailable", "Incident store is unavailable."), 503
	}
	body := fallback(reason, "Resolved.")
	return s.addNote(ctx, actorID, incidentID, "resolve", body, result, summary)
}

func (s *Service) Override(ctx context.Context, actorID, incidentID, classification, reason string) (*contract.AdminIncident, *contract.APIError, int) {
	item, errResp, status := s.requireIncident(ctx, incidentID)
	if errResp != nil {
		return nil, errResp, status
	}
	classification = strings.TrimSpace(classification)
	reason = strings.TrimSpace(reason)
	if classification == "" || reason == "" {
		return nil, apiErr("validation_error", "Override requires a classification and a reason."), 400
	}
	ov := contract.IncidentOverride{
		Classification: classification, Reason: reason, ActorID: actorID, At: s.stamp(),
	}
	if err := s.Store.SetIncidentOverride(ctx, item.ID, ov); err != nil {
		return nil, apiErr("unavailable", "Incident overlay is unavailable."), 503
	}
	s.Audit(ctx, actorID, "incident.override", "incident:"+incidentID, "overridden", "An operator overrode automated classification. Reason: "+reason)
	return s.addNote(ctx, actorID, incidentID, "override", "Classification override: "+classification+". "+reason, "overridden", "Automated classification was overridden.")
}

func (s *Service) requireIncident(ctx context.Context, incidentID string) (*contract.Incident, *contract.APIError, int) {
	if s.Diagnoses == nil {
		return nil, apiErr("not_found", "not found"), 404
	}
	item, err := s.Diagnoses.GetIncident(ctx, incidentID)
	if err != nil {
		return nil, apiErr("unavailable", "Incident store is unavailable."), 503
	}
	if item == nil {
		return nil, apiErr("not_found", "not found"), 404
	}
	copy := *item
	return &copy, nil, 200
}

func (s *Service) addNote(ctx context.Context, actorID, incidentID, kind, body, result, summary string) (*contract.AdminIncident, *contract.APIError, int) {
	if strings.TrimSpace(body) == "" {
		return nil, apiErr("validation_error", "A note body is required."), 400
	}
	note := contract.IncidentNote{
		ID: id.New(), IncidentID: incidentID, Kind: kind, Body: strings.TrimSpace(body),
		ActorID: actorID, At: s.stamp(),
	}
	if err := s.Store.AddIncidentNote(ctx, note); err != nil {
		return nil, apiErr("unavailable", "Incident notes are unavailable."), 503
	}
	if kind != "resolve" && kind != "override" {
		s.Audit(ctx, actorID, "incident."+kind, "incident:"+incidentID, result, summary)
	}
	return s.GetIncident(ctx, incidentID)
}

func fallback(value, or string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return or
	}
	return value
}

func nextSeverity(current string) string {
	switch strings.ToLower(strings.TrimSpace(current)) {
	case "low", "info":
		return "medium"
	case "medium", "warning":
		return "high"
	case "high", "error":
		return "critical"
	case "critical":
		return "critical"
	default:
		return "high"
	}
}
