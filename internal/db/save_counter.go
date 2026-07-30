package db

import (
	"context"
	"metrics/internal/retry"
)

func (p *DB) SaveCounter(ctx context.Context, name string, value int64) error {
	query := `
	INSERT INTO counter (name, delta) VALUES ($1, $2)
	ON CONFLICT (name) DO UPDATE SET delta = counter.delta + $2;
	`

	return retry.WithRetry(func() error {
		_, err := p.db.ExecContext(ctx, query, name, value)
		return err
	}, isRetriableDB)
}
