package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/admin"
	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/contract"
)

type adminActor struct {
	UserID      string
	Email       string
	Permissions []string
}

func (s *Server) requireAdmin(w http.ResponseWriter, r *http.Request, perm string) *adminActor {
	if s.Admin == nil || s.Accounts == nil {
		write(w, http.StatusServiceUnavailable, contract.Envelope{Error: &contract.APIError{Code: "unavailable", Message: "Admin platform is unavailable."}})
		return nil
	}
	_, user, errResp, status := s.Accounts.Require(r.Context(), SessionToken(r))
	if errResp != nil {
		write(w, status, contract.Envelope{Error: errResp})
		return nil
	}
	op, errResp, status := s.Admin.OperatorFor(r.Context(), user.ID, user.Email)
	if errResp != nil {
		write(w, status, contract.Envelope{Error: errResp})
		return nil
	}
	if perm != "" && !admin.HasPerm(op.Permissions, perm) {
		write(w, http.StatusForbidden, contract.Envelope{Error: &contract.APIError{Code: "forbidden", Message: "Missing permission."}})
		return nil
	}
	return &adminActor{UserID: user.ID, Email: user.Email, Permissions: op.Permissions}
}

func (s *Server) recordAbuse(r *http.Request, kind, actor, resource, result, summary string) {
	if s.Admin == nil {
		return
	}
	s.Admin.RecordAbuse(r.Context(), kind, actor, ClientIP(r), resource, result, summary)
}

func (s *Server) diagnoseLimit(r *http.Request) int {
	if s.Admin != nil {
		return s.Admin.DiagnoseLimit(r.Context())
	}
	return s.Cfg.RateLimitPerMin
}

func (s *Server) adminMe(w http.ResponseWriter, r *http.Request) {
	actor := s.requireAdmin(w, r, "")
	if actor == nil {
		return
	}
	write(w, http.StatusOK, contract.Envelope{OK: true, Operator: &contract.Operator{
		UserID: actor.UserID, Email: actor.Email, Role: admin.RoleOperator, Permissions: actor.Permissions,
	}})
}

func (s *Server) adminSystem(w http.ResponseWriter, r *http.Request) {
	if s.requireAdmin(w, r, admin.PermSystemRead) == nil {
		return
	}
	item, errResp, status := s.Admin.System(r.Context())
	if errResp != nil {
		write(w, status, contract.Envelope{Error: errResp})
		return
	}
	write(w, status, contract.Envelope{OK: true, AdminSystem: &item})
}

func (s *Server) adminUsers(w http.ResponseWriter, r *http.Request) {
	if s.requireAdmin(w, r, admin.PermUsersRead) == nil {
		return
	}
	items, errResp, status := s.Admin.ListUsers(r.Context())
	if errResp != nil {
		write(w, status, contract.Envelope{Error: errResp})
		return
	}
	write(w, status, contract.Envelope{OK: true, AdminUsers: items})
}

func (s *Server) adminUser(w http.ResponseWriter, r *http.Request) {
	if s.requireAdmin(w, r, admin.PermUsersRead) == nil {
		return
	}
	item, errResp, status := s.Admin.GetUser(r.Context(), r.PathValue("id"))
	if errResp != nil {
		write(w, status, contract.Envelope{Error: errResp})
		return
	}
	write(w, status, contract.Envelope{OK: true, AdminUser: item})
}

func (s *Server) adminOrgs(w http.ResponseWriter, r *http.Request) {
	if s.requireAdmin(w, r, admin.PermOrgsRead) == nil {
		return
	}
	items, errResp, status := s.Admin.ListOrganizations(r.Context())
	if errResp != nil {
		write(w, status, contract.Envelope{Error: errResp})
		return
	}
	write(w, status, contract.Envelope{OK: true, Organizations: items})
}

func (s *Server) adminOrg(w http.ResponseWriter, r *http.Request) {
	if s.requireAdmin(w, r, admin.PermOrgsRead) == nil {
		return
	}
	item, errResp, status := s.Admin.GetOrganization(r.Context(), r.PathValue("id"))
	if errResp != nil {
		write(w, status, contract.Envelope{Error: errResp})
		return
	}
	write(w, status, contract.Envelope{OK: true, Organization: item})
}

