package ws

import (
	"sync"
	"time"
)

// rateLimiter is a simple token-bucket rate limiter for per-client message throttling.
// It is goroutine-safe and has zero external dependencies.
type rateLimiter struct {
	mu       sync.Mutex
	tokens   float64 // current available tokens
	maxToken float64 // bucket capacity
	refill   float64 // tokens added per nanosecond
	last     time.Time
}

// newRateLimiter creates a limiter that allows ratePerPeriod events in each period.
// Example: newRateLimiter(10, time.Second) → max 10 messages per second.
func newRateLimiter(ratePerPeriod int, period time.Duration) *rateLimiter {
	return &rateLimiter{
		tokens:   float64(ratePerPeriod),
		maxToken: float64(ratePerPeriod),
		refill:   float64(ratePerPeriod) / float64(period.Nanoseconds()),
		last:     time.Now(),
	}
}

// Allow returns true if the request is within the rate limit and consumes one token.
func (r *rateLimiter) Allow() bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	elapsed := float64(now.Sub(r.last).Nanoseconds())
	r.last = now

	r.tokens += elapsed * r.refill
	if r.tokens > r.maxToken {
		r.tokens = r.maxToken
	}

	if r.tokens < 1 {
		return false
	}
	r.tokens--
	return true
}
