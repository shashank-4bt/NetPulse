package memory

import (
	"context"

	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/contract"
)

func (s *Store) CreateOrg(_ context.Context, org contract.Organization) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.orgs[org.ID] = org
	return nil
}

func (s *Store) GetOrg(_ context.Context, id string) (*contract.Organization, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	org, ok := s.orgs[id]
	if !ok {
		return nil, nil
	}
	copy := org
	return &copy, nil
}

func (s *Store) UpdateOrg(_ context.Context, org contract.Organization) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.orgs[org.ID] = org
	return nil
}

func (s *Store) ListOrgsForUser(_ context.Context, userID string) ([]contract.Organization, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := []contract.Organization{}
	for _, member := range s.members {
		if member.UserID != userID {
			continue
		}
		if org, ok := s.orgs[member.OrgID]; ok {
			org.Role = member.Role
			out = append(out, org)
		}
	}
	return out, nil
}

func (s *Store) CreateMember(_ context.Context, member contract.Member) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.members[member.ID] = member
	return nil
}

func (s *Store) GetMember(_ context.Context, orgID, memberID string) (*contract.Member, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	member, ok := s.members[memberID]
	if !ok || member.OrgID != orgID {
		return nil, nil
	}
	copy := member
	return &copy, nil
}

func (s *Store) GetMemberByUser(_ context.Context, orgID, userID string) (*contract.Member, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, member := range s.members {
		if member.OrgID == orgID && member.UserID == userID {
			copy := member
			return &copy, nil
		}
	}
	return nil, nil
}

func (s *Store) ListMembers(_ context.Context, orgID string) ([]contract.Member, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := []contract.Member{}
	for _, member := range s.members {
		if member.OrgID == orgID {
			out = append(out, member)
		}
	}
	return out, nil
}

func (s *Store) UpdateMember(_ context.Context, member contract.Member) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.members[member.ID] = member
	return nil
}

func (s *Store) DeleteMember(_ context.Context, orgID, memberID string) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	member, ok := s.members[memberID]
	if !ok || member.OrgID != orgID {
		return false, nil
	}
	delete(s.members, memberID)
	return true, nil
}

func (s *Store) CreateInvite(_ context.Context, invite contract.OrgInvite) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.invites[invite.ID] = invite
	return nil
}

func (s *Store) ListInvites(_ context.Context, orgID string) ([]contract.OrgInvite, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := []contract.OrgInvite{}
	for _, invite := range s.invites {
		if invite.OrgID == orgID {
			out = append(out, invite)
		}
	}
	return out, nil
}

func (s *Store) ListInvitesByEmail(_ context.Context, email string) ([]contract.OrgInvite, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := []contract.OrgInvite{}
	for _, invite := range s.invites {
		if invite.Email == email {
			out = append(out, invite)
		}
	}
	return out, nil
}

func (s *Store) DeleteInvite(_ context.Context, orgID, inviteID string) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	invite, ok := s.invites[inviteID]
	if !ok || invite.OrgID != orgID {
		return false, nil
	}
	delete(s.invites, inviteID)
	return true, nil
}

func (s *Store) CreateTeam(_ context.Context, team contract.Team) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.teams[team.ID] = team
	return nil
}

func (s *Store) GetTeam(_ context.Context, id string) (*contract.Team, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	team, ok := s.teams[id]
	if !ok {
		return nil, nil
	}
	copy := team
	return &copy, nil
}

func (s *Store) ListTeams(_ context.Context, orgID string) ([]contract.Team, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := []contract.Team{}
	for _, team := range s.teams {
		if team.OrgID == orgID {
			out = append(out, team)
		}
	}
	return out, nil
}

func (s *Store) UpdateTeam(_ context.Context, team contract.Team) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.teams[team.ID] = team
	return nil
}

func (s *Store) DeleteTeam(_ context.Context, orgID, id string) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	team, ok := s.teams[id]
	if !ok || team.OrgID != orgID {
		return false, nil
	}
	delete(s.teams, id)
	return true, nil
}

func (s *Store) CreateOrgDevice(_ context.Context, item contract.OrgDevice) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.orgDevices[item.ID] = item
	return nil
}

