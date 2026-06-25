package service

func (s *MetricsService) GetGauge(name string) (float64, error) {
	return s.metricsStorage.GetGauge(name)
}

func (s *MetricsService) GetCounter(name string) (int64, error) {
	return s.metricsStorage.GetCounter(name)
}

func (s *MetricsService) GetAll() (map[string]float64, map[string]int64) {
	return s.metricsStorage.GetAll()
}
