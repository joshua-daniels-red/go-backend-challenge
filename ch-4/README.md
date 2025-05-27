# Chapter 4: CI/CD Pipeline & Publishing

This chapter builds on previous modules by introducing a continuous integration and deployment (CI/CD) pipeline for the application. It ensures all changes are validated through automated testing and linting, and builds are published to a container registry like `ghcr.io` or DockerHub once verified.

---

## ✅ Features

* Full CI/CD pipeline using GitHub Actions
* Test automation for PRs to `main`
* Linting and code quality checks
* Docker image creation and publishing
* Support for local pipeline testing using [`act`](https://github.com/nektos/act)

---

---

## 📁 Project Structure

```
ch-4/
├── cmd/server/
│   ├── main.go                  # Entry point: sets up config and starts the server
│   └── main_test.go             # Tests for top-level application startup logic
│
├── db/
│   └── init.cql                 # CQL script to initialize Cassandra schema
│
├── internal/
│   ├── config/
│   │   ├── loader.go            # Loads configuration from env vars or JSON
│   │   └── loader_test.go       # Unit tests for config loader
│   │
│   ├── server/
│   │   ├── auth.go              # JWT auth: token creation and validation
│   │   ├── auth_test.go         # Tests for auth logic
│   │   ├── middleware.go        # Middleware to enforce JWT on protected routes
│   │   ├── middleware_test.go   # Tests for middleware
│   │   └── server.go            # Defines routes and HTTP handlers
│   │   └── server_test.go       # Tests for HTTP routes
│   │
│   └── stream/
│       ├── cassandra.go         # Cassandra-based stats store (implements StatsStore)
│       ├── client.go            # Wikipedia stream client 
│       ├── client_test.go       # Tests for stream client logic
│       ├── stats.go             # Stats logic and in-memory store
│       ├── stats_test.go        # Tests for stats logic (both memory and Cassandra)
│       ├── types.go             # Shared data structures (e.g., Event, Snapshot)
│       ├── user.go              # Dummy user login validation logic
│       └── user_test.go         # Tests for user login logic
│
├── Dockerfile                   # Builds Go app container
├── README.md                    # Chapter-specific documentation
├── config.json                  # Optional: static config file (used in some envs)
├── docker-compose.yml           # Spins up Go app, Cassandra, and DB initializer
├── go.mod / go.sum              # Module and dependency tracking

```

---

## ↻ CI/CD Workflow

### On Pull Request to `main`, the pipeline will:

1. **Run Unit Tests**

   * All unit tests within the project must pass.

2. **Run Integration Tests**

   * Test the application against the database (using `docker-compose`).

3. **Run `go vet`**

   * Ensure Go code adheres to good practices.

4. **Run `golangci-lint`**

   * Perform static code analysis. See: [https://github.com/golangci/golangci-lint](https://github.com/golangci/golangci-lint)

5. **Build Docker Image**

   * Create a new image as defined in the `Dockerfile`.

6. **Publish Image**

   * Push the image to GitHub Container Registry (`ghcr.io`) or DockerHub when all checks pass.

---

## 🐳 Docker

### Dockerfile

A multi-stage build that compiles the Go app and produces a minimal final image.

### docker-compose.yml

Used for integration testing with services like Cassandra.

---

## 🧪 Local Testing

To simulate the GitHub Actions pipeline locally, you can use [`act`](https://github.com/nektos/act):

```bash
# Run all GitHub Actions workflows locally
act pull_request #there are issues running locally, I verified this in the github workflow
```

Ensure Docker is running, and all services (e.g. Cassandra) are defined in your `docker-compose.yml`.

---

## 🧹 Linting

Install and run `golangci-lint` locally:

```bash
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

golangci-lint run ./...
```

---

## 🤪 Unit + Integration Testing

```bash
go test ./...       # Unit tests
go test -race ./... # Race condition detection
```

For integration tests (against DB), use:

```bash
docker-compose up -d
go test ./... 
```

---

## ✍️ Author

Joshua Daniels
Senior Software Engineer
