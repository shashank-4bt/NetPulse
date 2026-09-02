package accounts

import (
	"context"
	"errors"
	"sort"
	"strings"
	"time"

	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/auth"
	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/contract"
	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/id"
	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/storage"
)

const (
	TokenVerify = "verify_email"
	TokenReset  = "reset_password"
)

type Service struct {
	Accounts  storage.AccountStore
	Diagnoses storage.DiagnoseStore
	DevTokens bool
	SessionTTL time.Duration
	Now       func() time.Time
}

type RegisterInput struct {
	Email       string
	Password    string
	DisplayName string
	UserAgent   string
	IP          string
}

type LoginInput struct {
	Email     string
	Password  string
	UserAgent string
	IP        string
}

type AuthResult struct {
	User         contract.User
	Session      contract.Session
	SessionToken string
	Auth         contract.AuthState
}

func (s *Service) Register(ctx context.Context, in RegisterInput) (*AuthResult, *contract.APIError, int) {
	email := auth.NormalizeEmail(in.Email)
	if msg := auth.ValidateEmail(email); msg != "" {
		return nil, apiErr("validation_error", msg), 400
	}
	if msg := auth.ValidatePassword(in.Password, email); msg != "" {
		return nil, apiErr("validation_error", msg), 400
	}
	hash, err := auth.HashPassword(in.Password)
	if err != nil {
		return nil, apiErr("internal", "Could not store the password."), 500
	}
	now := s.now()
	user := contract.User{
		ID:             id.New(),
		Email:          email,
		DisplayName:    strings.TrimSpace(in.DisplayName),
		EmailVerified:  false,
		CreatedAt:      now.UTC().Format(time.RFC3339),
		TelemetryOptIn: false,
	}
	if user.DisplayName == "" {
		user.DisplayName = strings.Split(email, "@")[0]
	}
	if err := s.Accounts.CreateUser(ctx, storage.UserRecord{
		User:         user,
		PasswordHash: hash,
		Alerts:       contract.EmptyAlerts(),
	}); err != nil {
		if errors.Is(err, storage.ErrEmailTaken) {
			return nil, apiErr("validation_error", "An account with this email already exists."), 409
		}
		return nil, apiErr("unavailable", "Account store is unavailable."), 503
	}
	session, token, err := s.createSession(ctx, user.ID, in.UserAgent, in.IP, now)
	if err != nil {
		return nil, apiErr("unavailable", "Session store is unavailable."), 503
	}
	devToken := s.issueToken(ctx, TokenVerify, user.ID, now.Add(24*time.Hour))
	_ = s.Accounts.AddEvent(ctx, contract.SecurityEvent{
		ID: id.New(), UserID: user.ID, Kind: "registered", At: now.UTC().Format(time.RFC3339),
		Summary: "Account created. Email is not verified.", IP: auth.CoarseIP(in.IP),
	})
	return &AuthResult{User: user, Session: *session, SessionToken: token, Auth: s.authState(devToken)}, nil, 201
}

func (s *Service) Login(ctx context.Context, in LoginInput) (*AuthResult, *contract.APIError, int) {
	email := auth.NormalizeEmail(in.Email)
	rec, err := s.Accounts.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, apiErr("unavailable", "Account store is unavailable."), 503
	}
	if rec == nil || !auth.CheckPassword(rec.PasswordHash, in.Password) {
		if rec != nil {
			_ = s.Accounts.AddEvent(ctx, contract.SecurityEvent{
				ID: id.New(), UserID: rec.User.ID, Kind: "login_failed",
				At: s.now().UTC().Format(time.RFC3339), Summary: "Sign-in failed.", IP: auth.CoarseIP(in.IP),
			})
		}
		return nil, apiErr("unauthorized", "Email or password is incorrect."), 401
	}
	now := s.now()
	session, token, err := s.createSession(ctx, rec.User.ID, in.UserAgent, in.IP, now)
	if err != nil {
		return nil, apiErr("unavailable", "Session store is unavailable."), 503
	}
	_ = s.Accounts.AddEvent(ctx, contract.SecurityEvent{
		ID: id.New(), UserID: rec.User.ID, Kind: "logged_in",
		At: now.UTC().Format(time.RFC3339), Summary: "Signed in.", IP: auth.CoarseIP(in.IP),
	})
	return &AuthResult{User: rec.User, Session: *session, SessionToken: token, Auth: s.authState("")}, nil, 200
}

