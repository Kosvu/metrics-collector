package main

import (
	"metrics/internal/config"
	"metrics/internal/handler"
	"metrics/internal/logger"
	"metrics/internal/middleware"
	"metrics/internal/repository"
	"metrics/internal/service"
	"net/http"

	"github.com/go-chi/chi"
)

func main() {
	metricsStorage := repository.NewMemStorage()
	metricsService := service.NewMetricsService(metricsStorage)
	metricsHTTP := handler.NewMetricsHTTPHandlers(metricsService)
	log := logger.NewSugarLogger()
	defer log.Sync()

	r := chi.NewRouter()
	// http://<АДРЕС_СЕРВЕРА>/update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>

	r.Use(middleware.WithLogging(log))

	r.Get("/value/{type}/{name}", metricsHTTP.GetMetrics)
	r.Get("/", metricsHTTP.GetAll)
	r.Post("/update/{type}/{name}/{value}", metricsHTTP.SaveMetrics)

	cfg := config.NewServerConfig()

	err := http.ListenAndServe(cfg.Addr, r)
	if err != nil {
		panic(err)
	}
}
