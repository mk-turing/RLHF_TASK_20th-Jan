package _28194

import (
	"math/rand"
	"runtime"
	"strconv"
	"sync"
	"testing"
	"time"
)

// Parallel sum function
func parallelSum(nums []int, numGoroutines int) int {
	chunkSize := (len(nums) + numGoroutines - 1) / numGoroutines
	var sum int
	var wg sync.WaitGroup
	var mu sync.Mutex

	for i := 0; i < len(nums); i += chunkSize {
		wg.Add(1)
		go func(start int) {
			defer wg.Done()
			var localSum int
			end := start + chunkSize
			if end > len(nums) {
				end = len(nums)
			}
			for _, v := range nums[start:end] {
				localSum += v
			}
			mu.Lock()
			sum += localSum
			mu.Unlock()
		}(i)
	}

	wg.Wait()
	return sum
}

// Generate random slice
func generateRandomSlice(size int) []int {
	nums := make([]int, size)
	for i := range nums {
		nums[i] = rand.Intn(100)
	}
	return nums
}

// Benchmark for parallelSum with different goroutine counts
func BenchmarkParallelSum(b *testing.B) {
	nums := generateRandomSlice(1e6)
	for _, numGoroutines := range []int{1, 5, 10, 20, 50, 100} {
		b.Run("Goroutines_"+strconv.Itoa(numGoroutines), func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				_ = parallelSum(nums, numGoroutines)
			}
		})
	}
}

// Measure scheduling effects
func measureSchedulingEffects() {
	nums := generateRandomSlice(1e6)
	numGoroutines := 10

	for _, maxProcs := range []int{1, 2, runtime.NumCPU()} {
		runtime.GOMAXPROCS(maxProcs)
		start := time.Now()
		parallelSum(nums, numGoroutines)
		elapsed := time.Since(start)
		println("GOMAXPROCS:", maxProcs, "Time:", elapsed.Milliseconds(), "ms")
	}
}
