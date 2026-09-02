package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

func SignWebhook(secret, timestamp, eventID, body string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(timestamp))
	mac.Write([]byte("."))
	mac.Write([]byte(eventID))
	mac.Write([]byte("."))
	mac.Write([]byte(body))
	return "sha256=" + hex.EncodeToString(mac.Sum(nil))
}

func VerifyWebhook(secret, timestamp, eventID, body, header string) bool {
	expected := SignWebhook(secret, timestamp, eventID, body)
	return SecretsEqual(expected, strings.TrimSpace(header))
}
