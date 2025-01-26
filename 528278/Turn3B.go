package main

import (
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"runtime/pprof"
	"sync"
	"time"
)

const (
	numNodes    = 3       // Number of nodes in the distributed system
	profilePort = ":8080" // Port for listening to profiling requests
)

// startCPUProfile begins CPU profiling and returns a function to stop it.
func startCPUProfile() func() {
	f, err := os.Create("cpu.prof")
	if err != nil {
		fmt.Println("Could not create CPU profile:", err)
		return func() {}
	}
	if err := pprof.StartCPUProfile(f); err != nil {
		fmt.Println("Could not start CPU profile:", err)
		return func() {}
	}
	return func() {
		pprof.StopCPUProfile()
		f.Close()
	}
}

// simulateHighFrequencyUpdates simulates frequent UI changes.
func simulateHighFrequencyUpdates(size int) {
	data := make([]int, size)
	for i := range data {
		data[i] = i * 2
		time.Sleep(time.Millisecond) // Simulate processing time per element
	}
}

// simulateConcurrentBackgroundTasks simulates background tasks.
func simulateConcurrentBackgroundTasks(wg *sync.WaitGroup, size int) {
	defer wg.Done()
	data := make([]int, size)
	for i := range data {
		data[i] = rand.Int() % 100   // Simulate computational workload
		time.Sleep(time.Microsecond) // Shorter delay to mimic async background work
	}
}

// optimizedConcurrentLogic applies optimization with goroutines.
func optimizedConcurrentLogic(wg *sync.WaitGroup, size int) {
	defer wg.Done()
	chunkSize := 10000
	data := make([]int, size)
	for i := 0; i < len(data); i += chunkSize {
		go func(start, end int) {
			if end > len(data) {
				end = len(data)
			}
			for j := start; j < end; j++ {
				data[j] = data[j] * 2
			}
		}(i, i+chunkSize)
	}
}

// AggregatedProfile represents the aggregated performance metrics across nodes.
type AggregatedProfile struct {
	CPUProfile    []byte
	MemoryProfile []byte
}

// aggregateProfiles aggregates performance metrics from multiple nodes.
func aggregateProfiles(nodes map[int]*AggregatedProfile) *AggregatedProfile {
	aggregatedProfile := &AggregatedProfile{}
	for _, profile := range nodes {
		aggregatedProfile.CPUProfile = append(aggregatedProfile.CPUProfile, profile.CPUProfile...)
		aggregatedProfile.MemoryProfile = append(aggregatedProfile.MemoryProfile, profile.MemoryProfile...)
	}
	return aggregatedProfile
}

// collectProfile collects performance metrics for the current node.
func collectProfile() *AggregatedProfile {
	profile := &AggregatedProfile{}

	stopCPUProfile := startCPUProfile()
	defer stopCPUProfile()

	simulateHighFrequencyUpdates(100000)
	simulateConcurrentBackgroundTasks(&sync.WaitGroup{}, 100000)
	optimizedConcurrentLogic(&sync.WaitGroup{}, 100000)

	runtime.GC()
	pprof.Lookup("heap").WriteTo(profile.MemoryProfile, 0)

	return profile
}

// syncProfiles synchronizes performance metrics among nodes.
func syncProfiles(nodes map[int]chan *AggregatedProfile) {
	var wg sync.WaitGroup
	for id, ch := range nodes {
		wg.Add(1)
		go func(id int, ch chan<- *AggregatedProfile) {
			defer wg.Done()
			profile := collectProfile()
			ch <- profile
		}(id, ch)
	}
	wg.Wait()
}
