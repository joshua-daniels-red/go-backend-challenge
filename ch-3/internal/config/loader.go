package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Port          string `json:"port"`
	StreamURL     string `json:"stream_url"`
	Storage       string `json:"storage"`
	CassandraHost string `json:"cassandra_host"`
	JWTSecret     string `json:"jwt_secret"`
}

func LoadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var cfg Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
