package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

type Metric struct {
	NodeID string `json:"node_id"`
	Data   string `json:"data"`
}

func simulateWorkload(nodeID string) string {
	// Simulate a workload and range loop inefficiency
	time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
	return fmt.Sprintf("Node %s: CPU: %d%%, Memory: %dMB", nodeID, rand.Intn(100), rand.Intn(2000))
}

func sendMetrics(serverURL, nodeID string) {
	for {
		// Simulate workload and collect metrics
		metrics := Metric{
			NodeID: nodeID,
			Data:   simulateWorkload(nodeID),
		}

		// Convert the metric to JSON
		body, err := json.Marshal(metrics)
		if err != nil {
			fmt.Println("Failed to serialize metrics:", err)
			return
		}

		// Send metrics to the server
		resp, err := http.Post(serverURL, "application/json", bytes.NewBuffer(body))
		if err != nil {
			fmt.Println("Failed to send metrics:", err)
			return
		}
		resp.Body.Close()

		time.Sleep(1 * time.Second) // Send metrics every second
	}
}

func main() {
	serverURL := "http://localhost:8080/metrics"
	nodeID := "node-1"

	fmt.Printf("Node %s is sending metrics to %s...\n", nodeID, serverURL)
	sendMetrics(serverURL, nodeID)
}
