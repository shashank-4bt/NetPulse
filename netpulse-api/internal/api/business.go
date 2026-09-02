package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/business"
	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/contract"
)

func (s *Server) requireOrg(w http.ResponseWriter, r *http.Request, orgID string) *business.Actor {
	if s.Business == nil || s.Accounts == nil {
		write(w, http.StatusServiceUnavailable, contract.Envelope{Error: &contract.APIError{Code: "unavailable", Message: "Business platform is unavailable."}})
		return nil
	}
	if token := SessionToken(r); token != "" {
		_, user, errResp, status := s.Accounts.Require(r.Context(), token)
		if errResp != nil {
			write(w, status, contract.Envelope{Error: errResp})
			return nil
		}
		s.Business.ConsumeInvites(r.Context(), user.ID, user.Email)
		member, org, errResp, status := s.Business.MemberFor(r.Context(), orgID, user.ID)
		if errResp != nil {
			write(w, status, contract.Envelope{Error: errResp})
			return nil
		}
		return &business.Actor{
			UserID: user.ID, Email: user.Email, DisplayName: user.DisplayName,
			Org: *org, Member: member, Perms: member.Permissions, SessionOnly: true,
		}
	}
	raw := APIKeyToken(r)
	if raw == "" {
		write(w, http.StatusUnauthorized, contract.Envelope{Error: &contract.APIError{Code: "unauthorized", Message: "Sign in or supply an organization API key."}})
		return nil
	}
	key, org, errResp, status := s.Business.LookupOrgKey(r.Context(), raw)
	if errResp != nil {
		write(w, status, contract.Envelope{Error: errResp})
		return nil
	}
	if org.ID != orgID {
		write(w, http.StatusNotFound, contract.Envelope{Error: &contract.APIError{Code: "not_found", Message: "not found"}})
		return nil
	}
	limit := key.RateLimitPerMin
	if limit < 1 {
		limit = 60
	}
	if s.Limiter != nil && !s.Limiter.Allow("orgkey:"+key.ID, limit) {
		write(w, http.StatusTooManyRequests, contract.Envelope{Error: &contract.APIError{Code: "rate_limited", Message: "API key rate limit exceeded."}})
		return nil
	}
	return &business.Actor{UserID: "key:" + key.ID, Org: *org, Key: key, Perms: key.Scopes, SessionOnly: false}
}

func (s *Server) listOrgs(w http.ResponseWriter, r *http.Request) {
	if s.Business == nil || s.Accounts == nil {
		write(w, http.StatusServiceUnavailable, contract.Envelope{Error: &contract.APIError{Code: "unavailable", Message: "Business platform is unavailable."}})
		return
	}
	_, user, errResp, status := s.Accounts.Require(r.Context(), SessionToken(r))
	if errResp != nil {
		write(w, status, contract.Envelope{Error: errResp})
		return
	}
	items, errResp, status := s.Business.ListOrgs(r.Context(), user.ID, user.Email)
	if errResp != nil {
		write(w, status, contract.Envelope{Error: errResp})
		return
	}
	write(w, status, contract.Envelope{OK: true, Organizations: items})
}

func (s *Server) createOrg(w http.ResponseWriter, r *http.Request) {
	if s.Business == nil || s.Accounts == nil {
		write(w, http.StatusServiceUnavailable, contract.Envelope{Error: &contract.APIError{Code: "unavailable", Message: "Business platform is unavailable."}})
		return
	}
	_, user, errResp, status := s.Accounts.Require(r.Context(), SessionToken(r))
	if errResp != nil {
		write(w, status, contract.Envelope{Error: errResp})
		return
	}
	var body struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 8<<10)).Decode(&body); err != nil {
		write(w, http.StatusBadRequest, contract.Envelope{Error: &contract.APIError{Code: "validation_error", Message: "Expected JSON."}})
		return
	}
	item, errResp, status := s.Business.CreateOrg(r.Context(), user.ID, user.Email, user.DisplayName, body.Name)
	if errResp != nil {
		write(w, status, contract.Envelope{Error: errResp})
		return
	}
	write(w, status, contract.Envelope{OK: true, Organization: item})
}

