package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSaveMetrics(t *testing.T) {

	ts := httptest.NewServer(newTestRouter(stubService{}))
	defer ts.Close()

	tests := []struct {
		name string

		path       string
		wantStatus int
	}{
		// {name: "known gauge", path: "/value/gauge/Alloc", wantStatus: 200, wantBody: "42.5"},
		// update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>
		{
			name:       "empty name",
			path:       "/update/gauge//1",
			wantStatus: 404,
		},
		{
			name:       "invalid type",
			path:       "/update/unknown/Alloc/1",
			wantStatus: 400,
		},
		{
			name:       "invalid gauge value",
			path:       "/update/gauge/Alloc/abc",
			wantStatus: 400,
		},
		{
			name:       "invalid counter value",
			path:       "/update/counter/PollCount/abc",
			wantStatus: 400,
		},
		{
			name:       "valid gauge test",
			path:       "/update/gauge/Alloc/42.5",
			wantStatus: 200,
		},
		{
			name:       "valid counter test",
			path:       "/update/counter/PollCount/10",
			wantStatus: 200,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resp, _ := testRequest(t, ts, http.MethodPost, test.path)
			assert.Equal(t, test.wantStatus, resp.StatusCode)
		})
	}
}
