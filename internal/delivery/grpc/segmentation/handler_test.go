package grpc

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	pb "github.com/mahsandr/arman-challenge/api/proto/generated/segmentation/v1"
	"github.com/mahsandr/arman-challenge/internal/domain/models"
	mock_usecases "github.com/mahsandr/arman-challenge/internal/mocks/usecases"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestHandler_AddUserSegment(t *testing.T) {
	tests := []struct {
		name          string
		request       *pb.AddUserSegmentRequest
		mockService   func(m *mock_usecases.MockSegmentService)
		expectedError error
	}{
		{
			name: "success",
			request: &pb.AddUserSegmentRequest{
				UserSegment: &pb.UserSegment{
					UserId:  "user1",
					Segment: "segment1",
				},
			},
			mockService: func(m *mock_usecases.MockSegmentService) {
				m.EXPECT().AddUserSegment(gomock.Any(), &models.UserSegment{
					UserID:  "user1",
					Segment: "segment1",
				}).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "validation error",
			request: &pb.AddUserSegmentRequest{
				UserSegment: &pb.UserSegment{
					UserId:  "",
					Segment: "segment1",
				},
			},
			mockService:   func(m *mock_usecases.MockSegmentService) {},
			expectedError: status.Error(codes.InvalidArgument, "Key: 'UserSegment.UserID' Error:Field validation for 'UserID' failed on the 'required' tag"),
		},
		{
			name: "service error",
			request: &pb.AddUserSegmentRequest{
				UserSegment: &pb.UserSegment{
					UserId:  "user2",
					Segment: "segment2",
				},
			},
			mockService: func(m *mock_usecases.MockSegmentService) {
				m.EXPECT().AddUserSegment(gomock.Any(), gomock.Any()).
					Return(errors.New("service error"))
			},
			expectedError: status.Error(codes.Internal, "internal error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockService := mock_usecases.NewMockSegmentService(ctrl)
			logger := zap.NewNop()
			handler := NewHandler(logger, mockService)

			tt.mockService(mockService)

			_, err := handler.AddUserSegment(context.Background(), tt.request)
			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
