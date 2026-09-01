package validation

import (
	"fmt"
	"net"
	"net/url"
	"strings"
	"unicode"

	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/ssrf"
)

const (
	maxInput    = 2048
	maxHostname = 253
)

type Target struct {
	Raw         string `json:"raw"`
	Hostname    string `json:"hostname"`
	Kind        string `json:"kind"`
	ServiceSlug string `json:"serviceSlug"`
}

type Result struct {
	Target Target
	Err    error
	Code   string
}

var knownServices = map[string]string{
	"google":        "google.com",
	"youtube":       "youtube.com",
	"cloudflare":    "cloudflare.com",
	"github":        "github.com",
	"microsoft-365": "office.com",
	"microsoft 365": "office.com",
	"slack":         "slack.com",
	"aws":           "aws.amazon.com",
	"zoom":          "zoom.us",
	"instagram":     "instagram.com",
}

func ParseTarget(raw string) Result {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return Result{Err: fmt.Errorf("enter a hostname, URL, or known service"), Code: "validation_error"}
	}
	if len(trimmed) > maxInput {
		return Result{Err: fmt.Errorf("input is too long"), Code: "validation_error"}
	}
	if strings.ContainsAny(trimmed, "<>'\"\\") {
		return Result{Err: fmt.Errorf("input contains unsafe characters"), Code: "validation_error"}
	}

	lowered := strings.ToLower(trimmed)
	if host, ok := resolveKnown(lowered); ok {
		target := Target{Raw: trimmed, Hostname: host, Kind: "known_service", ServiceSlug: knownSlug(lowered, host)}
		if err := ssrf.DefaultPolicy().CheckHost(target.Hostname); err != nil {
			return Result{Err: err, Code: "ssrf_blocked"}
		}
		return Result{Target: target}
	}

	for _, scheme := range []string{"javascript:", "data:", "file:", "ftp:", "ws:", "wss:"} {
		if strings.HasPrefix(lowered, scheme) {
			return Result{Err: fmt.Errorf("only http and https URLs are accepted"), Code: "validation_error"}
		}
	}

	hostname := trimmed
	kind := "domain"
	if strings.Contains(trimmed, "://") {
		parsed, err := url.Parse(trimmed)
		if err != nil {
			return Result{Err: fmt.Errorf("enter a valid hostname or URL"), Code: "validation_error"}
		}
		if parsed.Scheme != "http" && parsed.Scheme != "https" {
			return Result{Err: fmt.Errorf("only http and https URLs are accepted"), Code: "validation_error"}
		}
		if parsed.User != nil {
			return Result{Err: fmt.Errorf("URLs with credentials are not accepted"), Code: "validation_error"}
		}
		hostname = parsed.Hostname()
		kind = "url"
	}

	hostname = strings.Trim(strings.ToLower(hostname), "[]")
	hostname = strings.TrimSuffix(hostname, ".")
	if hostname == "" || strings.ContainsAny(hostname, " /") || len(hostname) > maxHostname {
		return Result{Err: fmt.Errorf("enter a valid hostname or URL"), Code: "validation_error"}
	}
	if !isHostnameOrIP(hostname) {
		return Result{Err: fmt.Errorf("enter a fully qualified hostname"), Code: "validation_error"}
	}
	if err := ssrf.DefaultPolicy().CheckHost(hostname); err != nil {
		return Result{Err: err, Code: "ssrf_blocked"}
	}
	slug := ""
	kindOut := kind
	if host, ok := resolveKnown(hostname); ok && host == hostname {
		kindOut = "known_service"
		slug = knownSlug(hostname, host)
	}
	return Result{Target: Target{Raw: trimmed, Hostname: hostname, Kind: kindOut, ServiceSlug: slug}}
}

func resolveKnown(value string) (string, bool) {
	if host, ok := knownServices[value]; ok {
		return host, true
	}
	for slug, host := range knownServices {
		if value == host || value == "www."+host || value == slug {
			return host, true
		}
	}
	return "", false
}

func knownSlug(value, host string) string {
	for slug, mapped := range knownServices {
		if value == slug || value == mapped || host == mapped {
			if !strings.Contains(slug, " ") {
				return slug
			}
		}
	}
	return ""
}

func isHostnameOrIP(hostname string) bool {
	if ip := net.ParseIP(hostname); ip != nil {
		return true
	}
	labels := strings.Split(hostname, ".")
	if len(labels) < 2 {
		return false
	}
	for _, label := range labels {
		if len(label) == 0 || len(label) > 63 || strings.HasPrefix(label, "-") || strings.HasSuffix(label, "-") {
			return false
		}
		for _, r := range label {
			if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '-' {
				return false
			}
		}
	}
	return true
}
