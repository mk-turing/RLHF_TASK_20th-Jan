package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const (
	batchSize    = 10  // Initial batch size for processing
	streamSize   = 500 // Number of data points to simulate for each real-time stream
	maxQueueSize = 10  // Maximum buffer size for batching streamed data
	totalBatches = 5   // Limit the number of batches to process per stream
	numStreams   = 3   // Number of concurrent data streams
)

// Simulating real-time data using a generator
func realTimeDataGenerator(dataStream chan<- int, streamId int) {
	for i := 0; i < streamSize; i++ {
		select {
		case dataStream <- rand.Intn(1000):
			// Simulating dynamic data input with random intervals
			fmt.Printf("Stream %d - Data generated\n", streamId)
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(50)))
		}
	}
	close(dataStream)
}

func processDataBatch(data []int, streamId, batchId int, wg *sync.WaitGroup, results chan<- int) {
	defer wg.Done()
	sum := 0
	for _, value := range data {
		sum += value // Simple processing step
	}
	fmt.Printf("Stream %d - Batch %d processed, sum: %d\n", streamId, batchId, sum)
	results <- sum
}

func handleStream(streamId int, results chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()

	dataStream := make(chan int, maxQueueSize)
	go realTimeDataGenerator(dataStream, streamId)

	var localWg sync.WaitGroup
	batchCounter := 0
	currentBatch := make([]int, 0, batchSize)

	for batchCounter < totalBatches {
		select {
		case dataPoint, ok := <-dataStream:
			if !ok {
				return // Exit if the data stream is closed
			}
			// Add to batch
			if len(currentBatch) < batchSize {
				currentBatch = append(currentBatch, dataPoint)
			}

			// Once a batch is ready, process it
			if len(currentBatch) == batchSize {
				localWg.Add(1)
				// Pass a copy of the current batch to avoid race conditions
				go processDataBatch(append([]int(nil), currentBatch...), streamId, batchCounter, &localWg, results)
				batchCounter++
				currentBatch = currentBatch[:0] // Reset batch
			}
		}
	}

	// Wait for all batches within this stream to finish
	localWg.Wait()
}

func main() {
	rand.Seed(time.Now().UnixNano())
	results := make(chan int, maxQueueSize*numStreams)
	var wg sync.WaitGroup

	for streamId := 0; streamId < numStreams; streamId++ {
		wg.Add(1)
		go handleStream(streamId, results, &wg)
	}

	// Wait for all streams to complete
	go func() {
		wg.Wait()
		close(results)
	}()

	// Sum up the results
	totalSum := 0
	for result := range results {
		totalSum += result
	}

	// Final output after all streams and batches are processed
	fmt.Printf("Total sum of all processed batches across all streams: %d\n", totalSum)
}