func (s *Service) Logout(ctx context.Context, sessionID, ip string) (*contract.APIError, int) {
	session, user, errResp, status := s.Require(ctx, sessionID)
	if errResp != nil {
		return errResp, status
	}
	_, err := s.Accounts.RevokeSession(ctx, user.ID, session.ID)
	if err != nil {
		return apiErr("unavailable", "Session store is unavailable."), 503
	}
	_ = s.Accounts.AddEvent(ctx, contract.SecurityEvent{
		ID: id.New(), UserID: user.ID, Kind: "logged_out",
		At: s.now().UTC().Format(time.RFC3339), Summary: "Signed out.", IP: auth.CoarseIP(ip),
	})
	return nil, 200
}

func (s *Service) Require(ctx context.Context, sessionToken string) (*contract.Session, *contract.User, *contract.APIError, int) {
	if sessionToken == "" {
		return nil, nil, apiErr("unauthorized", "Sign in required."), 401
	}
	session, err := s.Accounts.GetSessionByTokenHash(ctx, auth.HashSecret(sessionToken))
	if err != nil {
		return nil, nil, apiErr("unavailable", "Session store is unavailable."), 503
	}
	if session == nil || session.Revoked {
		return nil, nil, apiErr("unauthorized", "Session is not valid."), 401
	}
	exp, err := time.Parse(time.RFC3339, session.ExpiresAt)
	if err != nil || !s.now().Before(exp) {
		return nil, nil, apiErr("unauthorized", "Session is not valid."), 401
	}
	rec, err := s.Accounts.GetUserByID(ctx, session.UserID)
	if err != nil {
		return nil, nil, apiErr("unavailable", "Account store is unavailable."), 503
	}
	if rec == nil {
		return nil, nil, apiErr("unauthorized", "Session is not valid."), 401
	}
	session.LastSeenAt = s.now().UTC().Format(time.RFC3339)
	_ = s.Accounts.UpdateSession(ctx, *session)
	user := rec.User
	return session, &user, nil, 200
}

func (s *Service) VerifyEmail(ctx context.Context, raw string) (*contract.User, *contract.APIError, int) {
	rec, errResp, status := s.consumeToken(ctx, TokenVerify, raw)
	if errResp != nil {
		return nil, errResp, status
	}
	userRec, err := s.Accounts.GetUserByID(ctx, rec.UserID)
	if err != nil || userRec == nil {
		return nil, apiErr("not_found", "Verification target was not found."), 404
	}
	userRec.User.EmailVerified = true
	if err := s.Accounts.UpdateUser(ctx, *userRec); err != nil {
		return nil, apiErr("unavailable", "Account store is unavailable."), 503
	}
	_ = s.Accounts.AddEvent(ctx, contract.SecurityEvent{
		ID: id.New(), UserID: userRec.User.ID, Kind: "email_verified",
		At: s.now().UTC().Format(time.RFC3339), Summary: "Email verified.",
	})
	return &userRec.User, nil, 200
}

func (s *Service) ResendVerification(ctx context.Context, sessionID string) (*contract.AuthState, *contract.APIError, int) {
	_, user, errResp, status := s.Require(ctx, sessionID)
	if errResp != nil {
		return nil, errResp, status
	}
	if user.EmailVerified {
		state := s.authState("")
		state.EmailReason = "Email is already verified."
		return &state, nil, 200
	}
	token := s.issueToken(ctx, TokenVerify, user.ID, s.now().Add(24*time.Hour))
	state := s.authState(token)
	return &state, nil, 200
}

func (s *Service) ForgotPassword(ctx context.Context, email, ip string) (*contract.AuthState, *contract.APIError, int) {
	state := s.authState("")
	state.EmailReason = "If an account exists, a reset path was created. Email delivery is not configured."
	rec, err := s.Accounts.GetUserByEmail(ctx, auth.NormalizeEmail(email))
	if err != nil {
		return nil, apiErr("unavailable", "Account store is unavailable."), 503
	}
	if rec == nil {
		return &state, nil, 200
	}
	token := s.issueToken(ctx, TokenReset, rec.User.ID, s.now().Add(time.Hour))
	_ = s.Accounts.AddEvent(ctx, contract.SecurityEvent{
		ID: id.New(), UserID: rec.User.ID, Kind: "password_reset_requested",
		At: s.now().UTC().Format(time.RFC3339), Summary: "Password reset requested.", IP: auth.CoarseIP(ip),
	})
	if s.DevTokens {
		state.DevToken = token
	}
	return &state, nil, 200
}

