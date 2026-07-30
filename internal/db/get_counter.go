package db

import (
	"context"
	"metrics/internal/retry"
)

func (p *DB) GetCounter(ctx context.Context, name string) (int64, error) {
	query := `
	SELECT delta FROM counter
	WHERE name = $1
	`

	var result int64

	err := retry.WithRetry(func() error {
		row := p.db.QueryRowContext(ctx, query, name)
		return row.Scan(&result)
	}, isRetriableDB)

	if err != nil {
		return 0, err
	}

	return result, nil
}
