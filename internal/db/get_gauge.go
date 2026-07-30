package db

import (
	"context"
	"metrics/internal/retry"
)

func (p *DB) GetGauge(ctx context.Context, name string) (float64, error) {
	query := `
	SELECT value FROM gauge
	WHERE name = $1
	`

	var result float64

	err := retry.WithRetry(func() error {
		row := p.db.QueryRowContext(ctx, query, name)
		return row.Scan(&result)
	}, isRetriableDB)

	if err != nil {
		return 0, err
	}

	return result, nil
}
