package stream

import (
	"bufio"
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

const streamURL = "https://stream.wikimedia.org/v2/stream/recentchange"

type WikipediaClient struct {
	client *http.Client
	stats  *Stats
}

func NewWikipediaClient(stats *Stats) *WikipediaClient {
	return &WikipediaClient{
		client: &http.Client{},
		stats:  stats,
	}
}

func (wc *WikipediaClient) Connect() error {
	resp, err := wc.client.Get(streamURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) < 6 || !strings.HasPrefix(line, "data:") {
			continue
		}

		line = strings.TrimPrefix(line, "data:")
		var change ChangeEvent
		if err := json.Unmarshal([]byte(line), &change); err != nil {
			log.Printf("failed to parse event: %v", err)
			continue
		}

		wc.stats.Record(change)
	}

	return scanner.Err()
}
