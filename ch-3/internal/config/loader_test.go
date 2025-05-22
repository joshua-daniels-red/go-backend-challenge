package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig_Success(t *testing.T) {
	content := `{
		"port": "7000",
		"stream_url": "https://example.com/stream",
		"storage": "in-memory",
		"cassandra_host": "localhost",
		"jwt_secret": "mysecret"
	}`

	tmpfile, err := os.CreateTemp("", "config*.json")
	assert.NoError(t, err)
	defer os.Remove(tmpfile.Name())

	_, err = tmpfile.Write([]byte(content))
	assert.NoError(t, err)
	tmpfile.Close()

	cfg, err := LoadConfig(tmpfile.Name())
	assert.NoError(t, err)
	assert.Equal(t, "7000", cfg.Port)
	assert.Equal(t, "https://example.com/stream", cfg.StreamURL)
	assert.Equal(t, "in-memory", cfg.Storage)
	assert.Equal(t, "localhost", cfg.CassandraHost)
	assert.Equal(t, "mysecret", cfg.JWTSecret)
}

func TestLoadConfig_FileNotFound(t *testing.T) {
	_, err := LoadConfig("nonexistent.json")
	assert.Error(t, err)
}

func TestLoadConfig_InvalidJSON(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "badconfig*.json")
	assert.NoError(t, err)
	defer os.Remove(tmpfile.Name())

	_, err = tmpfile.Write([]byte(`{bad json}`))
	assert.NoError(t, err)
	tmpfile.Close()

	_, err = LoadConfig(tmpfile.Name())
	assert.Error(t, err)
}
