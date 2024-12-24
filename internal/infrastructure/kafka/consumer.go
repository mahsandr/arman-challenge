package kafka

import (
	"context"
	"sync"
	"time"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

// KafkaConsumer wraps the Kafka consumer functionality.
type KafkaConsumer struct {
	topic      string
	segmentCh  chan [][]byte
	consumer   sarama.ConsumerGroup
	logger     *zap.Logger
	waitGroup  sync.WaitGroup
	bufferSize int
}

// NewConsumer creates a new Kafka consumer instance.
func NewConsumer(logger *zap.Logger, brokers []string, groupID, topic string,
	minBytes, maxBytes int32, pollInterval time.Duration, bufferSize int) (*KafkaConsumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRange()
	config.Consumer.Fetch.Min = minBytes
	config.Consumer.Fetch.Max = maxBytes
	config.Consumer.MaxWaitTime = pollInterval

	consumer, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		return nil, err
	}

	return &KafkaConsumer{
		bufferSize: bufferSize,
		logger:     logger,
		topic:      topic,
		segmentCh:  make(chan [][]byte, 100), // Buffered channel to handle bursts
		consumer:   consumer,
	}, nil
}

// Consume starts consuming messages from the Kafka topic.
func (kc *KafkaConsumer) Consume(ctx context.Context) (<-chan [][]byte, error) {
	handler := &consumerGroupHandler{
		segmentCh:  kc.segmentCh,
		logger:     kc.logger,
		bufferSize: kc.bufferSize,
		pool: &sync.Pool{
			New: func() interface{} {
				return make([]byte, 0, 1024) // Adjust the buffer size as needed
			},
		},
	}

	kc.waitGroup.Add(1)
	go func() {
		defer kc.waitGroup.Done()
		for {
			if ctx.Err() != nil {
				return
			}

			if err := kc.consumer.Consume(ctx, []string{kc.topic}, handler); err != nil {
				kc.logger.Error("Error consuming Kafka messages", zap.Error(err))
				if err == sarama.ErrClosedConsumerGroup {
					return
				}
			}
		}
	}()

	return kc.segmentCh, nil
}

// Close shuts down the Kafka consumer gracefully.
func (kc *KafkaConsumer) Close() error {
	err := kc.consumer.Close()
	kc.waitGroup.Wait()
	close(kc.segmentCh)
	return err
}

type consumerGroupHandler struct {
	segmentCh  chan [][]byte
	pool       *sync.Pool
	logger     *zap.Logger
	bufferSize int
}

func (c *consumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (c *consumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (c *consumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	batch := make([][]byte, 0, c.bufferSize)
	defer func() {
		if len(batch) > 0 {
			_ = c.sendBatch(sess.Context(), batch)
		}
	}()
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				return nil
			}

			msgCopy := c.pool.Get().([]byte)[:len(message.Value)]
			copy(msgCopy, message.Value)
			batch = append(batch, msgCopy)
			c.pool.Put(&msgCopy)

			sess.MarkMessage(message, "")

			if len(batch) >= c.bufferSize {
				if err := c.sendBatch(sess.Context(), batch); err != nil {
					c.logger.Error("Failed to send batch", zap.Error(err))
				}
				batch = batch[:0] // Reset the batch
			}

		case <-sess.Context().Done():
			return nil
		}
	}
}

func (c *consumerGroupHandler) sendBatch(ctx context.Context, msgs [][]byte) error {
	batch := make([][]byte, len(msgs))
	copy(batch, msgs)

	select {
	case c.segmentCh <- batch:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
