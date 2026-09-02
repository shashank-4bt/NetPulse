package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Addr              string
	CORSOrigin        string
	WorkerEmbedded    bool
	WorkerConcurrency int
	RateLimitPerMin   int
	DatabaseURL       string
	RedisURL          string
	ClickHouseURL     string
	EngineVersion     string
	AuthDevTokens     bool
	SessionTTLHours   int
}

func FromEnv() Config {
	return Config{
		Addr:              env("NETPULSE_API_ADDR", ":8080"),
		CORSOrigin:        env("NETPULSE_CORS_ORIGIN", "http://localhost:3000"),
		WorkerEmbedded:    envBool("NETPULSE_WORKER_EMBEDDED", true),
		WorkerConcurrency: envInt("NETPULSE_WORKER_CONCURRENCY", 2),
		RateLimitPerMin:   envInt("NETPULSE_RATE_LIMIT_PER_MIN", 10),
		DatabaseURL:       os.Getenv("NETPULSE_DATABASE_URL"),
		RedisURL:          os.Getenv("NETPULSE_REDIS_URL"),
		ClickHouseURL:     os.Getenv("NETPULSE_CLICKHOUSE_URL"),
		EngineVersion:     "0.11.0",
		AuthDevTokens:     envBool("NETPULSE_AUTH_DEV_TOKENS", false),
		SessionTTLHours:   envInt("NETPULSE_SESSION_TTL_HOURS", 168),
	}
}

func env(key, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}

func envBool(key string, fallback bool) bool {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}
	return parsed
}

func envInt(key string, fallback int) int {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}
