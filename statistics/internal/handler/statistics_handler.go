package handler

import (
	"context"

	"github.com/baccala1010/e-commerce/statistics/internal/usecase"
	"github.com/baccala1010/e-commerce/statistics/pkg/pb"
	"github.com/google/uuid"
)

type StatisticsHandler struct {
	uc                                      usecase.StatisticsUsecase
	pb.UnimplementedStatisticsServiceServer // Embed to satisfy gRPC server interface
}

func NewStatisticsHandler(uc usecase.StatisticsUsecase) *StatisticsHandler {
	return &StatisticsHandler{
		uc:                                   uc,
		UnimplementedStatisticsServiceServer: pb.UnimplementedStatisticsServiceServer{},
	}
}

func (h *StatisticsHandler) GetUserOrdersStatistics(ctx context.Context, req *pb.UserOrderStatisticsRequest) (*pb.UserOrderStatisticsResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, err
	}
	count, err := h.uc.GetUserOrderCount(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &pb.UserOrderStatisticsResponse{OrderCount: int32(count)}, nil
}

func (h *StatisticsHandler) GetAllUserOrdersStatistics(ctx context.Context, req *pb.GetAllUserOrdersStatisticsRequest) (*pb.GetAllUserOrdersStatisticsResponse, error) {
	// Set default pagination parameters if not provided
	page := int(req.Page)
	if page < 1 {
		page = 1
	}

	pageSize := int(req.PageSize)
	if pageSize < 1 {
		pageSize = 10 // Default page size
	}

	// Call the usecase
	stats, total, err := h.uc.GetAllUserOrderStatistics(ctx, page, pageSize)
	if err != nil {
		return nil, err
	}

	// Map to protobuf response
	protoStats := make([]*pb.UserStatistic, len(stats))
	for i, stat := range stats {
		protoStats[i] = &pb.UserStatistic{
			UserId:     stat.UserID,
			OrderCount: int32(stat.OrderCount),
		}
	}

	return &pb.GetAllUserOrdersStatisticsResponse{
		Statistics: protoStats,
		TotalCount: int32(total),
	}, nil
}
