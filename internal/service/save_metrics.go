package service

import "context"

func (s *MetricsService) SaveGauge(ctx context.Context, name string, value float64) error {
	err := s.metricsStorage.SaveGauge(ctx, name, value)
	return err
}

func (s *MetricsService) SaveCounter(ctx context.Context, name string, value int64) error {
	err := s.metricsStorage.SaveCounter(ctx, name, value)
	return err
}
