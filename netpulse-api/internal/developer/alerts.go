package developer

import (
	"context"
	"strings"

	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/contract"
	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/id"
)

type AlertInput struct {
	Kind      string
	MonitorID *string
	Threshold float64
	Enabled   *bool
}

func (s *Service) ListAlertRules(ctx context.Context, workspaceID string) ([]contract.AlertRule, *contract.APIError, int) {
	items, err := s.Store.ListAlertRules(ctx, workspaceID)
	if err != nil {
		return nil, apiErr("unavailable", "Alert store is unavailable."), 503
	}
	if items == nil {
		items = []contract.AlertRule{}
	}
	return items, nil, 200
}

func (s *Service) CreateAlertRule(ctx context.Context, workspaceID string, in AlertInput) (*contract.AlertRule, *contract.APIError, int) {
	kind := strings.ToLower(strings.TrimSpace(in.Kind))
	if kind != "availability" && kind != "latency" && kind != "error" && kind != "incident" {
		return nil, apiErr("validation_error", "Alert kind must be availability, latency, error, or incident."), 400
	}
	if kind != "incident" && in.Threshold < 0 {
		return nil, apiErr("validation_error", "Alert threshold must be zero or greater."), 400
	}
	if in.MonitorID != nil && *in.MonitorID != "" {
		item, err := s.Store.GetMonitor(ctx, *in.MonitorID)
		if err != nil {
			return nil, apiErr("unavailable", "Monitor store is unavailable."), 503
		}
		if item == nil || item.WorkspaceID != workspaceID {
			return nil, apiErr("not_found", "monitor not found"), 404
		}
	} else {
		in.MonitorID = nil
	}
	enabled := true
	if in.Enabled != nil {
		enabled = *in.Enabled
	}
	rule := contract.AlertRule{
		ID:             id.New(),
		WorkspaceID:    workspaceID,
		Kind:           kind,
		MonitorID:      in.MonitorID,
		Threshold:      in.Threshold,
		Enabled:        enabled,
		DeliveredCount: 0,
		Summary:        "No alert deliveries are stored. Email is not sent.",
		CreatedAt:      s.now().UTC().Format("2006-01-02T15:04:05Z07:00"),
	}
	if err := s.Store.CreateAlertRule(ctx, rule); err != nil {
		return nil, apiErr("unavailable", "Alert store is unavailable."), 503
	}
	return &rule, nil, 201
}

func (s *Service) UpdateAlertRule(ctx context.Context, workspaceID, ruleID string, in AlertInput) (*contract.AlertRule, *contract.APIError, int) {
	rules, err := s.Store.ListAlertRules(ctx, workspaceID)
	if err != nil {
		return nil, apiErr("unavailable", "Alert store is unavailable."), 503
	}
	var existing *contract.AlertRule
	for i := range rules {
		if rules[i].ID == ruleID {
			copy := rules[i]
			existing = &copy
			break
		}
	}
	if existing == nil {
		return nil, apiErr("not_found", "alert rule not found"), 404
	}
	if in.Kind != "" {
		existing.Kind = strings.ToLower(strings.TrimSpace(in.Kind))
	}
	if in.Enabled != nil {
		existing.Enabled = *in.Enabled
	}
	if in.Threshold != 0 || in.Kind != "" {
		existing.Threshold = in.Threshold
	}
	if err := s.Store.UpdateAlertRule(ctx, *existing); err != nil {
		return nil, apiErr("unavailable", "Alert store is unavailable."), 503
	}
	return existing, nil, 200
}

func (s *Service) DeleteAlertRule(ctx context.Context, workspaceID, ruleID string) *contract.APIError {
	ok, err := s.Store.DeleteAlertRule(ctx, workspaceID, ruleID)
	if err != nil {
		return apiErr("unavailable", "Alert store is unavailable.")
	}
	if !ok {
		return apiErr("not_found", "alert rule not found")
	}
	return nil
}

func (s *Service) applyCheckSideEffects(ctx context.Context, monitor contract.Monitor, check contract.MonitorCheck) {
	s.syncIncident(ctx, monitor, check)
	s.evaluateThresholds(ctx, monitor)
}

