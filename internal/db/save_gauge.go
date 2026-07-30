package db

import (
	"context"
	"metrics/internal/retry"
)

func (p *DB) SaveGauge(ctx context.Context, name string, value float64) error {

	query :=
		`INSERT INTO gauge (name,value) VALUES ($1, $2)
		ON CONFLICT (name) DO UPDATE SET value=$2;`

	return retry.WithRetry(func() error {
		_, err := p.db.ExecContext(ctx, query, name, value)
		return err
	}, isRetriableDB)
}
