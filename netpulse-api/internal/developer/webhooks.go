package developer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/auth"
	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/contract"
	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/id"
)

type WebhookInput struct {
	URL    string
	Events []string
}

func (s *Service) ListWebhooks(ctx context.Context, workspaceID string) ([]contract.Webhook, *contract.APIError, int) {
	items, err := s.Store.ListWebhooks(ctx, workspaceID)
	if err != nil {
		return nil, apiErr("unavailable", "Webhook store is unavailable."), 503
	}
	if items == nil {
		items = []contract.Webhook{}
	}
	return items, nil, 200
}

func (s *Service) CreateWebhook(ctx context.Context, workspaceID string, in WebhookInput) (*contract.Webhook, string, *contract.APIError, int) {
	item, secret, errResp := buildWebhook(workspaceID, in)
	if errResp != nil {
		return nil, "", errResp, statusFor(errResp)
	}
	item.CreatedAt = s.now().UTC().Format(time.RFC3339)
	if err := s.Store.CreateWebhook(ctx, *item); err != nil {
		return nil, "", apiErr("unavailable", "Webhook store is unavailable."), 503
	}
	return item, secret, nil, 201
}

func (s *Service) RotateWebhook(ctx context.Context, workspaceID, webhookID string) (*contract.Webhook, string, *contract.APIError, int) {
	existing, err := s.Store.GetWebhook(ctx, webhookID)
	if err != nil {
		return nil, "", apiErr("unavailable", "Webhook store is unavailable."), 503
	}
	if existing == nil || existing.WorkspaceID != workspaceID {
		return nil, "", apiErr("not_found", "webhook not found"), 404
	}
	raw, hint := auth.NewWebhookSecret()
	existing.Secret = raw
	existing.SecretHint = hint
	if err := s.Store.UpdateWebhook(ctx, *existing); err != nil {
		return nil, "", apiErr("unavailable", "Webhook store is unavailable."), 503
	}
	return existing, raw, nil, 200
}

func (s *Service) DeleteWebhook(ctx context.Context, workspaceID, webhookID string) *contract.APIError {
	ok, err := s.Store.DeleteWebhook(ctx, workspaceID, webhookID)
	if err != nil {
		return apiErr("unavailable", "Webhook store is unavailable.")
	}
	if !ok {
		return apiErr("not_found", "webhook not found")
	}
	return nil
}

func (s *Service) ListDeliveries(ctx context.Context, workspaceID, webhookID string) ([]contract.WebhookDelivery, *contract.APIError, int) {
	hook, err := s.Store.GetWebhook(ctx, webhookID)
	if err != nil {
		return nil, apiErr("unavailable", "Webhook store is unavailable."), 503
	}
	if hook == nil || hook.WorkspaceID != workspaceID {
		return nil, apiErr("not_found", "webhook not found"), 404
	}
	items, err := s.Store.ListDeliveries(ctx, workspaceID, webhookID)
	if err != nil {
		return nil, apiErr("unavailable", "Delivery store is unavailable."), 503
	}
	if items == nil {
		items = []contract.WebhookDelivery{}
	}
	return items, nil, 200
}

func (s *Service) RetryDeliveries(ctx context.Context, workspaceID, webhookID string) ([]contract.WebhookDelivery, *contract.APIError, int) {
	items, errResp, status := s.ListDeliveries(ctx, workspaceID, webhookID)
	if errResp != nil {
		return nil, errResp, status
	}
	out := []contract.WebhookDelivery{}
	for _, item := range items {
		if item.Status != "retryable" || item.Attempt >= 3 {
			continue
		}
		updated := s.attemptDelivery(ctx, item)
		out = append(out, updated)
	}
	return out, nil, 200
}

func (s *Service) EnqueueEvent(ctx context.Context, workspaceID, event string, data map[string]any) {
	if !KnownWebhookEvent(event) {
		return
	}
	hooks, err := s.Store.ListWebhooks(ctx, workspaceID)
	if err != nil {
		return
	}
	eventID := id.New()
	timestamp := s.now().UTC().Format(time.RFC3339)
	payload := map[string]any{
		"event":     event,
		"eventId":   eventID,
		"timestamp": timestamp,
		"data":      data,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return
	}
	for _, hook := range hooks {
		if hook.Disabled || !contains(hook.Events, event) {
			continue
		}
		s.enqueueHook(ctx, hook, event, eventID, timestamp, body)
	}
}

