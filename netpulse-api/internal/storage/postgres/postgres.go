package postgres

import (
	"context"
	"errors"

	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/contract"
	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/storage"
)

var ErrNotConfigured = errors.New("postgres is not configured")

// Store is the PostgreSQL adapter. Domain code depends on storage.DiagnoseStore,
// not on this type. A driver is linked only when a DSN is supplied by cmd/api.
type Store struct {
	dsn string
}

func Open(dsn string) (*Store, error) {
	if dsn == "" {
		return nil, ErrNotConfigured
	}
	return nil, errors.New("postgres DSN is set but no SQL driver is linked in this build; keep NETPULSE_DATABASE_URL empty to use the memory adapter")
}

func (s *Store) CreateDiagnosis(context.Context, storage.DiagnosisRecord) error {
	return ErrNotConfigured
}
func (s *Store) GetDiagnosis(context.Context, string) (*storage.DiagnosisRecord, error) {
	return nil, ErrNotConfigured
}
func (s *Store) UpdateDiagnosis(context.Context, storage.DiagnosisRecord) error {
	return ErrNotConfigured
}
func (s *Store) ListDiagnosesByUser(context.Context, string) ([]storage.DiagnosisRecord, error) {
	return nil, ErrNotConfigured
}
func (s *Store) DeleteDiagnosis(context.Context, string) error {
	return ErrNotConfigured
}
func (s *Store) ListIncidents(context.Context) ([]contract.Incident, error) {
	return nil, ErrNotConfigured
}
func (s *Store) GetIncident(context.Context, string) (*contract.Incident, error) {
	return nil, ErrNotConfigured
}
func (s *Store) ListServices(context.Context) ([]contract.Service, error) {
	return nil, ErrNotConfigured
}
func (s *Store) GetService(context.Context, string) (*contract.Service, error) {
	return nil, ErrNotConfigured
}
func (s *Store) Backend() string { return "postgres-unlinked" }

func (s *Store) CreateUser(context.Context, storage.UserRecord) error {
	return ErrNotConfigured
}
func (s *Store) GetUserByID(context.Context, string) (*storage.UserRecord, error) {
	return nil, ErrNotConfigured
}
func (s *Store) GetUserByEmail(context.Context, string) (*storage.UserRecord, error) {
	return nil, ErrNotConfigured
}
func (s *Store) UpdateUser(context.Context, storage.UserRecord) error {
	return ErrNotConfigured
}
func (s *Store) DeleteUser(context.Context, string) error { return ErrNotConfigured }
func (s *Store) CreateSession(context.Context, contract.Session) error {
	return ErrNotConfigured
}
func (s *Store) GetSession(context.Context, string) (*contract.Session, error) {
	return nil, ErrNotConfigured
}
func (s *Store) GetSessionByTokenHash(context.Context, string) (*contract.Session, error) {
	return nil, ErrNotConfigured
}
func (s *Store) ListSessions(context.Context, string) ([]contract.Session, error) {
	return nil, ErrNotConfigured
}
func (s *Store) UpdateSession(context.Context, contract.Session) error {
	return ErrNotConfigured
}
func (s *Store) RevokeSession(context.Context, string, string) (bool, error) {
	return false, ErrNotConfigured
}
func (s *Store) RevokeAllSessions(context.Context, string) error { return ErrNotConfigured }
func (s *Store) CreateToken(context.Context, storage.TokenRecord) error {
	return ErrNotConfigured
}
func (s *Store) GetToken(context.Context, string, string) (*storage.TokenRecord, error) {
	return nil, ErrNotConfigured
}
func (s *Store) MarkTokenUsed(context.Context, string, string) error { return ErrNotConfigured }
func (s *Store) AddEvent(context.Context, contract.SecurityEvent) error {
	return ErrNotConfigured
}
func (s *Store) ListEvents(context.Context, string) ([]contract.SecurityEvent, error) {
	return nil, ErrNotConfigured
}
func (s *Store) SaveService(context.Context, string, contract.SavedService) error {
	return ErrNotConfigured
}
func (s *Store) ListSavedServices(context.Context, string) ([]contract.SavedService, error) {
	return nil, ErrNotConfigured
}
func (s *Store) DeleteSavedService(context.Context, string, string) error {
	return ErrNotConfigured
}
func (s *Store) CreateShare(context.Context, storage.ShareRecord) error {
	return ErrNotConfigured
}
func (s *Store) GetShare(context.Context, string) (*storage.ShareRecord, error) {
	return nil, ErrNotConfigured
}
func (s *Store) ListSharesByUser(context.Context, string) ([]storage.ShareRecord, error) {
	return nil, ErrNotConfigured
}
func (s *Store) DeleteSharesForDiagnosis(context.Context, string) error {
	return ErrNotConfigured
}
