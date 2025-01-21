package _93898

import (
	"fmt"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Sample function to be benchmarked
func sampleFunction(n int) int {
	sum := 0
	for i := 0; i < n; i++ {
		sum += i
	}
	return sum
}

// Benchmark function using "testing" package
func BenchmarkSampleFunctionTesting(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sampleFunction(1000000) // Calling the sample function with a large input size
	}
}

// Benchmark function using "testify" package
func BenchmarkSampleFunctionTestify(b *testing.B) {
	for i := 0; i < b.N; i++ {
		assert.Equal(b, 499999500000, sampleFunction(1000000)) // Using testify's assert for validation
	}
}

func measureMemoryUsage() uint64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return m.Alloc
}

func main() {
	const numIterations = 5
	fmt.Println("\nBenchmark Results:")
	fmt.Printf("%-25s %-20s %-20s\n", "Framework", "Execution Time (ms)", "Memory Usage (bytes)")
	fmt.Println("---------------------------------------------------------------")
	// Benchmarking using "testing" package
	totalTimeTesting := 0.0
	totalMemUsageTesting := uint64(0)
	for i := 0; i < numIterations; i++ {
		start := time.Now()
		testing.Benchmark(BenchmarkSampleFunctionTesting)
		end := time.Now()
		duration := end.Sub(start)
		totalTimeTesting += duration.Seconds() * 1000
		totalMemUsageTesting += measureMemoryUsage()
	}
	avgTimeTesting := totalTimeTesting / float64(numIterations)
	avgMemUsageTesting := totalMemUsageTesting / uint64(numIterations)
	fmt.Printf("%-25s %-20.2f %-20d\n", "testing", avgTimeTesting, avgMemUsageTesting)

	// Benchmarking using "testify" package
	totalTimeTestify := 0.0
	totalMemUsageTestify := uint64(0)
	for i := 0; i < numIterations; i++ {
		start := time.Now()
		testing.Benchmark(BenchmarkSampleFunctionTestify)
		end := time.Now()
		duration := end.Sub(start)
		totalTimeTestify += duration.Seconds() * 1000
		totalMemUsageTestify += measureMemoryUsage()
	}
	avgTimeTestify := totalTimeTestify / float64(numIterations)
	avgMemUsageTestify := totalMemUsageTestify / uint64(numIterations)
	fmt.Printf("%-25s %-20.2f %-20d\n", "testify", avgTimeTestify, avgMemUsageTestify)

}
