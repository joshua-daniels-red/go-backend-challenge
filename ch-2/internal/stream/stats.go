package stream

import (
	"encoding/json"
	"net/http"
	"sync"
)

type Stats struct {
	mu           sync.RWMutex
	Total        int
	Users        map[string]struct{}
	BotCount     int
	NonBotCount  int
	ServerCounts map[string]int
}

func NewStats() *Stats {
	return &Stats{
		Users:        make(map[string]struct{}),
		ServerCounts: make(map[string]int),
	}
}

func (s *Stats) Record(ev ChangeEvent) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Total++
	s.Users[ev.User] = struct{}{}
	if ev.Bot {
		s.BotCount++
	} else {
		s.NonBotCount++
	}
	s.ServerCounts[ev.ServerURL]++
}

func (s *Stats) Handler(w http.ResponseWriter, r *http.Request) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	resp := struct {
		Messages      int            `json:"messages"`
		DistinctUsers int            `json:"distinct_users"`
		Bots          int            `json:"bots"`
		NonBots       int            `json:"non_bots"`
		ByServer      map[string]int `json:"by_server_url"`
	}{
		Messages:      s.Total,
		DistinctUsers: len(s.Users),
		Bots:          s.BotCount,
		NonBots:       s.NonBotCount,
		ByServer:      s.ServerCounts,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
