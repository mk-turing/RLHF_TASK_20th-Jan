package main

import (
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
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

// simulateUILogic demonstrates a range loop with potential inefficiencies.
func simulateUILogic(size int) {
	data := make([]int, size)
	for i := range data {
		data[i] = i * 2 // Simulate some work
	}
}

// optimizedUILogic applies optimization strategies for range loops.
func optimizedUILogic(size int) {
	data := make([]int, size)
	chunkSize := 10000
	for i := 0; i < len(data); i += chunkSize {
		end := i + chunkSize
		if end > len(data) {
			end = len(data)
		}
		processChunk(data[i:end])
	}
}

// processChunk simulates processing a chunk of data.
func processChunk(chunk []int) {
	for i := range chunk {
		chunk[i] = chunk[i] * 2 // Simulate optimized processing
	}
}

func main() {
	// Start CPU profiling
	stopCPUProfile := startCPUProfile("cpu.prof")
	defer stopCPUProfile()

	// Simulate initial UI logic (unoptimized range loop)
	simulateUILogic(1000000)

	// Wait to ensure profiling captures sufficient data
	time.Sleep(2 * time.Second)

	// Capture memory profile
	captureMemoryProfile("mem.prof")

	// Apply optimizations and profile again
	stopCPUProfileOptimized := startCPUProfile("cpu_optimized.prof")
	defer stopCPUProfileOptimized()

	optimizedUILogic(1000000)

	// Wait again to ensure sufficient profiling
	time.Sleep(2 * time.Second)

	// Capture memory profile for optimized logic
	captureMemoryProfile("mem_optimized.prof")

	fmt.Println("Profiling completed. Use `go tool pprof` to analyze the profiles.")
}
