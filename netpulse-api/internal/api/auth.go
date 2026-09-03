package api

import (
	"encoding/json"
	"net/http"

	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/accounts"
	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/contract"
)

func (s *Server) register(w http.ResponseWriter, r *http.Request) {
	if !s.allowAuth(w, r) {
		return
	}
	var body struct {
		Email       string `json:"email"`
		Password    string `json:"password"`
		DisplayName string `json:"displayName"`
	}
	if !decodeJSON(w, r, &body) {
		return
	}
	result, apiErr, status := s.Accounts.Register(r.Context(), accounts.RegisterInput{
		Email: body.Email, Password: body.Password, DisplayName: body.DisplayName,
		UserAgent: r.UserAgent(), IP: s.clientIP(r),
	})
	writeAuth(w, status, result, apiErr)
}

func (s *Server) login(w http.ResponseWriter, r *http.Request) {
	if !s.allowAuth(w, r) {
		return
	}
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if !decodeJSON(w, r, &body) {
		return
	}
	result, apiErr, status := s.Accounts.Login(r.Context(), accounts.LoginInput{
		Email: body.Email, Password: body.Password, UserAgent: r.UserAgent(), IP: s.clientIP(r),
	})
	writeAuth(w, status, result, apiErr)
}

func (s *Server) logout(w http.ResponseWriter, r *http.Request) {
	if !s.requireAccounts(w) {
		return
	}
	if apiErr, status := s.Accounts.Logout(r.Context(), SessionToken(r), s.clientIP(r)); apiErr != nil {
		write(w, status, contract.Envelope{Error: apiErr})
		return
	}
	write(w, http.StatusOK, contract.Envelope{OK: true})
}

func (s *Server) me(w http.ResponseWriter, r *http.Request) {
	if !s.requireAccounts(w) {
		return
	}
	_, user, apiErr, status := s.Accounts.Require(r.Context(), SessionToken(r))
	if apiErr != nil {
		write(w, status, contract.Envelope{Error: apiErr})
		return
	}
	write(w, http.StatusOK, contract.Envelope{OK: true, User: user})
}

func (s *Server) verifyEmail(w http.ResponseWriter, r *http.Request) {
	if !s.requireAccounts(w) {
		return
	}
	var body struct {
		Token string `json:"token"`
	}
	if !decodeJSON(w, r, &body) {
		return
	}
	user, apiErr, status := s.Accounts.VerifyEmail(r.Context(), body.Token)
	if apiErr != nil {
		write(w, status, contract.Envelope{Error: apiErr})
		return
	}
	write(w, status, contract.Envelope{OK: true, User: user})
}

func (s *Server) resendVerification(w http.ResponseWriter, r *http.Request) {
	if !s.requireAccounts(w) {
		return
	}
	state, apiErr, status := s.Accounts.ResendVerification(r.Context(), SessionToken(r))
	if apiErr != nil {
		write(w, status, contract.Envelope{Error: apiErr})
		return
	}
	write(w, status, contract.Envelope{OK: true, Auth: state})
}

func (s *Server) forgotPassword(w http.ResponseWriter, r *http.Request) {
	if !s.allowAuth(w, r) {
		return
	}
	var body struct {
		Email string `json:"email"`
	}
	if !decodeJSON(w, r, &body) {
		return
	}
	state, apiErr, status := s.Accounts.ForgotPassword(r.Context(), body.Email, s.clientIP(r))
	if apiErr != nil {
		write(w, status, contract.Envelope{Error: apiErr})
		return
	}
	write(w, status, contract.Envelope{OK: true, Auth: state})
}

func (s *Server) resetPassword(w http.ResponseWriter, r *http.Request) {
	if !s.allowAuth(w, r) {
		return
	}
	var body struct {
		Token    string `json:"token"`
		Password string `json:"password"`
	}
	if !decodeJSON(w, r, &body) {
		return
	}
	if apiErr, status := s.Accounts.ResetPassword(r.Context(), body.Token, body.Password); apiErr != nil {
		write(w, status, contract.Envelope{Error: apiErr})
		return
	}
	write(w, http.StatusOK, contract.Envelope{OK: true})
}

func (s *Server) changePassword(w http.ResponseWriter, r *http.Request) {
	if !s.requireAccounts(w) {
		return
	}
	var body struct {
		CurrentPassword string `json:"currentPassword"`
		NewPassword     string `json:"newPassword"`
	}
	if !decodeJSON(w, r, &body) {
		return
	}
	if apiErr, status := s.Accounts.ChangePassword(r.Context(), SessionToken(r), body.CurrentPassword, body.NewPassword, s.clientIP(r)); apiErr != nil {
		write(w, status, contract.Envelope{Error: apiErr})
		return
	}
	write(w, http.StatusOK, contract.Envelope{OK: true})
}

