package business

import (
	"context"
	"strings"

	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/auth"
	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/contract"
	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/id"
)

func (s *Service) ListMembers(ctx context.Context, actor *Actor, orgID string) ([]contract.Member, *contract.APIError, int) {
	if actor.Org.ID != orgID {
		return nil, apiErr("not_found", "organization not found"), 404
	}
	items, err := s.Store.ListMembers(ctx, orgID)
	if err != nil {
		return nil, apiErr("unavailable", "Member store is unavailable."), 503
	}
	for i := range items {
		items[i].Permissions = PermissionsFor(items[i].Role)
	}
	return items, nil, 200
}

func (s *Service) Invite(ctx context.Context, actor *Actor, orgID, email, role string) (*contract.Member, *contract.OrgInvite, *contract.APIError, int) {
	if errResp, status := s.require(actor, orgID, PermTeamManage); errResp != nil {
		return nil, nil, errResp, status
	}
	email = auth.NormalizeEmail(email)
	if msg := auth.ValidateEmail(email); msg != "" {
		return nil, nil, apiErr("validation_error", msg), 400
	}
	role = strings.ToLower(strings.TrimSpace(role))
	if role == "" {
		role = RoleViewer
	}
	if !KnownRole(role) {
		return nil, nil, apiErr("validation_error", "Unknown role."), 400
	}
	if role == RoleOwner && actor.Member != nil && actor.Member.Role != RoleOwner {
		return nil, nil, apiErr("forbidden", "Only an owner can grant owner."), 403
	}
	rec, err := s.Accounts.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, nil, apiErr("unavailable", "Account store is unavailable."), 503
	}
	if rec != nil {
		if existing, _ := s.Store.GetMemberByUser(ctx, orgID, rec.User.ID); existing != nil {
			return nil, nil, apiErr("validation_error", "That person is already a member."), 400
		}
		member := contract.Member{
			ID: id.New(), OrgID: orgID, UserID: rec.User.ID, Email: rec.User.Email,
			DisplayName: rec.User.DisplayName, Role: role, Permissions: PermissionsFor(role),
			CreatedAt: s.now().UTC().Format("2006-01-02T15:04:05Z07:00"),
		}
		if err := s.Store.CreateMember(ctx, member); err != nil {
			return nil, nil, apiErr("unavailable", "Member store is unavailable."), 503
		}
		s.audit(ctx, orgID, actor.UserID, "member.invited", "A member was added to the organization.")
		return &member, nil, nil, 201
	}
	invite := contract.OrgInvite{
		ID: id.New(), OrgID: orgID, Email: email, Role: role,
		CreatedAt: s.now().UTC().Format("2006-01-02T15:04:05Z07:00"),
		Summary:   "Invite recorded. A matching account is not confirmed in this response.",
	}
	if err := s.Store.CreateInvite(ctx, invite); err != nil {
		return nil, nil, apiErr("unavailable", "Invite store is unavailable."), 503
	}
	s.audit(ctx, orgID, actor.UserID, "invite.created", "An invite was stored.")
	return nil, &invite, nil, 201
}

func (s *Service) ChangeRole(ctx context.Context, actor *Actor, orgID, memberID, role string) (*contract.Member, *contract.APIError, int) {
	if errResp, status := s.require(actor, orgID, PermTeamManage); errResp != nil {
		return nil, errResp, status
	}
	role = strings.ToLower(strings.TrimSpace(role))
	if !KnownRole(role) {
		return nil, apiErr("validation_error", "Unknown role."), 400
	}
	member, err := s.Store.GetMember(ctx, orgID, memberID)
	if err != nil {
		return nil, apiErr("unavailable", "Member store is unavailable."), 503
	}
	if member == nil {
		return nil, apiErr("not_found", "member not found"), 404
	}
	if member.Role == RoleOwner && role != RoleOwner {
		owners := s.ownerCount(ctx, orgID)
		if owners <= 1 {
			return nil, apiErr("validation_error", "The last owner cannot be demoted."), 400
		}
	}
	if member.Role == RoleOwner && actor.Member != nil && actor.Member.Role != RoleOwner {
		return nil, apiErr("forbidden", "Only an owner can change an owner."), 403
	}
	member.Role = role
	member.Permissions = PermissionsFor(role)
	if err := s.Store.UpdateMember(ctx, *member); err != nil {
		return nil, apiErr("unavailable", "Member store is unavailable."), 503
	}
	s.audit(ctx, orgID, actor.UserID, "member.role", "A member role was changed.")
	return member, nil, 200
}

