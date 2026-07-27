package service

import "context"

type MetricsService struct {
	metricsStorage MetricsStorage
}

type MetricsStorage interface {
	SaveGauge(ctx context.Context, name string, value float64) error
	SaveCounter(ctx context.Context, name string, value int64) error
	GetGauge(ctx context.Context, name string) (float64, error)
	GetCounter(ctx context.Context, name string) (int64, error)
	GetAll(ctx context.Context) (map[string]float64, map[string]int64, error)
}

func NewMetricsService(metricsStorage MetricsStorage) *MetricsService {
	return &MetricsService{
		metricsStorage: metricsStorage,
	}
}
