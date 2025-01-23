package _28239

import (
	"fmt"
	"log"
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

// UserRateLimiter holds rate limits for multiple users.
type UserRateLimiter struct {
	rateLimits map[string]*RateLimiter
	mu         sync.Mutex
}

type FaultyMicroservice struct {
	name        string
	rateLimiter *UserRateLimiter
	userID      string
	failures    []string
}

// NewFaultyMicroservice creates a new microservice with fault injection capabilities.
func NewFaultyMicroservice(name string, rateLimiter *UserRateLimiter, userID string, failures []string) *FaultyMicroservice {
	return &FaultyMicroservice{name: name, rateLimiter: rateLimiter, userID: userID, failures: failures}
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

// CallAPI simulates calling the API with fault injection.
func (ms *FaultyMicroservice) CallAPI(wg *sync.WaitGroup, results chan<- string) {
	defer wg.Done()

	simulateFaults(ms.failures)

	if ms.rateLimiter.Allow(ms.userID) {
		results <- fmt.Sprintf("%s: Request allowed", ms.name)
	} else {
		results <- fmt.Sprintf("%s: Request rate limited", ms.name)
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

func simulateFaults(faults []string) {
	for _, fault := range faults {
		switch fault {
		case "network_partition":
			time.Sleep(200 * time.Millisecond) // Simulate network delay
			log.Println("Simulating network partition")
		case "api_failure":
			time.Sleep(200 * time.Millisecond) // Simulate API being temporarily unavailable
			log.Println("Simulating API failure")
		case "high_latency":
			time.Sleep(500 * time.Millisecond) // Simulate high latency conditions
			log.Println("Simulating high latency")
		}
	}
}

// NewUserRateLimiter initializes a new UserRateLimiter.
func NewUserRateLimiter() *UserRateLimiter {
	return &UserRateLimiter{
		rateLimits: make(map[string]*RateLimiter),
	}
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

// Integration test with fault injection.
func TestFaultInjectionRateLimiting(t *testing.T) {
	rateLimiter := NewUserRateLimiter()
	rateLimiter.SetRateLimit("api_user", 5) // 5 requests per second

	microservices := []*FaultyMicroservice{
		NewFaultyMicroservice("Service1", rateLimiter, "api_user", nil),
		NewFaultyMicroservice("Service2", rateLimiter, "api_user", []string{"network_partition"}),
		NewFaultyMicroservice("Service3", rateLimiter, "api_user", []string{"api_failure"}),
		NewFaultyMicroservice("Service4", rateLimiter, "api_user", []string{"high_latency"}),
	}

	var wg sync.WaitGroup
	results := make(chan string, 20)

	// Simulate 5 requests within a second with fault injection
	wg.Add(5)
	for i := 0; i < 5; i++ {
		go func(ms *FaultyMicroservice) {
			ms.CallAPI(&wg, results)
		}(microservices[i%len(microservices)])
		time.Sleep(200 * time.Millisecond)
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

	if allowedCount < 3 { // Expecting at least 3 requests to pass through due to faults
		t.Errorf("Expected at least 3 allowed requests, but got %d", allowedCount)
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