func (s *Service) RemoveMember(ctx context.Context, actor *Actor, orgID, memberID string) *contract.APIError {
	if errResp, _ := s.require(actor, orgID, PermTeamManage); errResp != nil {
		return errResp
	}
	member, err := s.Store.GetMember(ctx, orgID, memberID)
	if err != nil {
		return apiErr("unavailable", "Member store is unavailable.")
	}
	if member == nil {
		return apiErr("not_found", "member not found")
	}
	if member.Role == RoleOwner && s.ownerCount(ctx, orgID) <= 1 {
		return apiErr("validation_error", "The last owner cannot be removed.")
	}
	ok, err := s.Store.DeleteMember(ctx, orgID, memberID)
	if err != nil {
		return apiErr("unavailable", "Member store is unavailable.")
	}
	if !ok {
		return apiErr("not_found", "member not found")
	}
	s.audit(ctx, orgID, actor.UserID, "member.removed", "A member was removed.")
	return nil
}

func (s *Service) ListInvites(ctx context.Context, actor *Actor, orgID string) ([]contract.OrgInvite, *contract.APIError, int) {
	if errResp, status := s.require(actor, orgID, PermTeamManage); errResp != nil {
		return nil, errResp, status
	}
	items, err := s.Store.ListInvites(ctx, orgID)
	if err != nil {
		return nil, apiErr("unavailable", "Invite store is unavailable."), 503
	}
	if items == nil {
		items = []contract.OrgInvite{}
	}
	return items, nil, 200
}

func (s *Service) ListTeams(ctx context.Context, actor *Actor, orgID string) ([]contract.Team, *contract.APIError, int) {
	if actor.Org.ID != orgID {
		return nil, apiErr("not_found", "organization not found"), 404
	}
	items, err := s.Store.ListTeams(ctx, orgID)
	if err != nil {
		return nil, apiErr("unavailable", "Team store is unavailable."), 503
	}
	if items == nil {
		items = []contract.Team{}
	}
	return items, nil, 200
}

func (s *Service) CreateTeam(ctx context.Context, actor *Actor, orgID, name string, memberIDs []string) (*contract.Team, *contract.APIError, int) {
	if errResp, status := s.require(actor, orgID, PermTeamManage); errResp != nil {
		return nil, errResp, status
	}
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, apiErr("validation_error", "Team name is required."), 400
	}
	if memberIDs == nil {
		memberIDs = []string{}
	}
	team := contract.Team{
		ID: id.New(), OrgID: orgID, Name: name, MemberIDs: memberIDs,
		CreatedAt: s.now().UTC().Format("2006-01-02T15:04:05Z07:00"),
		Summary:   "Stored team. This is not a discovered group.",
	}
	if err := s.Store.CreateTeam(ctx, team); err != nil {
		return nil, apiErr("unavailable", "Team store is unavailable."), 503
	}
	s.audit(ctx, orgID, actor.UserID, "team.created", "A team was created.")
	return &team, nil, 201
}

func (s *Service) UpdateTeam(ctx context.Context, actor *Actor, orgID, teamID, name string, memberIDs []string) (*contract.Team, *contract.APIError, int) {
	if errResp, status := s.require(actor, orgID, PermTeamManage); errResp != nil {
		return nil, errResp, status
	}
	team, err := s.Store.GetTeam(ctx, teamID)
	if err != nil {
		return nil, apiErr("unavailable", "Team store is unavailable."), 503
	}
	if team == nil || team.OrgID != orgID {
		return nil, apiErr("not_found", "team not found"), 404
	}
	if strings.TrimSpace(name) != "" {
		team.Name = strings.TrimSpace(name)
	}
	if memberIDs != nil {
		team.MemberIDs = memberIDs
	}
	if err := s.Store.UpdateTeam(ctx, *team); err != nil {
		return nil, apiErr("unavailable", "Team store is unavailable."), 503
	}
	return team, nil, 200
}

func (s *Service) DeleteTeam(ctx context.Context, actor *Actor, orgID, teamID string) *contract.APIError {
	if errResp, _ := s.require(actor, orgID, PermTeamManage); errResp != nil {
		return errResp
	}
	ok, err := s.Store.DeleteTeam(ctx, orgID, teamID)
	if err != nil {
		return apiErr("unavailable", "Team store is unavailable.")
	}
	if !ok {
		return apiErr("not_found", "team not found")
	}
	return nil
}

func (s *Service) ownerCount(ctx context.Context, orgID string) int {
	items, err := s.Store.ListMembers(ctx, orgID)
	if err != nil {
		return 0
	}
	n := 0
	for _, item := range items {
		if item.Role == RoleOwner {
			n++
		}
	}
	return n
}
