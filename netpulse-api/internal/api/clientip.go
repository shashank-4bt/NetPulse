package api

import (
	"net"
	"net/http"
	"net/netip"
	"strings"
)

func (s *Server) clientIP(r *http.Request) string {
	trustProxy := false
	if s != nil {
		trustProxy = s.Cfg.TrustProxy
	}
	return ParseClientIP(r, trustProxy)
}

// ParseClientIP returns the peer address used for rate limits and audit.
// X-Forwarded-For is ignored unless trustProxy is true and the first hop is a
// valid IP. Loopback peers may supply X-NetPulse-Client-IP (set by the Next.js
// BFF when NETPULSE_WEB_TRUST_PROXY is enabled).
func ParseClientIP(r *http.Request, trustProxy bool) string {
	remote := peerIP(r.RemoteAddr)
	if isLoopbackIP(remote) {
		if client := strings.TrimSpace(r.Header.Get("X-NetPulse-Client-IP")); client != "" {
			if ip, ok := parseStandaloneIP(client); ok {
				return ip
			}
		}
	}
	if trustProxy {
		if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
			first := strings.TrimSpace(strings.Split(forwarded, ",")[0])
			if ip, ok := parseStandaloneIP(first); ok {
				return ip
			}
		}
	}
	if remote != "" {
		return remote
	}
	return r.RemoteAddr
}

func ClientIP(r *http.Request) string {
	return ParseClientIP(r, false)
}

func peerIP(remoteAddr string) string {
	host, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		host = remoteAddr
	}
	host = strings.TrimSpace(host)
	if host == "" {
		return ""
	}
	if ip, err := netip.ParseAddr(host); err == nil {
		return ip.Unmap().String()
	}
	return host
}

func parseStandaloneIP(raw string) (string, bool) {
	value := strings.TrimSpace(raw)
	if value == "" || strings.Contains(value, ",") {
		return "", false
	}
	if ip, err := netip.ParseAddr(value); err == nil {
		return ip.Unmap().String(), true
	}
	return "", false
}

func isLoopbackIP(raw string) bool {
	ip, err := netip.ParseAddr(raw)
	if err != nil {
		return raw == "localhost"
	}
	return ip.Unmap().IsLoopback()
}
