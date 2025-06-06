package main

import (
	"context"
	"errors"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/gocql/gocql"
	"github.com/joshua-daniels-red/go-backend-challenge/ch-7/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/twmb/franz-go/pkg/kgo"
)

func TestRun_ConfigFails(t *testing.T) {
	configLoadFunc = func() (*config.Config, error) {
		return nil, errors.New("config load failed")
	}
	err := run()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to load config")
}

func TestRun_KafkaClientFails(t *testing.T) {
	configLoadFunc = func() (*config.Config, error) {
		return &config.Config{
			RedpandaBroker:     "any",
			WikipediaStreamURL: "any",
			Storage:            "in-memory",
		}, nil
	}

	newKafkaClientFunc = func(...kgo.Opt) (*kgo.Client, error) {
		return nil, errors.New("kafka client failed")
	}

	err := run()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create Kafka client")
}

func TestHandleShutdown(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	sigCh := make(chan os.Signal, 1)
	go handleShutdown(cancel, sigCh)
	sigCh <- syscall.SIGTERM

	time.Sleep(100 * time.Millisecond)
	select {
	case <-ctx.Done():
		assert.True(t, true)
	default:
		t.Fatal("Context was not cancelled")
	}
}

func TestRun_CassandraSessionFails(t *testing.T) {
	configLoadFunc = func() (*config.Config, error) {
		return &config.Config{
			RedpandaBroker:     "b",
			WikipediaStreamURL: "u",
			Storage:            "cassandra",
		}, nil
	}

	newKafkaClientFunc = func(opts ...kgo.Opt) (*kgo.Client, error) {
		return kgo.NewClient(
			kgo.SeedBrokers("localhost:12345"),     // bogus broker
			kgo.DialTimeout(10*time.Millisecond),   // fail fast
			kgo.ProducerLinger(5*time.Millisecond), // fast retry config
		)
	}

	newCassandraSessionFn = func() (*gocql.Session, error) {
		return nil, errors.New("cassandra boom")
	}

	err := run()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to connect to Cassandra")
}
