package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// MockServer simulates an external service with configurable delay and failure rate.
func MockServer(responseDelay time.Duration, failRate int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(responseDelay)
		if failRate > 0 && time.Now().UnixNano()%int64(failRate) == 0 {
			http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Mock response"))
	}))
}

// SimulateAPICall makes a request to the mock server and measures response time.
func SimulateAPICall(serverURL string) (int, error) {
	start := time.Now()
	resp, err := http.Get(serverURL)
	elapsed := time.Since(start)

	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	return int(elapsed.Milliseconds()), err
}

// BenchmarkNetworkImpact measures the impact of network conditions on responsiveness.
func BenchmarkNetworkImpact(b *testing.B) {
	responseDelay := 100 * time.Millisecond
	failRate := 5 // Induces a failure every 5th request

	server := MockServer(responseDelay, failRate)
	defer server.Close()

	var totalTime int
	var failureCount int

	for i := 0; i < b.N; i++ {
		responseTime, err := SimulateAPICall(server.URL)
		if err != nil {
			failureCount++
		} else {
			totalTime += responseTime
		}
	}

	averageTime := totalTime / (b.N - failureCount)
	fmt.Printf("Average Response Time: %dms, Failures: %d\n", averageTime, failureCount)
}

func main() {
	// Typically, tests are run using the 'go test' command, but for demonstration purposes:
	benchmark := testing.Benchmark(BenchmarkNetworkImpact)
	fmt.Println(benchmark)
}
