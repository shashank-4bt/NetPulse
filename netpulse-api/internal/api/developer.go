package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/contract"
	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/developer"
)

type devActor struct {
	Workspace   contract.Workspace
	Scopes      []string
	SessionOnly bool
}

func (s *Server) requireDev(w http.ResponseWriter, r *http.Request, scope string) *devActor {
	if s.Developer == nil || s.Accounts == nil {
		write(w, http.StatusServiceUnavailable, contract.Envelope{Error: &contract.APIError{Code: "unavailable", Message: "Developer platform is unavailable."}})
		return nil
	}
	if token := SessionToken(r); token != "" {
		_, user, errResp, status := s.Accounts.Require(r.Context(), token)
		if errResp != nil {
			write(w, status, contract.Envelope{Error: errResp})
			return nil
		}
		ws, errResp, status := s.Developer.WorkspaceForOwner(r.Context(), user.ID, user.DisplayName)
		if errResp != nil {
			write(w, status, contract.Envelope{Error: errResp})
			return nil
		}
		return &devActor{Workspace: ws, Scopes: developer.SessionScopes(), SessionOnly: true}
	}
	raw := APIKeyToken(r)
	if raw == "" {
		write(w, http.StatusUnauthorized, contract.Envelope{Error: &contract.APIError{Code: "unauthorized", Message: "Sign in or supply an API key."}})
		return nil
	}
	key, ws, errResp, status := s.Developer.LookupAPIKey(r.Context(), raw)
	if errResp != nil {
		write(w, status, contract.Envelope{Error: errResp})
		return nil
	}
	limit := key.RateLimitPerMin
	if limit < 1 {
		limit = 60
	}
	if s.Limiter != nil && !s.Limiter.Allow("key:"+key.ID, limit) {
		s.recordAbuse(r, "api_abuse", "key:"+key.ID, "dev", "blocked", "API key rate limit exceeded.")
		write(w, http.StatusTooManyRequests, contract.Envelope{Error: &contract.APIError{Code: "rate_limited", Message: "API key rate limit exceeded."}})
		return nil
	}
	if scope != "" && !developer.HasScope(key.Scopes, scope) {
		write(w, http.StatusForbidden, contract.Envelope{Error: &contract.APIError{Code: "forbidden", Message: "API key is missing the required scope."}})
		return nil
	}
	return &devActor{Workspace: *ws, Scopes: key.Scopes, SessionOnly: false}
}

func (s *Server) requireSessionDev(w http.ResponseWriter, r *http.Request) *devActor {
	actor := s.requireDev(w, r, "")
	if actor == nil {
		return nil
	}
	if !actor.SessionOnly {
		write(w, http.StatusForbidden, contract.Envelope{Error: &contract.APIError{Code: "forbidden", Message: "API keys cannot manage API keys."}})
		return nil
	}
	return actor
}

func APIKeyToken(r *http.Request) string {
	header := r.Header.Get("Authorization")
	if strings.HasPrefix(strings.ToLower(header), "bearer ") {
		return strings.TrimSpace(header[7:])
	}
	return strings.TrimSpace(r.Header.Get("X-NetPulse-Key"))
}

func (s *Server) devWorkspace(w http.ResponseWriter, r *http.Request) {
	actor := s.requireDev(w, r, "")
	if actor == nil {
		return
	}
	write(w, http.StatusOK, contract.Envelope{OK: true, Workspace: &actor.Workspace})
}

func (s *Server) getDevWorkspace(w http.ResponseWriter, r *http.Request) {
	actor := s.requireDev(w, r, "")
	if actor == nil {
		return
	}
	ws, errResp, status := s.Developer.GetWorkspace(r.Context(), actor.Workspace.ID, r.PathValue("id"))
	if errResp != nil {
		write(w, status, contract.Envelope{Error: errResp})
		return
	}
	write(w, http.StatusOK, contract.Envelope{OK: true, Workspace: ws})
}

func (s *Server) devDashboard(w http.ResponseWriter, r *http.Request) {
	actor := s.requireDev(w, r, developer.ScopeDashboardRead)
	if actor == nil {
		return
	}
	dash, errResp, status := s.Developer.Dashboard(r.Context(), actor.Workspace.ID)
	if errResp != nil {
		write(w, status, contract.Envelope{Error: errResp})
		return
	}
	write(w, status, contract.Envelope{OK: true, DevDashboard: &dash})
}

