package handler

import (
	"net/http"
	"strconv"
)

func (h *MetricsHTTPHandlers) SaveMetrics(w http.ResponseWriter, r *http.Request) {
	// спарсить название тип и значение
	// вызвать нужный метод в service
	// отдать нужный http ответ

	pathType := r.PathValue("type")
	pathName := r.PathValue("name")
	pathValue := r.PathValue("value")

	if pathName == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if pathType != "gauge" && pathType != "counter" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if pathType == "gauge" {
		floatValue, err := strconv.ParseFloat(pathValue, 64)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		h.metricsService.SaveGauge(pathName, floatValue)
		w.WriteHeader(http.StatusOK)
	}

	if pathType == "counter" {
		intValue, err := strconv.ParseInt(pathValue, 10, 64)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		h.metricsService.SaveCounter(pathName, intValue)
		w.WriteHeader(http.StatusOK)

	}

}
