package stream

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestStatsHandler(t *testing.T) {
	stats := NewStats()

	// Simulate some events
	stats.Record(ChangeEvent{User: "user1", Bot: false, ServerURL: "https://en.wikipedia.org"})
	stats.Record(ChangeEvent{User: "bot1", Bot: true, ServerURL: "https://commons.wikimedia.org"})
	stats.Record(ChangeEvent{User: "user1", Bot: false, ServerURL: "https://en.wikipedia.org"})

	// Simulate HTTP request
	req := httptest.NewRequest(http.MethodGet, "/stats", nil)
	rr := httptest.NewRecorder()

	stats.Handler(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", rr.Code)
	}

	var response struct {
		Messages      int            `json:"messages"`
		DistinctUsers int            `json:"distinct_users"`
		Bots          int            `json:"bots"`
		NonBots       int            `json:"non_bots"`
		ByServer      map[string]int `json:"by_server_url"`
	}

	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode JSON: %v", err)
	}

	if response.Messages != 3 {
		t.Errorf("Expected 3 messages, got %d", response.Messages)
	}
	if response.DistinctUsers != 2 {
		t.Errorf("Expected 2 distinct users, got %d", response.DistinctUsers)
	}
	if response.Bots != 1 || response.NonBots != 2 {
		t.Errorf("Expected 1 bot and 2 non-bots, got %d and %d", response.Bots, response.NonBots)
	}
	if response.ByServer["https://en.wikipedia.org"] != 2 {
		t.Errorf("Expected 2 messages from en.wikipedia.org")
	}
}
