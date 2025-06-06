package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joshua-daniels-red/go-backend-challenge/ch-7/internal/config"
	"github.com/joshua-daniels-red/go-backend-challenge/ch-7/internal/stream"
)

var (
	configLoadFunc            = config.Load
	streamWikipediaEventsFunc = stream.StreamWikipediaEvents
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("fatal error: %v", err)
	}
}

func run() error {
	// Avoid duplicate registration in tests
	if flag.Lookup("test.v") == nil {
		stream.RegisterMetrics()
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go handleShutdown(cancel)

	cfg, err := configLoadFunc()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	log.Printf("📥 PRODUCER TOPIC: %s", cfg.WikipediaTopic)

	if err := streamWikipediaEventsFunc(ctx, cfg.RedpandaBroker, cfg.WikipediaStreamURL, cfg.WikipediaTopic); err != nil {
		return fmt.Errorf("streaming failed: %w", err)
	}

	return nil
}

func handleShutdown(cancel context.CancelFunc) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
	cancel()
}
