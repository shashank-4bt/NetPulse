package developer

import (
	"context"
	"fmt"
	"math"
	"sort"

	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/contract"
)

func (s *Service) Dashboard(ctx context.Context, workspaceID string) (contract.DeveloperDashboard, *contract.APIError, int) {
	checks, incidents, errResp, status := s.series(ctx, workspaceID)
	if errResp != nil {
		return contract.DeveloperDashboard{}, errResp, status
	}
	if len(checks) == 0 {
		empty := contract.EmptyDeveloperDashboard()
		empty.Incidents = incidents
		return empty, nil, 200
	}
	return contract.DeveloperDashboard{
		Availability: availabilityOf(checks),
		Latency:      latencyOf(checks),
		Incidents:    incidents,
		Regional:     regionalOf(checks),
		Summary:      "Figures are computed from stored worker checks only. They are not a user-path SLA.",
	}, nil, 200
}

func (s *Service) SLA(ctx context.Context, workspaceID string) (contract.SLAReport, *contract.APIError, int) {
	checks, incidents, errResp, status := s.series(ctx, workspaceID)
	if errResp != nil {
		return contract.SLAReport{}, errResp, status
	}
	if len(checks) == 0 {
		empty := contract.EmptySLA()
		empty.Incidents = incidents
		return empty, nil, 200
	}
	avail := availabilityOf(checks)
	return contract.SLAReport{
		Availability: avail,
		Downtime:     downtimeOf(checks),
		Latency:      latencyOf(checks),
		Incidents:    incidents,
		Regional:     regionalOf(checks),
		Summary:      "SLA figures are failed-versus-stored-check ratios, not calendar uptime and not a billed guarantee.",
	}, nil, 200
}

func (s *Service) Usage(ctx context.Context, workspaceID string) (contract.Usage, *contract.APIError, int) {
	usage, err := s.Store.GetUsage(ctx, workspaceID)
	if err != nil {
		return contract.Usage{}, apiErr("unavailable", "Usage store is unavailable."), 503
	}
	monitors, err := s.Store.ListMonitors(ctx, workspaceID)
	if err != nil {
		return contract.Usage{}, apiErr("unavailable", "Monitor store is unavailable."), 503
	}
	hooks, err := s.Store.ListWebhooks(ctx, workspaceID)
	if err != nil {
		return contract.Usage{}, apiErr("unavailable", "Webhook store is unavailable."), 503
	}
	checks, err := s.Store.ListChecks(ctx, workspaceID, "")
	if err != nil {
		return contract.Usage{}, apiErr("unavailable", "Check store is unavailable."), 503
	}
	regions := map[string]struct{}{}
	for _, item := range monitors {
		for _, region := range item.Regions {
			regions[region] = struct{}{}
		}
	}
	for _, check := range checks {
		if check.Region != "" {
			regions[check.Region] = struct{}{}
		}
	}
	usage.Monitors = len(monitors)
	usage.Webhooks = len(hooks)
	usage.Regions = len(regions)
	usage.Summary = "Usage counts are stored totals. Traffic is not estimated."
	return usage, nil, 200
}

func (s *Service) ListIncidents(ctx context.Context, workspaceID string) ([]contract.DeveloperIncident, *contract.APIError, int) {
	items, err := s.Store.ListDevIncidents(ctx, workspaceID)
	if err != nil {
		return nil, apiErr("unavailable", "Incident store is unavailable."), 503
	}
	if items == nil {
		items = []contract.DeveloperIncident{}
	}
	return items, nil, 200
}

func (s *Service) GetIncident(ctx context.Context, workspaceID, incidentID string) (*contract.DeveloperIncident, *contract.APIError, int) {
	item, err := s.Store.GetDevIncident(ctx, incidentID)
	if err != nil {
		return nil, apiErr("unavailable", "Incident store is unavailable."), 503
	}
	if item == nil || item.WorkspaceID != workspaceID {
		return nil, apiErr("not_found", "incident not found"), 404
	}
	return item, nil, 200
}