func (s *Server) adminServices(w http.ResponseWriter, r *http.Request) {
	if s.requireAdmin(w, r, admin.PermServicesRead) == nil {
		return
	}
	items, errResp, status := s.Admin.ListServices(r.Context())
	if errResp != nil {
		write(w, status, contract.Envelope{Error: errResp})
		return
	}
	write(w, status, contract.Envelope{OK: true, Services: items})
}

func (s *Server) adminIncidents(w http.ResponseWriter, r *http.Request) {
	if s.requireAdmin(w, r, admin.PermIncidentsRead) == nil {
		return
	}
	items, errResp, status := s.Admin.ListIncidents(r.Context())
	if errResp != nil {
		write(w, status, contract.Envelope{Error: errResp})
		return
	}
	write(w, status, contract.Envelope{OK: true, AdminIncidents: items})
}

func (s *Server) adminIncident(w http.ResponseWriter, r *http.Request) {
	if s.requireAdmin(w, r, admin.PermIncidentsRead) == nil {
		return
	}
	item, errResp, status := s.Admin.GetIncident(r.Context(), r.PathValue("id"))
	if errResp != nil {
		write(w, status, contract.Envelope{Error: errResp})
		return
	}
	write(w, status, contract.Envelope{OK: true, AdminIncident: item})
}

func (s *Server) adminIncidentAnnotate(w http.ResponseWriter, r *http.Request) {
	s.adminIncidentAction(w, r, "annotate")
}

func (s *Server) adminIncidentInvestigate(w http.ResponseWriter, r *http.Request) {
	s.adminIncidentAction(w, r, "investigate")
}

func (s *Server) adminIncidentEscalate(w http.ResponseWriter, r *http.Request) {
	s.adminIncidentAction(w, r, "escalate")
}

func (s *Server) adminIncidentAction(w http.ResponseWriter, r *http.Request, kind string) {
	actor := s.requireAdmin(w, r, admin.PermIncidentsOperate)
	if actor == nil {
		return
	}
	var body struct {
		Body   string `json:"body"`
		Reason string `json:"reason"`
	}
	_ = json.NewDecoder(http.MaxBytesReader(w, r.Body, 8<<10)).Decode(&body)
	note := strings.TrimSpace(body.Body)
	if note == "" {
		note = strings.TrimSpace(body.Reason)
	}
	var item *contract.AdminIncident
	var errResp *contract.APIError
	var status int
	switch kind {
	case "annotate":
		item, errResp, status = s.Admin.Annotate(r.Context(), actor.UserID, r.PathValue("id"), note)
	case "investigate":
		item, errResp, status = s.Admin.Investigate(r.Context(), actor.UserID, r.PathValue("id"), note)
	case "escalate":
		item, errResp, status = s.Admin.Escalate(r.Context(), actor.UserID, r.PathValue("id"), note)
	}
	if errResp != nil {
		write(w, status, contract.Envelope{Error: errResp})
		return
	}
	write(w, status, contract.Envelope{OK: true, AdminIncident: item})
}

func (s *Server) adminIncidentResolve(w http.ResponseWriter, r *http.Request) {
	actor := s.requireAdmin(w, r, admin.PermIncidentsOperate)
	if actor == nil {
		return
	}
	var body struct {
		Reason          string `json:"reason"`
		Override        bool   `json:"override"`
		Recoveries      int    `json:"recoveries"`
		IdentifiedCause bool   `json:"identifiedCause"`
	}
	_ = json.NewDecoder(http.MaxBytesReader(w, r.Body, 8<<10)).Decode(&body)
	item, errResp, status := s.Admin.Resolve(r.Context(), actor.UserID, r.PathValue("id"), body.Reason, body.Override, body.Recoveries, body.IdentifiedCause)
	if errResp != nil {
		write(w, status, contract.Envelope{Error: errResp})
		return
	}
	write(w, status, contract.Envelope{OK: true, AdminIncident: item})
}

func (s *Server) adminIncidentOverride(w http.ResponseWriter, r *http.Request) {
	actor := s.requireAdmin(w, r, admin.PermIncidentsOperate)
	if actor == nil {
		return
	}
	var body struct {
		Classification string `json:"classification"`
		Reason         string `json:"reason"`
	}
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 8<<10)).Decode(&body); err != nil {
		write(w, http.StatusBadRequest, contract.Envelope{Error: &contract.APIError{Code: "validation_error", Message: "Expected JSON."}})
		return
	}
	item, errResp, status := s.Admin.Override(r.Context(), actor.UserID, r.PathValue("id"), body.Classification, body.Reason)
	if errResp != nil {
		write(w, status, contract.Envelope{Error: errResp})
		return
	}
	write(w, status, contract.Envelope{OK: true, AdminIncident: item})
}

