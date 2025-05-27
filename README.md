# Go Backend Challenge

This repository is a progressive, hands-on backend engineering challenge using Go. Each chapter builds upon the previous to construct a multi-threaded stream processing service with pluggable storage and robust deployment workflows. This project emphasizes use of the Go standard library, clean architecture, containerization, and CI/CD best practices.

---

## 📚 Project Overview

The goal is to incrementally develop a backend service that:

* Listens for HTTP requests
* Consumes Wikimedia's RecentChange stream
* Processes and aggregates streaming data
* Exposes in-memory and persistent statistics endpoints
* Supports JWT-based authentication
* Is containerized and CI/CD ready

---

## 📁 Repository Structure

```
go-backend-challenge/
├── ch-1/          # Basic HTTP server with /status endpoint on port 7000
├── ch-2/          # Dockerized version of ch-1 with externalized configuration
├── ch-3/          # Adds streaming, Cassandra DB, and JWT-based authentication
├── ch-4/          # CI/CD pipeline: test, lint, docker build, publish
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

## 🧪 Testing

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
* **GitHub Actions** for CI/CD

---

## 🛠 Requirements

* Go 1.21+
* Docker + Docker Compose
* GitHub account for publishing images (optional)

---

## 👤 Author
Joshua Daniels 