func (s *Store) GetOrgDevice(_ context.Context, id string) (*contract.OrgDevice, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	item, ok := s.orgDevices[id]
	if !ok {
		return nil, nil
	}
	copy := item
	return &copy, nil
}

func (s *Store) ListOrgDevices(_ context.Context, orgID string) ([]contract.OrgDevice, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := []contract.OrgDevice{}
	for _, item := range s.orgDevices {
		if item.OrgID == orgID {
			out = append(out, item)
		}
	}
	return out, nil
}

func (s *Store) UpdateOrgDevice(_ context.Context, item contract.OrgDevice) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.orgDevices[item.ID] = item
	return nil
}

func (s *Store) DeleteOrgDevice(_ context.Context, orgID, id string) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	item, ok := s.orgDevices[id]
	if !ok || item.OrgID != orgID {
		return false, nil
	}
	delete(s.orgDevices, id)
	return true, nil
}

func (s *Store) CreateOrgNetwork(_ context.Context, item contract.OrgNetwork) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.orgNetworks[item.ID] = item
	return nil
}

func (s *Store) GetOrgNetwork(_ context.Context, id string) (*contract.OrgNetwork, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	item, ok := s.orgNetworks[id]
	if !ok {
		return nil, nil
	}
	copy := item
	return &copy, nil
}

func (s *Store) ListOrgNetworks(_ context.Context, orgID string) ([]contract.OrgNetwork, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := []contract.OrgNetwork{}
	for _, item := range s.orgNetworks {
		if item.OrgID == orgID {
			out = append(out, item)
		}
	}
	return out, nil
}

func (s *Store) UpdateOrgNetwork(_ context.Context, item contract.OrgNetwork) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.orgNetworks[item.ID] = item
	return nil
}

func (s *Store) DeleteOrgNetwork(_ context.Context, orgID, id string) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	item, ok := s.orgNetworks[id]
	if !ok || item.OrgID != orgID {
		return false, nil
	}
	delete(s.orgNetworks, id)
	return true, nil
}

func (s *Store) CreateOrgService(_ context.Context, item contract.OrgService) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.orgServices[item.ID] = item
	return nil
}

func (s *Store) GetOrgService(_ context.Context, id string) (*contract.OrgService, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	item, ok := s.orgServices[id]
	if !ok {
		return nil, nil
	}
	copy := item
	return &copy, nil
}

func (s *Store) ListOrgServices(_ context.Context, orgID string) ([]contract.OrgService, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := []contract.OrgService{}
	for _, item := range s.orgServices {
		if item.OrgID == orgID {
			out = append(out, item)
		}
	}
	return out, nil
}

func (s *Store) UpdateOrgService(_ context.Context, item contract.OrgService) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.orgServices[item.ID] = item
	return nil
}

func (s *Store) DeleteOrgService(_ context.Context, orgID, id string) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	item, ok := s.orgServices[id]
	if !ok || item.OrgID != orgID {
		return false, nil
	}
	delete(s.orgServices, id)
	return true, nil
}

func (s *Store) CreateOrgMonitor(_ context.Context, item contract.Monitor) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.orgMonitors[item.ID] = item
	return nil
}

func (s *Store) GetOrgMonitor(_ context.Context, id string) (*contract.Monitor, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	item, ok := s.orgMonitors[id]
	if !ok {
		return nil, nil
	}
	copy := item
	return &copy, nil
}

func (s *Store) ListOrgMonitors(_ context.Context, orgID string) ([]contract.Monitor, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := []contract.Monitor{}
	for _, item := range s.orgMonitors {
		if item.OrgID == orgID {
			out = append(out, item)
		}
	}
	return out, nil
}

func (s *Store) UpdateOrgMonitor(_ context.Context, item contract.Monitor) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.orgMonitors[item.ID] = item
	return nil
}

func (s *Store) DeleteOrgMonitor(_ context.Context, orgID, id string) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	item, ok := s.orgMonitors[id]
	if !ok || item.OrgID != orgID {
		return false, nil
	}
	delete(s.orgMonitors, id)
	return true, nil
}

func (s *Store) AddOrgCheck(_ context.Context, check contract.MonitorCheck) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.orgChecks = append(s.orgChecks, check)
	return nil
}

