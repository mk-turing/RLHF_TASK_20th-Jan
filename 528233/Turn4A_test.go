package _28233

import (
	"log"
	"os"
	"reflect"
	"testing"
	"time"
)

// Sample implementation of ProcessData
func ProcessData(input interface{}) interface{} {
	// Example processing logic
	return input
}

// Logger setup
var performanceLog *log.Logger

func init() {
	file, err := os.OpenFile("performance.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	// Initialize logger
	performanceLog = log.New(file, "PERF ", log.Ldate|log.Ltime|log.Lshortfile)
}

// Benchmark function enhanced with logging
func benchmarkProcessData(b *testing.B, input interface{}) {
	b.Helper() // Mark function as a helper

	for i := 0; i < b.N; i++ {
		start := time.Now()

		// Example reflection use case
		val := reflect.ValueOf(input)
		ProcessData(val.Interface())

		duration := time.Since(start)

		// Log execution time
		performanceLog.Printf("Input: %v, Type: %T, ExecutionTime: %v", input, input, duration)
	}
}

func BenchmarkFullSuite(b *testing.B) {
	testCases := []interface{}{123, "example", []int{1, 2, 3}, map[string]int{"key": 1}}
	for _, tc := range testCases {
		b.Run(reflect.TypeOf(tc).String(), func(b *testing.B) {
			benchmarkProcessData(b, tc)
		})
	}
}
