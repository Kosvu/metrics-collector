package main

import (
	"flag"
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

	addr := flag.String("a", ":8080", "The flag specifies the server address")

	flag.Parse()

	err := http.ListenAndServe(*addr, r)
	if err != nil {
		panic(err)
	}
}
