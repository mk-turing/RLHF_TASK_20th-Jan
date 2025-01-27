package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

// Define the functions you want to benchmark
func FunctionOne() {
	time.Sleep(50 * time.Millisecond) // Simulate some work
}

func FunctionTwo() {
	time.Sleep(100 * time.Millisecond) // Simulate some work
}

// Create a type for function signatures
type BenchmarkFunc struct {
	Name string
	Fn   func()
}

// List all functions to benchmark
var benchmarks = []BenchmarkFunc{
	{"FunctionOne", FunctionOne},
	{"FunctionTwo", FunctionTwo},
}

// Execute benchmark on given functions
func benchmarkFunc(b *testing.B, bf BenchmarkFunc) {
	for n := 0; n < b.N; n++ {
		bf.Fn()
	}
}

// Run the benchmark for Go 1.23 (assumed installed on the system)
func runGoVersionBenchmark() error {
	// Ensure Go 1.23 is installed and being used
	cmd := exec.Command("go", "version")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to get Go version: %v", err)
	}

	// Check if Go 1.23 is the installed version
	if !containsVersion(string(output), "go1.23") {
		return fmt.Errorf("Go 1.23 is not installed or not in use: %s", string(output))
	}

	// Run the benchmarks after confirming Go 1.23 is the version
	for _, bm := range benchmarks {
		result := testing.Benchmark(func(b *testing.B) { benchmarkFunc(b, bm) })
		fmt.Printf("Go 1.23 - %s: %s\n", bm.Name, result)
	}

	// Run code coverage analysis
	cmd = exec.Command("go", "test", "-cover")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to run coverage analysis: %v", err)
	}

	return nil
}

// Check if the Go version string contains the required version
func containsVersion(goOutput, version string) bool {
	return strings.Contains(goOutput, version)
}

func main() {
	// Run benchmarking and coverage for Go 1.23
	err := runGoVersionBenchmark()
	if err != nil {
		fmt.Printf("Error with Go 1.23: %v\n", err)
	}
}
