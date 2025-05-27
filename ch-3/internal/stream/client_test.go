package stream

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockStatsStore struct {
	Recorded []ChangeEvent
}

func (m *mockStatsStore) Record(ev ChangeEvent) {
	m.Recorded = append(m.Recorded, ev)
}

func (m *mockStatsStore) GetSnapshot() StatsSnapshot {
	return StatsSnapshot{} 
}

func TestWikipediaClient_Connect(t *testing.T) {
	fakeData := `data: {"user":"bob","bot":false,"server_url":"en.wikipedia.org"}
data: {"user":"alice","bot":true,"server_url":"fr.wikipedia.org"}`

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		for _, line := range strings.Split(fakeData, "\n") {
			_, _ = w.Write([]byte(line + "\n"))
		}
	}))
	defer ts.Close()

	mockStore := &mockStatsStore{}
	client := NewWikipediaClient(mockStore, ts.URL)

	err := client.Connect()
	assert.NoError(t, err)

	assert.Len(t, mockStore.Recorded, 2)

	assert.Equal(t, "bob", mockStore.Recorded[0].User)
	assert.Equal(t, false, mockStore.Recorded[0].Bot)
	assert.Equal(t, "en.wikipedia.org", mockStore.Recorded[0].ServerURL)

	assert.Equal(t, "alice", mockStore.Recorded[1].User)
	assert.Equal(t, true, mockStore.Recorded[1].Bot)
	assert.Equal(t, "fr.wikipedia.org", mockStore.Recorded[1].ServerURL)
}

func TestWikipediaClient_Connect_ErrorPaths(t *testing.T) {
	t.Run("HTTP error from streamURL", func(t *testing.T) {
		client := NewWikipediaClient(&mockStatsStore{}, "http://invalid.invalid")
		err := client.Connect()
		assert.Error(t, err) 
	})

	t.Run("Unmarshal error for invalid JSON", func(t *testing.T) {
		fakeData := `data: {invalid json`

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/event-stream")
			_, _ = w.Write([]byte(fakeData + "\n"))
		}))
		defer ts.Close()

		mockStore := &mockStatsStore{}
		client := NewWikipediaClient(mockStore, ts.URL)

		err := client.Connect()
		assert.NoError(t, err) 
		assert.Len(t, mockStore.Recorded, 0)
	})
}

func TestWikipediaClient_Connect_IgnoresNonDataLines(t *testing.T) {
	fakeData := `event: ping`

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		_, _ = w.Write([]byte(fakeData + "\n"))
	}))
	defer ts.Close()

	mockStore := &mockStatsStore{}
	client := NewWikipediaClient(mockStore, ts.URL)

	err := client.Connect()
	assert.NoError(t, err)

	assert.Len(t, mockStore.Recorded, 0)
}
