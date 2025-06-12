package stream

import (
	"sync"
)

// Snapshot represents an aggregate view of stats
type StatsSnapshot struct {
	ByDomain map[string]int `json:"by_domain"`
	ByUser   map[string]int `json:"by_user"`
}

// StatsStore defines an interface for tracking and retrieving stats
type StatsStore interface {
	Record(event Event)
	RecordMany([]Event)
	GetSnapshot() StatsSnapshot
}

type InMemoryStats struct {
	mu       sync.RWMutex
	domainCt map[string]int
	userCt   map[string]int
}

func NewInMemoryStats() *InMemoryStats {
	return &InMemoryStats{
		domainCt: make(map[string]int),
		userCt:   make(map[string]int),
	}
}

func (s *InMemoryStats) Record(event Event) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.domainCt[event.Domain]++
	s.userCt[event.User]++
}

func (s *InMemoryStats) RecordMany(events []Event) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, event := range events {
		s.domainCt[event.Domain]++
		s.userCt[event.User]++
	}
}

func (s *InMemoryStats) GetSnapshot() StatsSnapshot {
	s.mu.RLock()
	defer s.mu.RUnlock()

	domainCopy := make(map[string]int, len(s.domainCt))
	for k, v := range s.domainCt {
		domainCopy[k] = v
	}

	userCopy := make(map[string]int, len(s.userCt))
	for k, v := range s.userCt {
		userCopy[k] = v
	}

	return StatsSnapshot{
		ByDomain: domainCopy,
		ByUser:   userCopy,
	}
}
