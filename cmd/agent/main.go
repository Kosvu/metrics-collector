package main

import (
	"metrics/internal/agent"
	"metrics/internal/config"
	"net/http"
	"time"
)

func main() {
	client := &http.Client{}
	cfg := config.NewAgentConfig()

	storage := agent.NewAgentStorage()
	collector := agent.NewCollector(storage)
	sender := agent.NewSender(storage, cfg.Addr, cfg.Key, client)

	go func() {
		for {
			collector.Poll()
			time.Sleep(time.Duration(cfg.PollInterval) * time.Second)
		}
	}()

	for {
		sender.Send()
		time.Sleep(time.Duration(cfg.ReportInterval) * time.Second)
	}
}
