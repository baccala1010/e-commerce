FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install git for go mod download
RUN apk add --no-cache git

# Copy go mod and sum files
COPY api-gateway/go.mod api-gateway/go.sum ./

# Copy local modules for monorepo build
COPY inventory ./inventory
COPY order ./order

# Download dependencies
RUN go mod download

# Copy the source code
COPY api-gateway/. .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o api-gateway ./cmd/api-gateway

# Use a smaller image for the final build
FROM alpine:latest

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/api-gateway .
COPY --from=builder /app/config ./config

# Set the Docker-specific config file
ENV CONFIG_PATH=/app/config/config.docker.yaml

# Expose the port
EXPOSE 8080

# Run the service
CMD ["./api-gateway"]
