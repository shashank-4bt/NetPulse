package incidents

import "testing"

func TestCannotResolveFromOneRecovery(t *testing.T) {
	got := CanMarkResolved(ResolutionInput{RecoverySampleCount: 1, IdentifiedCause: true})
	if got.OK {
		t.Fatal("a single recovered measurement must not resolve an incident")
	}
}

func TestCannotResolveWithoutIdentifiedCause(t *testing.T) {
	got := CanMarkResolved(ResolutionInput{RecoverySampleCount: 4, IdentifiedCause: false})
	if got.OK {
		t.Fatal("recoveries without an identified cause must not resolve")
	}
}

func TestResolveNeedsMultipleRecoveriesAndCause(t *testing.T) {
	got := CanMarkResolved(ResolutionInput{RecoverySampleCount: 2, IdentifiedCause: true})
	if !got.OK {
		t.Fatal(got.Reason)
	}
}
