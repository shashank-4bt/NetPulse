package api

import (
	"encoding/json"
	"log/slog"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/accounts"
	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/admin"
	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/business"
	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/config"
	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/contract"
	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/developer"
	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/diagnostics"
	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/geo"
	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/incidents"
	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/storage"
)

type Server struct {
	Cfg         config.Config
	Log         *slog.Logger
	Diagnostics *diagnostics.Service
	Accounts    *accounts.Service
	Developer   *developer.Service
	Business    *business.Service
	Admin       *admin.Service
	Limiter     storage.RateLimiter
	StorageInfo map[string]string
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /v1/health", s.health)
	mux.HandleFunc("POST /v1/diagnoses", s.createDiagnosis)
	mux.HandleFunc("GET /v1/diagnoses/{id}", s.getDiagnosis)
	mux.HandleFunc("GET /v1/services", s.listServices)
	mux.HandleFunc("GET /v1/services/{slug}", s.getService)
	mux.HandleFunc("GET /v1/incidents", s.listIncidents)
	mux.HandleFunc("GET /v1/incidents/{id}", s.getIncident)
	mux.HandleFunc("GET /v1/map/aggregates", s.mapAggregates)
	mux.HandleFunc("POST /v1/auth/register", s.register)
	mux.HandleFunc("POST /v1/auth/login", s.login)
	mux.HandleFunc("POST /v1/auth/logout", s.logout)
	mux.HandleFunc("GET /v1/auth/me", s.me)
	mux.HandleFunc("POST /v1/auth/verify-email", s.verifyEmail)
	mux.HandleFunc("POST /v1/auth/resend-verification", s.resendVerification)
	mux.HandleFunc("POST /v1/auth/forgot-password", s.forgotPassword)
	mux.HandleFunc("POST /v1/auth/reset-password", s.resetPassword)
	mux.HandleFunc("POST /v1/auth/change-password", s.changePassword)
	mux.HandleFunc("GET /v1/auth/sessions", s.listSessions)
	mux.HandleFunc("POST /v1/auth/sessions/{id}/revoke", s.revokeSession)
	mux.HandleFunc("POST /v1/auth/sessions/revoke-others", s.revokeOtherSessions)
	mux.HandleFunc("GET /v1/auth/events", s.listEvents)
	mux.HandleFunc("GET /v1/auth/methods", s.authMethods)
	mux.HandleFunc("POST /v1/auth/oauth/{provider}", s.unsupportedFactor)
	mux.HandleFunc("POST /v1/auth/passkeys", s.unsupportedFactor)
	mux.HandleFunc("POST /v1/auth/mfa", s.unsupportedFactor)
	mux.HandleFunc("GET /v1/me/profile", s.meProfile)
	mux.HandleFunc("PATCH /v1/me/profile", s.patchProfile)
	mux.HandleFunc("GET /v1/me/privacy", s.mePrivacy)
	mux.HandleFunc("PUT /v1/me/privacy", s.putPrivacy)
	mux.HandleFunc("POST /v1/me/deletion", s.deleteAccount)
	mux.HandleFunc("GET /v1/me/dashboard", s.dashboard)
	mux.HandleFunc("GET /v1/me/diagnoses", s.myDiagnoses)
	mux.HandleFunc("GET /v1/me/reports", s.myReports)
	mux.HandleFunc("POST /v1/me/reports/{id}/share", s.shareReport)
	mux.HandleFunc("DELETE /v1/me/reports/{id}", s.deleteReport)
	mux.HandleFunc("GET /v1/me/saved-services", s.listSaved)
	mux.HandleFunc("PUT /v1/me/saved-services", s.saveService)
	mux.HandleFunc("DELETE /v1/me/saved-services/{slug}", s.deleteSaved)
	mux.HandleFunc("GET /v1/me/devices", s.devices)
	mux.HandleFunc("GET /v1/me/alerts", s.meAlerts)
	mux.HandleFunc("PUT /v1/me/alerts", s.putAlerts)
	mux.HandleFunc("GET /v1/me/billing", s.meBilling)
	mux.HandleFunc("GET /v1/users/{id}/billing", s.userBilling)
	mux.HandleFunc("GET /v1/organizations/{id}", s.organization)
	mux.HandleFunc("GET /v1/shares/{token}", s.readShare)
	mux.HandleFunc("GET /v1/dev/workspace", s.devWorkspace)
	mux.HandleFunc("GET /v1/dev/workspaces/{id}", s.getDevWorkspace)
	mux.HandleFunc("GET /v1/dev/dashboard", s.devDashboard)
	mux.HandleFunc("GET /v1/dev/monitors", s.listDevMonitors)
	mux.HandleFunc("POST /v1/dev/monitors", s.createDevMonitor)
	mux.HandleFunc("GET /v1/dev/monitors/{id}", s.getDevMonitor)
	mux.HandleFunc("PATCH /v1/dev/monitors/{id}", s.patchDevMonitor)
	mux.HandleFunc("DELETE /v1/dev/monitors/{id}", s.deleteDevMonitor)
	mux.HandleFunc("POST /v1/dev/monitors/{id}/run", s.runDevMonitor)
	mux.HandleFunc("GET /v1/dev/incidents", s.listDevIncidents)
	mux.HandleFunc("GET /v1/dev/incidents/{id}", s.getDevIncident)
	mux.HandleFunc("GET /v1/dev/keys", s.listDevKeys)
	mux.HandleFunc("POST /v1/dev/keys", s.createDevKey)
	mux.HandleFunc("POST /v1/dev/keys/{id}/rotate", s.rotateDevKey)
	mux.HandleFunc("POST /v1/dev/keys/{id}/revoke", s.revokeDevKey)
	mux.HandleFunc("GET /v1/dev/webhooks", s.listDevWebhooks)
	mux.HandleFunc("POST /v1/dev/webhooks", s.createDevWebhook)
	mux.HandleFunc("POST /v1/dev/webhooks/{id}/rotate", s.rotateDevWebhook)
	mux.HandleFunc("DELETE /v1/dev/webhooks/{id}", s.deleteDevWebhook)
	mux.HandleFunc("GET /v1/dev/webhooks/{id}/deliveries", s.listDevDeliveries)
	mux.HandleFunc("POST /v1/dev/webhooks/{id}/retry", s.retryDevDeliveries)
	mux.HandleFunc("GET /v1/dev/alerts", s.listDevAlerts)
	mux.HandleFunc("POST /v1/dev/alerts", s.createDevAlert)
	mux.HandleFunc("PUT /v1/dev/alerts/{id}", s.putDevAlert)
	mux.HandleFunc("DELETE /v1/dev/alerts/{id}", s.deleteDevAlert)
	mux.HandleFunc("GET /v1/dev/usage", s.devUsage)
	mux.HandleFunc("GET /v1/dev/sla", s.devSLA)
	mux.HandleFunc("GET /v1/orgs", s.listOrgs)
	mux.HandleFunc("POST /v1/orgs", s.createOrg)
	mux.HandleFunc("GET /v1/orgs/{orgId}", s.getOrg)
	mux.HandleFunc("PATCH /v1/orgs/{orgId}", s.patchOrg)
	mux.HandleFunc("GET /v1/orgs/{orgId}/dashboard", s.orgDashboard)
	mux.HandleFunc("GET /v1/orgs/{orgId}/analytics", s.orgAnalytics)
	mux.HandleFunc("GET /v1/orgs/{orgId}/members", s.listOrgMembers)
	mux.HandleFunc("POST /v1/orgs/{orgId}/members", s.inviteOrgMember)
	mux.HandleFunc("PATCH /v1/orgs/{orgId}/members/{id}", s.patchOrgMember)
	mux.HandleFunc("DELETE /v1/orgs/{orgId}/members/{id}", s.deleteOrgMember)
	mux.HandleFunc("GET /v1/orgs/{orgId}/invites", s.listOrgInvites)
	mux.HandleFunc("GET /v1/orgs/{orgId}/teams", s.listOrgTeams)
	mux.HandleFunc("POST /v1/orgs/{orgId}/teams", s.createOrgTeam)
	mux.HandleFunc("PATCH /v1/orgs/{orgId}/teams/{id}", s.patchOrgTeam)
	mux.HandleFunc("DELETE /v1/orgs/{orgId}/teams/{id}", s.deleteOrgTeam)
	mux.HandleFunc("GET /v1/orgs/{orgId}/devices", s.listOrgDevices)
	mux.HandleFunc("POST /v1/orgs/{orgId}/devices", s.createOrgDevice)
	mux.HandleFunc("DELETE /v1/orgs/{orgId}/devices/{id}", s.deleteOrgDevice)
	mux.HandleFunc("GET /v1/orgs/{orgId}/networks", s.listOrgNetworks)
	mux.HandleFunc("POST /v1/orgs/{orgId}/networks", s.createOrgNetwork)
	mux.HandleFunc("DELETE /v1/orgs/{orgId}/networks/{id}", s.deleteOrgNetwork)
	mux.HandleFunc("GET /v1/orgs/{orgId}/services", s.listOrgServices)
	mux.HandleFunc("POST /v1/orgs/{orgId}/services", s.createOrgService)
	mux.HandleFunc("DELETE /v1/orgs/{orgId}/services/{id}", s.deleteOrgService)
	mux.HandleFunc("GET /v1/orgs/{orgId}/monitors", s.listOrgMonitors)
	mux.HandleFunc("POST /v1/orgs/{orgId}/monitors", s.createOrgMonitor)
	mux.HandleFunc("GET /v1/orgs/{orgId}/monitors/{id}", s.getOrgMonitor)
	mux.HandleFunc("DELETE /v1/orgs/{orgId}/monitors/{id}", s.deleteOrgMonitor)
	mux.HandleFunc("GET /v1/orgs/{orgId}/incidents", s.listOrgIncidents)
	mux.HandleFunc("GET /v1/orgs/{orgId}/incidents/{id}", s.getOrgIncident)
	mux.HandleFunc("PATCH /v1/orgs/{orgId}/incidents/{id}", s.patchOrgIncident)
	mux.HandleFunc("GET /v1/orgs/{orgId}/reports", s.listOrgReports)
	mux.HandleFunc("POST /v1/orgs/{orgId}/reports", s.createOrgReport)
	mux.HandleFunc("GET /v1/orgs/{orgId}/reports/{id}", s.getOrgReport)
	mux.HandleFunc("GET /v1/orgs/{orgId}/diagnoses", s.listOrgDiagnoses)
	mux.HandleFunc("POST /v1/orgs/{orgId}/diagnoses", s.createOrgDiagnosis)
	mux.HandleFunc("GET /v1/orgs/{orgId}/keys", s.listOrgKeys)
	mux.HandleFunc("POST /v1/orgs/{orgId}/keys", s.createOrgKey)
	mux.HandleFunc("POST /v1/orgs/{orgId}/keys/{id}/rotate", s.rotateOrgKey)
	mux.HandleFunc("POST /v1/orgs/{orgId}/keys/{id}/revoke", s.revokeOrgKey)
	mux.HandleFunc("GET /v1/orgs/{orgId}/billing", s.orgBilling)
	mux.HandleFunc("GET /v1/orgs/{orgId}/audit", s.orgAudit)
	mux.HandleFunc("GET /v1/admin/me", s.adminMe)
	mux.HandleFunc("GET /v1/admin/system", s.adminSystem)
	mux.HandleFunc("GET /v1/admin/users", s.adminUsers)
	mux.HandleFunc("GET /v1/admin/users/{id}", s.adminUser)
	mux.HandleFunc("GET /v1/admin/organizations", s.adminOrgs)
	mux.HandleFunc("GET /v1/admin/organizations/{id}", s.adminOrg)
	mux.HandleFunc("GET /v1/admin/services", s.adminServices)
	mux.HandleFunc("GET /v1/admin/incidents", s.adminIncidents)
	mux.HandleFunc("GET /v1/admin/incidents/{id}", s.adminIncident)
	mux.HandleFunc("POST /v1/admin/incidents/{id}/annotate", s.adminIncidentAnnotate)
	mux.HandleFunc("POST /v1/admin/incidents/{id}/investigate", s.adminIncidentInvestigate)
	mux.HandleFunc("POST /v1/admin/incidents/{id}/escalate", s.adminIncidentEscalate)
	mux.HandleFunc("POST /v1/admin/incidents/{id}/resolve", s.adminIncidentResolve)
	mux.HandleFunc("POST /v1/admin/incidents/{id}/override", s.adminIncidentOverride)
	mux.HandleFunc("GET /v1/admin/measurements", s.adminMeasurements)
	mux.HandleFunc("GET /v1/admin/diagnostics", s.adminDiagnostics)
	mux.HandleFunc("GET /v1/admin/rules", s.adminRules)
	mux.HandleFunc("POST /v1/admin/rules/labels", s.adminRuleLabel)
	mux.HandleFunc("GET /v1/admin/abuse", s.adminAbuse)
	mux.HandleFunc("GET /v1/admin/audit", s.adminAudit)
	mux.HandleFunc("GET /v1/admin/flags", s.adminFlags)
	mux.HandleFunc("POST /v1/admin/flags", s.adminCreateFlag)
	mux.HandleFunc("GET /v1/admin/flags/{id}", s.adminGetFlag)
	mux.HandleFunc("PATCH /v1/admin/flags/{id}", s.adminPatchFlag)
	mux.HandleFunc("GET /v1/admin/config", s.adminConfig)
	mux.HandleFunc("PUT /v1/admin/config", s.adminPutConfig)
	return s.middleware(mux)
}

