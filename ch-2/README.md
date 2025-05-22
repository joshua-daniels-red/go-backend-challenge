# Chapter 2: Dockerized Go Backend

This folder (`ch-2/`) contains a Dockerized version of the Go backend from `ch-1`, structured to support multi-environment deployment and runtime configuration.

## ✅ Features
- Identical functionality as `ch-1`
- Reads dynamic config from `config.json`
- Built using a multi-stage Docker build
- Supports minimal final image via `scratch`

## 📁 Structure
```
ch-2/
├── go.mod
├── config.json
├── Dockerfile
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── config/
│   │   └── loader.go
│   ├── server/
│   │   └── server.go
│   └── stream/
│       ├── client.go
│       ├── stats.go
│       └── types.go
```

## ⚙️ config.json
```json
{
  "port": "7000",
  "stream_url": "https://stream.wikimedia.org/v2/stream/recentchange"
}
```

## 🐳 Dockerfile
```Dockerfile
# --- Stage 1: Builder ---
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod tidy && go build -o server ./cmd/server

# --- Stage 2: Minimal image ---
FROM scratch
COPY --from=builder /app/server /server
COPY --from=builder /app/config.json /config.json
ENTRYPOINT ["/server"]
```

## 🚀 Build and Run
### Local build
```bash
cd ch-2
GOOS=linux GOARCH=amd64 go build -o server ./cmd/server
./server
```

### Docker build
```bash
# Build Docker image
docker build -t ch2-server .

# Run container on port 7000
docker run -p 7000:7000 ch2-server
```

### Test URLs
- [http://localhost:7000/status](http://localhost:7000/status)
- [http://localhost:7000/stats](http://localhost:7000/stats)

---

