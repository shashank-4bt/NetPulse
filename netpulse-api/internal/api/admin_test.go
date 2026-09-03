package api_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/admin"
	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/contract"
)

func TestAdminNonOperatorIsNotFound(t *testing.T) {
	ts, _, _ := setup(t)
	_, token := registerUser(t, ts, "user@example.com", "correct-horse")
	env, status := sessionDo(t, ts, http.MethodGet, token, "/v1/admin/me", "")
	if status != 404 || env.Operator != nil || env.AdminSystem != nil {
		t.Fatalf("signed-in non-operator must be 404 without admin bodies, got %d %#v", status, env)
	}
	if env.Error == nil || env.Error.Code != "not_found" {
		t.Fatalf("non-operator body must be not_found, got %#v", env.Error)
	}
	if env.Error != nil && (strings.Contains(strings.ToLower(env.Error.Message), "admin") || strings.Contains(strings.ToLower(env.Error.Message), "operator")) {
		t.Fatal("404 must not say the caller is not an admin")
	}
}

func TestAdminUnauthenticated(t *testing.T) {
	ts, _, _ := setup(t)
	res, err := http.Get(ts.URL + "/v1/admin/system")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 401 {
		t.Fatalf("unauthenticated admin must be 401, got %d", res.StatusCode)
	}
}

func TestAdminSystemStaysUnmeasuredWithoutSamples(t *testing.T) {
	ts, store, _ := setup(t)
	userID, token := registerUser(t, ts, "ops@example.com", "correct-horse")
	store.AddOperator(userID, "ops@example.com", admin.RoleOperator, nil)
	env, status := sessionDo(t, ts, http.MethodGet, token, "/v1/admin/system", "")
	if status != 200 || env.AdminSystem == nil {
		t.Fatalf("system %d %#v", status, env.Error)
	}
	sys := env.AdminSystem
	if sys.API.Status != "up" || !sys.API.Measured {
		t.Fatal("API health must be this process")
	}
	if sys.Worker.Measured {
		t.Fatal("worker without a heartbeat must stay unmeasured")
	}
	if !sys.Queue.Measured || !strings.Contains(sys.Queue.Detail, "Queue depth: 0") {
		t.Fatalf("queue depth must be observed, got %s", sys.Queue.Detail)
	}
	if sys.ErrorRate.Measured || sys.ErrorRate.Value != nil {
		t.Fatal("empty error rate must stay unmeasured")
	}
	if sys.Latency.P50 != nil || sys.Latency.P95 != nil || sys.Latency.P99 != nil {
		t.Fatal("empty latency must not invent percentiles")
	}
	if sys.MeasurementFailures.Measured {
		t.Fatal("empty measurement failures must stay unmeasured")
	}
	blob := sys.Summary + sys.ErrorRate.Summary + sys.Latency.Summary + sys.API.Detail
	if strings.Contains(blob, "99.9") || strings.Contains(blob, "87%") {
		t.Fatal("system dashboard must not invent percentages")
	}
}

func TestAdminRestrictedPermissionIsForbidden(t *testing.T) {
	ts, store, _ := setup(t)
	userID, token := registerUser(t, ts, "limited@example.com", "correct-horse")
	store.AddOperator(userID, "limited@example.com", admin.RoleOperator, []string{admin.PermSystemRead})
	env, status := sessionDo(t, ts, http.MethodGet, token, "/v1/admin/users", "")
	if status != 403 || env.AdminUsers != nil {
		t.Fatalf("missing admin permission must be 403, got %d %#v", status, env)
	}
	ok, status := sessionDo(t, ts, http.MethodGet, token, "/v1/admin/system", "")
	if status != 200 || ok.AdminSystem == nil {
		t.Fatalf("permitted system read must work, got %d", status)
	}
}

func TestAdminRateLimitCreatesAbuseEvent(t *testing.T) {
	ts, store, _ := setup(t)
	userID, token := registerUser(t, ts, "ops-abuse@example.com", "correct-horse")
	store.AddOperator(userID, "ops-abuse@example.com", admin.RoleOperator, nil)
	for i := 0; i < 41; i++ {
		res, err := http.Post(ts.URL+"/v1/diagnoses", "application/json", strings.NewReader(`{"target":"http://127.0.0.1/"}`))
		if err != nil {
			t.Fatal(err)
		}
		res.Body.Close()
		if i == 40 && res.StatusCode != 429 {
			t.Fatalf("41st diagnose must be 429, got %d", res.StatusCode)
		}
	}
	env, status := sessionDo(t, ts, http.MethodGet, token, "/v1/admin/abuse", "")
	if status != 200 {
		t.Fatalf("abuse %d %#v", status, env.Error)
	}
	foundLimit, foundSSRF := false, false
	for _, event := range env.AbuseEvents {
		if event.Kind == "rate_limit" {
			foundLimit = true
		}
		if event.Kind == "ssrf" {
			foundSSRF = true
		}
	}
	if !foundLimit {
		t.Fatal("rate-limit 429 must store an abuse event")
	}
	if !foundSSRF {
		t.Fatal("ssrf_blocked diagnose must store an abuse event")
	}
}

