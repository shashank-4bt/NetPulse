package opsconfig

import (
	"strconv"
	"strings"
	"time"

	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/config"
	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/contract"
)

const (
	DiagnoseTimeoutSeconds         = "diagnose.timeoutSeconds"
	DiagnoseRateLimitPerMin        = "diagnose.rateLimitPerMin"
	WorkerConcurrency              = "worker.concurrency"
	SessionTTLHours                = "session.ttlHours"
	SLAMinLatencySamples           = "sla.minLatencySamples"
	IncidentMinRecoveries          = "incident.minRecoveries"
	MonitorDefaultFrequencySeconds = "monitor.defaultFrequencySeconds"
	MonitorDefaultTimeoutSeconds   = "monitor.defaultTimeoutSeconds"
)

func Known() []string {
	return []string{
		DiagnoseTimeoutSeconds,
		DiagnoseRateLimitPerMin,
		WorkerConcurrency,
		SessionTTLHours,
		SLAMinLatencySamples,
		IncidentMinRecoveries,
		MonitorDefaultFrequencySeconds,
		MonitorDefaultTimeoutSeconds,
	}
}

func KnownKey(key string) bool {
	key = strings.TrimSpace(key)
	for _, item := range Known() {
		if item == key {
			return true
		}
	}
	return false
}

func Defaults(cfg config.Config) []contract.RemoteConfigEntry {
	stamp := time.Now().UTC().Format(time.RFC3339)
	rate := cfg.RateLimitPerMin
	if rate < 1 {
		rate = 10
	}
	concurrency := cfg.WorkerConcurrency
	if concurrency < 1 {
		concurrency = 2
	}
	ttl := cfg.SessionTTLHours
	if ttl < 1 {
		ttl = 168
	}
	return []contract.RemoteConfigEntry{
		entry(DiagnoseTimeoutSeconds, "20", stamp, "Worker probe timeout in seconds."),
		entry(DiagnoseRateLimitPerMin, strconv.Itoa(rate), stamp, "Diagnose requests allowed per IP per minute."),
		entry(WorkerConcurrency, strconv.Itoa(concurrency), stamp, "Embedded worker goroutines."),
		entry(SessionTTLHours, strconv.Itoa(ttl), stamp, "Signed-in session lifetime in hours."),
		entry(SLAMinLatencySamples, "5", stamp, "Minimum stored latencies before percentiles are computed."),
		entry(IncidentMinRecoveries, "2", stamp, "Independent recoveries required to resolve without an override."),
		entry(MonitorDefaultFrequencySeconds, "300", stamp, "Default monitor interval in seconds."),
		entry(MonitorDefaultTimeoutSeconds, "10", stamp, "Default monitor timeout in seconds."),
	}
}

func entry(key, value, at, summary string) contract.RemoteConfigEntry {
	return contract.RemoteConfigEntry{Key: key, Value: value, UpdatedAt: at, Summary: summary}
}

func Int(entries []contract.RemoteConfigEntry, key string, fallback int) int {
	for _, item := range entries {
		if item.Key != key {
			continue
		}
		parsed, err := strconv.Atoi(strings.TrimSpace(item.Value))
		if err != nil {
			return fallback
		}
		return parsed
	}
	return fallback
}

func LooksLikeSecret(key, value string) bool {
	blob := strings.ToLower(key + " " + value)
	for _, token := range []string{"password", "secret", "token", "private_key", "api_key"} {
		if strings.Contains(blob, token) {
			return true
		}
	}
	return false
}
