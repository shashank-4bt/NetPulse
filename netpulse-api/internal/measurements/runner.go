package measurements

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/contract"
	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/ssrf"
	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/validation"
)

const (
	dnsTimeout     = 5 * time.Second
	tcpTimeout     = 5 * time.Second
	tlsTimeout     = 8 * time.Second
	httpTimeout    = 8 * time.Second
	maxResponse    = 64 << 10
	maxRedirects   = 3
	maxConnections = 4
)

type Resolver interface {
	LookupIP(ctx context.Context, host string) ([]net.IP, error)
}

type Dialer interface {
	DialContext(ctx context.Context, network, address string) (net.Conn, error)
}

type Runner struct {
	Policy   ssrf.Policy
	Resolver Resolver
	Dialer   Dialer
	Now      func() time.Time
}

type defaultResolver struct{}

func (defaultResolver) LookupIP(ctx context.Context, host string) ([]net.IP, error) {
	ctx, cancel := context.WithTimeout(ctx, dnsTimeout)
	defer cancel()
	return net.DefaultResolver.LookupIP(ctx, "ip", host)
}

func NewRunner() *Runner {
	return &Runner{
		Policy:   ssrf.DefaultPolicy(),
		Resolver: defaultResolver{},
		Dialer:   &net.Dialer{Timeout: tcpTimeout, KeepAlive: -1},
		Now:      time.Now,
	}
}

func (r *Runner) Run(ctx context.Context, target validation.Target) ([]contract.Measurement, error) {
	if err := r.Policy.CheckHost(target.Hostname); err != nil {
		return nil, err
	}

	now := r.Now().UTC().Format(time.RFC3339)
	measurements := make([]contract.Measurement, 0, 4)

	ips, dnsErr := r.Resolver.LookupIP(ctx, target.Hostname)
	public := ssrf.PublicIPs(ips)
	dnsOK := dnsErr == nil && len(public) > 0
	measurements = append(measurements, measurement("dns", "DNS", strconv.Itoa(len(public)), "public_addresses", dnsOK, now, nil, dnsSummary(dnsErr, len(ips), len(public))))

	if dnsErr == nil && len(ips) > 0 {
		if err := r.Policy.CheckResolved(target.Hostname, ips); err != nil {
			return measurements, err
		}
	}
	if !dnsOK {
		measurements = append(measurements,
			unmeasured("tcp", "TCP", now),
			unmeasured("tls", "TLS", now),
			unmeasured("http", "HTTP", now),
		)
		return measurements, nil
	}

	ip := public[0]
	port := 443
	addr := net.JoinHostPort(ip.String(), strconv.Itoa(port))
	start := r.Now()
	tcpCtx, cancel := context.WithTimeout(ctx, tcpTimeout)
	conn, tcpErr := r.safeDial(tcpCtx, "tcp", addr)
	cancel()
	rtt := int(r.Now().Sub(start).Milliseconds())
	tcpOK := tcpErr == nil
	measurements = append(measurements, measurement("tcp", "TCP", tcpValue(tcpOK, rtt, port), "ms", tcpOK, now, ptr("service"), tcpSummary(tcpErr, port, rtt)))
	if conn != nil {
		_ = conn.Close()
	}

	tlsOK := false
	tlsVersion := ""
	if tcpOK {
		tlsOK, tlsVersion = r.handshakeTLS(ctx, addr, target.Hostname)
	}
	measurements = append(measurements, measurement("tls", "TLS", tlsLabel(tlsOK, tlsVersion), nil, tlsOK, now, ptr("service"), tlsSummary(tlsOK, tlsVersion)))

	httpOK := false
	httpSummary := "HTTP not attempted"
	if tlsOK || tcpOK {
		httpOK, httpSummary = r.fetchHTTP(ctx, target.Hostname)
	}
	measurements = append(measurements, measurement("http", "HTTP", httpLabel(httpOK), nil, httpOK, now, ptr("service"), httpSummary))
	return measurements, nil
}

func (r *Runner) handshakeTLS(ctx context.Context, addr, serverName string) (bool, string) {
	ctx, cancel := context.WithTimeout(ctx, tlsTimeout)
	defer cancel()
	raw, err := r.safeDial(ctx, "tcp", addr)
	if err != nil {
		return false, ""
	}
	tlsConn := tls.Client(raw, &tls.Config{ServerName: serverName, MinVersion: tls.VersionTLS12})
	defer tlsConn.Close()
	if err := tlsConn.HandshakeContext(ctx); err != nil {
		return false, ""
	}
	return true, tls.VersionName(tlsConn.ConnectionState().Version)
}

