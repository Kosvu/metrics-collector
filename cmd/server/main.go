package main

import (
	"context"
	"metrics/internal/config"
	"metrics/internal/db"
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
	ctx := context.Background()
	cfg := config.NewServerConfig()

	log := logger.NewSugarLogger()
	defer log.Sync()

	var metricsStorage service.MetricsStorage
	var pinger handler.Pinger

	if cfg.Dsn != "" {
		dbStorage, err := db.NewDB(ctx, cfg.Dsn)
		if err != nil {
			log.Fatal(err)
		}
		metricsStorage = dbStorage
		pinger = dbStorage
	} else {
		metricsStorage = repository.NewMemStorage()
	}
	metricsService := service.NewMetricsService(metricsStorage)
	producer := storage.NewProducer(metricsStorage, cfg.FileStoragePath)
	metricsHTTP := handler.NewMetricsHTTPHandlers(
		metricsService,
		producer,
		cfg.StoreInterval == 0,
		pinger,
	)

	r := chi.NewRouter()
	r.Use(middleware.GZipHandle())
	r.Use(middleware.WithLogging(log))

	r.Get("/value/{type}/{name}", metricsHTTP.GetMetrics)
	r.Get("/", metricsHTTP.GetAll)
	r.Post("/update/{type}/{name}/{value}", metricsHTTP.SaveMetrics)
	r.Post("/update/", metricsHTTP.Update)
	r.Post("/value/", metricsHTTP.GetJSON)
	r.Get("/ping", metricsHTTP.Ping)

	if cfg.Restore {
		if err := producer.Load(ctx); err != nil {
			log.Fatal(err)
		}
	}
	if cfg.StoreInterval > 0 {
		ticker := time.NewTicker(time.Duration(cfg.StoreInterval) * time.Second)
		defer ticker.Stop()
		go func() {
			for range ticker.C {
				if err := producer.Save(ctx); err != nil {
					log.Errorln("failed to save metrics:", err)
				}
			}
		}()
	}

	go func() {
		if err := http.ListenAndServe(cfg.Addr, r); err != nil {
			panic(err)
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
	if err := producer.Save(ctx); err != nil {
		log.Errorln("failed to save on shutdown:", err)
	}
}
