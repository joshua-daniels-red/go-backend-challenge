package stream_test

import (
	"testing"

	"github.com/joshua-daniels-red/go-backend-challenge/ch-6/internal/stream"
	"github.com/stretchr/testify/assert"
)

func TestNewInMemoryStats(t *testing.T) {
	store := stream.NewInMemoryStats()
	assert.NotNil(t, store)
}

func TestInMemoryStats_RecordSingleEvent(t *testing.T) {
	store := stream.NewInMemoryStats()

	event := stream.Event{
		Domain: "en.wikipedia.org",
		User:   "alice",
	}
	store.Record(event)

	snapshot := store.GetSnapshot()
	assert.Equal(t, 1, snapshot.ByDomain["en.wikipedia.org"])
	assert.Equal(t, 1, snapshot.ByUser["alice"])
}

func TestInMemoryStats_RecordMultipleEvents(t *testing.T) {
	store := stream.NewInMemoryStats()

	events := []stream.Event{
		{Domain: "en.wikipedia.org", User: "alice"},
		{Domain: "en.wikipedia.org", User: "bob"},
		{Domain: "de.wikipedia.org", User: "alice"},
	}

	for _, e := range events {
		store.Record(e)
	}

	snapshot := store.GetSnapshot()
	assert.Equal(t, 2, snapshot.ByUser["alice"])
	assert.Equal(t, 1, snapshot.ByUser["bob"])
	assert.Equal(t, 2, snapshot.ByDomain["en.wikipedia.org"])
	assert.Equal(t, 1, snapshot.ByDomain["de.wikipedia.org"])
}

func TestInMemoryStats_SnapshotIsCopy(t *testing.T) {
	store := stream.NewInMemoryStats()

	store.Record(stream.Event{Domain: "test.com", User: "user1"})
	snapshot := store.GetSnapshot()

	// Modify snapshot
	snapshot.ByDomain["test.com"] = 999
	snapshot.ByUser["user1"] = 999

	// Get new snapshot and ensure original data is unchanged
	newSnapshot := store.GetSnapshot()
	assert.Equal(t, 1, newSnapshot.ByDomain["test.com"])
	assert.Equal(t, 1, newSnapshot.ByUser["user1"])
}
