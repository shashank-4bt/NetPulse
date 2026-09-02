package api

import (
	"net/http"

	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/contract"
)

func (s *Server) meProfile(w http.ResponseWriter, r *http.Request) {
	if !s.requireAccounts(w) {
		return
	}
	_, user, apiErr, status := s.Accounts.Require(r.Context(), SessionToken(r))
	if apiErr != nil {
		write(w, status, contract.Envelope{Error: apiErr})
		return
	}
	write(w, status, contract.Envelope{OK: true, User: user})
}

func (s *Server) patchProfile(w http.ResponseWriter, r *http.Request) {
	if !s.requireAccounts(w) {
		return
	}
	var body struct {
		DisplayName string `json:"displayName"`
	}
	if !decodeJSON(w, r, &body) {
		return
	}
	user, apiErr, status := s.Accounts.UpdateProfile(r.Context(), SessionToken(r), body.DisplayName)
	if apiErr != nil {
		write(w, status, contract.Envelope{Error: apiErr})
		return
	}
	write(w, status, contract.Envelope{OK: true, User: user})
}

func (s *Server) mePrivacy(w http.ResponseWriter, r *http.Request) {
	if !s.requireAccounts(w) {
		return
	}
	settings, apiErr, status := s.Accounts.Privacy(r.Context(), SessionToken(r))
	if apiErr != nil {
		write(w, status, contract.Envelope{Error: apiErr})
		return
	}
	write(w, status, contract.Envelope{OK: true, Privacy: settings})
}

func (s *Server) putPrivacy(w http.ResponseWriter, r *http.Request) {
	if !s.requireAccounts(w) {
		return
	}
	var body struct {
		TelemetryOptIn bool `json:"telemetryOptIn"`
	}
	if !decodeJSON(w, r, &body) {
		return
	}
	settings, apiErr, status := s.Accounts.UpdatePrivacy(r.Context(), SessionToken(r), body.TelemetryOptIn)
	if apiErr != nil {
		write(w, status, contract.Envelope{Error: apiErr})
		return
	}
	write(w, status, contract.Envelope{OK: true, Privacy: settings})
}

func (s *Server) deleteAccount(w http.ResponseWriter, r *http.Request) {
	if !s.requireAccounts(w) {
		return
	}
	if apiErr, status := s.Accounts.DeleteAccount(r.Context(), SessionToken(r)); apiErr != nil {
		write(w, status, contract.Envelope{Error: apiErr})
		return
	}
	write(w, http.StatusOK, contract.Envelope{OK: true})
}

func (s *Server) dashboard(w http.ResponseWriter, r *http.Request) {
	if !s.requireAccounts(w) {
		return
	}
	dash, apiErr, status := s.Accounts.Dashboard(r.Context(), SessionToken(r))
	if apiErr != nil {
		write(w, status, contract.Envelope{Error: apiErr})
		return
	}
	write(w, status, contract.Envelope{OK: true, Dashboard: dash})
}

func (s *Server) myDiagnoses(w http.ResponseWriter, r *http.Request) {
	if !s.requireAccounts(w) {
		return
	}
	q := r.URL.Query()
	items, apiErr, status := s.Accounts.ListDiagnoses(r.Context(), SessionToken(r), q.Get("q"), q.Get("status"), q.Get("target"), q.Get("from"), q.Get("to"))
	if apiErr != nil {
		write(w, status, contract.Envelope{Error: apiErr})
		return
	}
	write(w, status, contract.Envelope{OK: true, Diagnoses: items})
}

func (s *Server) myReports(w http.ResponseWriter, r *http.Request) {
	if !s.requireAccounts(w) {
		return
	}
	items, apiErr, status := s.Accounts.ListReports(r.Context(), SessionToken(r))
	if apiErr != nil {
		write(w, status, contract.Envelope{Error: apiErr})
		return
	}
	write(w, status, contract.Envelope{OK: true, Reports: items})
}

func (s *Server) shareReport(w http.ResponseWriter, r *http.Request) {
	if !s.requireAccounts(w) {
		return
	}
	share, apiErr, status := s.Accounts.ShareReport(r.Context(), SessionToken(r), r.PathValue("id"))
	if apiErr != nil {
		write(w, status, contract.Envelope{Error: apiErr})
		return
	}
	write(w, status, contract.Envelope{OK: true, Share: share})
}

func (s *Server) deleteReport(w http.ResponseWriter, r *http.Request) {
	if !s.requireAccounts(w) {
		return
	}
	if apiErr, status := s.Accounts.DeleteReport(r.Context(), SessionToken(r), r.PathValue("id")); apiErr != nil {
		write(w, status, contract.Envelope{Error: apiErr})
		return
	}
	write(w, http.StatusOK, contract.Envelope{OK: true})
}

