package stream_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/joshua-daniels-red/go-backend-challenge/ch-6/internal/stream"
	"github.com/stretchr/testify/assert"
	"github.com/twmb/franz-go/pkg/kgo"
)

type mockProducer struct {
	produced []*kgo.Record
	lock     sync.Mutex
}

func (m *mockProducer) Produce(_ context.Context, record *kgo.Record, _ func(*kgo.Record, error)) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.produced = append(m.produced, record)
}

func (m *mockProducer) Close() {}

func (m *mockProducer) Flush(_ context.Context) error {
	return nil
}

func (m *mockProducerWithError) Flush(_ context.Context) error {
	return nil
}

func (m *mockProducer) ProduceSync(_ context.Context, record *kgo.Record) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.produced = append(m.produced, record)
	return nil
}

func (m *mockProducerWithError) ProduceSync(_ context.Context, _ *kgo.Record) error {
	return errors.New("mock error")
}

type mockProducerWithError struct{}

func (m *mockProducerWithError) Produce(_ context.Context, record *kgo.Record, cb func(*kgo.Record, error)) {
	cb(record, errors.New("mock error"))
}

func (m *mockProducerWithError) Close() {}

func TestStreamWikipediaEvents_Success(t *testing.T) {
	event := map[string]interface{}{
		"title": "Test Page",
		"user":  "TestUser",
		"meta": map[string]interface{}{
			"domain": "en.wikipedia.org",
		},
	}
	eventJSON, _ := json.Marshal(event)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "data: %s\n", eventJSON)
	}))
	defer ts.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	mock := &mockProducer{}
	stream.SetKafkaClientForTest(mock)
	err := stream.StreamWikipediaEvents(ctx, "broker", ts.URL, "test.topic")
	assert.NoError(t, err)
	assert.Len(t, mock.produced, 1)
}

func TestStreamWikipediaEvents_ProduceError(t *testing.T) {
	event := map[string]interface{}{
		"title": "Test Page",
		"user":  "TestUser",
		"meta": map[string]interface{}{
			"domain": "en.wikipedia.org",
		},
	}
	eventJSON, _ := json.Marshal(event)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "data: %s\n", eventJSON)
	}))
	defer ts.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	stream.SetKafkaClientForTest(&mockProducerWithError{})
	err := stream.StreamWikipediaEvents(ctx, "broker", ts.URL, "test.topic")
	assert.NoError(t, err)
}

func TestStreamWikipediaEvents_MalformedJSON(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "data: {bad json}\n")
	}))
	defer ts.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	stream.SetKafkaClientForTest(&mockProducer{})
	err := stream.StreamWikipediaEvents(ctx, "fake", ts.URL, "test.topic")
	assert.NoError(t, err)
}

func TestStreamWikipediaEvents_BadConnection(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	err := stream.StreamWikipediaEvents(ctx, "fake", "http://127.0.0.1:0", "test.topic")
	assert.Error(t, err)
}

func TestStreamWikipediaEvents_GracefulShutdown(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		fmt.Fprint(w, "data: {}\n")
	}))
	defer ts.Close()

	ctx, cancel := context.WithCancel(context.Background())
	stream.SetKafkaClientForTest(&mockProducer{})
	cancel()

	err := stream.StreamWikipediaEvents(ctx, "unused", ts.URL, "test.topic")
	assert.NoError(t, err)
}

func TestStreamWikipediaEvents_ShutdownMidScan(t *testing.T) {
	event := map[string]interface{}{
		"title": "T", "user": "U", "meta": map[string]interface{}{"domain": "D"},
	}
	data, _ := json.Marshal(event)

	ctx, cancel := context.WithCancel(context.Background())
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "data: %s\n", data)
		time.AfterFunc(50*time.Millisecond, cancel)
	}))
	defer ts.Close()

	stream.SetKafkaClientForTest(&mockProducer{})
	err := stream.StreamWikipediaEvents(ctx, "unused", ts.URL, "test.topic")
	assert.NoError(t, err)
}

func TestStreamWikipediaEvents_MissingMeta(t *testing.T) {
	payload := map[string]interface{}{
		"title": "T", "user": "U",
	}
	data, _ := json.Marshal(payload)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "data: %s\n", data)
	}))
	defer ts.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	stream.SetKafkaClientForTest(&mockProducer{})
	_ = stream.StreamWikipediaEvents(ctx, "unused", ts.URL, "test.topic")
}

func TestStreamWikipediaEvents_IncompleteFields(t *testing.T) {
	payload := map[string]interface{}{
		"title": "", "user": "", "meta": map[string]interface{}{"domain": ""},
	}
	data, _ := json.Marshal(payload)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "data: %s\n", data)
	}))
	defer ts.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	stream.SetKafkaClientForTest(&mockProducer{})
	_ = stream.StreamWikipediaEvents(ctx, "unused", ts.URL, "test.topic")
}

type errReader struct{}

func (e *errReader) Read(p []byte) (int, error) {
	return 0, fmt.Errorf("read error")
}

func (e *errReader) Close() error {
	return nil
}

func TestStreamWikipediaEvents_ScannerError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.(http.Flusher).Flush()
	}))
	defer ts.Close()

	http.DefaultClient = &http.Client{
		Transport: roundTripperFunc(func(req *http.Request) (*http.Response, error) {
			resp, err := http.DefaultTransport.RoundTrip(req)
			if err != nil {
				return nil, err
			}
			resp.Body = &errReader{}
			return resp, nil
		}),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	stream.SetKafkaClientForTest(&mockProducer{})
	err := stream.StreamWikipediaEvents(ctx, "unused", ts.URL, "test.topic")
	assert.NoError(t, err)
}

type roundTripperFunc func(req *http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}
