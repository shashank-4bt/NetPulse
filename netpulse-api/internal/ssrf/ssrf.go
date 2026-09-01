package ssrf

import (
	"fmt"
	"net"
	"net/netip"
	"strings"
)

var (
	ErrBlocked = fmt.Errorf("ssrf_blocked")
)

var blockedHosts = map[string]struct{}{
	"localhost":                 {},
	"localhost.localdomain":     {},
	"metadata.google.internal":  {},
	"metadata.google.internal.": {},
	"metadata.aws.internal":     {},
}

type Policy struct{}

func DefaultPolicy() Policy {
	return Policy{}
}

func (p Policy) CheckHost(host string) error {
	hostname := normalizeHost(host)
	if hostname == "" {
		return fmt.Errorf("%w: empty host", ErrBlocked)
	}
	if _, blocked := blockedHosts[hostname]; blocked {
		return fmt.Errorf("%w: blocked host %s", ErrBlocked, hostname)
	}
	if strings.HasSuffix(hostname, ".local") || strings.HasSuffix(hostname, ".internal") {
		return fmt.Errorf("%w: internal host %s", ErrBlocked, hostname)
	}
	if looksLikeIPv4(hostname) {
		ip, err := netip.ParseAddr(hostname)
		if err != nil {
			return fmt.Errorf("%w: invalid ipv4 host %s", ErrBlocked, hostname)
		}
		return p.CheckIP(ip)
	}
	if ip, err := netip.ParseAddr(hostname); err == nil {
		return p.CheckIP(ip)
	}
	return nil
}

func (p Policy) CheckIP(ip netip.Addr) error {
	addr := ip.Unmap()
	if !addr.IsValid() {
		return fmt.Errorf("%w: invalid ip", ErrBlocked)
	}
	if addr.IsLoopback() || addr.IsPrivate() || addr.IsLinkLocalUnicast() || addr.IsLinkLocalMulticast() || addr.IsMulticast() || addr.IsUnspecified() {
		return fmt.Errorf("%w: non-public ip %s", ErrBlocked, addr)
	}
	if isCarrierGradeNAT(addr) || isDocumentation(addr) || isCloudMetadata(addr) {
		return fmt.Errorf("%w: reserved ip %s", ErrBlocked, addr)
	}
	return nil
}

func (p Policy) CheckResolved(host string, ips []net.IP) error {
	if err := p.CheckHost(host); err != nil {
		return err
	}
	if len(ips) == 0 {
		return fmt.Errorf("%w: no resolved addresses", ErrBlocked)
	}
	for _, ip := range ips {
		parsed, ok := netip.AddrFromSlice(ip)
		if !ok {
			return fmt.Errorf("%w: invalid resolved address", ErrBlocked)
		}
		if err := p.CheckIP(parsed); err != nil {
			return err
		}
	}
	return nil
}

func PublicIPs(ips []net.IP) []net.IP {
	policy := DefaultPolicy()
	out := make([]net.IP, 0, len(ips))
	for _, ip := range ips {
		parsed, ok := netip.AddrFromSlice(ip)
		if !ok {
			continue
		}
		if policy.CheckIP(parsed) == nil {
			out = append(out, ip)
		}
	}
	return out
}

func looksLikeIPv4(host string) bool {
	if host == "" {
		return false
	}
	for _, part := range strings.Split(host, ".") {
		if part == "" {
			return false
		}
		for _, c := range part {
			if c < '0' || c > '9' {
				return false
			}
		}
	}
	return true
}

func normalizeHost(host string) string {
	host = strings.TrimSpace(strings.ToLower(host))
	host = strings.TrimPrefix(host, "[")
	host = strings.TrimSuffix(host, "]")
	host = strings.TrimSuffix(host, ".")
	if h, _, err := net.SplitHostPort(host); err == nil {
		return normalizeHost(h)
	}
	return host
}

func isCarrierGradeNAT(ip netip.Addr) bool {
	if !ip.Is4() {
		return false
	}
	octets := ip.As4()
	return octets[0] == 100 && octets[1] >= 64 && octets[1] <= 127
}

func isDocumentation(ip netip.Addr) bool {
	if !ip.Is4() {
		return false
	}
	octets := ip.As4()
	return (octets[0] == 192 && octets[1] == 0 && octets[2] == 2) ||
		(octets[0] == 198 && octets[1] == 51 && octets[2] == 100) ||
		(octets[0] == 203 && octets[1] == 0 && octets[2] == 113)
}

func isCloudMetadata(ip netip.Addr) bool {
	if !ip.Is4() {
		return false
	}
	octets := ip.As4()
	return octets[0] == 169 && octets[1] == 254
}
