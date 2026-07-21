package handler

import (
	"encoding/json"
	models "metrics/internal/model"
	"net/http"
)

func (h *MetricsHTTPHandlers) Update(w http.ResponseWriter, r *http.Request) {
	var metric models.Metrics

	err := json.NewDecoder(r.Body).Decode(&metric)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch metric.MType {
	case models.Gauge:
		if metric.Value == nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		h.metricsService.SaveGauge(metric.ID, *metric.Value)
		if h.syncSave {
			if err := h.saver.Save(); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
		gauge, err := h.metricsService.GetGauge(metric.ID)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		updateM := models.Metrics{
			ID:    metric.ID,
			MType: models.Gauge,
			Delta: nil,
			Value: &gauge,
		}

		resp, err := json.MarshalIndent(updateM, "", "    ")
		w.Header().Set("Content-Type", "application/json")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(resp)
	case models.Counter:

		if metric.Delta == nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		h.metricsService.SaveCounter(metric.ID, *metric.Delta)
		if h.syncSave {
			if err := h.saver.Save(); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
		counter, err := h.metricsService.GetCounter(metric.ID)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		updateM := models.Metrics{
			ID:    metric.ID,
			MType: models.Counter,
			Value: nil,
			Delta: &counter,
		}

		resp, err := json.MarshalIndent(updateM, "", "    ")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(resp)
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}
