package storage

import (
	"context"
	"encoding/json"
	"io"
	models "metrics/internal/model"
	"os"
)

type MetricsStorage interface {
	GetAll(ctx context.Context) (map[string]float64, map[string]int64, error)
	SaveGauge(ctx context.Context, name string, value float64) error
	SaveCounter(ctx context.Context, name string, value int64) error
}

type Producer struct {
	storage MetricsStorage
	path    string
}

func NewProducer(storage MetricsStorage, path string) *Producer {
	return &Producer{
		storage: storage,
		path:    path,
	}
}

// Взять метрики из памяти и записать их в файл
func (p *Producer) Save(ctx context.Context) error {
	if p.path == "" {
		return nil
	}

	file, err := os.OpenFile(p.path, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0666)

	if err != nil {
		return err
	}

	defer file.Close()

	gauge, counter, err := p.storage.GetAll(ctx)

	if err != nil {
		return err
	}

	for name, value := range gauge {

		metric := models.Metrics{
			ID:    name,
			MType: models.Gauge,
			Delta: nil,
			Value: &value,
		}

		err := json.NewEncoder(file).Encode(&metric)

		if err != nil {
			return err
		}
	}

	for name, value := range counter {
		metric := models.Metrics{
			ID:    name,
			MType: models.Counter,
			Delta: &value,
			Value: nil,
		}

		err := json.NewEncoder(file).Encode(&metric)

		if err != nil {
			return err
		}
	}

	return nil
}

// При запуске прочитать метрики из файла

func (p *Producer) Load(ctx context.Context) error {
	if p.path == "" {
		return nil
	}

	file, err := os.OpenFile(p.path, os.O_RDONLY, 0666)

	if err != nil {
		if os.IsNotExist(err) {
			return nil
		} else {
			return err
		}
	}

	defer file.Close()

	decoder := json.NewDecoder(file)

	for {
		var metric models.Metrics
		err := decoder.Decode(&metric)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		switch metric.MType {
		case models.Counter:
			p.storage.SaveCounter(ctx, metric.ID, *metric.Delta)
		case models.Gauge:
			p.storage.SaveGauge(ctx, metric.ID, *metric.Value)
		}
	}

	return nil

}
