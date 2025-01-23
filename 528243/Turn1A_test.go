package _28243

import (
	"fmt"
	"runtime"
	"testing"
	"time"
)

func FunctionToTest(input int) int {
	// Example function logic
	return input * 2
}

func monitorPerformance(testFunc func(*testing.T), t *testing.T) {
	start := time.Now()

	var memStatsStart, memStatsEnd runtime.MemStats
	runtime.ReadMemStats(&memStatsStart)

	testFunc(t)

	runtime.ReadMemStats(&memStatsEnd)
	elapsed := time.Since(start)

	fmt.Printf("Test Execution Time: %s\n", elapsed)
	fmt.Printf("Memory Allocations: %d bytes\n", memStatsEnd.Alloc-memStatsStart.Alloc)
}

func TestFunctionToTest(t *testing.T) {
	monitorPerformance(func(t *testing.T) {
		result := FunctionToTest(2)
		if result != 4 {
			t.Errorf("Expected 4, but got %d", result)
		}
	}, t)
}