func (s *Server) getOrg(w http.ResponseWriter, r *http.Request) {
	actor := s.requireOrg(w, r, r.PathValue("orgId"))
	if actor == nil {
		return
	}
	item, errResp, status := s.Business.GetOrg(r.Context(), actor, r.PathValue("orgId"))
	if errResp != nil {
		write(w, status, contract.Envelope{Error: errResp})
		return
	}
	write(w, status, contract.Envelope{OK: true, Organization: item, Permissions: actor.Perms})
}

func (s *Server) patchOrg(w http.ResponseWriter, r *http.Request) {
	actor := s.requireOrg(w, r, r.PathValue("orgId"))
	if actor == nil {
		return
	}
	var body struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 8<<10)).Decode(&body); err != nil {
		write(w, http.StatusBadRequest, contract.Envelope{Error: &contract.APIError{Code: "validation_error", Message: "Expected JSON."}})
		return
	}
	item, errResp, status := s.Business.UpdateOrg(r.Context(), actor, r.PathValue("orgId"), body.Name)
	if errResp != nil {
		write(w, status, contract.Envelope{Error: errResp})
		return
	}
	write(w, status, contract.Envelope{OK: true, Organization: item})
}

func (s *Server) orgDashboard(w http.ResponseWriter, r *http.Request) {
	actor := s.requireOrg(w, r, r.PathValue("orgId"))
	if actor == nil {
		return
	}
	item, errResp, status := s.Business.Dashboard(r.Context(), actor, r.PathValue("orgId"))
	if errResp != nil {
		write(w, status, contract.Envelope{Error: errResp})
		return
	}
	write(w, status, contract.Envelope{OK: true, OrgDashboard: &item})
}

func (s *Server) orgAnalytics(w http.ResponseWriter, r *http.Request) {
	actor := s.requireOrg(w, r, r.PathValue("orgId"))
	if actor == nil {
		return
	}
	q := r.URL.Query()
	filters := map[string]string{}
	for _, key := range []string{"region", "network", "asn", "service", "endpoint", "device"} {
		if value := strings.TrimSpace(q.Get(key)); value != "" {
			filters[key] = value
		}
	}
	item, errResp, status := s.Business.Analytics(r.Context(), actor, r.PathValue("orgId"), filters)
	if errResp != nil {
		write(w, status, contract.Envelope{Error: errResp})
		return
	}
	write(w, status, contract.Envelope{OK: true, Analytics: &item})
}

func (s *Server) listOrgMembers(w http.ResponseWriter, r *http.Request) {
	s.writeMembers(w, r)
}

func (s *Server) writeMembers(w http.ResponseWriter, r *http.Request) {
	actor := s.requireOrg(w, r, r.PathValue("orgId"))
	if actor == nil {
		return
	}
	items, errResp, status := s.Business.ListMembers(r.Context(), actor, r.PathValue("orgId"))
	if errResp != nil {
		write(w, status, contract.Envelope{Error: errResp})
		return
	}
	write(w, status, contract.Envelope{OK: true, Members: items})
}

func (s *Server) inviteOrgMember(w http.ResponseWriter, r *http.Request) {
	actor := s.requireOrg(w, r, r.PathValue("orgId"))
	if actor == nil {
		return
	}
	var body struct {
		Email string `json:"email"`
		Role  string `json:"role"`
	}
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 8<<10)).Decode(&body); err != nil {
		write(w, http.StatusBadRequest, contract.Envelope{Error: &contract.APIError{Code: "validation_error", Message: "Expected JSON."}})
		return
	}
	member, invite, errResp, status := s.Business.Invite(r.Context(), actor, r.PathValue("orgId"), body.Email, body.Role)
	if errResp != nil {
		write(w, status, contract.Envelope{Error: errResp})
		return
	}
	env := contract.Envelope{OK: true, Member: member}
	if invite != nil {
		env.Invites = []contract.OrgInvite{*invite}
	}
	write(w, status, env)
}

