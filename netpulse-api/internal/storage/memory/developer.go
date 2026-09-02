package memory

import (
	"context"
	"time"

	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/contract"
	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/id"
)

func (s *Store) GetOrCreateWorkspace(_ context.Context, ownerID, name string) (contract.Workspace, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if existingID, ok := s.workspaceOwner[ownerID]; ok {
		return s.workspaces[existingID], nil
	}
	ws := contract.Workspace{
		ID:        id.New(),
		OwnerID:   ownerID,
		Name:      name,
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
	s.workspaces[ws.ID] = ws
	s.workspaceOwner[ownerID] = ws.ID
	s.usage[ws.ID] = contract.EmptyUsage()
	return ws, nil
}

func (s *Store) GetWorkspace(_ context.Context, id string) (*contract.Workspace, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	ws, ok := s.workspaces[id]
	if !ok {
		return nil, nil
	}
	copy := ws
	return &copy, nil
}

func (s *Store) GetWorkspaceByOwner(_ context.Context, ownerID string) (*contract.Workspace, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	id, ok := s.workspaceOwner[ownerID]
	if !ok {
		return nil, nil
	}
	ws := s.workspaces[id]
	copy := ws
	return &copy, nil
}

func (s *Store) CreateMonitor(_ context.Context, item contract.Monitor) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.monitors[item.ID] = item
	return nil
}

func (s *Store) GetMonitor(_ context.Context, id string) (*contract.Monitor, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	item, ok := s.monitors[id]
	if !ok {
		return nil, nil
	}
	copy := item
	return &copy, nil
}

func (s *Store) ListMonitors(_ context.Context, workspaceID string) ([]contract.Monitor, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := []contract.Monitor{}
	for _, item := range s.monitors {
		if item.WorkspaceID == workspaceID {
			out = append(out, item)
		}
	}
	return out, nil
}

func (s *Store) UpdateMonitor(_ context.Context, item contract.Monitor) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.monitors[item.ID] = item
	return nil
}

func (s *Store) DeleteMonitor(_ context.Context, workspaceID, id string) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	item, ok := s.monitors[id]
	if !ok || item.WorkspaceID != workspaceID {
		return false, nil
	}
	delete(s.monitors, id)
	return true, nil
}

func (s *Store) AddCheck(_ context.Context, check contract.MonitorCheck) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.checks = append(s.checks, check)
	return nil
}

func (s *Store) ListChecks(_ context.Context, workspaceID, monitorID string) ([]contract.MonitorCheck, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := []contract.MonitorCheck{}
	for _, check := range s.checks {
		if check.WorkspaceID != workspaceID {
			continue
		}
		if monitorID != "" && check.MonitorID != monitorID {
			continue
		}
		out = append(out, check)
	}
	return out, nil
}

func (s *Store) CreateAPIKey(_ context.Context, key contract.APIKey) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.apiKeys[key.ID] = key
	s.apiKeyHash[key.Hash] = key.ID
	return nil
}

func (s *Store) GetAPIKey(_ context.Context, id string) (*contract.APIKey, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	key, ok := s.apiKeys[id]
	if !ok {
		return nil, nil
	}
	copy := key
	return &copy, nil
}

func (s *Store) GetAPIKeyByHash(_ context.Context, hash string) (*contract.APIKey, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	id, ok := s.apiKeyHash[hash]
	if !ok {
		return nil, nil
	}
	key := s.apiKeys[id]
	copy := key
	return &copy, nil
}

func (s *Store) ListAPIKeys(_ context.Context, workspaceID string) ([]contract.APIKey, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := []contract.APIKey{}
	for _, key := range s.apiKeys {
		if key.WorkspaceID == workspaceID {
			out = append(out, key)
		}
	}
	return out, nil
}

func (s *Store) UpdateAPIKey(_ context.Context, key contract.APIKey) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if existing, ok := s.apiKeys[key.ID]; ok && existing.Hash != key.Hash {
		delete(s.apiKeyHash, existing.Hash)
	}
	s.apiKeys[key.ID] = key
	if key.Hash != "" {
		s.apiKeyHash[key.Hash] = key.ID
	}
	return nil
}

func (s *Store) CreateWebhook(_ context.Context, hook contract.Webhook) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.webhooks[hook.ID] = hook
	return nil
}

