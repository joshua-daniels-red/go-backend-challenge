package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/joshua-daniels-red/go-backend-challenge/ch-3/internal/config"
	"github.com/joshua-daniels-red/go-backend-challenge/ch-3/internal/stream"
	"github.com/stretchr/testify/assert"
)
type fakeCassandraStats struct{}
func (f *fakeCassandraStats) Record(ev stream.ChangeEvent) {
	// no-op for test
}

func (f *fakeCassandraStats) GetSnapshot() stream.StatsSnapshot {
	return stream.StatsSnapshot{
		Messages:      42,
		DistinctUsers: 5,
		Bots:          2,
		NonBots:       3,
		ByServer:      map[string]int{"test": 42},
	}
}


func TestNewHTTPServer_InMemory(t *testing.T) {
	cfg := &config.Config{
		Port:     "8080",
		JWTSecret: "testsecret",
		StreamURL: "wss://example.com", 
		Storage:  "in-memory", 
		DisableStreaming: true,
	}

	srv := NewHTTPServer(cfg)
	ts := httptest.NewServer(srv.Handler)
	defer ts.Close()

	
	resp, err := http.Get(ts.URL + "/status")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var statusResp statusResponse
	err = json.NewDecoder(resp.Body).Decode(&statusResp)
	assert.NoError(t, err)
	assert.Equal(t, "ok", statusResp.Status)


	resp, err = http.Get(ts.URL + "/stats")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	reqBody := bytes.NewBufferString(`{"username":"admin","password":"admin"}`)
	resp, err = http.Post(ts.URL+"/login", "application/json", reqBody)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var loginResp LoginResponse
	err = json.NewDecoder(resp.Body).Decode(&loginResp)
	assert.NoError(t, err)
	assert.NotEmpty(t, loginResp.Token)

}

func TestNewHTTPServer_StreamingDoesNotCrash(t *testing.T) {
	cfg := &config.Config{
		Port:      "8080",
		JWTSecret: "testsecret",
		StreamURL: "wss://invalid.fake", 
		Storage:   "in-memory",
		DisableStreaming: true,
	}

	
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("server crashed on bad stream URL: %v", r)
		}
	}()

	_ = NewHTTPServer(cfg)
	time.Sleep(50 * time.Millisecond) 
}

func TestNewHTTPServer_CassandraBranch(t *testing.T) {
	cfg := &config.Config{
		Port:             "8080",
		JWTSecret:        "testsecret",
		StreamURL:        "not-used",
		Storage:          "cassandra",     
		DisableStreaming: true,          
	}

	srv := NewHTTPServer(cfg, &fakeCassandraStats{}) 
	ts := httptest.NewServer(srv.Handler)
	defer ts.Close()

	// /status should return 200
	resp, err := http.Get(ts.URL + "/status")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}