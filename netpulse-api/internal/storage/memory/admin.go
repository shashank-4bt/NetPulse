package memory

import (
	"context"
	"sort"

	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/contract"
	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/storage"
)

func (s *Store) AddOperator(userID, email, role string, perms []string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(perms) == 0 {
		perms = append([]string{}, allAdminPerms()...)
	}
	s.operators[userID] = contract.Operator{
		UserID: userID, Email: email, Role: role, Permissions: append([]string{}, perms...),
	}
}

func allAdminPerms() []string {
	return []string{
		"users.read", "orgs.read", "services.read", "incidents.read", "incidents.operate",
		"measurements.read", "diagnostics.read", "rules.read", "abuse.read", "audit.read",
		"system.read", "flags.manage", "config.manage",
	}
}

func (s *Store) GetOperator(_ context.Context, userID string) (*contract.Operator, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	op, ok := s.operators[userID]
	if !ok {
		return nil, nil
	}
	copy := op
	copy.Permissions = append([]string{}, op.Permissions...)
	return &copy, nil
}

func (s *Store) UpsertOperator(_ context.Context, op contract.Operator) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	op.Permissions = append([]string{}, op.Permissions...)
	s.operators[op.UserID] = op
	return nil
}

func (s *Store) ListUsers(_ context.Context) ([]contract.AdminUser, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]contract.AdminUser, 0, len(s.users))
	for _, rec := range s.users {
		out = append(out, adminUserFrom(rec))
	}
	sort.Slice(out, func(i, j int) bool { return out[i].CreatedAt > out[j].CreatedAt })
	return out, nil
}

func (s *Store) GetAdminUser(_ context.Context, id string) (*contract.AdminUser, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	rec, ok := s.users[id]
	if !ok {
		return nil, nil
	}
	item := adminUserFrom(rec)
	return &item, nil
}

func adminUserFrom(rec storage.UserRecord) contract.AdminUser {
	return contract.AdminUser{
		ID:            rec.User.ID,
		Email:         rec.User.Email,
		DisplayName:   rec.User.DisplayName,
		EmailVerified: rec.User.EmailVerified,
		CreatedAt:     rec.User.CreatedAt,
		Summary:       "Account record only. Password hashes and session secrets are not included.",
	}
}

func (s *Store) ListAllOrgs(_ context.Context) ([]contract.Organization, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]contract.Organization, 0, len(s.orgs))
	for _, org := range s.orgs {
		copy := org
		copy.Role = ""
		out = append(out, copy)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].CreatedAt > out[j].CreatedAt })
	return out, nil
}

func (s *Store) ListAllDiagnoses(_ context.Context) ([]storage.DiagnosisRecord, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]storage.DiagnosisRecord, 0, len(s.diagnoses))
	for _, rec := range s.diagnoses {
		out = append(out, rec)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Diagnosis.Created > out[j].Diagnosis.Created })
	return out, nil
}

func (s *Store) ListAllMeasurements(_ context.Context) ([]contract.AdminMeasurement, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := []contract.AdminMeasurement{}
	for diagnosisID, items := range s.measurements {
		for _, item := range items {
			summary := ""
			if item.Summary != nil {
				summary = *item.Summary
			}
			out = append(out, contract.AdminMeasurement{
				DiagnosisID: diagnosisID,
				Key:         item.Key,
				Label:       item.Label,
				Measured:    item.Measured,
				Summary:     summary,
			})
		}
	}
	return out, nil
}

func (s *Store) ListAllChecks(_ context.Context) ([]contract.MonitorCheck, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := append([]contract.MonitorCheck{}, s.checks...)
	out = append(out, s.orgChecks...)
	return out, nil
}

func (s *Store) UpdateIncident(_ context.Context, item contract.Incident) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, existing := range s.incidents {
		if existing.ID == item.ID {
			s.incidents[i] = contract.NormalizeIncident(item)
			return nil
		}
	}
	s.incidents = append(s.incidents, contract.NormalizeIncident(item))
	return nil
}

