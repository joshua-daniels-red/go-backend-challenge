package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/joshua-daniels-red/go-backend-challenge/ch-6/internal/config"
	"github.com/joshua-daniels-red/go-backend-challenge/ch-6/internal/stream"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go handleShutdown(cancel)

	cfg, err := configLoadFunc()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	log.Printf("📥 PRODUCER TOPIC: %s", cfg.WikipediaTopic)

	stream.RegisterMetrics()
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		log.Println("Prometheus metrics available at :2112/metrics")
		http.ListenAndServe(":2112", nil)
	}()

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
