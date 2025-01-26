package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type Metric struct {
	NodeID string `json:"node_id"`
	Data   string `json:"data"`
}

var metrics sync.Map // Store metrics from nodes

// handleMetrics receives and stores metrics sent by nodes
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

	metrics.Store(metric.NodeID, metric.Data)
	fmt.Printf("Received metrics from Node %s: %s\n", metric.NodeID, metric.Data)
	w.WriteHeader(http.StatusOK)
}

// listMetrics allows you to view the aggregated metrics via an HTTP GET request
func listMetrics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed", http.StatusMethodNotAllowed)
		return
	}

	results := make(map[string]string)
	metrics.Range(func(key, value interface{}) bool {
		results[key.(string)] = value.(string)
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
