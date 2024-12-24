package kafka

import (
	"context"
	"fmt"
	"time"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

// KafkaProducer wraps the Kafka producer functionality.
type KafkaProducer struct {
	topic    string
	producer sarama.SyncProducer
	logger   *zap.Logger
}

func NewProducer(logger *zap.Logger, brokers []string, topic string, batchBytes int, flushInterval time.Duration) (*KafkaProducer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Compression = sarama.CompressionSnappy
	config.Producer.Flush.Bytes = batchBytes
	config.Producer.Flush.Frequency = flushInterval
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("error creating kafka producer: %v", err)
	}

	return &KafkaProducer{
		producer: producer,
		topic:    topic,
		logger:   logger,
	}, nil
}

// Produce sends messages to the Kafka topic.
func (kp *KafkaProducer) Produce(ctx context.Context, msg []byte) error {
	kafkaMsg := &sarama.ProducerMessage{
		Topic: kp.topic,
		Value: sarama.StringEncoder(msg),
	}

	if _, _, err := kp.producer.SendMessage(kafkaMsg); err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}
	return nil
}

// Close cleans up resources used by the Kafka producer.
func (kp *KafkaProducer) Close() error {
	return kp.producer.Close()
}
