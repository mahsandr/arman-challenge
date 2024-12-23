package kafka

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/IBM/sarama"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

type mockConsumerGroup struct {
	mock.Mock
}

func (m *mockConsumerGroup) Consume(ctx context.Context, topics []string, handler sarama.ConsumerGroupHandler) error {
	args := m.Called(ctx, topics, handler)
	return args.Error(0)
}

func (m *mockConsumerGroup) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *mockConsumerGroup) Pause(partitions map[string][]int32) {
	m.Called(partitions)
}

func (m *mockConsumerGroup) Resume(partitions map[string][]int32) {
	m.Called(partitions)
}

func (m *mockConsumerGroup) Errors() <-chan error {
	args := m.Called()
	return args.Get(0).(<-chan error)
}

func (m *mockConsumerGroup) ResumeAll() {
	m.Called()
}

func (m *mockConsumerGroup) PauseAll() {
	m.Called()
}

func TestNewKafkaConsumer(t *testing.T) {
	logger, _ := zap.NewProduction()
	brokers := []string{"localhost:9092"}
	groupID := "test-group"
	topic := "test-topic"
	minBytes := int32(1)
	maxBytes := int32(10)
	pollInterval := 100 * time.Millisecond
	readTimeout := 10 * time.Second
	bufferSize := 100

	consumer, err := NewKafkaConsumer(logger, brokers, groupID, topic, minBytes, maxBytes, pollInterval, readTimeout, bufferSize)
	require.NoError(t, err)
	assert.NotNil(t, consumer)
	assert.Equal(t, brokers, consumer.brokers)
	assert.Equal(t, groupID, consumer.groupID)
	assert.Equal(t, topic, consumer.topic)
	assert.NotNil(t, consumer.segmentCh)
	assert.NotNil(t, consumer.consumer)
}

func TestKafkaConsumer_Consume(t *testing.T) {
	logger, _ := zap.NewProduction()
	brokers := []string{"localhost:9092"}
	groupID := "test-group"
	topic := "test-topic"
	bufferSize := 100

	mockConsumer := new(mockConsumerGroup)
	consumer := &KafkaConsumer{
		logger:    logger,
		brokers:   brokers,
		topic:     topic,
		groupID:   groupID,
		segmentCh: make(chan [][]byte, bufferSize),
		consumer:  mockConsumer,
		waitGroup: &sync.WaitGroup{},
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mockConsumer.On("Consume", ctx, []string{topic}, mock.Anything).Return(nil).Once()

	segmentCh, err := consumer.Consume(ctx)
	require.NoError(t, err)
	assert.NotNil(t, segmentCh)
	consumer.waitGroup.Wait()
	mockConsumer.AssertExpectations(t)
}

func TestKafkaConsumer_Close(t *testing.T) {
	logger, _ := zap.NewProduction()
	brokers := []string{"localhost:9092"}
	groupID := "test-group"
	topic := "test-topic"
	bufferSize := 100

	mockConsumer := new(mockConsumerGroup)
	consumer := &KafkaConsumer{
		logger:    logger,
		brokers:   brokers,
		topic:     topic,
		groupID:   groupID,
		segmentCh: make(chan [][]byte, bufferSize),
		consumer:  mockConsumer,
		waitGroup: &sync.WaitGroup{},
	}

	mockConsumer.On("Close").Return(nil).Once()

	err := consumer.Close()
	require.NoError(t, err)

	mockConsumer.AssertExpectations(t)
}
