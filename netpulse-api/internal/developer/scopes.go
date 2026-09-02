package developer

import "strings"

const (
	ScopeMonitorsRead  = "monitors:read"
	ScopeMonitorsWrite = "monitors:write"
	ScopeIncidentsRead = "incidents:read"
	ScopeWebhooksRead  = "webhooks:read"
	ScopeWebhooksWrite = "webhooks:write"
	ScopeUsageRead     = "usage:read"
	ScopeSLARead       = "sla:read"
	ScopeAlertsRead    = "alerts:read"
	ScopeAlertsWrite   = "alerts:write"
	ScopeDashboardRead = "dashboard:read"
)

var AllowedScopes = []string{
	ScopeMonitorsRead,
	ScopeMonitorsWrite,
	ScopeIncidentsRead,
	ScopeWebhooksRead,
	ScopeWebhooksWrite,
	ScopeUsageRead,
	ScopeSLARead,
	ScopeAlertsRead,
	ScopeAlertsWrite,
	ScopeDashboardRead,
}

var AllowedWebhookEvents = []string{
	"incident.created",
	"incident.updated",
	"incident.resolved",
	"monitor.down",
	"monitor.recovered",
	"threshold.exceeded",
	"diagnosis.completed",
}

func SessionScopes() []string {
	out := make([]string, len(AllowedScopes))
	copy(out, AllowedScopes)
	return out
}

func HasScope(have []string, need string) bool {
	for _, scope := range have {
		if scope == need {
			return true
		}
	}
	return false
}

func NormalizeScopes(in []string) []string {
	if len(in) == 0 {
		return SessionScopes()
	}
	seen := map[string]struct{}{}
	out := []string{}
	allowed := map[string]struct{}{}
	for _, scope := range AllowedScopes {
		allowed[scope] = struct{}{}
	}
	for _, raw := range in {
		scope := strings.TrimSpace(raw)
		if _, ok := allowed[scope]; !ok {
			continue
		}
		if _, dup := seen[scope]; dup {
			continue
		}
		seen[scope] = struct{}{}
		out = append(out, scope)
	}
	return out
}

func KnownWebhookEvent(event string) bool {
	for _, item := range AllowedWebhookEvents {
		if item == event {
			return true
		}
	}
	return false
}
