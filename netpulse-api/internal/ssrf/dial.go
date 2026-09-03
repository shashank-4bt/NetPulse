package ssrf

import (
	"context"
	"fmt"
	"net"
	"net/netip"
	"time"
)

type LookupIPFunc func(ctx context.Context, host string) ([]net.IP, error)

type ContextDialer interface {
	DialContext(ctx context.Context, network, address string) (net.Conn, error)
}

func defaultLookup(ctx context.Context, host string) ([]net.IP, error) {
	return net.DefaultResolver.LookupIP(ctx, "ip", host)
}

func PinDial(ctx context.Context, policy Policy, lookup LookupIPFunc, dialer ContextDialer, network, address string) (net.Conn, error) {
	if lookup == nil {
		lookup = defaultLookup
	}
	if dialer == nil {
		dialer = &net.Dialer{Timeout: 5 * time.Second, KeepAlive: -1}
	}
	host, port, err := net.SplitHostPort(address)
	if err != nil {
		return nil, err
	}
	if err := policy.CheckHost(host); err != nil {
		return nil, err
	}

	dialAddr := address
	if parsed := net.ParseIP(host); parsed != nil {
		addr, ok := netip.AddrFromSlice(parsed)
		if !ok {
			return nil, fmt.Errorf("%w: invalid ip", ErrBlocked)
		}
		if err := policy.CheckIP(addr); err != nil {
			return nil, err
		}
		dialAddr = net.JoinHostPort(parsed.String(), port)
	} else {
		ips, lookupErr := lookup(ctx, host)
		if lookupErr != nil {
			return nil, lookupErr
		}
		if err := policy.CheckResolved(host, ips); err != nil {
			return nil, err
		}
		public := PublicIPs(ips)
		if len(public) == 0 {
			return nil, fmt.Errorf("%w: no public address", ErrBlocked)
		}
		dialAddr = net.JoinHostPort(public[0].String(), port)
	}
	return dialer.DialContext(ctx, network, dialAddr)
}
