package server_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/joshua-daniels-red/go-backend-challenge/ch-1/internal/server"
)

func TestStatusEndpoint(t *testing.T) {
	srv := server.NewHTTPServer(":7000")
	req := httptest.NewRequest("GET", "/status", nil)
	w := httptest.NewRecorder()

	srv.Handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
}

func TestStatsEndpoint(t *testing.T) {
	srv := server.NewHTTPServer(":7000")
	req := httptest.NewRequest("GET", "/stats", nil)
	w := httptest.NewRecorder()

	srv.Handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
	if len(w.Body.Bytes()) == 0 {
		t.Fatal("expected non-empty stats response")
	}
}