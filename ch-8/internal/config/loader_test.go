package config_test

import (
	"os"
	"testing"

	"github.com/joshua-daniels-red/go-backend-challenge/ch-8/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestLoad_ConfigWithAllEnvVars(t *testing.T) {
	os.Setenv("REDPANDA_BROKER", "localhost:9092")
	os.Setenv("WIKIPEDIA_STREAM_URL", "http://custom.wikimedia.stream")
	os.Setenv("STORAGE", "in-memory")

	cfg, err := config.Load()
	assert.NoError(t, err)

	assert.Equal(t, "localhost:9092", cfg.RedpandaBroker)
	assert.Equal(t, "http://custom.wikimedia.stream", cfg.WikipediaStreamURL)
	assert.Equal(t, "in-memory", cfg.Storage)
}

func TestLoad_DefaultWikipediaStreamURL(t *testing.T) {
	os.Setenv("REDPANDA_BROKER", "localhost:9092")
	os.Unsetenv("WIKIPEDIA_STREAM_URL")
	os.Setenv("STORAGE", "cassandra")

	cfg, err := config.Load()
	assert.NoError(t, err)

	assert.Equal(t, "localhost:9092", cfg.RedpandaBroker)
	assert.Equal(t, "https://stream.wikimedia.org/v2/stream/recentchange", cfg.WikipediaStreamURL)
	assert.Equal(t, "cassandra", cfg.Storage)
}

func TestLoad_MissingRedpandaBrokerReturnsError(t *testing.T) {
	os.Unsetenv("REDPANDA_BROKER")
	os.Setenv("STORAGE", "cassandra")

	cfg, err := config.Load()
	assert.Nil(t, cfg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "REDPANDA_BROKER must be set")
}
