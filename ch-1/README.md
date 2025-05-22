# Chapter 1: Go Backend Application

This directory (`ch-1/`) contains a self-contained Go application that:
- Listens on port **7000**
- Consumes the [Wikipedia recent change stream](https://stream.wikimedia.org/v2/stream/recentchange)
- Tracks custom analytics in-memory
- Exposes two endpoints: `/status` and `/stats`
- Implements graceful shutdown and unit tests

## ✅ Features

### `/status`
Returns a simple health check JSON:
```json
{ "status": "ok" }
```

### `/stats`
Returns real-time analytics computed from Wikipedia event stream:
```json
{
  "messages": 125,
  "distinct_users": 80,
  "bots": 34,
  "non_bots": 91,
  "by_server_url": {
    "https://en.wikipedia.org": 90,
    "https://commons.wikimedia.org": 35
  }
}
```

## 📁 File Structure

```
ch-1/
├── go.mod
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── server/
│   │   ├── server.go
│   │   └── server_test.go
│   └── stream/
│       ├── client.go
│       ├── stats.go
│       └── types.go
```

## 🚀 How to Run

```bash
cd ch-1
# Download dependencies
go mod tidy

# Run the server
go run ./cmd/server
```

Access the app:
- [http://localhost:7000/status](http://localhost:7000/status)
- [http://localhost:7000/stats](http://localhost:7000/stats)

## 🧪 Run Tests

```bash
# Unit tests
cd ch-1
go test ./...

# Race detector
go test -race ./...
```

## 🛠 Requirements
- Go 1.21+
- No external libraries

## ✨ Highlights
- Idiomatic Go with `cmd/` and `internal/` structure
- In-memory stat aggregation
- Concurrency-safe with `sync.RWMutex`
- Graceful shutdown using `os/signal` and `context`

---

**Author:** Joshua Daniels
