package main

import (
	"context"
	"errors"
	"testing"

	"github.com/joshua-daniels-red/go-backend-challenge/ch-5/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestRun_Success(t *testing.T) {
	originalLoad := configLoadFunc
	originalStream := streamWikipediaEventsFunc
	t.Cleanup(func() {
		configLoadFunc = originalLoad
		streamWikipediaEventsFunc = originalStream
	})

	configLoadFunc = func() (*config.Config, error) {
		return &config.Config{
			RedpandaBroker:     "test-broker",
			WikipediaStreamURL: "http://test-stream",
		}, nil
	}
	streamWikipediaEventsFunc = func(_ context.Context, _, _ string) error {
		return nil
	}

	err := run()
	assert.NoError(t, err)
}

func TestRun_ConfigLoadFails(t *testing.T) {
	original := configLoadFunc
	t.Cleanup(func() { configLoadFunc = original })

	configLoadFunc = func() (*config.Config, error) {
		return nil, errors.New("load error")
	}

	err := run()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to load config")
}

func TestRun_StreamFails(t *testing.T) {
	originalLoad := configLoadFunc
	originalStream := streamWikipediaEventsFunc
	t.Cleanup(func() {
		configLoadFunc = originalLoad
		streamWikipediaEventsFunc = originalStream
	})

	configLoadFunc = func() (*config.Config, error) {
		return &config.Config{
			RedpandaBroker:     "test-broker",
			WikipediaStreamURL: "http://test-stream",
		}, nil
	}
	streamWikipediaEventsFunc = func(_ context.Context, _, _ string) error {
		return errors.New("kafka error")
	}

	err := run()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "streaming failed")
}