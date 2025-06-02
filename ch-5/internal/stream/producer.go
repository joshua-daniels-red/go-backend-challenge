package stream

import (
	"bufio"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"
)

const topic = "wikipedia.changes"

// streamProducer defines the minimal Kafka interface needed for testing
type streamProducer interface {
	Produce(ctx context.Context, record *kgo.Record, cb func(*kgo.Record, error))
	Close()
}

// test hook override
var kafkaClientOverride streamProducer

// SetKafkaClientForTest allows test code to override the Kafka client
func SetKafkaClientForTest(p streamProducer) {
	kafkaClientOverride = p
}

func StreamWikipediaEvents(ctx context.Context, broker string, wikipediaURL string) error {
	var client streamProducer
	var err error

	if kafkaClientOverride != nil {
		client = kafkaClientOverride
	} else {
		client, err = kgo.NewClient(
			kgo.SeedBrokers(broker),
			kgo.ProducerLinger(100*time.Millisecond),
		)
		if err != nil {
			return err
		}
	}
	defer client.Close()

	resp, err := http.Get(wikipediaURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	log.Println("Connected to Wikipedia stream. Streaming events...")

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			log.Println("Producer shutting down...")
			return nil
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

			event := Event{
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

	return nil
}
