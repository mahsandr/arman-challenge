package kafka

import (
	"github.com/mahsandr/arman-challenge/internal/domain/repository"
)

var _ repository.MessageBroker = &Kafka{}

type Kafka struct {
	KafkaConsumer
	KafkaProducer
}