func (s *Service) enqueueHook(ctx context.Context, hook contract.Webhook, event, eventID, timestamp string, body []byte) {
	idem := eventID + ":" + hook.ID
	existing, err := s.Store.GetDeliveryByIdempotency(ctx, idem)
	if err != nil || existing != nil {
		return
	}
	sig := auth.SignWebhook(hook.Secret, timestamp, eventID, string(body))
	item := contract.WebhookDelivery{
		ID:             id.New(),
		WebhookID:      hook.ID,
		WorkspaceID:    hook.WorkspaceID,
		Event:          event,
		EventID:        eventID,
		Timestamp:      timestamp,
		IdempotencyKey: idem,
		Signature:      sig,
		Payload:        string(body),
		Attempt:        0,
		Status:         "queued",
		Summary:        "Delivery is queued. Success is not assumed.",
	}
	if err := s.Store.CreateDelivery(ctx, item); err != nil {
		return
	}
	s.attemptDelivery(ctx, item)
}

func (s *Service) attemptDelivery(ctx context.Context, item contract.WebhookDelivery) contract.WebhookDelivery {
	hook, err := s.Store.GetWebhook(ctx, item.WebhookID)
	if err != nil || hook == nil {
		item.Status = "failed"
		item.Summary = "Webhook is no longer stored."
		_ = s.Store.UpdateDelivery(ctx, item)
		return item
	}
	item.Attempt++
	if errResp := ValidateWebhookURL(hook.URL); errResp != nil {
		item.Status = "blocked"
		item.Summary = "Webhook URL failed SSRF checks. The request was not sent."
		_ = s.Store.UpdateDelivery(ctx, item)
		return item
	}
	if !s.AllowPrivateDelivery {
		if err := resolveWebhookHost(mustHostname(hook.URL)); err != nil {
			item.Status = "blocked"
			item.Summary = "Webhook host resolved to a blocked address. The request was not sent."
			_ = s.Store.UpdateDelivery(ctx, item)
			return item
		}
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, hook.URL, bytes.NewReader([]byte(item.Payload)))
	if err != nil {
		return s.markRetry(ctx, item, "Delivery request could not be created.")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-NetPulse-Signature", item.Signature)
	req.Header.Set("X-NetPulse-Event-Id", item.EventID)
	req.Header.Set("X-NetPulse-Timestamp", item.Timestamp)
	req.Header.Set("X-NetPulse-Idempotency-Key", item.IdempotencyKey)
	res, err := s.client().Do(req)
	if err != nil {
		return s.markRetry(ctx, item, "Delivery did not complete. A retry may be recorded.")
	}
	defer res.Body.Close()
	if res.StatusCode >= 200 && res.StatusCode < 300 {
		item.Status = "delivered"
		item.NextRetryAt = nil
		item.Summary = fmt.Sprintf("Receiver returned %d.", res.StatusCode)
		_ = s.Store.UpdateDelivery(ctx, item)
		return item
	}
	return s.markRetry(ctx, item, fmt.Sprintf("Receiver returned %d.", res.StatusCode))
}

func (s *Service) markRetry(ctx context.Context, item contract.WebhookDelivery, summary string) contract.WebhookDelivery {
	if item.Attempt >= 3 {
		item.Status = "failed"
		item.Summary = summary + " Retry limit reached."
		item.NextRetryAt = nil
	} else {
		item.Status = "retryable"
		next := s.now().UTC().Add(time.Duration(item.Attempt) * time.Minute).Format(time.RFC3339)
		item.NextRetryAt = &next
		item.Summary = summary
	}
	_ = s.Store.UpdateDelivery(ctx, item)
	return item
}

func buildWebhook(workspaceID string, in WebhookInput) (*contract.Webhook, string, *contract.APIError) {
	if errResp := ValidateWebhookURL(in.URL); errResp != nil {
		return nil, "", errResp
	}
	events := []string{}
	seen := map[string]struct{}{}
	for _, raw := range in.Events {
		event := strings.TrimSpace(raw)
		if !KnownWebhookEvent(event) {
			return nil, "", apiErr("validation_error", "Unknown webhook event.")
		}
		if _, ok := seen[event]; ok {
			continue
		}
		seen[event] = struct{}{}
		events = append(events, event)
	}
	if len(events) == 0 {
		return nil, "", apiErr("validation_error", "Select at least one webhook event.")
	}
	raw, hint := auth.NewWebhookSecret()
	return &contract.Webhook{
		ID:         id.New(),
		WorkspaceID: workspaceID,
		URL:        strings.TrimSpace(in.URL),
		Events:     events,
		Secret:     raw,
		SecretHint: hint,
	}, raw, nil
}

func mustHostname(raw string) string {
	parsed, err := url.Parse(raw)
	if err != nil {
		return ""
	}
	return parsed.Hostname()
}

func contains(items []string, want string) bool {
	for _, item := range items {
		if item == want {
			return true
		}
	}
	return false
}
