package storage

import (
	"context"
	"errors"
	"time"

	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/contract"
)

type DiagnosisRecord struct {
	Diagnosis contract.Diagnosis
	Target    contract.Target
	UserID    string
}

type Job struct {
	DiagnosisID string
	Target      contract.Target
	QueuedAt    time.Time
}

type DiagnoseStore interface {
	CreateDiagnosis(ctx context.Context, rec DiagnosisRecord) error
	GetDiagnosis(ctx context.Context, id string) (*DiagnosisRecord, error)
	UpdateDiagnosis(ctx context.Context, rec DiagnosisRecord) error
	ListDiagnosesByUser(ctx context.Context, userID string) ([]DiagnosisRecord, error)
	DeleteDiagnosis(ctx context.Context, id string) error
	ListIncidents(ctx context.Context) ([]contract.Incident, error)
	GetIncident(ctx context.Context, id string) (*contract.Incident, error)
	ListServices(ctx context.Context) ([]contract.Service, error)
	GetService(ctx context.Context, slug string) (*contract.Service, error)
	Backend() string
}

type MeasurementStore interface {
	Record(ctx context.Context, diagnosisID string, measurements []contract.Measurement) error
	Backend() string
}

type Queue interface {
	Enqueue(ctx context.Context, job Job) error
	Dequeue(ctx context.Context) (Job, bool)
	Backend() string
}

type RateLimiter interface {
	Allow(key string, limitPerMin int) bool
	Backend() string
}

type Cache interface {
	Get(key string) (string, bool)
	Set(key, value string, ttl time.Duration)
	Backend() string
}

type UserRecord struct {
	User         contract.User
	PasswordHash string
	Alerts       contract.AlertPreferences
}

type TokenRecord struct {
	Kind      string
	UserID    string
	Hash      string
	ExpiresAt time.Time
	Used      bool
}

type ShareRecord struct {
	Hash        string
	DiagnosisID string
	UserID      string
	CreatedAt   time.Time
}

var ErrEmailTaken = errors.New("email already registered")

type AccountStore interface {
	CreateUser(ctx context.Context, rec UserRecord) error
	GetUserByID(ctx context.Context, id string) (*UserRecord, error)
	GetUserByEmail(ctx context.Context, email string) (*UserRecord, error)
	UpdateUser(ctx context.Context, rec UserRecord) error
	DeleteUser(ctx context.Context, id string) error
	CreateSession(ctx context.Context, session contract.Session) error
	GetSession(ctx context.Context, id string) (*contract.Session, error)
	GetSessionByTokenHash(ctx context.Context, hash string) (*contract.Session, error)
	ListSessions(ctx context.Context, userID string) ([]contract.Session, error)
	UpdateSession(ctx context.Context, session contract.Session) error
	RevokeSession(ctx context.Context, userID, sessionID string) (bool, error)
	RevokeAllSessions(ctx context.Context, userID string) error
	CreateToken(ctx context.Context, rec TokenRecord) error
	GetToken(ctx context.Context, kind, hash string) (*TokenRecord, error)
	MarkTokenUsed(ctx context.Context, kind, hash string) error
	AddEvent(ctx context.Context, event contract.SecurityEvent) error
	ListEvents(ctx context.Context, userID string) ([]contract.SecurityEvent, error)
	SaveService(ctx context.Context, userID string, item contract.SavedService) error
	ListSavedServices(ctx context.Context, userID string) ([]contract.SavedService, error)
	DeleteSavedService(ctx context.Context, userID, slug string) error
	CreateShare(ctx context.Context, rec ShareRecord) error
	GetShare(ctx context.Context, hash string) (*ShareRecord, error)
	ListSharesByUser(ctx context.Context, userID string) ([]ShareRecord, error)
	DeleteSharesForDiagnosis(ctx context.Context, diagnosisID string) error
}

type DeveloperStore interface {
	GetOrCreateWorkspace(ctx context.Context, ownerID, name string) (contract.Workspace, error)
	GetWorkspace(ctx context.Context, id string) (*contract.Workspace, error)
	GetWorkspaceByOwner(ctx context.Context, ownerID string) (*contract.Workspace, error)
	CreateMonitor(ctx context.Context, item contract.Monitor) error
	GetMonitor(ctx context.Context, id string) (*contract.Monitor, error)
	ListMonitors(ctx context.Context, workspaceID string) ([]contract.Monitor, error)
	UpdateMonitor(ctx context.Context, item contract.Monitor) error
	DeleteMonitor(ctx context.Context, workspaceID, id string) (bool, error)
	AddCheck(ctx context.Context, check contract.MonitorCheck) error
	ListChecks(ctx context.Context, workspaceID, monitorID string) ([]contract.MonitorCheck, error)
	CreateAPIKey(ctx context.Context, key contract.APIKey) error
	GetAPIKey(ctx context.Context, id string) (*contract.APIKey, error)
	GetAPIKeyByHash(ctx context.Context, hash string) (*contract.APIKey, error)
	ListAPIKeys(ctx context.Context, workspaceID string) ([]contract.APIKey, error)
	UpdateAPIKey(ctx context.Context, key contract.APIKey) error
	CreateWebhook(ctx context.Context, hook contract.Webhook) error
	GetWebhook(ctx context.Context, id string) (*contract.Webhook, error)
	ListWebhooks(ctx context.Context, workspaceID string) ([]contract.Webhook, error)
	UpdateWebhook(ctx context.Context, hook contract.Webhook) error
	DeleteWebhook(ctx context.Context, workspaceID, id string) (bool, error)
	CreateDelivery(ctx context.Context, item contract.WebhookDelivery) error
	GetDeliveryByIdempotency(ctx context.Context, key string) (*contract.WebhookDelivery, error)
	ListDeliveries(ctx context.Context, workspaceID, webhookID string) ([]contract.WebhookDelivery, error)
	UpdateDelivery(ctx context.Context, item contract.WebhookDelivery) error
	CreateAlertRule(ctx context.Context, rule contract.AlertRule) error
	ListAlertRules(ctx context.Context, workspaceID string) ([]contract.AlertRule, error)
	UpdateAlertRule(ctx context.Context, rule contract.AlertRule) error
	DeleteAlertRule(ctx context.Context, workspaceID, id string) (bool, error)
	CreateDevIncident(ctx context.Context, item contract.DeveloperIncident) error
	ListDevIncidents(ctx context.Context, workspaceID string) ([]contract.DeveloperIncident, error)
	GetDevIncident(ctx context.Context, id string) (*contract.DeveloperIncident, error)
	UpdateDevIncident(ctx context.Context, item contract.DeveloperIncident) error
	IncrUsage(ctx context.Context, workspaceID, field string, delta int) error
	GetUsage(ctx context.Context, workspaceID string) (contract.Usage, error)
}
