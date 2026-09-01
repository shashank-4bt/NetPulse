package storage

import (
	"context"
	"time"

	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/contract"
)

type DiagnosisRecord struct {
	Diagnosis contract.Diagnosis
	Target    contract.Target
}

type Job struct {
	DiagnosisID string
	Target      contract.Target
	QueuedAt    time.Time
}

type DiagnoseStore interface {
	CreateDiagnosis(ctx context.Context, rec DiagnosisRecord) error
	GetDiagnosis(ctx context.Context, id string) (*DiagnosisRecord, error)
	UpdateDiagnosis(ctx context.Context, rec DiagnosisRecord) error
	ListIncidents(ctx context.Context) ([]contract.Incident, error)
	ListServices(ctx context.Context) ([]contract.Service, error)
	GetService(ctx context.Context, slug string) (*contract.Service, error)
	Backend() string
}

type MeasurementStore interface {
	Record(ctx context.Context, diagnosisID string, measurements []contract.Measurement) error
	Backend() string
}

type Queue interface {
	Enqueue(ctx context.Context, job Job) error
	Dequeue(ctx context.Context) (Job, bool)
	Backend() string
}

type RateLimiter interface {
	Allow(key string, limitPerMin int) bool
	Backend() string
}

type Cache interface {
	Get(key string) (string, bool)
	Set(key, value string, ttl time.Duration)
	Backend() string
}
