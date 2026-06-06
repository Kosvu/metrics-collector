package service

func (s *MetricsService) SaveGauge(name string, value float64) {
	s.metricsStorage.SaveGauge(name, value)
}

func (s *MetricsService) SaveCounter(name string, value int64) {
	s.metricsStorage.SaveCounter(name, value)
}
