package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gocql/gocql"
	"github.com/joshua-daniels-red/go-backend-challenge/ch-5/internal/config"
	"github.com/joshua-daniels-red/go-backend-challenge/ch-5/internal/stream"
	"github.com/twmb/franz-go/pkg/kgo"
)

const topic = "wikipedia.changes"

var (
	configLoadFunc           = config.Load
	newKafkaClientFunc       = kgo.NewClient
	newCassandraSessionFn    = defaultCassandraSessionFn
	streamWikipediaHandlerFn = http.ListenAndServe
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("consumer error: %v", err)
	}
}

func run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	sigCh := make(chan os.Signal, 1)
	go handleShutdown(cancel, sigCh)

	cfg, err := configLoadFunc()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	client, err := newKafkaClientFunc(
		kgo.SeedBrokers(cfg.RedpandaBroker),
		kgo.ConsumeTopics(topic),
		kgo.ConsumerGroup("wikipedia-consumer-group"),
		kgo.MaxConcurrentFetches(5),
	)
	if err != nil {
		return fmt.Errorf("failed to create Kafka client: %w", err)
	}
	defer client.Close()

	var store stream.StatsStore
	if cfg.Storage == "cassandra" {
		sess, err := newCassandraSessionFn()
		if err != nil {
			return fmt.Errorf("failed to connect to Cassandra: %w", err)
		}
		defer sess.Close()
		store = stream.NewCassandraStats(stream.NewCassandraSessionAdapter(sess))
	} else {
		store = stream.NewInMemoryStats()
	}

	go func() {
		http.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(store.GetSnapshot()); err != nil {
				http.Error(w, "failed to encode stats", http.StatusInternalServerError)
			}
		})
		log.Println("HTTP server listening on :8080")
		if err := streamWikipediaHandlerFn(":8080", nil); err != nil {
			log.Printf("HTTP server error: %v", err)
		}
	}()

	log.Println("Consumer started. Waiting for messages...")
	for {
		select {
		case <-ctx.Done():
			log.Println("Shutting down consumer...")
			return nil
		default:
			fetches := client.PollFetches(ctx)
			fetches.EachPartition(func(p kgo.FetchTopicPartition) {
				log.Printf("Fetched %d records from partition %s", len(p.Records), p.Topic)
				for _, record := range p.Records {
					var e stream.Event
					if err := json.Unmarshal(record.Value, &e); err != nil {
						log.Printf("invalid record: %v", err)
						continue
					}
					store.Record(e)
				}
				client.CommitRecords(ctx, p.Records...)
			})
		}
	}
}

func defaultCassandraSessionFn() (*gocql.Session, error) {
	cluster := gocql.NewCluster("cassandra")
	cluster.Keyspace = "goanalytics"
	cluster.Consistency = gocql.Quorum
	cluster.ConnectTimeout = 5 * time.Second
	return cluster.CreateSession()
}

func handleShutdown(cancel context.CancelFunc, sigCh chan os.Signal) {
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
	cancel()
}
