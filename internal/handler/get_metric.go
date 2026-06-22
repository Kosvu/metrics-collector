package handler

import (
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

func (h *MetricsHTTPHandlers) GetMetrics(w http.ResponseWriter, r *http.Request) {
	pathType := chi.URLParam(r, "type")
	pathName := chi.URLParam(r, "name")

	if pathType != "gauge" && pathType != "counter" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if pathType == "gauge" {
		metric, err := h.metricsService.GetGauge(pathName)

		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		metricString := strconv.FormatFloat(metric, 'f', -1, 64)

		w.Header().Set("Content-Type", "text/plain")
		io.WriteString(w, metricString)
	}

	if pathType == "counter" {
		metric, err := h.metricsService.GetCounter(pathName)

		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		metricString := strconv.FormatInt(metric, 10)

		w.Header().Set("Content-Type", "text/plain")
		io.WriteString(w, metricString)
	}
}
