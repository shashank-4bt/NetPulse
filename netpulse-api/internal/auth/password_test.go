package auth

import "testing"

func TestPasswordHashRoundTrip(t *testing.T) {
	hash, err := HashPassword("correct-horse")
	if err != nil {
		t.Fatal(err)
	}
	if hash == "correct-horse" || hash == "" {
		t.Fatal("must store a hash, not the password")
	}
	if !CheckPassword(hash, "correct-horse") {
		t.Fatal("expected match")
	}
	if CheckPassword(hash, "wrong-password") {
		t.Fatal("wrong password must fail")
	}
}

func TestPasswordValidation(t *testing.T) {
	if ValidatePassword("short", "a@b.com") == "" {
		t.Fatal("short password")
	}
	if ValidatePassword("user@example.com", "user@example.com") == "" {
		t.Fatal("password must not equal email")
	}
	if ValidateEmail("not-an-email") == "" {
		t.Fatal("invalid email")
	}
	if NormalizeEmail("  A@B.COM ") != "a@b.com" {
		t.Fatal("email normalize")
	}
}
