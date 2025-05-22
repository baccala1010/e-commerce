FROM golang:1.23 as builder

WORKDIR /app

# Install build dependencies
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    librdkafka-dev \
    && rm -rf /var/lib/apt/lists/*

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application for Linux/amd64 and statically link
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -a -o order-service ./cmd/order && ls -lh order-service

FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates libc6-compat

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/order-service /app/order-service
COPY --from=builder /app/config/config.docker.yaml /app/config/config.docker.yaml

# Set the Docker-specific config file
ENV CONFIG_PATH=/app/config/config.docker.yaml

EXPOSE 8082 8083 9082

CMD ["/app/order-service"]
