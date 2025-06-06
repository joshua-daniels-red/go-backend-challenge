package stream

import (
	"bufio"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	pb "github.com/joshua-daniels-red/go-backend-challenge/ch-7/proto"
	"github.com/twmb/franz-go/pkg/kgo"
	"google.golang.org/protobuf/proto"
)

type streamProducer interface {
	Produce(ctx context.Context, record *kgo.Record, cb func(*kgo.Record, error))
	ProduceSync(ctx context.Context, record *kgo.Record) error
	Flush(ctx context.Context) error
	Close()
}

type realKafkaClient struct {
	*kgo.Client
}

func (r *realKafkaClient) Produce(ctx context.Context, record *kgo.Record, cb func(*kgo.Record, error)) {
	r.Client.Produce(ctx, record, cb)
}

func (r *realKafkaClient) ProduceSync(ctx context.Context, record *kgo.Record) error {
	return r.Client.ProduceSync(ctx, record).FirstErr()
}

func (r *realKafkaClient) Flush(ctx context.Context) error {
	return r.Client.Flush(ctx)
}

func (r *realKafkaClient) Close() {
	r.Client.Close()
}

var kafkaClientOverride streamProducer

func SetKafkaClientForTest(p streamProducer) {
	kafkaClientOverride = p
}

func StreamWikipediaEvents(ctx context.Context, broker string, wikipediaURL string, topic string) error {
	var client streamProducer
	var err error

	if kafkaClientOverride != nil {
		client = kafkaClientOverride
	} else {
		c, err := kgo.NewClient(
			kgo.SeedBrokers(broker),
			kgo.ProducerLinger(100*time.Millisecond),
		)
		if err != nil {
			return err
		}
		client = &realKafkaClient{Client: c}
	}
	defer client.Close()

	resp, err := http.Get(wikipediaURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	log.Println("🌐 Connected to Wikipedia stream. Streaming events...")
	scanner := bufio.NewScanner(resp.Body)

	for scanner.Scan() {
		select {
		case <-ctx.Done():
			log.Println("🛑 Context cancelled, flushing pending messages before shutdown...")
			client.Flush(context.Background()) // flush with fresh context
			log.Println("✅ Producer flushed and shutting down.")
			return nil
		default:
			line := scanner.Text()

			if !strings.HasPrefix(line, "data: ") {
				continue
			}

			var raw map[string]interface{}
			if err := json.Unmarshal([]byte(line[6:]), &raw); err != nil {
				log.Printf("❌ Skipping malformed JSON event: %v", err)
				continue
			}

			metaMap, ok := raw["meta"].(map[string]interface{})
			if !ok {
				log.Println("❌ Skipping event: missing 'meta' field")
				continue
			}

			domain, _ := metaMap["domain"].(string)
			title, _ := raw["title"].(string)
			user, _ := raw["user"].(string)

			if domain == "" || title == "" || user == "" {
				pretty, _ := json.MarshalIndent(raw, "", "  ")
				log.Printf("❌ Skipping event - missing field(s): domain='%s', title='%s', user='%s'\n⚠️ Full event:\n%s",
					domain, title, user, pretty)
				continue
			}

			event := Event{
				Domain: domain,
				Title:  title,
				User:   user,
			}

			protoEvent := &pb.Event{
				Domain: event.Domain,
				Title:  event.Title,
				User:   event.User,
			}

			data, err := proto.Marshal(protoEvent)
			if err != nil {
				log.Printf("❌ Failed to marshal protobuf: %v", err)
				continue
			}

			record := &kgo.Record{
				Topic: topic,
				Value: data,
			}

			log.Printf("📦 Sending event to topic %s: %+v", topic, protoEvent)
			err = client.ProduceSync(ctx, record)
			if err != nil {
				log.Printf("❌ Failed to produce message: %v", err)
			} else {
				log.Printf("✅ Produced message to topic %s", record.Topic)
			}

		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("error reading stream: %v", err)
	}

	return nil
}
