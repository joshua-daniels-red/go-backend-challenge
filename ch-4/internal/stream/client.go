package stream

import (
	"bufio"
	"encoding/json"
	"net/http"
	"strings"
)

// WikipediaClient streams events from the Wikimedia API and pushes them to a StatsStore.
type WikipediaClient struct {
	client    *http.Client
	stats     StatsStore
	streamURL string
}

// NewWikipediaClient accepts any StatsStore (in-memory or Cassandra) and a stream URL.
func NewWikipediaClient(stats StatsStore, streamURL string) *WikipediaClient {
	return &WikipediaClient{
		client:    &http.Client{},
		stats:     stats,
		streamURL: streamURL,
	}
}

func (wc *WikipediaClient) Connect() error {
	resp, err := wc.client.Get(wc.streamURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "data:") {
			continue
		}

		line = strings.TrimPrefix(line, "data:")
		var ev ChangeEvent
		if err := json.Unmarshal([]byte(line), &ev); err != nil {
			continue
		}
		wc.stats.Record(ev)
	}

	return scanner.Err()
}
