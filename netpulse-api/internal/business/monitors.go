package business

import (
	"context"
	"strings"

	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/auth"
	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/contract"
	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/id"
	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/validation"
)

func (s *Service) ListMonitors(ctx context.Context, actor *Actor, orgID string) ([]contract.Monitor, *contract.APIError, int) {
	if errResp, status := s.requireRead(actor, orgID); errResp != nil {
		return nil, errResp, status
	}
	items, err := s.Store.ListOrgMonitors(ctx, orgID)
	if err != nil {
		return nil, apiErr("unavailable", "Monitor store is unavailable."), 503
	}
	if items == nil {
		items = []contract.Monitor{}
	}
	return items, nil, 200
}

func (s *Service) CreateMonitor(ctx context.Context, actor *Actor, orgID string, in struct {
	Name, Target, Type string
	Regions            []string
	FrequencyS         int
	TimeoutS           int
	DeviceID, NetworkID, ServiceID *string
}) (*contract.Monitor, *contract.APIError, int) {
	if errResp, status := s.require(actor, orgID, PermMonitorCreate); errResp != nil {
		return nil, errResp, status
	}
	name := strings.TrimSpace(in.Name)
	if name == "" {
		return nil, apiErr("validation_error", "Monitor name is required."), 400
	}
	kind := strings.ToLower(strings.TrimSpace(in.Type))
	if kind != "http" && kind != "dns" && kind != "tls" {
		return nil, apiErr("validation_error", "Monitor type must be http, dns, or tls."), 400
	}
	parsed := validation.ParseTarget(in.Target)
	if parsed.Err != nil {
		code := parsed.Code
		if code == "" {
			code = "validation_error"
		}
		status := 400
		if code == "ssrf_blocked" {
			status = 403
		}
		return nil, apiErr(code, parsed.Err.Error()), status
	}
	freq, timeout := in.FrequencyS, in.TimeoutS
	if freq == 0 {
		freq = 300
	}
	if timeout == 0 {
		timeout = 10
	}
	item := contract.Monitor{
		ID: id.New(), OrgID: orgID, Name: name, Target: parsed.Target.Raw, Type: kind,
		Regions: in.Regions, FrequencyS: freq, TimeoutS: timeout, DeviceID: in.DeviceID,
		NetworkID: in.NetworkID, ServiceID: in.ServiceID, Status: "unmeasured",
		Summary: "No checks are stored. Status is not estimated.", CreatedAt: s.stamp(), UpdatedAt: s.stamp(),
	}
	if item.Regions == nil {
		item.Regions = []string{}
	}
	if err := s.Store.CreateOrgMonitor(ctx, item); err != nil {
		return nil, apiErr("unavailable", "Monitor store is unavailable."), 503
	}
	s.audit(ctx, orgID, actor.UserID, "monitor.created", "An organization monitor was stored.")
	return &item, nil, 201
}

func (s *Service) GetMonitor(ctx context.Context, actor *Actor, orgID, monitorID string) (*contract.Monitor, []contract.MonitorCheck, *contract.APIError, int) {
	if errResp, status := s.requireRead(actor, orgID); errResp != nil {
		return nil, nil, errResp, status
	}
	item, err := s.Store.GetOrgMonitor(ctx, monitorID)
	if err != nil {
		return nil, nil, apiErr("unavailable", "Monitor store is unavailable."), 503
	}
	if item == nil || item.OrgID != orgID {
		return nil, nil, apiErr("not_found", "monitor not found"), 404
	}
	checks, err := s.Store.ListOrgChecks(ctx, orgID, monitorID)
	if err != nil {
		return nil, nil, apiErr("unavailable", "Check store is unavailable."), 503
	}
	if checks == nil {
		checks = []contract.MonitorCheck{}
	}
	return item, checks, nil, 200
}

func (s *Service) DeleteMonitor(ctx context.Context, actor *Actor, orgID, monitorID string) *contract.APIError {
	if errResp, _ := s.require(actor, orgID, PermMonitorDelete); errResp != nil {
		return errResp
	}
	ok, err := s.Store.DeleteOrgMonitor(ctx, orgID, monitorID)
	if err != nil {
		return apiErr("unavailable", "Monitor store is unavailable.")
	}
	if !ok {
		return apiErr("not_found", "monitor not found")
	}
	s.audit(ctx, orgID, actor.UserID, "monitor.deleted", "An organization monitor was deleted.")
	return nil
}

