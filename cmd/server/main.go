package main

import (
	"metrics/internal/handler"
	"metrics/internal/repository"
	"metrics/internal/service"
	"net/http"
)

func main() {
	metricsStorage := repository.NewMemStorage()
	metricsService := service.NewMetricsService(metricsStorage)
	metricsHTTP := handler.NewMetricsHTTPHandlers(metricsService)

	mux := http.NewServeMux()
	// http://<АДРЕС_СЕРВЕРА>/update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>

	mux.HandleFunc("POST /update/{type}/{name}/{value}", metricsHTTP.SaveMetrics)

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
