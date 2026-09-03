package admin

import (
	"context"
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/contract"
	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/opsconfig"
)

func (s *Service) System(ctx context.Context) (contract.AdminSystem, *contract.APIError, int) {
	out := contract.EmptyAdminSystem()
	out.API = contract.HealthComponent{
		Name: "API", Status: "up", Measured: true,
		Detail: "This process answered the request. Engine version " + s.Cfg.EngineVersion + ".",
	}
	out.Worker = s.workerHealth()
	out.Queue = s.queueHealth()
	out.Database = backendHealth("Database", s.storageName("postgres"))
	out.Cache = backendHealth("Cache", s.storageName("redis"))

	recs, err := s.Store.ListAllDiagnoses(ctx)
	if err != nil {
		return contract.AdminSystem{}, apiErr("unavailable", "Diagnosis store is unavailable."), 503
	}
	checks, err := s.Store.ListAllChecks(ctx)
	if err != nil {
		return contract.AdminSystem{}, apiErr("unavailable", "Check store is unavailable."), 503
	}

	failedDiag := 0
	for _, rec := range recs {
		switch rec.Diagnosis.Status {
		case "failed", "unavailable":
			failedDiag++
		}
	}
	failedChecks := 0
	for _, check := range checks {
		if !check.OK {
			failedChecks++
		}
	}
	failureSamples := failedDiag + failedChecks
	totalSamples := len(recs) + len(checks)
	if totalSamples == 0 {
		out.MeasurementFailures = contract.UnmeasuredObservation("Measurement failures")
		out.ErrorRate = contract.UnmeasuredObservation("Error rate")
	} else {
		unit := "count"
		out.MeasurementFailures = contract.Observation{
			Value: failureSamples, Unit: &unit, Measured: true, SampleCount: totalSamples,
			Summary: fmt.Sprintf("Stored failed or unavailable diagnoses: %d. Failed stored checks: %d. This is not a live SLO.", failedDiag, failedChecks),
		}
		ratio := float64(failureSamples) / float64(totalSamples)
		ratioUnit := "ratio"
		out.ErrorRate = contract.Observation{
			Value: ratio, Unit: &ratioUnit, Measured: true, SampleCount: totalSamples,
			Summary: fmt.Sprintf("Failed stored outcomes: %d of %d. This is a stored-sample ratio, not a billed error budget.", failureSamples, totalSamples),
		}
	}

	minSamples := s.ConfigInt(ctx, opsconfig.SLAMinLatencySamples, 5)
	if minSamples < 1 {
		minSamples = 5
	}
	latencies := []float64{}
	for _, check := range checks {
		if check.LatencyMs != nil {
			latencies = append(latencies, float64(*check.LatencyMs))
		}
	}
	out.Latency = latencyPercentiles(latencies, minSamples)
	out.Summary = "System figures use this process, stored worker heartbeats, queue depth, configured backends, and stored measurements. Empty series stay unmeasured."
	return out, nil, 200
}

func (s *Service) workerHealth() contract.HealthComponent {
	if s.Worker == nil {
		return contract.HealthComponent{
			Name: "Worker", Status: "unmeasured", Measured: false,
			Detail: "No worker is wired to this API process.",
		}
	}
	snap := s.Worker.Snapshot()
	if !snap.Running || snap.LastHeartbeat == nil {
		return contract.HealthComponent{
			Name: "Worker", Status: "unmeasured", Measured: false,
			Detail: "No worker heartbeat is stored.",
		}
	}
	return contract.HealthComponent{
		Name: "Worker", Status: "running", Measured: true,
		Detail: fmt.Sprintf("Last heartbeat %s. Processed jobs: %d. Concurrency: %d.", *snap.LastHeartbeat, snap.Processed, snap.Concurrency),
	}
}

func (s *Service) queueHealth() contract.HealthComponent {
	depth := 0
	backend := "unknown"
	if s.Queue != nil {
		depth = s.Queue.Depth()
		backend = s.Queue.Backend()
	}
	return contract.HealthComponent{
		Name: "Queue", Status: "observed", Measured: true,
		Detail: fmt.Sprintf("Queue depth: %d. Backend: %s.", depth, backend),
	}
}

func (s *Service) storageName(key string) string {
	if s.StorageInfo == nil {
		return ""
	}
	return s.StorageInfo[key]
}

func backendHealth(name, backend string) contract.HealthComponent {
	if backend == "" {
		return contract.HealthComponent{
			Name: name, Status: "unmeasured", Measured: false,
			Detail: "No " + strings.ToLower(name) + " backend is configured.",
		}
	}
	return contract.HealthComponent{
		Name: name, Status: "configured", Measured: true,
		Detail: name + " adapter: " + backend + ". This is the configured backend name, not a live probe or uptime percentage.",
	}
}

func latencyPercentiles(samples []float64, minSamples int) contract.PercentilePoint {
	if len(samples) < minSamples {
		return contract.PercentilePoint{
			SampleCount: len(samples),
			Summary:     fmt.Sprintf("Observed latency sample count: %d. Percentiles require at least %d stored latencies.", len(samples), minSamples),
		}
	}
	sort.Float64s(samples)
	p50 := nearestRank(samples, 0.50)
	p95 := nearestRank(samples, 0.95)
	p99 := nearestRank(samples, 0.99)
	return contract.PercentilePoint{
		P50: &p50, P95: &p95, P99: &p99, SampleCount: len(samples),
		Summary: fmt.Sprintf("Nearest-rank percentiles from %d stored latencies.", len(samples)),
	}
}

func nearestRank(sorted []float64, p float64) float64 {
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
