package xqueue

import (
	"context"
	"time"
)

type QueueService interface {
	PublishMessage(ctx context.Context, msgType MessageType, payload interface{}) error
}

// MessageType represents the type of message in the queue
type MessageType string

// MessageHandler is a function that processes a message
type MessageHandler func(context.Context, interface{}) error

// QueueConfig contains the configuration for the queue
type QueueConfig struct {
	Workers    int           // number of workers
	QueueSize  int           // size of the queue
	RetryLimit int           // number of maximum retries
	RetryDelay time.Duration // time delay between retries
}

// Message represents a message in the queue
type Message struct {
	ID        string
	Type      MessageType
	Payload   interface{}
	Attempts  int
	Timestamp time.Time
}