func (s *Service) ResetPassword(ctx context.Context, raw, password string) (*contract.APIError, int) {
	rec, errResp, status := s.consumeToken(ctx, TokenReset, raw)
	if errResp != nil {
		return errResp, status
	}
	userRec, err := s.Accounts.GetUserByID(ctx, rec.UserID)
	if err != nil || userRec == nil {
		return apiErr("not_found", "Reset target was not found."), 404
	}
	if msg := auth.ValidatePassword(password, userRec.User.Email); msg != "" {
		return apiErr("validation_error", msg), 400
	}
	hash, err := auth.HashPassword(password)
	if err != nil {
		return apiErr("internal", "Could not store the password."), 500
	}
	userRec.PasswordHash = hash
	if err := s.Accounts.UpdateUser(ctx, *userRec); err != nil {
		return apiErr("unavailable", "Account store is unavailable."), 503
	}
	_ = s.Accounts.RevokeAllSessions(ctx, userRec.User.ID)
	_ = s.Accounts.AddEvent(ctx, contract.SecurityEvent{
		ID: id.New(), UserID: userRec.User.ID, Kind: "password_reset",
		At: s.now().UTC().Format(time.RFC3339), Summary: "Password was reset. Other sessions were revoked.",
	})
	return nil, 200
}

func (s *Service) ChangePassword(ctx context.Context, sessionID, current, next, ip string) (*contract.APIError, int) {
	_, user, errResp, status := s.Require(ctx, sessionID)
	if errResp != nil {
		return errResp, status
	}
	rec, err := s.Accounts.GetUserByID(ctx, user.ID)
	if err != nil || rec == nil {
		return apiErr("unavailable", "Account store is unavailable."), 503
	}
	if !auth.CheckPassword(rec.PasswordHash, current) {
		return apiErr("unauthorized", "Current password is incorrect."), 401
	}
	if msg := auth.ValidatePassword(next, user.Email); msg != "" {
		return apiErr("validation_error", msg), 400
	}
	hash, err := auth.HashPassword(next)
	if err != nil {
		return apiErr("internal", "Could not store the password."), 500
	}
	rec.PasswordHash = hash
	if err := s.Accounts.UpdateUser(ctx, *rec); err != nil {
		return apiErr("unavailable", "Account store is unavailable."), 503
	}
	_ = s.Accounts.AddEvent(ctx, contract.SecurityEvent{
		ID: id.New(), UserID: user.ID, Kind: "password_changed",
		At: s.now().UTC().Format(time.RFC3339), Summary: "Password changed.", IP: auth.CoarseIP(ip),
	})
	return nil, 200
}

func (s *Service) ListSessions(ctx context.Context, sessionID string) ([]contract.Session, *contract.APIError, int) {
	current, user, errResp, status := s.Require(ctx, sessionID)
	if errResp != nil {
		return nil, errResp, status
	}
	items, err := s.Accounts.ListSessions(ctx, user.ID)
	if err != nil {
		return nil, apiErr("unavailable", "Session store is unavailable."), 503
	}
	active := []contract.Session{}
	for _, item := range items {
		if item.Revoked {
			continue
		}
		exp, err := time.Parse(time.RFC3339, item.ExpiresAt)
		if err != nil || !s.now().Before(exp) {
			continue
		}
		item.Current = item.ID == current.ID
		active = append(active, item)
	}
	sort.SliceStable(active, func(i, j int) bool { return active[i].LastSeenAt > active[j].LastSeenAt })
	return active, nil, 200
}

func (s *Service) RevokeSession(ctx context.Context, sessionID, targetID, ip string) (*contract.APIError, int) {
	_, user, errResp, status := s.Require(ctx, sessionID)
	if errResp != nil {
		return errResp, status
	}
	ok, err := s.Accounts.RevokeSession(ctx, user.ID, targetID)
	if err != nil {
		return apiErr("unavailable", "Session store is unavailable."), 503
	}
	if !ok {
		return apiErr("not_found", "Session was not found."), 404
	}
	_ = s.Accounts.AddEvent(ctx, contract.SecurityEvent{
		ID: id.New(), UserID: user.ID, Kind: "session_revoked",
		At: s.now().UTC().Format(time.RFC3339), Summary: "A session was revoked.", IP: auth.CoarseIP(ip),
	})
	return nil, 200
}

