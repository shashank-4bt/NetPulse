package admin

import "strings"

const (
	PermUsersRead        = "users.read"
	PermOrgsRead         = "orgs.read"
	PermServicesRead     = "services.read"
	PermIncidentsRead    = "incidents.read"
	PermIncidentsOperate = "incidents.operate"
	PermMeasurementsRead = "measurements.read"
	PermDiagnosticsRead  = "diagnostics.read"
	PermRulesRead        = "rules.read"
	PermAbuseRead        = "abuse.read"
	PermAuditRead        = "audit.read"
	PermSystemRead       = "system.read"
	PermFlagsManage      = "flags.manage"
	PermConfigManage     = "config.manage"

	RoleOperator = "operator"
)

func AllPermissions() []string {
	return []string{
		PermUsersRead, PermOrgsRead, PermServicesRead, PermIncidentsRead, PermIncidentsOperate,
		PermMeasurementsRead, PermDiagnosticsRead, PermRulesRead, PermAbuseRead, PermAuditRead,
		PermSystemRead, PermFlagsManage, PermConfigManage,
	}
}

func HasPerm(have []string, need string) bool {
	for _, item := range have {
		if item == need {
			return true
		}
	}
	return false
}

func NormalizePerms(in []string) []string {
	if len(in) == 0 {
		return AllPermissions()
	}
	allowed := map[string]struct{}{}
	for _, perm := range AllPermissions() {
		allowed[perm] = struct{}{}
	}
	seen := map[string]struct{}{}
	out := []string{}
	for _, raw := range in {
		perm := strings.TrimSpace(raw)
		if _, ok := allowed[perm]; !ok {
			continue
		}
		if _, dup := seen[perm]; dup {
			continue
		}
		seen[perm] = struct{}{}
		out = append(out, perm)
	}
	return out
}
