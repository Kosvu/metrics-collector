package service

type MetricsService struct {
	metricsStorage MetricsStorage
}

type MetricsStorage interface {
	SaveGauge(name string, value float64)
	SaveCounter(name string, value int64)
	GetGauge(name string) (float64, error)
	GetCounter(name string) (int64, error)
	GetAll() (map[string]float64, map[string]int64)
}

func NewMetricsService(metricsStorage MetricsStorage) *MetricsService {
	return &MetricsService{
		metricsStorage: metricsStorage,
	}
}
