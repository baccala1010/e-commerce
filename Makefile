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
	protoc --go_out=paths=source_relative:events/pkg/pb --go-grpc_out=paths=source_relative:events/pkg/pb \
		-I proto/ proto/events/events.proto
	@echo "Moving generated files to correct locations..."
	cp inventory/pkg/pb/inventory/*.pb.go inventory/pkg/pb/
	cp order/pkg/pb/order/*.pb.go order/pkg/pb/
	cp statistics/pkg/pb/statistics/*.pb.go statistics/pkg/pb/
	cp events/pkg/pb/events/*.pb.go events/pkg/pb/
	rm -rf inventory/pkg/pb/inventory
	rm -rf order/pkg/pb/order
	rm -rf statistics/pkg/pb/statistics
	rm -rf events/pkg/pb/events

# Clean generated files
clean:
	@echo "Cleaning generated files..."
	rm -rf inventory/pkg/pb
	rm -rf order/pkg/pb
	rm -rf statistics/pkg/pb
	rm -rf events/pkg/pb

# Build all services
build:
	@echo "Building inventory service..."
	cd inventory && go build -o bin/inventory cmd/inventory/main.go
	@echo "Building order service..."
	cd order && go build -o bin/order cmd/order/main.go
	@echo "Building statistics service..."
	cd statistics && go build -o bin/statistics cmd/statistics/main.go
	@echo "Building API gateway..."
	cd api-gateway && go build -o bin/api-gateway cmd/api-gateway/main.go

# Run all services
run:
	@echo "Running services..."
	docker-compose up -d
