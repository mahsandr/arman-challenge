package config

import (
	"fmt"
	"log"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type KafkaConfig struct {
	Brokers        []string      `env:"KAFKA_BROKERS,notEmpty" envSeparator:","`
	GroupID        string        `env:"KAFKA_GROUP_ID"`
	Topic          string        `env:"KAFKA_TOPIC"`
	MinBytes       int32         `env:"KAFKA_MIN_BYTES" envDefault:"100"`
	MaxBytes       int32         `env:"KAFKA_MAX_BYTES" envDefault:"100"`
	PollInterval   time.Duration `env:"KAFKA_POLL_INTERVAL_MS" envDefault:"10s"`
	ReadTimeout    time.Duration `env:"KAFKA_READ_TIMEOUT_MS" envDefault:"10s"`
	Retries        int           `env:"KAFKA_RETRIES" envDefault:"3"`
	BatchBytes     int           `env:"KAFKA_BATCH_BYTES" envDefault:"1048576"`
	FlushInterval  time.Duration `env:"KAFKA_FLUSH_INTERVAL"  envDefault:"30s"`
	ConsumerBuffer int           `env:"KAFKA_CONSUMER_BUFFER" envDefault:"1000"`
}

func GetKafkaConfig(options ...string) (*KafkaConfig, error) {
	envPath := ".env" // default .env file path
	if len(options) > 0 {
		envPath = options[0]
	}
	// Load .env file
	if err := godotenv.Load(envPath); err != nil {
		log.Printf("error loading .env file from %s, using default config: %v", envPath, err)
	}
	envs := &KafkaConfig{}
	err := env.Parse(envs)
	if err != nil {
		return nil, fmt.Errorf("error parsing env: %v", err)
	}
	return envs, err
}
