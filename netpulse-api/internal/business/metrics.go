package business

import (
	"context"
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/contract"
	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/id"
)

func (s *Service) Dashboard(ctx context.Context, actor *Actor, orgID string) (contract.OrgDashboard, *contract.APIError, int) {
	if errResp, status := s.requireRead(actor, orgID); errResp != nil {
		return contract.OrgDashboard{}, errResp, status
	}
	checks, _ := s.Store.ListOrgChecks(ctx, orgID, "")
	incidents, _ := s.Store.ListOrgIncidents(ctx, orgID)
	devices, _ := s.Store.ListOrgDevices(ctx, orgID)
	networks, _ := s.Store.ListOrgNetworks(ctx, orgID)
	services, _ := s.Store.ListOrgServices(ctx, orgID)
	if incidents == nil {
		incidents = []contract.OrgIncident{}
	}
	if devices == nil {
		devices = []contract.OrgDevice{}
	}
	if networks == nil {
		networks = []contract.OrgNetwork{}
	}
	if services == nil {
		services = []contract.OrgService{}
	}
	dash := contract.EmptyOrgDashboard()
	dash.Incidents = incidents
	dash.Networks = networks
	dash.Services = services
	affected := []contract.OrgDevice{}
	seen := map[string]struct{}{}
	for _, inc := range incidents {
		for _, deviceID := range inc.DeviceIDs {
			if _, ok := seen[deviceID]; ok {
				continue
			}
			for _, device := range devices {
				if device.ID == deviceID {
					affected = append(affected, device)
					seen[deviceID] = struct{}{}
				}
			}
		}
	}
	dash.Devices = affected
	if len(checks) == 0 {
		return dash, nil, 200
	}
	dash.Availability = availabilityOf(checks)
	dash.Regions = regionalOf(checks)
	dash.OverallHealth = fmt.Sprintf("Observed from %d stored worker checks. This is not a live fleet score.", len(checks))
	dash.Summary = "Figures use stored organization checks only."
	return dash, nil, 200
}

func (s *Service) Analytics(ctx context.Context, actor *Actor, orgID string, filters map[string]string) (contract.OrgAnalytics, *contract.APIError, int) {
	if errResp, status := s.requireRead(actor, orgID); errResp != nil {
		return contract.OrgAnalytics{}, errResp, status
	}
	checks, err := s.Store.ListOrgChecks(ctx, orgID, "")
	if err != nil {
		return contract.OrgAnalytics{}, apiErr("unavailable", "Check store is unavailable."), 503
	}
	incidents, _ := s.Store.ListOrgIncidents(ctx, orgID)
	filtered := filterChecks(checks, filters)
	out := contract.EmptyOrgAnalytics()
	out.Filters = filters
	if incidents == nil {
		incidents = []contract.OrgIncident{}
	}
	out.Incidents = filterIncidents(incidents, filters)
	out.SampleCount = len(filtered)
	if len(filtered) == 0 {
		return out, nil, 200
	}
	out.Availability = availabilityOf(filtered)
	out.Latency = latencyOf(filtered)
	out.Summary = fmt.Sprintf("Filtered stored checks: %d. Percentiles stay unset below 5 latency samples.", len(filtered))
	return out, nil, 200
}

func (s *Service) ListReports(ctx context.Context, actor *Actor, orgID string) ([]contract.OrgReport, *contract.APIError, int) {
	if errResp, status := s.requireRead(actor, orgID); errResp != nil {
		return nil, errResp, status
	}
	items, err := s.Store.ListOrgReports(ctx, orgID)
	if err != nil {
		return nil, apiErr("unavailable", "Report store is unavailable."), 503
	}
	if items == nil {
		items = []contract.OrgReport{}
	}
	return items, nil, 200
}

func (s *Service) GetReport(ctx context.Context, actor *Actor, orgID, reportID string) (*contract.OrgReport, *contract.APIError, int) {
	if errResp, status := s.requireRead(actor, orgID); errResp != nil {
		return nil, errResp, status
	}
	item, err := s.Store.GetOrgReport(ctx, reportID)
	if err != nil {
		return nil, apiErr("unavailable", "Report store is unavailable."), 503
	}
	if item == nil || item.OrgID != orgID {
		return nil, apiErr("not_found", "report not found"), 404
	}
	return item, nil, 200
}