func (s *Store) AddIncidentNote(_ context.Context, note contract.IncidentNote) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.incidentNotes = append(s.incidentNotes, note)
	return nil
}

func (s *Store) ListIncidentNotes(_ context.Context, incidentID string) ([]contract.IncidentNote, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := []contract.IncidentNote{}
	for _, note := range s.incidentNotes {
		if note.IncidentID == incidentID {
			out = append(out, note)
		}
	}
	return out, nil
}

func (s *Store) SetIncidentOverride(_ context.Context, incidentID string, ov contract.IncidentOverride) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.overrides[incidentID] = ov
	return nil
}

func (s *Store) GetIncidentOverride(_ context.Context, incidentID string) (*contract.IncidentOverride, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	ov, ok := s.overrides[incidentID]
	if !ok {
		return nil, nil
	}
	copy := ov
	return &copy, nil
}

func (s *Store) AddAbuse(_ context.Context, ev contract.AbuseEvent) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.abuse = append(s.abuse, ev)
	return nil
}

func (s *Store) ListAbuse(_ context.Context) ([]contract.AbuseEvent, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]contract.AbuseEvent, len(s.abuse))
	copy(out, s.abuse)
	sort.Slice(out, func(i, j int) bool { return out[i].At > out[j].At })
	return out, nil
}

func (s *Store) AddAdminAudit(_ context.Context, ev contract.AdminAudit) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.adminAudit = append(s.adminAudit, ev)
	return nil
}

func (s *Store) ListAdminAudit(_ context.Context) ([]contract.AdminAudit, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]contract.AdminAudit, len(s.adminAudit))
	copy(out, s.adminAudit)
	sort.Slice(out, func(i, j int) bool { return out[i].At > out[j].At })
	return out, nil
}

func (s *Store) ListFlags(_ context.Context) ([]contract.FeatureFlag, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]contract.FeatureFlag, 0, len(s.flags))
	for _, flag := range s.flags {
		out = append(out, copyFlag(flag))
	}
	sort.Slice(out, func(i, j int) bool { return out[i].UpdatedAt > out[j].UpdatedAt })
	return out, nil
}

func (s *Store) GetFlag(_ context.Context, id string) (*contract.FeatureFlag, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	flag, ok := s.flags[id]
	if !ok {
		return nil, nil
	}
	copy := copyFlag(flag)
	return &copy, nil
}

func (s *Store) UpsertFlag(_ context.Context, flag contract.FeatureFlag) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.flags[flag.ID] = copyFlag(flag)
	return nil
}

func copyFlag(flag contract.FeatureFlag) contract.FeatureFlag {
	flag.UserIDs = append([]string{}, flag.UserIDs...)
	flag.OrgIDs = append([]string{}, flag.OrgIDs...)
	if flag.UserIDs == nil {
		flag.UserIDs = []string{}
	}
	if flag.OrgIDs == nil {
		flag.OrgIDs = []string{}
	}
	return flag
}

func (s *Store) ListRemoteConfig(_ context.Context) ([]contract.RemoteConfigEntry, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]contract.RemoteConfigEntry, 0, len(s.remoteConfig))
	for _, entry := range s.remoteConfig {
		out = append(out, entry)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Key < out[j].Key })
	return out, nil
}

func (s *Store) GetRemoteConfig(_ context.Context, key string) (*contract.RemoteConfigEntry, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	entry, ok := s.remoteConfig[key]
	if !ok {
		return nil, nil
	}
	copy := entry
	return &copy, nil
}

func (s *Store) UpsertRemoteConfig(_ context.Context, entry contract.RemoteConfigEntry) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.remoteConfig[entry.Key] = entry
	return nil
}

func (s *Store) AddRuleLabel(_ context.Context, label contract.RuleLabel) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.ruleLabels = append(s.ruleLabels, label)
	return nil
}

func (s *Store) ListRuleLabels(_ context.Context) ([]contract.RuleLabel, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]contract.RuleLabel, len(s.ruleLabels))
	copy(out, s.ruleLabels)
	return out, nil
}
