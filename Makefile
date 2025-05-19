.PHONY: generate-proto clean

# Protobuf generation
generate-proto:
	@echo "Generating protobuf code..."
	protoc --go_out=paths=source_relative:inventory/pkg/pb --go-grpc_out=paths=source_relative:inventory/pkg/pb \
		-I proto/ proto/inventory/inventory.proto
	protoc --go_out=paths=source_relative:order/pkg/pb --go-grpc_out=paths=source_relative:order/pkg/pb \
		-I proto/ proto/order/order.proto
	protoc --go_out=paths=source_relative:statistics/pkg/pb --go-grpc_out=paths=source_relative:statistics/pkg/pb \
		-I proto/ proto/statistics/statistics.proto

# Clean generated files
clean:
	@echo "Cleaning generated files..."
	rm -rf inventory/pkg/pb
	rm -rf order/pkg/pb


# Build all services
build:
	@echo "Building inventory service..."
	cd inventory && go build -o bin/inventory cmd/inventory/main.go
	@echo "Building order service..."
	cd order && go build -o bin/order cmd/order/main.go
	@echo "Building API gateway..."
	cd api-gateway && go build -o bin/api-gateway cmd/api-gateway/main.go

# Run all services
run:
	@echo "Running services..."
	docker-compose up -d
