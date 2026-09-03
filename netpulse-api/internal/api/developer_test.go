package api_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/auth"
	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/contract"
	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/id"
)

func TestDeveloperTenantIsolationAndKeys(t *testing.T) {
	ts, store, _ := setup(t)
	_, tokenA := registerUser(t, ts, "dev-a@example.com", "correct-horse")
	_, tokenB := registerUser(t, ts, "dev-b@example.com", "correct-horse")

	created, status := sessionDo(t, ts, http.MethodPost, tokenA, "/v1/dev/monitors", `{
		"name":"Example HTTP","target":"example.com","type":"http","regions":["us-east"],"frequencySeconds":300,"timeoutSeconds":10
	}`)
	if status != 201 || created.Monitor == nil {
		t.Fatalf("create monitor: %d %#v", status, created.Error)
	}
	monitorID := created.Monitor.ID

	foreign, status := sessionDo(t, ts, http.MethodGet, tokenB, "/v1/dev/monitors/"+monitorID, "")
	if status != 404 || foreign.Monitor != nil || foreign.Error == nil || foreign.Error.Code != "not_found" {
		t.Fatalf("cross-tenant monitor must be 404, got %d %#v", status, foreign)
	}
	patched, status := sessionDo(t, ts, http.MethodPatch, tokenB, "/v1/dev/monitors/"+monitorID, `{"name":"stolen"}`)
	if status != 404 || patched.Monitor != nil {
		t.Fatalf("cross-tenant patch must be 404, got %d", status)
	}
	deleted, status := sessionDo(t, ts, http.MethodDelete, tokenB, "/v1/dev/monitors/"+monitorID, "")
	if status != 404 || deleted.Error == nil || deleted.Error.Code != "not_found" {
		t.Fatalf("cross-tenant delete must be 404, got %d", status)
	}

	wsA, status := sessionDo(t, ts, http.MethodGet, tokenA, "/v1/dev/workspace", "")
	if status != 200 || wsA.Workspace == nil {
		t.Fatal("owner workspace required")
	}
	stolenWS, status := sessionDo(t, ts, http.MethodGet, tokenB, "/v1/dev/workspaces/"+wsA.Workspace.ID, "")
	if status != 404 || stolenWS.Workspace != nil {
		t.Fatalf("foreign workspace must be 404, got %d", status)
	}

	keyEnv, status := sessionDo(t, ts, http.MethodPost, tokenA, "/v1/dev/keys", `{"name":"ci","scopes":["monitors:read","usage:read","sla:read","dashboard:read"],"rateLimitPerMin":30}`)
	if status != 201 || keyEnv.APIKey == nil || keyEnv.KeySecret == "" || !strings.HasPrefix(keyEnv.KeySecret, "npk_") {
		t.Fatalf("create key: %d %#v", status, keyEnv)
	}
	secret := keyEnv.KeySecret
	keyID := keyEnv.APIKey.ID
	listed, status := sessionDo(t, ts, http.MethodGet, tokenA, "/v1/dev/keys", "")
	if status != 200 {
		t.Fatalf("list keys %d", status)
	}
	raw, _ := json.Marshal(listed)
	if bytes.Contains(raw, []byte(secret)) || bytes.Contains(raw, []byte(`"keySecret"`)) {
		t.Fatal("key list must not include the raw secret")
	}

	rotateB, status := sessionDo(t, ts, http.MethodPost, tokenB, "/v1/dev/keys/"+keyID+"/rotate", "")
	if status != 404 || rotateB.KeySecret != "" {
		t.Fatalf("foreign rotate must be 404, got %d", status)
	}

	readA, status := bearerDo(t, ts, http.MethodGet, secret, "/v1/dev/monitors/"+monitorID, "")
	if status != 200 || readA.Monitor == nil || readA.Monitor.ID != monitorID {
		t.Fatalf("own key should read own monitor, got %d", status)
	}
	_, status = bearerDo(t, ts, http.MethodPost, secret, "/v1/dev/monitors", `{
		"name":"blocked","target":"example.com","type":"dns"
	}`)
	if status != 403 {
		t.Fatalf("key without write scope must be 403, got %d", status)
	}
	ownWS, status := bearerDo(t, ts, http.MethodGet, secret, "/v1/dev/workspaces/"+wsA.Workspace.ID, "")
	if status != 200 || ownWS.Workspace == nil {
		t.Fatalf("own workspace via key: %d", status)
	}
	wsB, _ := sessionDo(t, ts, http.MethodGet, tokenB, "/v1/dev/workspace", "")
	stolen, status := bearerDo(t, ts, http.MethodGet, secret, "/v1/dev/workspaces/"+wsB.Workspace.ID, "")
	if status != 404 || stolen.Workspace != nil {
		t.Fatalf("key must not read another workspace, got %d", status)
	}

	_, status = bearerDo(t, ts, http.MethodPost, secret, "/v1/dev/keys", `{"name":"nope"}`)
	if status != 403 {
		t.Fatalf("keys cannot mint keys, got %d", status)
	}

	revoked, status := sessionDo(t, ts, http.MethodPost, tokenA, "/v1/dev/keys/"+keyID+"/revoke", "")
	if status != 200 || revoked.APIKey == nil || !revoked.APIKey.Revoked {
		t.Fatalf("revoke: %d", status)
	}
	_, status = bearerDo(t, ts, http.MethodGet, secret, "/v1/dev/monitors", "")
	if status != 401 {
		t.Fatalf("revoked key must be 401, got %d", status)
	}

	usageB, status := sessionDo(t, ts, http.MethodGet, tokenB, "/v1/dev/usage", "")
	if status != 200 || usageB.Usage == nil || usageB.Usage.Monitors != 0 {
		t.Fatalf("other tenant usage must stay empty, got %#v", usageB.Usage)
	}

	_ = store
}

