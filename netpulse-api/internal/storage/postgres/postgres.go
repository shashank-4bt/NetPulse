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
