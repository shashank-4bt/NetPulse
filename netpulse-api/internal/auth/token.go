package auth

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"net"
	"strings"

	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/id"
)

func NewSecret() string {
	return id.RandomHex(32)
}

func HashSecret(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}

func SecretsEqual(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}

func CoarseIP(raw string) string {
	host, _, err := net.SplitHostPort(raw)
	if err != nil {
		host = raw
	}
	host = strings.TrimSpace(host)
	ip := net.ParseIP(host)
	if ip == nil {
		return ""
	}
	if v4 := ip.To4(); v4 != nil {
		return net.IPv4(v4[0], v4[1], v4[2], 0).String()
	}
	return host
}

func SessionLabel(userAgent string) string {
	ua := strings.TrimSpace(userAgent)
	if ua == "" {
		return "Unknown device"
	}
	if len(ua) > 80 {
		return ua[:80]
	}
	return ua
}
