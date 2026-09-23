package limiter

import (
	"sync"
	"testing"
	"time"
)

// TestBurstLimit ensures a bucket consumes up to capacity, and the next call fails.
func TestBurstLimit(t *testing.T) {
	capacity := 3
	bucket := newBucket(capacity, 1*time.Second)

	for i := 0; i < capacity; i++ {
		if !bucket.Allow() {
			t.Fatalf("expected request %d to be allowed", i+1)
		}
	}

	if bucket.Allow() {
		t.Fatal("expected 4th request to be blocked by rate limiter")
	}
}

// TestRefill ensures tokens regenerate correctly over time.
func TestRefill(t *testing.T) {
	bucket := newBucket(1, 100*time.Millisecond)

	if !bucket.Allow() {
		t.Fatal("expected first request to pass")
	}

	if bucket.Allow() {
		t.Fatal("expected immediate second request to fail")
	}

	time.Sleep(110 * time.Millisecond)

	if !bucket.Allow() {
		t.Fatal("expected request to pass after refill interval")
	}
}

// TestConcurrentAccess verifies mutex thread safety under concurrent goroutines.
func TestConcurrentAccess(t *testing.T) {
	bucket := newBucket(10, 100*time.Millisecond)
	var wg sync.WaitGroup

	numGoroutines := 50
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			_ = bucket.Allow()
		}()
	}

	wg.Wait()
}

// TestIPRateLimiterIsolation verifies one IP's burst does not impact another IP.
func TestIPRateLimiterIsolation(t *testing.T) {
	ipLimiter := NewIPRateLimiter(2, 1*time.Second)

	ipA := "192.168.1.1"
	ipB := "192.168.1.2"

	// Exhaust tokens for IP A
	if !ipLimiter.Allow(ipA) {
		t.Fatal("expected IP A request 1 to pass")
	}
	if !ipLimiter.Allow(ipA) {
		t.Fatal("expected IP A request 2 to pass")
	}
	if ipLimiter.Allow(ipA) {
		t.Fatal("expected IP A request 3 to be blocked")
	}

	// IP B should still have all tokens available
	if !ipLimiter.Allow(ipB) {
		t.Fatal("expected IP B request 1 to pass despite IP A being blocked")
	}
	if !ipLimiter.Allow(ipB) {
		t.Fatal("expected IP B request 2 to pass")
	}
}
