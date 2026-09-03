package postgres

import (
	"context"

	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/contract"
	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/storage"
)

func (s *Store) GetOperator(context.Context, string) (*contract.Operator, error) {
	return nil, ErrNotConfigured
}
func (s *Store) UpsertOperator(context.Context, contract.Operator) error { return ErrNotConfigured }
func (s *Store) ListUsers(context.Context) ([]contract.AdminUser, error) {
	return nil, ErrNotConfigured
}
func (s *Store) GetAdminUser(context.Context, string) (*contract.AdminUser, error) {
	return nil, ErrNotConfigured
}
func (s *Store) ListAllOrgs(context.Context) ([]contract.Organization, error) {
	return nil, ErrNotConfigured
}
func (s *Store) ListAllDiagnoses(context.Context) ([]storage.DiagnosisRecord, error) {
	return nil, ErrNotConfigured
}
func (s *Store) ListAllMeasurements(context.Context) ([]contract.AdminMeasurement, error) {
	return nil, ErrNotConfigured
}
func (s *Store) ListAllChecks(context.Context) ([]contract.MonitorCheck, error) {
	return nil, ErrNotConfigured
}
func (s *Store) UpdateIncident(context.Context, contract.Incident) error { return ErrNotConfigured }
func (s *Store) AddIncidentNote(context.Context, contract.IncidentNote) error {
	return ErrNotConfigured
}
func (s *Store) ListIncidentNotes(context.Context, string) ([]contract.IncidentNote, error) {
	return nil, ErrNotConfigured
}
func (s *Store) SetIncidentOverride(context.Context, string, contract.IncidentOverride) error {
	return ErrNotConfigured
}
func (s *Store) GetIncidentOverride(context.Context, string) (*contract.IncidentOverride, error) {
	return nil, ErrNotConfigured
}
func (s *Store) AddAbuse(context.Context, contract.AbuseEvent) error { return ErrNotConfigured }
func (s *Store) ListAbuse(context.Context) ([]contract.AbuseEvent, error) {
	return nil, ErrNotConfigured
}
func (s *Store) AddAdminAudit(context.Context, contract.AdminAudit) error { return ErrNotConfigured }
func (s *Store) ListAdminAudit(context.Context) ([]contract.AdminAudit, error) {
	return nil, ErrNotConfigured
}
func (s *Store) ListFlags(context.Context) ([]contract.FeatureFlag, error) {
	return nil, ErrNotConfigured
}
func (s *Store) GetFlag(context.Context, string) (*contract.FeatureFlag, error) {
	return nil, ErrNotConfigured
}
func (s *Store) UpsertFlag(context.Context, contract.FeatureFlag) error { return ErrNotConfigured }
func (s *Store) ListRemoteConfig(context.Context) ([]contract.RemoteConfigEntry, error) {
	return nil, ErrNotConfigured
}
func (s *Store) GetRemoteConfig(context.Context, string) (*contract.RemoteConfigEntry, error) {
	return nil, ErrNotConfigured
}
func (s *Store) UpsertRemoteConfig(context.Context, contract.RemoteConfigEntry) error {
	return ErrNotConfigured
}
func (s *Store) AddRuleLabel(context.Context, contract.RuleLabel) error { return ErrNotConfigured }
func (s *Store) ListRuleLabels(context.Context) ([]contract.RuleLabel, error) {
	return nil, ErrNotConfigured
}
