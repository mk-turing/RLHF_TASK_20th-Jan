package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Simulating a large dataset
func generateLargeDataset(size int) []int {
	dataset := make([]int, size)
	for i := range dataset {
		dataset[i] = rand.Intn(1000) // Random numbers between 0 and 999
	}
	return dataset
}

// Process data with a range loop
func processData(data []int, wg *sync.WaitGroup, results chan<- int) {
	defer wg.Done()
	sum := 0
	for _, value := range data {
		sum += value // Simple processing: summing values
	}
	results <- sum
}

func main() {
	rand.Seed(time.Now().UnixNano())
	datasetSize := 1000000 // Simulating a large dataset
	dataset := generateLargeDataset(datasetSize)

	// Using channels for efficient communication and synchronization
	results := make(chan int, 10)
	var wg sync.WaitGroup
	chunkSize := 10000

	start := time.Now()

	// Splitting the dataset into chunks and processing each chunk concurrently
	for i := 0; i < len(dataset); i += chunkSize {
		end := i + chunkSize
		if end > len(dataset) {
			end = len(dataset)
		}
		wg.Add(1)
		go processData(dataset[i:end], &wg, results)
	}

	// Closing the results channel once all goroutines are done
	go func() {
		wg.Wait()
		close(results)
	}()

	// Aggregating results
	totalSum := 0
	for result := range results {
		totalSum += result
	}

	elapsed := time.Since(start)
	fmt.Printf("Total sum: %d\n", totalSum)
	fmt.Printf("Processing took %s\n", elapsed)
}
