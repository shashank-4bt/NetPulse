package ssrf

import (
	"context"
	"errors"
	"net"
	"testing"
)

type recordDialer struct {
	addrs []string
}

func (d *recordDialer) DialContext(_ context.Context, _, address string) (net.Conn, error) {
	d.addrs = append(d.addrs, address)
	return nil, errors.New("dial disabled in test")
}

func TestPinDialPinsPublicIP(t *testing.T) {
	dialer := &recordDialer{}
	lookup := func(context.Context, string) ([]net.IP, error) {
		return []net.IP{net.ParseIP("8.8.8.8")}, nil
	}
	_, _ = PinDial(context.Background(), DefaultPolicy(), lookup, dialer, "tcp", "example.com:443")
	if len(dialer.addrs) != 1 || dialer.addrs[0] != "8.8.8.8:443" {
		t.Fatalf("expected pinned public IP, dialed %v", dialer.addrs)
	}
}

func TestPinDialDoesNotConnectToLoopback(t *testing.T) {
	dialer := &recordDialer{}
	_, err := PinDial(context.Background(), DefaultPolicy(), nil, dialer, "tcp", "127.0.0.1:443")
	if err == nil || !errors.Is(err, ErrBlocked) {
		t.Fatalf("expected ssrf block, got %v", err)
	}
	if len(dialer.addrs) != 0 {
		t.Fatalf("loopback must not be dialed: %v", dialer.addrs)
	}
}

func TestPinDialDoesNotConnectAfterDNSRebind(t *testing.T) {
	dialer := &recordDialer{}
	calls := 0
	lookup := func(context.Context, string) ([]net.IP, error) {
		calls++
		if calls == 1 {
			return []net.IP{net.ParseIP("1.1.1.1")}, nil
		}
		return []net.IP{net.ParseIP("127.0.0.1")}, nil
	}
	_, _ = PinDial(context.Background(), DefaultPolicy(), lookup, dialer, "tcp", "example.com:443")
	if len(dialer.addrs) != 1 || dialer.addrs[0] != "1.1.1.1:443" {
		t.Fatalf("first lookup must pin a public IP, dialed %v", dialer.addrs)
	}
	_, err := PinDial(context.Background(), DefaultPolicy(), lookup, dialer, "tcp", "example.com:443")
	if err == nil || !errors.Is(err, ErrBlocked) {
		t.Fatalf("rebinding to loopback must be blocked, got %v", err)
	}
	if len(dialer.addrs) != 1 {
		t.Fatalf("rebind must not dial a second address: %v", dialer.addrs)
	}
}
