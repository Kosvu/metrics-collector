package repository

import "context"

func (r *MemStorage) SaveGauge(ctx context.Context, name string, value float64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.gaugeMap[name] = value
	return nil
}

func (r *MemStorage) SaveCounter(ctx context.Context, name string, value int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.counterMap[name] += value
	return nil
}
