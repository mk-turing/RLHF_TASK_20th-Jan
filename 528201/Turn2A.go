package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"sync"
	"time"
)

type Metrics struct {
	TotalRequests  int           `json:"total_requests"`
	FailedRequests int           `json:"failed_requests"`
	ErrorRate      float64       `json:"error_rate"`
	CPUUsage       float64       `json:"cpu_usage"`
	MemoryUsage    uint64        `json:"memory_usage"`
	Goroutines     int           `json:"goroutines"`
	Concurrency    int           `json:"concurrency"`
	ElapsedTime    time.Duration `json:"elapsed_time"`
}

func stressTest(url string, maxConcurrency int, step int, duration time.Duration) []Metrics {
	var results []Metrics
	for concurrency := step; concurrency <= maxConcurrency; concurrency += step {
		log.Printf("Starting stress test with concurrency: %d", concurrency)
		metrics := runTest(url, concurrency, duration)
		results = append(results, metrics)
	}
	return results
}

func runTest(url string, concurrency int, duration time.Duration) Metrics {
	var wg sync.WaitGroup
	metrics := Metrics{Concurrency: concurrency}
	var mu sync.Mutex
	startTime := time.Now()

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				if time.Since(startTime) > duration {
					return
				}
				resp, err := http.Get(url)
				mu.Lock()
				metrics.TotalRequests++
				if err != nil || resp.StatusCode != 200 {
					metrics.FailedRequests++
				}
				mu.Unlock()
				if resp != nil {
					_ = resp.Body.Close()
				}
			}
		}()
	}

	wg.Wait()

	// Record resource metrics
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	metrics.MemoryUsage = memStats.Alloc
	metrics.CPUUsage = float64(runtime.NumCPU())
	metrics.Goroutines = runtime.NumGoroutine()
	metrics.ErrorRate = float64(metrics.FailedRequests) / float64(metrics.TotalRequests) * 100
	metrics.ElapsedTime = time.Since(startTime)

	return metrics
}

func saveMetricsToFile(metrics []Metrics, filename string) {
	file, err := os.Create(filename)
	if err != nil {
		log.Fatalf("Failed to create metrics file: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(metrics); err != nil {
		log.Fatalf("Failed to write metrics to file: %v", err)
	}
}

func main() {
	url := "http://your-api-endpoint.com/resource" // Replace with your application endpoint
	maxConcurrency := 500                          // Maximum concurrency to test
	step := 50                                     // Increment step for concurrency
	duration := 15 * time.Second                   // Duration for each concurrency level

	// Run the stress tests
	results := stressTest(url, maxConcurrency, step, duration)

	// Save results to a file for visualization
	saveMetricsToFile(results, "stress_test_results.json")

	fmt.Println("Stress test completed. Results saved to stress_test_results.json")
}
