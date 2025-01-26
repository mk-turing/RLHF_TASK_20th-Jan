package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"runtime"
	"time"
)

type Metric struct {
	NodeID      string `json:"node_id"`
	CPUUsage    int    `json:"cpu_usage"`
	MemoryUsage int    `json:"memory_usage"`
	Interaction string `json:"interaction"`
	Timestamp   int64  `json:"timestamp"`
}

// simulateUIActions to mimic frequency of UI operations like resizing and drag-and-drop.
func simulateUIActions(action string) {
	switch action {
	case "resize":
		// Simulate resizing action
		fmt.Println("Simulating resizing...")
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
	case "drag_and_drop":
		// Simulate drag-and-drop action
		fmt.Println("Simulating drag-and-drop...")
		time.Sleep(time.Duration(rand.Intn(50)) * time.Millisecond)
	}
}

// collectMetrics collects CPU and Memory usage statistics.
func collectMetrics(nodeID, interaction string) Metric {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return Metric{
		NodeID:      nodeID,
		CPUUsage:    rand.Intn(100), // Placeholder for actual CPU usage retrieval
		MemoryUsage: int(m.Alloc / 1024 / 1024),
		Interaction: interaction,
		Timestamp:   time.Now().Unix(),
	}
}

// sendMetrics posts collected metrics to the central server.
func sendMetrics(serverURL string, metric Metric) {
	body, err := json.Marshal(metric)
	if err != nil {
		fmt.Println("Failed to serialize metrics:", err)
		return
	}

	resp, err := http.Post(serverURL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("Failed to send metrics:", err)
		return
	}
	resp.Body.Close()
}

func main() {
	serverURL := "http://localhost:8080/metrics"
	nodeID := "node-1"
	interactions := []string{"resize", "drag_and_drop"}

	for {
		for _, interaction := range interactions {
			simulateUIActions(interaction)
			metric := collectMetrics(nodeID, interaction)
			sendMetrics(serverURL, metric)
			time.Sleep(1 * time.Second)
		}
	}
}
