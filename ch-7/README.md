# Chapter 7: Multi-Threaded Stream Processing with Batching

This chapter evolves the backend application into a high-throughput, non-blocking, multi-threaded system capable of consuming messages from Kafka (Redpanda), batching them, and persisting efficiently to Cassandra. It builds directly on the foundation laid in Chapter 6.

---

## ✅ Key Additions in This Chapter

### 🧵 Multi-Threaded Kafka Consumers

* The application now spawns multiple Kafka consumers concurrently within a single process.
* Each consumer runs in its own goroutine and shares access to a common stats store.
* Graceful shutdown is handled via context cancellation and sync.WaitGroup.

### 💾 Batch Writing to Storage

* Events are now buffered and flushed to the underlying store (Cassandra or in-memory) in batches.
* Batching is triggered either by reaching a threshold size or a flush interval timeout.
* This significantly improves write throughput and efficiency.

### 🔄 Commit After Flush

* Kafka offsets are only committed **after** a successful batch flush.
* This ensures that events are not lost if the application crashes or is interrupted mid-batch.

### 🧪 Concurrency-Safe Testing

* Additional tests were added to verify `RecordMany()` and concurrent writes to the store.
* All tests are run with the Go `-race` detector enabled to catch data races early.

### 🐳 CI/CD Enhancements

* The GitHub Actions pipeline now builds, tests, and verifies both `/metrics` and `/stats` endpoints.
* Integration steps spin up all services with Docker Compose, test functionality, and then tear them down.
* `.env` configuration is injected dynamically during CI to avoid leaking credentials or requiring committed env files.

---

## 🗂️ Project Structure (Additions)

```
ch-7/
├── cmd/
│   └── consumer/               # Multi-threaded Kafka consumer logic
├── internal/
│   └── stream/
│       ├── batcher.go          # Batch buffering and flushing logic
│       └── stats.go            # Updated with RecordMany interface
├── .env                        # Created dynamically in CI
├── docker-compose.yml
└── .github/workflows/ci-ch7.yml
```

---

## 📈 How to Run Locally

```bash
# Start the full stack
cd ch-7
cp .env.example .env  # or manually create .env with required vars

# Then
docker-compose up --build
```

Access:

* Stats endpoint: [http://localhost:8080/stats](http://localhost:8080/stats)
* Prometheus metrics: [http://localhost:8080/metrics](http://localhost:8080/metrics)
* Cassandra CLI: `docker exec -it ch_7_cassandra cqlsh`

---

## 🧪 To Run Tests

```bash
cd ch-7
go test -race ./...
```

---

## ✅ Summary

This chapter focused on production-readiness: concurrency, batching, fault tolerance, and CI/CD. These additions transform the application into a resilient, efficient analytics pipeline ready for real-world traffic.
