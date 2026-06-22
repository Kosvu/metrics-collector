package main

import (
	"metrics/internal/handler"
	"metrics/internal/repository"
	"metrics/internal/service"
	"net/http"

	"github.com/go-chi/chi"
)

func main() {
	metricsStorage := repository.NewMemStorage()
	metricsService := service.NewMetricsService(metricsStorage)
	metricsHTTP := handler.NewMetricsHTTPHandlers(metricsService)

	r := chi.NewRouter()
	// http://<АДРЕС_СЕРВЕРА>/update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>

	r.Get("/value/{type}/{name}", metricsHTTP.GetMetrics)
	r.Get("/", metricsHTTP.GetAll)
	r.Post("/update/{type}/{name}/{value}", metricsHTTP.SaveMetrics)

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		panic(err)
	}
}
