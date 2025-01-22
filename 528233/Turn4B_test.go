package _28233

import (
	"fmt"
	"log"
	"reflect"
	"runtime"
	"testing"
	"time"
)

func ProcessData(input interface{}) interface{} {
	// Sample implementation
	return input // Just a placeholder
}

// Logger interface for performance and error logging
type Logger interface {
	Info(format string, v ...interface{})
	Warn(format string, v ...interface{})
	Error(format string, v ...interface{})
	Fatal(format string, v ...interface{})
}

// TestLogger - An example logger implementation
type TestLogger struct {
	*testing.T
}

func (tl TestLogger) Info(format string, v ...interface{}) {
	tl.Logf("INFO: "+format, v...)
}

func (tl TestLogger) Warn(format string, v ...interface{}) {
	tl.Logf("WARN: "+format, v...)
}

func (tl TestLogger) Error(format string, v ...interface{}) {
	tl.Logf("ERROR: "+format, v...)
}

func (tl TestLogger) Fatal(format string, v ...interface{}) {
	tl.Fatalf("FATAL: "+format, v...)
}

// BenchmarkReport - Holds performance metrics
type BenchmarkReport struct {
	ExecTime       time.Duration
	MemoryAllocs   uint64
	MemoryAllocsBytes uint64
}

func measurePerf(t *testing.T, f func(), describe ...interface{}) *BenchmarkReport {
	// Before execution, get the memory stats
	var m0 runtime.MemStats
	runtime.ReadMemStats(&m0)
	t0 := time.Now()

	f()

	t1 := time.Now()
	// After execution, get the memory stats
	var m1 runtime.MemStats
	runtime.ReadMemStats(&m1)

	result := BenchmarkReport{
		ExecTime:      t1.Sub(t0),
		MemoryAllocs:  m1.Mallocs - m0.Mallocs,
		MemoryAllocsBytes: m1.Alloc - m0.Alloc,
	}

	// Log performance metrics
	logger := TestLogger{T: t}
	logger.Info("Benchmark Performance Report: %v", describe)
	logger.Info("Execution Time: %v", result.ExecTime)
	logger.Info("Memory Allocs: %v", result.MemoryAllocs)
	logger.Info("Memory Allocs Bytes: %v", result.MemoryAllocsBytes)
	return &result
}

func TestProcessDataWithLogging(t *testing.T) {
	// Define test cases
	testCases := []struct {
		input       interface{}
		expected    interface{}
	}{
		{123, 123},
		{"string", "string"},
		{[]int{1, 2, 3}, []int{1, 2, 3}},
		{map[string]int{"one": 1}, map[string]int{"one": 1}},
	}

	for _, tc := range testCases {
		testName := fmt.Sprintf("TestProcessData_%v", tc.input)
		t.Run(testName, func(t *testing.T) {
			// Measure performance before assertion
			report := measurePerf(t, func() {
				result := ProcessData(tc.input)
				if !reflect.DeepEqual(result, tc.expected) {
					t.Errorf("Mismatch for input %v: got %v, expected %v", tc.input, result, tc.expected)
				}
			}, testName)

			// Display the performance report for a specific testcase (optional)
			// t.Logf("Performance Report: %v", report)
		})
	}