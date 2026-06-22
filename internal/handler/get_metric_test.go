package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetMetric(t *testing.T) {
	tests := []struct {
		name         string
		path         string
		gaugeValue   float64
		counterValue int64
		err          error // общая ошибка для 404 (или раздельные, если хочешь)
		wantStatus   int
		wantBody     string
	}{
		{name: "gauge success", path: "/value/gauge/Alloc", gaugeValue: 40, wantStatus: 200, wantBody: "40"},
		{name: "counter success", path: "/value/counter/Poll", counterValue: 5, wantStatus: 200, wantBody: "5"},
		{name: "not found", path: "/value/gauge/unknown", err: errors.New("nf"), wantStatus: 404},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			stub := stubService{GaugeValue: test.gaugeValue, CounterValue: test.counterValue, MetError: test.err}
			ts := httptest.NewServer(newTestRouter(stub))
			defer ts.Close()
			resp, body := testRequest(t, ts, http.MethodGet, test.path)
			if test.wantBody != "" {
				assert.Equal(t, test.wantBody, body)
			}
			assert.Equal(t, test.wantStatus, resp.StatusCode)
		})
	}
}