func (s *Server) middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if origin := s.Cfg.CORSOrigin; origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-NetPulse-Session, X-NetPulse-Key")
			w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
		}
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func (s *Server) health(w http.ResponseWriter, _ *http.Request) {
	write(w, http.StatusOK, contract.Envelope{
		OK: true,
		Health: &contract.Health{
			Status:  "ok",
			Version: s.Cfg.EngineVersion,
			Storage: s.StorageInfo,
		},
	})
}

func (s *Server) createDiagnosis(w http.ResponseWriter, r *http.Request) {
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	if ip == "" {
		ip = r.RemoteAddr
	}
	if s.Limiter != nil && !s.Limiter.Allow("diag:"+ip, s.diagnoseLimit(r)) {
		s.recordAbuse(r, "rate_limit", ip, "diagnoses", "blocked", "Diagnose rate limit exceeded.")
		write(w, http.StatusTooManyRequests, contract.Envelope{Error: &contract.APIError{Code: "rate_limited", Message: "Too many diagnose requests"}})
		return
	}

	var body struct {
		Target string `json:"target"`
	}
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 8<<10)).Decode(&body); err != nil {
		write(w, http.StatusBadRequest, contract.Envelope{Error: &contract.APIError{Code: "validation_error", Message: "Expected JSON with a target"}})
		return
	}
	userID := ""
	if s.Accounts != nil {
		if _, user, _, _ := s.Accounts.Require(r.Context(), SessionToken(r)); user != nil {
			userID = user.ID
		}
	}
	diag, apiErr, status := s.Diagnostics.Create(r.Context(), body.Target, userID)
	if apiErr != nil {
		if apiErr.Code == "ssrf_blocked" {
			s.recordAbuse(r, "ssrf", userID, "diagnoses", "blocked", "Diagnose target was blocked by SSRF policy.")
		}
		write(w, status, contract.Envelope{Error: apiErr})
		return
	}
	write(w, status, contract.Envelope{OK: true, Diagnosis: diag})
}

