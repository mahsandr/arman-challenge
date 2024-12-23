package grpc

import (
	"context"

	pb "github.com/mahsandr/arman-challenge/api/proto/generated/segmentation/v1"
	"github.com/mahsandr/arman-challenge/internal/domain/models"
	"github.com/mahsandr/arman-challenge/internal/domain/usecases"
	"github.com/mahsandr/arman-challenge/internal/validation"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ pb.SegmentationServiceServer = &Handler{}

type Handler struct {
	segmentService usecases.SegmentService
	logger         *zap.Logger
	validator      *validation.Validator
	pb.UnimplementedSegmentationServiceServer
}

func NewHandler(logger *zap.Logger, segmentService usecases.SegmentService) *Handler {
	return &Handler{
		segmentService: segmentService,
		logger:         logger,
		validator:      validation.NewValidator(),
	}
}
func (h *Handler) AddUserSegment(ctx context.Context, in *pb.AddUserSegmentRequest) (*pb.AddUserSegmentResponse, error) {
	userSegment := &models.UserSegment{
		UserID:  in.UserSegment.UserId,
		Segment: in.UserSegment.Segment,
	}
	if validateErr := h.validator.ValidateStruct(userSegment); validateErr != "" {
		return nil, status.Error(codes.InvalidArgument, validateErr)
	}
	err := h.segmentService.AddUserSegment(ctx, userSegment)
	if err != nil {
		h.logger.Error("error on AddUserSegment", zap.Error(err))
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &pb.AddUserSegmentResponse{}, nil
}
