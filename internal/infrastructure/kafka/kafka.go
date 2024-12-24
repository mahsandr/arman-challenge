package kafka

import (
	"time"

	"github.com/mahsandr/arman-challenge/internal/domain/repository"
	"go.uber.org/zap"
)

var _ repository.MessageBroker = &Kafka{}

type Kafka struct {
	*KafkaConsumer
	*KafkaProducer
}

func NewKafkaProducer(logger *zap.Logger, brokers []string, topic string, batchBytes int, flushInterval time.Duration) (*Kafka, error) {
	producer, err := NewProducer(logger, brokers, topic, batchBytes, flushInterval)
	if err != nil {
		return nil, err
	}
	return &Kafka{KafkaProducer: producer}, nil
}
func NewKafkaConsumer(logger *zap.Logger, brokers []string, groupID,
	topic string, minBytes, maxBytes int32, pollInterval time.Duration, bufferSize int) (*Kafka, error) {
	consumer, err := NewConsumer(logger, brokers, groupID, topic, minBytes, maxBytes, pollInterval, bufferSize)
	if err != nil {
		return nil, err
	}
	return &Kafka{KafkaConsumer: consumer}, nil
}

func (k *Kafka) Close() {
	if k.KafkaConsumer != nil {
		k.KafkaConsumer.Close()
	}
	if k.KafkaProducer != nil {
		k.KafkaProducer.Close()
	}
}
