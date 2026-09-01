package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/config"
	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/logging"
	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/storage/memory"
	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/worker"
)

func main() {
	cfg := config.FromEnv()
	log := logging.New("worker")
	store := memory.New()
	log.Info("standalone worker uses a process-local queue; run the API with NETPULSE_WORKER_EMBEDDED=true for local development, or configure Redis later")

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	w := &worker.Worker{
		Store:        store,
		Measurements: store,
		Queue:        store,
		Runner:       worker.DefaultRunner(),
		Log:          log,
		Concurrency:  cfg.WorkerConcurrency,
	}
	w.Start(ctx)
	<-ctx.Done()
}
