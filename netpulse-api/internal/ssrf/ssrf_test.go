package ssrf

import (
	"net"
	"net/http"
	"net/netip"
	"net/url"
	"testing"
)

func TestCheckHostBlocksLocalAndMetadata(t *testing.T) {
	policy := DefaultPolicy()
	blocked := []string{
		"localhost",
		"foo.localhost",
		"metadata.google.internal",
		"metadata",
		"printer.local",
		"db.internal",
		"127.0.0.1",
		"10.0.0.8",
		"192.168.1.1",
		"169.254.169.254",
		"168.63.129.16",
		"100.64.0.1",
		"255.255.255.255",
		"240.0.0.1",
		"127.1",
		"2130706433",
		"0x7f000001",
		"::1",
		"[::1]",
		"::ffff:127.0.0.1",
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

func TestCheckIPBlocksCloudMetadataAndIPv6Specials(t *testing.T) {
	policy := DefaultPolicy()
	blocked := []string{
		"168.63.129.16",
		"169.254.169.254",
		"fd00:ec2::254",
		"2001:db8::1",
		"2002:c0a8:0101::",
		"64:ff9b::a00:1",
		"fe80::1",
		"0100::1",
	}
	for _, raw := range blocked {
		ip := netip.MustParseAddr(raw)
		if err := policy.CheckIP(ip); err == nil {
			t.Fatalf("expected block for %s", raw)
		}
	}
	if err := policy.CheckIP(netip.MustParseAddr("1.1.1.1")); err != nil {
		t.Fatalf("public v4 should pass: %v", err)
	}
}

func TestCheckURLBlocksCredentialsAndHTTPWhenRequired(t *testing.T) {
	policy := DefaultPolicy()
	private, _ := url.Parse("https://127.0.0.1/hook")
	if err := CheckURL(policy, private, true); err == nil {
		t.Fatal("loopback https must be blocked")
	}
	httpURL, _ := url.Parse("http://example.com/hook")
	if err := CheckURL(policy, httpURL, true); err == nil {
		t.Fatal("http must be blocked when https-only")
	}
	ok, _ := url.Parse("https://example.com/hook")
	if err := CheckURL(policy, ok, true); err != nil {
		t.Fatalf("public https should pass: %v", err)
	}
}

func TestHTTPSClientRejectsPrivateRedirect(t *testing.T) {
	client := NewHTTPSClient()
	req, err := http.NewRequest(http.MethodPost, "https://127.0.0.1/hook", nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := client.CheckRedirect(req, []*http.Request{req}); err == nil {
		t.Fatal("redirect to loopback must fail")
	}
	private, err := http.NewRequest(http.MethodPost, "https://10.1.2.3/hook", nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := client.CheckRedirect(private, []*http.Request{private}); err == nil {
		t.Fatal("redirect to RFC1918 must fail")
	}
	azure, err := http.NewRequest(http.MethodPost, "https://168.63.129.16/", nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := client.CheckRedirect(azure, []*http.Request{azure}); err == nil {
		t.Fatal("redirect to Azure IMDS must fail")
	}
	plain, err := http.NewRequest(http.MethodPost, "http://example.com/hook", nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := client.CheckRedirect(plain, []*http.Request{plain}); err == nil {
		t.Fatal("http redirect must fail for webhook client")
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