func (s *Store) GetWebhook(_ context.Context, id string) (*contract.Webhook, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	hook, ok := s.webhooks[id]
	if !ok {
		return nil, nil
	}
	copy := hook
	return &copy, nil
}

func (s *Store) ListWebhooks(_ context.Context, workspaceID string) ([]contract.Webhook, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := []contract.Webhook{}
	for _, hook := range s.webhooks {
		if hook.WorkspaceID == workspaceID {
			out = append(out, hook)
		}
	}
	return out, nil
}

func (s *Store) UpdateWebhook(_ context.Context, hook contract.Webhook) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.webhooks[hook.ID] = hook
	return nil
}

func (s *Store) DeleteWebhook(_ context.Context, workspaceID, id string) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	hook, ok := s.webhooks[id]
	if !ok || hook.WorkspaceID != workspaceID {
		return false, nil
	}
	delete(s.webhooks, id)
	return true, nil
}

func (s *Store) CreateDelivery(_ context.Context, item contract.WebhookDelivery) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.deliveries[item.ID] = item
	s.deliveryIdem[item.IdempotencyKey] = item.ID
	return nil
}

func (s *Store) GetDeliveryByIdempotency(_ context.Context, key string) (*contract.WebhookDelivery, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	id, ok := s.deliveryIdem[key]
	if !ok {
		return nil, nil
	}
	item := s.deliveries[id]
	copy := item
	return &copy, nil
}

func (s *Store) ListDeliveries(_ context.Context, workspaceID, webhookID string) ([]contract.WebhookDelivery, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := []contract.WebhookDelivery{}
	for _, item := range s.deliveries {
		if item.WorkspaceID == workspaceID && (webhookID == "" || item.WebhookID == webhookID) {
			out = append(out, item)
		}
	}
	return out, nil
}

func (s *Store) UpdateDelivery(_ context.Context, item contract.WebhookDelivery) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.deliveries[item.ID] = item
	return nil
}

func (s *Store) CreateAlertRule(_ context.Context, rule contract.AlertRule) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.alertRules[rule.ID] = rule
	return nil
}

func (s *Store) ListAlertRules(_ context.Context, workspaceID string) ([]contract.AlertRule, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := []contract.AlertRule{}
	for _, rule := range s.alertRules {
		if rule.WorkspaceID == workspaceID {
			out = append(out, rule)
		}
	}
	return out, nil
}

func (s *Store) UpdateAlertRule(_ context.Context, rule contract.AlertRule) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.alertRules[rule.ID] = rule
	return nil
}

func (s *Store) DeleteAlertRule(_ context.Context, workspaceID, id string) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	rule, ok := s.alertRules[id]
	if !ok || rule.WorkspaceID != workspaceID {
		return false, nil
	}
	delete(s.alertRules, id)
	return true, nil
}

func (s *Store) CreateDevIncident(_ context.Context, item contract.DeveloperIncident) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.devIncidents[item.ID] = item
	return nil
}

func (s *Store) ListDevIncidents(_ context.Context, workspaceID string) ([]contract.DeveloperIncident, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := []contract.DeveloperIncident{}
	for _, item := range s.devIncidents {
		if item.WorkspaceID == workspaceID {
			out = append(out, item)
		}
	}
	return out, nil
}

func (s *Store) GetDevIncident(_ context.Context, id string) (*contract.DeveloperIncident, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	item, ok := s.devIncidents[id]
	if !ok {
		return nil, nil
	}
	copy := item
	return &copy, nil
}

func (s *Store) UpdateDevIncident(_ context.Context, item contract.DeveloperIncident) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.devIncidents[item.ID] = item
	return nil
}

func (s *Store) IncrUsage(_ context.Context, workspaceID, field string, delta int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	usage := s.usage[workspaceID]
	switch field {
	case "requests":
		usage.Requests += delta
	case "measurements":
		usage.Measurements += delta
	}
	usage.Summary = contract.EmptyUsage().Summary
	s.usage[workspaceID] = usage
	return nil
}

func (s *Store) GetUsage(_ context.Context, workspaceID string) (contract.Usage, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	usage := s.usage[workspaceID]
	usage.Summary = contract.EmptyUsage().Summary
	return usage, nil
}
