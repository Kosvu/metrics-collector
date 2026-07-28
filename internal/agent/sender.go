package agent

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	models "metrics/internal/model"
	"net/http"
)

type Sender struct {
	metricsReader MetricsReader
	serverAddr    string
	client        *http.Client
}

type MetricsReader interface {
	GetAll() (map[string]float64, map[string]int64)
}

func NewSender(metricsReader MetricsReader, serverAddr string, client *http.Client) *Sender {
	return &Sender{
		metricsReader: metricsReader,
		serverAddr:    serverAddr,
		client:        client,
	}
}

func (s *Sender) Send() {

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

	if len(metrics) == 0 {
		return
	}

	var bf bytes.Buffer

	gz := gzip.NewWriter(&bf)

	if err := json.NewEncoder(gz).Encode(metrics); err != nil {
		return
	}

	gz.Close()

	url := fmt.Sprintf("http://%s/updates/", s.serverAddr)
	req, err := http.NewRequest(http.MethodPost, url, &bf)

	if err != nil {
		return
	}

	req.Header.Set("Content-Encoding", "gzip")
	resp, err := s.client.Do(req)
	if err != nil {
		return
	}

	resp.Body.Close()
}
