package measurements

import (
	"context"
	"errors"
	"net"
	"net/url"
	"testing"
	"time"

	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/ssrf"
	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/validation"
)

type recordDialer struct {
	addrs []string
}

func (d *recordDialer) DialContext(_ context.Context, _, address string) (net.Conn, error) {
	d.addrs = append(d.addrs, address)
	return nil, errors.New("dial disabled in test")
}

type fakeResolver struct {
	ips []net.IP
	err error
}

func (f fakeResolver) LookupIP(context.Context, string) ([]net.IP, error) {
	return f.ips, f.err
}

type rejectDialer struct{}

func (rejectDialer) DialContext(context.Context, string, string) (net.Conn, error) {
	return nil, errors.New("dial disabled in test")
}

func TestRunBlocksPrivateResolution(t *testing.T) {
	runner := NewRunner()
	runner.Resolver = fakeResolver{ips: []net.IP{net.ParseIP("127.0.0.1")}}
	runner.Dialer = rejectDialer{}
	_, err := runner.Run(context.Background(), validation.Target{Hostname: "example.com"})
	if err == nil || !errors.Is(err, ssrf.ErrBlocked) {
		t.Fatalf("expected ssrf block, got %v", err)
	}
}

func TestRunDoesNotDialWhenDNSFails(t *testing.T) {
	runner := NewRunner()
	runner.Resolver = fakeResolver{err: errors.New("nxdomain")}
	runner.Dialer = rejectDialer{}
	runner.Now = func() time.Time { return time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC) }
	out, err := runner.Run(context.Background(), validation.Target{Hostname: "example.com"})
	if err != nil {
		t.Fatal(err)
	}
	if len(out) != 4 {
		t.Fatalf("expected 4 measurement slots, got %d", len(out))
	}
	if out[1].Measured || out[2].Measured || out[3].Measured {
		t.Fatal("tcp/tls/http must stay unmeasured when DNS fails")
	}
}

func TestHostCheckBlocksLocalhostBeforeLookup(t *testing.T) {
	runner := NewRunner()
	_, err := runner.Run(context.Background(), validation.Target{Hostname: "localhost"})
	if err == nil {
		t.Fatal("localhost must not be probed")
	}
}

func TestSafeDialPinsResolvedPublicIP(t *testing.T) {
	dialer := &recordDialer{}
	runner := NewRunner()
	runner.Resolver = fakeResolver{ips: []net.IP{net.ParseIP("8.8.8.8")}}
	runner.Dialer = dialer
	_, _ = runner.safeDial(context.Background(), "tcp", "example.com:443")
	if len(dialer.addrs) != 1 || dialer.addrs[0] != "8.8.8.8:443" {
		t.Fatalf("expected pinned public IP, dialed %v", dialer.addrs)
	}
}

func TestSafeDialDoesNotConnectToLoopback(t *testing.T) {
	dialer := &recordDialer{}
	runner := NewRunner()
	runner.Dialer = dialer
	_, err := runner.safeDial(context.Background(), "tcp", "127.0.0.1:443")
	if err == nil || !errors.Is(err, ssrf.ErrBlocked) {
		t.Fatalf("expected ssrf block, got %v", err)
	}
	if len(dialer.addrs) != 0 {
		t.Fatalf("loopback must not be dialed: %v", dialer.addrs)
	}
}

func TestRevalidateBlocksPrivateRedirect(t *testing.T) {
	runner := NewRunner()
	runner.Resolver = fakeResolver{ips: []net.IP{net.ParseIP("10.0.0.8")}}
	err := runner.revalidate(context.Background(), &url.URL{Scheme: "https", Host: "evil.example"})
	if err == nil || !errors.Is(err, ssrf.ErrBlocked) {
		t.Fatalf("expected redirect revalidation block, got %v", err)
	}
}

func TestRevalidateBlocksLoopbackURL(t *testing.T) {
	runner := NewRunner()
	err := runner.revalidate(context.Background(), &url.URL{Scheme: "https", Host: "127.0.0.1"})
	if err == nil {
		t.Fatal("redirect to loopback must be blocked")
	}
}

func TestSafeDialDoesNotConnectAfterRebindToPrivate(t *testing.T) {
	dialer := &recordDialer{}
	calls := 0
	runner := NewRunner()
	runner.Dialer = dialer
	runner.Resolver = lookupFunc(func(context.Context, string) ([]net.IP, error) {
		calls++
		if calls == 1 {
			return []net.IP{net.ParseIP("1.1.1.1")}, nil
		}
		return []net.IP{net.ParseIP("10.0.0.8")}, nil
	})
	_, _ = runner.safeDial(context.Background(), "tcp", "example.com:443")
	_, err := runner.safeDial(context.Background(), "tcp", "example.com:443")
	if err == nil || !errors.Is(err, ssrf.ErrBlocked) {
		t.Fatalf("expected rebind block, got %v", err)
	}
	if len(dialer.addrs) != 1 || dialer.addrs[0] != "1.1.1.1:443" {
		t.Fatalf("private rebind must not be dialed: %v", dialer.addrs)
	}
}

type lookupFunc func(context.Context, string) ([]net.IP, error)

func (f lookupFunc) LookupIP(ctx context.Context, host string) ([]net.IP, error) {
	return f(ctx, host)
}
