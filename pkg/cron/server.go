package xcron

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/robfig/cron/v3"

	xlogger "thomas.vn/hr_recruitment/pkg/logger"
)

type Server struct {
	cron     *cron.Cron
	logger   *xlogger.Logger
	jobs     []Job
	jobIDs   map[string]cron.EntryID
	stopChan chan struct{}
	wg       sync.WaitGroup
}

func NewCronServer(logger *xlogger.Logger, jobs []Job) *Server {
	return &Server{
		logger:   logger,
		jobs:     jobs,
		jobIDs:   make(map[string]cron.EntryID),
		stopChan: make(chan struct{}),
	}
}

// Start starts the cron server and schedules all jobs.
func (s *Server) Start() error {
	cronOpts := []cron.Option{
		cron.WithSeconds(),
	}

	s.cron = cron.New(cronOpts...)

	// Register all jobs
	for _, job := range s.jobs {
		if err := s.registerJob(job); err != nil {
			return err
		}
	}

	s.cron.Start()

	s.logger.Info("Cron server started")
	return nil
}

// Stop stops the cron server and waits for all jobs to complete.
func (s *Server) Stop(ctx context.Context) error {
	// Stop accepting new jobs
	s.cron.Stop()

	// Wait for running jobs with context timeout
	doneChan := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(doneChan)
	}()

	select {
	case <-ctx.Done():
		return fmt.Errorf("timeout waiting for jobs to complete: %w", ctx.Err())
	case <-doneChan:
	}

	s.logger.Info("Cron server stopped gracefully")
	return nil
}

// registerJob adds a single job to the cron scheduler.
func (s *Server) registerJob(job Job) error {
	jobName := job.Name()

	if !job.Enabled() {
		s.logger.Info("Skipping", xlogger.String("job", jobName))
		return nil
	}

	schedule := job.Schedule()
	if schedule == "" {
		return fmt.Errorf("schedule not configured for job %s", jobName)
	}

	// Create a wrapper to handle context and logging
	wrapper := createJobWrapper(s, job, 4*time.Hour)

	id, err := s.cron.AddFunc(schedule, wrapper)
	if err != nil {
		return fmt.Errorf("failed to add job %s: %w", jobName, err)
	}
	s.jobIDs[jobName] = id
	return nil
}

// createJobWrapper creates a function wrapper for the job execution.
func createJobWrapper(s *Server, job Job, timeout time.Duration) func() {
	return func() {
		s.wg.Add(1)
		defer s.wg.Done()

		// Create context with timeout
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		s.logger.Info("Starting", xlogger.String("job", job.Name()))
		if err := job.Execute(ctx); err != nil {
			s.logger.Error("Job failed", xlogger.String("job", job.Name()), xlogger.Error(err))
		}
		s.logger.Info("Completed", xlogger.String("job", job.Name()))
	}
}