func (s *Server) listSaved(w http.ResponseWriter, r *http.Request) {
	if !s.requireAccounts(w) {
		return
	}
	items, apiErr, status := s.Accounts.ListSaved(r.Context(), SessionToken(r))
	if apiErr != nil {
		write(w, status, contract.Envelope{Error: apiErr})
		return
	}
	write(w, status, contract.Envelope{OK: true, Saved: items})
}

func (s *Server) saveService(w http.ResponseWriter, r *http.Request) {
	if !s.requireAccounts(w) {
		return
	}
	var body struct {
		Slug string `json:"slug"`
	}
	if !decodeJSON(w, r, &body) {
		return
	}
	items, apiErr, status := s.Accounts.SaveService(r.Context(), SessionToken(r), body.Slug)
	if apiErr != nil {
		write(w, status, contract.Envelope{Error: apiErr})
		return
	}
	write(w, status, contract.Envelope{OK: true, Saved: items})
}

func (s *Server) deleteSaved(w http.ResponseWriter, r *http.Request) {
	if !s.requireAccounts(w) {
		return
	}
	items, apiErr, status := s.Accounts.DeleteSaved(r.Context(), SessionToken(r), r.PathValue("slug"))
	if apiErr != nil {
		write(w, status, contract.Envelope{Error: apiErr})
		return
	}
	write(w, status, contract.Envelope{OK: true, Saved: items})
}

func (s *Server) devices(w http.ResponseWriter, r *http.Request) {
	if !s.requireAccounts(w) {
		return
	}
	items, apiErr, status := s.Accounts.Devices(r.Context(), SessionToken(r))
	if apiErr != nil {
		write(w, status, contract.Envelope{Error: apiErr})
		return
	}
	write(w, status, contract.Envelope{OK: true, Devices: items})
}

func (s *Server) meAlerts(w http.ResponseWriter, r *http.Request) {
	if !s.requireAccounts(w) {
		return
	}
	alerts, apiErr, status := s.Accounts.Alerts(r.Context(), SessionToken(r))
	if apiErr != nil {
		write(w, status, contract.Envelope{Error: apiErr})
		return
	}
	write(w, status, contract.Envelope{OK: true, Alerts: alerts})
}

func (s *Server) putAlerts(w http.ResponseWriter, r *http.Request) {
	if !s.requireAccounts(w) {
		return
	}
	var body struct {
		EmailEnabled   bool `json:"emailEnabled"`
		IncidentAlerts bool `json:"incidentAlerts"`
	}
	if !decodeJSON(w, r, &body) {
		return
	}
	alerts, apiErr, status := s.Accounts.UpdateAlerts(r.Context(), SessionToken(r), body.EmailEnabled, body.IncidentAlerts)
	if apiErr != nil {
		write(w, status, contract.Envelope{Error: apiErr})
		return
	}
	write(w, status, contract.Envelope{OK: true, Alerts: alerts})
}

func (s *Server) meBilling(w http.ResponseWriter, r *http.Request) {
	if !s.requireAccounts(w) {
		return
	}
	billing, apiErr, status := s.Accounts.Billing(r.Context(), SessionToken(r))
	if apiErr != nil {
		write(w, status, contract.Envelope{Error: apiErr})
		return
	}
	write(w, status, contract.Envelope{OK: true, Billing: billing})
}

func (s *Server) userBilling(w http.ResponseWriter, r *http.Request) {
	if !s.requireAccounts(w) {
		return
	}
	if apiErr, status := s.Accounts.ForeignBilling(r.Context(), SessionToken(r), r.PathValue("id")); apiErr != nil {
		write(w, status, contract.Envelope{Error: apiErr})
		return
	}
	billing, apiErr, status := s.Accounts.Billing(r.Context(), SessionToken(r))
	if apiErr != nil {
		write(w, status, contract.Envelope{Error: apiErr})
		return
	}
	write(w, status, contract.Envelope{OK: true, Billing: billing})
}

func (s *Server) organization(w http.ResponseWriter, r *http.Request) {
	if !s.requireAccounts(w) {
		return
	}
	apiErr, status := s.Accounts.Organization(r.Context(), SessionToken(r), r.PathValue("id"))
	write(w, status, contract.Envelope{Error: apiErr})
}

func (s *Server) readShare(w http.ResponseWriter, r *http.Request) {
	if !s.requireAccounts(w) {
		return
	}
	diag, apiErr, status := s.Accounts.ReadShare(r.Context(), r.PathValue("token"))
	if apiErr != nil {
		write(w, status, contract.Envelope{Error: apiErr})
		return
	}
	write(w, status, contract.Envelope{OK: true, Diagnosis: diag})
}
