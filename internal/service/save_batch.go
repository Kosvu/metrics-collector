package service

import (
	"context"
	models "metrics/internal/model"
)

func (s *MetricsService) SaveBatch(ctx context.Context, metrics []models.Metrics) error {
	return s.metricsStorage.SaveBatch(ctx, metrics)
}
