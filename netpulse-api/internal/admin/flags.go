package admin

import (
	"context"
	"fmt"
	"strings"

	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/contract"
	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/id"
	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/opsconfig"
)

func EvaluateFlag(flag contract.FeatureFlag, env, userID, orgID string) bool {
	if !flag.Enabled {
		return false
	}
	if flag.Environment != "" && !strings.EqualFold(flag.Environment, env) {
		return false
	}
	for _, item := range flag.UserIDs {
		if item != "" && item == userID {
			return true
		}
	}
	for _, item := range flag.OrgIDs {
		if item != "" && item == orgID {
			return true
		}
	}
	if flag.Percentage <= 0 {
		return false
	}
	if flag.Percentage >= 100 {
		return true
	}
	subject := userID
	if subject == "" {
		subject = orgID
	}
	if subject == "" {
		return false
	}
	h := 0
	for _, c := range subject {
		h = (h*31 + int(c)) % 100
	}
	if h < 0 {
		h = -h
	}
	return h < flag.Percentage
}

func (s *Service) ListFlags(ctx context.Context) ([]contract.FeatureFlag, *contract.APIError, int) {
	items, err := s.Store.ListFlags(ctx)
	if err != nil {
		return nil, apiErr("unavailable", "Feature flag store is unavailable."), 503
	}
	if items == nil {
		items = []contract.FeatureFlag{}
	}
	return items, nil, 200
}

func (s *Service) GetFlag(ctx context.Context, flagID, env, userID, orgID string) (*contract.FeatureFlag, *contract.APIError, int) {
	item, err := s.Store.GetFlag(ctx, flagID)
	if err != nil {
		return nil, apiErr("unavailable", "Feature flag store is unavailable."), 503
	}
	if item == nil {
		return nil, apiErr("not_found", "not found"), 404
	}
	if env != "" || userID != "" || orgID != "" {
		if env == "" {
			env = s.Cfg.Environment
		}
		match := EvaluateFlag(*item, env, userID, orgID)
		item.TargetMatch = &match
		item.Summary = flagSummary(*item) + fmt.Sprintf(" Target match for the supplied identifiers: %t. This is not a user count.", match)
	}
	return item, nil, 200
}

func (s *Service) CreateFlag(ctx context.Context, actorID string, flag contract.FeatureFlag) (*contract.FeatureFlag, *contract.APIError, int) {
	name := strings.TrimSpace(flag.Name)
	if name == "" {
		return nil, apiErr("validation_error", "A flag name is required."), 400
	}
	if flag.Percentage < 0 || flag.Percentage > 100 {
		return nil, apiErr("validation_error", "Percentage must be between 0 and 100."), 400
	}
	created := contract.FeatureFlag{
		ID: id.New(), Name: name, Environment: strings.TrimSpace(flag.Environment),
		Enabled: flag.Enabled, Percentage: flag.Percentage,
		UserIDs: append([]string{}, flag.UserIDs...), OrgIDs: append([]string{}, flag.OrgIDs...),
		UpdatedAt: s.stamp(), Summary: flagSummary(flag),
	}
	if created.UserIDs == nil {
		created.UserIDs = []string{}
	}
	if created.OrgIDs == nil {
		created.OrgIDs = []string{}
	}
	if err := s.Store.UpsertFlag(ctx, created); err != nil {
		return nil, apiErr("unavailable", "Feature flag store is unavailable."), 503
	}
	s.Audit(ctx, actorID, "flag.create", "flag:"+created.ID, "created", "A feature flag was stored.")
	return &created, nil, 201
}

func (s *Service) PatchFlag(ctx context.Context, actorID, flagID string, flag contract.FeatureFlag, hasEnabled, hasPercentage bool) (*contract.FeatureFlag, *contract.APIError, int) {
	existing, errResp, status := s.GetFlag(ctx, flagID, "", "", "")
	if errResp != nil {
		return nil, errResp, status
	}
	if name := strings.TrimSpace(flag.Name); name != "" {
		existing.Name = name
	}
	if env := strings.TrimSpace(flag.Environment); env != "" {
		existing.Environment = env
	}
	if hasEnabled {
		existing.Enabled = flag.Enabled
	}
	if hasPercentage {
		if flag.Percentage < 0 || flag.Percentage > 100 {
			return nil, apiErr("validation_error", "Percentage must be between 0 and 100."), 400
		}
		existing.Percentage = flag.Percentage
	}
	if flag.UserIDs != nil {
		existing.UserIDs = append([]string{}, flag.UserIDs...)
	}
	if flag.OrgIDs != nil {
		existing.OrgIDs = append([]string{}, flag.OrgIDs...)
	}
	existing.UpdatedAt = s.stamp()
	existing.Summary = flagSummary(*existing)
	existing.TargetMatch = nil
	if err := s.Store.UpsertFlag(ctx, *existing); err != nil {
		return nil, apiErr("unavailable", "Feature flag store is unavailable."), 503
	}
	s.Audit(ctx, actorID, "flag.update", "flag:"+flagID, "updated", "A feature flag was updated.")
	return existing, nil, 200
}

func flagSummary(flag contract.FeatureFlag) string {
	return fmt.Sprintf(
		"Environment %s. Percentage %d. User targets %d. Organization targets %d. Counts of users currently in rollout are not estimated.",
		emptyEnv(flag.Environment), flag.Percentage, len(flag.UserIDs), len(flag.OrgIDs),
	)
}

func emptyEnv(env string) string {
	if strings.TrimSpace(env) == "" {
		return "any"
	}
	return env
}

func (s *Service) ListConfig(ctx context.Context) ([]contract.RemoteConfigEntry, *contract.APIError, int) {
	s.Seed(ctx)
	items, err := s.Store.ListRemoteConfig(ctx)
	if err != nil {
		return nil, apiErr("unavailable", "Configuration store is unavailable."), 503
	}
	if items == nil {
		items = []contract.RemoteConfigEntry{}
	}
	return items, nil, 200
}

func (s *Service) PutConfig(ctx context.Context, actorID, key, value string) ([]contract.RemoteConfigEntry, *contract.APIError, int) {
	key = strings.TrimSpace(key)
	value = strings.TrimSpace(value)
	if !opsconfig.KnownKey(key) {
		return nil, apiErr("validation_error", "Unknown configuration key."), 400
	}
	if opsconfig.LooksLikeSecret(key, value) {
		return nil, apiErr("validation_error", "Secrets are not stored in remote configuration."), 400
	}
	if value == "" {
		return nil, apiErr("validation_error", "A configuration value is required."), 400
	}
	existing, _ := s.Store.GetRemoteConfig(ctx, key)
	summary := "Operational setting."
	if existing != nil && existing.Summary != "" {
		summary = existing.Summary
	} else {
		for _, item := range opsconfig.Defaults(s.Cfg) {
			if item.Key == key {
				summary = item.Summary
				break
			}
		}
	}
	entry := contract.RemoteConfigEntry{Key: key, Value: value, UpdatedAt: s.stamp(), Summary: summary}
	if err := s.Store.UpsertRemoteConfig(ctx, entry); err != nil {
		return nil, apiErr("unavailable", "Configuration store is unavailable."), 503
	}
	s.Audit(ctx, actorID, "config.update", "config:"+key, "updated", "Remote configuration was updated.")
	return s.ListConfig(ctx)
}
