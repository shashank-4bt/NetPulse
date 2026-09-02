package developer

import (
	"context"
	"strings"

	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/auth"
	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/contract"
	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/id"
)

type APIKeyInput struct {
	Name            string
	Scopes          []string
	RateLimitPerMin int
}

func (s *Service) ListKeys(ctx context.Context, workspaceID string) ([]contract.APIKey, *contract.APIError, int) {
	items, err := s.Store.ListAPIKeys(ctx, workspaceID)
	if err != nil {
		return nil, apiErr("unavailable", "Key store is unavailable."), 503
	}
	if items == nil {
		items = []contract.APIKey{}
	}
	return items, nil, 200
}

func (s *Service) CreateKey(ctx context.Context, workspaceID string, in APIKeyInput) (*contract.APIKey, string, *contract.APIError, int) {
	stamp := s.now().UTC().Format("2006-01-02T15:04:05Z07:00")
	item, secret, errResp := buildKey(workspaceID, "", in, stamp)
	if errResp != nil {
		return nil, "", errResp, statusFor(errResp)
	}
	if err := s.Store.CreateAPIKey(ctx, *item); err != nil {
		return nil, "", apiErr("unavailable", "Key store is unavailable."), 503
	}
	return item, secret, nil, 201
}

func (s *Service) RotateKey(ctx context.Context, workspaceID, keyID string) (*contract.APIKey, string, *contract.APIError, int) {
	existing, err := s.Store.GetAPIKey(ctx, keyID)
	if err != nil {
		return nil, "", apiErr("unavailable", "Key store is unavailable."), 503
	}
	if existing == nil || existing.WorkspaceID != workspaceID || existing.Revoked {
		return nil, "", apiErr("not_found", "api key not found"), 404
	}
	raw, prefix, last4, hash := auth.NewAPIKeySecret()
	existing.Prefix = prefix
	existing.Last4 = last4
	existing.Hash = hash
	if err := s.Store.UpdateAPIKey(ctx, *existing); err != nil {
		return nil, "", apiErr("unavailable", "Key store is unavailable."), 503
	}
	return existing, raw, nil, 200
}

func (s *Service) RevokeKey(ctx context.Context, workspaceID, keyID string) (*contract.APIKey, *contract.APIError, int) {
	existing, err := s.Store.GetAPIKey(ctx, keyID)
	if err != nil {
		return nil, apiErr("unavailable", "Key store is unavailable."), 503
	}
	if existing == nil || existing.WorkspaceID != workspaceID {
		return nil, apiErr("not_found", "api key not found"), 404
	}
	existing.Revoked = true
	if err := s.Store.UpdateAPIKey(ctx, *existing); err != nil {
		return nil, apiErr("unavailable", "Key store is unavailable."), 503
	}
	return existing, nil, 200
}

func buildKey(workspaceID, existingID string, in APIKeyInput, createdAt string) (*contract.APIKey, string, *contract.APIError) {
	name := strings.TrimSpace(in.Name)
	if name == "" {
		name = "API key"
	}
	scopes := NormalizeScopes(in.Scopes)
	if len(in.Scopes) > 0 && len(scopes) == 0 {
		return nil, "", apiErr("validation_error", "No recognized scopes were supplied.")
	}
	limit := in.RateLimitPerMin
	if limit == 0 {
		limit = 60
	}
	if limit < 1 || limit > 600 {
		return nil, "", apiErr("validation_error", "Rate limit must be between 1 and 600 requests per minute.")
	}
	raw, prefix, last4, hash := auth.NewAPIKeySecret()
	itemID := existingID
	if itemID == "" {
		itemID = id.New()
	}
	return &contract.APIKey{
		ID:              itemID,
		WorkspaceID:     workspaceID,
		Name:            name,
		Prefix:          prefix,
		Last4:           last4,
		Hash:            hash,
		Scopes:          scopes,
		RateLimitPerMin: limit,
		CreatedAt:       createdAt,
	}, raw, nil
}
