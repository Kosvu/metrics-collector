package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetJSON(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		wantStatus int
		wantBody   string
		stub       stubService
	}{
		{
			name:       "simple gauge test",
			body:       `{"id":"Alloc","type":"gauge"}`,
			wantStatus: 200,
			wantBody:   `{"id":"Alloc","type":"gauge","value":15}`,
			stub:       stubService{GaugeValue: 15, CounterValue: 15},
		},
		{
			name:       "simple counter test",
			body:       `{"id":"Poll","type":"counter"}`,
			wantStatus: 200,
			wantBody:   `{"id":"Poll","type":"counter","delta":15}`,
			stub:       stubService{GaugeValue: 15, CounterValue: 15},
		},
		{
			name:       "invalid type",
			body:       `{"id":"Poll","type":"unknown"}`,
			wantStatus: 400,
			stub:       stubService{},
		},
		{
			name:       "invalid id",
			body:       `{"id":"unknown","type":"gauge"}`,
			wantStatus: 404,
			stub:       stubService{MetError: errors.New("not found")},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ts := httptest.NewServer(newTestRouter(test.stub))
			resp, body := testRequestBody(t, ts, http.MethodPost, "/value/", test.body)
			if test.wantBody != "" {
				assert.JSONEq(t, test.wantBody, body)
			}
			assert.Equal(t, test.wantStatus, resp.StatusCode)
		})
	}
}
