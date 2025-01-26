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

// startCPUProfile begins CPU profiling and returns a function to stop it.
func startCPUProfile(filename string) func() {
	f, err := os.Create(filename)
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

// captureMemoryProfile captures the memory profile and writes it to a file.
func captureMemoryProfile(filename string) {
	f, err := os.Create(filename)
	if err != nil {
		fmt.Println("Could not create memory profile:", err)
		return
	}
	defer f.Close()
	runtime.GC() // Collect garbage to get an updated memory profile
	if err := pprof.WriteHeapProfile(f); err != nil {
		fmt.Println("Could not write memory profile:", err)
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

func main() {
	// Begin CPU profiling
	stopCPUProfile := startCPUProfile("cpu_workload.prof")
	defer stopCPUProfile()

	numItems := 100000
	var wg sync.WaitGroup

	// Scenario 1: Simulate high frequency updates (resource-intensive range loop)
	wg.Add(1)
	go func() {
		defer wg.Done()
		simulateHighFrequencyUpdates(numItems)
	}()

	// Scenario 2: Simulate concurrent background tasks
	for i := 0; i < 4; i++ { // Launch multiple background tasks simultaneously
		wg.Add(1)
		go simulateConcurrentBackgroundTasks(&wg, numItems)
	}

	// Wait before capturing memory profile
	wg.Wait()
	time.Sleep(2 * time.Second)

	// Capture memory profile of current loads
	captureMemoryProfile("mem_workload.prof")

	// Apply and test optimized concurrent logic
	wg.Add(1)
	go optimizedConcurrentLogic(&wg, numItems)

	wg.Wait()
	time.Sleep(2 * time.Second)

	captureMemoryProfile("mem_optimized_workload.prof")
	fmt.Println("Profiling complete. Analyze the profiles using `go tool pprof`.")
}
