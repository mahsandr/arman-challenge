package estimation

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	mock_repository "github.com/mahsandr/arman-challenge/internal/mocks/repository"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func TestSegmentConsumer_StartConsuming(t *testing.T) {
	tests := []struct {
		name          string
		consumeErr    error
		segments      [][]byte
		unmarshalErr  error
		expectedError error
		initMock      func(segmentChannel chan [][]byte,
			mockConsumer *mock_repository.MockMessageBroker,
			mockRepo *mock_repository.MockSegmentRepository)
		logEntry   []*zapcore.Entry
		logContext []map[string]interface{}
	}{
		{
			name: "success",
			segments: [][]byte{
				[]byte(`{"UserID":"user1","Segment":"segment1"}`),
			},
			initMock: func(segmentChannel chan [][]byte,
				mockConsumer *mock_repository.MockMessageBroker,
				mockRepo *mock_repository.MockSegmentRepository) {
				mockConsumer.EXPECT().Consume(gomock.Any()).Return(segmentChannel, nil).Times(1)
				mockRepo.EXPECT().SaveUserSegments(gomock.Any(), gomock.Any()).Return(nil).Times(1)
			},
		},
		{
			name:          "consume error",
			consumeErr:    errors.New("consume error"),
			expectedError: errors.New("consume error"),
			initMock: func(segmentChannel chan [][]byte,
				mockConsumer *mock_repository.MockMessageBroker,
				mockRepo *mock_repository.MockSegmentRepository) {
				mockConsumer.EXPECT().Consume(gomock.Any()).Return(nil, errors.New("consume error")).Times(1)
			},
			logEntry: []*zapcore.Entry{
				{
					Level:   zapcore.ErrorLevel,
					Message: "error on consuming",
				},
			},
			logContext: []map[string]interface{}{
				{
					"error": "consume error",
				},
			},
		},
		{
			name: "unmarshal error",
			segments: [][]byte{
				[]byte(`invalid json`),
			},
			initMock: func(segmentChannel chan [][]byte,
				mockConsumer *mock_repository.MockMessageBroker,
				mockRepo *mock_repository.MockSegmentRepository) {
				mockConsumer.EXPECT().Consume(gomock.Any()).Return(segmentChannel, nil).Times(1)
			},
			logEntry: []*zapcore.Entry{
				{
					Level:   zapcore.ErrorLevel,
					Message: "error unmarshaling segment",
				},
			},
			logContext: []map[string]interface{}{
				{
					//nolint: lll
					"error": "readObjectStart: expect { or n, but found i, error found in #1 byte of ...|invalid jso|..., bigger context ...|invalid json|...",
				},
			},
		},
		{
			name: "save error",
			segments: [][]byte{
				[]byte(`{"UserID":"user1","Segment":"segment1"}`),
			},
			initMock: func(segmentChannel chan [][]byte,
				mockConsumer *mock_repository.MockMessageBroker,
				mockRepo *mock_repository.MockSegmentRepository) {
				mockConsumer.EXPECT().Consume(gomock.Any()).Return(segmentChannel, nil).Times(1)
				mockRepo.EXPECT().SaveUserSegments(gomock.Any(), gomock.Any()).Return(errors.New("save error")).Times(1)
			},
			logEntry: []*zapcore.Entry{
				{
					Level:   zapcore.ErrorLevel,
					Message: "error saving segments",
				},
			},
			logContext: []map[string]interface{}{
				{
					"error": "save error",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockConsumer := mock_repository.NewMockMessageBroker(ctrl)
			mockConsumer.EXPECT().Close().AnyTimes()

			mockRepo := mock_repository.NewMockSegmentRepository(ctrl)

			observerLog, out := observer.New(zapcore.InfoLevel)
			logger := zap.New(observerLog)
			uc := NewSegmentConsumer(mockRepo, mockConsumer, logger)

			segmentChannel := make(chan [][]byte, 1)
			go func() {
				defer close(segmentChannel)
				segmentChannel <- tt.segments
				time.Sleep(time.Second * 1)
			}()
			tt.initMock(segmentChannel, mockConsumer, mockRepo)

			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()

			err := uc.StartConsuming(ctx)
			uc.waitGroup.Wait()

			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
			logs := out.TakeAll()
			for i, log := range logs {
				assert.Equal(t, tt.logEntry[i].Level, log.Level)
				assert.Equal(t, tt.logEntry[i].Message, log.Message)
				assert.Equal(t, tt.logContext[i], log.ContextMap())
			}
		})
	}
}
