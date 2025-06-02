# Chapter 5: Event-Driven Architecture with Redpanda

This module introduces an event-driven architecture using [Redpanda](https://redpanda.com/) to build a resilient streaming pipeline. The application is split into a **producer** and a **consumer**, each containerized and deployed independently. Events from the Wikimedia stream are ingested and persisted through Redpanda for downstream analytics.

---

## ✅ Objectives

1. **Split the application** into:

   * `cmd/producer`: connects to the Wikimedia stream and pushes events to Redpanda.
   * `cmd/consumer`: reads events from Redpanda and stores aggregate stats.

2. **Ingest + Produce**: The producer reads the Wikipedia stream and produces structured events to the `wikipedia.changes` topic.

3. **Consume + Aggregate**: The consumer reads batches of messages from Redpanda, unmarshals them, and updates stats by domain and user.

4. **Graceful Shutdown**: Both services use `context.WithCancel()` and listen for `SIGINT` and `SIGTERM` signals to terminate cleanly.

5. **High-Throughput Streaming**:

   * The consumer is configured with `MaxConcurrentFetches=5` to handle batches concurrently.
   * Acknowledgements are handled via `client.CommitRecords(...)`.

6. **Complete Dockerization**: Everything runs fully within Docker using `docker-compose`.

---

## 🧱 Architecture

```
+-----------------+       +-------------+        +-------------------+
| Wikipedia Stream| --->  |  Producer   | -----> |    Redpanda Topic |
+-----------------+       +-------------+        +-------------------+
                                                  |
                                                  v
                                          +----------------+
                                          |   Consumer     |
                                          | Aggregates +   |
                                          | Exposes /stats |
                                          +----------------+
```

---

## 🚀 Getting Started

### 1. Build & Run with Docker Compose

```bash
cd ch-5
docker-compose up --build
```

This will launch:

* Redpanda broker
* The producer (`cmd/producer`)
* The consumer (`cmd/consumer`)
* HTTP server at [http://localhost:8080/stats](http://localhost:8080/stats)

### 2. Environment Variables

You can override default environment variables with a `.env` file in the `ch-5` root:

```
# .env
REDPANDA_BROKER=redpanda:9092
WIKIPEDIA_STREAM_URL=https://stream.wikimedia.org/v2/stream/recentchange
STORAGE=cassandra
```

These are already embedded into the Dockerfiles for local testing.

---

## 🔬 API

### `/stats` – Returns aggregated stats in JSON

```json
{
  "by_domain": {
    "en.wikipedia.org": 31,
    "de.wikipedia.org": 12
  },
  "by_user": {
    "Alice": 20,
    "Bob": 15
  }
}
```

---

## 🦖 Running Tests

```bash
cd ch-5
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

Each component has isolated unit tests. The `cmd/consumer` and `cmd/producer` packages are testable with mocked Kafka and Cassandra clients.

---

## 🐳 Docker Images

Both services are built with multi-stage Dockerfiles:

* `ch-5/Dockerfile.producer`
* `ch-5/Dockerfile.consumer`

They can be independently built and pushed via the CI/CD pipeline.

---

## 🧪 CI/CD (GitHub Actions)

The GitHub Actions workflow in `.github/workflows/ch5.yml`:

* Runs `go vet` and `golangci-lint`.
* Runs unit tests with coverage.
* Builds `producer` and `consumer` images.
* Optionally pushes to a container registry.

---

## 📊 Technologies

* Go 1.21+
* Redpanda (via franz-go client)
* Docker & Compose
* Cassandra or in-memory storage
* GitHub Actions (CI)

---

## 📁 Project Structure

```
ch-5/
├── cmd/
│   ├── consumer/
│   │   ├── main.go
│   │   └── main_test.go
│   └── producer/
│       ├── main.go
│       └── main_test.go
├── db/
│   └── init.cql
├── internal/
│   ├── config/
│   │   ├── loader.go
│   │   └── loader_test.go
│   └── stream/
│       ├── mocks/
│       │   └── cassandra_mocks.go
│       ├── cassandra.go
│       ├── cassandra_adapter.go
│       ├── cassandra_test.go
│       ├── producer.go
│       ├── producer_test.go
│       ├── stats.go
│       ├── stats_test.go
│       └── types.go
├── docker-compose.yml
├── Dockerfile.consumer
├── Dockerfile.producer
├── .env
├── .env.example
├── .dockerignore
├── go.mod
├── go.sum
└── README.md
```

---

## 📍 Summary

Chapter 5 transitions the architecture into a decoupled, event-driven model powered by Redpanda. The system is now resilient, testable, and ready for real-time analytics at scale.
