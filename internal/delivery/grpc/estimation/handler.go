package grpc

import (
	"context"

	pb "github.com/mahsandr/arman-challenge/api/proto/generated/estimation/v1"
	"github.com/mahsandr/arman-challenge/internal/domain/usecases"
	"github.com/mahsandr/arman-challenge/internal/validation"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ pb.EstimationServiceServer = &Handler{}

type Handler struct {
	estimationService usecases.EstimationService
	logger            *zap.Logger
	validator         *validation.Validator
	pb.UnimplementedEstimationServiceServer
}

func NewHandler(logger *zap.Logger, estimationService usecases.EstimationService) *Handler {
	return &Handler{
		estimationService: estimationService,
		logger:            logger,
		validator:         validation.NewValidator(),
	}
}
func (h *Handler) EstimateUsers(ctx context.Context, in *pb.EstimateUsersRequest) (*pb.EstimateUsersResponse, error) {
	if in.Segment == "" {
		return nil, status.Error(codes.InvalidArgument, "segment is required")
	}
	userCount, err := h.estimationService.GetSegmentUsersCount(ctx, in.Segment)
	if err != nil {
		h.logger.Error("error on GetSegmentUsersCount", zap.Error(err))
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &pb.EstimateUsersResponse{
		Count: userCount,
	}, nil
}
