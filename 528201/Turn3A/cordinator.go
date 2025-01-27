package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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

var (
	metricsStore []Metrics
	mu           sync.Mutex
)

func reportHandler(w http.ResponseWriter, r *http.Request) {
	var metrics Metrics
	err := json.NewDecoder(r.Body).Decode(&metrics)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	mu.Lock()
	metricsStore = append(metricsStore, metrics)
	mu.Unlock()

	fmt.Fprintf(w, "Metrics received from %s", metrics.AgentID)
}

func aggregateMetrics() Metrics {
	var totalMetrics Metrics
	for _, m := range metricsStore {
		totalMetrics.TotalRequests += m.TotalRequests
		totalMetrics.FailedRequests += m.FailedRequests
		totalMetrics.CPUUsage += m.CPUUsage
		totalMetrics.MemoryUsage += m.MemoryUsage
	}

	totalMetrics.ErrorRate = float64(totalMetrics.FailedRequests) / float64(totalMetrics.TotalRequests) * 100
	totalMetrics.CPUUsage /= float64(len(metricsStore))
	totalMetrics.MemoryUsage /= uint64(len(metricsStore))

	return totalMetrics
}

func main() {
	http.HandleFunc("/report", reportHandler)

	go func() {
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	fmt.Println("Coordinator is running on port 8080. Waiting for metrics...")

	// Simple aggregation example, in a real-world scenario you might want to persist this data
	time.Sleep(1 * time.Minute)
	aggregatedMetrics := aggregateMetrics()
	fmt.Printf("Aggregated Metrics: %+v\n", aggregatedMetrics)
}
