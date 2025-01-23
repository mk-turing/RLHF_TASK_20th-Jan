package _28239

import (
	"os"
	"sync"
	"testing"
	"time"
)

// RateLimiter is a simple token bucket rate limiter.
type RateLimiter struct {
	rate      int
	bucket    float64
	lastCheck time.Time
	mu        sync.Mutex
}

// NewRateLimiter creates a new RateLimiter with specified tokens per second.
func NewRateLimiter(rate int) *RateLimiter {
	return &RateLimiter{
		rate:      rate,
		bucket:    float64(rate),
		lastCheck: time.Now(),
	}
}

// Allow checks if a request is allowed under the rate limit.
func (rl *RateLimiter) Allow() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(rl.lastCheck).Seconds()
	rl.bucket += elapsed * float64(rl.rate)

	// Ensure bucket doesn't exceed its capacity
	if rl.bucket > float64(rl.rate) {
		rl.bucket = float64(rl.rate)
	}
	rl.lastCheck = now

	if rl.bucket >= 1 {
		rl.bucket--
		return true
	}
	return false
}

// Setup and teardown processes for testing.
func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func setup() {
	// Custom setup actions for tests
}

func teardown() {
	// Custom teardown actions
}

// TestRateLimiter tests the RateLimiter's basic functionality.
func TestRateLimiter(t *testing.T) {
	rl := NewRateLimiter(2) // 2 tokens per second

	tests := []struct {
		name    string
		wait    time.Duration
		allowed bool
	}{
		{"Immediate request", 0, true},
		{"Immediate next request", 0, true},
		{"Immediate third request fails", 0, false},
		{"After waiting 1 second", 1 * time.Second, true},
		{"After waiting 2 seconds", 2 * time.Second, true},
		{"Bursts of requests", 500 * time.Millisecond, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			time.Sleep(tt.wait)
			if rl.Allow() != tt.allowed {
				t.Errorf("expected allowed = %v, got = %v", tt.allowed, rl.Allow())
			}
		})
	}
}

// TestBurstyTraffic simulates bursty traffic scenarios.
func TestBurstyTraffic(t *testing.T) {
	rl := NewRateLimiter(5) // 5 tokens per second

	for i := 0; i < 5; i++ {
		if !rl.Allow() {
			t.Errorf("expected request %d to be allowed", i)
		}
	}
	if rl.Allow() {
		t.Error("expected request 6 to be rejected due to rate limit")
	}

	time.Sleep(1 * time.Second)

	if !rl.Allow() {
		t.Error("expected request to be allowed after 1 second")
	}
}

// TestUnevenDistribution tests handling uneven traffic.
func TestUnevenDistribution(t *testing.T) {
	rl := NewRateLimiter(3) // 3 tokens per second

	for i := 0; i < 3; i++ {
		if !rl.Allow() {
			t.Errorf("expected request %d to be allowed", i)
		}
	}

	time.Sleep(300 * time.Millisecond)

	if rl.Allow() {
		t.Error("expected request to be rejected due to partial refill")
	}

	time.Sleep(700 * time.Millisecond)

	if !rl.Allow() {
		t.Error("expected request to be allowed after full second")
	}
}

// TestRateLimitReset simulates handling of rate-limit resets.
func TestRateLimitReset(t *testing.T) {
	rl := NewRateLimiter(1) // 1 token per second

	if !rl.Allow() {
		t.Error("expected initial request to be allowed")
	}

	time.Sleep(1 * time.Second)

	if !rl.Allow() {
		t.Error("expected request to be allowed after reset period")
	}
}