func (s *Server) patchOrgMember(w http.ResponseWriter, r *http.Request) {
	actor := s.requireOrg(w, r, r.PathValue("orgId"))
	if actor == nil {
		return
	}
	var body struct {
		Role string `json:"role"`
	}
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 8<<10)).Decode(&body); err != nil {
		write(w, http.StatusBadRequest, contract.Envelope{Error: &contract.APIError{Code: "validation_error", Message: "Expected JSON."}})
		return
	}
	item, errResp, status := s.Business.ChangeRole(r.Context(), actor, r.PathValue("orgId"), r.PathValue("id"), body.Role)
	if errResp != nil {
		write(w, status, contract.Envelope{Error: errResp})
		return
	}
	write(w, status, contract.Envelope{OK: true, Member: item})
}

func (s *Server) deleteOrgMember(w http.ResponseWriter, r *http.Request) {
	actor := s.requireOrg(w, r, r.PathValue("orgId"))
	if actor == nil {
		return
	}
	if errResp := s.Business.RemoveMember(r.Context(), actor, r.PathValue("orgId"), r.PathValue("id")); errResp != nil {
		write(w, statusForAPI(errResp), contract.Envelope{Error: errResp})
		return
	}
	write(w, http.StatusOK, contract.Envelope{OK: true})
}

func (s *Server) listOrgInvites(w http.ResponseWriter, r *http.Request) {
	actor := s.requireOrg(w, r, r.PathValue("orgId"))
	if actor == nil {
		return
	}
	items, errResp, status := s.Business.ListInvites(r.Context(), actor, r.PathValue("orgId"))
	if errResp != nil {
		write(w, status, contract.Envelope{Error: errResp})
		return
	}
	write(w, status, contract.Envelope{OK: true, Invites: items})
}

func (s *Server) listOrgTeams(w http.ResponseWriter, r *http.Request) {
	actor := s.requireOrg(w, r, r.PathValue("orgId"))
	if actor == nil {
		return
	}
	items, errResp, status := s.Business.ListTeams(r.Context(), actor, r.PathValue("orgId"))
	if errResp != nil {
		write(w, status, contract.Envelope{Error: errResp})
		return
	}
	write(w, status, contract.Envelope{OK: true, Teams: items})
}

func (s *Server) createOrgTeam(w http.ResponseWriter, r *http.Request) {
	actor := s.requireOrg(w, r, r.PathValue("orgId"))
	if actor == nil {
		return
	}
	var body struct {
		Name      string   `json:"name"`
		MemberIDs []string `json:"memberIds"`
	}
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 8<<10)).Decode(&body); err != nil {
		write(w, http.StatusBadRequest, contract.Envelope{Error: &contract.APIError{Code: "validation_error", Message: "Expected JSON."}})
		return
	}
	item, errResp, status := s.Business.CreateTeam(r.Context(), actor, r.PathValue("orgId"), body.Name, body.MemberIDs)
	if errResp != nil {
		write(w, status, contract.Envelope{Error: errResp})
		return
	}
	write(w, status, contract.Envelope{OK: true, Team: item})
}

func (s *Server) patchOrgTeam(w http.ResponseWriter, r *http.Request) {
	actor := s.requireOrg(w, r, r.PathValue("orgId"))
	if actor == nil {
		return
	}
	var body struct {
		Name      string   `json:"name"`
		MemberIDs []string `json:"memberIds"`
	}
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 8<<10)).Decode(&body); err != nil {
		write(w, http.StatusBadRequest, contract.Envelope{Error: &contract.APIError{Code: "validation_error", Message: "Expected JSON."}})
		return
	}
	item, errResp, status := s.Business.UpdateTeam(r.Context(), actor, r.PathValue("orgId"), r.PathValue("id"), body.Name, body.MemberIDs)
	if errResp != nil {
		write(w, status, contract.Envelope{Error: errResp})
		return
	}
	write(w, status, contract.Envelope{OK: true, Team: item})
}

func (s *Server) deleteOrgTeam(w http.ResponseWriter, r *http.Request) {
	actor := s.requireOrg(w, r, r.PathValue("orgId"))
	if actor == nil {
		return
	}
	if errResp := s.Business.DeleteTeam(r.Context(), actor, r.PathValue("orgId"), r.PathValue("id")); errResp != nil {
		write(w, statusForAPI(errResp), contract.Envelope{Error: errResp})
		return
	}
	write(w, http.StatusOK, contract.Envelope{OK: true})
}

