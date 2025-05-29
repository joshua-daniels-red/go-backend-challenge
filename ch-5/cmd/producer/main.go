package main

import (
	"bufio"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/joshua-daniels-red/go-backend-challenge/ch-5/internal/config"
	"github.com/joshua-daniels-red/go-backend-challenge/ch-5/internal/stream"
	"github.com/twmb/franz-go/pkg/kgo"
)

const topic = "wikipedia.changes"

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go handleShutdown(cancel)

	cfg := config.Load()
	log.Printf("Loaded REDPANDA_BROKER: %s", cfg.RedpandaBroker)

	// Create Redpanda (Kafka-compatible) producer client
	client, err := kgo.NewClient(
		kgo.SeedBrokers(cfg.RedpandaBroker),
		kgo.ProducerLinger(100*time.Millisecond),
	)
	if err != nil {
		log.Fatalf("failed to create Redpanda client: %v", err)
	}
	defer client.Close()

	// Connect to Wikipedia stream
	resp, err := http.Get(cfg.WikipediaStreamURL)
	if err != nil {
		log.Fatalf("failed to connect to Wikipedia stream: %v", err)
	}
	defer resp.Body.Close()

	log.Println("Connected to Wikipedia stream. Streaming events...")

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			log.Println("Producer shutting down...")
			return
		default:
			line := scanner.Text()
			if !strings.HasPrefix(line, "data: ") {
				continue
			}

			var raw map[string]interface{}
			if err := json.Unmarshal([]byte(line[6:]), &raw); err != nil {
				log.Printf("skipping malformed event: %v", err)
				continue
			}

			metaMap, ok := raw["meta"].(map[string]interface{})
			if !ok {
				log.Println("skipping event: missing meta")
				continue
			}

			domain, _ := metaMap["domain"].(string)
			title, _ := raw["title"].(string)
			user, _ := raw["user"].(string)

			if domain == "" || title == "" || user == "" {
				log.Println("skipping incomplete event")
				continue
			}

			event := stream.Event{
				Domain: domain,
				Title:  title,
				User:   user,
			}

			data, err := json.Marshal(event)
			if err != nil {
				log.Printf("failed to marshal event: %v", err)
				continue
			}

			record := &kgo.Record{
				Topic: topic,
				Value: data,
			}

			client.Produce(ctx, record, func(_ *kgo.Record, err error) {
				if err != nil {
					log.Printf("failed to produce message: %v", err)
				}
			})
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("error reading stream: %v", err)
	}
}

func handleShutdown(cancel context.CancelFunc) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
	cancel()
}