func (s *Server) listDevMonitors(w http.ResponseWriter, r *http.Request) {
	actor := s.requireDev(w, r, developer.ScopeMonitorsRead)
	if actor == nil {
		return
	}
	items, errResp, status := s.Developer.ListMonitors(r.Context(), actor.Workspace.ID)
	if errResp != nil {
		write(w, status, contract.Envelope{Error: errResp})
		return
	}
	write(w, status, contract.Envelope{OK: true, Monitors: items})
}

func (s *Server) createDevMonitor(w http.ResponseWriter, r *http.Request) {
	actor := s.requireDev(w, r, developer.ScopeMonitorsWrite)
	if actor == nil {
		return
	}
	in, ok := decodeMonitor(w, r)
	if !ok {
		return
	}
	item, errResp, status := s.Developer.CreateMonitor(r.Context(), actor.Workspace.ID, in)
	if errResp != nil {
		write(w, status, contract.Envelope{Error: errResp})
		return
	}
	write(w, status, contract.Envelope{OK: true, Monitor: item})
}

func (s *Server) getDevMonitor(w http.ResponseWriter, r *http.Request) {
	actor := s.requireDev(w, r, developer.ScopeMonitorsRead)
	if actor == nil {
		return
	}
	item, checks, errResp, status := s.Developer.GetMonitor(r.Context(), actor.Workspace.ID, r.PathValue("id"))
	if errResp != nil {
		write(w, status, contract.Envelope{Error: errResp})
		return
	}
	write(w, status, contract.Envelope{OK: true, Monitor: item, Checks: checks})
}

func (s *Server) patchDevMonitor(w http.ResponseWriter, r *http.Request) {
	actor := s.requireDev(w, r, developer.ScopeMonitorsWrite)
	if actor == nil {
		return
	}
	in, ok := decodeMonitor(w, r)
	if !ok {
		return
	}
	item, errResp, status := s.Developer.UpdateMonitor(r.Context(), actor.Workspace.ID, r.PathValue("id"), in)
	if errResp != nil {
		write(w, status, contract.Envelope{Error: errResp})
		return
	}
	write(w, status, contract.Envelope{OK: true, Monitor: item})
}

func (s *Server) deleteDevMonitor(w http.ResponseWriter, r *http.Request) {
	actor := s.requireDev(w, r, developer.ScopeMonitorsWrite)
	if actor == nil {
		return
	}
	if errResp := s.Developer.DeleteMonitor(r.Context(), actor.Workspace.ID, r.PathValue("id")); errResp != nil {
		write(w, statusForAPI(errResp), contract.Envelope{Error: errResp})
		return
	}
	write(w, http.StatusOK, contract.Envelope{OK: true})
}

func (s *Server) runDevMonitor(w http.ResponseWriter, r *http.Request) {
	actor := s.requireDev(w, r, developer.ScopeMonitorsWrite)
	if actor == nil {
		return
	}
	var body struct {
		Region string `json:"region"`
	}
	_ = json.NewDecoder(http.MaxBytesReader(w, r.Body, 8<<10)).Decode(&body)
	check, errResp, status := s.Developer.RunMonitor(r.Context(), actor.Workspace.ID, r.PathValue("id"), body.Region)
	if errResp != nil {
		write(w, status, contract.Envelope{Error: errResp})
		return
	}
	write(w, status, contract.Envelope{OK: true, Checks: []contract.MonitorCheck{*check}})
}

func (s *Server) listDevIncidents(w http.ResponseWriter, r *http.Request) {
	actor := s.requireDev(w, r, developer.ScopeIncidentsRead)
	if actor == nil {
		return
	}
	items, errResp, status := s.Developer.ListIncidents(r.Context(), actor.Workspace.ID)
	if errResp != nil {
		write(w, status, contract.Envelope{Error: errResp})
		return
	}
	write(w, status, contract.Envelope{OK: true, DevIncidents: items})
}

