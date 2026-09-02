package business

import (
	"context"
	"strings"

	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/contract"
	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/id"
	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/validation"
)

func (s *Service) ListDevices(ctx context.Context, actor *Actor, orgID string) ([]contract.OrgDevice, *contract.APIError, int) {
	if errResp, status := s.requireRead(actor, orgID); errResp != nil {
		return nil, errResp, status
	}
	items, err := s.Store.ListOrgDevices(ctx, orgID)
	if err != nil {
		return nil, apiErr("unavailable", "Device store is unavailable."), 503
	}
	if items == nil {
		items = []contract.OrgDevice{}
	}
	return items, nil, 200
}

func (s *Service) CreateDevice(ctx context.Context, actor *Actor, orgID, name, label, region string) (*contract.OrgDevice, *contract.APIError, int) {
	if errResp, status := s.requireWrite(actor, orgID); errResp != nil {
		return nil, errResp, status
	}
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, apiErr("validation_error", "Device name is required."), 400
	}
	item := contract.OrgDevice{
		ID: id.New(), OrgID: orgID, Name: name, Label: strings.TrimSpace(label), Region: strings.TrimSpace(region),
		Summary: "Stored device record. This is not a discovered hardware inventory.", CreatedAt: s.stamp(),
	}
	if err := s.Store.CreateOrgDevice(ctx, item); err != nil {
		return nil, apiErr("unavailable", "Device store is unavailable."), 503
	}
	return &item, nil, 201
}

func (s *Service) DeleteDevice(ctx context.Context, actor *Actor, orgID, id string) *contract.APIError {
	if errResp, _ := s.requireDelete(actor, orgID); errResp != nil {
		return errResp
	}
	ok, err := s.Store.DeleteOrgDevice(ctx, orgID, id)
	if err != nil {
		return apiErr("unavailable", "Device store is unavailable.")
	}
	if !ok {
		return apiErr("not_found", "device not found")
	}
	return nil
}

func (s *Service) ListNetworks(ctx context.Context, actor *Actor, orgID string) ([]contract.OrgNetwork, *contract.APIError, int) {
	if errResp, status := s.requireRead(actor, orgID); errResp != nil {
		return nil, errResp, status
	}
	items, err := s.Store.ListOrgNetworks(ctx, orgID)
	if err != nil {
		return nil, apiErr("unavailable", "Network store is unavailable."), 503
	}
	if items == nil {
		items = []contract.OrgNetwork{}
	}
	return items, nil, 200
}

func (s *Service) CreateNetwork(ctx context.Context, actor *Actor, orgID, name, asn, region string) (*contract.OrgNetwork, *contract.APIError, int) {
	if errResp, status := s.requireWrite(actor, orgID); errResp != nil {
		return nil, errResp, status
	}
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, apiErr("validation_error", "Network name is required."), 400
	}
	item := contract.OrgNetwork{
		ID: id.New(), OrgID: orgID, Name: name, ASN: strings.TrimSpace(asn), Region: strings.TrimSpace(region),
		Summary: "Stored network label. This is not a live routing table.", CreatedAt: s.stamp(),
	}
	if err := s.Store.CreateOrgNetwork(ctx, item); err != nil {
		return nil, apiErr("unavailable", "Network store is unavailable."), 503
	}
	return &item, nil, 201
}

func (s *Service) DeleteNetwork(ctx context.Context, actor *Actor, orgID, id string) *contract.APIError {
	if errResp, _ := s.requireDelete(actor, orgID); errResp != nil {
		return errResp
	}
	ok, err := s.Store.DeleteOrgNetwork(ctx, orgID, id)
	if err != nil {
		return apiErr("unavailable", "Network store is unavailable.")
	}
	if !ok {
		return apiErr("not_found", "network not found")
	}
	return nil
}

func (s *Service) ListServices(ctx context.Context, actor *Actor, orgID string) ([]contract.OrgService, *contract.APIError, int) {
	if errResp, status := s.requireRead(actor, orgID); errResp != nil {
		return nil, errResp, status
	}
	items, err := s.Store.ListOrgServices(ctx, orgID)
	if err != nil {
		return nil, apiErr("unavailable", "Service store is unavailable."), 503
	}
	if items == nil {
		items = []contract.OrgService{}
	}
	return items, nil, 200
}

func (s *Service) CreateService(ctx context.Context, actor *Actor, orgID, name, slug, endpoint string) (*contract.OrgService, *contract.APIError, int) {
	if errResp, status := s.requireWrite(actor, orgID); errResp != nil {
		return nil, errResp, status
	}
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, apiErr("validation_error", "Service name is required."), 400
	}
	if endpoint != "" {
		if parsed := validation.ParseTarget(endpoint); parsed.Err != nil {
			code := parsed.Code
			if code == "" {
				code = "validation_error"
			}
			status := 400
			if code == "ssrf_blocked" {
				status = 403
			}
			return nil, apiErr(code, parsed.Err.Error()), status
		}
	}
	item := contract.OrgService{
		ID: id.New(), OrgID: orgID, Name: name, Slug: strings.TrimSpace(slug), Endpoint: strings.TrimSpace(endpoint),
		Summary: "Stored service record. This is not live health.", CreatedAt: s.stamp(),
	}
	if err := s.Store.CreateOrgService(ctx, item); err != nil {
		return nil, apiErr("unavailable", "Service store is unavailable."), 503
	}
	return &item, nil, 201
}

func (s *Service) DeleteService(ctx context.Context, actor *Actor, orgID, id string) *contract.APIError {
	if errResp, _ := s.requireDelete(actor, orgID); errResp != nil {
		return errResp
	}
	ok, err := s.Store.DeleteOrgService(ctx, orgID, id)
	if err != nil {
		return apiErr("unavailable", "Service store is unavailable.")
	}
	if !ok {
		return apiErr("not_found", "service not found")
	}
	return nil
}

func (s *Service) requireRead(actor *Actor, orgID string) (*contract.APIError, int) {
	if actor.Org.ID != orgID {
		return apiErr("not_found", "not found"), 404
	}
	if !CanReadOrg(actor.Perms) {
		return apiErr("forbidden", "Missing permission."), 403
	}
	return nil, 200
}

func (s *Service) requireWrite(actor *Actor, orgID string) (*contract.APIError, int) {
	if actor.Org.ID != orgID {
		return apiErr("not_found", "not found"), 404
	}
	if !HasPerm(actor.Perms, PermMonitorCreate) && !HasPerm(actor.Perms, PermTeamManage) {
		return apiErr("forbidden", "Missing permission."), 403
	}
	return nil, 200
}

func (s *Service) requireDelete(actor *Actor, orgID string) (*contract.APIError, int) {
	if actor.Org.ID != orgID {
		return apiErr("not_found", "not found"), 404
	}
	if !HasPerm(actor.Perms, PermMonitorDelete) && !HasPerm(actor.Perms, PermTeamManage) {
		return apiErr("forbidden", "Missing permission."), 403
	}
	return nil, 200
}

func (s *Service) stamp() string {
	return s.now().UTC().Format("2006-01-02T15:04:05Z07:00")
}
