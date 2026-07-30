package db

import (
	"context"
	"metrics/internal/retry"
)

type gauge struct {
	name  string
	value float64
}

type counter struct {
	name  string
	delta int64
}

func (p *DB) GetAll(ctx context.Context) (map[string]float64, map[string]int64, error) {
	resultGauge := make(map[string]float64)
	resultCounter := make(map[string]int64)

	queryGauge :=
		`
	SELECT name,value FROM gauge
	`

	queryCounter :=
		`
	SELECT name,delta FROM counter
	`

	err := retry.WithRetry(func() error {
		resultGauge = make(map[string]float64)
		rowsGauge, err := p.db.QueryContext(ctx, queryGauge)
		if err != nil {
			return err
		}
		defer rowsGauge.Close()

		for rowsGauge.Next() {
			var g gauge

			if err := rowsGauge.Scan(&g.name, &g.value); err != nil {
				return err
			}

			resultGauge[g.name] = g.value
		}

		return rowsGauge.Err()

	}, isRetriableDB)

	if err != nil {
		return nil, nil, err
	}

	err = retry.WithRetry(func() error {
		resultCounter = make(map[string]int64)
		rowsCounter, err := p.db.QueryContext(ctx, queryCounter)
		if err != nil {
			return err
		}
		defer rowsCounter.Close()

		for rowsCounter.Next() {
			var c counter

			if err := rowsCounter.Scan(&c.name, &c.delta); err != nil {
				return err
			}

			resultCounter[c.name] = c.delta
		}

		return rowsCounter.Err()
	}, isRetriableDB)

	if err != nil {
		return nil, nil, err
	}

	return resultGauge, resultCounter, nil
}