func (s *Server) getDevIncident(w http.ResponseWriter, r *http.Request) {
	actor := s.requireDev(w, r, developer.ScopeIncidentsRead)
	if actor == nil {
		return
	}
	item, errResp, status := s.Developer.GetIncident(r.Context(), actor.Workspace.ID, r.PathValue("id"))
	if errResp != nil {
		write(w, status, contract.Envelope{Error: errResp})
		return
	}
	write(w, status, contract.Envelope{OK: true, DevIncident: item})
}

func (s *Server) listDevKeys(w http.ResponseWriter, r *http.Request) {
	actor := s.requireSessionDev(w, r)
	if actor == nil {
		return
	}
	items, errResp, status := s.Developer.ListKeys(r.Context(), actor.Workspace.ID)
	if errResp != nil {
		write(w, status, contract.Envelope{Error: errResp})
		return
	}
	write(w, status, contract.Envelope{OK: true, APIKeys: items})
}

func (s *Server) createDevKey(w http.ResponseWriter, r *http.Request) {
	actor := s.requireSessionDev(w, r)
	if actor == nil {
		return
	}
	var body struct {
		Name            string   `json:"name"`
		Scopes          []string `json:"scopes"`
		RateLimitPerMin int      `json:"rateLimitPerMin"`
	}
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 8<<10)).Decode(&body); err != nil {
		write(w, http.StatusBadRequest, contract.Envelope{Error: &contract.APIError{Code: "validation_error", Message: "Expected JSON."}})
		return
	}
	item, secret, errResp, status := s.Developer.CreateKey(r.Context(), actor.Workspace.ID, developer.APIKeyInput{
		Name: body.Name, Scopes: body.Scopes, RateLimitPerMin: body.RateLimitPerMin,
	})
	if errResp != nil {
		write(w, status, contract.Envelope{Error: errResp})
		return
	}
	write(w, status, contract.Envelope{OK: true, APIKey: item, KeySecret: secret})
}

func (s *Server) rotateDevKey(w http.ResponseWriter, r *http.Request) {
	actor := s.requireSessionDev(w, r)
	if actor == nil {
		return
	}
	item, secret, errResp, status := s.Developer.RotateKey(r.Context(), actor.Workspace.ID, r.PathValue("id"))
	if errResp != nil {
		write(w, status, contract.Envelope{Error: errResp})
		return
	}
	write(w, status, contract.Envelope{OK: true, APIKey: item, KeySecret: secret})
}

func (s *Server) revokeDevKey(w http.ResponseWriter, r *http.Request) {
	actor := s.requireSessionDev(w, r)
	if actor == nil {
		return
	}
	item, errResp, status := s.Developer.RevokeKey(r.Context(), actor.Workspace.ID, r.PathValue("id"))
	if errResp != nil {
		write(w, status, contract.Envelope{Error: errResp})
		return
	}
	write(w, status, contract.Envelope{OK: true, APIKey: item})
}

func (s *Server) listDevWebhooks(w http.ResponseWriter, r *http.Request) {
	actor := s.requireDev(w, r, developer.ScopeWebhooksRead)
	if actor == nil {
		return
	}
	items, errResp, status := s.Developer.ListWebhooks(r.Context(), actor.Workspace.ID)
	if errResp != nil {
		write(w, status, contract.Envelope{Error: errResp})
		return
	}
	write(w, status, contract.Envelope{OK: true, Webhooks: items})
}

func (s *Server) createDevWebhook(w http.ResponseWriter, r *http.Request) {
	actor := s.requireDev(w, r, developer.ScopeWebhooksWrite)
	if actor == nil {
		return
	}
	var body struct {
		URL    string   `json:"url"`
		Events []string `json:"events"`
	}
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 8<<10)).Decode(&body); err != nil {
		write(w, http.StatusBadRequest, contract.Envelope{Error: &contract.APIError{Code: "validation_error", Message: "Expected JSON."}})
		return
	}
	item, secret, errResp, status := s.Developer.CreateWebhook(r.Context(), actor.Workspace.ID, developer.WebhookInput{URL: body.URL, Events: body.Events})
	if errResp != nil {
		if errResp.Code == "ssrf_blocked" {
			s.recordAbuse(r, "ssrf", actor.Workspace.ID, "webhooks", "blocked", "Webhook URL was blocked by SSRF policy.")
		}
		write(w, status, contract.Envelope{Error: errResp})
		return
	}
	write(w, status, contract.Envelope{OK: true, Webhook: item, WebhookSecret: secret})
}

