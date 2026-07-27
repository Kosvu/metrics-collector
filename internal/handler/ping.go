package handler

import "net/http"

func (h *MetricsHTTPHandlers) Ping(w http.ResponseWriter, r *http.Request) {
	if err := h.ping.Ping(r.Context()); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
