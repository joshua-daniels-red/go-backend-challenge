package stream

import (
	"bufio"
	"encoding/json"
	"log"
	"net/http"
)

const streamURL = "https://stream.wikimedia.org/v2/stream/recentchange"

// WikipediaClient manages the connection to the stream
type WikipediaClient struct {
	client *http.Client
}

func NewWikipediaClient() *WikipediaClient {
	return &WikipediaClient{client: &http.Client{}}
}

func (wc *WikipediaClient) ConnectAndLog() error {
	resp, err := wc.client.Get(streamURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) < 6 || line[:5] != "data:" {
			continue
		}

		line = line[5:] // trim "data:"
		var change ChangeEvent
		if err := json.Unmarshal([]byte(line), &change); err != nil {
			log.Printf("failed to parse event: %v", err)
			continue
		}

		log.Printf("[%s] %s - %s", change.Type, change.User, change.Title)
	}

	return scanner.Err()
}
