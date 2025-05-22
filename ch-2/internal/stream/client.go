package stream

import (
	"bufio"
	"encoding/json"
	"net/http"
	"strings"
)

type WikipediaClient struct {
	client    *http.Client
	stats     *Stats
	streamURL string
}

func NewWikipediaClient(stats *Stats, streamURL string) *WikipediaClient {
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