package server

import (
	"encoding/json"
	"net/http"
)

type statusResponse struct {
	Status string `json:"status"`
}

func NewHTTPServer(addr string) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/status", statusHandler)

	return &http.Server{
		Addr:    addr,
		Handler: mux,
	}
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(statusResponse{Status: "ok"})
}
