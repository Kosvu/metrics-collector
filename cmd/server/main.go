package main

import (
	"database/sql"
	"metrics/internal/config"
	"metrics/internal/handler"
	"metrics/internal/logger"
	"metrics/internal/middleware"
	"metrics/internal/repository"
	"metrics/internal/service"
	"metrics/internal/storage"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	metricsStorage := repository.NewMemStorage()
	metricsService := service.NewMetricsService(metricsStorage)
	cfg := config.NewServerConfig()
	db, err := sql.Open("pgx", cfg.Dsn)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	producer := storage.NewProducer(metricsStorage, cfg.FileStoragePath)
	metricsHTTP := handler.NewMetricsHTTPHandlers(metricsService, producer, cfg.StoreInterval == 0, db)
	log := logger.NewSugarLogger()
	defer log.Sync()

	r := chi.NewRouter()
	// http://<АДРЕС_СЕРВЕРА>/update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>

	r.Use(middleware.GZipHandle())
	r.Use(middleware.WithLogging(log))

	r.Get("/value/{type}/{name}", metricsHTTP.GetMetrics)
	r.Get("/", metricsHTTP.GetAll)
	r.Post("/update/{type}/{name}/{value}", metricsHTTP.SaveMetrics)
	r.Post("/update/", metricsHTTP.Update)
	r.Post("/value/", metricsHTTP.GetJSON)
	r.Get("/ping", metricsHTTP.Ping)

	if cfg.Restore {
		err := producer.Load()

		if err != nil {
			log.Fatal(err)
		}
	}

	if cfg.StoreInterval > 0 {
		ticker := time.NewTicker(time.Duration(cfg.StoreInterval) * time.Second)
		defer ticker.Stop()

		go func() {
			for range ticker.C {
				err := producer.Save()
				if err != nil {
					//логирование ошибки
					log.Errorln("failed to save metrics:", err)
					continue
				}
			}
		}()
	}

	go func() {
		err := http.ListenAndServe(cfg.Addr, r)
		if err != nil {
			panic(err)
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
	if err := producer.Save(); err != nil {
		log.Errorln("failed to save on shutdown:", err)
	}
}