// safeDial resolves hostnames, re-checks every address against the SSRF
// policy, and dials only a public IP. HTTP redirects must not inherit a
// previous DNS answer (rebinding).
func (r *Runner) safeDial(ctx context.Context, network, address string) (net.Conn, error) {
	lookup := func(ctx context.Context, host string) ([]net.IP, error) {
		return r.Resolver.LookupIP(ctx, host)
	}
	return ssrf.PinDial(ctx, r.Policy, lookup, r.Dialer, network, address)
}

func (r *Runner) fetchHTTP(ctx context.Context, hostname string) (bool, string) {
	redirects := 0
	client := &http.Client{
		Timeout: httpTimeout,
		Transport: &http.Transport{
			DisableKeepAlives:   true,
			MaxConnsPerHost:     maxConnections,
			TLSHandshakeTimeout: tlsTimeout,
			DialContext:         r.safeDial,
		},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			redirects++
			if redirects > maxRedirects {
				return fmt.Errorf("redirect limit exceeded")
			}
			return r.revalidate(ctx, req.URL)
		},
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://"+hostname+"/", nil)
	if err != nil {
		return false, "request build failed"
	}
	req.Header.Set("User-Agent", "NetPulse-Measurement/0.6")
	resp, err := client.Do(req)
	if err != nil {
		return false, "http error: " + sanitizeErr(err)
	}
	defer resp.Body.Close()
	n, _ := io.Copy(io.Discard, io.LimitReader(resp.Body, maxResponse+1))
	limited := n > maxResponse
	summary := fmt.Sprintf("status=%d redirects=%d bytes_read=%d truncated=%t final_host=%s", resp.StatusCode, redirects, min64(n, maxResponse), limited, resp.Request.URL.Hostname())
	return resp.StatusCode > 0, summary
}

func (r *Runner) revalidate(ctx context.Context, next *url.URL) error {
	if err := ssrf.CheckURL(r.Policy, next, false); err != nil {
		return err
	}
	host := next.Hostname()
	ips, err := r.Resolver.LookupIP(ctx, host)
	if err != nil {
		return err
	}
	return r.Policy.CheckResolved(host, ips)
}

func measurement(key, label string, value any, unit any, measured bool, at string, layer *string, summary string) contract.Measurement {
	var unitPtr *string
	if text, ok := unit.(string); ok {
		unitPtr = &text
	}
	return contract.Measurement{
		ID:         "measurement-" + key,
		Key:        key,
		Label:      label,
		Value:      value,
		Unit:       unitPtr,
		Measured:   measured,
		MeasuredAt: &at,
		Layer:      layer,
		Summary:    &summary,
	}
}

func unmeasured(key, label, at string) contract.Measurement {
	summary := "Not measured because DNS did not yield a public address."
	return contract.Measurement{
		ID:         "measurement-" + key,
		Key:        key,
		Label:      label,
		Value:      nil,
		Measured:   false,
		MeasuredAt: nil,
		Summary:    &summary,
	}
}

func ptr(value string) *string { return &value }

func dnsSummary(err error, resolved, public int) string {
	if err != nil {
		return "dns lookup failed"
	}
	return fmt.Sprintf("resolved=%d public=%d private_filtered=%d", resolved, public, resolved-public)
}

func tcpValue(ok bool, rtt, port int) string {
	if !ok {
		return "connect_failed"
	}
	return fmt.Sprintf("%d", rtt)
}

func tcpSummary(err error, port, rtt int) string {
	if err != nil {
		return fmt.Sprintf("tcp port %d failed", port)
	}
	return fmt.Sprintf("tcp port %d rtt_ms=%d", port, rtt)
}

func tlsLabel(ok bool, version string) string {
	if !ok {
		return "handshake_failed"
	}
	return version
}

func tlsSummary(ok bool, version string) string {
	if !ok {
		return "tls handshake failed"
	}
	return "tls handshake ok version=" + version
}

func httpLabel(ok bool) string {
	if !ok {
		return "request_failed"
	}
	return "completed"
}

func sanitizeErr(err error) string {
	text := err.Error()
	if strings.Contains(text, "ssrf") {
		return "ssrf_blocked"
	}
	if len(text) > 160 {
		return text[:160]
	}
	return text
}

func min64(value, limit int64) int64 {
	if value > limit {
		return limit
	}
	return value
}