func (s *Server) getDiagnosis(w http.ResponseWriter, r *http.Request) {
	rec, apiErr, status := s.Diagnostics.GetRecord(r.Context(), r.PathValue("id"))
	if apiErr != nil {
		write(w, status, contract.Envelope{Error: apiErr})
		return
	}
	viewerID := ""
	if s.Accounts != nil {
		if _, user, _, _ := s.Accounts.Require(r.Context(), SessionToken(r)); user != nil {
			viewerID = user.ID
		}
		if !s.Accounts.MayReadDiagnosis(r.Context(), rec, viewerID, r.URL.Query().Get("share")) {
			write(w, http.StatusNotFound, contract.Envelope{Error: &contract.APIError{Code: "not_found", Message: "diagnosis not found"}})
			return
		}
	} else if rec.UserID != "" {
		write(w, http.StatusNotFound, contract.Envelope{Error: &contract.APIError{Code: "not_found", Message: "diagnosis not found"}})
		return
	}
	copy := rec.Diagnosis
	write(w, status, contract.Envelope{OK: true, Diagnosis: &copy})
}

func (s *Server) listServices(w http.ResponseWriter, r *http.Request) {
	items, apiErr, status := s.Diagnostics.ListServices(r.Context())
	if apiErr != nil {
		write(w, status, contract.Envelope{Error: apiErr})
		return
	}
	write(w, status, contract.Envelope{OK: true, Services: items})
}

