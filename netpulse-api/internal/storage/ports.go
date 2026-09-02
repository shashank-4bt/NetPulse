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
