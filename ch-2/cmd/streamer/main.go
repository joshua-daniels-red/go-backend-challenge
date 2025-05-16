package main

import (
	"log"

	"github.com/joshua-daniels-red/go-backend-challenge/ch-2/internal/stream"
)

func main() {
	client := stream.NewWikipediaClient()
	if err := client.ConnectAndLog(); err != nil {
		log.Fatalf("streaming failed: %v", err)
	}
}
