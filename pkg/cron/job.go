package xcron

import (
	"context"
)

type Job interface {
	// Name returns the unique identifier of the job.
	Name() string

	// Schedule returns the cron schedule.
	Schedule() string

	// Enabled returns whether the job is enabled.
	Enabled() bool

	// Execute runs the job with context.
	Execute(ctx context.Context) error
}
