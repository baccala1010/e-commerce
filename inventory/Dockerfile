FROM golang:1.23 as builder

WORKDIR /app

# Install build dependencies
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    librdkafka-dev \
    && rm -rf /var/lib/apt/lists/*

COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Build the application for Linux/amd64 and statically link
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -a -o inventory-service ./cmd/inventory && ls -lh inventory-service

FROM alpine:latest

RUN apk --no-cache add ca-certificates libc6-compat
WORKDIR /app
COPY --from=builder /app/inventory-service /app/inventory-service
COPY --from=builder /app/config/config.docker.yaml /app/config/config.docker.yaml

EXPOSE 8081 9081
CMD ["/app/inventory-service"]
