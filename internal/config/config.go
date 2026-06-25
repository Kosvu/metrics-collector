package config

import (
	"flag"
	"log"
	"os"
	"strconv"
)

type ServerConfig struct {
	Addr string
}

type AgentConfig struct {
	Addr           string
	ReportInterval int
	PollInterval   int
}

func NewServerConfig() *ServerConfig {

	cfg := &ServerConfig{}
	flag.StringVar(&cfg.Addr, "a", "localhost:8080", "server address")
	flag.Parse()

	if envAddr := os.Getenv("ADDRESS"); envAddr != "" {
		cfg.Addr = envAddr
	}

	return cfg
}

func NewAgentConfig() *AgentConfig {
	cfg := &AgentConfig{}

	flag.StringVar(&cfg.Addr, "a", "localhost:8080", "server address")
	flag.IntVar(&cfg.ReportInterval, "r", 10, "report interval, sec")
	flag.IntVar(&cfg.PollInterval, "p", 2, "poll interval, sec")

	flag.Parse()

	if v := os.Getenv("ADDRESS"); v != "" {
		cfg.Addr = v
	}
	if envRep := os.Getenv("REPORT_INTERVAL"); envRep != "" {
		envRepInt, err := strconv.Atoi(envRep)

		if err != nil {
			log.Fatal(err)
		}

		cfg.ReportInterval = envRepInt
	}
	if envPoll := os.Getenv("POLL_INTERVAL"); envPoll != "" {
		envPollInt, err := strconv.Atoi(envPoll)

		if err != nil {
			log.Fatal(err)
		}

		cfg.PollInterval = envPollInt
	}

	return cfg
}
