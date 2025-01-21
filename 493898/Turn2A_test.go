package _93898

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func ConcurrentSum(numbers []int) int {
	const numWorkers = 4
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

func BenchmarkConcurrentSum(b *testing.B) {
	numbers := make([]int, 1000000)
	for i := range numbers {
		numbers[i] = i
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ConcurrentSum(numbers)
	}
}

func BenchmarkConcurrentSumWithTestify(b *testing.B) {
	numbers := make([]int, 1000000)
	for i := range numbers {
		numbers[i] = i
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		result := ConcurrentSum(numbers)
		assert.NotNil(b, result)
		assert.Greater(b, result, 0)
	}
}
