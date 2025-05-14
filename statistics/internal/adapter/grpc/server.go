package grpc

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/baccala1010/e-commerce/statistics/internal/config"
	"github.com/baccala1010/e-commerce/statistics/internal/handler"
	"github.com/baccala1010/e-commerce/statistics/pkg/pb" // добавлен импорт pb
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Server represents the gRPC server for the statistics service
type Server struct {
	server   *grpc.Server
	cfg      *config.Config
	listener net.Listener
}

// NewServer creates a new gRPC server
func NewServer(cfg *config.Config, statisticsHandler *handler.StatisticsHandler) (*Server, error) {
	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("failed to listen on %s: %w", addr, err)
	}

	server := grpc.NewServer()
	pb.RegisterStatisticsServiceServer(server, statisticsHandler)

	// Enable reflection for tools like grpcurl
	reflection.Register(server)

	return &Server{
		server:   server,
		cfg:      cfg,
		listener: lis,
	}, nil
}

// Start starts the gRPC server
func (s *Server) Start(ctx context.Context) error {
	log.Printf("gRPC server started on %s", s.listener.Addr().String())

	errCh := make(chan error)
	go func() {
		if err := s.server.Serve(s.listener); err != nil {
			errCh <- fmt.Errorf("failed to serve gRPC: %w", err)
		}
	}()

	select {
	case <-ctx.Done():
		log.Println("Shutting down gRPC server...")
		s.server.GracefulStop()
		return nil
	case err := <-errCh:
		return err
	}
}

// Stop stops the gRPC server
func (s *Server) Stop() {
	if s.server != nil {
		s.server.GracefulStop()
	}
}
