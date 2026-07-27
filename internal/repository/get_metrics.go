package repository

import (
	"context"
	"fmt"
)

func (r *MemStorage) GetGauge(ctx context.Context, name string) (float64, error) {

	r.mu.RLock()
	defer r.mu.RUnlock()

	value, ok := r.gaugeMap[name]

	if !ok {
		return 0, fmt.Errorf("metric not found")
	}

	return value, nil
}

func (r *MemStorage) GetCounter(ctx context.Context, name string) (int64, error) {

	r.mu.RLock()
	defer r.mu.RUnlock()

	value, ok := r.counterMap[name]

	if !ok {
		return 0, fmt.Errorf("metric not found")
	}

	return value, nil
}

func (r *MemStorage) GetAll(ctx context.Context) (map[string]float64, map[string]int64, error) {

	r.mu.RLock()
	defer r.mu.RUnlock()

	tmpGauge := make(map[string]float64)
	tmpCounter := make(map[string]int64)

	for key, value := range r.counterMap {
		tmpCounter[key] = value
	}

	for key, value := range r.gaugeMap {
		tmpGauge[key] = value
	}

	return tmpGauge, tmpCounter, nil
}
