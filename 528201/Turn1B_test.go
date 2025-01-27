package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// The actual realTimeHandler from the real application
func realTimeHandler(w http.ResponseWriter, r *http.Request) {
	// ... (implementation unchanged)
}

func BenchmarkLatency(b *testing.B) {
	// Create a test server with the realTimeHandler
	ts := httptest.NewServer(http.HandlerFunc(realTimeHandler))
	defer ts.Close()

	client := &http.Client{}
	for i := 0; i < b.N; i++ {
		// Start measuring the duration
		start := time.Now()
		req, err := http.NewRequest("GET", ts.URL, nil)
		if err != nil {
			b.Fatalf("Failed to create request: %v", err)
		}

		resp, err := client.Do(req)
		if err != nil {
			b.Fatalf("Failed to make request: %v", err)
		}
		_ = resp.Body.Close()

		// End the timing and record the result
		b.SetBytes(int64(len(resp.Header)))
		b.ReportMetric(float64(time.Since(start).Milliseconds()), "ms")
	}
}

func BenchmarkThroughput(b *testing.B) {
	// Create a test server with the realTimeHandler
	ts := httptest.NewServer(http.HandlerFunc(realTimeHandler))
	defer ts.Close()

	client := &http.Client{}
	for i := 0; i < b.N; i++ {
		req, err := http.NewRequest("GET", ts.URL, nil)
		if err != nil {
			b.Fatalf("Failed to create request: %v", err)
		}

		resp, err := client.Do(req)
		if err != nil {
			b.Fatalf("Failed to make request: %v", err)
		}
		_ = resp.Body.Close()
	}
}
