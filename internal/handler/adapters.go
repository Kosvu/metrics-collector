package handler

import (
	"context"
)

type MetricsHTTPHandlers struct {
	metricsService MetricsService
	saver          Saver
	syncSave       bool
	ping           Pinger
}

type Pinger interface {
	Ping(ctx context.Context) error
}

type Saver interface {
	Save(ctx context.Context) error
}

type MetricsService interface {
	SaveGauge(ctx context.Context, name string, value float64) error
	SaveCounter(ctx context.Context, name string, value int64) error
	GetGauge(ctx context.Context, name string) (float64, error)
	GetCounter(ctx context.Context, name string) (int64, error)
	GetAll(ctx context.Context) (map[string]float64, map[string]int64, error)
}

func NewMetricsHTTPHandlers(metricService MetricsService, saver Saver, syncSave bool, ping Pinger) *MetricsHTTPHandlers {
	return &MetricsHTTPHandlers{
		metricsService: metricService,
		saver:          saver,
		syncSave:       syncSave,
		ping:           ping,
	}
}
