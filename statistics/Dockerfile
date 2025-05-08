FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install necessary packages for compilation
RUN apk add --no-cache gcc musl-dev

# Copy and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -a -o statistics ./cmd/statistics/main.go

FROM alpine:latest

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/statistics .
COPY --from=builder /app/config ./config

# Set environment variable to use Docker configuration
ENV DOCKER=true

# Run the binary
CMD ["./statistics"]