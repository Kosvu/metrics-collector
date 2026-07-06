package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdate(t *testing.T) {
	ts := httptest.NewServer(newTestRouter(stubService{GaugeValue: 15, CounterValue: 15}))
	defer ts.Close()

	tests := []struct {
		name       string
		body       string
		wantStatus int
		wantBody   string
	}{
		{
			name:       "simple gauge test",
			body:       `{"id":"Alloc","type":"gauge","value":15}`,
			wantStatus: 200,
			wantBody:   `{"id":"Alloc","type":"gauge","value":15}`,
		},
		{
			name:       "simple counter test",
			body:       `{"id":"Poll","type":"counter","delta":15}`,
			wantStatus: 200,
			wantBody:   `{"id":"Poll","type":"counter","delta":15}`,
		},
		{
			name:       "invalid type",
			body:       `{"id":"Alloc","type":"unknown","value":13}`,
			wantStatus: 400,
		},
		{
			name:       "nil value",
			body:       `{"id":"Alloc","type":"gauge","delta":13}`,
			wantStatus: 400,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resp, body := testRequestBody(t, ts, http.MethodPost, "/update/", test.body)
			if test.wantBody != "" {
				assert.JSONEq(t, test.wantBody, body)
			}
			assert.Equal(t, test.wantStatus, resp.StatusCode)
		})
	}
}
