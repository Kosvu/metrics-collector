package handler

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/require"
)

type stubService struct {
	GaugeValue   float64
	CounterValue int64
	MetError     error
}

func (s *stubService) SaveGauge(name string, value float64) {
}

func (s *stubService) SaveCounter(name string, value int64) {
}

func (s *stubService) GetGauge(name string) (float64, error) {
	return s.GaugeValue, s.MetError
}
func (s *stubService) GetCounter(name string) (int64, error) {
	return s.CounterValue, s.MetError
}

func (s *stubService) GetAll() (map[string]float64, map[string]int64) {
	stubGaugeMap := make(map[string]float64)
	stubCounterMap := make(map[string]int64)

	stubGaugeMap["Alloc"] = s.GaugeValue
	stubCounterMap["Poll"] = s.CounterValue

	return stubGaugeMap, stubCounterMap
}

func newTestRouter(s stubService) chi.Router {
	h := NewMetricsHTTPHandlers(&s)
	r := chi.NewRouter()
	r.Get("/value/{type}/{name}", h.GetMetrics)
	r.Get("/", h.GetAll)
	r.Post("/update/{type}/{name}/{value}", h.SaveMetrics)

	return r
}

func testRequest(t *testing.T, ts *httptest.Server, method, path string) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, nil)
	require.NoError(t, err) // require — фатальный шаг, стопаем
	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	return resp, string(body)
}
