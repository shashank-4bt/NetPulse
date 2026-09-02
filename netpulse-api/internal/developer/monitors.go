package developer

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/contract"
	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/id"
	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/validation"
)

type MonitorInput struct {
	Name       string
	Target     string
	Type       string
	Regions    []string
	FrequencyS int
	TimeoutS   int
	Thresholds contract.MonitorThresholds
}

func (s *Service) ListMonitors(ctx context.Context, workspaceID string) ([]contract.Monitor, *contract.APIError, int) {
	items, err := s.Store.ListMonitors(ctx, workspaceID)
	if err != nil {
		return nil, apiErr("unavailable", "Monitor store is unavailable."), 503
	}
	if items == nil {
		items = []contract.Monitor{}
	}
	return items, nil, 200
}

func (s *Service) GetMonitor(ctx context.Context, workspaceID, monitorID string) (*contract.Monitor, []contract.MonitorCheck, *contract.APIError, int) {
	item, err := s.Store.GetMonitor(ctx, monitorID)
	if err != nil {
		return nil, nil, apiErr("unavailable", "Monitor store is unavailable."), 503
	}
	if item == nil || item.WorkspaceID != workspaceID {
		return nil, nil, apiErr("not_found", "monitor not found"), 404
	}
	checks, err := s.Store.ListChecks(ctx, workspaceID, monitorID)
	if err != nil {
		return nil, nil, apiErr("unavailable", "Check store is unavailable."), 503
	}
	if checks == nil {
		checks = []contract.MonitorCheck{}
	}
	return item, checks, nil, 200
}

func (s *Service) CreateMonitor(ctx context.Context, workspaceID string, in MonitorInput) (*contract.Monitor, *contract.APIError, int) {
	item, errResp := buildMonitor(workspaceID, "", in, s.now())
	if errResp != nil {
		return nil, errResp, statusFor(errResp)
	}
	if err := s.Store.CreateMonitor(ctx, *item); err != nil {
		return nil, apiErr("unavailable", "Monitor store is unavailable."), 503
	}
	return item, nil, 201
}

func (s *Service) UpdateMonitor(ctx context.Context, workspaceID, monitorID string, in MonitorInput) (*contract.Monitor, *contract.APIError, int) {
	existing, err := s.Store.GetMonitor(ctx, monitorID)
	if err != nil {
		return nil, apiErr("unavailable", "Monitor store is unavailable."), 503
	}
	if existing == nil || existing.WorkspaceID != workspaceID {
		return nil, apiErr("not_found", "monitor not found"), 404
	}
	merged := MonitorInput{
		Name:       firstNonEmpty(in.Name, existing.Name),
		Target:     firstNonEmpty(in.Target, existing.Target),
		Type:       firstNonEmpty(in.Type, existing.Type),
		Regions:    existing.Regions,
		FrequencyS: existing.FrequencyS,
		TimeoutS:   existing.TimeoutS,
		Thresholds: existing.Thresholds,
	}
	if in.Regions != nil {
		merged.Regions = in.Regions
	}
	if in.FrequencyS != 0 {
		merged.FrequencyS = in.FrequencyS
	}
	if in.TimeoutS != 0 {
		merged.TimeoutS = in.TimeoutS
	}
	if in.Thresholds.AvailabilityBelow != nil || in.Thresholds.LatencyMsAbove != nil || in.Thresholds.ErrorRateAbove != nil {
		merged.Thresholds = in.Thresholds
	}
	item, errResp := buildMonitor(workspaceID, existing.ID, merged, s.now())
	if errResp != nil {
		return nil, errResp, statusFor(errResp)
	}
	item.CreatedAt = existing.CreatedAt
	item.CheckCount = existing.CheckCount
	item.Status = existing.Status
	item.Summary = existing.Summary
	if err := s.Store.UpdateMonitor(ctx, *item); err != nil {
		return nil, apiErr("unavailable", "Monitor store is unavailable."), 503
	}
	return item, nil, 200
}

