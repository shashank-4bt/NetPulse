package business

import "strings"

const (
	RoleOwner         = "owner"
	RoleAdmin         = "admin"
	RoleSecurityAdmin = "security_admin"
	RoleDeveloper     = "developer"
	RoleAnalyst       = "analyst"
	RoleViewer        = "viewer"
	RoleBillingAdmin  = "billing_admin"

	PermDiagnosisRead   = "diagnosis.read"
	PermDiagnosisCreate = "diagnosis.create"
	PermIncidentRead    = "incident.read"
	PermIncidentManage  = "incident.manage"
	PermMonitorCreate   = "monitor.create"
	PermMonitorDelete   = "monitor.delete"
	PermTeamManage      = "team.manage"
	PermBillingManage   = "billing.manage"
	PermAPIManage       = "api.manage"
	PermAuditRead       = "audit.read"
)

var AllPermissions = []string{
	PermDiagnosisRead,
	PermDiagnosisCreate,
	PermIncidentRead,
	PermIncidentManage,
	PermMonitorCreate,
	PermMonitorDelete,
	PermTeamManage,
	PermBillingManage,
	PermAPIManage,
	PermAuditRead,
}

var AllRoles = []string{
	RoleOwner, RoleAdmin, RoleSecurityAdmin, RoleDeveloper, RoleAnalyst, RoleViewer, RoleBillingAdmin,
}

func PermissionsFor(role string) []string {
	switch strings.ToLower(strings.TrimSpace(role)) {
	case RoleOwner:
		return append([]string{}, AllPermissions...)
	case RoleAdmin:
		return []string{
			PermDiagnosisRead, PermDiagnosisCreate, PermIncidentRead, PermIncidentManage,
			PermMonitorCreate, PermMonitorDelete, PermTeamManage, PermAPIManage, PermAuditRead,
		}
	case RoleSecurityAdmin:
		return []string{
			PermDiagnosisRead, PermDiagnosisCreate, PermIncidentRead, PermIncidentManage,
			PermMonitorCreate, PermMonitorDelete, PermAPIManage, PermAuditRead,
		}
	case RoleDeveloper:
		return []string{PermDiagnosisRead, PermDiagnosisCreate, PermIncidentRead, PermMonitorCreate, PermMonitorDelete}
	case RoleAnalyst:
		return []string{PermDiagnosisRead, PermIncidentRead, PermAuditRead}
	case RoleViewer:
		return []string{PermDiagnosisRead, PermIncidentRead}
	case RoleBillingAdmin:
		return []string{PermBillingManage, PermDiagnosisRead, PermIncidentRead}
	default:
		return []string{}
	}
}

func KnownRole(role string) bool {
	role = strings.ToLower(strings.TrimSpace(role))
	for _, item := range AllRoles {
		if item == role {
			return true
		}
	}
	return false
}

func HasPerm(have []string, need string) bool {
	for _, item := range have {
		if item == need {
			return true
		}
	}
	return false
}

func CanReadOrg(have []string) bool {
	return HasPerm(have, PermDiagnosisRead) || HasPerm(have, PermIncidentRead) || HasPerm(have, PermTeamManage)
}

func NormalizePerms(in []string) []string {
	if len(in) == 0 {
		return append([]string{}, AllPermissions...)
	}
	allowed := map[string]struct{}{}
	for _, perm := range AllPermissions {
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
