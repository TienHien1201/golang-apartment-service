package service

import "context"

// QueueService defines the async messaging contract for the domain and usecase layers.
// The msgType parameter uses string (via the type-alias consts.MessageType) so that
// infrastructure implementations (e.g. InMemoryQueue) satisfy this interface directly
// after pkg/queue.MessageType is declared as a string type alias.
type QueueService interface {
	PublishMessage(ctx context.Context, msgType string, payload interface{}) error
}
