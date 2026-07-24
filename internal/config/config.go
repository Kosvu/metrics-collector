package config

import (
	"flag"
	"log"
	"os"
	"strconv"
)

type ServerConfig struct {
	Addr            string
	StoreInterval   int
	FileStoragePath string
	Restore         bool
	Dsn             string
}

type AgentConfig struct {
	Addr           string
	ReportInterval int
	PollInterval   int
}

func NewServerConfig() *ServerConfig {

	cfg := &ServerConfig{StoreInterval: 300, FileStoragePath: "/tmp/metrics-db.json", Restore: true}
	flag.StringVar(&cfg.Addr, "a", "localhost:8080", "server address")
	flag.IntVar(&cfg.StoreInterval, "i", 300, "store interval")
	flag.StringVar(&cfg.FileStoragePath, "f", "/tmp/metrics-db.json", "file storage path")
	flag.BoolVar(&cfg.Restore, "r", true, "restore")
	flag.StringVar(&cfg.Dsn, "d", "", "dsn for connect to database")
	flag.Parse()

	if envAddr := os.Getenv("ADDRESS"); envAddr != "" {
		cfg.Addr = envAddr
	}
	if envStoreInterval := os.Getenv("STORE_INTERVAL"); envStoreInterval != "" {
		envStoreIntervalInt, err := strconv.Atoi(envStoreInterval)
		if err != nil {
			log.Fatal(err)
		}
		cfg.StoreInterval = envStoreIntervalInt
	}
	if envFileStoragePath, ok := os.LookupEnv("FILE_STORAGE_PATH"); ok {
		cfg.FileStoragePath = envFileStoragePath
	}
	if envRestore := os.Getenv("RESTORE"); envRestore != "" {
		envRestoreBool, err := strconv.ParseBool(envRestore)
		if err != nil {
			log.Fatal(err)
		}
		cfg.Restore = envRestoreBool
	}
	if envDsn := os.Getenv("DATABASE_DSN"); envDsn != "" {
		cfg.Dsn = envDsn
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
