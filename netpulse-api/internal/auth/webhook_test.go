package auth

import "testing"

func TestWebhookSignatureRoundTrip(t *testing.T) {
	secret := "nwh_test-secret"
	timestamp := "2026-09-02T04:30:00Z"
	eventID := "evt-1"
	body := `{"event":"monitor.down"}`
	sig := SignWebhook(secret, timestamp, eventID, body)
	if !stringsHasPrefix(sig, "sha256=") {
		t.Fatalf("signature prefix: %s", sig)
	}
	if !VerifyWebhook(secret, timestamp, eventID, body, sig) {
		t.Fatal("valid signature must verify")
	}
	if VerifyWebhook(secret, timestamp, eventID, `{"event":"other"}`, sig) {
		t.Fatal("tampered body must not verify")
	}
	if VerifyWebhook("other", timestamp, eventID, body, sig) {
		t.Fatal("wrong secret must not verify")
	}
	if VerifyWebhook(secret, timestamp, eventID, body, "") {
		t.Fatal("empty signature must not verify")
	}
}

func TestAPIKeySecretIsHashedAndPrefixed(t *testing.T) {
	raw, prefix, last4, hash := NewAPIKeySecret()
	if len(raw) < 12 || raw[:4] != "npk_" {
		t.Fatalf("raw key format: %s", raw)
	}
	if prefix != raw[:8] || last4 != raw[len(raw)-4:] {
		t.Fatal("display fields must come from the raw secret")
	}
	if hash == raw || hash != HashSecret(raw) {
		t.Fatal("only the hash should be stored")
	}
}

func stringsHasPrefix(value, prefix string) bool {
	return len(value) >= len(prefix) && value[:len(prefix)] == prefix
}
