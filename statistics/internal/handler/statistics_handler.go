package handler

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/baccala1010/e-commerce/statistics/internal/usecase"
	"github.com/baccala1010/e-commerce/statistics/pkg/pb"
	"github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// StatisticsHandler implements the gRPC StatisticsService server interface
type StatisticsHandler struct {
	pb.UnimplementedStatisticsServiceServer
	statisticsUsecase usecase.StatisticsUsecase
}

// NewStatisticsHandler creates a new statistics gRPC handler
func NewStatisticsHandler(statisticsUsecase usecase.StatisticsUsecase) *StatisticsHandler {
	return &StatisticsHandler{
		statisticsUsecase: statisticsUsecase,
	}
}

// GetUserOrdersStatistics retrieves order statistics for a specific user
func (h *StatisticsHandler) GetUserOrdersStatistics(ctx context.Context, req *pb.UserOrdersStatisticsRequest) (*pb.UserOrdersStatisticsResponse, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	// Call use case to get user order statistics
	stats, err := h.statisticsUsecase.GetUserOrdersStatistics(ctx, req.UserId)
	if err != nil {
		log.Printf("Failed to get user order statistics: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to get statistics: %v", err)
	}

	// Convert domain model to protobuf response
	response := &pb.UserOrdersStatisticsResponse{
		UserId:           stats.UserID,
		TotalOrders:      int32(stats.TotalOrders),
		TotalSpent:       stats.TotalSpent,
		AverageOrderValue: stats.AverageOrderValue,
	}

	// Add time-related fields if available
	if !stats.FirstOrderAt.IsZero() {
		response.FirstOrderAt = &timestamp.Timestamp{
			Seconds: stats.FirstOrderAt.Unix(),
			Nanos:   int32(stats.FirstOrderAt.Nanosecond()),
		}
	}

	if !stats.LastOrderAt.IsZero() {
		response.LastOrderAt = &timestamp.Timestamp{
			Seconds: stats.LastOrderAt.Unix(),
			Nanos:   int32(stats.LastOrderAt.Nanosecond()),
		}
	}

	// Add order time distribution
	for _, item := range stats.OrderTimeDistribution {
		response.OrderTimeDistribution = append(response.OrderTimeDistribution, &pb.OrderTimeOfDay{
			Hour:       item.Hour,
			OrderCount: int32(item.OrderCount),
		})
	}

	return response, nil
}

// GetUserStatistics retrieves general user statistics
func (h *StatisticsHandler) GetUserStatistics(ctx context.Context, req *pb.UserStatisticsRequest) (*pb.UserStatisticsResponse, error) {
	// Call use case to get general user statistics
	stats, err := h.statisticsUsecase.GetUserStatistics(ctx)
	if err != nil {
		log.Printf("Failed to get user statistics: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to get user statistics: %v", err)
	}

	response := &pb.UserStatisticsResponse{
		TotalRegisteredUsers: int32(stats.TotalRegisteredUsers),
	}

	return response, nil
}