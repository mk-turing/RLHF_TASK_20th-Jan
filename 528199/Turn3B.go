package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const (
	batchSize    = 10
	numStreams   = 3 // Number of concurrent data streams
	maxQueueSize = 10
	totalBatches = 5
)

// Simulate real-time data for multiple streams
func realTimeDataGenerator(streamID int, dataStream chan<- int) {
	for range totalBatches {
		for i := 0; i < batchSize; i++ {
			data := rand.Intn(1000)
			select {
			case dataStream <- data:
				fmt.Printf("Stream %d: Data sent: %d\n", streamID, data)
				// Simulate variable data input intervals
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(50)))
			}
		}
		time.Sleep(time.Second) // Wait between batches
	}
}

func processDataBatch(data []int, streamID int, wg *sync.WaitGroup, results chan<- int) {
	defer wg.Done()
	sum := 0
	for _, value := range data {
		sum += value
	}
	fmt.Printf("Stream %d: Batch processed, sum: %d\n", streamID, sum)
	results <- sum
}

func main() {
	rand.Seed(time.Now().UnixNano())
	results := make(chan int, numStreams)
	var wg sync.WaitGroup

	// Process each data stream in a separate goroutine
	for i := 0; i < numStreams; i++ {
		dataStream := make(chan int, maxQueueSize)
		go func(streamID int) {
			defer close(dataStream)
			realTimeDataGenerator(streamID, dataStream)
		}(i)

		wg.Add(1)
		go func(streamID int, dataStream <-chan int) {
			defer wg.Done()
			currentBatch := make([]int, 0, batchSize)

			for dataPoint := range dataStream {
				currentBatch = append(currentBatch, dataPoint)
				if len(currentBatch) == batchSize {
					wg.Add(1)
					go processDataBatch(append([]int(nil), currentBatch...), streamID, &wg, results)
					currentBatch = currentBatch[:0]
				}
			}

			// Process the remaining data in the last batch
			if len(currentBatch) > 0 {
				wg.Add(1)
				go processDataBatch(append([]int(nil), currentBatch...), streamID, &wg, results)
			}
		}(i, dataStream)
	}

	// Wait for all processing goroutines to finish
	go func() {
		wg.Wait()
		close(results)
	}()

	// Sum up the results from all streams
	totalSum := 0
	for result := range results {
		totalSum += result
	}

	fmt.Printf("Total sum of all processed data: %d\n", totalSum)
}
