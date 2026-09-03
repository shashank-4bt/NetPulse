package ssrf

import (
	"fmt"
	"net"
	"net/netip"
	"net/url"
	"strings"
)

var (
	ErrBlocked = fmt.Errorf("ssrf_blocked")
)

var blockedHosts = map[string]struct{}{
	"localhost":                            {},
	"localhost.localdomain":                {},
	"metadata":                             {},
	"metadata.google.internal":             {},
	"metadata.google.internal.":            {},
	"metadata.aws.internal":                {},
	"instance-data":                        {},
	"kubernetes":                           {},
	"kubernetes.default":                   {},
	"kubernetes.default.svc":               {},
	"kubernetes.default.svc.cluster.local": {},
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
	if strings.HasSuffix(hostname, ".local") ||
		strings.HasSuffix(hostname, ".internal") ||
		strings.HasSuffix(hostname, ".localhost") ||
		strings.HasSuffix(hostname, ".lan") ||
		strings.HasSuffix(hostname, ".corp") {
		return fmt.Errorf("%w: internal host %s", ErrBlocked, hostname)
	}
	if strings.Contains(hostname, "0x") {
		return fmt.Errorf("%w: hex-encoded host %s", ErrBlocked, hostname)
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
	if parsed := net.ParseIP(hostname); parsed != nil {
		addr, ok := netip.AddrFromSlice(parsed)
		if !ok {
			return fmt.Errorf("%w: invalid ip host %s", ErrBlocked, hostname)
		}
		return p.CheckIP(addr)
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
	if isCarrierGradeNAT(addr) || isDocumentation(addr) || isCloudMetadata(addr) || isReservedIPv4(addr) || isBlockedIPv6(addr) {
		return fmt.Errorf("%w: reserved ip %s", ErrBlocked, addr)
	}
	if embedded, ok := nat64Embedded(addr); ok {
		if err := p.CheckIP(embedded); err != nil {
			return fmt.Errorf("%w: nat64 embedded %s", ErrBlocked, addr)
		}
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

func CheckURL(policy Policy, next *url.URL, httpsOnly bool) error {
	if next == nil {
		return fmt.Errorf("%w: empty url", ErrBlocked)
	}
	if httpsOnly {
		if next.Scheme != "https" {
			return fmt.Errorf("%w: scheme", ErrBlocked)
		}
	} else if next.Scheme != "http" && next.Scheme != "https" {
		return fmt.Errorf("%w: scheme", ErrBlocked)
	}
	if next.User != nil {
		return fmt.Errorf("%w: credentials", ErrBlocked)
	}
	return policy.CheckHost(next.Hostname())
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
	if ip.Is4() {
		octets := ip.As4()
		return (octets[0] == 192 && octets[1] == 0 && octets[2] == 2) ||
			(octets[0] == 198 && octets[1] == 51 && octets[2] == 100) ||
			(octets[0] == 203 && octets[1] == 0 && octets[2] == 113)
	}
	if ip.Is6() {
		a := ip.As16()
		return a[0] == 0x20 && a[1] == 0x01 && a[2] == 0x0d && a[3] == 0xb8
	}
	return false
}

func isCloudMetadata(ip netip.Addr) bool {
	if !ip.Is4() {
		return false
	}
	octets := ip.As4()
	if octets[0] == 169 && octets[1] == 254 {
		return true
	}
	return octets[0] == 168 && octets[1] == 63 && octets[2] == 129 && octets[3] == 16
}

func isReservedIPv4(ip netip.Addr) bool {
	if !ip.Is4() {
		return false
	}
	octets := ip.As4()
	if octets[0] == 0 {
		return true
	}
	if octets[0] >= 240 {
		return true
	}
	if octets[0] == 192 && octets[1] == 0 && octets[2] == 0 {
		return true
	}
	if octets[0] == 192 && octets[1] == 88 && octets[2] == 99 {
		return true
	}
	if octets[0] == 198 && (octets[1] == 18 || octets[1] == 19) {
		return true
	}
	return false
}

func isBlockedIPv6(ip netip.Addr) bool {
	if !ip.Is6() {
		return false
	}
	a := ip.As16()
	if a[0] == 0x20 && a[1] == 0x02 {
		return true
	}
	if a[0] == 0x20 && a[1] == 0x01 && a[2] == 0x00 && a[3] == 0x00 {
		return true
	}
	if a[0] == 0x01 && a[1] == 0x00 {
		for i := 2; i < 8; i++ {
			if a[i] != 0 {
				return false
			}
		}
		return true
	}
	if a[0] == 0xfe && a[1]&0xc0 == 0xc0 {
		return true
	}
	return false
}

func nat64Embedded(ip netip.Addr) (netip.Addr, bool) {
	if !ip.Is6() {
		return netip.Addr{}, false
	}
	a := ip.As16()
	if a[0] != 0x00 || a[1] != 0x64 || a[2] != 0xff || a[3] != 0x9b {
		return netip.Addr{}, false
	}
	for i := 4; i < 12; i++ {
		if a[i] != 0 {
			return netip.Addr{}, false
		}
	}
	return netip.AddrFrom4([4]byte{a[12], a[13], a[14], a[15]}), true
}
