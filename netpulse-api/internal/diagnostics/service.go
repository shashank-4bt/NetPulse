package diagnostics

import (
	"context"
	"fmt"
	"time"

	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/contract"
	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/id"
	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/storage"
	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/validation"
)

type Service struct {
	Store storage.DiagnoseStore
	Queue storage.Queue
	Now   func() time.Time
}

func (s *Service) Create(ctx context.Context, raw string) (*contract.Diagnosis, *contract.APIError, int) {
	parsed := validation.ParseTarget(raw)
	if parsed.Err != nil {
		code := parsed.Code
		if code == "" {
			code = "validation_error"
		}
		status := 400
		if code == "ssrf_blocked" {
			status = 403
		}
		return nil, &contract.APIError{Code: code, Message: parsed.Err.Error()}, status
	}

	now := s.now()
	diag := contract.Diagnosis{
		ID:      id.New(),
		Status:  "queued",
		Created: now.UTC().Format(time.RFC3339),
	}
	rec := storage.DiagnosisRecord{Diagnosis: diag, Target: contract.Target{
		Raw: parsed.Target.Raw, Hostname: parsed.Target.Hostname, Kind: parsed.Target.Kind, ServiceSlug: parsed.Target.ServiceSlug,
	}}
	if err := s.Store.CreateDiagnosis(ctx, rec); err != nil {
		return nil, &contract.APIError{Code: "unavailable", Message: "diagnose store is unavailable"}, 503
	}
	if err := s.Queue.Enqueue(ctx, storage.Job{DiagnosisID: diag.ID, Target: rec.Target, QueuedAt: now}); err != nil {
		return nil, &contract.APIError{Code: "unavailable", Message: "measurement queue is unavailable"}, 503
	}
	return &diag, nil, 201
}

func (s *Service) Get(ctx context.Context, diagnosisID string) (*contract.Diagnosis, *contract.APIError, int) {
	rec, err := s.Store.GetDiagnosis(ctx, diagnosisID)
	if err != nil {
		return nil, &contract.APIError{Code: "unavailable", Message: "diagnose store is unavailable"}, 503
	}
	if rec == nil {
		return nil, &contract.APIError{Code: "not_found", Message: "diagnosis not found"}, 404
	}
	copy := rec.Diagnosis
	return &copy, nil, 200
}

func (s *Service) ListServices(ctx context.Context) ([]contract.Service, *contract.APIError, int) {
	items, err := s.Store.ListServices(ctx)
	if err != nil {
		return nil, &contract.APIError{Code: "unavailable", Message: "service catalog is unavailable"}, 503
	}
	return items, nil, 200
}

func (s *Service) GetService(ctx context.Context, slug string) (*contract.Service, *contract.APIError, int) {
	item, err := s.Store.GetService(ctx, slug)
	if err != nil {
		return nil, &contract.APIError{Code: "unavailable", Message: "service catalog is unavailable"}, 503
	}
	if item == nil {
		return nil, &contract.APIError{Code: "not_found", Message: "service not found"}, 404
	}
	return item, nil, 200
}

func (s *Service) ListIncidents(ctx context.Context) ([]contract.Incident, *contract.APIError, int) {
	items, err := s.Store.ListIncidents(ctx)
	if err != nil {
		return nil, &contract.APIError{Code: "unavailable", Message: "incident store is unavailable"}, 503
	}
	if items == nil {
		items = []contract.Incident{}
	}
	return items, nil, 200
}

func (s *Service) now() time.Time {
	if s.Now != nil {
		return s.Now()
	}
	return time.Now()
}

func UnknownError() *contract.APIError {
	return &contract.APIError{Code: "internal", Message: "an unexpected error occurred"}
}

func FormatStoreBackends(diagnose, measurements, queue string) map[string]string {
	return map[string]string{
		"postgres":   diagnose,
		"clickhouse": measurements,
		"redis":      queue,
	}
}

func EnsureQueued(rec *storage.DiagnosisRecord) error {
	if rec == nil {
		return fmt.Errorf("missing diagnosis")
	}
	return nil
}
