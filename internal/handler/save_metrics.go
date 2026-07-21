package handler

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

func (h *MetricsHTTPHandlers) SaveMetrics(w http.ResponseWriter, r *http.Request) {
	// спарсить название тип и значение
	// вызвать нужный метод в service
	// отдать нужный http ответ

	pathType := chi.URLParam(r, "type")
	pathName := chi.URLParam(r, "name")
	pathValue := chi.URLParam(r, "value")

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
		if h.syncSave {
			if err := h.saver.Save(); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
		w.WriteHeader(http.StatusOK)
	}

	if pathType == "counter" {
		intValue, err := strconv.ParseInt(pathValue, 10, 64)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		h.metricsService.SaveCounter(pathName, intValue)
		if h.syncSave {
			if err := h.saver.Save(); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
		w.WriteHeader(http.StatusOK)

	}

}