func (s *Server) listOrgDevices(w http.ResponseWriter, r *http.Request) {
	actor := s.requireOrg(w, r, r.PathValue("orgId"))
	if actor == nil {
		return
	}
	items, errResp, status := s.Business.ListDevices(r.Context(), actor, r.PathValue("orgId"))
	if errResp != nil {
		write(w, status, contract.Envelope{Error: errResp})
		return
	}
	write(w, status, contract.Envelope{OK: true, OrgDevices: items})
}

func (s *Server) createOrgDevice(w http.ResponseWriter, r *http.Request) {
	actor := s.requireOrg(w, r, r.PathValue("orgId"))
	if actor == nil {
		return
	}
	var body struct {
		Name   string `json:"name"`
		Label  string `json:"label"`
		Region string `json:"region"`
	}
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 8<<10)).Decode(&body); err != nil {
		write(w, http.StatusBadRequest, contract.Envelope{Error: &contract.APIError{Code: "validation_error", Message: "Expected JSON."}})
		return
	}
	item, errResp, status := s.Business.CreateDevice(r.Context(), actor, r.PathValue("orgId"), body.Name, body.Label, body.Region)
	if errResp != nil {
		write(w, status, contract.Envelope{Error: errResp})
		return
	}
	write(w, status, contract.Envelope{OK: true, OrgDevice: item})
}

func (s *Server) deleteOrgDevice(w http.ResponseWriter, r *http.Request) {
	actor := s.requireOrg(w, r, r.PathValue("orgId"))
	if actor == nil {
		return
	}
	if errResp := s.Business.DeleteDevice(r.Context(), actor, r.PathValue("orgId"), r.PathValue("id")); errResp != nil {
		write(w, statusForAPI(errResp), contract.Envelope{Error: errResp})
		return
	}
	write(w, http.StatusOK, contract.Envelope{OK: true})
}

func (s *Server) listOrgNetworks(w http.ResponseWriter, r *http.Request) {
	actor := s.requireOrg(w, r, r.PathValue("orgId"))
	if actor == nil {
		return
	}
	items, errResp, status := s.Business.ListNetworks(r.Context(), actor, r.PathValue("orgId"))
	if errResp != nil {
		write(w, status, contract.Envelope{Error: errResp})
		return
	}
	write(w, status, contract.Envelope{OK: true, Networks: items})
}

func (s *Server) createOrgNetwork(w http.ResponseWriter, r *http.Request) {
	actor := s.requireOrg(w, r, r.PathValue("orgId"))
	if actor == nil {
		return
	}
	var body struct {
		Name   string `json:"name"`
		ASN    string `json:"asn"`
		Region string `json:"region"`
	}
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 8<<10)).Decode(&body); err != nil {
		write(w, http.StatusBadRequest, contract.Envelope{Error: &contract.APIError{Code: "validation_error", Message: "Expected JSON."}})
		return
	}
	item, errResp, status := s.Business.CreateNetwork(r.Context(), actor, r.PathValue("orgId"), body.Name, body.ASN, body.Region)
	if errResp != nil {
		write(w, status, contract.Envelope{Error: errResp})
		return
	}
	write(w, status, contract.Envelope{OK: true, Network: item})
}

func (s *Server) deleteOrgNetwork(w http.ResponseWriter, r *http.Request) {
	actor := s.requireOrg(w, r, r.PathValue("orgId"))
	if actor == nil {
		return
	}
	if errResp := s.Business.DeleteNetwork(r.Context(), actor, r.PathValue("orgId"), r.PathValue("id")); errResp != nil {
		write(w, statusForAPI(errResp), contract.Envelope{Error: errResp})
		return
	}
	write(w, http.StatusOK, contract.Envelope{OK: true})
}

func (s *Server) listOrgServices(w http.ResponseWriter, r *http.Request) {
	actor := s.requireOrg(w, r, r.PathValue("orgId"))
	if actor == nil {
		return
	}
	items, errResp, status := s.Business.ListServices(r.Context(), actor, r.PathValue("orgId"))
	if errResp != nil {
		write(w, status, contract.Envelope{Error: errResp})
		return
	}
	write(w, status, contract.Envelope{OK: true, OrgServices: items})
}

