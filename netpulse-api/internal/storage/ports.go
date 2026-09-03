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
	Depth() int
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

type BusinessStore interface {
	CreateOrg(ctx context.Context, org contract.Organization) error
	GetOrg(ctx context.Context, id string) (*contract.Organization, error)
	UpdateOrg(ctx context.Context, org contract.Organization) error
	ListOrgsForUser(ctx context.Context, userID string) ([]contract.Organization, error)
	CreateMember(ctx context.Context, member contract.Member) error
	GetMember(ctx context.Context, orgID, memberID string) (*contract.Member, error)
	GetMemberByUser(ctx context.Context, orgID, userID string) (*contract.Member, error)
	ListMembers(ctx context.Context, orgID string) ([]contract.Member, error)
	UpdateMember(ctx context.Context, member contract.Member) error
	DeleteMember(ctx context.Context, orgID, memberID string) (bool, error)
	CreateInvite(ctx context.Context, invite contract.OrgInvite) error
	ListInvites(ctx context.Context, orgID string) ([]contract.OrgInvite, error)
	ListInvitesByEmail(ctx context.Context, email string) ([]contract.OrgInvite, error)
	DeleteInvite(ctx context.Context, orgID, inviteID string) (bool, error)
	CreateTeam(ctx context.Context, team contract.Team) error
	GetTeam(ctx context.Context, id string) (*contract.Team, error)
	ListTeams(ctx context.Context, orgID string) ([]contract.Team, error)
	UpdateTeam(ctx context.Context, team contract.Team) error
	DeleteTeam(ctx context.Context, orgID, id string) (bool, error)
	CreateOrgDevice(ctx context.Context, item contract.OrgDevice) error
	GetOrgDevice(ctx context.Context, id string) (*contract.OrgDevice, error)
	ListOrgDevices(ctx context.Context, orgID string) ([]contract.OrgDevice, error)
	UpdateOrgDevice(ctx context.Context, item contract.OrgDevice) error
	DeleteOrgDevice(ctx context.Context, orgID, id string) (bool, error)
	CreateOrgNetwork(ctx context.Context, item contract.OrgNetwork) error
	GetOrgNetwork(ctx context.Context, id string) (*contract.OrgNetwork, error)
	ListOrgNetworks(ctx context.Context, orgID string) ([]contract.OrgNetwork, error)
	UpdateOrgNetwork(ctx context.Context, item contract.OrgNetwork) error
	DeleteOrgNetwork(ctx context.Context, orgID, id string) (bool, error)
	CreateOrgService(ctx context.Context, item contract.OrgService) error
	GetOrgService(ctx context.Context, id string) (*contract.OrgService, error)
	ListOrgServices(ctx context.Context, orgID string) ([]contract.OrgService, error)
	UpdateOrgService(ctx context.Context, item contract.OrgService) error
	DeleteOrgService(ctx context.Context, orgID, id string) (bool, error)
	CreateOrgMonitor(ctx context.Context, item contract.Monitor) error
	GetOrgMonitor(ctx context.Context, id string) (*contract.Monitor, error)
	ListOrgMonitors(ctx context.Context, orgID string) ([]contract.Monitor, error)
	UpdateOrgMonitor(ctx context.Context, item contract.Monitor) error
	DeleteOrgMonitor(ctx context.Context, orgID, id string) (bool, error)
	AddOrgCheck(ctx context.Context, check contract.MonitorCheck) error
	ListOrgChecks(ctx context.Context, orgID, monitorID string) ([]contract.MonitorCheck, error)
	CreateOrgIncident(ctx context.Context, item contract.OrgIncident) error
	GetOrgIncident(ctx context.Context, id string) (*contract.OrgIncident, error)
	ListOrgIncidents(ctx context.Context, orgID string) ([]contract.OrgIncident, error)
	UpdateOrgIncident(ctx context.Context, item contract.OrgIncident) error
	CreateOrgDiagnosis(ctx context.Context, item contract.OrgDiagnosis) error
	ListOrgDiagnoses(ctx context.Context, orgID string) ([]contract.OrgDiagnosis, error)
	GetOrgDiagnosis(ctx context.Context, id string) (*contract.OrgDiagnosis, error)
	CreateOrgReport(ctx context.Context, item contract.OrgReport) error
	GetOrgReport(ctx context.Context, id string) (*contract.OrgReport, error)
	ListOrgReports(ctx context.Context, orgID string) ([]contract.OrgReport, error)
	CreateOrgKey(ctx context.Context, key contract.APIKey) error
	GetOrgKey(ctx context.Context, id string) (*contract.APIKey, error)
	GetOrgKeyByHash(ctx context.Context, hash string) (*contract.APIKey, error)
	ListOrgKeys(ctx context.Context, orgID string) ([]contract.APIKey, error)
	UpdateOrgKey(ctx context.Context, key contract.APIKey) error
	AddAudit(ctx context.Context, event contract.AuditEvent) error
	ListAudit(ctx context.Context, orgID string) ([]contract.AuditEvent, error)
}

type AdminStore interface {
	GetOperator(ctx context.Context, userID string) (*contract.Operator, error)
	UpsertOperator(ctx context.Context, op contract.Operator) error
	ListUsers(ctx context.Context) ([]contract.AdminUser, error)
	GetAdminUser(ctx context.Context, id string) (*contract.AdminUser, error)
	ListAllOrgs(ctx context.Context) ([]contract.Organization, error)
	ListAllDiagnoses(ctx context.Context) ([]DiagnosisRecord, error)
	ListAllMeasurements(ctx context.Context) ([]contract.AdminMeasurement, error)
	ListAllChecks(ctx context.Context) ([]contract.MonitorCheck, error)
	UpdateIncident(ctx context.Context, item contract.Incident) error
	AddIncidentNote(ctx context.Context, note contract.IncidentNote) error
	ListIncidentNotes(ctx context.Context, incidentID string) ([]contract.IncidentNote, error)
	SetIncidentOverride(ctx context.Context, incidentID string, ov contract.IncidentOverride) error
	GetIncidentOverride(ctx context.Context, incidentID string) (*contract.IncidentOverride, error)
	AddAbuse(ctx context.Context, ev contract.AbuseEvent) error
	ListAbuse(ctx context.Context) ([]contract.AbuseEvent, error)
	AddAdminAudit(ctx context.Context, ev contract.AdminAudit) error
	ListAdminAudit(ctx context.Context) ([]contract.AdminAudit, error)
	ListFlags(ctx context.Context) ([]contract.FeatureFlag, error)
	GetFlag(ctx context.Context, id string) (*contract.FeatureFlag, error)
	UpsertFlag(ctx context.Context, flag contract.FeatureFlag) error
	ListRemoteConfig(ctx context.Context) ([]contract.RemoteConfigEntry, error)
	GetRemoteConfig(ctx context.Context, key string) (*contract.RemoteConfigEntry, error)
	UpsertRemoteConfig(ctx context.Context, entry contract.RemoteConfigEntry) error
	AddRuleLabel(ctx context.Context, label contract.RuleLabel) error
	ListRuleLabels(ctx context.Context) ([]contract.RuleLabel, error)
}