func (s *Server) getService(w http.ResponseWriter, r *http.Request) {
	item, intel, apiErr, status := s.Diagnostics.GetService(r.Context(), r.PathValue("slug"))
	if apiErr != nil {
		write(w, status, contract.Envelope{Error: apiErr})
		return
	}
	write(w, status, contract.Envelope{OK: true, Service: item, Intelligence: intel})
}

func (s *Server) listIncidents(w http.ResponseWriter, r *http.Request) {
	query := incidents.ParseQuery(map[string]string{
		"service":  r.URL.Query().Get("service"),
		"region":   r.URL.Query().Get("region"),
		"network":  r.URL.Query().Get("network"),
		"severity": r.URL.Query().Get("severity"),
		"status":   r.URL.Query().Get("status"),
		"q":        r.URL.Query().Get("q"),
		"sort":     r.URL.Query().Get("sort"),
		"time":     r.URL.Query().Get("time"),
		"page":     r.URL.Query().Get("page"),
		"pageSize": r.URL.Query().Get("pageSize"),
	}, time.Now().UTC())
	items, page, apiErr, status := s.Diagnostics.ListIncidents(r.Context(), query)
	if apiErr != nil {
		write(w, status, contract.Envelope{Error: apiErr})
		return
	}
	write(w, status, contract.Envelope{OK: true, Incidents: items, Page: page})
}

