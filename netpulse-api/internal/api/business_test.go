package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

func TestBusinessTenantIsolationAndRBAC(t *testing.T) {
	ts, _, _ := setup(t)
	userA, tokenA := registerUser(t, ts, "org-a@example.com", "correct-horse")
	_, tokenB := registerUser(t, ts, "org-b@example.com", "correct-horse")
	_, tokenC := registerUser(t, ts, "org-c@example.com", "correct-horse")

	created, status := sessionDo(t, ts, http.MethodPost, tokenA, "/v1/orgs", `{"name":"Acme"}`)
	if status != 201 || created.Organization == nil {
		t.Fatalf("create org: %d %#v", status, created.Error)
	}
	orgA := created.Organization.ID

	_, status = sessionDo(t, ts, http.MethodPost, tokenB, "/v1/orgs", `{"name":"Other Co"}`)
	if status != 201 {
		t.Fatalf("create org b: %d", status)
	}

	foreign, status := sessionDo(t, ts, http.MethodGet, tokenB, "/v1/orgs/"+orgA, "")
	if status != 404 || foreign.Organization != nil {
		t.Fatalf("cross-org get must be 404, got %d", status)
	}
	legacy, status := sessionDo(t, ts, http.MethodGet, tokenB, "/v1/organizations/"+orgA, "")
	if status != 404 || legacy.Organization != nil {
		t.Fatalf("legacy org path must be 404 for non-members, got %d", status)
	}
	dash, status := sessionDo(t, ts, http.MethodGet, tokenB, "/v1/orgs/"+orgA+"/dashboard", "")
	if status != 404 || dash.OrgDashboard != nil {
		t.Fatalf("cross-org dashboard must be 404, got %d", status)
	}

	legacySelf, status := sessionDo(t, ts, http.MethodGet, tokenA, "/v1/organizations/"+userA, "")
	if status != 404 || legacySelf.Organization != nil {
		t.Fatalf("user id must not be treated as an organization, got %d", status)
	}
	legacyMember, status := sessionDo(t, ts, http.MethodGet, tokenA, "/v1/organizations/"+orgA, "")
	if status != 200 || legacyMember.Organization == nil || legacyMember.Organization.ID != orgA {
		t.Fatalf("member should read legacy org path, got %d", status)
	}

	invite, status := sessionDo(t, ts, http.MethodPost, tokenA, "/v1/orgs/"+orgA+"/members", `{"email":"org-c@example.com","role":"viewer"}`)
	if status != 201 || invite.Member == nil || invite.Member.Role != "viewer" {
		t.Fatalf("invite viewer: %d %#v", status, invite)
	}
	viewerID := invite.Member.ID

	mine, status := sessionDo(t, ts, http.MethodGet, tokenC, "/v1/orgs/"+orgA, "")
	if status != 200 || mine.Organization == nil {
		t.Fatalf("viewer should read org, got %d", status)
	}

	createMon, status := sessionDo(t, ts, http.MethodPost, tokenC, "/v1/orgs/"+orgA+"/monitors", `{"name":"blocked","target":"example.com","type":"http"}`)
	if status != 403 {
		t.Fatalf("viewer cannot create monitor, got %d %#v", status, createMon.Error)
	}
	billing, status := sessionDo(t, ts, http.MethodGet, tokenC, "/v1/orgs/"+orgA+"/billing", "")
	if status != 403 || billing.Billing != nil {
		t.Fatalf("viewer billing must be 403 without invoices, got %d", status)
	}
	ownerBilling, status := sessionDo(t, ts, http.MethodGet, tokenA, "/v1/orgs/"+orgA+"/billing", "")
	if status != 200 || ownerBilling.Billing == nil || ownerBilling.Billing.HasAccount || len(ownerBilling.Billing.Invoices) != 0 {
		t.Fatalf("owner billing must stay empty, got %d %#v", status, ownerBilling.Billing)
	}

	monitor, status := sessionDo(t, ts, http.MethodPost, tokenA, "/v1/orgs/"+orgA+"/monitors", `{"name":"HTTP","target":"example.com","type":"http","regions":["us-east"]}`)
	if status != 201 || monitor.Monitor == nil || monitor.Monitor.Status != "unmeasured" {
		t.Fatalf("owner create monitor: %d %#v", status, monitor.Error)
	}
	foreignMon, status := sessionDo(t, ts, http.MethodGet, tokenB, "/v1/orgs/"+orgA+"/monitors/"+monitor.Monitor.ID, "")
	if status != 404 || foreignMon.Monitor != nil {
		t.Fatalf("cross-org monitor must be 404, got %d", status)
	}

	report, status := sessionDo(t, ts, http.MethodPost, tokenA, "/v1/orgs/"+orgA+"/reports", `{"kind":"availability"}`)
	if status != 201 || report.OrgReport == nil || report.OrgReport.Availability.Measured || report.OrgReport.Latency.P95 != nil {
		t.Fatalf("empty report must stay unmeasured, got %#v", report.OrgReport)
	}
	if strings.Contains(report.OrgReport.Summary, "99.9") || strings.Contains(report.OrgReport.Summary, "87") {
		t.Fatal("report must not invent availability copy")
	}
	stolenReport, status := sessionDo(t, ts, http.MethodGet, tokenB, "/v1/orgs/"+orgA+"/reports/"+report.OrgReport.ID, "")
	if status != 404 || stolenReport.OrgReport != nil {
		t.Fatalf("cross-org report must be 404, got %d", status)
	}

	keyEnv, status := sessionDo(t, ts, http.MethodPost, tokenA, "/v1/orgs/"+orgA+"/keys", `{"name":"ci","scopes":["incident.read","diagnosis.read"]}`)
	if status != 201 || keyEnv.APIKey == nil || !strings.HasPrefix(keyEnv.KeySecret, "npo_") {
		t.Fatalf("org key: %d %#v", status, keyEnv.Error)
	}
	secret := keyEnv.KeySecret
	listed, _ := sessionDo(t, ts, http.MethodGet, tokenA, "/v1/orgs/"+orgA+"/keys", "")
	raw, _ := json.Marshal(listed)
	if bytes.Contains(raw, []byte(secret)) {
		t.Fatal("org key list must not include the raw secret")
	}
	viaKey, status := bearerDo(t, ts, http.MethodGet, secret, "/v1/orgs/"+orgA+"/incidents", "")
	if status != 200 || viaKey.Error != nil {
		t.Fatalf("org key with incident.read should list incidents, got %d", status)
	}
	wsB, _ := sessionDo(t, ts, http.MethodGet, tokenB, "/v1/orgs", "")
	if len(wsB.Organizations) == 0 {
		t.Fatal("user b should have an org")
	}
	crossKey, status := bearerDo(t, ts, http.MethodGet, secret, "/v1/orgs/"+wsB.Organizations[0].ID+"/dashboard", "")
	if status != 404 || crossKey.OrgDashboard != nil {
		t.Fatalf("org key must not read another org, got %d", status)
	}
	mint, status := bearerDo(t, ts, http.MethodPost, secret, "/v1/orgs/"+orgA+"/keys", `{"name":"nope"}`)
	if status != 403 {
		t.Fatalf("org keys cannot mint keys, got %d %#v", status, mint.Error)
	}

	role, status := sessionDo(t, ts, http.MethodPatch, tokenB, "/v1/orgs/"+orgA+"/members/"+viewerID, `{"role":"admin"}`)
	if status != 404 {
		t.Fatalf("foreign role change must be 404, got %d %#v", status, role.Error)
	}
	removed, status := sessionDo(t, ts, http.MethodDelete, tokenB, "/v1/orgs/"+orgA+"/members/"+viewerID, "")
	if status != 404 || removed.Member != nil {
		t.Fatalf("foreign remove must be 404, got %d", status)
	}

	emptyDash, status := sessionDo(t, ts, http.MethodGet, tokenA, "/v1/orgs/"+orgA+"/dashboard", "")
	if status != 200 || emptyDash.OrgDashboard == nil || emptyDash.OrgDashboard.Availability.Measured {
		t.Fatal("empty org dashboard must stay unmeasured")
	}
	if strings.Contains(emptyDash.OrgDashboard.OverallHealth, "%") {
		t.Fatal("overall health must not invent a percentage")
	}

	blockedInvite, status := sessionDo(t, ts, http.MethodPost, tokenC, "/v1/orgs/"+orgA+"/members", `{"email":"nobody@example.com","role":"analyst"}`)
	if status != 403 {
		t.Fatalf("viewer cannot invite, got %d %#v", status, blockedInvite.Error)
	}
	members, status := sessionDo(t, ts, http.MethodGet, tokenA, "/v1/orgs/"+orgA+"/members", "")
	if status != 200 {
		t.Fatalf("list members: %d", status)
	}
	ownerID := ""
	for _, item := range members.Members {
		if item.Role == "owner" {
			ownerID = item.ID
		}
	}
	if ownerID == "" {
		t.Fatal("owner member missing")
	}
	lastOwner, status := sessionDo(t, ts, http.MethodDelete, tokenA, "/v1/orgs/"+orgA+"/members/"+ownerID, "")
	if status != 400 {
		t.Fatalf("last owner cannot be removed, got %d %#v", status, lastOwner.Error)
	}

	revoked, status := sessionDo(t, ts, http.MethodPost, tokenA, "/v1/orgs/"+orgA+"/keys/"+keyEnv.APIKey.ID+"/revoke", "")
	if status != 200 || revoked.APIKey == nil || !revoked.APIKey.Revoked {
		t.Fatalf("revoke org key: %d %#v", status, revoked.Error)
	}
	dead, status := bearerDo(t, ts, http.MethodGet, secret, "/v1/orgs/"+orgA+"/incidents", "")
	if status != 401 || dead.OrgIncidents != nil {
		t.Fatalf("revoked org key must be 401, got %d", status)
	}
}

func TestBusinessInviteUnknownEmailDoesNotEnumerate(t *testing.T) {
	ts, _, _ := setup(t)
	_, token := registerUser(t, ts, "owner@example.com", "correct-horse")
	org, status := sessionDo(t, ts, http.MethodPost, token, "/v1/orgs", `{"name":" stewards "}`)
	if status != 201 {
		t.Fatalf("org %d", status)
	}
	env, status := sessionDo(t, ts, http.MethodPost, token, "/v1/orgs/"+org.Organization.ID+"/members", `{"email":"missing@example.com","role":"analyst"}`)
	if status != 201 || env.Member != nil || len(env.Invites) != 1 {
		t.Fatalf("unknown email should store an invite without a member body, got %d %#v", status, env)
	}
}