func (s *Server) createOrgService(w http.ResponseWriter, r *http.Request) {
	actor := s.requireOrg(w, r, r.PathValue("orgId"))
	if actor == nil {
		return
	}
	var body struct {
		Name     string `json:"name"`
		Slug     string `json:"slug"`
		Endpoint string `json:"endpoint"`
	}
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 8<<10)).Decode(&body); err != nil {
		write(w, http.StatusBadRequest, contract.Envelope{Error: &contract.APIError{Code: "validation_error", Message: "Expected JSON."}})
		return
	}
	item, errResp, status := s.Business.CreateService(r.Context(), actor, r.PathValue("orgId"), body.Name, body.Slug, body.Endpoint)
	if errResp != nil {
		write(w, status, contract.Envelope{Error: errResp})
		return
	}
	write(w, status, contract.Envelope{OK: true, OrgService: item})
}

func (s *Server) deleteOrgService(w http.ResponseWriter, r *http.Request) {
	actor := s.requireOrg(w, r, r.PathValue("orgId"))
	if actor == nil {
		return
	}
	if errResp := s.Business.DeleteService(r.Context(), actor, r.PathValue("orgId"), r.PathValue("id")); errResp != nil {
		write(w, statusForAPI(errResp), contract.Envelope{Error: errResp})
		return
	}
	write(w, http.StatusOK, contract.Envelope{OK: true})
}

func (s *Server) listOrgMonitors(w http.ResponseWriter, r *http.Request) {
	actor := s.requireOrg(w, r, r.PathValue("orgId"))
	if actor == nil {
		return
	}
	items, errResp, status := s.Business.ListMonitors(r.Context(), actor, r.PathValue("orgId"))
	if errResp != nil {
		write(w, status, contract.Envelope{Error: errResp})
		return
	}
	write(w, status, contract.Envelope{OK: true, Monitors: items})
}

func (s *Server) createOrgMonitor(w http.ResponseWriter, r *http.Request) {
	actor := s.requireOrg(w, r, r.PathValue("orgId"))
	if actor == nil {
		return
	}
	var body struct {
		Name             string   `json:"name"`
		Target           string   `json:"target"`
		Type             string   `json:"type"`
		Regions          []string `json:"regions"`
		FrequencySeconds int      `json:"frequencySeconds"`
		TimeoutSeconds   int      `json:"timeoutSeconds"`
		DeviceID         *string  `json:"deviceId"`
		NetworkID        *string  `json:"networkId"`
		ServiceID        *string  `json:"serviceId"`
	}
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 16<<10)).Decode(&body); err != nil {
		write(w, http.StatusBadRequest, contract.Envelope{Error: &contract.APIError{Code: "validation_error", Message: "Expected JSON."}})
		return
	}
	item, errResp, status := s.Business.CreateMonitor(r.Context(), actor, r.PathValue("orgId"), struct {
		Name, Target, Type string
		Regions            []string
		FrequencyS         int
		TimeoutS           int
		DeviceID, NetworkID, ServiceID *string
	}{body.Name, body.Target, body.Type, body.Regions, body.FrequencySeconds, body.TimeoutSeconds, body.DeviceID, body.NetworkID, body.ServiceID})
	if errResp != nil {
		write(w, status, contract.Envelope{Error: errResp})
		return
	}
	write(w, status, contract.Envelope{OK: true, Monitor: item})
}

func (s *Server) getOrgMonitor(w http.ResponseWriter, r *http.Request) {
	actor := s.requireOrg(w, r, r.PathValue("orgId"))
	if actor == nil {
		return
	}
	item, checks, errResp, status := s.Business.GetMonitor(r.Context(), actor, r.PathValue("orgId"), r.PathValue("id"))
	if errResp != nil {
		write(w, status, contract.Envelope{Error: errResp})
		return
	}
	write(w, status, contract.Envelope{OK: true, Monitor: item, Checks: checks})
}