func (s *Service) DeleteMonitor(ctx context.Context, workspaceID, monitorID string) *contract.APIError {
	ok, err := s.Store.DeleteMonitor(ctx, workspaceID, monitorID)
	if err != nil {
		return apiErr("unavailable", "Monitor store is unavailable.")
	}
	if !ok {
		return apiErr("not_found", "monitor not found")
	}
	return nil
}

func (s *Service) RunMonitor(ctx context.Context, workspaceID, monitorID, region string) (*contract.MonitorCheck, *contract.APIError, int) {
	item, _, errResp, status := s.GetMonitor(ctx, workspaceID, monitorID)
	if errResp != nil {
		return nil, errResp, status
	}
	if s.Runner == nil {
		return nil, apiErr("unavailable", "No probe runner is configured. A check was not invented."), 503
	}
	parsed := validation.ParseTarget(item.Target)
	if parsed.Err != nil {
		return nil, apiErr(parsed.Code, parsed.Err.Error()), 400
	}
	results, err := s.Runner.Run(ctx, parsed.Target)
	if err != nil {
		return nil, apiErr("unavailable", "The probe did not complete. No synthetic check was stored."), 503
	}
	ok, latency, summary := outcomeForType(item.Type, results)
	if region == "" {
		if len(item.Regions) > 0 {
			region = item.Regions[0]
		} else {
			region = "unspecified"
		}
	}
	check, errResp, status := s.RecordCheck(ctx, workspaceID, monitorID, region, ok, latency, summary)
	if errResp != nil {
		return nil, errResp, status
	}
	return check, nil, 200
}

func (s *Service) RecordCheck(ctx context.Context, workspaceID, monitorID, region string, ok bool, latency *int, summary string) (*contract.MonitorCheck, *contract.APIError, int) {
	item, err := s.Store.GetMonitor(ctx, monitorID)
	if err != nil {
		return nil, apiErr("unavailable", "Monitor store is unavailable."), 503
	}
	if item == nil || item.WorkspaceID != workspaceID {
		return nil, apiErr("not_found", "monitor not found"), 404
	}
	now := s.now().UTC().Format(time.RFC3339)
	if region == "" {
		region = "unspecified"
	}
	if summary == "" {
		if ok {
			summary = "Worker vantage check succeeded. This is not a user-path measurement."
		} else {
			summary = "Worker vantage check failed. This is not a user-path measurement."
		}
	}
	check := contract.MonitorCheck{
		ID:          id.New(),
		MonitorID:   monitorID,
		WorkspaceID: workspaceID,
		Region:      region,
		OK:          ok,
		LatencyMs:   latency,
		At:          now,
		Summary:     summary,
	}
	if err := s.Store.AddCheck(ctx, check); err != nil {
		return nil, apiErr("unavailable", "Check store is unavailable."), 503
	}
	_ = s.Store.IncrUsage(ctx, workspaceID, "measurements", 1)
	item.CheckCount++
	item.UpdatedAt = now
	if ok {
		item.Status = "up"
		item.Summary = "Last stored check succeeded from a worker vantage."
	} else {
		item.Status = "down"
		item.Summary = "Last stored check failed from a worker vantage."
	}
	_ = s.Store.UpdateMonitor(ctx, *item)
	s.applyCheckSideEffects(ctx, *item, check)
	return &check, nil, 201
}

