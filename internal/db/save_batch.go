package db

import (
	"context"
	models "metrics/internal/model"
	"metrics/internal/retry"
)

func (p *DB) SaveBatch(ctx context.Context, metrics []models.Metrics) error {

	return retry.WithRetry(func() error {
		tx, err := p.db.BeginTx(ctx, nil)
		if err != nil {
			return err
		}
		defer tx.Rollback()

		for _, m := range metrics {
			switch m.MType {
			case models.Counter:
				if m.Delta != nil {

					query := `
					INSERT INTO counter (name, delta) VALUES ($1, $2)
					ON CONFLICT (name) DO UPDATE SET delta = counter.delta + $2;
					`
					if _, err := tx.ExecContext(ctx, query, m.ID, m.Delta); err != nil {
						return err
					}

				}
			case models.Gauge:
				if m.Value != nil {
					query :=
						`INSERT INTO gauge (name,value) VALUES ($1, $2)
					ON CONFLICT (name) DO UPDATE SET value=$2;`

					if _, err := tx.ExecContext(ctx, query, m.ID, m.Value); err != nil {
						return err
					}
				}
			}
		}

		return tx.Commit()
	}, isRetriableDB)
}
