package estimation

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	mock_repository "github.com/mahsandr/arman-challenge/internal/mocks/repository"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestService_GetSegmentUsersCount(t *testing.T) {
	tests := []struct {
		name          string
		segment       string
		count         uint32
		repoErr       error
		expectedCount uint32
		expectedError error
		initMock      func(mockRepo *mock_repository.MockSegmentRepository)
	}{
		{
			name:          "success",
			segment:       "segment1",
			count:         10,
			expectedCount: 10,
			initMock: func(mockRepo *mock_repository.MockSegmentRepository) {
				mockRepo.EXPECT().GetSegmentUsersCount(gomock.Any(), "segment1").Return(uint32(10), nil).Times(1)
			},
		},
		{
			name:          "repository error",
			segment:       "segment2",
			repoErr:       errors.New("repository error"),
			expectedError: errors.New("repository error"),
			initMock: func(mockRepo *mock_repository.MockSegmentRepository) {
				mockRepo.EXPECT().GetSegmentUsersCount(gomock.Any(), "segment2").Return(uint32(0), errors.New("repository error")).Times(1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mock_repository.NewMockSegmentRepository(ctrl)
			logger := zap.NewNop()

			service := NewService(mockRepo, logger)

			tt.initMock(mockRepo)

			count, err := service.GetSegmentUsersCount(context.Background(), tt.segment)
			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
				assert.Equal(t, tt.expectedCount, count)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedCount, count)
			}
		})
	}
}