func (s *Service) RecordCheck(ctx context.Context, orgID, monitorID, region string, ok bool, latency *int, summary string) (*contract.MonitorCheck, *contract.APIError, int) {
	item, err := s.Store.GetOrgMonitor(ctx, monitorID)
	if err != nil {
		return nil, apiErr("unavailable", "Monitor store is unavailable."), 503
	}
	if item == nil || item.OrgID != orgID {
		return nil, apiErr("not_found", "monitor not found"), 404
	}
	if region == "" {
		region = "unspecified"
	}
	if summary == "" {
		summary = "Worker vantage check stored. This is not a user-path measurement."
	}
	check := contract.MonitorCheck{
		ID: id.New(), MonitorID: monitorID, OrgID: orgID, Region: region, OK: ok, LatencyMs: latency,
		At: s.stamp(), Summary: summary, Endpoint: item.Target,
	}
	if item.DeviceID != nil {
		if device, _ := s.Store.GetOrgDevice(ctx, *item.DeviceID); device != nil && device.OrgID == orgID {
			check.Device = device.Name
		}
	}
	if item.NetworkID != nil {
		if network, _ := s.Store.GetOrgNetwork(ctx, *item.NetworkID); network != nil && network.OrgID == orgID {
			check.Network = network.Name
			check.ASN = network.ASN
		}
	}
	if item.ServiceID != nil {
		if svc, _ := s.Store.GetOrgService(ctx, *item.ServiceID); svc != nil && svc.OrgID == orgID {
			check.Service = svc.Name
		}
	}
	if err := s.Store.AddOrgCheck(ctx, check); err != nil {
		return nil, apiErr("unavailable", "Check store is unavailable."), 503
	}
	item.CheckCount++
	item.UpdatedAt = check.At
	if ok {
		item.Status = "up"
		item.Summary = "Last stored check succeeded from a worker vantage."
	} else {
		item.Status = "down"
		item.Summary = "Last stored check failed from a worker vantage."
	}
	_ = s.Store.UpdateOrgMonitor(ctx, *item)
	s.syncIncident(ctx, *item, check)
	return &check, nil, 201
}

func (s *Service) ListIncidents(ctx context.Context, actor *Actor, orgID string) ([]contract.OrgIncident, *contract.APIError, int) {
	if errResp, status := s.require(actor, orgID, PermIncidentRead); errResp != nil {
		return nil, errResp, status
	}
	items, err := s.Store.ListOrgIncidents(ctx, orgID)
	if err != nil {
		return nil, apiErr("unavailable", "Incident store is unavailable."), 503
	}
	if items == nil {
		items = []contract.OrgIncident{}
	}
	return items, nil, 200
}

func (s *Service) GetIncident(ctx context.Context, actor *Actor, orgID, incidentID string) (*contract.OrgIncident, *contract.APIError, int) {
	if errResp, status := s.require(actor, orgID, PermIncidentRead); errResp != nil {
		return nil, errResp, status
	}
	item, err := s.Store.GetOrgIncident(ctx, incidentID)
	if err != nil {
		return nil, apiErr("unavailable", "Incident store is unavailable."), 503
	}
	if item == nil || item.OrgID != orgID {
		return nil, apiErr("not_found", "incident not found"), 404
	}
	return item, nil, 200
}

func (s *Service) UpdateIncident(ctx context.Context, actor *Actor, orgID, incidentID, status string) (*contract.OrgIncident, *contract.APIError, int) {
	if errResp, st := s.require(actor, orgID, PermIncidentManage); errResp != nil {
		return nil, errResp, st
	}
	item, errResp, code := s.GetIncident(ctx, actor, orgID, incidentID)
	if errResp != nil {
		return nil, errResp, code
	}
	status = strings.ToLower(strings.TrimSpace(status))
	if status != "open" && status != "resolved" {
		return nil, apiErr("validation_error", "Incident status must be open or resolved."), 400
	}
	item.Status = status
	if status == "resolved" {
		now := s.stamp()
		item.ResolvedAt = &now
		item.Summary = "Marked resolved from a stored organization action. This is not calendar SLA proof."
	}
	if err := s.Store.UpdateOrgIncident(ctx, *item); err != nil {
		return nil, apiErr("unavailable", "Incident store is unavailable."), 503
	}
	return item, nil, 200
}

