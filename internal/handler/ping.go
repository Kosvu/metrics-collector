package handler

import "net/http"

func (h *MetricsHTTPHandlers) Ping(w http.ResponseWriter, r *http.Request) {
	err := h.db.PingContext(r.Context())

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