func TestAdminOverrideIsAudited(t *testing.T) {
	ts, store, _ := setup(t)
	userID, token := registerUser(t, ts, "ops-inc@example.com", "correct-horse")
	store.AddOperator(userID, "ops-inc@example.com", admin.RoleOperator, nil)
	store.ReplaceIncidents([]contract.Incident{{
		ID:               "22222222-2222-4222-8222-222222222222",
		Title:            "Elevated connectivity failures observed",
		Severity:         "high",
		Status:           "investigating",
		Scope:            "youtube",
		StartedAt:        "2026-09-01T04:00:00Z",
		LastUpdatedAt:    "2026-09-01T05:00:00Z",
		AffectedServices: []string{"youtube"},
		SampleCount:      1,
	}})
	blocked, status := sessionDo(t, ts, http.MethodPost, token, "/v1/admin/incidents/22222222-2222-4222-8222-222222222222/resolve", `{"reason":"looks fine"}`)
	if status != 400 {
		t.Fatalf("resolve without evidence must be 400, got %d %#v", status, blocked.Error)
	}
	over, status := sessionDo(t, ts, http.MethodPost, token, "/v1/admin/incidents/22222222-2222-4222-8222-222222222222/override", `{"classification":"isp","reason":"manual review of stored evidence"}`)
	if status != 200 || over.AdminIncident == nil || over.AdminIncident.Override == nil {
		t.Fatalf("override %d %#v", status, over.Error)
	}
	if over.AdminIncident.Override.Classification != "isp" {
		t.Fatal("override classification must be stored")
	}
	audit, status := sessionDo(t, ts, http.MethodGet, token, "/v1/admin/audit", "")
	if status != 200 {
		t.Fatalf("audit %d", status)
	}
	found := false
	for _, event := range audit.AdminAudit {
		if event.Action == "incident.override" && event.ActorID == userID && event.Resource == "incident:22222222-2222-4222-8222-222222222222" && event.Result == "overridden" {
			found = true
		}
	}
	if !found {
		t.Fatalf("override must write actor, action, resource, result, got %#v", audit.AdminAudit)
	}
}

func TestAdminRulesFalsePositivesStayZeroWithoutLabels(t *testing.T) {
	ts, store, _ := setup(t)
	userID, token := registerUser(t, ts, "ops-rules@example.com", "correct-horse")
	store.AddOperator(userID, "ops-rules@example.com", admin.RoleOperator, nil)
	env, status := sessionDo(t, ts, http.MethodGet, token, "/v1/admin/rules", "")
	if status != 200 || env.RuleOutcomes == nil || len(env.AdminRules) == 0 {
		t.Fatalf("rules %d %#v", status, env.Error)
	}
	if env.RuleOutcomes.FalsePositives != 0 || env.RuleOutcomes.FalseNegatives != 0 {
		t.Fatal("unlabeled diagnoses must not invent false positives or false negatives")
	}
}

func TestAdminFeatureFlagsTargetingAndConfig(t *testing.T) {
	ts, store, _ := setup(t)
	userID, token := registerUser(t, ts, "ops-flags@example.com", "correct-horse")
	store.AddOperator(userID, "ops-flags@example.com", admin.RoleOperator, nil)
	created, status := sessionDo(t, ts, http.MethodPost, token, "/v1/admin/flags", `{"name":"ops.beta","environment":"test","enabled":true,"percentage":0,"userIds":["`+userID+`"],"orgIds":[]}`)
	if status != 201 || created.FeatureFlag == nil {
		t.Fatalf("create flag %d %#v", status, created.Error)
	}
	hit, status := sessionDo(t, ts, http.MethodGet, token, "/v1/admin/flags/"+created.FeatureFlag.ID+"?environment=test&userId="+userID, "")
	if status != 200 || hit.FeatureFlag == nil || hit.FeatureFlag.TargetMatch == nil || !*hit.FeatureFlag.TargetMatch {
		t.Fatalf("user targeting must match, got %d %#v", status, hit.FeatureFlag)
	}
	miss, status := sessionDo(t, ts, http.MethodGet, token, "/v1/admin/flags/"+created.FeatureFlag.ID+"?environment=test&userId=other-user", "")
	if status != 200 || miss.FeatureFlag == nil || miss.FeatureFlag.TargetMatch == nil || *miss.FeatureFlag.TargetMatch {
		t.Fatal("unknown user must not match a user-targeted 0% flag")
	}
	orgFlag, status := sessionDo(t, ts, http.MethodPost, token, "/v1/admin/flags", `{"name":"ops.org","environment":"test","enabled":true,"percentage":100,"userIds":[],"orgIds":["org-1"]}`)
	if status != 201 {
		t.Fatalf("org flag %d %#v", status, orgFlag.Error)
	}
	orgHit, status := sessionDo(t, ts, http.MethodGet, token, "/v1/admin/flags/"+orgFlag.FeatureFlag.ID+"?environment=test&orgId=org-1", "")
	if status != 200 || orgHit.FeatureFlag.TargetMatch == nil || !*orgHit.FeatureFlag.TargetMatch {
		t.Fatal("organization targeting must match")
	}
	cfg, status := sessionDo(t, ts, http.MethodPut, token, "/v1/admin/config", `{"key":"diagnose.timeoutSeconds","value":"25"}`)
	if status != 200 {
		t.Fatalf("config put %d %#v", status, cfg.Error)
	}
	found := false
	for _, item := range cfg.RemoteConfig {
		if item.Key == "diagnose.timeoutSeconds" && item.Value == "25" {
			found = true
		}
		if strings.Contains(strings.ToLower(item.Value), "secret") || strings.Contains(strings.ToLower(item.Key), "password") {
			t.Fatal("remote config must not store secrets")
		}
	}
	if !found {
		t.Fatal("updated timeout must be visible")
	}
	denied, status := sessionDo(t, ts, http.MethodPut, token, "/v1/admin/config", `{"key":"diagnose.timeoutSeconds","value":"super-secret-token"}`)
	if status != 400 {
		t.Fatalf("secret-like config values must be rejected, got %d", status)
	}
	_ = denied
}
