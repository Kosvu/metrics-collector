package agent

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
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
		url := fmt.Sprintf("http://%s/update/%s/%s/%s", s.serverAddr, "gauge", name, strconv.FormatFloat(value, 'f', -1, 64))
		resp, err := s.client.Post(url, "text/plain", nil)
		if err != nil {
			log.Printf("failed to get response: %v", err)
			//continue потому что если тело не получили, то resp.Body.Close() запаникует
			continue
		}
		resp.Body.Close()
	}

	for name, value := range counter {
		url := fmt.Sprintf("http://%s/update/%s/%s/%s", s.serverAddr, "counter", name, strconv.FormatInt(value, 10))
		resp, err := s.client.Post(url, "text/plain", nil)
		if err != nil {
			log.Printf("failed to get response: %v", err)
			continue
		}
		resp.Body.Close()
	}
}
