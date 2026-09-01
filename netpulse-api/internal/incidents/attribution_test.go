package incidents

import "testing"

func TestDefaultTitleDoesNotBlameAnISP(t *testing.T) {
	title := Title("inferred_hypothesis", "isp", "AS123")
	if title != "Elevated connectivity failures observed" {
		t.Fatalf("got %q", title)
	}
	if RejectsCausalOverclaim("ISP X caused the outage") == false {
		t.Fatal("causal overclaim must be flagged")
	}
}

func TestMeasuredIsolationMayNameANetworkWithoutSayingItCausedAnOutage(t *testing.T) {
	title := Title("measured_fact", "isp", "AS64500")
	if title != "Measured failures isolated to network AS64500" {
		t.Fatalf("got %q", title)
	}
	if RejectsCausalOverclaim(title) {
		t.Fatal("isolation language is not a causal overclaim")
	}
}
