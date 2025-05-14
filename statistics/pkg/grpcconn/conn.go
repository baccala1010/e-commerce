package grpcconn

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Options represents gRPC client connection options
type Options struct {
	Address     string
	DialTimeout time.Duration
}

// Connect establishes a gRPC client connection
func Connect(ctx context.Context, opts Options) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(ctx, opts.DialTimeout)
	defer cancel()

	conn, err := grpc.DialContext(
		ctx,
		opts.Address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to gRPC server at %s: %w", opts.Address, err)
	}

	return conn, nil
}
