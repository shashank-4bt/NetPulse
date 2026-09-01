package redisx

import (
	"context"
	"errors"
	"time"

	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/storage"
)

var ErrNotConfigured = errors.New("redis is not configured")

// Store covers cache, rate limiting, and queue support. Domain logic uses
// storage.Queue, storage.RateLimiter, and storage.Cache.
type Store struct{}

func Open(url string) (*Store, error) {
	if url == "" {
		return nil, ErrNotConfigured
	}
	return nil, errors.New("redis URL is set but no client is linked; leave NETPULSE_REDIS_URL empty")
}

func (s *Store) Enqueue(context.Context, storage.Job) error { return ErrNotConfigured }
func (s *Store) Dequeue(context.Context) (storage.Job, bool) {
	return storage.Job{}, false
}
func (s *Store) Allow(string, int) bool            { return false }
func (s *Store) Get(string) (string, bool)         { return "", false }
func (s *Store) Set(string, string, time.Duration) {}
func (s *Store) Backend() string                   { return "redis-unlinked" }
