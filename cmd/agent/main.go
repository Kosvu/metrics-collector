package main

import (
	"context"
	"log"
	"metrics/internal/agent"
	"metrics/internal/config"
	models "metrics/internal/model"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	client := &http.Client{}
	cfg := config.NewAgentConfig()
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	storage := agent.NewAgentStorage()
	collector := agent.NewCollector(storage)
	metricCh := make(chan models.Metrics)
	sender := agent.NewSender(storage, cfg.Addr, cfg.Key, client)
	systemCollector := agent.NewSystemCollector(storage)

	wg.Add(3 + cfg.RateLimit)

	for i := 0; i < cfg.RateLimit; i++ {
		go func() {
			defer wg.Done()
			for m := range metricCh {
				err := sender.Send(m)

				if err != nil {
					log.Println("send metric error", err)
				}
			}
		}()
	}

	go func() {
		defer wg.Done()
		ticker := time.NewTicker(time.Duration(cfg.PollInterval) * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				collector.Poll()
			case <-ctx.Done():
				return
			}
		}
	}()

	go func() {
		defer wg.Done()
		ticker := time.NewTicker(time.Duration(cfg.PollInterval) * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				systemCollector.Poll()
			case <-ctx.Done():
				return
			}
		}
	}()

	go func() {
		defer close(metricCh)
		defer wg.Done()
		ticker := time.NewTicker(time.Duration(cfg.ReportInterval) * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				metArr := sender.CollectorMetrics()

				for _, el := range metArr {
					select {
					case metricCh <- el:
						continue
					case <-ctx.Done():
						return
					}
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
	cancel()
	wg.Wait()

}
