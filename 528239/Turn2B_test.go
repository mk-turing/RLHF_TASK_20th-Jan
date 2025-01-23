package _28239

import (
	"fmt"
	"math/rand"
	"os"
	"sync"
	"testing"
	"time"
)

// UserRateLimiter manages rate limits for individual users.
type UserRateLimiter struct {
	mu      sync.Mutex
	limits  map[string]*RateLimiter
	updates chan update
}

type update struct {
	user   string
	rate   int
	action string // "add" or "remove"
}

// NewUserRateLimiter creates a new UserRateLimiter.
func NewUserRateLimiter() *UserRateLimiter {
	return &UserRateLimiter{
		limits:  make(map[string]*RateLimiter),
		updates: make(chan update),
	}
}

// Run starts the background goroutine to handle rate-limit updates.
func (url *UserRateLimiter) Run() {
	go func() {
		for update := range url.updates {
			url.mu.Lock()
			switch update.action {
			case "add":
				url.limits[update.user] = NewRateLimiter(update.rate)
			case "remove":
				delete(url.limits, update.user)
			}
			url.mu.Unlock()
		}
	}()
}

// Allow checks if a request for a specific user is allowed under their rate limit.
func (url *UserRateLimiter) Allow(user string) bool {
	url.mu.Lock()
	defer url.mu.Unlock()

	rl, ok := url.limits[user]
	if !ok {
		// User not found, allow by default (or you can block)
		return true
	}
	return rl.Allow()
}

// UpdateRate updates the rate limit for a specific user.
func (url *UserRateLimiter) UpdateRate(user string, rate int) {
	url.updates <- update{user, rate, "add"}
}

// RemoveUser removes a user from the rate-limiting mechanism.
func (url *UserRateLimiter) RemoveUser(user string) {
	url.updates <- update{user, 0, "remove"}
}

// SimulateConcurrency tests the rate limiter with multiple concurrent users and services.
func TestSimulateConcurrency(t *testing.T) {
	numUsers := 10
	numServices := 5
	rateLimitPerUser := 5 // 5 requests per second per user
	totalRequests := 1000

	url := NewUserRateLimiter()
	url.Run()

	// Add users with rate limits
	for i := 1; i <= numUsers; i++ {
		user := fmt.Sprintf("user%d", i)
		url.UpdateRate(user, rateLimitPerUser)
	}

	// Create services with random request rates
	var wg sync.WaitGroup
	for i := 1; i <= numServices; i++ {
		service := fmt.Sprintf("service%d", i)
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < totalRequests/numServices; j++ {
				// Simulate random user
				user := fmt.Sprintf("user%d", rand.Intn(numUsers)+1)
				if !url.Allow(user) {
					t.Errorf("service %s: user %s exceeded rate limit", service, user)
				}
				time.Sleep(time.Duration(rand.Intn(200)) * time.Millisecond)
			}
		}()
	}

	wg.Wait()
}

// TestDynamicRateLimit adjusts rate limits during runtime and verifies the changes.
func TestDynamicRateLimit(t *testing.T) {
	url := NewUserRateLimiter()
	url.Run()

	user := "testuser"
	initialRate := 2
	updatedRate := 10

	url.UpdateRate(user, initialRate)