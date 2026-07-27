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

		if err := h.metricsService.SaveGauge(r.Context(), metric.ID, *metric.Value); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if h.syncSave {
			if err := h.saver.Save(r.Context()); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
		gauge, err := h.metricsService.GetGauge(r.Context(), metric.ID)

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
		if err := h.metricsService.SaveCounter(r.Context(), metric.ID, *metric.Delta); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if h.syncSave {
			if err := h.saver.Save(r.Context()); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
		counter, err := h.metricsService.GetCounter(r.Context(), metric.ID)

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
