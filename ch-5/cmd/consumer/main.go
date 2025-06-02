package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joshua-daniels-red/go-backend-challenge/ch-5/internal/config"
	"github.com/joshua-daniels-red/go-backend-challenge/ch-5/internal/stream"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/gocql/gocql"
)

const topic = "wikipedia.changes"

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go handleShutdown(cancel)

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	log.Printf("Consumer connecting to broker: %s", cfg.RedpandaBroker)

	client, err := kgo.NewClient(
		kgo.SeedBrokers(cfg.RedpandaBroker),
		kgo.ConsumeTopics(topic),
		kgo.ConsumerGroup("wikipedia-consumer-group"),
		kgo.MaxConcurrentFetches(5), // batch/fetch tuning
	)
	if err != nil {
		log.Fatalf("failed to create Redpanda client: %v", err)
	}
	defer func() {
		log.Println("Flushing and closing Redpanda client...")
		client.Close()
	}()

	var store stream.StatsStore

	if cfg.Storage == "cassandra" {
		cluster := gocql.NewCluster("cassandra")
		cluster.Keyspace = "goanalytics"
		cluster.Consistency = gocql.Quorum
		cluster.ConnectTimeout = 5 * time.Second
		session, err := cluster.CreateSession()
		if err != nil {
			log.Fatalf("failed to connect to Cassandra: %v", err)
		}
		defer session.Close()
		store = stream.NewCassandraStats(session)
	} else {
		store = stream.NewInMemoryStats()
	}

	go startHTTPServer(store)

	log.Println("Consumer started. Waiting for messages...")

	for {
		select {
		case <-ctx.Done():
			log.Println("Shutting down consumer...")
			return
		default:
			fetches := client.PollFetches(ctx)
			fetches.EachPartition(func(p kgo.FetchTopicPartition) {
				log.Printf("Fetched %d records from partition %s", len(p.Records), p.Topic)
				for _, record := range p.Records {
					var event stream.Event
					if err := json.Unmarshal(record.Value, &event); err != nil {
						log.Printf("failed to decode event: %v", err)
						continue
					}
					store.Record(event)
				}
				client.CommitRecords(ctx, p.Records...)
			})
		}
	}
}

func startHTTPServer(store stream.StatsStore) {
	http.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		snapshot := store.GetSnapshot()
		if err := json.NewEncoder(w).Encode(snapshot); err != nil {
			http.Error(w, "failed to encode stats", http.StatusInternalServerError)
		}
	})

	port := ":8080"
	log.Printf("HTTP server listening on %s", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("HTTP server failed: %v", err)
	}
}

func handleShutdown(cancel context.CancelFunc) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
	cancel()
}
