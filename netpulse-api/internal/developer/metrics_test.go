package developer

import (
	"testing"

	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/contract"
)

func TestPercentileNeedsFiveSamples(t *testing.T) {
	point := latencyOf(nil)
	if point.P50 != nil || point.P95 != nil || point.P99 != nil {
		t.Fatal("empty series must not invent percentiles")
	}
	one := 10
	fromOne := latencyOf([]contract.MonitorCheck{{LatencyMs: &one}})
	if fromOne.P95 != nil || fromOne.SampleCount != 1 {
		t.Fatal("one sample is not a distribution")
	}
}

func TestNearestRankPercentile(t *testing.T) {
	values := []float64{10, 20, 30, 40, 50}
	if percentile(values, 0.50) != 30 {
		t.Fatalf("p50=%v", percentile(values, 0.50))
	}
	if percentile(values, 0.95) != 50 {
		t.Fatalf("p95=%v", percentile(values, 0.95))
	}
}
