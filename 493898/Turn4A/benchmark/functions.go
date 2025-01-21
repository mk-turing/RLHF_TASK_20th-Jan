package benchmark

import (
	"errors"
	"math/rand"
	"sync"
)

// Regular function that performs a sum.
func Sum(numbers []int) int {
	sum := 0
	for _, num := range numbers {
		sum += num
	}
	return sum
}

// Concurrency function that performs sum using Goroutines.
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

// Error-injection function simulating potential errors.
func FaultySum(numbers []int, failRate float64) (int, error) {
	if rand.Float64() < failRate {
		return 0, errors.New("simulated computation error")
	}
	return Sum(numbers), nil
}