func (s *Service) GenerateReport(ctx context.Context, actor *Actor, orgID, kind string) (*contract.OrgReport, *contract.APIError, int) {
	if errResp, status := s.requireRead(actor, orgID); errResp != nil {
		return nil, errResp, status
	}
	kind = strings.ToLower(strings.TrimSpace(kind))
	switch kind {
	case "availability", "latency", "incidents", "regions", "network", "findings":
	default:
		return nil, apiErr("validation_error", "Report kind must be availability, latency, incidents, regions, network, or findings."), 400
	}
	checks, _ := s.Store.ListOrgChecks(ctx, orgID, "")
	incidents, _ := s.Store.ListOrgIncidents(ctx, orgID)
	networks, _ := s.Store.ListOrgNetworks(ctx, orgID)
	if incidents == nil {
		incidents = []contract.OrgIncident{}
	}
	if networks == nil {
		networks = []contract.OrgNetwork{}
	}
	findings := []string{}
	for _, check := range checks {
		if check.Summary != "" {
			findings = append(findings, check.Summary)
		}
	}
	if len(findings) > 20 {
		findings = findings[:20]
	}
	report := contract.OrgReport{
		ID: id.New(), OrgID: orgID, Kind: kind, Title: "Organization " + kind + " report",
		Availability: availabilityOf(checks), Latency: latencyOf(checks), Incidents: incidents,
		Regions: regionalOf(checks), Networks: networks, Findings: findings, SampleCount: len(checks),
		CreatedAt: s.stamp(),
		Summary:   "Generated from stored organization records. Empty fields stay empty.",
	}
	if len(checks) == 0 {
		report.Availability = contract.UnmeasuredObservation("Availability")
		report.Latency = contract.EmptyPercentiles()
		report.Summary = "No stored checks. This report does not invent availability or latency."
	}
	if err := s.Store.CreateOrgReport(ctx, report); err != nil {
		return nil, apiErr("unavailable", "Report store is unavailable."), 503
	}
	s.audit(ctx, orgID, actor.UserID, "report.generated", "An organization report was generated from stored records.")
	return &report, nil, 201
}

func filterChecks(checks []contract.MonitorCheck, filters map[string]string) []contract.MonitorCheck {
	out := []contract.MonitorCheck{}
	for _, check := range checks {
		if match := filters["region"]; match != "" && !strings.EqualFold(check.Region, match) {
			continue
		}
		if match := filters["network"]; match != "" && !strings.EqualFold(check.Network, match) {
			continue
		}
		if match := filters["asn"]; match != "" && !strings.EqualFold(check.ASN, match) {
			continue
		}
		if match := filters["service"]; match != "" && !strings.EqualFold(check.Service, match) {
			continue
		}
		if match := filters["endpoint"]; match != "" && !strings.EqualFold(check.Endpoint, match) {
			continue
		}
		if match := filters["device"]; match != "" && !strings.EqualFold(check.Device, match) {
			continue
		}
		out = append(out, check)
	}
	return out
}

func filterIncidents(items []contract.OrgIncident, filters map[string]string) []contract.OrgIncident {
	out := []contract.OrgIncident{}
	for _, item := range items {
		if match := filters["region"]; match != "" && !containsFold(item.Regions, match) {
			continue
		}
		out = append(out, item)
	}
	return out
}

func containsFold(items []string, want string) bool {
	for _, item := range items {
		if strings.EqualFold(item, want) {
			return true
		}
	}
	return false
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
		Value: ratio, Unit: &unit, Measured: true, SampleCount: len(checks),
		Summary: fmt.Sprintf("Successful stored checks: %d of %d.", ok, len(checks)),
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
	return contract.PercentilePoint{P50: &p50, P95: &p95, P99: &p99, SampleCount: len(samples), Summary: "Nearest-rank percentiles from stored worker latencies."}
}

func regionalOf(checks []contract.MonitorCheck) []contract.RegionalSlice {
	grouped := map[string]int{}
	for _, check := range checks {
		region := check.Region
		if region == "" {
			region = "unspecified"
		}
		grouped[region]++
	}
	keys := make([]string, 0, len(grouped))
	for key := range grouped {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	out := []contract.RegionalSlice{}
	for _, region := range keys {
		out = append(out, contract.RegionalSlice{
			Region: region, SampleCount: grouped[region], Status: "observed",
			Summary: "Requested region label from stored checks.",
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
