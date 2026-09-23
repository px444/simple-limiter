package limiter

import (
	"sync"
	"time"
)

// RateLimiter controls how frequently events can happen.
type RateLimiter struct {
	mu         sync.Mutex    // Protects fields from concurrent read/writes
	capacity   int           // Max tokens the bucket can hold
	tokens     int           // Current tokens available
	refillRate time.Duration // Time required to add 1 token
	lastRefill time.Time     // Last time tokens were calculated
}

// New creates and initializes a RateLimiter.
// - capacity: max burst allowed
// - refillEvery: interval to restore 1 token (e.g., 500ms means 2 tokens/sec)
func New(capacity int, refillEvery time.Duration) *RateLimiter {
	return &RateLimiter{
		capacity:   capacity,
		tokens:     capacity, // start with a full bucket
		refillRate: refillEvery,
		lastRefill: time.Now(),
	}
}

// Allow checks if a request is allowed to proceed.
func (rl *RateLimiter) Allow() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	// 1. Calculate tokens earned based on elapsed time
	now := time.Now()
	elapsed := now.Sub(rl.lastRefill)
	tokensToAdd := int(elapsed / rl.refillRate)

	if tokensToAdd > 0 {
		rl.tokens += tokensToAdd
		if rl.tokens > rl.capacity {
			rl.tokens = rl.capacity
		}
		// Advance lastRefill by the exact duration consumed
		rl.lastRefill = rl.lastRefill.Add(time.Duration(tokensToAdd) * rl.refillRate)
	}

	// 2. Consume a token if available
	if rl.tokens > 0 {
		rl.tokens--
		return true
	}

	return false
}
