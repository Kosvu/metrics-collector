package db

import "context"

func (p *DB) SaveGauge(ctx context.Context, name string, value float64) error {

	query :=
		`INSERT INTO gauge (name,value) VALUES ($1, $2)
		ON CONFLICT (name) DO UPDATE SET value=$2;`

	if _, err := p.db.ExecContext(ctx, query, name, value); err != nil {
		return err
	}

	return nil

}
