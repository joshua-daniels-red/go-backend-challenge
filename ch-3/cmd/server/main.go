package main

import (
	"log"
	"net/http"

	"github.com/joshua-daniels-red/go-backend-challenge/ch-3/internal/stream"
)

func main() {
	stats := stream.NewStats()
	client := stream.NewWikipediaClient(stats)

	go func() {
		if err := client.Connect(); err != nil {
			log.Fatalf("streaming failed: %v", err)
		}
	}()

	http.HandleFunc("/stats", stats.Handler)
	log.Println("HTTP server listening on :7000")
	log.Fatal(http.ListenAndServe(":7000", nil))
}