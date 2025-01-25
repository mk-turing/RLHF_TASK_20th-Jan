package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"
)

const (
	batchSize      = 10
	streamSize     = 500
	maxQueueSize   = 10
	totalBatches   = 5
	numStreams     = 3
	maxRetry       = 3 // Max retries for failed batches
	errorThreshold = 5 // Max errors before terminating a stream
)

type dataError struct {
	msg string
}

func (e *dataError) Error() string {
	return e.msg
}

func realTimeDataGenerator(dataStream chan<- int, streamId int) {
	defer close(dataStream)
	for i := 0; i < streamSize; i++ {
		if rand.Intn(100) < 2 { // Simulated 2% chance of error
			log.Printf("Stream %d - Simulated network interruption\n", streamId)
			return
		}
		dataStream <- rand.Intn(1000)
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(50)))
	}
}

func processDataBatch(data []int, streamId, batchId, retryCount int, results chan<- int, errors chan<- error) bool {
	sum := 0
	for _, value := range data {
		if rand.Intn(100) < 1 { // Simulated 1% chance of error
			errors <- &dataError{msg: fmt.Sprintf("Stream %d - Batch %d contains an error (retry %d)", streamId, batchId, retryCount)}
			return false
		}
		sum += value
	}
	log.Printf("Stream %d - Batch %d processed successfully, sum: %d\n", streamId, batchId, sum)
	results <- sum
	return true
}

func handleStream(streamId int, results chan<- int, wg *sync.WaitGroup, errors chan<- error) {
	defer wg.Done()
	dataStream := make(chan int, maxQueueSize)
	go realTimeDataGenerator(dataStream, streamId)

	var localWg sync.WaitGroup
	batchCounter, errorCount := 0, 0
	currentBatch := make([]int, 0, batchSize)

	for batchCounter < totalBatches && errorCount < errorThreshold {
		dataPoint, ok := <-dataStream
		if !ok {
			if len(currentBatch) > 0 {
				localWg.Add(1)
				go func(batch []int, batchId int) {
					defer localWg.Done()
					processDataBatch(batch, streamId, batchId, 0, results, errors)
				}(append([]int(nil), currentBatch...), batchCounter)
				batchCounter++
			}
			break
		}
		currentBatch = append(currentBatch, dataPoint)
		if len(currentBatch) == batchSize {
			localWg.Add(1)
			go func(batch []int, batchId int) {
				defer localWg.Done()
				retries := 0
				for !processDataBatch(batch, streamId, batchId, retries, results, errors) && retries < maxRetry {
					retries++
				}
				if retries == maxRetry {
					errors <- fmt.Errorf("Stream %d - Batch %d failed after %d retries", streamId, batchId, retries)
				}
			}(append([]int(nil), currentBatch...), batchCounter)
			batchCounter++
			currentBatch = currentBatch[:0]
		}
	}
	localWg.Wait()
}

func setupLogging() {
	logFile, err := os.OpenFile("system.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}
	log.SetOutput(logFile)
	log.Println("Logging initiated")
}

func main() {
	rand.Seed(time.Now().UnixNano())
	setupLogging()

	results := make(chan int, maxQueueSize*numStreams)
	errors := make(chan error)
	var wg sync.WaitGroup

	for streamId := 0; streamId < numStreams; streamId++ {
		wg.Add(1)
		go handleStream(streamId, results, &wg, errors)
	}

	go func() {
		for err := range errors {
			log.Printf("Error occurred: %v", err)
		}
	}()

	go func() {
		wg.Wait()
		close(results)
		close(errors)
	}()

	totalSum := 0
	for result := range results {
		totalSum += result
	}

	log.Printf("Total sum of all processed batches across all streams: %d\n", totalSum)
	fmt.Printf("Total sum of all processed batches across all streams: %d\n", totalSum)
}
