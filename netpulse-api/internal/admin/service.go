package admin

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/config"
	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/contract"
	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/id"
	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/opsconfig"
	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/storage"
	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/worker"
)

type Service struct {
	Store       storage.AdminStore
	Diagnoses   storage.DiagnoseStore
	Queue       storage.Queue
	Cfg         config.Config
	StorageInfo map[string]string
	Worker      *worker.Worker
}

func (s *Service) Seed(ctx context.Context) {
	if s.Store == nil {
		return
	}
	existing, err := s.Store.ListRemoteConfig(ctx)
	if err != nil {
		return
	}
	have := map[string]struct{}{}
	for _, item := range existing {
		have[item.Key] = struct{}{}
	}
	for _, item := range opsconfig.Defaults(s.Cfg) {
		if _, ok := have[item.Key]; ok {
			continue
		}
		_ = s.Store.UpsertRemoteConfig(ctx, item)
	}
}

func (s *Service) stamp() string {
	return time.Now().UTC().Format(time.RFC3339)
}

func (s *Service) ConfigInt(ctx context.Context, key string, fallback int) int {
	if s.Store == nil {
		return fallback
	}
	entry, err := s.Store.GetRemoteConfig(ctx, key)
	if err != nil || entry == nil {
		return opsconfig.Int(opsconfig.Defaults(s.Cfg), key, fallback)
	}
	return opsconfig.Int([]contract.RemoteConfigEntry{*entry}, key, fallback)
}

func (s *Service) DiagnoseLimit(ctx context.Context) int {
	limit := s.ConfigInt(ctx, opsconfig.DiagnoseRateLimitPerMin, s.Cfg.RateLimitPerMin)
	if limit < 1 {
		return s.Cfg.RateLimitPerMin
	}
	return limit
}

func (s *Service) OperatorFor(ctx context.Context, userID, email string) (*contract.Operator, *contract.APIError, int) {
	if s.Store == nil {
		return nil, apiErr("unavailable", "Admin store is unavailable."), 503
	}
	op, err := s.Store.GetOperator(ctx, userID)
	if err != nil {
		return nil, apiErr("unavailable", "Operator store is unavailable."), 503
	}
	if op != nil {
		op.Permissions = NormalizePerms(op.Permissions)
		return op, nil, 200
	}
	email = strings.ToLower(strings.TrimSpace(email))
	for _, allowed := range s.Cfg.AdminEmails {
		if strings.ToLower(strings.TrimSpace(allowed)) == email {
			created := contract.Operator{
				UserID: userID, Email: email, Role: RoleOperator, Permissions: AllPermissions(),
			}
			if err := s.Store.UpsertOperator(ctx, created); err != nil {
				return nil, apiErr("unavailable", "Operator store is unavailable."), 503
			}
			return &created, nil, 200
		}
	}
	return nil, apiErr("not_found", "not found"), 404
}

func (s *Service) Audit(ctx context.Context, actorID, action, resource, result, summary string) {
	if s.Store == nil {
		return
	}
	_ = s.Store.AddAdminAudit(ctx, contract.AdminAudit{
		ID: id.New(), ActorID: actorID, Action: action, Resource: resource,
		At: s.stamp(), Result: result, Summary: summary,
	})
}

func (s *Service) RecordAbuse(ctx context.Context, kind, actor, ip, resource, result, summary string) {
	if s.Store == nil {
		return
	}
	_ = s.Store.AddAbuse(ctx, contract.AbuseEvent{
		ID: id.New(), Kind: kind, Actor: actor, IP: ip, Resource: resource,
		At: s.stamp(), Result: result, Summary: summary,
	})
}

func (s *Service) ListUsers(ctx context.Context) ([]contract.AdminUser, *contract.APIError, int) {
	items, err := s.Store.ListUsers(ctx)
	if err != nil {
		return nil, apiErr("unavailable", "User store is unavailable."), 503
	}
	if items == nil {
		items = []contract.AdminUser{}
	}
	return items, nil, 200
}

