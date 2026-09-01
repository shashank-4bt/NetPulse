package worker

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/contract"
	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/diagnostics"
	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/measurements"
	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/ssrf"
	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/storage"
	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/validation"
)

type ProbeRunner interface {
	Run(ctx context.Context, target validation.Target) ([]contract.Measurement, error)
}

type Worker struct {
	Store        storage.DiagnoseStore
	Measurements storage.MeasurementStore
	Queue        storage.Queue
	Runner       ProbeRunner
	Log          *slog.Logger
	Concurrency  int
}

func (w *Worker) Start(ctx context.Context) {
	n := w.Concurrency
	if n < 1 {
		n = 1
	}
	for i := 0; i < n; i++ {
		go w.loop(ctx)
	}
}

func (w *Worker) loop(ctx context.Context) {
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			job, ok := w.Queue.Dequeue(ctx)
			if !ok {
				continue
			}
			w.process(ctx, job)
		}
	}
}

func (w *Worker) ProcessOne(ctx context.Context, job storage.Job) {
	w.process(ctx, job)
}

func (w *Worker) process(ctx context.Context, job storage.Job) {
	rec, err := w.Store.GetDiagnosis(ctx, job.DiagnosisID)
	if err != nil || rec == nil {
		return
	}
	rec.Diagnosis.Status = "running"
	_ = w.Store.UpdateDiagnosis(ctx, *rec)

	target := validation.Target{
		Raw: job.Target.Raw, Hostname: job.Target.Hostname, Kind: job.Target.Kind, ServiceSlug: job.Target.ServiceSlug,
	}
	probeCtx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	results, runErr := w.Runner.Run(probeCtx, target)
	if runErr != nil && errors.Is(runErr, ssrf.ErrBlocked) {
		rec.Diagnosis.Status = "unavailable"
		rec.Diagnosis.Error = &contract.APIError{Code: "ssrf_blocked", Message: runErr.Error()}
		rec.Diagnosis.Report = nil
		_ = w.Store.UpdateDiagnosis(ctx, *rec)
		return
	}

	if runErr != nil && errors.Is(probeCtx.Err(), context.DeadlineExceeded) {
		rec.Diagnosis.Status = "partial"
	} else if runErr != nil {
		rec.Diagnosis.Status = "partial"
		rec.Diagnosis.Error = &contract.APIError{Code: "unavailable", Message: "one or more probes did not complete"}
	} else {
		rec.Diagnosis.Status = "insufficient_evidence"
	}

	if w.Measurements != nil && len(results) > 0 {
		_ = w.Measurements.Record(ctx, job.DiagnosisID, results)
	}

	report := diagnostics.Analyze(job.DiagnosisID, target, results, time.Now())
	if runErr != nil && errors.Is(probeCtx.Err(), context.DeadlineExceeded) {
		report.Outcome = "timeout"
	}
	rec.Diagnosis.Report = &report
	_ = w.Store.UpdateDiagnosis(ctx, *rec)
	if w.Log != nil {
		w.Log.Info("diagnosis processed", "id", job.DiagnosisID, "status", rec.Diagnosis.Status)
	}
}

func DefaultRunner() ProbeRunner {
	return measurements.NewRunner()
}
