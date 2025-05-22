#!/bin/bash

# Navigate to the root directory of the project
cd "$(dirname "$0")/../.." || exit

# Ensure the output directory exists
mkdir -p statistics/pkg/pb

# Generate Go code from the protobuf definitions
protoc --go_out=. --go_opt=paths=source_relative \
  --go-grpc_out=. --go-grpc_opt=paths=source_relative \
  proto/statistics/statistics.proto

# Move the generated files to the correct location if they're not already there
mkdir -p statistics/pkg/pb
[ -f proto/statistics/statistics.pb.go ] && mv proto/statistics/statistics.pb.go statistics/pkg/pb/
[ -f proto/statistics/statistics_grpc.pb.go ] && mv proto/statistics/statistics_grpc.pb.go statistics/pkg/pb/

echo "Proto files generated successfully in statistics/pkg/pb!"