func (s *Service) RevokeOtherSessions(ctx context.Context, sessionToken, ip string) (*contract.APIError, int) {
	current, user, errResp, status := s.Require(ctx, sessionToken)
	if errResp != nil {
		return errResp, status
	}
	items, err := s.Accounts.ListSessions(ctx, user.ID)
	if err != nil {
		return apiErr("unavailable", "Session store is unavailable."), 503
	}
	for _, item := range items {
		if item.ID == current.ID || item.Revoked {
			continue
		}
		if _, err := s.Accounts.RevokeSession(ctx, user.ID, item.ID); err != nil {
			return apiErr("unavailable", "Session store is unavailable."), 503
		}
	}
	_ = s.Accounts.AddEvent(ctx, contract.SecurityEvent{
		ID: id.New(), UserID: user.ID, Kind: "sessions_revoked",
		At: s.now().UTC().Format(time.RFC3339), Summary: "Other sessions were revoked.", IP: auth.CoarseIP(ip),
	})
	return nil, 200
}

func (s *Service) ListEvents(ctx context.Context, sessionID string) ([]contract.SecurityEvent, *contract.APIError, int) {
	_, user, errResp, status := s.Require(ctx, sessionID)
	if errResp != nil {
		return nil, errResp, status
	}
	items, err := s.Accounts.ListEvents(ctx, user.ID)
	if err != nil {
		return nil, apiErr("unavailable", "Account store is unavailable."), 503
	}
	sort.SliceStable(items, func(i, j int) bool { return items[i].At > items[j].At })
	if len(items) > 50 {
		items = items[:50]
	}
	return items, nil, 200
}

func (s *Service) UpdateProfile(ctx context.Context, sessionID, displayName string) (*contract.User, *contract.APIError, int) {
	_, user, errResp, status := s.Require(ctx, sessionID)
	if errResp != nil {
		return nil, errResp, status
	}
	rec, err := s.Accounts.GetUserByID(ctx, user.ID)
	if err != nil || rec == nil {
		return nil, apiErr("unavailable", "Account store is unavailable."), 503
	}
	name := strings.TrimSpace(displayName)
	if name == "" {
		return nil, apiErr("validation_error", "Display name is required."), 400
	}
	if len(name) > 80 {
		return nil, apiErr("validation_error", "Display name is too long."), 400
	}
	rec.User.DisplayName = name
	if err := s.Accounts.UpdateUser(ctx, *rec); err != nil {
		return nil, apiErr("unavailable", "Account store is unavailable."), 503
	}
	return &rec.User, nil, 200
}

func (s *Service) Privacy(ctx context.Context, sessionID string) (*contract.PrivacySettings, *contract.APIError, int) {
	_, user, errResp, status := s.Require(ctx, sessionID)
	if errResp != nil {
		return nil, errResp, status
	}
	settings := contract.DefaultPrivacy()
	settings.TelemetryOptIn = user.TelemetryOptIn
	return &settings, nil, 200
}

func (s *Service) UpdatePrivacy(ctx context.Context, sessionID string, telemetry bool) (*contract.PrivacySettings, *contract.APIError, int) {
	_, user, errResp, status := s.Require(ctx, sessionID)
	if errResp != nil {
		return nil, errResp, status
	}
	rec, err := s.Accounts.GetUserByID(ctx, user.ID)
	if err != nil || rec == nil {
		return nil, apiErr("unavailable", "Account store is unavailable."), 503
	}
	rec.User.TelemetryOptIn = telemetry
	if err := s.Accounts.UpdateUser(ctx, *rec); err != nil {
		return nil, apiErr("unavailable", "Account store is unavailable."), 503
	}
	settings := contract.DefaultPrivacy()
	settings.TelemetryOptIn = telemetry
	return &settings, nil, 200
}

func (s *Service) DeleteAccount(ctx context.Context, sessionID string) (*contract.APIError, int) {
	_, user, errResp, status := s.Require(ctx, sessionID)
	if errResp != nil {
		return errResp, status
	}
	if err := s.Accounts.DeleteUser(ctx, user.ID); err != nil {
		return apiErr("unavailable", "Account store is unavailable."), 503
	}
	return nil, 200
}

func (s *Service) Alerts(ctx context.Context, sessionID string) (*contract.AlertPreferences, *contract.APIError, int) {
	_, user, errResp, status := s.Require(ctx, sessionID)
	if errResp != nil {
		return nil, errResp, status
	}
	rec, err := s.Accounts.GetUserByID(ctx, user.ID)
	if err != nil || rec == nil {
		return nil, apiErr("unavailable", "Account store is unavailable."), 503
	}
	alerts := rec.Alerts
	if alerts.Summary == "" {
		alerts = contract.EmptyAlerts()
	}
	return &alerts, nil, 200
}

