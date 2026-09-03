package ssrf

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"
)

type ClientOptions struct {
	Timeout      time.Duration
	MaxRedirects int
	HTTPSOnly    bool
	Lookup       LookupIPFunc
	Dialer       ContextDialer
	Policy       Policy
}

func NewClient(opts ClientOptions) *http.Client {
	if opts.Timeout <= 0 {
		opts.Timeout = 5 * time.Second
	}
	if opts.MaxRedirects <= 0 {
		opts.MaxRedirects = 3
	}
	policy := opts.Policy
	lookup := opts.Lookup
	dialer := opts.Dialer
	httpsOnly := opts.HTTPSOnly
	maxRedirects := opts.MaxRedirects
	return &http.Client{
		Timeout: opts.Timeout,
		Transport: &http.Transport{
			DisableKeepAlives:   true,
			TLSHandshakeTimeout: 5 * time.Second,
			DialContext: func(ctx context.Context, network, address string) (net.Conn, error) {
				return PinDial(ctx, policy, lookup, dialer, network, address)
			},
		},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= maxRedirects {
				return fmt.Errorf("redirect limit exceeded")
			}
			return CheckURL(policy, req.URL, httpsOnly)
		},
	}
}

func NewHTTPSClient() *http.Client {
	return NewClient(ClientOptions{HTTPSOnly: true, Timeout: 5 * time.Second})
}
