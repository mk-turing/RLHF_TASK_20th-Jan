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

type WorkloadType int

const (
	SmallWorkload WorkloadType = iota
	MediumWorkload
	LargeWorkload
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

// simulateUILogic demonstrates a range loop with potential inefficiencies, varying workload intensity.
func simulateUILogic(workload WorkloadType) {
	data := generateData(workload)
	for i := range data {
		data[i] = data[i] * 2 // Simulate some work
		// Add random delay to simulate different workload intensities
		time.Sleep(time.Duration(rand.Intn(100)) * time.Microsecond)
	}
}

// generateData creates sample data based on the workload type.
func generateData(workload WorkloadType) []int {
	size := 0
	switch workload {
	case SmallWorkload:
		size = 10000
	case MediumWorkload:
		size = 100000
	case LargeWorkload:
		size = 1000000
	}
	data := make([]int, size)
	for i := range data {
		data[i] = i
	}
	return data
}

// optimizedUILogic applies optimization strategies for range loops, varying workload intensity.
func optimizedUILogic(workload WorkloadType) {
	data := generateData(workload)
	chunkSize := 10000
	switch workload {
	case SmallWorkload:
		chunkSize = 1000
	case MediumWorkload:
		chunkSize = 10000
	case LargeWorkload:
		chunkSize = 100000
	}
	for i := 0; i < len(data); i += chunkSize {
		end := i + chunkSize
		if end > len(data) {
			end = len(data)
		}
		processChunk(data[i:end])
	}
}

// processChunk simulates processing a chunk of data with varying workload intensity.
func processChunk(chunk []int) {
	for i := range chunk {
		chunk[i] = chunk[i] * 2 // Simulate optimized processing
		// Add random delay to simulate different workload intensities
		time.Sleep(time.Duration(rand.Intn(100)) * time.Microsecond)
	}
}

func concurrentWorkload(workload WorkloadType, wg *sync.WaitGroup) {
	defer wg.Done()
	// Simulate concurrent background tasks
	for i := 0; i < 5; i++ {
		simulateUILogic(workload)
		time.Sleep(time.Millisecond * 100)
	}
}

func main() {
	// Start CPU profiling