func TestDeveloperWebhookSSRFAndEmptySLA(t *testing.T) {
	ts, store, _ := setup(t)
	_, token := registerUser(t, ts, "hooks@example.com", "correct-horse")

	blocked, status := sessionDo(t, ts, http.MethodPost, token, "/v1/dev/webhooks", `{
		"url":"https://127.0.0.1/hook","events":["monitor.down"]
	}`)
	if status != 403 || blocked.Webhook != nil || blocked.WebhookSecret != "" {
		t.Fatalf("loopback webhook must be blocked, got %d %#v", status, blocked.Error)
	}
	private, status := sessionDo(t, ts, http.MethodPost, token, "/v1/dev/webhooks", `{
		"url":"https://169.254.169.254/latest","events":["monitor.down"]
	}`)
	if status != 403 || private.Webhook != nil {
		t.Fatalf("metadata webhook must be blocked, got %d", status)
	}
	azure, status := sessionDo(t, ts, http.MethodPost, token, "/v1/dev/webhooks", `{
		"url":"https://168.63.129.16/metadata","events":["monitor.down"]
	}`)
	if status != 403 || azure.Webhook != nil {
		t.Fatalf("Azure IMDS webhook must be blocked, got %d", status)
	}
	httpHook, status := sessionDo(t, ts, http.MethodPost, token, "/v1/dev/webhooks", `{
		"url":"http://example.com/hook","events":["monitor.down"]
	}`)
	if status != 400 {
		t.Fatalf("http webhook must be rejected, got %d", status)
	}
	_ = httpHook

	okHook, status := sessionDo(t, ts, http.MethodPost, token, "/v1/dev/webhooks", `{
		"url":"https://example.com/hooks","events":["monitor.down","incident.created"]
	}`)
	if status != 201 || okHook.Webhook == nil || okHook.WebhookSecret == "" || !strings.HasPrefix(okHook.WebhookSecret, "nwh_") {
		t.Fatalf("public https webhook: %d %#v", status, okHook.Error)
	}
	listed, status := sessionDo(t, ts, http.MethodGet, token, "/v1/dev/webhooks", "")
	if status != 200 {
		t.Fatalf("list webhooks %d", status)
	}
	raw, _ := json.Marshal(listed)
	if bytes.Contains(raw, []byte(okHook.WebhookSecret)) || bytes.Contains(raw, []byte(`"webhookSecret"`)) {
		t.Fatal("webhook list must not include the signing secret")
	}

	sla, status := sessionDo(t, ts, http.MethodGet, token, "/v1/dev/sla", "")
	if status != 200 || sla.SLA == nil {
		t.Fatal("empty sla required")
	}
	if sla.SLA.Availability.Measured || sla.SLA.Latency.P50 != nil || sla.SLA.Latency.P95 != nil || sla.SLA.Latency.P99 != nil {
		t.Fatal("empty sla must not invent percentiles")
	}
	if strings.Contains(sla.SLA.Summary, "99.9") || strings.Contains(sla.SLA.Summary, "87") {
		t.Fatal("empty sla must not invent availability copy")
	}
	dash, status := sessionDo(t, ts, http.MethodGet, token, "/v1/dev/dashboard", "")
	if status != 200 || dash.DevDashboard == nil || dash.DevDashboard.Availability.Measured || dash.DevDashboard.Latency.P95 != nil {
		t.Fatal("empty dashboard must stay unmeasured")
	}

	monitor, status := sessionDo(t, ts, http.MethodPost, token, "/v1/dev/monitors", `{
		"name":"DNS","target":"example.com","type":"dns","regions":["eu-west"]
	}`)
	if status != 201 {
		t.Fatalf("monitor %d", status)
	}
	ws, _ := sessionDo(t, ts, http.MethodGet, token, "/v1/dev/workspace", "")
	lat := 12
	_ = store.AddCheck(context.Background(), contract.MonitorCheck{
		ID: id.New(), MonitorID: monitor.Monitor.ID, WorkspaceID: ws.Workspace.ID,
		Region: "eu-west", OK: true, LatencyMs: &lat, At: "2026-09-02T00:00:00Z",
		Summary: "test sample",
	})
	sla2, _ := sessionDo(t, ts, http.MethodGet, token, "/v1/dev/sla", "")
	if sla2.SLA.Latency.P95 != nil || sla2.SLA.Latency.SampleCount != 1 {
		t.Fatalf("one sample must not produce p95, got %#v", sla2.SLA.Latency)
	}

	body := `{"event":"monitor.down"}`
	sig := auth.SignWebhook(okHook.WebhookSecret, "2026-09-02T00:00:00Z", "evt-1", body)
	if !auth.VerifyWebhook(okHook.WebhookSecret, "2026-09-02T00:00:00Z", "evt-1", body, sig) {
		t.Fatal("hmac round-trip")
	}
}

func sessionDo(t *testing.T, ts *httptest.Server, method, token, path, body string) (contract.Envelope, int) {
	t.Helper()
	return doAuth(t, ts, method, "Session "+token, "", path, body)
}

func bearerDo(t *testing.T, ts *httptest.Server, method, key, path, body string) (contract.Envelope, int) {
	t.Helper()
	return doAuth(t, ts, method, "Bearer "+key, "", path, body)
}

func doAuth(t *testing.T, ts *httptest.Server, method, authorization, extraHeader, path, body string) (contract.Envelope, int) {
	t.Helper()
	var reader *bytes.Buffer
	if body != "" {
		reader = bytes.NewBufferString(body)
	} else {
		reader = bytes.NewBuffer(nil)
	}
	req, err := http.NewRequest(method, ts.URL+path, reader)
	if err != nil {
		t.Fatal(err)
	}
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	if authorization != "" {
		req.Header.Set("Authorization", authorization)
	}
	if extraHeader != "" {
		req.Header.Set("X-NetPulse-Key", extraHeader)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	var env contract.Envelope
	if err := json.NewDecoder(res.Body).Decode(&env); err != nil {
		t.Fatal(err)
	}
	return env, res.StatusCode
}