func (s *Service) CreateDiagnosis(ctx context.Context, actor *Actor, orgID, target string) (*contract.OrgDiagnosis, *contract.APIError, int) {
	if errResp, status := s.require(actor, orgID, PermDiagnosisCreate); errResp != nil {
		return nil, errResp, status
	}
	if s.Diagnoses == nil {
		return nil, apiErr("unavailable", "Diagnosis service is unavailable."), 503
	}
	diag, errResp, status := s.Diagnoses.Create(ctx, target, actor.UserID)
	if errResp != nil {
		return nil, errResp, status
	}
	item := contract.OrgDiagnosis{
		ID: id.New(), OrgID: orgID, DiagnosisID: diag.ID, Target: target, Status: diag.Status,
		CreatedAt: diag.Created, Summary: "Organization-linked diagnosis. Worker vantage only.",
	}
	if err := s.Store.CreateOrgDiagnosis(ctx, item); err != nil {
		return nil, apiErr("unavailable", "Diagnosis store is unavailable."), 503
	}
	return &item, nil, 201
}

func (s *Service) ListDiagnoses(ctx context.Context, actor *Actor, orgID string) ([]contract.OrgDiagnosis, *contract.APIError, int) {
	if errResp, status := s.require(actor, orgID, PermDiagnosisRead); errResp != nil {
		return nil, errResp, status
	}
	items, err := s.Store.ListOrgDiagnoses(ctx, orgID)
	if err != nil {
		return nil, apiErr("unavailable", "Diagnosis store is unavailable."), 503
	}
	if items == nil {
		items = []contract.OrgDiagnosis{}
	}
	return items, nil, 200
}

func (s *Service) ListKeys(ctx context.Context, actor *Actor, orgID string) ([]contract.APIKey, *contract.APIError, int) {
	if errResp, status := s.require(actor, orgID, PermAPIManage); errResp != nil {
		return nil, errResp, status
	}
	if !actor.SessionOnly {
		return nil, apiErr("forbidden", "API keys cannot manage API keys."), 403
	}
	items, err := s.Store.ListOrgKeys(ctx, orgID)
	if err != nil {
		return nil, apiErr("unavailable", "Key store is unavailable."), 503
	}
	if items == nil {
		items = []contract.APIKey{}
	}
	return items, nil, 200
}

func (s *Service) CreateKey(ctx context.Context, actor *Actor, orgID, name string, scopes []string, limit int) (*contract.APIKey, string, *contract.APIError, int) {
	if errResp, status := s.require(actor, orgID, PermAPIManage); errResp != nil {
		return nil, "", errResp, status
	}
	if !actor.SessionOnly {
		return nil, "", apiErr("forbidden", "API keys cannot manage API keys."), 403
	}
	if strings.TrimSpace(name) == "" {
		name = "Organization key"
	}
	perms := NormalizePerms(scopes)
	if len(scopes) > 0 && len(perms) == 0 {
		return nil, "", apiErr("validation_error", "No recognized permissions were supplied."), 400
	}
	if limit == 0 {
		limit = 60
	}
	raw, prefix, last4, hash := auth.NewOrgAPIKeySecret()
	key := contract.APIKey{
		ID: id.New(), OrgID: orgID, Name: name, Prefix: prefix, Last4: last4, Hash: hash,
		Scopes: perms, RateLimitPerMin: limit, CreatedAt: s.stamp(),
	}
	if err := s.Store.CreateOrgKey(ctx, key); err != nil {
		return nil, "", apiErr("unavailable", "Key store is unavailable."), 503
	}
	s.audit(ctx, orgID, actor.UserID, "key.created", "An organization API key was created.")
	return &key, raw, nil, 201
}

