# Go Backend Challenge

This project is a multi-stage backend system written in idiomatic Go, focused on real-time stream processing using the [Wikipedia Recent Changes EventStream](https://stream.wikimedia.org/v2/stream/recentchange). Each chapter is self-contained, independently runnable, and builds on the previous to demonstrate increasing complexity.

The application is built entirely using the Go standard library, with clean modular design and production-ready practices.

## 📁 Project Structure

```
go-backend-challenge/
├── ch-1/  # Basic HTTP server with /status
├── ch-2/  # Real-time stream logger from Wikipedia
├── ch-3/  # In-memory stats API + config loader, tests, Dockerfile
```

## ✅ Features by Chapter

### `ch-1`: Basic HTTP Server
- Listens on port `7000`
- Exposes `/status` endpoint returning `{ "status": "ok" }`
- Implements graceful shutdown via `context`
- Organized into `cmd/` and `internal/` packages

### `ch-2`: Wikipedia Stream Logger
- Connects to the Wikimedia EventStream API
- Consumes newline-delimited streaming JSON
- Logs real-time edit events to `stdout`

### `ch-3`: Real-time Stats + API + Dockerized Runtime
- Tracks and exposes the following stats via `/stats`:
  - Total messages consumed
  - Number of distinct users
  - Number of bot vs non-bot users
  - Count of events per distinct `server_url`
- Loads port and stream URL from a `config.json` file using a custom config loader
- Implements thread-safe in-memory state tracking using `sync.RWMutex`
- Fully unit-tested with `go test` and `-race` flag for concurrency safety
- Includes Docker support:
  - Multi-stage image based on Alpine
  - Minimal `scratch` image with root certificates for TLS

## 🚀 Running Each Chapter

### 🧩 ch-1: Basic HTTP Server

```bash
cd ch-1
go mod tidy
go run ./cmd/server
```

Visit: http://localhost:7000/status

### 🧩 ch-2: Wikipedia Stream Logger

```bash
cd ch-2
go mod tidy
go run ./cmd/streamer
```

Observe stdout logs for real-time Wikipedia edits.

### 🧩 ch-3: Stats Endpoint + Analytics

```bash
cd ch-3
go mod tidy
go run ./cmd/server
```

Visit: http://localhost:7000/stats  
Returns JSON like:

```json
{
  "messages": 42,
  "distinct_users": 35,
  "bots": 18,
  "non_bots": 24,
  "by_server_url": {
    "https://en.wikipedia.org": 30,
    "https://commons.wikimedia.org": 12
  }
}
```

## ⚙️ Configuration File (ch-3)

A `config.json` file is used to provide runtime configuration:

```json
{
  "port": ":7000",
  "stream_url": "https://stream.wikimedia.org/v2/stream/recentchange"
}
```

Place it in the root of the `ch-3` directory. The application loads this file on startup.

## 🐳 Docker Support (ch-3)

### Build and Run (Alpine-based)

```bash
cd ch-3
docker build -t ch3-server .
docker run -p 7000:7000 ch3-server
```

### Build and Run (Scratch image)

```bash
docker build -t ch3-scratch -f Dockerfile .
docker run -p 7000:7000 ch3-scratch
```

The `scratch` image includes only the Go binary, `config.json`, and CA certificates to allow secure HTTPS streaming.

## 🧪 Testing (`ch-3`)

From inside the `ch-3` directory:

```bash
go test ./...
go test -race ./...
```

This verifies the thread safety of concurrent writes and reads in the in-memory metrics system.

## 📦 Requirements

- Go 1.21+
- Docker (for containerized builds)
- No third-party dependencies
- Works on macOS, Linux, and Windows

## 🛠️ Author

Joshua Daniels  
Senior Software Engineer
