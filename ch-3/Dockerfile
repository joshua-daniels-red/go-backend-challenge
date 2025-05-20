# Stage 1: Build static Go binary
FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY ./cmd ./cmd
COPY ./internal ./internal
COPY go.mod .

RUN CGO_ENABLED=0 go build -o server ./cmd/server

# Stage 2: Create scratch image with certs
FROM scratch

WORKDIR /app

# Copy binary and config
COPY --from=builder /app/server .
COPY config.json .

# ✅ Copy root certificates from Alpine image
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT ["./server"]
