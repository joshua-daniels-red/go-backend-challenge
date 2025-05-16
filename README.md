# Go Backend Challenge

This project is a multi-stage backend system written in idiomatic Go, focused on real-time stream processing using the [Wikipedia Recent Changes EventStream](https://stream.wikimedia.org/v2/stream/recentchange). Each chapter is self-contained, independently runnable, and builds on the previous to demonstrate increasing complexity.

The application is built entirely using the Go standard library, with clean modular design and production-ready practices.

## 📁 Project Structure

```
go-backend-challenge/
├── ch-1/  # Basic HTTP server with /status
├── ch-2/  # Real-time stream logger from Wikipedia
├── ch-3/  # In-memory stats API + tests
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

### `ch-3`: Real-time Stats + API
- Tracks and exposes the following stats via `/stats`:
  - Total messages consumed
  - Number of distinct users
  - Number of bot vs non-bot users
  - Count of events per distinct `server_url`
- Implements thread-safe in-memory state tracking using `sync.RWMutex`
- Fully unit-tested with `go test` and `-race` flag for concurrency safety

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

## 🧪 Testing (`ch-3`)

From inside the `ch-3` directory:

### Run Unit Tests

```bash
go test ./...
```

### Run Tests with Race Detector

```bash
go test -race ./...
```

This verifies the thread safety of concurrent writes and reads in the in-memory metrics system.

## 📦 Requirements

- Go 1.21+
- No third-party dependencies
- Works on macOS, Linux, and Windows

## 🛠️ Author

Joshua Daniels  
Senior Software Engineer  
