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
│       ├── client.go            # Wikipedia stream client (stubbed/disabled in ch-4)
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
with in memory the password is admin
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
