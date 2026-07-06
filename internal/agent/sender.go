package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
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
	gauge, counter := s.metricsReader.GetAll()

	for name, value := range gauge {

		var buf bytes.Buffer

		url := fmt.Sprintf("http://%s/update/", s.serverAddr)

		metricS := models.Metrics{
			ID:    name,
			MType: models.Gauge,
			Delta: nil,
			Value: &value,
		}

		err := json.NewEncoder(&buf).Encode(metricS)
		if err != nil {
			log.Printf("failed model encode")
			continue
		}

		resp, err := s.client.Post(url, "application/json", &buf)
		if err != nil {
			log.Printf("failed to get response: %v", err)
			//continue потому что если тело не получили, то resp.Body.Close() запаникует
			continue
		}
		resp.Body.Close()
	}

	for name, value := range counter {
		var buf bytes.Buffer

		url := fmt.Sprintf("http://%s/update/", s.serverAddr)

		metricS := models.Metrics{
			ID:    name,
			MType: models.Counter,
			Delta: &value,
			Value: nil,
		}

		err := json.NewEncoder(&buf).Encode(metricS)
		if err != nil {
			log.Printf("failed model encode")
			continue
		}

		resp, err := s.client.Post(url, "application/json", &buf)
		if err != nil {
			log.Printf("failed to get response: %v", err)
			//continue потому что если тело не получили, то resp.Body.Close() запаникует
			continue
		}
		resp.Body.Close()
	}
}