func (s *Service) RotateKey(ctx context.Context, actor *Actor, orgID, keyID string) (*contract.APIKey, string, *contract.APIError, int) {
	if errResp, status := s.require(actor, orgID, PermAPIManage); errResp != nil {
		return nil, "", errResp, status
	}
	if !actor.SessionOnly {
		return nil, "", apiErr("forbidden", "API keys cannot manage API keys."), 403
	}
	key, err := s.Store.GetOrgKey(ctx, keyID)
	if err != nil {
		return nil, "", apiErr("unavailable", "Key store is unavailable."), 503
	}
	if key == nil || key.OrgID != orgID || key.Revoked {
		return nil, "", apiErr("not_found", "api key not found"), 404
	}
	raw, prefix, last4, hash := auth.NewOrgAPIKeySecret()
	key.Prefix, key.Last4, key.Hash = prefix, last4, hash
	if err := s.Store.UpdateOrgKey(ctx, *key); err != nil {
		return nil, "", apiErr("unavailable", "Key store is unavailable."), 503
	}
	return key, raw, nil, 200
}

func (s *Service) RevokeKey(ctx context.Context, actor *Actor, orgID, keyID string) (*contract.APIKey, *contract.APIError, int) {
	if errResp, status := s.require(actor, orgID, PermAPIManage); errResp != nil {
		return nil, errResp, status
	}
	if !actor.SessionOnly {
		return nil, apiErr("forbidden", "API keys cannot manage API keys."), 403
	}
	key, err := s.Store.GetOrgKey(ctx, keyID)
	if err != nil {
		return nil, apiErr("unavailable", "Key store is unavailable."), 503
	}
	if key == nil || key.OrgID != orgID {
		return nil, apiErr("not_found", "api key not found"), 404
	}
	key.Revoked = true
	if err := s.Store.UpdateOrgKey(ctx, *key); err != nil {
		return nil, apiErr("unavailable", "Key store is unavailable."), 503
	}
	return key, nil, 200
}

func (s *Service) Billing(ctx context.Context, actor *Actor, orgID string) (*contract.Billing, *contract.APIError, int) {
	if errResp, status := s.require(actor, orgID, PermBillingManage); errResp != nil {
		return nil, errResp, status
	}
	billing := contract.EmptyOrgBilling(orgID)
	return &billing, nil, 200
}

func (s *Service) ListAudit(ctx context.Context, actor *Actor, orgID string) ([]contract.AuditEvent, *contract.APIError, int) {
	if errResp, status := s.require(actor, orgID, PermAuditRead); errResp != nil {
		return nil, errResp, status
	}
	items, err := s.Store.ListAudit(ctx, orgID)
	if err != nil {
		return nil, apiErr("unavailable", "Audit store is unavailable."), 503
	}
	if items == nil {
		items = []contract.AuditEvent{}
	}
	return items, nil, 200
}

func (s *Service) syncIncident(ctx context.Context, monitor contract.Monitor, check contract.MonitorCheck) {
	incidents, err := s.Store.ListOrgIncidents(ctx, monitor.OrgID)
	if err != nil {
		return
	}
	var open *contract.OrgIncident
	for i := range incidents {
		if incidents[i].MonitorID == monitor.ID && incidents[i].Status != "resolved" {
			copy := incidents[i]
			open = &copy
			break
		}
	}
	devices, networks, services := []string{}, []string{}, []string{}
	if monitor.DeviceID != nil {
		devices = []string{*monitor.DeviceID}
	}
	if monitor.NetworkID != nil {
		networks = []string{*monitor.NetworkID}
	}
	if monitor.ServiceID != nil {
		services = []string{*monitor.ServiceID}
	}
	if !check.OK && open == nil {
		item := contract.OrgIncident{
			ID: id.New(), OrgID: monitor.OrgID, MonitorID: monitor.ID,
			Title: monitor.Name + " is down from a stored worker check", Status: "open",
			StartedAt: check.At, DeviceIDs: devices, NetworkIDs: networks, ServiceIDs: services,
			Regions: []string{check.Region}, SampleCount: 1,
			Summary: "Opened from a stored failed check. This is not a global outage claim.",
		}
		_ = s.Store.CreateOrgIncident(ctx, item)
		return
	}
	if !check.OK && open != nil {
		open.SampleCount++
		_ = s.Store.UpdateOrgIncident(ctx, *open)
		return
	}
	if check.OK && open != nil {
		open.Status = "resolved"
		open.ResolvedAt = &check.At
		open.SampleCount++
		open.Summary = "Resolved after a stored successful check."
		_ = s.Store.UpdateOrgIncident(ctx, *open)
	}
}