func (s *Service) series(ctx context.Context, workspaceID string) ([]contract.MonitorCheck, []contract.DeveloperIncident, *contract.APIError, int) {
	checks, err := s.Store.ListChecks(ctx, workspaceID, "")
	if err != nil {
		return nil, nil, apiErr("unavailable", "Check store is unavailable."), 503
	}
	incidents, err := s.Store.ListDevIncidents(ctx, workspaceID)
	if err != nil {
		return nil, nil, apiErr("unavailable", "Incident store is unavailable."), 503
	}
	if incidents == nil {
		incidents = []contract.DeveloperIncident{}
	}
	return checks, incidents, nil, 200
}

func availabilityOf(checks []contract.MonitorCheck) contract.Observation {
	if len(checks) == 0 {
		return contract.UnmeasuredObservation("Availability")
	}
	ok := 0
	for _, check := range checks {
		if check.OK {
			ok++
		}
	}
	ratio := float64(ok) / float64(len(checks))
	unit := "ratio"
	return contract.Observation{
		Value:       ratio,
		Unit:        &unit,
		Measured:    true,
		SampleCount: len(checks),
		Summary:     fmt.Sprintf("Successful stored checks: %d of %d. This is a worker-vantage ratio, not calendar availability.", ok, len(checks)),
	}
}

func downtimeOf(checks []contract.MonitorCheck) contract.Observation {
	if len(checks) == 0 {
		return contract.UnmeasuredObservation("Downtime")
	}
	failed := 0
	for _, check := range checks {
		if !check.OK {
			failed++
		}
	}
	ratio := float64(failed) / float64(len(checks))
	unit := "ratio"
	return contract.Observation{
		Value:       ratio,
		Unit:        &unit,
		Measured:    true,
		SampleCount: len(checks),
		Summary:     fmt.Sprintf("Failed stored checks: %d of %d. This is not calendar downtime.", failed, len(checks)),
	}
}

func latencyOf(checks []contract.MonitorCheck) contract.PercentilePoint {
	samples := []float64{}
	for _, check := range checks {
		if check.LatencyMs != nil {
			samples = append(samples, float64(*check.LatencyMs))
		}
	}
	if len(samples) < 5 {
		return contract.PercentilePoint{
			SampleCount: len(samples),
			Summary:     fmt.Sprintf("Observed latency sample count: %d. Percentiles require at least 5 stored latencies.", len(samples)),
		}
	}
	sort.Float64s(samples)
	p50 := percentile(samples, 0.50)
	p95 := percentile(samples, 0.95)
	p99 := percentile(samples, 0.99)
	return contract.PercentilePoint{
		P50:         &p50,
		P95:         &p95,
		P99:         &p99,
		SampleCount: len(samples),
		Summary:     fmt.Sprintf("Nearest-rank percentiles from %d stored worker latencies.", len(samples)),
	}
}

func regionalOf(checks []contract.MonitorCheck) []contract.RegionalSlice {
	grouped := map[string][]contract.MonitorCheck{}
	for _, check := range checks {
		region := check.Region
		if region == "" {
			region = "unspecified"
		}
		grouped[region] = append(grouped[region], check)
	}
	keys := make([]string, 0, len(grouped))
	for key := range grouped {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	out := make([]contract.RegionalSlice, 0, len(keys))
	for _, region := range keys {
		items := grouped[region]
		ok := 0
		for _, item := range items {
			if item.OK {
				ok++
			}
		}
		out = append(out, contract.RegionalSlice{
			Region:      region,
			SampleCount: len(items),
			Status:      "observed",
			Summary:     fmt.Sprintf("Stored checks for requested region %s: %d succeeded of %d. Region labels are requested vantages, not proof a region ran.", region, ok, len(items)),
		})
	}
	return out
}

func percentile(sorted []float64, p float64) float64 {
	if len(sorted) == 0 {
		return 0
	}
	rank := int(math.Ceil(p*float64(len(sorted)))) - 1
	if rank < 0 {
		rank = 0
	}
	if rank >= len(sorted) {
		rank = len(sorted) - 1
	}
	return sorted[rank]
}
