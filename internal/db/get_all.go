package db

import "context"

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

	rowsGauge, err := p.db.QueryContext(ctx, queryGauge)

	if err != nil {
		return nil, nil, err
	}

	for rowsGauge.Next() {
		var g gauge

		if err := rowsGauge.Scan(&g.name, &g.value); err != nil {
			return nil, nil, err
		}

		resultGauge[g.name] = g.value
	}

	if err := rowsGauge.Err(); err != nil {
		return nil, nil, err
	}

	rowsGauge.Close()

	rowsCounter, err := p.db.QueryContext(ctx, queryCounter)

	if err != nil {
		return nil, nil, err
	}

	for rowsCounter.Next() {
		var c counter

		if err := rowsCounter.Scan(&c.name, &c.delta); err != nil {
			return nil, nil, err
		}

		resultCounter[c.name] = c.delta
	}

	if err := rowsCounter.Err(); err != nil {
		return nil, nil, err
	}

	rowsCounter.Close()

	return resultGauge, resultCounter, nil
}
