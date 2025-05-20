package main

import (
	"log"
	"net/http"

	"github.com/joshua-daniels-red/go-backend-challenge/ch-3/internal/config"
	"github.com/joshua-daniels-red/go-backend-challenge/ch-3/internal/stream"
)

func main() {
	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	stats := stream.NewStats()
	client := stream.NewWikipediaClient(stats, cfg.StreamURL)

	go func() {
		if err := client.Connect(); err != nil {
			log.Fatalf("streaming failed: %v", err)
		}
	}()

	http.HandleFunc("/stats", stats.Handler)
	log.Printf("HTTP server listening on %s", cfg.Port)
	log.Fatal(http.ListenAndServe(cfg.Port, nil))
}
