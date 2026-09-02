package auth

import "testing"

func TestSecretHashIsNotReversible(t *testing.T) {
	raw := NewSecret()
	if len(raw) < 32 {
		t.Fatal("secret too short")
	}
	hashed := HashSecret(raw)
	if hashed == raw {
		t.Fatal("must not store the raw secret")
	}
	if !SecretsEqual(hashed, HashSecret(raw)) {
		t.Fatal("same secret must hash equal")
	}
}

func TestCoarseIPStripsHostIdentifier(t *testing.T) {
	if got := CoarseIP("203.0.113.44:51234"); got != "203.0.113.0" {
		t.Fatalf("got %s", got)
	}
	if CoarseIP("127.0.0.1:1") != "127.0.0.0" {
		t.Fatal("loopback should still be coarsened, not published as a precise socket")
	}
}
