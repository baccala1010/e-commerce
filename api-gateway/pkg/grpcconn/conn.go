package grpcconn

import (
	"context"
	"fmt"
	"time"

	"github.com/baccala1010/e-commerce/api-gateway/internal/config"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Connection represents a gRPC client connection
type Connection struct {
	ServiceName string
	Client      *grpc.ClientConn
}

// ConnectionManager manages gRPC connections to various services
type ConnectionManager struct {
	connections map[string]*Connection
}

// NewConnectionManager creates a new connection manager
func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		connections: make(map[string]*Connection),
	}
}

// Connect establishes a connection to a gRPC service
func (m *ConnectionManager) Connect(ctx context.Context, serviceName string, cfg config.ServiceConfig) (*grpc.ClientConn, error) {
	// Check if connection already exists
	if conn, exists := m.connections[serviceName]; exists && conn.Client != nil {
		return conn.Client, nil
	}

	// Create connection options
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	}

	// Connect to the service
	address := fmt.Sprintf("%s:%d", cfg.GRPCHost, cfg.GRPCPort)
	logrus.Infof("Connecting to %s gRPC service at %s", serviceName, address)

	// Create a context with timeout
	timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(timeoutCtx, address, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to %s service: %w", serviceName, err)
	}

	// Store the connection
	m.connections[serviceName] = &Connection{
		ServiceName: serviceName,
		Client:      conn,
	}

	return conn, nil
}

// GetConnection returns an existing connection or creates a new one
func (m *ConnectionManager) GetConnection(ctx context.Context, serviceName string, cfg config.ServiceConfig) (*grpc.ClientConn, error) {
	if conn, exists := m.connections[serviceName]; exists && conn.Client != nil {
		return conn.Client, nil
	}
	return m.Connect(ctx, serviceName, cfg)
}

// Close closes all connections
func (m *ConnectionManager) Close() {
	for name, conn := range m.connections {
		if conn.Client != nil {
			logrus.Infof("Closing connection to %s service", name)
			conn.Client.Close()
		}
	}
}
