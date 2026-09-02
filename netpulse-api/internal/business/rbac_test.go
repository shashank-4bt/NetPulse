package business

import "testing"

func TestRolePermissions(t *testing.T) {
	if !HasPerm(PermissionsFor(RoleOwner), PermBillingManage) {
		t.Fatal("owner must have billing.manage")
	}
	if HasPerm(PermissionsFor(RoleAdmin), PermBillingManage) {
		t.Fatal("admin must not have billing.manage")
	}
	if HasPerm(PermissionsFor(RoleViewer), PermTeamManage) || HasPerm(PermissionsFor(RoleViewer), PermMonitorCreate) {
		t.Fatal("viewer must not manage team or create monitors")
	}
	if !CanReadOrg(PermissionsFor(RoleViewer)) || !CanReadOrg(PermissionsFor(RoleBillingAdmin)) {
		t.Fatal("viewer and billing admin must be able to read org dashboards")
	}
	if HasPerm(PermissionsFor(RoleAnalyst), PermIncidentManage) {
		t.Fatal("analyst must not manage incidents")
	}
	if !KnownRole("security_admin") || KnownRole("superuser") {
		t.Fatal("role catalog mismatch")
	}
}
