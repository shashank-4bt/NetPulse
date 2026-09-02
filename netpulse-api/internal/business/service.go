package business

import (
	"context"
	"strings"
	"time"

	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/auth"
	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/contract"
	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/diagnostics"
	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/id"
	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/storage"
)

type Actor struct {
	UserID      string
	Email       string
	DisplayName string
	Org         contract.Organization
	Member      *contract.Member
	Key         *contract.APIKey
	Perms       []string
	SessionOnly bool
}

type Service struct {
	Store       storage.BusinessStore
	Accounts    storage.AccountStore
	Diagnoses   *diagnostics.Service
	Now         func() time.Time
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

func (s *Service) LookupOrgKey(ctx context.Context, raw string) (*contract.APIKey, *contract.Organization, *contract.APIError, int) {
	raw = strings.TrimSpace(raw)
	if raw == "" || !strings.HasPrefix(raw, "npo_") {
		return nil, nil, apiErr("unauthorized", "API key is not valid."), 401
	}
	key, err := s.Store.GetOrgKeyByHash(ctx, auth.HashSecret(raw))
	if err != nil {
		return nil, nil, apiErr("unavailable", "Key store is unavailable."), 503
	}
	if key == nil || key.Revoked {
		return nil, nil, apiErr("unauthorized", "API key is not valid."), 401
	}
	org, err := s.Store.GetOrg(ctx, key.OrgID)
	if err != nil || org == nil {
		return nil, nil, apiErr("unauthorized", "API key is not valid."), 401
	}
	used := s.now().UTC().Format(time.RFC3339)
	key.LastUsedAt = &used
	_ = s.Store.UpdateOrgKey(ctx, *key)
	return key, org, nil, 200
}

func (s *Service) ConsumeInvites(ctx context.Context, userID, email string) {
	email = auth.NormalizeEmail(email)
	invites, err := s.Store.ListInvitesByEmail(ctx, email)
	if err != nil {
		return
	}
	for _, invite := range invites {
		if existing, _ := s.Store.GetMemberByUser(ctx, invite.OrgID, userID); existing != nil {
			_, _ = s.Store.DeleteInvite(ctx, invite.OrgID, invite.ID)
			continue
		}
		rec, err := s.Accounts.GetUserByID(ctx, userID)
		name := email
		if err == nil && rec != nil {
			name = rec.User.DisplayName
		}
		_ = s.Store.CreateMember(ctx, contract.Member{
			ID: id.New(), OrgID: invite.OrgID, UserID: userID, Email: email, DisplayName: name,
			Role: invite.Role, Permissions: PermissionsFor(invite.Role), CreatedAt: s.now().UTC().Format(time.RFC3339),
		})
		_, _ = s.Store.DeleteInvite(ctx, invite.OrgID, invite.ID)
		s.audit(ctx, invite.OrgID, userID, "invite.accepted", "Pending invite was applied after sign-in.")
	}
}

func (s *Service) MemberFor(ctx context.Context, orgID, userID string) (*contract.Member, *contract.Organization, *contract.APIError, int) {
	org, err := s.Store.GetOrg(ctx, orgID)
	if err != nil {
		return nil, nil, apiErr("unavailable", "Organization store is unavailable."), 503
	}
	if org == nil {
		return nil, nil, apiErr("not_found", "organization not found"), 404
	}
	member, err := s.Store.GetMemberByUser(ctx, orgID, userID)
	if err != nil {
		return nil, nil, apiErr("unavailable", "Member store is unavailable."), 503
	}
	if member == nil {
		return nil, nil, apiErr("not_found", "organization not found"), 404
	}
	member.Permissions = PermissionsFor(member.Role)
	org.Role = member.Role
	return member, org, nil, 200
}

func (s *Service) ListOrgs(ctx context.Context, userID, email string) ([]contract.Organization, *contract.APIError, int) {
	s.ConsumeInvites(ctx, userID, email)
	items, err := s.Store.ListOrgsForUser(ctx, userID)
	if err != nil {
		return nil, apiErr("unavailable", "Organization store is unavailable."), 503
	}
	if items == nil {
		items = []contract.Organization{}
	}
	return items, nil, 200
}

func (s *Service) CreateOrg(ctx context.Context, userID, email, displayName, name string) (*contract.Organization, *contract.APIError, int) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, apiErr("validation_error", "Organization name is required."), 400
	}
	now := s.now().UTC().Format(time.RFC3339)
	org := contract.Organization{
		ID: id.New(), Name: name, CreatedAt: now, UpdatedAt: now, Role: RoleOwner,
		Summary: "Stored organization. Membership is required to read any resource.",
	}
	if err := s.Store.CreateOrg(ctx, org); err != nil {
		return nil, apiErr("unavailable", "Organization store is unavailable."), 503
	}
	if err := s.Store.CreateMember(ctx, contract.Member{
		ID: id.New(), OrgID: org.ID, UserID: userID, Email: auth.NormalizeEmail(email),
		DisplayName: displayName, Role: RoleOwner, Permissions: PermissionsFor(RoleOwner), CreatedAt: now,
	}); err != nil {
		return nil, apiErr("unavailable", "Member store is unavailable."), 503
	}
	s.audit(ctx, org.ID, userID, "org.created", "Organization created.")
	return &org, nil, 201
}

func (s *Service) GetOrg(ctx context.Context, actor *Actor, orgID string) (*contract.Organization, *contract.APIError, int) {
	if actor.Org.ID != orgID {
		return nil, apiErr("not_found", "organization not found"), 404
	}
	org := actor.Org
	org.Role = ""
	if actor.Member != nil {
		org.Role = actor.Member.Role
	}
	return &org, nil, 200
}

func (s *Service) UpdateOrg(ctx context.Context, actor *Actor, orgID, name string) (*contract.Organization, *contract.APIError, int) {
	if errResp, status := s.require(actor, orgID, PermTeamManage); errResp != nil {
		return nil, errResp, status
	}
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, apiErr("validation_error", "Organization name is required."), 400
	}
	org := actor.Org
	org.Name = name
	org.UpdatedAt = s.now().UTC().Format(time.RFC3339)
	if err := s.Store.UpdateOrg(ctx, org); err != nil {
		return nil, apiErr("unavailable", "Organization store is unavailable."), 503
	}
	s.audit(ctx, orgID, actor.UserID, "org.updated", "Organization settings were updated.")
	return &org, nil, 200
}

func (s *Service) require(actor *Actor, orgID, perm string) (*contract.APIError, int) {
	if actor.Org.ID != orgID {
		return apiErr("not_found", "not found"), 404
	}
	if perm != "" && !HasPerm(actor.Perms, perm) {
		return apiErr("forbidden", "Missing permission."), 403
	}
	return nil, 200
}

func (s *Service) audit(ctx context.Context, orgID, actorID, kind, summary string) {
	_ = s.Store.AddAudit(ctx, contract.AuditEvent{
		ID: id.New(), OrgID: orgID, ActorID: actorID, Kind: kind,
		At: s.now().UTC().Format(time.RFC3339), Summary: summary,
	})
}
