package grpc

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	pb "github.com/mahsandr/arman-challenge/api/proto/generated/estimation/v1"
	mock_usecases "github.com/mahsandr/arman-challenge/internal/mocks/usecases"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestHandler_EstimateUsers(t *testing.T) {
	tests := []struct {
		name          string
		request       *pb.EstimateUsersRequest
		mockService   func(m *mock_usecases.MockEstimationService)
		expectedError error
		expectedCount uint32
	}{
		{
			name: "success",
			request: &pb.EstimateUsersRequest{
				Segment: "segment1",
			},
			mockService: func(m *mock_usecases.MockEstimationService) {
				m.EXPECT().GetSegmentUsersCount(gomock.Any(), "segment1").Return(uint32(100), nil)
			},
			expectedError: nil,
			expectedCount: 100,
		},
		{
			name: "validation error",
			request: &pb.EstimateUsersRequest{
				Segment: "",
			},
			mockService:   func(m *mock_usecases.MockEstimationService) {},
			expectedError: status.Error(codes.InvalidArgument, "segment is required"),
			expectedCount: 0,
		},
		{
			name: "service error",
			request: &pb.EstimateUsersRequest{
				Segment: "segment2",
			},
			mockService: func(m *mock_usecases.MockEstimationService) {
				m.EXPECT().GetSegmentUsersCount(gomock.Any(), "segment2").Return(uint32(0), errors.New("service error"))
			},
			expectedError: status.Error(codes.Internal, "internal error"),
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockService := mock_usecases.NewMockEstimationService(ctrl)
			logger := zap.NewNop()
			handler := NewHandler(logger, mockService)

			tt.mockService(mockService)

			resp, err := handler.EstimateUsers(context.Background(), tt.request)
			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, tt.expectedCount, resp.Count)
			}
		})
	}
}
