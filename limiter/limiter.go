package limiter

import (
	"sync"
	"time"
)

// RateLimiter represents a single token bucket for one client.
type RateLimiter struct {
	mu         sync.Mutex
	capacity   int
	tokens     int
	refillRate time.Duration
	lastRefill time.Time
}

// newBucket initializes a single token bucket.
func newBucket(capacity int, refillEvery time.Duration) *RateLimiter {
	return &RateLimiter{
		capacity:   capacity,
		tokens:     capacity,
		refillRate: refillEvery,
		lastRefill: time.Now(),
	}
}

// Allow checks if the single bucket has an available token.
func (rl *RateLimiter) Allow() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(rl.lastRefill)
	tokensToAdd := int(elapsed / rl.refillRate)

	if tokensToAdd > 0 {
		rl.tokens += tokensToAdd
		if rl.tokens > rl.capacity {
			rl.tokens = rl.capacity
		}
		rl.lastRefill = rl.lastRefill.Add(time.Duration(tokensToAdd) * rl.refillRate)
	}

	if rl.tokens > 0 {
		rl.tokens--
		return true
	}

	return false
}

// IPRateLimiter manages a dedicated token bucket for each IP address.
type IPRateLimiter struct {
	mu          sync.RWMutex
	ips         map[string]*RateLimiter
	capacity    int
	refillEvery time.Duration
}

// NewIPRateLimiter creates an IP-based rate limit manager.
func NewIPRateLimiter(capacity int, refillEvery time.Duration) *IPRateLimiter {
	return &IPRateLimiter{
		ips:         make(map[string]*RateLimiter),
		capacity:    capacity,
		refillEvery: refillEvery,
	}
}

// getLimiter retrieves or initializes the bucket for a given IP.
func (i *IPRateLimiter) getLimiter(ip string) *RateLimiter {
	// 1. Fast read-lock to check if IP already exists
	i.mu.RLock()
	limiter, exists := i.ips[ip]
	i.mu.RUnlock()

	if exists {
		return limiter
	}

	// 2. IP does not exist: acquire write-lock to add it
	i.mu.Lock()
	defer i.mu.Unlock()

	// Double-check existence in case another goroutine created it in the meantime
	limiter, exists = i.ips[ip]
	if !exists {
		limiter = newBucket(i.capacity, i.refillEvery)
		i.ips[ip] = limiter
	}

	return limiter
}

// Allow checks if the given IP address is allowed to proceed.
func (i *IPRateLimiter) Allow(ip string) bool {
	limiter := i.getLimiter(ip)
	return limiter.Allow()
}