func (s *Server) deleteOrgMonitor(w http.ResponseWriter, r *http.Request) {
	actor := s.requireOrg(w, r, r.PathValue("orgId"))
	if actor == nil {
		return
	}
	if errResp := s.Business.DeleteMonitor(r.Context(), actor, r.PathValue("orgId"), r.PathValue("id")); errResp != nil {
		write(w, statusForAPI(errResp), contract.Envelope{Error: errResp})
		return
	}
	write(w, http.StatusOK, contract.Envelope{OK: true})
}

func (s *Server) listOrgIncidents(w http.ResponseWriter, r *http.Request) {
	actor := s.requireOrg(w, r, r.PathValue("orgId"))
	if actor == nil {
		return
	}
	items, errResp, status := s.Business.ListIncidents(r.Context(), actor, r.PathValue("orgId"))
	if errResp != nil {
		write(w, status, contract.Envelope{Error: errResp})
		return
	}
	write(w, status, contract.Envelope{OK: true, OrgIncidents: items})
}

func (s *Server) getOrgIncident(w http.ResponseWriter, r *http.Request) {
	actor := s.requireOrg(w, r, r.PathValue("orgId"))
	if actor == nil {
		return
	}
	item, errResp, status := s.Business.GetIncident(r.Context(), actor, r.PathValue("orgId"), r.PathValue("id"))
	if errResp != nil {
		write(w, status, contract.Envelope{Error: errResp})
		return
	}
	write(w, status, contract.Envelope{OK: true, OrgIncident: item})
}

func (s *Server) patchOrgIncident(w http.ResponseWriter, r *http.Request) {
	actor := s.requireOrg(w, r, r.PathValue("orgId"))
	if actor == nil {
		return
	}
	var body struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 8<<10)).Decode(&body); err != nil {
		write(w, http.StatusBadRequest, contract.Envelope{Error: &contract.APIError{Code: "validation_error", Message: "Expected JSON."}})
		return
	}
	item, errResp, status := s.Business.UpdateIncident(r.Context(), actor, r.PathValue("orgId"), r.PathValue("id"), body.Status)
	if errResp != nil {
		write(w, status, contract.Envelope{Error: errResp})
		return
	}
	write(w, status, contract.Envelope{OK: true, OrgIncident: item})
}

func (s *Server) listOrgReports(w http.ResponseWriter, r *http.Request) {
	actor := s.requireOrg(w, r, r.PathValue("orgId"))
	if actor == nil {
		return
	}
	items, errResp, status := s.Business.ListReports(r.Context(), actor, r.PathValue("orgId"))
	if errResp != nil {
		write(w, status, contract.Envelope{Error: errResp})
		return
	}
	write(w, status, contract.Envelope{OK: true, OrgReports: items})
}

func (s *Server) createOrgReport(w http.ResponseWriter, r *http.Request) {
	actor := s.requireOrg(w, r, r.PathValue("orgId"))
	if actor == nil {
		return
	}
	var body struct {
		Kind string `json:"kind"`
	}
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 8<<10)).Decode(&body); err != nil {
		write(w, http.StatusBadRequest, contract.Envelope{Error: &contract.APIError{Code: "validation_error", Message: "Expected JSON."}})
		return
	}
	item, errResp, status := s.Business.GenerateReport(r.Context(), actor, r.PathValue("orgId"), body.Kind)
	if errResp != nil {
		write(w, status, contract.Envelope{Error: errResp})
		return
	}
	write(w, status, contract.Envelope{OK: true, OrgReport: item})
}

func (s *Server) getOrgReport(w http.ResponseWriter, r *http.Request) {
	actor := s.requireOrg(w, r, r.PathValue("orgId"))
	if actor == nil {
		return
	}
	item, errResp, status := s.Business.GetReport(r.Context(), actor, r.PathValue("orgId"), r.PathValue("id"))
	if errResp != nil {
		write(w, status, contract.Envelope{Error: errResp})
		return
	}
	write(w, status, contract.Envelope{OK: true, OrgReport: item})
}

