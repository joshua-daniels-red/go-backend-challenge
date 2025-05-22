package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/joshua-daniels-red/go-backend-challenge/ch-1/internal/stream"
)

type statusResponse struct {
	Status string `json:"status"`
}

func NewHTTPServer(addr string) *http.Server {
	stats := stream.NewStats()
	client := stream.NewWikipediaClient(stats)

	// Start streaming in background
	go func() {
		if err := client.Connect(); err != nil {
			log.Fatalf("streaming failed: %v", err)
		}
	}()

	mux := http.NewServeMux()

	// Status endpoint
	mux.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(statusResponse{Status: "ok"})
	})

	// ✅ Register /stats endpoint
	mux.HandleFunc("/stats", stats.Handler)

	return &http.Server{
		Addr:    addr,
		Handler: mux,
	}
}
