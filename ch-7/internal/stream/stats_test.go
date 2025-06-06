package stream_test

import (
	"sync"
	"testing"

	"github.com/joshua-daniels-red/go-backend-challenge/ch-7/internal/stream"
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

func TestInMemoryStats_RecordMany(t *testing.T) {
	store := stream.NewInMemoryStats()

	events := []stream.Event{
		{Domain: "en.wikipedia.org", User: "alice"},
		{Domain: "en.wikipedia.org", User: "bob"},
		{Domain: "de.wikipedia.org", User: "alice"},
	}

	store.RecordMany(events)

	snapshot := store.GetSnapshot()
	assert.Equal(t, 2, snapshot.ByUser["alice"])
	assert.Equal(t, 1, snapshot.ByUser["bob"])
	assert.Equal(t, 2, snapshot.ByDomain["en.wikipedia.org"])
	assert.Equal(t, 1, snapshot.ByDomain["de.wikipedia.org"])
}

func TestInMemoryStats_ConcurrentWrites(t *testing.T) {
	store := stream.NewInMemoryStats()
	events := []stream.Event{
		{Domain: "en.wikipedia.org", User: "alice"},
		{Domain: "en.wikipedia.org", User: "bob"},
		{Domain: "de.wikipedia.org", User: "charlie"},
	}

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for _, e := range events {
				store.Record(e)
			}
		}()
	}
	wg.Wait()

	snapshot := store.GetSnapshot()
	total := snapshot.ByUser["alice"] + snapshot.ByUser["bob"] + snapshot.ByUser["charlie"]
	assert.Equal(t, 10*len(events), total)
}

func TestInMemoryStats_ConcurrentRecordMany(t *testing.T) {
	store := stream.NewInMemoryStats()
	events := []stream.Event{
		{Domain: "en.wikipedia.org", User: "alice"},
		{Domain: "de.wikipedia.org", User: "bob"},
	}

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			store.RecordMany(events)
		}()
	}
	wg.Wait()

	snapshot := store.GetSnapshot()
	assert.Equal(t, 10*2, snapshot.ByUser["alice"]+snapshot.ByUser["bob"])
	assert.Equal(t, 10, snapshot.ByDomain["en.wikipedia.org"])
	assert.Equal(t, 10, snapshot.ByDomain["de.wikipedia.org"])
}
