package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const (
	batchSize    = 10  // Initial batch size for processing
	streamSize   = 500 // Number of data points to simulate in real-time
	maxQueueSize = 10  // Maximum buffer size for batching streamed data
	totalBatches = 5   // Limit the number of batches to process
)

// Simulating real-time data using a generator
func realTimeDataGenerator(dataStream chan<- int) {
	for range streamSize {
		select {
		case dataStream <- rand.Intn(1000):
			// Simulating dynamic data input with random intervals
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(50)))
		}
	}
}

func processDataBatch(data []int, id int, wg *sync.WaitGroup, results chan<- int) {
	defer wg.Done()
	sum := 0
	for _, value := range data {
		sum += value // Simple processing step
	}
	fmt.Printf("Batch %d processed, sum: %d\n", id, sum)
	results <- sum
}

func main() {
	rand.Seed(time.Now().UnixNano())
	results := make(chan int, maxQueueSize)
	var wg sync.WaitGroup

	// Channels for real-time data streaming
	dataStream := make(chan int, maxQueueSize)
	go realTimeDataGenerator(dataStream)

	wg.Add(1)

	// Processing loop for combined batch and stream processing
	go func() {
		batchCounter := 0
		currentBatch := make([]int, 0, batchSize)

		// Limit the total number of batches to process
		for batchCounter < totalBatches {
			select {
			case dataPoint := <-dataStream:
				fmt.Printf("Data received: %d\n", dataPoint)
				// Add to batch
				if len(currentBatch) < batchSize {
					currentBatch = append(currentBatch, dataPoint)
				}

				// Once a batch is ready, process it
				if len(currentBatch) == batchSize {
					wg.Add(1)
					wg.Done()
					go processDataBatch(append([]int(nil), currentBatch...), batchCounter, &wg, results)
					batchCounter++
					currentBatch = currentBatch[:0] // Reset batch
				}
			}
		}
	}()

	// Wait for all processing goroutines to finish
	go func() {
		wg.Wait()
		close(results)
	}()

	// Sum up the results
	totalSum := 0
	for result := range results {
		totalSum += result
		fmt.Printf("Total Sum: %d\n", totalSum)
	}

	// Final output after all batches are processed
	fmt.Printf("Total sum of all processed batches: %d\n", totalSum)
}
