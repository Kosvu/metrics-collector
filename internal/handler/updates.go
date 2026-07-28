package handler

import (
	"encoding/json"
	models "metrics/internal/model"
	"net/http"
)

func (h *MetricsHTTPHandlers) Updates(w http.ResponseWriter, r *http.Request) {
	var metrics []models.Metrics

	if err := json.NewDecoder(r.Body).Decode(&metrics); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := h.metricsService.SaveBatch(r.Context(), metrics); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
