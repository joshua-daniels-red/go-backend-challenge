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
