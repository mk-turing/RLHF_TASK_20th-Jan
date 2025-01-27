package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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
	Concurrency    int           `json:"concurrency"`
	ElapsedTime    time.Duration `json:"elapsed_time"`
	AgentID        string        `json:"agent_id"`
}

func runTest(url string, concurrency int, duration time.Duration) Metrics {
	var wg sync.WaitGroup
	metrics := Metrics{Concurrency: concurrency, AgentID: "agent-1"} // Assign unique ID per agent
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
	metrics.ErrorRate = float64(metrics.FailedRequests) / float64(metrics.TotalRequests) * 100
	metrics.ElapsedTime = time.Since(startTime)

	return metrics
}

func reportMetrics(metrics Metrics) {
	url := "http://localhost:8080/report" // Replace with your coordinator endpoint
	metricsJSON, err := json.Marshal(metrics)
	if err != nil {
		log.Fatalf("Failed to serialize metrics: %v", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(metricsJSON))
	if err != nil {
		log.Fatalf("Failed to report metrics: %v", err)
	}
	defer resp.Body.Close()
}

func main() {
	url := "http://your-api-endpoint.com/resource" // Replace with your application endpoint
	concurrency := 100                             // Concurrency level
	duration := 15 * time.Second                   // Duration for the test

	// Run the stress test
	metrics := runTest(url, concurrency, duration)

	// Report the metrics to the coordinator
	reportMetrics(metrics)

	fmt.Println("Test completed. Metrics reported.")
}
