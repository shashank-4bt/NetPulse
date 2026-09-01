package clickhouse

import (
	"context"
	"errors"

	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/contract"
)

var ErrNotConfigured = errors.New("clickhouse is not configured")

// Store records high-volume measurements. Domain code uses storage.MeasurementStore.
type Store struct{}

func Open(dsn string) (*Store, error) {
	if dsn == "" {
		return nil, ErrNotConfigured
	}
	return nil, errors.New("clickhouse DSN is set but no driver is linked; leave NETPULSE_CLICKHOUSE_URL empty")
}

func (s *Store) Record(context.Context, string, []contract.Measurement) error {
	return ErrNotConfigured
}
func (s *Store) Backend() string { return "clickhouse-unlinked" }
