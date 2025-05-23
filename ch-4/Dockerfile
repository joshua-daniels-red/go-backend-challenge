# --- Stage 1: Build binary ---
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod tidy && go build -o server ./cmd/server

# --- Stage 2: Lightweight runtime with certs ---
FROM alpine:latest
RUN apk add --no-cache ca-certificates
COPY --from=builder /app/server /server
COPY --from=builder /app/config.json /config.json
ENTRYPOINT ["/server"]
