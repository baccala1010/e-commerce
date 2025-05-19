FROM golang:1.20-alpine
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o statistics ./cmd/statistics
CMD ["/app/statistics"]
