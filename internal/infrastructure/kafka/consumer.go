package kafka

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

// KafkaConsumer wraps the Kafka consumer functionality.
type KafkaConsumer struct {
	brokers   []string
	topic     string
	groupID   string
	segmentCh chan [][]byte
	consumer  sarama.ConsumerGroup
	logger    *zap.Logger
	pool      *sync.Pool
	waitGroup *sync.WaitGroup
}

// NewKafkaConsumer creates a new Kafka consumer instance.
func NewKafkaConsumer(logger *zap.Logger, brokers []string, groupID, topic string, minBytes, maxBytes int32, pollInterval, readTimeout time.Duration, bufferSize int) (*KafkaConsumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRange()
	config.Consumer.Fetch.Min = minBytes
	config.Consumer.Fetch.Max = maxBytes
	config.Consumer.MaxWaitTime = pollInterval
	config.Consumer.Group.Session.Timeout = readTimeout

	consumer, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		return nil, err
	}

	return &KafkaConsumer{
		logger:    logger,
		brokers:   brokers,
		topic:     topic,
		groupID:   groupID,
		segmentCh: make(chan [][]byte, bufferSize),
		consumer:  consumer,
		waitGroup: &sync.WaitGroup{},
		pool: &sync.Pool{
			New: func() interface{} {
				return make([][]byte, 0, 10)
			},
		},
	}, nil
}

func (kc *KafkaConsumer) Consume(ctx context.Context) (<-chan [][]byte, error) {
	handler := &consumerGroupHandler{segmentCh: kc.segmentCh, pool: kc.pool}
	defer kc.waitGroup.Add(1)
	go func() {
		defer kc.waitGroup.Done()
		for {
			if err := kc.consumer.Consume(ctx, []string{kc.topic}, handler); err != nil {
				log.Printf("Error consuming Kafka messages: %v", err)
			}
			if ctx.Err() != nil {
				return
			}
		}
	}()
	return kc.segmentCh, nil
}

func (kc *KafkaConsumer) Close() error {
	close(kc.segmentCh)
	err := kc.consumer.Close()
	if err != nil {
		return err
	}
	kc.waitGroup.Wait()
	return nil
}

type consumerGroupHandler struct {
	segmentCh chan [][]byte
	pool      *sync.Pool
}

func (c *consumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (c *consumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (c *consumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	batch := c.pool.Get().([][]byte)
	defer c.pool.Put(batch[:0])

	for message := range claim.Messages() {
		batch = append(batch, message.Value)
	}
	if len(batch) > 0 {
		select {
		case c.segmentCh <- batch:
			for message := range claim.Messages() {
				sess.MarkMessage(message, "")
			}
		default:
			return fmt.Errorf("error on sending message to channel")
		}
	}
	return nil
}
