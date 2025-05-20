package stream

import (
	"bufio"
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type WikipediaClient struct {
	client *http.Client
	stats  *Stats
	streamURL string
}

func NewWikipediaClient(stats *Stats,streamURL string) *WikipediaClient {
	return &WikipediaClient{
		client: &http.Client{},
		stats:  stats,
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
