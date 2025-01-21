package _93898

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"testing"
)

// ConcurrentSum computes the sum of numbers concurrently using Goroutines.
func ConcurrentSum(numbers []int, numWorkers int) int {
	partSize := (len(numbers) + numWorkers - 1) / numWorkers

	results := make(chan int, numWorkers)
	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(start int) {
			defer wg.Done()
			end := start + partSize
			if end > len(numbers) {
				end = len(numbers)
			}
			sum := 0
			for _, num := range numbers[start:end] {
				sum += num
			}
			results <- sum
		}(i * partSize)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	finalSum := 0
	for sum := range results {
		finalSum += sum
	}

	return finalSum
}

// BenchmarkConcurrentSum benchmarks ConcurrentSum using the testing package.
func BenchmarkConcurrentSum(b *testing.B) {
	// Get environment variables for N (size) and G (number of Goroutines)
	n, _ := strconv.Atoi(os.Getenv("N"))
	g, _ := strconv.Atoi(os.Getenv("G"))

	// Generate a large slice of random numbers for benchmarking.
	numbers := make([]int, n)
	for i := range numbers {
		numbers[i] = rand.Intn(100)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ConcurrentSum(numbers, g)
	}
}

// BenchmarkConcurrentSumWithTestify benchmarks ConcurrentSum using the testify package.
func BenchmarkConcurrentSumWithTestify(b *testing.B) {
	// Get environment variables for N (size) and G (number of Goroutines)
	n, _ := strconv.Atoi(os.Getenv("N"))
	g, _ := strconv.Atoi(os.Getenv("G"))

	// Generate a large slice of random numbers for benchmarking.
	numbers := make([]int, n)
	for i := range numbers {
		numbers[i] = rand.Intn(100)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := ConcurrentSum(numbers, g)
		assert.NotNil(b, result)
		assert.Greater(b, result, 0)
	}
}
