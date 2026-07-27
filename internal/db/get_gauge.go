package db

import "context"

func (p *DB) GetGauge(ctx context.Context, name string) (float64, error) {
	query := `
	SELECT value FROM gauge
	WHERE name = $1
	`

	var result float64

	row := p.db.QueryRowContext(ctx, query, name)

	if err := row.Scan(&result); err != nil {
		return 0, err
	}

	return result, nil
}
