package config

import (
	"fmt"
	"log"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

// Config holds the configuration values
type ServiceConfig struct {
	Host                 string `env:"SERVICE_HOST" envDefault:"0.0.0.0"`
	Port                 string `env:"SERVICE_PORT" envDefault:"3000"`
	ClickHouseAddr       string `env:"CLICKHOUSE_ADDR" envDefault:"host=127.0.0.1"`
	UserSegmentTableName string `env:"USER_SEGMENT_TABLE_NAME" envDefault:"user_segments"`
	SegmentsViewName     string `env:"SEGMENTS_VIEW_NAME" envDefault:"segment_counts"`
}

// NewConfig loads the environment variables from the specified .env file and returns a Config struct
func GetConfig(options ...string) (*ServiceConfig, error) {
	envPath := ".env" // default .env file path
	if len(options) > 0 {
		envPath = options[0]
	}
	// Load .env file
	if err := godotenv.Load(envPath); err != nil {
		log.Printf("Error loading .env file from %s, using default config: %v", envPath, err)
	}
	envs := &ServiceConfig{}
	err := env.Parse(envs)
	if err != nil {
		return nil, fmt.Errorf("error parsing env: %v", err)
	}
	return envs, err
}
