package validation

import "testing"

func TestParseTargetAcceptsPublicHosts(t *testing.T) {
	cases := []string{"example.com", "https://Status.Example.com/path", "youtube.com", "Google"}
	for _, raw := range cases {
		result := ParseTarget(raw)
		if result.Err != nil {
			t.Fatalf("%s: %v", raw, result.Err)
		}
	}
}

func TestParseTargetRejectsUnsafe(t *testing.T) {
	cases := []string{"", "localhost", "http://127.0.0.1/", "http://127.1/", "http://169.254.169.254/", "javascript:alert(1)", "https://user:pass@example.com", "<script>example.com"}
	for _, raw := range cases {
		result := ParseTarget(raw)
		if result.Err == nil {
			t.Fatalf("expected rejection for %q", raw)
		}
	}
}
