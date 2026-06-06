package handler

type MetricsHTTPHandlers struct {
	metricsService SaveMetricsService
}

type SaveMetricsService interface {
	SaveGauge(name string, value float64)
	SaveCounter(name string, value int64)
}

func NewMetricsHTTPHandlers(metricService SaveMetricsService) *MetricsHTTPHandlers {
	return &MetricsHTTPHandlers{
		metricsService: metricService,
	}
}
