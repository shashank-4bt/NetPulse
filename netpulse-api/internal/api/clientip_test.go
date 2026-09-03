package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestParseClientIPIgnoresSpoofedXForwardedFor(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/v1/auth/login", nil)
	req.RemoteAddr = "203.0.113.9:55555"
	req.Header.Set("X-Forwarded-For", "8.8.8.8")
	if got := ParseClientIP(req, false); got != "203.0.113.9" {
		t.Fatalf("untrusted XFF must not replace peer IP, got %q", got)
	}
}

func TestParseClientIPHonorsTrustedProxy(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/v1/auth/login", nil)
	req.RemoteAddr = "192.0.2.10:443"
	req.Header.Set("X-Forwarded-For", "198.51.100.20, 192.0.2.10")
	if got := ParseClientIP(req, true); got != "198.51.100.20" {
		t.Fatalf("trusted proxy should use first XFF hop, got %q", got)
	}
}

func TestParseClientIPLoopbackAllowsBFFHeader(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/v1/diagnoses", nil)
	req.RemoteAddr = "127.0.0.1:4321"
	req.Header.Set("X-NetPulse-Client-IP", "203.0.113.44")
	req.Header.Set("X-Forwarded-For", "8.8.8.8")
	if got := ParseClientIP(req, false); got != "203.0.113.44" {
		t.Fatalf("loopback BFF header should win, got %q", got)
	}
}

func TestParseClientIPRejectsSpoofedBFFHeaderFromInternet(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/v1/auth/login", nil)
	req.RemoteAddr = "198.51.100.9:443"
	req.Header.Set("X-NetPulse-Client-IP", "8.8.8.8")
	if got := ParseClientIP(req, false); got != "198.51.100.9" {
		t.Fatalf("non-loopback peer must ignore BFF header, got %q", got)
	}
}

func TestClientIPStripsPort(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/v1/health", nil)
	req.RemoteAddr = "[::1]:9999"
	if got := ClientIP(req); got != "::1" {
		t.Fatalf("expected ::1, got %q", got)
	}
}
