package _28194

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// SharedCounter with Mutex
type SharedCounterMutex struct {
	count int
	mutex sync.Mutex
}

func (sc *SharedCounterMutex) Increment() {
	sc.mutex.Lock()
	sc.count++
	sc.mutex.Unlock()
}

func (sc *SharedCounterMutex) Decrement() {
	sc.mutex.Lock()
	sc.count--
	sc.mutex.Unlock()
}

func (sc *SharedCounterMutex) Count() int {
	sc.mutex.Lock()
	defer sc.mutex.Unlock()
	return sc.count
}

// SharedCounter with Atomic
type SharedCounterAtomic struct {
	count int64
}

func (sc *SharedCounterAtomic) Increment() {
	atomic.AddInt64(&sc.count, 1)
}

func (sc *SharedCounterAtomic) Decrement() {
	atomic.AddInt64(&sc.count, -1)
}

func (sc *SharedCounterAtomic) Count() int64 {
	return atomic.LoadInt64(&sc.count)
}

// Parallel Counter Operations
func parallelCounterOperations(nums []int, numGoroutines int, sharedCounter interface{}) {
	chunkSize := (len(nums) + numGoroutines - 1) / numGoroutines
	var wg sync.WaitGroup

	for i := 0; i < len(nums); i += chunkSize {
		wg.Add(1)
		go func(start int) {
			defer wg.Done()
			end := start + chunkSize
			if end > len(nums) {
				end = len(nums)
			}
			for range nums[start:end] {
				if rand.Intn(2) == 0 {
					switch counter := sharedCounter.(type) {
					case *SharedCounterMutex:
						counter.Increment()
					case *SharedCounterAtomic:
						counter.Increment()
					}
				} else {
					switch counter := sharedCounter.(type) {
					case *SharedCounterMutex:
						counter.Decrement()
					case *SharedCounterAtomic:
						counter.Decrement()
					}
				}
			}
		}(i)
	}

	wg.Wait()
}

func generateRandomSlice(size int) []int {
	nums := make([]int, size)
	for i := range nums {
		nums[i] = rand.Intn(100)
	}
	return nums
}

// Benchmark with Throughput and Latency Metrics
func BenchmarkParallelCounter(b *testing.B) {
	rand.Seed(time.Now().UnixNano())
	nums := generateRandomSlice(1e6)

	for _, numGoroutines := range []int{1, 5, 10, 20, 50, 100} {
		b.Run("Mutex_Goroutines_"+strconv.Itoa(numGoroutines), func(b *testing.B) {
			sharedCounter := &SharedCounterMutex{}
			b.ResetTimer()
			for n := 0; n < b.N; n++ {
				start := time.Now()
				parallelCounterOperations(nums, numGoroutines, sharedCounter)
				elapsed := time.Since(start)
				throughput := float64(len(nums)) / elapsed.Seconds()
				fmt.Printf("Mutex: Goroutines: %d, Throughput: %.2f ops/sec, Latency: %.2f µs/op\n",
					numGoroutines, throughput, elapsed.Seconds()*1e6/float64(len(nums)))
			}
		})

		b.Run("Atomic_Goroutines_"+strconv.Itoa(numGoroutines), func(b *testing.B) {
			sharedCounter := &SharedCounterAtomic{}
			b.ResetTimer()
			for n := 0; n < b.N; n++ {
				start := time.Now()
				parallelCounterOperations(nums, numGoroutines, sharedCounter)
				elapsed := time.Since(start)
				throughput := float64(len(nums)) / elapsed.Seconds()
				fmt.Printf("Atomic: Goroutines: %d, Throughput: %.2f ops/sec, Latency: %.2f µs/op\n",
					numGoroutines, throughput, elapsed.Seconds()*1e6/float64(len(nums)))
			}
		})
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	nums := generateRandomSlice(1e6)

	sharedCounter := &SharedCounterMutex{}
	parallelCounterOperations(nums, 10, sharedCounter)
	fmt.Println("Final Counter (Mutex):", sharedCounter.Count())

	sharedCounterAtomic := &SharedCounterAtomic{}
	parallelCounterOperations(nums, 10, sharedCounterAtomic)
	fmt.Println("Final Counter (Atomic):", sharedCounterAtomic.Count())
}