func (s *Server) rotateDevWebhook(w http.ResponseWriter, r *http.Request) {
	actor := s.requireDev(w, r, developer.ScopeWebhooksWrite)
	if actor == nil {
		return
	}
	item, secret, errResp, status := s.Developer.RotateWebhook(r.Context(), actor.Workspace.ID, r.PathValue("id"))
	if errResp != nil {
		write(w, status, contract.Envelope{Error: errResp})
		return
	}
	write(w, status, contract.Envelope{OK: true, Webhook: item, WebhookSecret: secret})
}

func (s *Server) deleteDevWebhook(w http.ResponseWriter, r *http.Request) {
	actor := s.requireDev(w, r, developer.ScopeWebhooksWrite)
	if actor == nil {
		return
	}
	if errResp := s.Developer.DeleteWebhook(r.Context(), actor.Workspace.ID, r.PathValue("id")); errResp != nil {
		write(w, statusForAPI(errResp), contract.Envelope{Error: errResp})
		return
	}
	write(w, http.StatusOK, contract.Envelope{OK: true})
}

func (s *Server) listDevDeliveries(w http.ResponseWriter, r *http.Request) {
	actor := s.requireDev(w, r, developer.ScopeWebhooksRead)
	if actor == nil {
		return
	}
	items, errResp, status := s.Developer.ListDeliveries(r.Context(), actor.Workspace.ID, r.PathValue("id"))
	if errResp != nil {
		write(w, status, contract.Envelope{Error: errResp})
		return
	}
	write(w, status, contract.Envelope{OK: true, Deliveries: items})
}

func (s *Server) retryDevDeliveries(w http.ResponseWriter, r *http.Request) {
	actor := s.requireDev(w, r, developer.ScopeWebhooksWrite)
	if actor == nil {
		return
	}
	items, errResp, status := s.Developer.RetryDeliveries(r.Context(), actor.Workspace.ID, r.PathValue("id"))
	if errResp != nil {
		write(w, status, contract.Envelope{Error: errResp})
		return
	}
	write(w, status, contract.Envelope{OK: true, Deliveries: items})
}

func (s *Server) listDevAlerts(w http.ResponseWriter, r *http.Request) {
	actor := s.requireDev(w, r, developer.ScopeAlertsRead)
	if actor == nil {
		return
	}
	items, errResp, status := s.Developer.ListAlertRules(r.Context(), actor.Workspace.ID)
	if errResp != nil {
		write(w, status, contract.Envelope{Error: errResp})
		return
	}
	write(w, status, contract.Envelope{OK: true, AlertRules: items})
}

func (s *Server) createDevAlert(w http.ResponseWriter, r *http.Request) {
	actor := s.requireDev(w, r, developer.ScopeAlertsWrite)
	if actor == nil {
		return
	}
	in, ok := decodeAlert(w, r)
	if !ok {
		return
	}
	item, errResp, status := s.Developer.CreateAlertRule(r.Context(), actor.Workspace.ID, in)
	if errResp != nil {
		write(w, status, contract.Envelope{Error: errResp})
		return
	}
	write(w, status, contract.Envelope{OK: true, AlertRule: item})
}

func (s *Server) putDevAlert(w http.ResponseWriter, r *http.Request) {
	actor := s.requireDev(w, r, developer.ScopeAlertsWrite)
	if actor == nil {
		return
	}
	in, ok := decodeAlert(w, r)
	if !ok {
		return
	}
	item, errResp, status := s.Developer.UpdateAlertRule(r.Context(), actor.Workspace.ID, r.PathValue("id"), in)
	if errResp != nil {
		write(w, status, contract.Envelope{Error: errResp})
		return
	}
	write(w, status, contract.Envelope{OK: true, AlertRule: item})
}

func (s *Server) deleteDevAlert(w http.ResponseWriter, r *http.Request) {
	actor := s.requireDev(w, r, developer.ScopeAlertsWrite)
	if actor == nil {
		return
	}
	if errResp := s.Developer.DeleteAlertRule(r.Context(), actor.Workspace.ID, r.PathValue("id")); errResp != nil {
		write(w, statusForAPI(errResp), contract.Envelope{Error: errResp})
		return
	}
	write(w, http.StatusOK, contract.Envelope{OK: true})
}

