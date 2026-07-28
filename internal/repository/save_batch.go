package repository

import (
	"context"
	models "metrics/internal/model"
)

func (m *MemStorage) SaveBatch(ctx context.Context, metrics []models.Metrics) error {
	for _, mt := range metrics {
		switch mt.MType {
		case models.Gauge:
			if mt.Value != nil {
				if err := m.SaveGauge(ctx, mt.ID, *mt.Value); err != nil {
					return err
				}
			}
		case models.Counter:
			if mt.Delta != nil {
				if err := m.SaveCounter(ctx, mt.ID, *mt.Delta); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
