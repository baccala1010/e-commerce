FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o statistics ./cmd/statistics

FROM alpine:latest

RUN apk --no-cache add ca-certificates net-tools

WORKDIR /app
COPY --from=builder /app/statistics /app/
COPY --from=builder /app/config/config.docker.yaml /app/config/config.docker.yaml

EXPOSE 50053
# Using a simpler health check that doesn't rely on HTTP endpoints
# since this service only exposes a gRPC port
HEALTHCHECK --interval=5s --timeout=3s --start-period=5s --retries=3 \
    CMD netstat -tulpn | grep 50053 || exit 1

CMD ["/app/statistics"]
