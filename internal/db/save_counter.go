package db

import "context"

func (p *DB) SaveCounter(ctx context.Context, name string, value int64) error {
	query := `
	INSERT INTO counter (name, delta) VALUES ($1, $2)
	ON CONFLICT (name) DO UPDATE SET delta = counter.delta + $2;
	`

	if _, err := p.db.ExecContext(ctx, query, name, value); err != nil {
		return err
	}

	return nil
}
