package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Port 	  string `json:"port"`
	StreamURL string `json:"stream_url"`
}

func LoadConfig (path string) (*Config, error){
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	var cfg Config
	if err := json.NewDecoder(file).Decode(&cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &cfg, nil
}