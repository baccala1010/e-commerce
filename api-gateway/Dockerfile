FROM golang:1.20-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

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