func (s *Store) ListOrgChecks(_ context.Context, orgID, monitorID string) ([]contract.MonitorCheck, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := []contract.MonitorCheck{}
	for _, check := range s.orgChecks {
		if check.OrgID != orgID {
			continue
		}
		if monitorID != "" && check.MonitorID != monitorID {
			continue
		}
		out = append(out, check)
	}
	return out, nil
}

func (s *Store) CreateOrgIncident(_ context.Context, item contract.OrgIncident) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.orgIncidents[item.ID] = item
	return nil
}

func (s *Store) GetOrgIncident(_ context.Context, id string) (*contract.OrgIncident, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	item, ok := s.orgIncidents[id]
	if !ok {
		return nil, nil
	}
	copy := item
	return &copy, nil
}

func (s *Store) ListOrgIncidents(_ context.Context, orgID string) ([]contract.OrgIncident, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := []contract.OrgIncident{}
	for _, item := range s.orgIncidents {
		if item.OrgID == orgID {
			out = append(out, item)
		}
	}
	return out, nil
}

func (s *Store) UpdateOrgIncident(_ context.Context, item contract.OrgIncident) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.orgIncidents[item.ID] = item
	return nil
}

func (s *Store) CreateOrgDiagnosis(_ context.Context, item contract.OrgDiagnosis) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.orgDiagnoses[item.ID] = item
	return nil
}

func (s *Store) ListOrgDiagnoses(_ context.Context, orgID string) ([]contract.OrgDiagnosis, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := []contract.OrgDiagnosis{}
	for _, item := range s.orgDiagnoses {
		if item.OrgID == orgID {
			out = append(out, item)
		}
	}
	return out, nil
}

func (s *Store) GetOrgDiagnosis(_ context.Context, id string) (*contract.OrgDiagnosis, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	item, ok := s.orgDiagnoses[id]
	if !ok {
		return nil, nil
	}
	copy := item
	return &copy, nil
}

func (s *Store) CreateOrgReport(_ context.Context, item contract.OrgReport) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.orgReports[item.ID] = item
	return nil
}

func (s *Store) GetOrgReport(_ context.Context, id string) (*contract.OrgReport, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	item, ok := s.orgReports[id]
	if !ok {
		return nil, nil
	}
	copy := item
	return &copy, nil
}

func (s *Store) ListOrgReports(_ context.Context, orgID string) ([]contract.OrgReport, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := []contract.OrgReport{}
	for _, item := range s.orgReports {
		if item.OrgID == orgID {
			out = append(out, item)
		}
	}
	return out, nil
}

func (s *Store) CreateOrgKey(_ context.Context, key contract.APIKey) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.orgKeys[key.ID] = key
	s.orgKeyHash[key.Hash] = key.ID
	return nil
}

func (s *Store) GetOrgKey(_ context.Context, id string) (*contract.APIKey, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	key, ok := s.orgKeys[id]
	if !ok {
		return nil, nil
	}
	copy := key
	return &copy, nil
}

func (s *Store) GetOrgKeyByHash(_ context.Context, hash string) (*contract.APIKey, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	id, ok := s.orgKeyHash[hash]
	if !ok {
		return nil, nil
	}
	key := s.orgKeys[id]
	copy := key
	return &copy, nil
}

func (s *Store) ListOrgKeys(_ context.Context, orgID string) ([]contract.APIKey, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := []contract.APIKey{}
	for _, key := range s.orgKeys {
		if key.OrgID == orgID {
			out = append(out, key)
		}
	}
	return out, nil
}

func (s *Store) UpdateOrgKey(_ context.Context, key contract.APIKey) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if existing, ok := s.orgKeys[key.ID]; ok && existing.Hash != key.Hash {
		delete(s.orgKeyHash, existing.Hash)
	}
	s.orgKeys[key.ID] = key
	if key.Hash != "" {
		s.orgKeyHash[key.Hash] = key.ID
	}
	return nil
}

func (s *Store) AddAudit(_ context.Context, event contract.AuditEvent) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.audit = append(s.audit, event)
	return nil
}

func (s *Store) ListAudit(_ context.Context, orgID string) ([]contract.AuditEvent, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := []contract.AuditEvent{}
	for _, event := range s.audit {
		if event.OrgID == orgID {
			out = append(out, event)
		}
	}
	return out, nil
}
