package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/accounts"
	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/admin"
	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/api"
	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/business"
	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/config"
	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/developer"
	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/diagnostics"
	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/logging"
	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/measurements"
	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/opsconfig"
	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/storage/clickhouse"
	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/storage/memory"
	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/storage/postgres"
	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/storage/redisx"
	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/worker"
)

func main() {
	cfg := config.FromEnv()
	log := logging.New("api")
	store := memory.New()
	storageInfo := map[string]string{
		"postgres":   store.Backend(),
		"clickhouse": store.Backend(),
		"redis":      store.Backend(),
	}

	if _, err := postgres.Open(cfg.DatabaseURL); err == nil {
		storageInfo["postgres"] = "postgres"
	} else if cfg.DatabaseURL != "" {
		log.Warn("postgres adapter not linked; using memory", "err", err)
	}
	if _, err := clickhouse.Open(cfg.ClickHouseURL); err == nil {
		storageInfo["clickhouse"] = "clickhouse"
	} else if cfg.ClickHouseURL != "" {
		log.Warn("clickhouse adapter not linked; using memory", "err", err)
	}
	if _, err := redisx.Open(cfg.RedisURL); err == nil {
		storageInfo["redis"] = "redis"
	} else if cfg.RedisURL != "" {
		log.Warn("redis adapter not linked; using memory", "err", err)
	}

	svc := &diagnostics.Service{Store: store, Queue: store}
	accountSvc := &accounts.Service{
		Accounts:   store,
		Diagnoses:  store,
		DevTokens:  cfg.AuthDevTokens,
		SessionTTL: time.Duration(cfg.SessionTTLHours) * time.Hour,
	}
	developerSvc := &developer.Service{
		Store:  store,
		Runner: measurements.NewRunner(),
	}
	businessSvc := &business.Service{
		Store:     store,
		Accounts:  store,
		Diagnoses: svc,
	}
	adminSvc := &admin.Service{
		Store:       store,
		Diagnoses:   store,
		Queue:       store,
		Cfg:         cfg,
		StorageInfo: storageInfo,
	}
	adminSvc.Seed(context.Background())
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if cfg.WorkerEmbedded {
		concurrency := adminSvc.ConfigInt(context.Background(), opsconfig.WorkerConcurrency, cfg.WorkerConcurrency)
		timeout := time.Duration(adminSvc.ConfigInt(context.Background(), opsconfig.DiagnoseTimeoutSeconds, 20)) * time.Second
		w := &worker.Worker{
			Store:        store,
			Measurements: store,
			Queue:        store,
			Runner:       worker.DefaultRunner(),
			Completions:  developerSvc,
			Log:          log,
			Concurrency:  concurrency,
			Timeout:      timeout,
		}
		w.Start(ctx)
		adminSvc.Worker = w
		log.Info("embedded worker started", "concurrency", concurrency)
	}

	server := &api.Server{
		Cfg:         cfg,
		Log:         log,
		Diagnostics: svc,
		Accounts:    accountSvc,
		Developer:   developerSvc,
		Business:    businessSvc,
		Admin:       adminSvc,
		Limiter:     store,
		StorageInfo: storageInfo,
	}
	httpServer := api.NewHTTPServer(cfg.Addr, server.Handler())
	go func() {
		<-ctx.Done()
		_ = httpServer.Shutdown(context.Background())
	}()
	log.Info("api listening", "addr", cfg.Addr)
	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Error("api stopped", slog.Any("err", err))
		os.Exit(1)
	}
}
