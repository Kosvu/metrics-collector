package handler

type MetricsHTTPHandlers struct {
	metricsService MetricsService
	saver          Saver
	syncSave       bool
}

type Saver interface {
	Save() error
}

type MetricsService interface {
	SaveGauge(name string, value float64)
	SaveCounter(name string, value int64)
	GetGauge(name string) (float64, error)
	GetCounter(name string) (int64, error)
	GetAll() (map[string]float64, map[string]int64)
}

func NewMetricsHTTPHandlers(metricService MetricsService, saver Saver, syncSave bool) *MetricsHTTPHandlers {
	return &MetricsHTTPHandlers{
		metricsService: metricService,
		saver:          saver,
		syncSave:       syncSave,
	}
}
