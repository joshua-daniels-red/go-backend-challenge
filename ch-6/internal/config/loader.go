package config

import (
	"fmt"
	"os"
)

type Config struct {
	RedpandaBroker     string
	WikipediaStreamURL string
	Storage            string
	WikipediaTopic     string
}

func Load() (*Config, error) {
	cfg := &Config{
		RedpandaBroker:     os.Getenv("REDPANDA_BROKER"),
		WikipediaStreamURL: os.Getenv("WIKIPEDIA_STREAM_URL"),
		Storage:            os.Getenv("STORAGE"),
		WikipediaTopic:     os.Getenv("WIKIPEDIA_TOPIC"),
	}

	if cfg.RedpandaBroker == "" {
		return nil, fmt.Errorf("REDPANDA_BROKER must be set")
	}
	if cfg.WikipediaStreamURL == "" {
		cfg.WikipediaStreamURL = "https://stream.wikimedia.org/v2/stream/recentchange"
	}

	if cfg.WikipediaTopic == "" {
		cfg.WikipediaTopic = "wikipedia.changes"
	}

	return cfg, nil
}
