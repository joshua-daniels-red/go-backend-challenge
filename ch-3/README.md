# Chapter 3: Persistent Stats Backend with Auth

This module builds upon `ch-2` by adding persistent storage using Cassandra and integrating basic authentication with JWT. It demonstrates an end-to-end backend with:

* Real-time processing of Wikipedia change streams
* Stats aggregation
* Optional storage in memory or Cassandra
* Auth-protected `/stats` endpoint

---

## ✅ Features

* Full Docker Compose stack
* Switchable backend storage (in-memory or Cassandra)
* Basic login via `/login` with JWT token issuance
* Bearer token middleware for protected endpoints
* `init.cql` script to auto-provision Cassandra schema
* Clean graceful shutdown

---

## 📁 Project Structure

```
ch-3/
├── cmd/
│   └── server/
│       └── main.go              # Entry point
├── db/
│   └── init.cql                 # Cassandra schema + seed data
├── internal/
│   ├── config/
│   │   └── loader.go            # Loads JSON config
│   ├── server/
│   │   ├── auth.go              # Login logic
│   │   ├── middleware.go        # JWT validation
│   │   ├── server.go            # HTTP handlers
│   └── stream/
│       ├── cassandra.go        # Cassandra implementation of StatsStore
│       ├── client.go           # Wikipedia stream consumer
│       ├── stats.go            # In-memory StatsStore
│       ├── types.go            # Common types
│       └── user.go             # UserStore interface + Cassandra implementation
├── config.json                 # Runtime config
├── docker-compose.yml          # App + Cassandra orchestration
├── Dockerfile                  # Multi-stage build
├── go.mod / go.sum             # Go dependencies
└── README.md                   # This file
```

---

## ⚙️ Config: `config.json`

```json
{
  "port": "7000",
  "stream_url": "https://stream.wikimedia.org/v2/stream/recentchange",
  "storage": "cassandra",
  "cassandra_host": "cassandra-db",
  "jwt_secret": "secret123"
}
```

Set `"storage": "in-memory"` to switch off database usage.

---

## 🐳 Docker

### Dockerfile

Multi-stage build with a final `scratch` image for minimal footprint.

### docker-compose.yml

```yaml
version: '3.8'

services:
  cassandra-db:
    image: cassandra:4.1
    container_name: cassandra-db
    ports:
      - "9042:9042"
    volumes:
      - ./db/init.cql:/docker-entrypoint-initdb.d/init.cql:ro
    healthcheck:
      test: ["CMD-SHELL", "cqlsh -e 'describe keyspaces' || exit 1"]
      interval: 10s
      timeout: 5s
      retries: 10

  app:
    build: .
    depends_on:
      cassandra-db:
        condition: service_healthy
    ports:
      - "7000:7000"
    volumes:
      - ./config.json:/config.json
    environment:
      - CASSANDRA_HOST=cassandra-db
```

---

## 🚀 Usage

### Local run (in-memory)

```bash
cd ch-3
CONFIG_FILE=config.json go run ./cmd/server
```

### Docker run

```bash
# Build and start app + Cassandra
docker-compose up --build
```

---

## 🔐 Auth Endpoints

### POST `/login`

```bash
curl -X POST http://localhost:7000/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "password123"}'
```

Returns:

```json
{"token": "<JWT token>"}
```

### GET `/stats` (Requires Bearer token)

```bash
curl -H "Authorization: Bearer <token>" http://localhost:7000/stats
```

---

## 🧪 Testing

### Unit tests

```bash
cd ch-3
go test ./...
```

### Race condition testing

```bash
go test -race ./...
```

---

## 🛠️ Seeded Cassandra Schema: `init.cql`

```sql
CREATE KEYSPACE IF NOT EXISTS goanalytics WITH REPLICATION = {
  'class': 'SimpleStrategy', 'replication_factor': 1
};
USE goanalytics;

CREATE TABLE IF NOT EXISTS stats_summary (
  id TEXT PRIMARY KEY,
  total_messages COUNTER,
  bot_count COUNTER,
  non_bot_count COUNTER
);

CREATE TABLE IF NOT EXISTS unique_users (
  username TEXT PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS server_counts (
  server_url TEXT PRIMARY KEY,
  count COUNTER
);

CREATE TABLE IF NOT EXISTS users (
  username TEXT PRIMARY KEY,
  password TEXT
);

INSERT INTO users (username, password) VALUES ('admin', 'password123');
```

---

## ✍️ Author

Joshua Daniels
Senior Software Engineer
