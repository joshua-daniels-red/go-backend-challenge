package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gocql/gocql"
	"github.com/joshua-daniels-red/go-backend-challenge/ch-3/internal/config"
	"github.com/joshua-daniels-red/go-backend-challenge/ch-3/internal/stream"
)

type statusResponse struct {
	Status string `json:"status"`
}

func NewHTTPServer(cfg *config.Config, injectedStore ...stream.StatsStore) *http.Server {
	var store stream.StatsStore
	var cassSession *gocql.Session
	var err error

	if cfg.Storage == "cassandra" {
		if len(injectedStore) > 0 {
			store = injectedStore[0]
			log.Println("Using injected Cassandra store (test mode)")
		} else {
			store, err = stream.NewCassandraStats(cfg.CassandraHost)
			if err != nil {
				log.Fatalf("failed to init Cassandra store: %v", err)
			}

			cluster := gocql.NewCluster(cfg.CassandraHost)
			cluster.Keyspace = "goanalytics"
			cassSession, err = cluster.CreateSession()
			if err != nil {
				log.Fatalf("could not establish session for auth: %v", err)
			}

			log.Println("Using Cassandra storage")
		}
	} else {
		store = stream.NewStats()
		log.Println("Using in-memory storage")
	}

	if !cfg.DisableStreaming {
		client := stream.NewWikipediaClient(store, cfg.StreamURL)
		go func() {
			if err := client.Connect(); err != nil {
				log.Fatalf("streaming failed: %v", err)
			}
		}()
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(statusResponse{Status: "ok"}); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	})

	var userStore stream.UserStore
	if cassSession != nil {
		userStore = stream.NewUserStore(cassSession)
	} else {
		userStore = stream.NewInMemoryUserStore()
	}
	mux.HandleFunc("/login", LoginHandler(userStore, cfg.JWTSecret))


	mux.HandleFunc("/stats", AuthMiddleware(cfg.JWTSecret, func(w http.ResponseWriter, r *http.Request) {
		snapshot := store.GetSnapshot()
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(snapshot); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}))

	return &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: mux,
	}
}
