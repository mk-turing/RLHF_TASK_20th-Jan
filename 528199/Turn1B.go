package main

import (
	"fmt"
	"math/rand"
	"sync"
)

// Our dummy large dataset
var dataset = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10} // Replace this with a large slice
const (
	numProcessors  = 4        // Number of worker goroutines
	chunkSize      = 3        // Size of chunks processed by each worker
	iterationCount = 10000000 // Number of iterations to simulate large processing
)

func main() {
	inputChan := make(chan []int)
	outputChan := make(chan []int)
	wg := sync.WaitGroup{}

	// Spawn worker goroutines to process the data
	for i := 0; i < numProcessors; i++ {
		wg.Add(1)
		go worker(inputChan, outputChan, &wg)
	}

	// Generate and stream data chunks into the input channel
	go func() {
		defer close(inputChan)
		for i := 0; i < iterationCount; i++ {
			chunk := make([]int, 0, chunkSize)
			for len(chunk) < chunkSize {
				idx := rand.Intn(len(dataset))
				chunk = append(chunk, dataset[idx])
			}
			inputChan <- chunk
		}
	}()

	// Aggregate results from the output channel
	go func() {
		wg.Wait()
		close(outputChan)
	}()

	var result []int
	for data := range outputChan {
		result = append(result, data...)
	}

	// Visualize the processed result
	fmt.Printf("Processed result: %v\n", result)
}

func worker(input <-chan []int, output chan<- []int, wg *sync.WaitGroup) {
	defer wg.Done()
	for chunk := range input {
		processedChunk := make([]int, len(chunk))
		// Simulate some processing on each element of the chunk using a range loop
		for i, val := range chunk {
			processedChunk[i] = val * 2 // Replace this with your actual processing logic
		}
		output <- processedChunk
	}
}
