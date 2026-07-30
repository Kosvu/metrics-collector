package db

import (
	"context"
	"metrics/internal/retry"
)

func (p *DB) Ping(ctx context.Context) error {
	retry.WithRetry(func() error {
		err := p.db.PingContext(ctx)
		return err
	}, isRetriableDB)

	return nil
}
