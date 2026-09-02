package accounts

import (
	"context"
	"sort"
	"strings"
	"time"

	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/auth"
	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/contract"
	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/storage"
)

func (s *Service) Dashboard(ctx context.Context, sessionToken string) (*contract.Dashboard, *contract.APIError, int) {
	_, user, errResp, status := s.Require(ctx, sessionToken)
	if errResp != nil {
		return nil, errResp, status
	}
	summaries, reports, errResp, status := s.ownedLists(ctx, user.ID, "", "", "", "", "")
	if errResp != nil {
		return nil, errResp, status
	}
	if len(summaries) > 8 {
		summaries = summaries[:8]
	}
	if len(reports) > 8 {
		reports = reports[:8]
	}
	saved, err := s.Accounts.ListSavedServices(ctx, user.ID)
	if err != nil {
		return nil, apiErr("unavailable", "Account store is unavailable."), 503
	}
	if saved == nil {
		saved = []contract.SavedService{}
	}
	incidents := []contract.Incident{}
	if s.Diagnoses != nil {
		listed, listErr := s.Diagnoses.ListIncidents(ctx)
		if listErr == nil && listed != nil {
			for _, item := range listed {
				normalized := contract.NormalizeIncident(item)
				if strings.EqualFold(normalized.Status, "resolved") {
					continue
				}
				incidents = append(incidents, normalized)
			}
		}
	}
	if len(incidents) > 8 {
		incidents = incidents[:8]
	}
	alerts, errResp, status := s.Alerts(ctx, sessionToken)
	if errResp != nil {
		return nil, errResp, status
	}
	return &contract.Dashboard{
		InternetHealth: "Internet Health on this page is an observatory snapshot, not your device, Wi-Fi, or ISP path. No live score is published here.",
		NetworkInfo:    "NetPulse does not collect ISP, Wi-Fi, or browser network details. Network information stays unavailable until you submit a diagnosis target.",
		Diagnoses:      summaries,
		SavedServices:  saved,
		Incidents:      incidents,
		Reports:        reports,
		Alerts:         *alerts,
	}, nil, 200
}

func (s *Service) ListDiagnoses(ctx context.Context, sessionToken, q, statusFilter, target, from, to string) ([]contract.DiagnosisSummary, *contract.APIError, int) {
	_, user, errResp, status := s.Require(ctx, sessionToken)
	if errResp != nil {
		return nil, errResp, status
	}
	summaries, _, errResp, status := s.ownedLists(ctx, user.ID, q, statusFilter, target, from, to)
	if errResp != nil {
		return nil, errResp, status
	}
	return summaries, nil, 200
}

func (s *Service) ListReports(ctx context.Context, sessionToken string) ([]contract.UserReport, *contract.APIError, int) {
	_, user, errResp, status := s.Require(ctx, sessionToken)
	if errResp != nil {
		return nil, errResp, status
	}
	_, reports, errResp, status := s.ownedLists(ctx, user.ID, "", "", "", "", "")
	if errResp != nil {
		return nil, errResp, status
	}
	return reports, nil, 200
}

func (s *Service) ShareReport(ctx context.Context, sessionToken, diagnosisID string) (*contract.ShareLink, *contract.APIError, int) {
	_, user, errResp, status := s.Require(ctx, sessionToken)
	if errResp != nil {
		return nil, errResp, status
	}
	rec, errResp, status := s.ownedDiagnosis(ctx, user.ID, diagnosisID)
	if errResp != nil {
		return nil, errResp, status
	}
	raw := auth.NewSecret()
	if err := s.Accounts.CreateShare(ctx, storage.ShareRecord{
		Hash:        auth.HashSecret(raw),
		DiagnosisID: rec.Diagnosis.ID,
		UserID:      user.ID,
		CreatedAt:   s.now(),
	}); err != nil {
		return nil, apiErr("unavailable", "Account store is unavailable."), 503
	}
	return &contract.ShareLink{
		Token:   raw,
		Path:    "/share/" + raw,
		Summary: "Anyone with this token can view this report. The token is shown once.",
	}, nil, 201
}

func (s *Service) DeleteReport(ctx context.Context, sessionToken, diagnosisID string) (*contract.APIError, int) {
	_, user, errResp, status := s.Require(ctx, sessionToken)
	if errResp != nil {
		return errResp, status
	}
	if _, errResp, status = s.ownedDiagnosis(ctx, user.ID, diagnosisID); errResp != nil {
		return errResp, status
	}
	_ = s.Accounts.DeleteSharesForDiagnosis(ctx, diagnosisID)
	if err := s.Diagnoses.DeleteDiagnosis(ctx, diagnosisID); err != nil {
		return apiErr("unavailable", "Diagnose store is unavailable."), 503
	}
	return nil, 200
}

