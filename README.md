# go-backend-challenge

A modular Go project structured into chapters, each showcasing different backend concepts, best practices, and architecture patterns. The goal of this repository is to provide a learning playground and real-world examples for building production-grade backend services using Go.

## 🚀 Project Overview

This repository is organized into chapters, with each chapter focusing on a specific backend development concept or challenge. You can explore each chapter independently to learn, test, or build upon it.

### ✨ Features Across Chapters

* Clean project architecture using Go best practices
* Real-world backend patterns: HTTP servers, JWT auth, DB integrations, etc.
* Modular directory structure for scalability and clarity
* GitHub Actions CI for automated linting and builds

## 📁 Repository Structure

```
go-backend-challenge/
├── ch-1/               # Basic HTTP server with simple endpoints
├── ch-2/               # Middleware and routing enhancements
├── ch-3/               # JSON handling and request validation
├── ch-4/               # Stats service with pluggable storage backends (e.g., Cassandra)
├── ch-5/               # 
├── go.mod
├── go.sum
└── README.md           # This file
```

Each chapter contains its own `README.md` with setup and usage instructions specific to that module.

## 🛠️ Getting Started

To get started with any chapter:

1. Clone the repository:

   ```bash
   git clone https://github.com/your-username/go-backend-challenge.git
   cd go-backend-challenge
   ```

2. Navigate to the desired chapter:

   ```bash
   cd ch-4
   ```

3. Install dependencies:

   ```bash
   go mod tidy
   ```

4. Run the server:

   ```bash
   go run cmd/server/main.go
   ```

## ✅ Requirements

* Go 1.20+
* Docker (optional, for running services like Cassandra locally)
* Make (optional, for workflows)

## 🧪 CI/CD

The project uses GitHub Actions to:

* Run linters (`golangci-lint`)
* Validate builds
* Ensure consistent code formatting

## 📚 Learning Goals

This repository is structured to help developers:

* Understand idiomatic Go practices
* Learn how to write scalable and testable backend code
* Practice real-world concepts like authentication, metrics, error handling, and observability


---

Built with ❤️ to master Go backend development.
