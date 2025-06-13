# Go Backend Challenge

This repository is a progressive, hands-on backend engineering challenge using Go. Each chapter builds upon the previous to construct a multi-threaded stream processing service with pluggable storage and robust deployment workflows. This project emphasizes use of the Go standard library, clean architecture, containerization, and CI/CD best practices.

---

## 📚 Project Overview

The goal is to incrementally develop a backend service that:

* Listens for HTTP requests
* Consumes Wikimedia's RecentChange stream
* Processes and aggregates streaming data
* Exposes in-memory and persistent statistics endpoints
* Supports JWT-based authentication (Ch 3)
* Uses protobuf for structured messaging (Ch 6)
* Supports Prometheus + Grafana observability (Ch 6)
* Is containerized, concurrent, and CI/CD ready (Ch 4, 7)
* Deploys to Kubernetes with CI/CD verification (Ch 8)

---

## 📁 Repository Structure

```
go-backend-challenge/
├── ch-1/          # Basic HTTP server with /status endpoint on port 7000
├── ch-2/          # Dockerized version of ch-1 with externalized configuration
├── ch-3/          # Adds streaming, Cassandra DB, and JWT-based authentication
├── ch-4/          # CI/CD pipeline: test, lint, docker build, publish
├── ch-5/          # Producer/Consumer with Redpanda, full containerized workflow
├── ch-6/          # Protobuf serialization and Prometheus/Grafana observability
├── ch-7/          # Multi-threaded consumers with batched writes and race-safety
├── ch-8/          # Full Kubernetes deployment with CI/CD verification
├── go.mod
├── go.sum
└── README.md      # This file
```

---

## 🔧 Chapters Breakdown

### ✅ Chapter 1: Basic HTTP Server

* Go HTTP server listening on port `7000`
* Implements `/status` endpoint for health checks
* Laid out using idiomatic Go project structure (`cmd/`, `internal/`)

### ✅ Chapter 2: Dockerization

* Introduces `Dockerfile` and containerizes the application
* Moves configuration (e.g., port, stream URL) to external `.env` or config file
* Build using multi-stage and `scratch` base image
* Adds `docker-compose.yml` for running locally

### ✅ Chapter 3: Streaming, Cassandra & Auth

* Connects to the Wikimedia RecentChange stream
* Implements `/stats` endpoint with in-memory or Cassandra-backed storage
* Tracks:

  * Total messages consumed
  * Distinct users
  * Bots vs Non-bots
  * Count by `server_name`
* Introduces `StatsStore` interface to support pluggable storage
* Adds Cassandra implementation (`CassandraStats`) using `gocql`
* Adds `/login` endpoint to issue JWTs
* Secures `/stats` with Bearer token authentication
* Configurable via a centralized config file
* Integration tested with `docker-compose`

### ✅ Chapter 4: CI/CD Pipeline

* Adds GitHub Actions workflow to automate:

  * Running all unit and integration tests
  * `go vet` and `golangci-lint`
  * Building and pushing Docker image
* Local CI simulation tested with [`nektos/act`](https://github.com/nektos/act)
* Ensures standards compliance and publishing readiness

### ✅ Chapter 5: Event Streaming with Redpanda

* Splits application into:

  * A **producer** that streams data from Wikimedia → Redpanda
  * A **consumer** that reads from Redpanda and records statistics
* Uses the `franz-go` client to produce/consume messages
* Statistics are stored using pluggable backends (in-memory or Cassandra)
* Both services are containerized:

  * `Dockerfile.producer`
  * `Dockerfile.consumer`
* `docker-compose.yml` spins up:

  * Redpanda
  * Cassandra
  * Producer and Consumer services
* Tests and CI workflow added for `ch-5`, reusing existing practices
* Includes `.env.example` for local dev:

```
REDPANDA_BROKER=redpanda:9092
WIKIPEDIA_STREAM_URL=https://stream.wikimedia.org/v2/stream/recentchange
STORAGE=cassandra
```

To run:

```bash
docker compose up --build
```

Stats available at: `GET http://localhost:8080/stats`

### ✅ Chapter 6: Protobuf & Observability

* Introduces `protobuf` for efficient message serialization
* Adds `.proto` schema + generated Go code for type-safe event encoding
* Redpanda messages are now serialized using protobuf
* Adds Prometheus metrics collection:

  * Total events consumed
  * Events persisted, succeeded, failed
* Adds Grafana dashboard with preconfigured panels
* Metrics exposed at: `GET /metrics`
* CI updated to include `/metrics` verification

### ✅ Chapter 7: Multi-Threaded Consumers & Batching

* Consumer now supports multiple concurrent workers using goroutines
* Shared store access with race-safe logic (`sync.Mutex` and `-race` tested)
* Adds batching layer:

  * Events are buffered and flushed to Cassandra/in-memory storage
  * Flush occurs by size threshold or time interval
* Kafka offsets are committed **only after** successful batch flush
* Prevents message loss mid-batch
* Pipeline updated to dynamically create `.env`, run end-to-end services, and validate `/stats` + `/metrics`

### ✅ Chapter 8: Kubernetes Deployment & CI/CD

* Deploys all components (Cassandra, Redpanda, Producer, Consumer, Prometheus, Grafana) into Kubernetes using Minikube
* Fully automated `setup.sh` brings the cluster online with working endpoints and observability
* Services are accessible using `minikube service` or Minikube IP + NodePort
* Pre-provisioned Grafana dashboard available at `http://<minikube-ip>:30300`
* CI/CD pipeline uses **KinD** in GitHub Actions to:
  * Build Docker images
  * Spin up Kubernetes cluster
  * Apply manifests
  * Verify `/stats` and `/metrics` endpoints
* Teardown script (`teardown.sh`) deletes all resources and persistent volumes

---

## 📊 Statistics Endpoint (`/stats`)

* Total messages consumed
* Distinct users
* Bots vs Non-bots
* Count by distinct `server_name`

---

## 🔐 Authentication

* `POST /login` issues JWTs with configured secret
* `/stats` is secured via Bearer token
* Configurable secret key via env/config file

---

## 🥯 Testing

* Unit tests (`go test ./...`)
* Integration tests (with Cassandra via Docker Compose)
* Race detector enabled (`go test -race`)
* Mocked Cassandra interactions for fast unit test coverage

---

## 🐳 Docker

* Multi-stage Dockerfile with `scratch` base for production
* `docker-compose.yml` spins up app and Cassandra cluster
* All configs externalized for portability

---

## 📦 Tech Stack

* **Go** (standard library preferred)
* **Cassandra/Scylla** for persistent analytics
* **Docker** & **Docker Compose**
* **JWT** for auth
* **Protobuf** for message encoding
* **Redpanda** for streaming
* **Prometheus + Grafana** for observability
* **GitHub Actions** for CI/CD
* **Kubernetes + Minikube** for orchestration
* **KinD** for CI-based Kubernetes cluster testing

---

## 🛠️ Requirements

* Go 1.21+
* Docker + Docker Compose
* Minikube or KinD for local/CI Kubernetes deployment
* GitHub account for publishing images (optional)

---

## 👤 Author

Joshua Daniels