func (s *Server) adminMeasurements(w http.ResponseWriter, r *http.Request) {
	if s.requireAdmin(w, r, admin.PermMeasurementsRead) == nil {
		return
	}
	items, errResp, status := s.Admin.ListMeasurements(r.Context())
	if errResp != nil {
		write(w, status, contract.Envelope{Error: errResp})
		return
	}
	write(w, status, contract.Envelope{OK: true, AdminMeasurements: items})
}

func (s *Server) adminDiagnostics(w http.ResponseWriter, r *http.Request) {
	if s.requireAdmin(w, r, admin.PermDiagnosticsRead) == nil {
		return
	}
	items, errResp, status := s.Admin.ListDiagnoses(r.Context())
	if errResp != nil {
		write(w, status, contract.Envelope{Error: errResp})
		return
	}
	write(w, status, contract.Envelope{OK: true, AdminDiagnoses: items})
}

func (s *Server) adminRules(w http.ResponseWriter, r *http.Request) {
	if s.requireAdmin(w, r, admin.PermRulesRead) == nil {
		return
	}
	rules, outcomes, errResp, status := s.Admin.Rules(r.Context())
	if errResp != nil {
		write(w, status, contract.Envelope{Error: errResp})
		return
	}
	write(w, status, contract.Envelope{OK: true, AdminRules: rules, RuleOutcomes: outcomes})
}

func (s *Server) adminRuleLabel(w http.ResponseWriter, r *http.Request) {
	actor := s.requireAdmin(w, r, admin.PermRulesRead)
	if actor == nil {
		return
	}
	var body struct {
		DiagnosisID string `json:"diagnosisId"`
		Kind        string `json:"kind"`
	}
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 8<<10)).Decode(&body); err != nil {
		write(w, http.StatusBadRequest, contract.Envelope{Error: &contract.APIError{Code: "validation_error", Message: "Expected JSON."}})
		return
	}
	outcomes, errResp, status := s.Admin.LabelDiagnosis(r.Context(), actor.UserID, body.DiagnosisID, body.Kind)
	if errResp != nil {
		write(w, status, contract.Envelope{Error: errResp})
		return
	}
	write(w, status, contract.Envelope{OK: true, RuleOutcomes: outcomes})
}

func (s *Server) adminAbuse(w http.ResponseWriter, r *http.Request) {
	if s.requireAdmin(w, r, admin.PermAbuseRead) == nil {
		return
	}
	items, errResp, status := s.Admin.ListAbuse(r.Context())
	if errResp != nil {
		write(w, status, contract.Envelope{Error: errResp})
		return
	}
	write(w, status, contract.Envelope{OK: true, AbuseEvents: items})
}

func (s *Server) adminAudit(w http.ResponseWriter, r *http.Request) {
	if s.requireAdmin(w, r, admin.PermAuditRead) == nil {
		return
	}
	items, errResp, status := s.Admin.ListAudit(r.Context())
	if errResp != nil {
		write(w, status, contract.Envelope{Error: errResp})
		return
	}
	write(w, status, contract.Envelope{OK: true, AdminAudit: items})
}

func (s *Server) adminFlags(w http.ResponseWriter, r *http.Request) {
	if s.requireAdmin(w, r, admin.PermFlagsManage) == nil {
		return
	}
	items, errResp, status := s.Admin.ListFlags(r.Context())
	if errResp != nil {
		write(w, status, contract.Envelope{Error: errResp})
		return
	}
	write(w, status, contract.Envelope{OK: true, FeatureFlags: items})
}

func (s *Server) adminCreateFlag(w http.ResponseWriter, r *http.Request) {
	actor := s.requireAdmin(w, r, admin.PermFlagsManage)
	if actor == nil {
		return
	}
	flag, ok := decodeFlag(w, r)
	if !ok {
		return
	}
	item, errResp, status := s.Admin.CreateFlag(r.Context(), actor.UserID, flag)
	if errResp != nil {
		write(w, status, contract.Envelope{Error: errResp})
		return
	}
	write(w, status, contract.Envelope{OK: true, FeatureFlag: item})
}

