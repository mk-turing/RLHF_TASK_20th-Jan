package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type Metric struct {
	NodeID      string `json:"node_id"`
	CPUUsage    int    `json:"cpu_usage"`
	MemoryUsage int    `json:"memory_usage"`
	Interaction string `json:"interaction"`
	Timestamp   int64  `json:"timestamp"`
}

var metricsStore sync.Map

func handleMetrics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed", http.StatusMethodNotAllowed)
		return
	}

	var metric Metric
	if err := json.NewDecoder(r.Body).Decode(&metric); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	metricsStore.Store(metric.Timestamp, metric)
	fmt.Printf("Received metrics: %+v\n", metric)
	w.WriteHeader(http.StatusOK)
}

func listMetrics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed", http.StatusMethodNotAllowed)
		return
	}

	results := []Metric{}
	metricsStore.Range(func(key, value interface{}) bool {
		results = append(results, value.(Metric))
		return true
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func main() {
	http.HandleFunc("/metrics", handleMetrics)
	http.HandleFunc("/metrics/list", listMetrics)

	fmt.Println("Server is running on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Failed to start server:", err)
	}
}
