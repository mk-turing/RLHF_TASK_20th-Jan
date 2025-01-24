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

// Dynamic workload generator
func dynamicWorkloadGenerator(numGoroutines int, workloadIntensity int, sharedCounter *SharedCounterAtomic, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < workloadIntensity; i++ {
		if rand.Intn(2) == 0 {
			sharedCounter.Increment()
		} else {
			sharedCounter.Decrement()
		}
	}
}

// Benchmark with dynamic workloads
func BenchmarkDynamicWorkload(b *testing.B) {
	rand.Seed(time.Now().UnixNano())
	sharedCounter := &SharedCounterAtomic{}
	workloadScenarios := []struct {
		numGoroutines     int
		workloadIntensity int
	}{
		{5, 1000},    // Low workload
		{50, 5000},   // Medium workload
		{100, 10000}, // High workload
	}

	for _, scenario := range workloadScenarios {
		b.Run("Goroutines_"+strconv.Itoa(scenario.numGoroutines), func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				var wg sync.WaitGroup
				wg.Add(scenario.numGoroutines)

				start := time.Now()
				for i := 0; i < scenario.numGoroutines; i++ {
					go dynamicWorkloadGenerator(scenario.numGoroutines, scenario.workloadIntensity, sharedCounter, &wg)
				}
				wg.Wait()
				elapsed := time.Since(start)

				throughput := float64(scenario.workloadIntensity*scenario.numGoroutines) / elapsed.Seconds()
				fmt.Printf("Goroutines: %d, Throughput: %.2f ops/sec, Latency: %.2f µs/op\n",
					scenario.numGoroutines, throughput, elapsed.Seconds()*1e6/float64(scenario.workloadIntensity*scenario.numGoroutines))
			}
		})
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Simulate dynamic workloads
	var wg sync.WaitGroup
	nums := []int{10, 100, 1000}
	workloadIntensities := []int{1000, 5000, 10000}

	for _, numGoroutines := range nums {
		for _, intensity := range workloadIntensities {
			wg.Add(1)
			go func(num, intensity int) {
				defer wg.Done()
				sharedCounter := &SharedCounterAtomic{}
				start := time.Now()
				dynamicWorkloadGenerator(num, intensity, sharedCounter, &wg)
				elapsed := time.Since(start)
				fmt.Printf("Goroutines: %d, Intensity: %d, Time: %v\n", num, intensity, elapsed)
			}(numGoroutines, intensity)
		}
	}

	wg.Wait()
}