func (s *Service) GetUser(ctx context.Context, id string) (*contract.AdminUser, *contract.APIError, int) {
	item, err := s.Store.GetAdminUser(ctx, id)
	if err != nil {
		return nil, apiErr("unavailable", "User store is unavailable."), 503
	}
	if item == nil {
		return nil, apiErr("not_found", "not found"), 404
	}
	return item, nil, 200
}

func (s *Service) ListOrganizations(ctx context.Context) ([]contract.Organization, *contract.APIError, int) {
	items, err := s.Store.ListAllOrgs(ctx)
	if err != nil {
		return nil, apiErr("unavailable", "Organization store is unavailable."), 503
	}
	if items == nil {
		items = []contract.Organization{}
	}
	return items, nil, 200
}

func (s *Service) GetOrganization(ctx context.Context, orgID string) (*contract.Organization, *contract.APIError, int) {
	items, errResp, status := s.ListOrganizations(ctx)
	if errResp != nil {
		return nil, errResp, status
	}
	for i := range items {
		if items[i].ID == orgID {
			return &items[i], nil, 200
		}
	}
	return nil, apiErr("not_found", "not found"), 404
}

func (s *Service) ListServices(ctx context.Context) ([]contract.Service, *contract.APIError, int) {
	if s.Diagnoses == nil {
		return []contract.Service{}, nil, 200
	}
	items, err := s.Diagnoses.ListServices(ctx)
	if err != nil {
		return nil, apiErr("unavailable", "Service catalog is unavailable."), 503
	}
	if items == nil {
		items = []contract.Service{}
	}
	return items, nil, 200
}

func (s *Service) ListMeasurements(ctx context.Context) ([]contract.AdminMeasurement, *contract.APIError, int) {
	items, err := s.Store.ListAllMeasurements(ctx)
	if err != nil {
		return nil, apiErr("unavailable", "Measurement store is unavailable."), 503
	}
	if items == nil {
		items = []contract.AdminMeasurement{}
	}
	return items, nil, 200
}

func (s *Service) ListDiagnoses(ctx context.Context) ([]contract.AdminDiagnosis, *contract.APIError, int) {
	recs, err := s.Store.ListAllDiagnoses(ctx)
	if err != nil {
		return nil, apiErr("unavailable", "Diagnosis store is unavailable."), 503
	}
	out := make([]contract.AdminDiagnosis, 0, len(recs))
	for _, rec := range recs {
		target := rec.Target.Hostname
		if target == "" {
			target = rec.Target.Raw
		}
		out = append(out, contract.AdminDiagnosis{
			ID: rec.Diagnosis.ID, Target: target, Status: rec.Diagnosis.Status,
			UserID: rec.UserID, CreatedAt: rec.Diagnosis.Created,
			Summary: fmt.Sprintf("Stored diagnosis %s. Status is the worker outcome, not a user-path proof.", rec.Diagnosis.Status),
		})
	}
	return out, nil, 200
}

func (s *Service) ListAudit(ctx context.Context) ([]contract.AdminAudit, *contract.APIError, int) {
	items, err := s.Store.ListAdminAudit(ctx)
	if err != nil {
		return nil, apiErr("unavailable", "Audit store is unavailable."), 503
	}
	if items == nil {
		items = []contract.AdminAudit{}
	}
	return items, nil, 200
}

func (s *Service) ListAbuse(ctx context.Context) ([]contract.AbuseEvent, *contract.APIError, int) {
	items, err := s.Store.ListAbuse(ctx)
	if err != nil {
		return nil, apiErr("unavailable", "Abuse store is unavailable."), 503
	}
	if items == nil {
		items = []contract.AbuseEvent{}
	}
	return items, nil, 200
}

func apiErr(code, message string) *contract.APIError {
	return &contract.APIError{Code: code, Message: message}
}