func (s *Service) ReadShare(ctx context.Context, raw string) (*contract.Diagnosis, *contract.APIError, int) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, apiErr("not_found", "Share link was not found."), 404
	}
	rec, err := s.Accounts.GetShare(ctx, auth.HashSecret(raw))
	if err != nil {
		return nil, apiErr("unavailable", "Account store is unavailable."), 503
	}
	if rec == nil {
		return nil, apiErr("not_found", "Share link was not found."), 404
	}
	diag, err := s.Diagnoses.GetDiagnosis(ctx, rec.DiagnosisID)
	if err != nil {
		return nil, apiErr("unavailable", "Diagnose store is unavailable."), 503
	}
	if diag == nil {
		return nil, apiErr("not_found", "Share link was not found."), 404
	}
	copy := diag.Diagnosis
	return &copy, nil, 200
}

func (s *Service) MayReadDiagnosis(ctx context.Context, rec *storage.DiagnosisRecord, viewerUserID, shareToken string) bool {
	if rec == nil {
		return false
	}
	if rec.UserID == "" {
		return true
	}
	if viewerUserID != "" && viewerUserID == rec.UserID {
		return true
	}
	shareToken = strings.TrimSpace(shareToken)
	if shareToken == "" {
		return false
	}
	share, err := s.Accounts.GetShare(ctx, auth.HashSecret(shareToken))
	if err != nil || share == nil {
		return false
	}
	return share.DiagnosisID == rec.Diagnosis.ID
}

func (s *Service) SaveService(ctx context.Context, sessionToken, slug string) ([]contract.SavedService, *contract.APIError, int) {
	_, user, errResp, status := s.Require(ctx, sessionToken)
	if errResp != nil {
		return nil, errResp, status
	}
	slug = strings.TrimSpace(strings.ToLower(slug))
	if slug == "" {
		return nil, apiErr("validation_error", "A service slug is required."), 400
	}
	if s.Diagnoses != nil {
		item, err := s.Diagnoses.GetService(ctx, slug)
		if err != nil {
			return nil, apiErr("unavailable", "Service catalog is unavailable."), 503
		}
		if item == nil {
			return nil, apiErr("not_found", "Service was not found."), 404
		}
	}
	if err := s.Accounts.SaveService(ctx, user.ID, contract.SavedService{
		Slug:      slug,
		CreatedAt: s.now().UTC().Format(time.RFC3339),
	}); err != nil {
		return nil, apiErr("unavailable", "Account store is unavailable."), 503
	}
	return s.listSaved(ctx, user.ID)
}

func (s *Service) ListSaved(ctx context.Context, sessionToken string) ([]contract.SavedService, *contract.APIError, int) {
	_, user, errResp, status := s.Require(ctx, sessionToken)
	if errResp != nil {
		return nil, errResp, status
	}
	return s.listSaved(ctx, user.ID)
}

func (s *Service) DeleteSaved(ctx context.Context, sessionToken, slug string) ([]contract.SavedService, *contract.APIError, int) {
	_, user, errResp, status := s.Require(ctx, sessionToken)
	if errResp != nil {
		return nil, errResp, status
	}
	if err := s.Accounts.DeleteSavedService(ctx, user.ID, strings.TrimSpace(slug)); err != nil {
		return nil, apiErr("unavailable", "Account store is unavailable."), 503
	}
	return s.listSaved(ctx, user.ID)
}

func (s *Service) Devices(ctx context.Context, sessionToken string) ([]contract.Device, *contract.APIError, int) {
	sessions, errResp, status := s.ListSessions(ctx, sessionToken)
	if errResp != nil {
		return nil, errResp, status
	}
	out := make([]contract.Device, 0, len(sessions))
	for _, session := range sessions {
		out = append(out, contract.Device{
			ID:        session.ID,
			Label:     session.Label,
			UserAgent: session.UserAgent,
			IP:        session.IP,
			LastSeen:  session.LastSeenAt,
			Current:   session.Current,
			Kind:      "session",
		})
	}
	return out, nil, 200
}

func (s *Service) listSaved(ctx context.Context, userID string) ([]contract.SavedService, *contract.APIError, int) {
	items, err := s.Accounts.ListSavedServices(ctx, userID)
	if err != nil {
		return nil, apiErr("unavailable", "Account store is unavailable."), 503
	}
	if items == nil {
		items = []contract.SavedService{}
	}
	return items, nil, 200
}

