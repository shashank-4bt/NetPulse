package admin

import (
	"testing"

	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/contract"
)

func TestEvaluateFlag(t *testing.T) {
	flag := contract.FeatureFlag{Enabled: true, Environment: "test", Percentage: 0, UserIDs: []string{"user-1"}, OrgIDs: []string{}}
	if !EvaluateFlag(flag, "test", "user-1", "") {
		t.Fatal("listed user must match")
	}
	if EvaluateFlag(flag, "test", "user-2", "") {
		t.Fatal("other user must not match a 0% targeted flag")
	}
	flag.UserIDs = []string{}
	flag.OrgIDs = []string{"org-9"}
	if !EvaluateFlag(flag, "test", "", "org-9") {
		t.Fatal("listed organization must match")
	}
	flag.OrgIDs = []string{}
	flag.Percentage = 100
	if !EvaluateFlag(flag, "test", "anyone", "") {
		t.Fatal("100 percent rollout must include the subject")
	}
	flag.Enabled = false
	if EvaluateFlag(flag, "test", "anyone", "") {
		t.Fatal("disabled flags must not match")
	}
	flag.Enabled = true
	flag.Percentage = 50
	if EvaluateFlag(flag, "production", "anyone", "") {
		t.Fatal("environment mismatch must not match")
	}
}

func TestHasPerm(t *testing.T) {
	if !HasPerm(AllPermissions(), PermSystemRead) {
		t.Fatal("all permissions must include system.read")
	}
	if HasPerm([]string{PermSystemRead}, PermUsersRead) {
		t.Fatal("missing permission must be false")
	}
}
