package main

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

// Simulate a complex UI interaction
func SimulateUserInteraction(id int, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("User %d starts interaction\n", id)
	// Simulate processing time
	time.Sleep(100 * time.Millisecond)
	fmt.Printf("User %d completes interaction\n", id)
}

// Stress test simulating multiple users
func StressTestConcurrentUsers(b *testing.B, userCount int) {
	for n := 0; n < b.N; n++ {
		var wg sync.WaitGroup
		for i := 0; i < userCount; i++ {
			wg.Add(1)
			go SimulateUserInteraction(i, &wg)
		}
		wg.Wait()
	}
}

// This is the entry point for the test
func BenchmarkStressTest(b *testing.B) {
	userCount := 1000 // Number of concurrent users
	StressTestConcurrentUsers(b, userCount)
}
