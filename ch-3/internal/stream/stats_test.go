package stream

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRecordAndSnapshot(t *testing.T) {
	stats := NewStats()

	events := []ChangeEvent{
		{User: "alice", Bot: false, ServerURL: "https://en.wikipedia.org"},
		{User: "bob", Bot: true, ServerURL: "https://en.wikipedia.org"},
		{User: "alice", Bot: false, ServerURL: "https://commons.wikimedia.org"},
		{User: "charlie", Bot: false, ServerURL: "https://en.wikipedia.org"},
		{User: "dave", Bot: true, ServerURL: "https://de.wikipedia.org"},
	}

	for _, ev := range events {
		stats.Record(ev)
	}

	snapshot := stats.GetSnapshot()

	assert.Equal(t, 5, snapshot.Messages)
	assert.Equal(t, 4, snapshot.DistinctUsers) // alice, bob, charlie, dave
	assert.Equal(t, 2, snapshot.Bots)
	assert.Equal(t, 3, snapshot.NonBots)
	assert.Equal(t, 3, len(snapshot.ByServer))
	assert.Equal(t, 3, snapshot.ByServer["https://en.wikipedia.org"])
	assert.Equal(t, 1, snapshot.ByServer["https://commons.wikimedia.org"])
	assert.Equal(t, 1, snapshot.ByServer["https://de.wikipedia.org"])
}

func TestRecordAndGetSnapshot(t *testing.T) {
	stats := NewStats()

	events := []ChangeEvent{
		{User: "alice", Bot: false, ServerURL: "https://en.wikipedia.org"},
		{User: "bob", Bot: true, ServerURL: "https://en.wikipedia.org"},
		{User: "alice", Bot: false, ServerURL: "https://commons.wikimedia.org"},
		{User: "carol", Bot: false, ServerURL: "https://en.wikipedia.org"},
	}

	for _, ev := range events {
		stats.Record(ev)
	}

	snapshot := stats.GetSnapshot()

	assert.Equal(t, 4, snapshot.Messages)
	assert.Equal(t, 3, snapshot.DistinctUsers)
	assert.Equal(t, 1, snapshot.Bots)
	assert.Equal(t, 3, snapshot.NonBots)
	assert.Equal(t, map[string]int{
		"https://en.wikipedia.org":     3,
		"https://commons.wikimedia.org": 1,
	}, snapshot.ByServer)
}