func (s *Server) listSessions(w http.ResponseWriter, r *http.Request) {
	if !s.requireAccounts(w) {
		return
	}
	items, apiErr, status := s.Accounts.ListSessions(r.Context(), SessionToken(r))
	if apiErr != nil {
		write(w, status, contract.Envelope{Error: apiErr})
		return
	}
	write(w, status, contract.Envelope{OK: true, Sessions: items})
}

func (s *Server) revokeSession(w http.ResponseWriter, r *http.Request) {
	if !s.requireAccounts(w) {
		return
	}
	if apiErr, status := s.Accounts.RevokeSession(r.Context(), SessionToken(r), r.PathValue("id"), s.clientIP(r)); apiErr != nil {
		write(w, status, contract.Envelope{Error: apiErr})
		return
	}
	write(w, http.StatusOK, contract.Envelope{OK: true})
}

func (s *Server) revokeOtherSessions(w http.ResponseWriter, r *http.Request) {
	if !s.requireAccounts(w) {
		return
	}
	if apiErr, status := s.Accounts.RevokeOtherSessions(r.Context(), SessionToken(r), s.clientIP(r)); apiErr != nil {
		write(w, status, contract.Envelope{Error: apiErr})
		return
	}
	write(w, http.StatusOK, contract.Envelope{OK: true})
}

func (s *Server) listEvents(w http.ResponseWriter, r *http.Request) {
	if !s.requireAccounts(w) {
		return
	}
	items, apiErr, status := s.Accounts.ListEvents(r.Context(), SessionToken(r))
	if apiErr != nil {
		write(w, status, contract.Envelope{Error: apiErr})
		return
	}
	write(w, status, contract.Envelope{OK: true, Events: items})
}

func (s *Server) authMethods(w http.ResponseWriter, _ *http.Request) {
	methods := contract.EmptyAuthMethods()
	if s.Accounts != nil {
		methods = s.Accounts.Methods()
	}
	write(w, http.StatusOK, contract.Envelope{OK: true, Auth: &contract.AuthState{Methods: methods}})
}

func (s *Server) unsupportedFactor(w http.ResponseWriter, _ *http.Request) {
	if !s.requireAccounts(w) {
		return
	}
	apiErr, status := s.Accounts.UnsupportedFactor()
	write(w, status, contract.Envelope{Error: apiErr})
}

func (s *Server) allowAuth(w http.ResponseWriter, r *http.Request) bool {
	if s.Accounts == nil {
		write(w, http.StatusServiceUnavailable, contract.Envelope{Error: &contract.APIError{Code: "unavailable", Message: "Account service is unavailable."}})
		return false
	}
	ip := s.clientIP(r)
	if s.Limiter != nil && !s.Limiter.Allow("auth:"+ip, s.Cfg.RateLimitPerMin) {
		s.recordAbuse(r, "rate_limit", ip, "auth", "blocked", "Authentication rate limit exceeded.")
		write(w, http.StatusTooManyRequests, contract.Envelope{Error: &contract.APIError{Code: "rate_limited", Message: "Too many authentication requests"}})
		return false
	}
	return true
}

func writeAuth(w http.ResponseWriter, status int, result *accounts.AuthResult, apiErr *contract.APIError) {
	if apiErr != nil {
		write(w, status, contract.Envelope{Error: apiErr})
		return
	}
	session := result.Session
	write(w, status, contract.Envelope{
		OK:           true,
		User:         &result.User,
		Session:      &session,
		SessionToken: result.SessionToken,
		Auth:         &result.Auth,
	})
}

func decodeJSON(w http.ResponseWriter, r *http.Request, dest any) bool {
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 8<<10)).Decode(dest); err != nil {
		write(w, http.StatusBadRequest, contract.Envelope{Error: &contract.APIError{Code: "validation_error", Message: "Expected JSON."}})
		return false
	}
	return true
}

func (s *Server) requireAccounts(w http.ResponseWriter) bool {
	if s.Accounts != nil {
		return true
	}
	write(w, http.StatusServiceUnavailable, contract.Envelope{Error: &contract.APIError{Code: "unavailable", Message: "Account service is unavailable."}})
	return false
}
