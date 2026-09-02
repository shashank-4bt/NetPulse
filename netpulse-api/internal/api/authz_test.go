package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/contract"
)

func registerUser(t *testing.T, ts *httptest.Server, email, password string) (userID, sessionToken string) {
	t.Helper()
	res, err := http.Post(ts.URL+"/v1/auth/register", "application/json", bytes.NewBufferString(
		`{"email":"`+email+`","password":"`+password+`","displayName":"Tester"}`,
	))
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 201 {
		t.Fatalf("register %s: %d", email, res.StatusCode)
	}
	var env contract.Envelope
	if err := json.NewDecoder(res.Body).Decode(&env); err != nil {
		t.Fatal(err)
	}
	if env.User == nil || env.SessionToken == "" || env.Session == nil {
		t.Fatal("register must return user, public session, and session token")
	}
	if env.Session.ID == env.SessionToken {
		t.Fatal("public session id must not be the cookie secret")
	}
	return env.User.ID, env.SessionToken
}

func sessionGet(t *testing.T, ts *httptest.Server, token, path string) *http.Response {
	t.Helper()
	req, err := http.NewRequest(http.MethodGet, ts.URL+path, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Session "+token)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	return res
}

func sessionJSON(t *testing.T, ts *httptest.Server, method, token, path, body string) contract.Envelope {
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
	req.Header.Set("Authorization", "Session "+token)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	var env contract.Envelope
	if err := json.NewDecoder(res.Body).Decode(&env); err != nil {
		t.Fatal(err)
	}
	if res.StatusCode >= 400 && env.Error == nil {
		t.Fatalf("%s %s: %d without error", method, path, res.StatusCode)
	}
	return env
}

func TestAuthBoundariesPreventCrossUserAccess(t *testing.T) {
	ts, _, _ := setup(t)
	userA, tokenA := registerUser(t, ts, "a@example.com", "correct-horse")
	_, tokenB := registerUser(t, ts, "b@example.com", "correct-horse")

	req, err := http.NewRequest(http.MethodPost, ts.URL+"/v1/diagnoses", bytes.NewBufferString(`{"target":"example.com"}`))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Session "+tokenA)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	var created contract.Envelope
	if err := json.NewDecoder(res.Body).Decode(&created); err != nil {
		t.Fatal(err)
	}
	res.Body.Close()
	if created.Diagnosis == nil {
		t.Fatal("expected owned diagnosis")
	}
	diagID := created.Diagnosis.ID

	res = sessionGet(t, ts, tokenB, "/v1/diagnoses/"+diagID)
	defer res.Body.Close()
	if res.StatusCode != 404 {
		t.Fatalf("cross-user diagnosis must be 404, got %d", res.StatusCode)
	}
	var hidden contract.Envelope
	if err := json.NewDecoder(res.Body).Decode(&hidden); err != nil {
		t.Fatal(err)
	}
	if hidden.Diagnosis != nil {
		t.Fatal("must not leak another user's diagnosis body")
	}

	env := sessionJSON(t, ts, http.MethodDelete, tokenB, "/v1/me/reports/"+diagID, "")
	if env.Error == nil || env.Error.Code != "not_found" {
		t.Fatal("cross-user delete must be not_found")
	}

	env = sessionJSON(t, ts, http.MethodGet, tokenB, "/v1/users/"+userA+"/billing", "")
	if env.Error == nil || env.Error.Code != "not_found" || env.Billing != nil {
		t.Fatal("cross-user billing must be 404 without invoices")
	}

	env = sessionJSON(t, ts, http.MethodGet, tokenB, "/v1/organizations/"+userA, "")
	if env.Error == nil || env.Error.Code != "not_found" {
		t.Fatal("organizations must stay not found")
	}

	mine := sessionJSON(t, ts, http.MethodGet, tokenA, "/v1/me/diagnoses", "")
	if len(mine.Diagnoses) != 1 || mine.Diagnoses[0].ID != diagID {
		t.Fatal("owner must see their diagnosis")
	}
	theirs := sessionJSON(t, ts, http.MethodGet, tokenB, "/v1/me/diagnoses", "")
	if len(theirs.Diagnoses) != 0 {
		t.Fatal("other user must not see foreign diagnoses")
	}

	sessions := sessionJSON(t, ts, http.MethodGet, tokenA, "/v1/auth/sessions", "")
	if len(sessions.Sessions) == 0 {
		t.Fatal("owner must see their session")
	}
	raw, _ := json.Marshal(sessions)
	if bytes.Contains(raw, []byte(tokenA)) {
		t.Fatal("session list must not include the cookie secret")
	}
	if bytes.Contains(raw, []byte(`"userId"`)) || bytes.Contains(raw, []byte(`"tokenHash"`)) {
		t.Fatal("session list must not expose internal owner fields")
	}

	foreignSession := sessions.Sessions[0].ID
	revoked := sessionJSON(t, ts, http.MethodPost, tokenB, "/v1/auth/sessions/"+foreignSession+"/revoke", "")
	if revoked.Error == nil || revoked.Error.Code != "not_found" {
		t.Fatal("cannot revoke another user's session")
	}
}

func TestLoginAndResetDoNotEnumerateAccounts(t *testing.T) {
	ts, _, _ := setup(t)
	_, _ = registerUser(t, ts, "known@example.com", "correct-horse")

	res, err := http.Post(ts.URL+"/v1/auth/login", "application/json", bytes.NewBufferString(
		`{"email":"known@example.com","password":"wrong-password"}`,
	))
	if err != nil {
		t.Fatal(err)
	}
	var failed contract.Envelope
	if err := json.NewDecoder(res.Body).Decode(&failed); err != nil {
		t.Fatal(err)
	}
	res.Body.Close()
	if res.StatusCode != 401 || failed.Error == nil || failed.Error.Message != "Email or password is incorrect." {
		t.Fatalf("generic login failure, got %d %v", res.StatusCode, failed.Error)
	}

	unknown, err := http.Post(ts.URL+"/v1/auth/login", "application/json", bytes.NewBufferString(
		`{"email":"missing@example.com","password":"wrong-password"}`,
	))
	if err != nil {
		t.Fatal(err)
	}
	var missing contract.Envelope
	if err := json.NewDecoder(unknown.Body).Decode(&missing); err != nil {
		t.Fatal(err)
	}
	unknown.Body.Close()
	if missing.Error == nil || missing.Error.Message != failed.Error.Message {
		t.Fatal("unknown email must use the same login error")
	}

	knownReset, err := http.Post(ts.URL+"/v1/auth/forgot-password", "application/json", bytes.NewBufferString(`{"email":"known@example.com"}`))
	if err != nil {
		t.Fatal(err)
	}
	var knownEnv contract.Envelope
	if err := json.NewDecoder(knownReset.Body).Decode(&knownEnv); err != nil {
		t.Fatal(err)
	}
	knownReset.Body.Close()
	unknownReset, err := http.Post(ts.URL+"/v1/auth/forgot-password", "application/json", bytes.NewBufferString(`{"email":"missing@example.com"}`))
	if err != nil {
		t.Fatal(err)
	}
	var unknownEnv contract.Envelope
	if err := json.NewDecoder(unknownReset.Body).Decode(&unknownEnv); err != nil {
		t.Fatal(err)
	}
	unknownReset.Body.Close()
	if knownEnv.Auth == nil || unknownEnv.Auth == nil || knownEnv.Auth.EmailReason != unknownEnv.Auth.EmailReason {
		t.Fatal("forgot-password must not enumerate emails")
	}
}

func TestShareTokenIsTheOnlyForeignReadPath(t *testing.T) {
	ts, _, _ := setup(t)
	_, tokenA := registerUser(t, ts, "share-a@example.com", "correct-horse")
	_, tokenB := registerUser(t, ts, "share-b@example.com", "correct-horse")
	req, _ := http.NewRequest(http.MethodPost, ts.URL+"/v1/diagnoses", bytes.NewBufferString(`{"target":"example.com"}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Session "+tokenA)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	var created contract.Envelope
	if err := json.NewDecoder(res.Body).Decode(&created); err != nil {
		t.Fatal(err)
	}
	res.Body.Close()

	shared := sessionJSON(t, ts, http.MethodPost, tokenA, "/v1/me/reports/"+created.Diagnosis.ID+"/share", "")
	if shared.Share == nil || shared.Share.Token == "" {
		t.Fatal("owner can create a share token")
	}
	res = sessionGet(t, ts, tokenB, "/v1/diagnoses/"+created.Diagnosis.ID+"?share="+shared.Share.Token)
	defer res.Body.Close()
	if res.StatusCode != 200 {
		t.Fatalf("valid share should read, got %d", res.StatusCode)
	}

	oauth, err := http.Post(ts.URL+"/v1/auth/oauth/google", "application/json", strings.NewReader(`{}`))
	if err != nil {
		t.Fatal(err)
	}
	if oauth.StatusCode != 501 {
		t.Fatalf("oauth must be unavailable, got %d", oauth.StatusCode)
	}
	oauth.Body.Close()
}

func TestAnonymousDiagnosisRemainsUUIDAccessible(t *testing.T) {
	ts, _, _ := setup(t)
	res, err := http.Post(ts.URL+"/v1/diagnoses", "application/json", bytes.NewBufferString(`{"target":"example.com"}`))
	if err != nil {
		t.Fatal(err)
	}
	var created contract.Envelope
	if err := json.NewDecoder(res.Body).Decode(&created); err != nil {
		t.Fatal(err)
	}
	res.Body.Close()
	res, err = http.Get(ts.URL + "/v1/diagnoses/" + created.Diagnosis.ID)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		t.Fatalf("anonymous diagnosis should stay readable by id, got %d", res.StatusCode)
	}
}

func TestDashboardDoesNotInventHealthOrBilling(t *testing.T) {
	ts, _, _ := setup(t)
	_, token := registerUser(t, ts, "dash@example.com", "correct-horse")
	dash := sessionJSON(t, ts, http.MethodGet, token, "/v1/me/dashboard", "")
	if dash.Dashboard == nil {
		t.Fatal("dashboard required")
	}
	if strings.Contains(dash.Dashboard.InternetHealth, "%") || strings.Contains(dash.Dashboard.InternetHealth, "87") {
		t.Fatal("must not invent a health percentage")
	}
	billing := sessionJSON(t, ts, http.MethodGet, token, "/v1/me/billing", "")
	if billing.Billing == nil || billing.Billing.HasAccount || len(billing.Billing.Invoices) != 0 {
		t.Fatal("billing must stay empty")
	}
	alerts := sessionJSON(t, ts, http.MethodGet, token, "/v1/me/alerts", "")
	if alerts.Alerts == nil || alerts.Alerts.DeliveredCount != 0 {
		t.Fatal("alerts must not invent deliveries")
	}
}
