package repository

func (r *MemStorage) SaveGauge(name string, value float64) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.gaugeMap[name] = value
}

func (r *MemStorage) SaveCounter(name string, value int64) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.counterMap[name] += value
}
