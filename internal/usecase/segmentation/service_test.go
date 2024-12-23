package segment

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/mahsandr/arman-challenge/internal/domain/models"
	mock_repository "github.com/mahsandr/arman-challenge/internal/mocks/repository"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestService_AddUserSegment(t *testing.T) {
	tests := []struct {
		name          string
		segment       *models.UserSegment
		produceErr    error
		expectedError error
		initMock      func(mockProducer *mock_repository.MockMessageBroker)
	}{
		{
			name: "success",
			segment: &models.UserSegment{
				UserID:  "user1",
				Segment: "segment1",
			},
			initMock: func(mockProducer *mock_repository.MockMessageBroker) {
				mockProducer.EXPECT().Produce(gomock.Any(), gomock.Any()).Return(nil).Times(1)
			},
		},
		{
			name: "produce error",
			segment: &models.UserSegment{
				UserID:  "user2",
				Segment: "segment2",
			},
			produceErr:    errors.New("produce error"),
			expectedError: errors.New("produce error"),
			initMock: func(mockProducer *mock_repository.MockMessageBroker) {
				mockProducer.EXPECT().Produce(gomock.Any(), gomock.Any()).Return(errors.New("produce error")).Times(1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockProducer := mock_repository.NewMockMessageBroker(ctrl)
			logger, _ := zap.NewProduction()

			service := NewService(nil, mockProducer, logger)

			tt.initMock(mockProducer)

			err := service.AddUserSegment(context.Background(), tt.segment)
			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
