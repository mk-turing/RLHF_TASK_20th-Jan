package _28194

import (
	"fmt"
	"math/rand"
	"runtime"
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

// Workload Generator
func workloadGenerator(b *testing.B, numGoroutinesCh chan int, intensityCh chan int) {
	for i := 0; i < b.N; i++ {
		// Generate random number of goroutines between 1 and 100
		numGoroutines := rand.Intn(99) + 1
		numGoroutinesCh <- numGoroutines

		// Generate random intensity between 1 and 100
		intensity := rand.Intn(99) + 1
		intensityCh <- intensity

		time.Sleep(time.Millisecond * 100) // Generate new workload every 100 ms
	}
	close(numGoroutinesCh)
	close(intensityCh)
}

// Benchmark with Dynamic Workload
func BenchmarkParallelCounterDynamicWorkload(b *testing.B) {
	rand.Seed(time.Now().UnixNano())
	nums := generateRandomSlice(1e6)