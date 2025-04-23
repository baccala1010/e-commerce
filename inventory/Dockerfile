FROM golang:1.20-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o inventory-service ./cmd/inventory

# Use a smaller image for the final build
FROM alpine:latest

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/inventory-service .
COPY --from=builder /app/config ./config

# Set the Docker-specific config file
ENV CONFIG_PATH=/app/config/config.docker.yaml

# Expose the port
EXPOSE 8081

# Run the service
CMD ["./inventory-service"] 