package agent

import (
	"bytes"
	"compress/gzip"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"metrics/internal/hash"
	models "metrics/internal/model"
	"metrics/internal/retry"
	"net/http"
)

type Sender struct {
	metricsReader MetricsReader
	serverAddr    string
	client        *http.Client
	key           string
}

type MetricsReader interface {
	GetAll() (map[string]float64, map[string]int64)
}

func NewSender(metricsReader MetricsReader, serverAddr string, key string, client *http.Client) *Sender {
	return &Sender{
		metricsReader: metricsReader,
		serverAddr:    serverAddr,
		client:        client,
		key:           key,
	}
}

func (s *Sender) CollectorMetrics() []models.Metrics {

	var metrics []models.Metrics
	gauge, counter := s.metricsReader.GetAll()

	for name, value := range gauge {
		v := &value
		metrics = append(metrics, models.Metrics{ID: name, MType: models.Gauge, Value: v})
	}

	for name, delta := range counter {
		d := &delta
		metrics = append(metrics, models.Metrics{ID: name, MType: models.Counter, Delta: d})
	}

	return metrics
}

func (s *Sender) Send(metric models.Metrics) error {
	b, err := json.Marshal(metric)

	if err != nil {
		return err
	}

	var signString string

	if s.key != "" {
		sign := hash.Sign(b, []byte(s.key))
		signString = hex.EncodeToString(sign)
	}

	return retry.WithRetry(func() error {
		var bf bytes.Buffer

		gz := gzip.NewWriter(&bf)

		_, err := gz.Write(b)
		if err != nil {
			return err
		}

		gz.Close()

		url := fmt.Sprintf("http://%s/update/", s.serverAddr)
		req, err := http.NewRequest(http.MethodPost, url, &bf)

		if err != nil {
			return err
		}

		req.Header.Set("Content-Encoding", "gzip")
		if signString != "" {
			req.Header.Set("HashSHA256", signString)
		}
		resp, err := s.client.Do(req)
		if err != nil {
			return err
		}

		resp.Body.Close()
		return nil
	}, isRetriableNet)
}
