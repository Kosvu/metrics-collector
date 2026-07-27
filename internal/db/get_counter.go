package db

import "context"

func (p *DB) GetCounter(ctx context.Context, name string) (int64, error) {
	query := `
	SELECT delta FROM counter
	WHERE name = $1
	`

	var result int64

	row := p.db.QueryRowContext(ctx, query, name)

	if err := row.Scan(&result); err != nil {
		return 0, err
	}

	return result, nil
}
