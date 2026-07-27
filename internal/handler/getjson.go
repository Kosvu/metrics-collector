package handler

import (
	"encoding/json"
	models "metrics/internal/model"
	"net/http"
)

func (h *MetricsHTTPHandlers) GetJSON(w http.ResponseWriter, r *http.Request) {
	var metric models.Metrics

	err := json.NewDecoder(r.Body).Decode(&metric)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch metric.MType {
	case models.Counter:
		val, err := h.metricsService.GetCounter(r.Context(), metric.ID)

		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		metric.Delta = &val
		res, err := json.MarshalIndent(metric, "", "    ")

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(res)
	case models.Gauge:
		val, err := h.metricsService.GetGauge(r.Context(), metric.ID)

		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		metric.Value = &val
		res, err := json.MarshalIndent(metric, "", "    ")

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(res)
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}
