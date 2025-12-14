package retry

import (
	"context"
	"time"
)

type Config struct {
	MaxAttempts int
	Delay       time.Duration
	Timeout     time.Duration
}

func DefaultConfig() *Config {
	return &Config{
		MaxAttempts: 3,
		Delay:       1 * time.Second,
		Timeout:     10 * time.Second,
	}
}

type Func[T any] func(ctx context.Context) (T, error)

func WithRetry[T any](ctx context.Context, config *Config, fn Func[T]) (T, error) {
	var lastErr error
	var result T

	if config == nil {
		config = DefaultConfig()
	}

	for attempt := 0; attempt < config.MaxAttempts; attempt++ {
		attemptCtx, cancel := context.WithTimeout(ctx, config.Timeout)
		result, lastErr = fn(attemptCtx)
		cancel()

		if lastErr == nil {
			return result, nil
		}

		if attempt == config.MaxAttempts-1 {
			break
		}

		select {
		case <-ctx.Done():
			return result, ctx.Err()
		case <-time.After(config.Delay):
		}
	}

	return result, lastErr
}