func (s *Service) ownedDiagnosis(ctx context.Context, userID, diagnosisID string) (*storage.DiagnosisRecord, *contract.APIError, int) {
	if s.Diagnoses == nil {
		return nil, apiErr("unavailable", "Diagnose store is unavailable."), 503
	}
	rec, err := s.Diagnoses.GetDiagnosis(ctx, diagnosisID)
	if err != nil {
		return nil, apiErr("unavailable", "Diagnose store is unavailable."), 503
	}
	if rec == nil || rec.UserID == "" || rec.UserID != userID {
		return nil, apiErr("not_found", "Report was not found."), 404
	}
	return rec, nil, 200
}

func (s *Service) ownedLists(ctx context.Context, userID, q, statusFilter, target, from, to string) ([]contract.DiagnosisSummary, []contract.UserReport, *contract.APIError, int) {
	if s.Diagnoses == nil {
		return nil, nil, apiErr("unavailable", "Diagnose store is unavailable."), 503
	}
	items, err := s.Diagnoses.ListDiagnosesByUser(ctx, userID)
	if err != nil {
		return nil, nil, apiErr("unavailable", "Diagnose store is unavailable."), 503
	}
	shares, err := s.Accounts.ListSharesByUser(ctx, userID)
	if err != nil {
		return nil, nil, apiErr("unavailable", "Account store is unavailable."), 503
	}
	shared := map[string]bool{}
	for _, share := range shares {
		shared[share.DiagnosisID] = true
	}
	q = strings.ToLower(strings.TrimSpace(q))
	statusFilter = strings.ToLower(strings.TrimSpace(statusFilter))
	target = strings.ToLower(strings.TrimSpace(target))
	fromTime, fromOK := parseBoundary(from, false)
	toTime, toOK := parseBoundary(to, true)
	summaries := []contract.DiagnosisSummary{}
	reports := []contract.UserReport{}
	for _, rec := range items {
		created, _ := time.Parse(time.RFC3339, rec.Diagnosis.Created)
		if fromOK && created.Before(fromTime) {
			continue
		}
		if toOK && created.After(toTime) {
			continue
		}
		rawTarget := strings.ToLower(rec.Target.Raw + " " + rec.Target.Hostname)
		if target != "" && !strings.Contains(rawTarget, target) {
			continue
		}
		if statusFilter != "" && !strings.EqualFold(rec.Diagnosis.Status, statusFilter) {
			continue
		}
		if q != "" {
			hay := rawTarget + " " + strings.ToLower(rec.Diagnosis.Status) + " " + strings.ToLower(rec.Diagnosis.ID)
			if rec.Diagnosis.Report != nil {
				hay += " " + strings.ToLower(rec.Diagnosis.Report.Outcome)
			}
			if !strings.Contains(hay, q) {
				continue
			}
		}
		var outcome *string
		if rec.Diagnosis.Report != nil && rec.Diagnosis.Report.Outcome != "" {
			value := rec.Diagnosis.Report.Outcome
			outcome = &value
		}
		summaries = append(summaries, contract.DiagnosisSummary{
			ID:        rec.Diagnosis.ID,
			Target:    rec.Target.Raw,
			Status:    rec.Diagnosis.Status,
			Outcome:   outcome,
			CreatedAt: rec.Diagnosis.Created,
		})
		reports = append(reports, contract.UserReport{
			ID:        rec.Diagnosis.ID,
			Target:    rec.Target.Raw,
			Status:    rec.Diagnosis.Status,
			Shared:    shared[rec.Diagnosis.ID],
			CreatedAt: rec.Diagnosis.Created,
			Outcome:   outcome,
		})
	}
	sort.SliceStable(summaries, func(i, j int) bool { return summaries[i].CreatedAt > summaries[j].CreatedAt })
	sort.SliceStable(reports, func(i, j int) bool { return reports[i].CreatedAt > reports[j].CreatedAt })
	return summaries, reports, nil, 200
}

func parseBoundary(raw string, endOfDay bool) (time.Time, bool) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return time.Time{}, false
	}
	if parsed, err := time.Parse(time.RFC3339, raw); err == nil {
		return parsed, true
	}
	if parsed, err := time.Parse("2006-01-02", raw); err == nil {
		if endOfDay {
			return parsed.Add(24*time.Hour - time.Nanosecond), true
		}
		return parsed, true
	}
	return time.Time{}, false
}