func (s *Server) mapAggregates(w http.ResponseWriter, r *http.Request) {
	query := geo.ParseQuery(map[string]string{
		"level":   r.URL.Query().Get("level"),
		"parent":  r.URL.Query().Get("parent"),
		"west":    r.URL.Query().Get("west"),
		"south":   r.URL.Query().Get("south"),
		"east":    r.URL.Query().Get("east"),
		"north":   r.URL.Query().Get("north"),
		"layers":  strings.Join(r.URL.Query()["layers"], ","),
		"q":       r.URL.Query().Get("q"),
		"service": r.URL.Query().Get("service"),
		"limit":   r.URL.Query().Get("limit"),
	})
	agg, apiErr, status := s.Diagnostics.ListMapAggregates(r.Context(), query)
	if apiErr != nil {
		write(w, status, contract.Envelope{Error: apiErr})
		return
	}
	write(w, status, contract.Envelope{OK: true, Map: agg})
}

func (s *Server) getIncident(w http.ResponseWriter, r *http.Request) {
	item, apiErr, status := s.Diagnostics.GetIncident(r.Context(), r.PathValue("id"))
	if apiErr != nil {
		write(w, status, contract.Envelope{Error: apiErr})
		return
	}
	write(w, status, contract.Envelope{OK: true, Incident: item})
}

func write(w http.ResponseWriter, status int, body contract.Envelope) {
	if body.Incidents == nil {
		body.Incidents = []contract.Incident{}
	}
	if body.Map != nil {
		normalized := contract.NormalizeMap(*body.Map)
		body.Map = &normalized
	}
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func ClientIP(r *http.Request) string {
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		return strings.TrimSpace(strings.Split(forwarded, ",")[0])
	}
	return r.RemoteAddr
}

func SessionToken(r *http.Request) string {
	if header := r.Header.Get("Authorization"); strings.HasPrefix(strings.ToLower(header), "session ") {
		return strings.TrimSpace(header[8:])
	}
	return strings.TrimSpace(r.Header.Get("X-NetPulse-Session"))
}

func NewHTTPServer(addr string, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:              addr,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       30 * time.Second,
	}
}
