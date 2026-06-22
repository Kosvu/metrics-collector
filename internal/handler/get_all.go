package handler

import (
	"net/http"
	"strconv"
	"strings"
)

func (h *MetricsHTTPHandlers) GetAll(w http.ResponseWriter, r *http.Request) {
	gMap, cMap := h.metricsService.GetAll()

	var b strings.Builder
	b.WriteString("<html><body>")
	for name, value := range gMap {
		valueString := strconv.FormatFloat(value, 'f', -1, 64)
		b.WriteString("<p>" + name + ": " + valueString + "</p>")
	}
	for name, value := range cMap {
		valueString := strconv.FormatInt(value, 10)
		b.WriteString("<p>" + name + ": " + valueString + "</p>")
	}
	b.WriteString("</body></html>")

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(b.String()))
}
