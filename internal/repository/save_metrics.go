package repository

func (r *MemStorage) SaveGauge(name string, value float64) {
	r.gaugeMap[name] = value
}

func (r *MemStorage) SaveCounter(name string, value int64) {
	r.counterMap[name] += value
}
