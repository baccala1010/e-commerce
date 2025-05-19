package grpc

import (
	"log"
	"net"

	"github.com/baccala1010/e-commerce/statistics/pkg/pb"
	"google.golang.org/grpc"
)

func RunGRPCServer(handler pb.StatisticsServiceServer, port string) error {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}
	grpcServer := grpc.NewServer()
	pb.RegisterStatisticsServiceServer(grpcServer, handler)
	log.Printf("gRPC server listening on %s", port)
	return grpcServer.Serve(lis)
}
