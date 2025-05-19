package handler

import (
	"context"

	"github.com/baccala1010/e-commerce/statistics/internal/usecase"
	"github.com/baccala1010/e-commerce/statistics/pkg/pb"
	"github.com/google/uuid"
)

type StatisticsHandler struct {
	uc usecase.StatisticsUsecase
	pb.UnimplementedStatisticsServiceServer
}

func NewStatisticsHandler(uc usecase.StatisticsUsecase) *StatisticsHandler {
	return &StatisticsHandler{uc: uc}
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
