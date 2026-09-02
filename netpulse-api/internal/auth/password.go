package auth

import (
	"strings"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

const (
	MinPasswordLen = 10
	MaxPasswordLen = 128
	bcryptCost     = 12
)

func HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func CheckPassword(hash, password string) bool {
	if hash == "" {
		return false
	}
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

func ValidatePassword(password, email string) string {
	if len(password) < MinPasswordLen {
		return "Password must be at least 10 characters."
	}
	if len(password) > MaxPasswordLen {
		return "Password is too long."
	}
	if email != "" && strings.EqualFold(strings.TrimSpace(password), strings.TrimSpace(email)) {
		return "Password must not be the email address."
	}
	return ""
}

func NormalizeEmail(raw string) string {
	return strings.ToLower(strings.TrimSpace(raw))
}

func ValidateEmail(raw string) string {
	email := NormalizeEmail(raw)
	if email == "" || len(email) > 254 || !strings.Contains(email, "@") {
		return "Enter a valid email address."
	}
	at := strings.LastIndex(email, "@")
	local, domain := email[:at], email[at+1:]
	if local == "" || domain == "" || !strings.Contains(domain, ".") {
		return "Enter a valid email address."
	}
	for _, r := range email {
		if unicode.IsSpace(r) {
			return "Enter a valid email address."
		}
	}
	return ""
}
