package developer

import (
	"net/http"
	"testing"
)

func TestDefaultWebhookClientBlocksPrivateAndHTTPRedirects(t *testing.T) {
	svc := &Service{}
	client := svc.client()
	if client.CheckRedirect == nil {
		t.Fatal("webhook client must revalidate redirects")
	}
	loopback, _ := http.NewRequest(http.MethodPost, "https://127.0.0.1/hook", nil)
	if err := client.CheckRedirect(loopback, []*http.Request{loopback}); err == nil {
		t.Fatal("redirect to loopback must be blocked")
	}
	rfc1918, _ := http.NewRequest(http.MethodPost, "https://10.0.0.8/hook", nil)
	if err := client.CheckRedirect(rfc1918, []*http.Request{rfc1918}); err == nil {
		t.Fatal("redirect to private IP must be blocked")
	}
	plain, _ := http.NewRequest(http.MethodPost, "http://example.com/hook", nil)
	if err := client.CheckRedirect(plain, []*http.Request{plain}); err == nil {
		t.Fatal("http redirect must be blocked")
	}
	ok, _ := http.NewRequest(http.MethodPost, "https://example.com/hook", nil)
	if err := client.CheckRedirect(ok, []*http.Request{ok}); err != nil {
		t.Fatalf("public https redirect host should pass CheckHost: %v", err)
	}
}
