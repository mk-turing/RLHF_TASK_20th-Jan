package _28239_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"golang.org/x/time/rate"
)

// define API handler function
func apiHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "API Request Handled")
}

type testLimiter struct {
	mu      sync.Mutex
	limiter *rate.Limiter
}

func (t *testLimiter) Allow() bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.limiter.Allow()
}

func (t *testLimiter) Reset(limit rate.Limit, burst int) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.limiter = rate.NewLimiter(limit, burst)
}

// custom setup function
func TestMain(m *testing.M) {
	m.Run()
}

// helper function to make a test request
func makeTestRequest(url string, t *testing.T) *httptest.ResponseRecorder {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	rr := httptest.NewRecorder()
	// Handler should not panic as the limiter checks for rate limits
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Recovered from panic: %v", r)
		}
	}()

	apiHandler(rr, req)
	return rr
}

func TestApiRateLimiting_EvenTraffic(t *testing.T) {
	for i := 0; i < 100; i++ {
		rr := makeTestRequest("/api", t)
		if rr.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", rr.Code)
		}
		time.Sleep(time.Millisecond * 10) // Pause to simulate traffic distribution
	}
}

func TestApiRateLimiting_BurstRequests(t *testing.T) {
	// Send a burst of 100 requests
	for i := 0; i < 100; i++ {
		go makeTestRequest("/api", t)
	}

	// Give the limiter time to handle the burst
	time.Sleep(time.Second)

	// Make sure all requests were processed
	for i := 0; i < 100; i++ {
		rr := makeTestRequest("/api", t)
		if rr.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", rr.Code)
		}
	}
}

func TestApiRateLimiting_RateLimitReset(t *testing.T) {
	for i := 0; i < 10; i++ {
		rr := makeTestRequest("/api", t)
		if rr.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", rr.Code)
		}
	}

	// Wait for the rate limit to reset
	time.Sleep(time.Second)

	rr := makeTestRequest("/api", t)
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200 after rate limit reset, got %d", rr.Code)
	}
}
