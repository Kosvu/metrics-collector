package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type stubService struct{}

func (s *stubService) SaveGauge(name string, value float64) {
}

func (s *stubService) SaveCounter(name string, value int64) {
}

func TestSaveMetrics(t *testing.T) {
	tests := []struct {
		name        string
		metricType  string
		metricName  string
		metricValue string
		want        int
	}{
		{
			name:        "empty name",
			metricType:  "gauge",
			metricName:  "",
			metricValue: "1",
			want:        404,
		},
		{
			name:        "invalid type",
			metricType:  "unknown",
			metricName:  "Alloc",
			metricValue: "1",
			want:        400,
		},
		{
			name:        "invalid gauge value",
			metricType:  "gauge",
			metricName:  "Alloc",
			metricValue: "abc",
			want:        400,
		},
		{
			name:        "invalid counter value",
			metricType:  "counter",
			metricName:  "PollCount",
			metricValue: "abc",
			want:        400,
		},
		{
			name:        "valid gauge test",
			metricType:  "gauge",
			metricName:  "Alloc",
			metricValue: "42.5",
			want:        200,
		},
		{
			name:        "valid counter test",
			metricType:  "counter",
			metricName:  "PollCount",
			metricValue: "10",
			want:        200,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			handler := NewMetricsHTTPHandlers(&stubService{})
			req := httptest.NewRequest(http.MethodPost, "/", nil)
			req.SetPathValue("type", test.metricType)
			req.SetPathValue("name", test.metricName)
			req.SetPathValue("value", test.metricValue)

			rec := httptest.NewRecorder()
			handler.SaveMetrics(rec, req)
			assert.Equal(t, test.want, rec.Code)
		})
	}
}
