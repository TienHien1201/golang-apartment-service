package xqueue

import (
	"context"
	"fmt"
	"sync"
	"time"

	xlogger "thomas.vn/hr_recruitment/pkg/logger"
)

// InMemoryQueue represents an in-memory queue server.
type InMemoryQueue struct {
	logger    *xlogger.Logger
	config    *QueueConfig
	jobs      map[MessageType]Job
	queue     chan Message
	wg        sync.WaitGroup
	mu        sync.RWMutex
	isRunning bool
	stopCh    chan struct{}
}

func NewInMemoryQueue(logger *xlogger.Logger, config *QueueConfig) *InMemoryQueue {
	if config == nil {
		config = &QueueConfig{}
	}
	if config.Workers <= 0 {
		config.Workers = 3
	}
	if config.QueueSize <= 0 {
		config.QueueSize = 1000
	}
	if config.RetryDelay <= 0 {
		config.RetryDelay = 5 * time.Second
	}

	return &InMemoryQueue{
		logger: logger,
		config: config,
		jobs:   make(map[MessageType]Job),
		queue:  make(chan Message, config.QueueSize),
		stopCh: make(chan struct{}),
	}
}

// RegisterJobs registers multiple jobs with the queue server.
func (s *InMemoryQueue) RegisterJobs(jobs []Job) {
	for _, job := range jobs {
		s.RegisterJob(job)
	}
}

// RegisterJob registers a job with the queue server.
func (s *InMemoryQueue) RegisterJob(job Job) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.jobs[job.Type()]; exists {
		s.logger.Error("Job already registered", xlogger.String("jobName", job.Name()))
		return
	}

	s.jobs[job.Type()] = job
	s.logger.Info("Job registered", xlogger.String("jobName", job.Name()))
}

// Start starts the queue server and begins processing messages.
func (s *InMemoryQueue) Start() error {
	s.mu.Lock()
	if s.isRunning {
		s.mu.Unlock()
		return fmt.Errorf("queue server already running")
	}
	s.isRunning = true
	s.mu.Unlock()

	for i := 0; i < s.config.Workers; i++ {
		s.wg.Add(1)
		go s.worker(i)
	}

	s.logger.Info("Queue server started", xlogger.Int("workers", s.config.Workers), xlogger.Int("queueSize", s.config.QueueSize))
	return nil
}

// Stop stops the queue server and waits for all workers to finish processing.
func (s *InMemoryQueue) Stop(ctx context.Context) error {
	s.mu.Lock()
	if !s.isRunning {
		s.mu.Unlock()
		return nil
	}
	s.isRunning = false
	close(s.stopCh)
	s.mu.Unlock()

	doneCh := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(doneCh)
	}()

	select {
	case <-ctx.Done():
		return fmt.Errorf("timeout waiting for queue workers to stop: %w", ctx.Err())
	case <-doneCh:
		s.logger.Info("Queue server stopped gracefully")
		return nil
	}
}

// Enqueue adds a message to the queue.
func (s *InMemoryQueue) Enqueue(ctx context.Context, msgType MessageType, payload interface{}) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if !s.isRunning {
		return fmt.Errorf("queue server not running")
	}

	if _, exists := s.jobs[msgType]; !exists {
		return fmt.Errorf("no job registered for message type: %s", msgType)
	}

	msg := Message{
		ID:        fmt.Sprintf("%d", time.Now().UnixNano()),
		Type:      msgType,
		Payload:   payload,
		Timestamp: time.Now(),
	}

	select {
	case <-ctx.Done():
		return fmt.Errorf("context cancelled: %w", ctx.Err())
	case s.queue <- msg:
		return nil
	default:
		return fmt.Errorf("queue is full")
	}
}

// PublishMessage publishes a message to the queue.
func (s *InMemoryQueue) PublishMessage(ctx context.Context, msgType MessageType, payload interface{}) error {
	return s.Enqueue(ctx, msgType, payload)
}

func (s *InMemoryQueue) worker(id int) {
	defer s.wg.Done()

	s.logger.Info("Queue worker started", xlogger.Int("workerID", id))

	for {
		select {
		case <-s.stopCh:
			s.logger.Info("Queue worker stopping", xlogger.Int("workerID", id))
			return
		case msg := <-s.queue:
			s.processMessage(msg)
		}
	}
}

func (s *InMemoryQueue) processMessage(msg Message) {
	job, exists := s.jobs[msg.Type]
	if !exists {
		s.logger.Error("No job found for message type",
			xlogger.String("messageID", msg.ID),
			xlogger.String("messageType", string(msg.Type)))
		return
	}

	s.logger.Info("Processing message",
		xlogger.String("messageID", msg.ID),
		xlogger.String("jobName", job.Name()))

	err := job.Handle(context.Background(), msg.Payload)
	if err != nil {
		s.handleProcessingError(msg, job, err)
	} else {
		s.logger.Info("Message processed successfully",
			xlogger.String("messageID", msg.ID),
			xlogger.String("jobName", job.Name()))
	}
}

func (s *InMemoryQueue) handleProcessingError(msg Message, job Job, err error) {
	s.logger.Error("Error processing message",
		xlogger.String("messageID", msg.ID),
		xlogger.String("jobName", job.Name()),
		xlogger.Error(err),
		xlogger.Int("attempt", msg.Attempts+1))

	// retry if not reached max attempts
	if msg.Attempts < s.config.RetryLimit {
		msg.Attempts++

		// delay before retry
		time.Sleep(s.config.RetryDelay)

		// re-push message to queue
		select {
		case s.queue <- msg:
			s.logger.Info("Requeued message for retry",
				xlogger.String("messageID", msg.ID),
				xlogger.String("jobName", job.Name()),
				xlogger.Int("attempt", msg.Attempts))
		default:
			s.logger.Error("Failed to requeue message, queue is full",
				xlogger.String("messageID", msg.ID),
				xlogger.String("jobName", job.Name()))
		}
	} else {
		s.logger.Error("Message processing failed after max retries",
			xlogger.String("messageID", msg.ID),
			xlogger.String("jobName", job.Name()),
			xlogger.Int("maxRetries", s.config.RetryLimit))
	}
}
