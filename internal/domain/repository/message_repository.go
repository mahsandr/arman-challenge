package repository

import (
	"context"
)

// Repository Interfaces
type MessageBroker interface {
	Produce(ctx context.Context, msg []byte) error
	Consume(ctx context.Context) (<-chan [][]byte, error)
	Close()
}