func buildMonitor(workspaceID, existingID string, in MonitorInput, now time.Time) (*contract.Monitor, *contract.APIError) {
	name := strings.TrimSpace(in.Name)
	if name == "" {
		return nil, apiErr("validation_error", "Monitor name is required.")
	}
	kind := strings.ToLower(strings.TrimSpace(in.Type))
	if kind != "http" && kind != "dns" && kind != "tls" {
		return nil, apiErr("validation_error", "Monitor type must be http, dns, or tls.")
	}
	parsed := validation.ParseTarget(in.Target)
	if parsed.Err != nil {
		code := parsed.Code
		if code == "" {
			code = "validation_error"
		}
		return nil, apiErr(code, parsed.Err.Error())
	}
	freq := in.FrequencyS
	if freq == 0 {
		freq = 300
	}
	if freq < 60 || freq > 86400 {
		return nil, apiErr("validation_error", "Frequency must be between 60 and 86400 seconds.")
	}
	timeout := in.TimeoutS
	if timeout == 0 {
		timeout = 10
	}
	if timeout < 1 || timeout > 30 {
		return nil, apiErr("validation_error", "Timeout must be between 1 and 30 seconds.")
	}
	regions, errResp := normalizeRegions(in.Regions)
	if errResp != nil {
		return nil, errResp
	}
	if errResp := validateThresholds(in.Thresholds); errResp != nil {
		return nil, errResp
	}
	stamp := now.UTC().Format(time.RFC3339)
	itemID := existingID
	if itemID == "" {
		itemID = id.New()
	}
	return &contract.Monitor{
		ID:          itemID,
		WorkspaceID: workspaceID,
		Name:        name,
		Target:      parsed.Target.Raw,
		Type:        kind,
		Regions:     regions,
		FrequencyS:  freq,
		TimeoutS:    timeout,
		Thresholds:  in.Thresholds,
		Status:      "unmeasured",
		Summary:     "No checks are stored. Status is not estimated.",
		CreatedAt:   stamp,
		UpdatedAt:   stamp,
	}, nil
}

func normalizeRegions(in []string) ([]string, *contract.APIError) {
	out := []string{}
	seen := map[string]struct{}{}
	for _, raw := range in {
		label := strings.ToLower(strings.TrimSpace(raw))
		if label == "" {
			continue
		}
		if len(label) > 32 {
			return nil, apiErr("validation_error", "Region labels must be 32 characters or fewer.")
		}
		for _, r := range label {
			if (r < 'a' || r > 'z') && (r < '0' || r > '9') && r != '-' {
				return nil, apiErr("validation_error", "Region labels may contain letters, digits, and hyphen.")
			}
		}
		if _, ok := seen[label]; ok {
			continue
		}
		seen[label] = struct{}{}
		out = append(out, label)
	}
	if len(out) > 8 {
		return nil, apiErr("validation_error", "At most 8 requested regions are stored.")
	}
	return out, nil
}

func validateThresholds(th contract.MonitorThresholds) *contract.APIError {
	if th.AvailabilityBelow != nil && (*th.AvailabilityBelow < 0 || *th.AvailabilityBelow > 1) {
		return apiErr("validation_error", "Availability threshold must be a ratio between 0 and 1.")
	}
	if th.ErrorRateAbove != nil && (*th.ErrorRateAbove < 0 || *th.ErrorRateAbove > 1) {
		return apiErr("validation_error", "Error threshold must be a ratio between 0 and 1.")
	}
	if th.LatencyMsAbove != nil && *th.LatencyMsAbove < 1 {
		return apiErr("validation_error", "Latency threshold must be a positive millisecond value.")
	}
	return nil
}

func outcomeForType(kind string, results []contract.Measurement) (bool, *int, string) {
	key := kind
	if key == "tls" {
		key = "tls"
	}
	for _, item := range results {
		if item.Key != key {
			continue
		}
		summary := "Worker vantage measurement stored."
		if item.Summary != nil {
			summary = *item.Summary
		}
		var latency *int
		if ms, ok := asInt(item.Value); ok {
			latency = &ms
		}
		return item.Measured, latency, summary
	}
	return false, nil, "No measurement of the requested type was stored."
}

func asInt(value any) (int, bool) {
	switch typed := value.(type) {
	case int:
		return typed, true
	case int64:
		return int(typed), true
	case float64:
		return int(typed), true
	case string:
		n, err := strconv.Atoi(typed)
		return n, err == nil
	default:
		return 0, false
	}
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}

func statusFor(err *contract.APIError) int {
	if err == nil {
		return 200
	}
	switch err.Code {
	case "ssrf_blocked":
		return 403
	case "not_found":
		return 404
	case "unauthorized":
		return 401
	case "unavailable":
		return 503
	default:
		return 400
	}
}
