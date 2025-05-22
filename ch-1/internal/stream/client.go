package stream

import (
	"bufio"
	"encoding/json"
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