func (s *Server) adminGetFlag(w http.ResponseWriter, r *http.Request) {
	if s.requireAdmin(w, r, admin.PermFlagsManage) == nil {
		return
	}
	q := r.URL.Query()
	item, errResp, status := s.Admin.GetFlag(r.Context(), r.PathValue("id"), q.Get("environment"), q.Get("userId"), q.Get("orgId"))
	if errResp != nil {
		write(w, status, contract.Envelope{Error: errResp})
		return
	}
	write(w, status, contract.Envelope{OK: true, FeatureFlag: item})
}

func (s *Server) adminPatchFlag(w http.ResponseWriter, r *http.Request) {
	actor := s.requireAdmin(w, r, admin.PermFlagsManage)
	if actor == nil {
		return
	}
	var raw map[string]json.RawMessage
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 16<<10)).Decode(&raw); err != nil {
		write(w, http.StatusBadRequest, contract.Envelope{Error: &contract.APIError{Code: "validation_error", Message: "Expected JSON."}})
		return
	}
	encoded, _ := json.Marshal(raw)
	var flag contract.FeatureFlag
	_ = json.Unmarshal(encoded, &flag)
	_, hasEnabled := raw["enabled"]
	_, hasPercentage := raw["percentage"]
	item, errResp, status := s.Admin.PatchFlag(r.Context(), actor.UserID, r.PathValue("id"), flag, hasEnabled, hasPercentage)
	if errResp != nil {
		write(w, status, contract.Envelope{Error: errResp})
		return
	}
	write(w, status, contract.Envelope{OK: true, FeatureFlag: item})
}

func (s *Server) adminConfig(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPut {
		s.adminPutConfig(w, r)
		return
	}
	if s.requireAdmin(w, r, admin.PermConfigManage) == nil {
		return
	}
	items, errResp, status := s.Admin.ListConfig(r.Context())
	if errResp != nil {
		write(w, status, contract.Envelope{Error: errResp})
		return
	}
	write(w, status, contract.Envelope{OK: true, RemoteConfig: items})
}

func (s *Server) adminPutConfig(w http.ResponseWriter, r *http.Request) {
	actor := s.requireAdmin(w, r, admin.PermConfigManage)
	if actor == nil {
		return
	}
	var body struct {
		Key     string `json:"key"`
		Value   string `json:"value"`
		Entries []struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		} `json:"entries"`
	}
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 16<<10)).Decode(&body); err != nil {
		write(w, http.StatusBadRequest, contract.Envelope{Error: &contract.APIError{Code: "validation_error", Message: "Expected JSON."}})
		return
	}
	if len(body.Entries) == 0 {
		items, errResp, status := s.Admin.PutConfig(r.Context(), actor.UserID, body.Key, body.Value)
		if errResp != nil {
			write(w, status, contract.Envelope{Error: errResp})
			return
		}
		write(w, status, contract.Envelope{OK: true, RemoteConfig: items})
		return
	}
	var items []contract.RemoteConfigEntry
	for _, entry := range body.Entries {
		var errResp *contract.APIError
		var status int
		items, errResp, status = s.Admin.PutConfig(r.Context(), actor.UserID, entry.Key, entry.Value)
		if errResp != nil {
			write(w, status, contract.Envelope{Error: errResp})
			return
		}
	}
	write(w, http.StatusOK, contract.Envelope{OK: true, RemoteConfig: items})
}

func decodeFlag(w http.ResponseWriter, r *http.Request) (contract.FeatureFlag, bool) {
	var body struct {
		Name        string   `json:"name"`
		Environment string   `json:"environment"`
		Enabled     bool     `json:"enabled"`
		Percentage  int      `json:"percentage"`
		UserIDs     []string `json:"userIds"`
		OrgIDs      []string `json:"orgIds"`
	}
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 16<<10)).Decode(&body); err != nil {
		write(w, http.StatusBadRequest, contract.Envelope{Error: &contract.APIError{Code: "validation_error", Message: "Expected JSON."}})
		return contract.FeatureFlag{}, false
	}
	return contract.FeatureFlag{
		Name: body.Name, Environment: body.Environment, Enabled: body.Enabled,
		Percentage: body.Percentage, UserIDs: body.UserIDs, OrgIDs: body.OrgIDs,
	}, true
}
