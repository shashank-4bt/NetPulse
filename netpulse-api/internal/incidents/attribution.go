package incidents

import (
	"strings"

	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/contract"
)

func Title(evidenceClass, isolatedLayer, namedNetwork string) string {
	if evidenceClass == "measured_fact" && isolatedLayer == "isp" && strings.TrimSpace(namedNetwork) != "" {
		return "Measured failures isolated to network " + namedNetwork
	}
	if evidenceClass == "measured_fact" && isolatedLayer == "service" {
		return "Elevated connectivity failures observed toward the service"
	}
	return contract.ObservedFailures
}

func RejectsCausalOverclaim(title string) bool {
	lower := strings.ToLower(title)
	return strings.Contains(lower, "caused the outage") || strings.Contains(lower, "is down for everyone")
}
