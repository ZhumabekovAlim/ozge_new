# syntax=docker/dockerfile:1

# Build stage
FROM golang:1.24 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/web

# Final stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/server /app/server
COPY config ./config
COPY migrations ./migrations
COPY static ./static
COPY DejaVuSans.ttf ./
EXPOSE 4000
ENTRYPOINT ["/app/server"]
