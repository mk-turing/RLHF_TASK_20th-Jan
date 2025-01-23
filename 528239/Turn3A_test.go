package _28239

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

// RateLimiter is a token bucket rate limiter.
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

// Allow checks if a request is allowed under the rate limit.
func (rl *RateLimiter) Allow() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(rl.lastCheck).Seconds()
	// Refill tokens based on elapsed time
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
		rateLimiter.mu.Lock()
		defer rateLimiter.mu.Unlock()
		rateLimiter.rate = float64(rate)
		rateLimiter.capacity = float64(rate)
		if rateLimiter.bucket > rateLimiter.capacity {
			rateLimiter.bucket = rateLimiter.capacity
		}
	} else {
		url.rateLimits[userID] = NewRateLimiter(rate)
	}
}

// Allow checks if a user is allowed to make a request.
func (url *UserRateLimiter) Allow(userID string) bool {
	url.mu.Lock()
	rateLimiter, exists := url.rateLimits[userID]
	url.mu.Unlock()

	if !exists {
		return false
	}
	return rateLimiter.Allow()
}

// Microservice represents a service interacting with the API.
type Microservice struct {
	name        string
	rateLimiter *UserRateLimiter
	userID      string
}

// NewMicroservice creates a new microservice instance.
func NewMicroservice(name string, rateLimiter *UserRateLimiter, userID string) *Microservice {
	return &Microservice{name: name, rateLimiter: rateLimiter, userID: userID}
}

// CallAPI simulates calling the API with rate limiting.
func (ms *Microservice) CallAPI(wg *sync.WaitGroup, results chan<- string) {
	defer wg.Done()

	if ms.rateLimiter.Allow(ms.userID) {
		results <- fmt.Sprintf("%s: Request allowed", ms.name)
	} else {
		results <- fmt.Sprintf("%s: Request rate limited", ms.name)
	}
}

// Integration test for distributed microservices architecture.
func TestMicroservicesRateLimiting(t *testing.T) {
	rateLimiter := NewUserRateLimiter()
	rateLimiter.SetRateLimit("api_user", 5) // 5 requests per second

	microservices := []*Microservice{
		NewMicroservice("Service1", rateLimiter, "api_user"),
		NewMicroservice("Service2", rateLimiter, "api_user"),
		NewMicroservice("Service3", rateLimiter, "api_user"),
	}

	var wg sync.WaitGroup
	results := make(chan string, 15)

	// Simulate 5 requests within a second
	wg.Add(5)
	for i := 0; i < 5; i++ {
		go func(ms *Microservice) {
			ms.CallAPI(&wg, results)
		}(microservices[i%len(microservices)])
		time.Sleep(200 * time.Millisecond) // Spread out requests slightly
	}

	wg.Wait()
	close(results)

	allowedCount := 0
	for result := range results {
		t.Log(result)
		if result[len(result)-7:] == "allowed" {
			allowedCount++
		}
	}

	if allowedCount != 5 {
		t.Errorf("Expected 5 allowed requests, but got %d", allowedCount)
	}
}

func TestMain(m *testing.M) {
	setup()
	m.Run()
	teardown()
}

func setup() {
	// Custom setup actions for tests
}

func teardown() {
	// Custom teardown actions
}
