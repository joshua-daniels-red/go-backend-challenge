package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/joshua-daniels-red/go-backend-challenge/ch-2/internal/config"
	"github.com/joshua-daniels-red/go-backend-challenge/ch-2/internal/stream"
)

type statusResponse struct {
	Status string `json:"status"`
}

func NewHTTPServer(cfg *config.Config) *http.Server {
	stats := stream.NewStats()
	client := stream.NewWikipediaClient(stats, cfg.StreamURL)

	go func() {
		if err := client.Connect(); err != nil {
			log.Fatalf("streaming failed: %v", err)
		}
	}()

	mux := http.NewServeMux()
	mux.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(statusResponse{Status: "ok"})
	})
	mux.HandleFunc("/stats", stats.Handler)

	return &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: mux,
	}
}
