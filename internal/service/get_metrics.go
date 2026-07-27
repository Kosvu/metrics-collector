package service

import "context"

func (s *MetricsService) GetGauge(ctx context.Context, name string) (float64, error) {
	return s.metricsStorage.GetGauge(ctx, name)
}

func (s *MetricsService) GetCounter(ctx context.Context, name string) (int64, error) {
	return s.metricsStorage.GetCounter(ctx, name)
}

func (s *MetricsService) GetAll(ctx context.Context) (map[string]float64, map[string]int64, error) {
	return s.metricsStorage.GetAll(ctx)
}
