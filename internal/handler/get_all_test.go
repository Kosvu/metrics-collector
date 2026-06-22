package handler

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAll(t *testing.T) {
	tests := []struct {
		name         string
		path         string
		gaugeValue   float64
		counterValue int64
		wantStatus   int
	}{
		{
			name:         "simple test",
			path:         "/",
			gaugeValue:   40,
			counterValue: 50,
			wantStatus:   200,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := stubService{
				GaugeValue:   test.gaugeValue,
				CounterValue: test.counterValue,
			}
			ts := httptest.NewServer(newTestRouter(s))

			resp, body := testRequest(t, ts, http.MethodGet, test.path)
			assert.Contains(t, body, strconv.FormatFloat(test.gaugeValue, 'f', -1, 64))
			assert.Contains(t, body, strconv.FormatInt(test.counterValue, 10))
			assert.Equal(t, test.wantStatus, resp.StatusCode)
		})
	}
}
