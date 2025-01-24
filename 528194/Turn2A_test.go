package _28194

import (
	"strconv"
	"sync"
	"testing"
	"time"
)

// Simulate a task with shared resource contention using a mutex
func incrementWithMutex(numIncrements, numGoroutines int) int {
	var counter int
	var mu sync.Mutex
	var wg sync.WaitGroup

	incrementsPerRoutine := numIncrements / numGoroutines

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < incrementsPerRoutine; j++ {
				mu.Lock()
				counter++
				mu.Unlock()
			}
		}()
	}

	wg.Wait()
	return counter
}

// Simulate a task with shared resource contention using a channel
func incrementWithChannel(numIncrements, numGoroutines int) int {
	counter := 0
	counterChan := make(chan int)
	var wg sync.WaitGroup

	incrementsPerRoutine := numIncrements / numGoroutines

	// Goroutine to gather increments
	go func() {
		for val := range counterChan {
			counter += val
		}
	}()

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			localCount := 0
			for j := 0; j < incrementsPerRoutine; j++ {
				localCount++
			}
			counterChan <- localCount
		}()
	}

	wg.Wait()
	close(counterChan)
	return counter
}

// Benchmark the shared resource contention
func BenchmarkSharedResourceContention(b *testing.B) {
	numIncrements := 1e6
	for _, numGoroutines := range []int{1, 5, 10, 20, 50, 100} {
		b.Run("Mutex_"+strconv.Itoa(numGoroutines), func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				_ = incrementWithMutex(numIncrements, numGoroutines)
			}
		})
		b.Run("Channel_"+strconv.Itoa(numGoroutines), func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				_ = incrementWithChannel(numIncrements, numGoroutines)
			}
		})
	}
}

// Function to analyze performance
func analyzePerformance() {
	numIncrements := 1e6
	numGoroutines := 10

	start := time.Now()
	_ = incrementWithMutex(numIncrements, numGoroutines)
	elapsedMutex := time.Since(start)

	start = time.Now()
	_ = incrementWithChannel(numIncrements, numGoroutines)
	elapsedChannel := time.Since(start)

	println("Mutex:", elapsedMutex.Milliseconds(), "ms, Channel:", elapsedChannel.Milliseconds(), "ms")
}
