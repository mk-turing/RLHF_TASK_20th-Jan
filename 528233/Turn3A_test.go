package _28233

import (
	"reflect"
	"testing"
)

// Sample implementation of ProcessData (placeholder)
func ProcessData(input interface{}) interface{} {
	return input // Just a placeholder
}

func BenchmarkProcessDataOptimized(b *testing.B) {
	// Set up test cases with reflection optimizations
	testCases := []interface{}{
		123,                      // Integer
		"optimized string",       // String
		[]int{1, 2, 3},           // Slice
		map[string]int{"one": 1}, // Map
	}

	for _, tc := range testCases {
		b.Run(reflect.TypeOf(tc).String(), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				// Calling an optimized version of ProcessData on input
				ProcessData(tc)
			}
		})
	}
}

func BenchmarkProcessDataUnoptimized(b *testing.B) {
	// Set up test cases without reflection optimizations
	richTestCases := []interface{}{
		123,                      // Integer
		"basic string",           // String
		[]int{1, 2, 3},           // Slice
		map[string]int{"key": 1}, // Map
	}

	for _, tc := range richTestCases {
		b.Run(reflect.TypeOf(tc).String(), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				// Directly scanning through each tc without optimization
				ProcessData(tc)
			}
		})
	}
}