func (s *Server) listOrgDiagnoses(w http.ResponseWriter, r *http.Request) {
	actor := s.requireOrg(w, r, r.PathValue("orgId"))
	if actor == nil {
		return
	}
	items, errResp, status := s.Business.ListDiagnoses(r.Context(), actor, r.PathValue("orgId"))
	if errResp != nil {
		write(w, status, contract.Envelope{Error: errResp})
		return
	}
	write(w, status, contract.Envelope{OK: true, OrgDiagnoses: items})
}

func (s *Server) createOrgDiagnosis(w http.ResponseWriter, r *http.Request) {
	actor := s.requireOrg(w, r, r.PathValue("orgId"))
	if actor == nil {
		return
	}
	var body struct {
		Target string `json:"target"`
	}
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 8<<10)).Decode(&body); err != nil {
		write(w, http.StatusBadRequest, contract.Envelope{Error: &contract.APIError{Code: "validation_error", Message: "Expected JSON."}})
		return
	}
	item, errResp, status := s.Business.CreateDiagnosis(r.Context(), actor, r.PathValue("orgId"), body.Target)
	if errResp != nil {
		write(w, status, contract.Envelope{Error: errResp})
		return
	}
	write(w, status, contract.Envelope{OK: true, OrgDiagnoses: []contract.OrgDiagnosis{*item}})
}

func (s *Server) listOrgKeys(w http.ResponseWriter, r *http.Request) {
	actor := s.requireOrg(w, r, r.PathValue("orgId"))
	if actor == nil {
		return
	}
	items, errResp, status := s.Business.ListKeys(r.Context(), actor, r.PathValue("orgId"))
	if errResp != nil {
		write(w, status, contract.Envelope{Error: errResp})
		return
	}
	write(w, status, contract.Envelope{OK: true, APIKeys: items})
}

func (s *Server) createOrgKey(w http.ResponseWriter, r *http.Request) {
	actor := s.requireOrg(w, r, r.PathValue("orgId"))
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
	item, secret, errResp, status := s.Business.CreateKey(r.Context(), actor, r.PathValue("orgId"), body.Name, body.Scopes, body.RateLimitPerMin)
	if errResp != nil {
		write(w, status, contract.Envelope{Error: errResp})
		return
	}
	write(w, status, contract.Envelope{OK: true, APIKey: item, KeySecret: secret})
}

func (s *Server) rotateOrgKey(w http.ResponseWriter, r *http.Request) {
	actor := s.requireOrg(w, r, r.PathValue("orgId"))
	if actor == nil {
		return
	}
	item, secret, errResp, status := s.Business.RotateKey(r.Context(), actor, r.PathValue("orgId"), r.PathValue("id"))
	if errResp != nil {
		write(w, status, contract.Envelope{Error: errResp})
		return
	}
	write(w, status, contract.Envelope{OK: true, APIKey: item, KeySecret: secret})
}

func (s *Server) revokeOrgKey(w http.ResponseWriter, r *http.Request) {
	actor := s.requireOrg(w, r, r.PathValue("orgId"))
	if actor == nil {
		return
	}
	item, errResp, status := s.Business.RevokeKey(r.Context(), actor, r.PathValue("orgId"), r.PathValue("id"))
	if errResp != nil {
		write(w, status, contract.Envelope{Error: errResp})
		return
	}
	write(w, status, contract.Envelope{OK: true, APIKey: item})
}

func (s *Server) orgBilling(w http.ResponseWriter, r *http.Request) {
	actor := s.requireOrg(w, r, r.PathValue("orgId"))
	if actor == nil {
		return
	}
	item, errResp, status := s.Business.Billing(r.Context(), actor, r.PathValue("orgId"))
	if errResp != nil {
		write(w, status, contract.Envelope{Error: errResp})
		return
	}
	write(w, status, contract.Envelope{OK: true, Billing: item})
}

func (s *Server) orgAudit(w http.ResponseWriter, r *http.Request) {
	actor := s.requireOrg(w, r, r.PathValue("orgId"))
	if actor == nil {
		return
	}
	items, errResp, status := s.Business.ListAudit(r.Context(), actor, r.PathValue("orgId"))
	if errResp != nil {
		write(w, status, contract.Envelope{Error: errResp})
		return
	}
	write(w, status, contract.Envelope{OK: true, AuditEvents: items})
}
