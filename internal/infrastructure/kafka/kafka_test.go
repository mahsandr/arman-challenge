package kafka

import (
	"context"
	"testing"
	"time"

	"github.com/IBM/sarama"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestKafkaProducer(t *testing.T) {
	broker := sarama.NewMockBroker(t, 1)
	defer broker.Close()

	mockMetadataResponse := sarama.NewMockMetadataResponse(t).
		SetBroker(broker.Addr(), broker.BrokerID()).
		SetLeader("test-topic", 0, broker.BrokerID())
	broker.SetHandlerByMap(map[string]sarama.MockResponse{
		"MetadataRequest": mockMetadataResponse,
		"ProduceRequest":  sarama.NewMockProduceResponse(t),
	})

	logger := zap.NewNop()
	defer logger.Sync()

	producer, err := NewProducer(logger, []string{broker.Addr()}, "test-topic", 100, 10*time.Millisecond)
	assert.NoError(t, err)
	defer producer.Close()

	err = producer.Produce(context.Background(), []byte("test-message"))
	assert.NoError(t, err)
}
func TestKafkaConsumer(t *testing.T) {
	// Mock broker setup
	broker := sarama.NewMockBroker(t, 1)
	defer broker.Close()

	// Create mock consumer group metadata response
	groupMetadataResponse := sarama.NewMockFindCoordinatorResponse(t).
		SetCoordinator(sarama.CoordinatorGroup, "group-id", broker)

	// Create mock offset fetch response
	offsetFetchResponse := sarama.NewMockOffsetFetchResponse(t).
		SetOffset("group-id", "test-topic", 0, 0, "", sarama.ErrNoError)

	// Create mock consumer metadata response
	metadataResponse := sarama.NewMockMetadataResponse(t).
		SetBroker(broker.Addr(), broker.BrokerID()).
		SetLeader("test-topic", 0, broker.BrokerID())

	// Create mock offset response
	offsetResponse := sarama.NewMockOffsetResponse(t).
		SetOffset("test-topic", 0, sarama.OffsetOldest, 0).
		SetOffset("test-topic", 0, sarama.OffsetNewest, 1)

	// Create mock fetch response with test message
	fetchResponse := sarama.NewMockFetchResponse(t, 1).
		SetMessage("test-topic", 0, 0, sarama.StringEncoder("test-message"))

	// Set up the mock handler responses
	broker.SetHandlerByMap(map[string]sarama.MockResponse{
		"MetadataRequest":        metadataResponse,
		"OffsetRequest":          offsetResponse,
		"FindCoordinatorRequest": groupMetadataResponse,
		"JoinGroupRequest": sarama.NewMockJoinGroupResponse(t).
			SetGroupProtocol(sarama.RangeBalanceStrategyName).
			SetMemberId("test-member-id").
			SetGenerationId(1).
			SetLeaderId("test-member-id"),
		"SyncGroupRequest": sarama.NewMockSyncGroupResponse(t).
			SetMemberAssignment(&sarama.ConsumerGroupMemberAssignment{
				UserData: []byte("test-topic"),
				Topics: map[string][]int32{
					"test-topic": {0},
				},
			}).
			SetError(sarama.ErrNoError),
		"HeartbeatRequest":    sarama.NewMockHeartbeatResponse(t),
		"OffsetFetchRequest":  offsetFetchResponse,
		"FetchRequest":        fetchResponse,
		"OffsetCommitRequest": sarama.NewMockOffsetCommitResponse(t),
	})

	// Logger setup
	logger := zap.NewNop()
	defer logger.Sync()

	// KafkaConsumer setup with modified configuration
	config := sarama.NewConfig()
	config.Version = sarama.V2_0_0_0 // Set an explicit version
	config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRange()
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Consumer.Return.Errors = true

	consumer, err := NewConsumer(logger, []string{broker.Addr()}, "group-id", "test-topic", 1, 1048576, 100*time.Millisecond, 1)
	assert.NoError(t, err)

	// Start consuming with a longer timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer consumer.Close()

	messageCh, err := consumer.Consume(ctx)
	assert.NoError(t, err)

	// Wait for messages
	for {
		select {
		case messages := <-messageCh:
			assert.Greater(t, len(messages), 0, "No messages received")
			assert.Equal(t, "test-message", string(messages[0]), "Incorrect message received")
			cancel()
			return
		case <-ctx.Done():
			t.Fatalf("Timeout waiting for messages: %v", ctx.Err())
		}
	}
}
