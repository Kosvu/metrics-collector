package db

import (
	"context"
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type DB struct {
	db  *sql.DB
	dsn string
}

func NewDB(ctx context.Context, dsn string) (*DB, error) {
	db, err := sql.Open("pgx", dsn)

	if err != nil {
		return nil, err
	}

	queryCounter := `
	CREATE TABLE IF NOT EXISTS counter (
		name TEXT PRIMARY KEY,
		delta bigint
	);
	`

	queryGauge := `
	CREATE TABLE IF NOT EXISTS gauge (
		name TEXT PRIMARY KEY,
		value double precision
	);
	`

	_, err = db.ExecContext(
		ctx,
		queryCounter,
	)

	if err != nil {
		return nil, err
	}

	_, err = db.ExecContext(
		ctx,
		queryGauge,
	)

	if err != nil {
		return nil, err
	}

	return &DB{db: db, dsn: dsn}, nil
}