func (s *Service) UpdateAlerts(ctx context.Context, sessionID string, emailEnabled, incidents bool) (*contract.AlertPreferences, *contract.APIError, int) {
	_, user, errResp, status := s.Require(ctx, sessionID)
	if errResp != nil {
		return nil, errResp, status
	}
	rec, err := s.Accounts.GetUserByID(ctx, user.ID)
	if err != nil || rec == nil {
		return nil, apiErr("unavailable", "Account store is unavailable."), 503
	}
	rec.Alerts.EmailEnabled = emailEnabled
	rec.Alerts.IncidentAlerts = incidents
	rec.Alerts.DeliveredCount = 0
	rec.Alerts.Summary = "Preferences saved. No alerts have been delivered. Notification delivery is not configured."
	if err := s.Accounts.UpdateUser(ctx, *rec); err != nil {
		return nil, apiErr("unavailable", "Account store is unavailable."), 503
	}
	return &rec.Alerts, nil, 200
}

func (s *Service) Billing(ctx context.Context, sessionID string) (*contract.Billing, *contract.APIError, int) {
	if _, _, errResp, status := s.Require(ctx, sessionID); errResp != nil {
		return nil, errResp, status
	}
	billing := contract.EmptyBilling()
	return &billing, nil, 200
}

func (s *Service) Organization(ctx context.Context, sessionID, orgID string) (*contract.APIError, int) {
	if _, _, errResp, status := s.Require(ctx, sessionID); errResp != nil {
		return errResp, status
	}
	return apiErr("not_found", "No organization exists with that id."), 404
}

func (s *Service) ForeignBilling(ctx context.Context, sessionID, userID string) (*contract.APIError, int) {
	_, user, errResp, status := s.Require(ctx, sessionID)
	if errResp != nil {
		return errResp, status
	}
	if user.ID != userID {
		return apiErr("not_found", "Billing information was not found."), 404
	}
	return nil, 200
}

func (s *Service) Methods() contract.AuthMethods {
	return contract.EmptyAuthMethods()
}

func (s *Service) UnsupportedFactor() (*contract.APIError, int) {
	return apiErr("unavailable", "OAuth, passkeys, and MFA are not enabled. No provider is configured."), 501
}

func (s *Service) createSession(ctx context.Context, userID, userAgent, ip string, now time.Time) (*contract.Session, string, error) {
	ttl := s.SessionTTL
	if ttl <= 0 {
		ttl = 7 * 24 * time.Hour
	}
	secret := auth.NewSecret()
	session := contract.Session{
		ID:         id.New(),
		UserID:     userID,
		TokenHash:  auth.HashSecret(secret),
		CreatedAt:  now.UTC().Format(time.RFC3339),
		LastSeenAt: now.UTC().Format(time.RFC3339),
		ExpiresAt:  now.Add(ttl).UTC().Format(time.RFC3339),
		UserAgent:  strings.TrimSpace(userAgent),
		IP:         auth.CoarseIP(ip),
		Label:      auth.SessionLabel(userAgent),
		Current:    true,
	}
	if err := s.Accounts.CreateSession(ctx, session); err != nil {
		return nil, "", err
	}
	return &session, secret, nil
}

func (s *Service) issueToken(ctx context.Context, kind, userID string, expires time.Time) string {
	raw := auth.NewSecret()
	_ = s.Accounts.CreateToken(ctx, storage.TokenRecord{
		Kind: kind, UserID: userID, Hash: auth.HashSecret(raw), ExpiresAt: expires,
	})
	if s.DevTokens {
		return raw
	}
	return ""
}

func (s *Service) consumeToken(ctx context.Context, kind, raw string) (*storage.TokenRecord, *contract.APIError, int) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, apiErr("validation_error", "A token is required."), 400
	}
	rec, err := s.Accounts.GetToken(ctx, kind, auth.HashSecret(raw))
	if err != nil {
		return nil, apiErr("unavailable", "Account store is unavailable."), 503
	}
	if rec == nil || rec.Used || !s.now().Before(rec.ExpiresAt) {
		return nil, apiErr("not_found", "This token is not valid."), 404
	}
	if err := s.Accounts.MarkTokenUsed(ctx, kind, rec.Hash); err != nil {
		return nil, apiErr("unavailable", "Account store is unavailable."), 503
	}
	return rec, nil, 200
}

func (s *Service) authState(devToken string) contract.AuthState {
	return contract.AuthState{
		EmailSent:   false,
		EmailReason: "Email delivery is not configured. A message was not sent.",
		DevToken:    devToken,
		Methods:     contract.EmptyAuthMethods(),
	}
}

func (s *Service) now() time.Time {
	if s.Now != nil {
		return s.Now()
	}
	return time.Now()
}

func apiErr(code, message string) *contract.APIError {
	return &contract.APIError{Code: code, Message: message}
}
