package developer

import (
	"context"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/auth"
	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/contract"
	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/ssrf"
	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/storage"
	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/validation"
)

type ProbeRunner interface {
	Run(ctx context.Context, target validation.Target) ([]contract.Measurement, error)
}

type Service struct {
	Store                storage.DeveloperStore
	Runner               ProbeRunner
	AllowPrivateDelivery bool
	HTTPClient           *http.Client
	Now                  func() time.Time
}

func (s *Service) now() time.Time {
	if s.Now != nil {
		return s.Now()
	}
	return time.Now()
}

func (s *Service) client() *http.Client {
	if s.HTTPClient != nil {
		return s.HTTPClient
	}
	return ssrf.NewHTTPSClient()
}

func (s *Service) WorkspaceForOwner(ctx context.Context, ownerID, name string) (contract.Workspace, *contract.APIError, int) {
	if name == "" {
		name = "Personal workspace"
	}
	ws, err := s.Store.GetOrCreateWorkspace(ctx, ownerID, name)
	if err != nil {
		return contract.Workspace{}, apiErr("unavailable", "Workspace store is unavailable."), 503
	}
	return ws, nil, 200
}

func (s *Service) GetWorkspace(ctx context.Context, actorWorkspaceID, id string) (*contract.Workspace, *contract.APIError, int) {
	if id == "" || id != actorWorkspaceID {
		return nil, apiErr("not_found", "workspace not found"), 404
	}
	ws, err := s.Store.GetWorkspace(ctx, id)
	if err != nil {
		return nil, apiErr("unavailable", "Workspace store is unavailable."), 503
	}
	if ws == nil {
		return nil, apiErr("not_found", "workspace not found"), 404
	}
	return ws, nil, 200
}

func (s *Service) NotifyDiagnosisCompleted(ctx context.Context, userID, diagnosisID string) {
	if userID == "" || diagnosisID == "" {
		return
	}
	ws, err := s.Store.GetWorkspaceByOwner(ctx, userID)
	if err != nil || ws == nil {
		return
	}
	s.EnqueueEvent(ctx, ws.ID, "diagnosis.completed", map[string]any{"diagnosisId": diagnosisID})
}

func (s *Service) LookupAPIKey(ctx context.Context, raw string) (*contract.APIKey, *contract.Workspace, *contract.APIError, int) {
	raw = strings.TrimSpace(raw)
	if raw == "" || !strings.HasPrefix(raw, "npk_") {
		return nil, nil, apiErr("unauthorized", "API key is not valid."), 401
	}
	key, err := s.Store.GetAPIKeyByHash(ctx, auth.HashSecret(raw))
	if err != nil {
		return nil, nil, apiErr("unavailable", "Key store is unavailable."), 503
	}
	if key == nil || key.Revoked {
		return nil, nil, apiErr("unauthorized", "API key is not valid."), 401
	}
	ws, err := s.Store.GetWorkspace(ctx, key.WorkspaceID)
	if err != nil || ws == nil {
		return nil, nil, apiErr("unauthorized", "API key is not valid."), 401
	}
	used := s.now().UTC().Format(time.RFC3339)
	key.LastUsedAt = &used
	_ = s.Store.UpdateAPIKey(ctx, *key)
	_ = s.Store.IncrUsage(ctx, ws.ID, "requests", 1)
	return key, ws, nil, 200
}

func apiErr(code, message string) *contract.APIError {
	return &contract.APIError{Code: code, Message: message}
}

func ValidateWebhookURL(raw string) *contract.APIError {
	parsed, err := url.Parse(strings.TrimSpace(raw))
	if err != nil || parsed.Host == "" {
		return apiErr("validation_error", "Enter an https webhook URL.")
	}
	if parsed.Scheme != "https" {
		return apiErr("validation_error", "Webhook URLs must use https.")
	}
	if parsed.User != nil {
		return apiErr("validation_error", "Webhook URLs with credentials are not accepted.")
	}
	if err := ssrf.DefaultPolicy().CheckHost(parsed.Hostname()); err != nil {
		return apiErr("ssrf_blocked", "Webhook URL is not allowed.")
	}
	return nil
}

func resolveWebhookHost(host string) error {
	ips, err := net.LookupIP(host)
	if err != nil {
		return err
	}
	return ssrf.DefaultPolicy().CheckResolved(host, ips)
}
