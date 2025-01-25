package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Simulating a data stream source with dynamic data size
type DataStream struct {
	sync.Mutex
	data chan []int
}

func NewDataStream(bufferSize int) *DataStream {
	return &DataStream{
		data: make(chan []int, bufferSize),
	}
}

// Function to add data to the data stream (this could be a TCP listener or any data source)
func (ds *DataStream) AddData(data []int) {
	ds.Lock()
	defer ds.Unlock()
	ds.data <- data
}

// Process data using a range loop (stream processing)
func processDataStream(ds *DataStream, wg *sync.WaitGroup, results chan<- int) {
	defer wg.Done()
	for data := range ds.data {
		sum := 0
		// Using range loop for batch processing as well
		for _, value := range data {
			sum += value // Simple processing: summing values
		}
		results <- sum
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	bufferSize := 100 // Buffer size for data stream
	dataStream := NewDataStream(bufferSize)

	results := make(chan int, bufferSize)
	var wg sync.WaitGroup

	// Spawn consumer goroutines
	wg.Add(1)
	go processDataStream(dataStream, &wg, results)

	// Generate random batches of data with fluctuating sizes and add them to the stream
	go func() {
		for {
			dataSize := rand.Intn(10) + 1 // Fluctuating data size between 1 and 100000
			data := make([]int, dataSize)
			for i := range data {
				data[i] = rand.Intn(1000)
			}
			dataStream.AddData(data)
			time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond) // Adding random delay
		}
	}()

	// Aggregating results
	totalSum := 0
	for result := range results {
		totalSum += result
	}

	// Wait for all goroutines to complete (should not be reached since data stream is continuous)
	wg.Wait()
	close(results)

	fmt.Printf("Total sum: %d\n", totalSum)
}
