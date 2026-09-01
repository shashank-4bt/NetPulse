package ssrf

import (
	"net"
	"net/netip"
	"testing"
)

func TestCheckHostBlocksLocalAndMetadata(t *testing.T) {
	policy := DefaultPolicy()
	blocked := []string{
		"localhost",
		"metadata.google.internal",
		"printer.local",
		"127.0.0.1",
		"10.0.0.8",
		"192.168.1.1",
		"169.254.169.254",
		"127.1",
		"2130706433",
		"::1",
	}
	for _, host := range blocked {
		if err := policy.CheckHost(host); err == nil {
			t.Fatalf("expected block for %s", host)
		}
	}
	if err := policy.CheckHost("example.com"); err != nil {
		t.Fatalf("public host should pass hostname check: %v", err)
	}
}

func TestCheckResolvedBlocksPrivateAnswers(t *testing.T) {
	policy := DefaultPolicy()
	err := policy.CheckResolved("example.com", []net.IP{net.ParseIP("10.1.2.3")})
	if err == nil {
		t.Fatal("resolved private IP must be blocked")
	}
	if err := policy.CheckResolved("example.com", []net.IP{net.ParseIP("8.8.8.8")}); err != nil {
		t.Fatalf("public resolution should pass: %v", err)
	}
}

func TestPublicIPsFiltersLoopback(t *testing.T) {
	filtered := PublicIPs([]net.IP{
		net.ParseIP("127.0.0.1"),
		net.ParseIP("1.1.1.1"),
	})
	if len(filtered) != 1 || filtered[0].String() != "1.1.1.1" {
		t.Fatalf("unexpected filter result: %v", filtered)
	}
}

func TestMappedIPv4Loopback(t *testing.T) {
	ip := netip.MustParseAddr("::ffff:127.0.0.1")
	if err := DefaultPolicy().CheckIP(ip); err == nil {
		t.Fatal("mapped loopback must be blocked")
	}
}
