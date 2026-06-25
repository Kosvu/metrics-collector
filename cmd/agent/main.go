package main

import (
	"flag"
	"metrics/internal/agent"
	"net/http"
	"time"
)

func main() {
	client := &http.Client{}

	addr := flag.String("a", ":8080", "The flag specifies the agent address")
	rep := flag.Int("r", 10, "frequency of metric submissions")
	repVal := *rep
	poll := flag.Int("p", 2, "metrics polling frequency")
	pollVal := *poll

	flag.Parse()

	storage := agent.NewAgentStorage()
	collector := agent.NewCollector(storage)
	sender := agent.NewSender(storage, *addr, client)

	go func() {
		for {
			collector.Poll()
			time.Sleep(time.Duration(repVal) * time.Second)
		}
	}()

	for {
		sender.Send()
		time.Sleep(time.Duration(pollVal) * time.Second)
	}
}
