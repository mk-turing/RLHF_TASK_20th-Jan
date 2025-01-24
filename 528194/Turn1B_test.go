package _28194

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

const (
	numOps        = 10000000
	numGoroutines = 100
)

var (
	wg sync.WaitGroup
)

func worker(id int, data []int, ch chan bool) {
	defer wg.Done()
	for i := 0; i < len(data); i++ {
		data[i] = data[i] + 1
		data[i] = data[i] - 1
	}
	ch <- true
}
func main() {
	// Set the number of goroutines to use
	runtime.GOMAXPROCS(numGoroutines)
	data := make([]int, numOps)
	ch := make(chan bool, numGoroutines)
	// Start the timer
	start := time.Now()
	// Launch workers to perform the arithmetic operations
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go worker(i, data[i*numOps/numGoroutines:(i+1)*numOps/numGoroutines], ch)
	}
	// Wait for all workers to complete
	wg.Wait()
	close(ch)
	// Calculate the latency
	latency := time.Since(start) / time.Duration(numOps)
	// Calculate the throughput
	throughput := time.Duration(numOps) / time.Since(start)
	fmt.Printf("Latency: %v\n", latency)
	fmt.Printf("Throughput: %v ops/sec\n", throughput)
}