func (s *Service) syncIncident(ctx context.Context, monitor contract.Monitor, check contract.MonitorCheck) {
	incidents, err := s.Store.ListDevIncidents(ctx, monitor.WorkspaceID)
	if err != nil {
		return
	}
	var open *contract.DeveloperIncident
	for i := range incidents {
		if incidents[i].MonitorID == monitor.ID && incidents[i].Status != "resolved" {
			copy := incidents[i]
			open = &copy
			break
		}
	}
	now := check.At
	if !check.OK && open == nil {
		item := contract.DeveloperIncident{
			ID:          id.New(),
			WorkspaceID: monitor.WorkspaceID,
			MonitorID:   monitor.ID,
			Title:       monitor.Name + " is down from a stored worker check",
			Status:      "open",
			StartedAt:   now,
			SampleCount: 1,
			Summary:     "Opened from a stored failed check. This is not a global outage claim.",
		}
		_ = s.Store.CreateDevIncident(ctx, item)
		s.EnqueueEvent(ctx, monitor.WorkspaceID, "incident.created", map[string]any{"incidentId": item.ID, "monitorId": monitor.ID})
		s.EnqueueEvent(ctx, monitor.WorkspaceID, "monitor.down", map[string]any{"monitorId": monitor.ID, "checkId": check.ID})
		s.bumpIncidentAlerts(ctx, monitor.WorkspaceID, monitor.ID)
		return
	}
	if !check.OK && open != nil {
		open.SampleCount++
		_ = s.Store.UpdateDevIncident(ctx, *open)
		s.EnqueueEvent(ctx, monitor.WorkspaceID, "incident.updated", map[string]any{"incidentId": open.ID, "monitorId": monitor.ID})
		s.EnqueueEvent(ctx, monitor.WorkspaceID, "monitor.down", map[string]any{"monitorId": monitor.ID, "checkId": check.ID})
		s.bumpIncidentAlerts(ctx, monitor.WorkspaceID, monitor.ID)
		return
	}
	if check.OK && open != nil {
		open.Status = "resolved"
		open.ResolvedAt = &now
		open.SampleCount++
		open.Summary = "Resolved after a stored successful check. A single recovery is not calendar SLA proof."
		_ = s.Store.UpdateDevIncident(ctx, *open)
		s.EnqueueEvent(ctx, monitor.WorkspaceID, "incident.resolved", map[string]any{"incidentId": open.ID, "monitorId": monitor.ID})
		s.EnqueueEvent(ctx, monitor.WorkspaceID, "monitor.recovered", map[string]any{"monitorId": monitor.ID, "checkId": check.ID})
		s.bumpIncidentAlerts(ctx, monitor.WorkspaceID, monitor.ID)
	}
}

func (s *Service) evaluateThresholds(ctx context.Context, monitor contract.Monitor) {
	checks, err := s.Store.ListChecks(ctx, monitor.WorkspaceID, monitor.ID)
	if err != nil || len(checks) == 0 {
		return
	}
	successes := 0
	var lastLatency *int
	for _, check := range checks {
		if check.OK {
			successes++
		}
		if check.LatencyMs != nil {
			lastLatency = check.LatencyMs
		}
	}
	n := float64(len(checks))
	avail := float64(successes) / n
	errRate := 1 - avail
	fired := false
	if monitor.Thresholds.AvailabilityBelow != nil && avail < *monitor.Thresholds.AvailabilityBelow {
		fired = true
	}
	if monitor.Thresholds.ErrorRateAbove != nil && errRate > *monitor.Thresholds.ErrorRateAbove {
		fired = true
	}
	if monitor.Thresholds.LatencyMsAbove != nil && lastLatency != nil && *lastLatency > *monitor.Thresholds.LatencyMsAbove {
		fired = true
	}
	rules, _ := s.Store.ListAlertRules(ctx, monitor.WorkspaceID)
	for _, rule := range rules {
		if !rule.Enabled {
			continue
		}
		if rule.MonitorID != nil && *rule.MonitorID != monitor.ID {
			continue
		}
		match := false
		switch rule.Kind {
		case "availability":
			match = avail < rule.Threshold
		case "error":
			match = errRate > rule.Threshold
		case "latency":
			match = lastLatency != nil && float64(*lastLatency) > rule.Threshold
		}
		if match {
			fired = true
			rule.DeliveredCount++
			rule.Summary = "A threshold event was recorded. Email was not sent."
			_ = s.Store.UpdateAlertRule(ctx, rule)
		}
	}
	if fired {
		s.EnqueueEvent(ctx, monitor.WorkspaceID, "threshold.exceeded", map[string]any{
			"monitorId":    monitor.ID,
			"availability": avail,
			"errorRate":    errRate,
			"sampleCount":  len(checks),
		})
	}
}

func (s *Service) bumpIncidentAlerts(ctx context.Context, workspaceID, monitorID string) {
	rules, err := s.Store.ListAlertRules(ctx, workspaceID)
	if err != nil {
		return
	}
	for _, rule := range rules {
		if !rule.Enabled || rule.Kind != "incident" {
			continue
		}
		if rule.MonitorID != nil && *rule.MonitorID != monitorID {
			continue
		}
		rule.DeliveredCount++
		rule.Summary = "An incident event was recorded. Email was not sent."
		_ = s.Store.UpdateAlertRule(ctx, rule)
	}
}
