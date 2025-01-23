package _28239

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

// Initialize global variables and mutex for rate limiter usage.
var (
	rateLimiter   = NewUserRateLimiter()
	rateLimitMu sync.Mutex
)

// simulateDelay introduces a random delay between 0 and 500ms to simulate delays in service responses.
func simulateDelay() {
	time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
}

// Ensure rate limiting for service handlers.
func rateLimitedHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate delay for edge cases.
		simulateDelay()

		userID := r.FormValue("user_id")
		if userID == "" {
			http.Error(w, "user_id is required", http.StatusBadRequest)
			return
		}

		rateLimitMu.Lock()
		defer rateLimitMu.Unlock()

		if !rateLimiter.Allow(userID) {
			http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
			return
		}
		// Pass the request to the underlying handler if not rate limited.
		h.ServeHTTP(w, r)
	})
}

// Microservice 1 handler.
func microservice1Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello from Microservice 1!")
}

// Microservice 2 handler.
func microservice2Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello from Microservice 2!")
}

// Combined microservice endpoints for distributed setup.
func microservicesHandler() http.Handler {
	m := http.NewServeMux()
	m.Handle("/microservice1", rateLimitedHandler(http.HandlerFunc(microservice1Handler)))
	m.Handle("/microservice2", rateLimitedHandler(http.HandlerFunc(microservice2Handler)))
	return m
}

// Integration Test
func TestDistributedRateLimiting(t *testing.T) {
	// Reset rate limiter before the test.
	rateLimiter = NewUserRateLimiter()

	// Set up distributed microservices server.
	server := httptest.NewServer(microservicesHandler())
	defer server.Close()

	baseURL := server.URL

	// Simulate normal usage from multiple microservices, ensuring rates are limited consistently.
	testUserID := "testuser"
	rateLimiter.SetRateLimit(testUserID, 2) // Allow 2 requests per second

	for i := 0; i < 5; i++ {
		resp, err := http.Get(fmt.Sprintf("%s/microservice1?user_id=%s", baseURL, testUserID))
		if err != nil {
			t.Errorf("Request failed: %v", err)
			continue
		}

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Got status code %d, expected %d", resp.StatusCode, http.StatusOK)
		}
		resp.Body.Close()
		time.Sleep(500 * time.Millisecond) // Wait to avoid overwhelming the service during tests
	}

	// Simulate micro-service1 being down for a short period.
	// We should still be rate-limited using microservice2.
	rateLimiter.SetRateLimit(testUserID, 1)

	if _, err := http.Get(fmt.Sprintf("%s/microservice1?user_id=%s", baseURL, testUserID)); err == nil {
		t.Error("Expecting an error when microservice1 is down.")
	}

	resp, err := http.Get(fmt.Sprintf("%s/microservice2?user_id=%s", baseURL, testUserID))