func (s *Server) devUsage(w http.ResponseWriter, r *http.Request) {
	actor := s.requireDev(w, r, developer.ScopeUsageRead)
	if actor == nil {
		return
	}
	item, errResp, status := s.Developer.Usage(r.Context(), actor.Workspace.ID)
	if errResp != nil {
		write(w, status, contract.Envelope{Error: errResp})
		return
	}
	write(w, status, contract.Envelope{OK: true, Usage: &item})
}

func (s *Server) devSLA(w http.ResponseWriter, r *http.Request) {
	actor := s.requireDev(w, r, developer.ScopeSLARead)
	if actor == nil {
		return
	}
	item, errResp, status := s.Developer.SLA(r.Context(), actor.Workspace.ID)
	if errResp != nil {
		write(w, status, contract.Envelope{Error: errResp})
		return
	}
	write(w, status, contract.Envelope{OK: true, SLA: &item})
}

func decodeMonitor(w http.ResponseWriter, r *http.Request) (developer.MonitorInput, bool) {
	var body struct {
		Name              string   `json:"name"`
		Target            string   `json:"target"`
		Type              string   `json:"type"`
		Regions           []string `json:"regions"`
		FrequencySeconds  int      `json:"frequencySeconds"`
		TimeoutSeconds    int      `json:"timeoutSeconds"`
		AvailabilityBelow *float64 `json:"availabilityBelow"`
		LatencyMsAbove    *int     `json:"latencyMsAbove"`
		ErrorRateAbove    *float64 `json:"errorRateAbove"`
		Thresholds        *struct {
			AvailabilityBelow *float64 `json:"availabilityBelow"`
			LatencyMsAbove    *int     `json:"latencyMsAbove"`
			ErrorRateAbove    *float64 `json:"errorRateAbove"`
		} `json:"thresholds"`
	}
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 16<<10)).Decode(&body); err != nil {
		write(w, http.StatusBadRequest, contract.Envelope{Error: &contract.APIError{Code: "validation_error", Message: "Expected JSON."}})
		return developer.MonitorInput{}, false
	}
	in := developer.MonitorInput{
		Name:       body.Name,
		Target:     body.Target,
		Type:       body.Type,
		Regions:    body.Regions,
		FrequencyS: body.FrequencySeconds,
		TimeoutS:   body.TimeoutSeconds,
		Thresholds: contract.MonitorThresholds{
			AvailabilityBelow: body.AvailabilityBelow,
			LatencyMsAbove:    body.LatencyMsAbove,
			ErrorRateAbove:    body.ErrorRateAbove,
		},
	}
	if body.Thresholds != nil {
		in.Thresholds = contract.MonitorThresholds{
			AvailabilityBelow: body.Thresholds.AvailabilityBelow,
			LatencyMsAbove:    body.Thresholds.LatencyMsAbove,
			ErrorRateAbove:    body.Thresholds.ErrorRateAbove,
		}
	}
	return in, true
}

func decodeAlert(w http.ResponseWriter, r *http.Request) (developer.AlertInput, bool) {
	var body struct {
		Kind      string   `json:"kind"`
		MonitorID *string  `json:"monitorId"`
		Threshold float64  `json:"threshold"`
		Enabled   *bool    `json:"enabled"`
	}
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 8<<10)).Decode(&body); err != nil {
		write(w, http.StatusBadRequest, contract.Envelope{Error: &contract.APIError{Code: "validation_error", Message: "Expected JSON."}})
		return developer.AlertInput{}, false
	}
	return developer.AlertInput{Kind: body.Kind, MonitorID: body.MonitorID, Threshold: body.Threshold, Enabled: body.Enabled}, true
}

func statusForAPI(err *contract.APIError) int {
	if err == nil {
		return http.StatusOK
	}
	switch err.Code {
	case "not_found":
		return http.StatusNotFound
	case "unauthorized":
		return http.StatusUnauthorized
	case "forbidden":
		return http.StatusForbidden
	case "ssrf_blocked":
		return http.StatusForbidden
	case "unavailable":
		return http.StatusServiceUnavailable
	default:
		return http.StatusBadRequest
	}
}
