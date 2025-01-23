package _28239

import (
	"sync"
	"testing"
	"time"
)

// RateLimiter is a simple token bucket rate limiter.
type RateLimiter struct {
	rate      float64   // Tokens added per second
	bucket    float64   // Current tokens in the bucket
	capacity  float64   // Maximum bucket capacity
	lastCheck time.Time // Last time tokens were added
	mu        sync.Mutex
}

// NewRateLimiter creates a new RateLimiter with a specified rate and capacity.
func NewRateLimiter(rate int) *RateLimiter {
	return &RateLimiter{
		rate:      float64(rate),
		bucket:    float64(rate),
		capacity:  float64(rate),
		lastCheck: time.Now(),
	}
}

// SetRate dynamically adjusts the rate and capacity of the rate limiter.
func (rl *RateLimiter) SetRate(rate int) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(rl.lastCheck).Seconds()

	// Refill bucket based on the previous rate before adjusting the rate
	rl.bucket += elapsed * rl.rate
	if rl.bucket > rl.capacity {
		rl.bucket = rl.capacity
	}

	// Adjust rate and capacity
	rl.rate = float64(rate)
	rl.capacity = float64(rate)
	if rl.bucket > rl.capacity {
		rl.bucket = rl.capacity
	}

	rl.lastCheck = now
}

// Allow checks if a request is allowed under the rate limit.
func (rl *RateLimiter) Allow() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(rl.lastCheck).Seconds()
	rl.bucket += elapsed * rl.rate
	if rl.bucket > rl.capacity {
		rl.bucket = rl.capacity
	}
	rl.lastCheck = now

	if rl.bucket >= 1 {
		rl.bucket--
		return true
	}
	return false
}

// UserRateLimiter holds rate limits for multiple users.
type UserRateLimiter struct {
	rateLimits map[string]*RateLimiter
	mu         sync.Mutex
}

// NewUserRateLimiter initializes a new UserRateLimiter.
func NewUserRateLimiter() *UserRateLimiter {
	return &UserRateLimiter{
		rateLimits: make(map[string]*RateLimiter),
	}
}

// SetRateLimit sets or updates a rate limit for a specific user.
func (url *UserRateLimiter) SetRateLimit(userID string, rate int) {
	url.mu.Lock()
	defer url.mu.Unlock()

	if rateLimiter, exists := url.rateLimits[userID]; exists {
		rateLimiter.SetRate(rate)
	} else {
		url.rateLimits[userID] = NewRateLimiter(rate)
	}
}

// Allow checks if a user is allowed to make a request.
func (url *UserRateLimiter) Allow(userID string) bool {
	url.mu.Lock()
	defer url.mu.Unlock()

	if rateLimiter, exists := url.rateLimits[userID]; exists {
		return rateLimiter.Allow()
	}
	return false // Default behavior if no rate limit is set for the user
}

// Tests
func TestDynamicRateLimitAdjustment(t *testing.T) {
	ul := NewUserRateLimiter()
	ul.SetRateLimit("user1", 1)

	if !ul.Allow("user1") {
		t.Error("user1: first request should be allowed")
	}
	if ul.Allow("user1") {
		t.Error("user1: second request should be rejected initially")
	}

	// Dynamically adjust rate limit
	ul.SetRateLimit("user1", 2)

	time.Sleep(500 * time.Millisecond) // half-second wait

	// After rate adjustment, partial refill should allow one token for "user1"
	if ul.Allow("user1") {
		t.Log("user1: partial refill allows one request")
	} else {
		t.Error("user1: partial refill should allow one request")
	}

	time.Sleep(500 * time.Millisecond) // another half-second wait to complete 1 second

	// Full refill should now allow another request
	if !ul.Allow("user1") {
		t.Error("user1: should be allowed after rate limit adjustment and full time passage")
	}
}

func TestMultipleUsers(t *testing.T) {
	ul := NewUserRateLimiter()
	ul.SetRateLimit("user1", 3)
	ul.SetRateLimit("user2", 1)

	if !ul.Allow("user1") || !ul.Allow("user1") || !ul.Allow("user1") {
		t.Error("user1: should allow 3 requests immediately")
	}
	if ul.Allow("user1") {
		t.Error("user1: fourth request should be rejected")
	}

	if !ul.Allow("user2") {
		t.Error("user2: first request should be allowed")
	}
	if ul.Allow("user2") {
		t.Error("user2: second request should be rejected")
	}
}

func TestConcurrency(t *testing.T) {
	ul := NewUserRateLimiter()
	ul.SetRateLimit("user1", 5)

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 5; i++ {
			if !ul.Allow("user1") {
				t.Errorf("user1: request %d should have been allowed", i)
			}
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 5; i++ {
			if ul.Allow("user1") {
				t.Logf("user1: concurrent request %d was allowed after refill", i)
			}
		}
	}()

	wg.Wait()
}
