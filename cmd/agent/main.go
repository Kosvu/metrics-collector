package main

import (
	"metrics/internal/agent"
	"net/http"
	"time"
)

func main() {
	client := &http.Client{}

	storage := agent.NewAgentStorage()
	collector := agent.NewCollector(storage)
	sender := agent.NewSender(storage, "localhost:8080", client)

	go func() {
		for {
			collector.Poll()
			time.Sleep(2 * time.Second)
		}
	}()

	for {
		sender.Send()
		time.Sleep(10 * time.Second)
	}
}
