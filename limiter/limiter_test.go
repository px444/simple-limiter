package limiter

import (
	"sync"
	"testing"
	"time"
)

// TestBurstLimit ensures we can consume up to capacity, and the next call fails.
func TestBurstLimit(t *testing.T) {
	capacity := 3
	rl := New(capacity, 1*time.Second)

	// Consume all 3 burst tokens
	for i := 0; i < capacity; i++ {
		if !rl.Allow() {
			t.Fatalf("expected request %d to be allowed", i+1)
		}
	}

	// 4th request must be rejected
	if rl.Allow() {
		t.Fatal("expected 4th request to be blocked by rate limiter")
	}
}

// TestRefill ensures tokens regenerate correctly over time.
func TestRefill(t *testing.T) {
	rl := New(1, 100*time.Millisecond)

	// Consume the single token
	if !rl.Allow() {
		t.Fatal("expected first request to pass")
	}

	// Immediate next request should fail
	if rl.Allow() {
		t.Fatal("expected immediate second request to fail")
	}

	// Wait for refill interval
	time.Sleep(110 * time.Millisecond)

	// Token should now be restored
	if !rl.Allow() {
		t.Fatal("expected request to pass after refill interval")
	}
}

// TestConcurrentAccess launches multiple goroutines to test mutex thread-safety.
func TestConcurrentAccess(t *testing.T) {
	rl := New(10, 100*time.Millisecond)
	var wg sync.WaitGroup

	numGoroutines := 50
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			// Calling Allow() concurrently must not panic or cause a data race
			_ = rl.Allow()
		}()
	}

	wg.Wait()
}
