package config

import (
	"log"
	"os"
)

type Config struct {
	RedpandaBroker     string
	WikipediaStreamURL string
	Storage			   string
}

func Load() *Config {
	cfg := &Config{
		RedpandaBroker:     os.Getenv("REDPANDA_BROKER"),
		WikipediaStreamURL: os.Getenv("WIKIPEDIA_STREAM_URL"),
	}

	if cfg.RedpandaBroker == "" {
		log.Fatal("REDPANDA_BROKER must be set")
	}
	if cfg.WikipediaStreamURL == "" {
		cfg.WikipediaStreamURL = "https://stream.wikimedia.org/v2/stream/recentchange"
	}

	return cfg
}
