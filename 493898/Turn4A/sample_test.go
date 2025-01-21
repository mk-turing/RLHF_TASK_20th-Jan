package main

import (
	"github.com/stretchr/testify/assert"
	"main/test/Turn4A/benchmark"
	"math/rand"
	"os"
	"strconv"
	"testing"
)

// setupNumbers creates and returns a slice of random numbers based on the environment variable N.
func setupNumbers() []int {
	n, _ := strconv.Atoi(os.Getenv("N"))
	numbers := make([]int, n)
	for i := range numbers {
		numbers[i] = rand.Intn(100)
	}
	return numbers
}

// BenchmarkSum benchmarks a simple sum function.
func BenchmarkSum(b *testing.B) {
	numbers := setupNumbers()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		benchmark.Sum(numbers) // Call the sum function from your custom benchmark package
	}
}

// BenchmarkConcurrentSum benchmarks a concurrent sum function with a number of Goroutines specified by G.
func BenchmarkConcurrentSum(b *testing.B) {
	numbers := setupNumbers()
	g, _ := strconv.Atoi(os.Getenv("G")) // Get number of Goroutines from environment variable G
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		benchmark.ConcurrentSum(numbers, g) // Call the concurrent sum function from your package
	}
}

// BenchmarkFaultySumWithTestify benchmarks a function that may return errors and uses Testify for assertions.
func BenchmarkFaultySumWithTestify(b *testing.B) {
	numbers := setupNumbers()
	failRate := 0.1 // 10% chance of failure
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		result, err := benchmark.FaultySum(numbers, failRate) // Call the faulty sum function from your package
		if err != nil {
			b.Log("Error recorded: ", err) // Log error details
		} else {
			assert.NotNil(b, result) // Assert that the result is not nil
		}
	}
}
