package db

import "context"

func (p *DB) Ping(ctx context.Context) error {
	err := p.db.PingContext(ctx)

	return err
}
