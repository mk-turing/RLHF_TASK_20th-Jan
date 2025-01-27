package _28204

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime"
	"runtime/pprof"
	"sync"
	"testing"
	"time"
)

// Simulate high-concurrency scenarios with dynamic scaling and memory tracking
func BenchmarkAppHighConcurrency(b *testing.B) {
	client := &http.Client{}
	url := "https://randomuser.me/"

	// Simulate dynamic scaling with varying concurrency levels
	concurrencyLevels := []int{100, 500, 1000}         // Number of goroutines to simulate
	workloadPatterns := []string{"burst", "sustained"} // Workload patterns to simulate

	// Memory profiling setup
	memProfile := "mem.prof"
	cpuProfile := "cpu.prof"

	for _, workload := range workloadPatterns {
		for _, concurrency := range concurrencyLevels {
			b.Run(fmt.Sprintf("Workload-%s-Concurrency-%d", workload, concurrency), func(b *testing.B) {
				var wg sync.WaitGroup
				var memStats runtime.MemStats

				// Start CPU profiling
				f, err := startCPUProfiling(cpuProfile)
				if err != nil {
					b.Fatalf("Failed to start CPU profiling: %v", err)
				}

				start := time.Now()

				// Run concurrent requests
				for i := 0; i < concurrency; i++ {
					wg.Add(1)
					go func(workerID int) {
						defer wg.Done()
						for j := 0; j < b.N; j++ {
							resp, err := client.Get(url)
							if err != nil {
								b.Errorf("Worker %d: failed to reach the application: %v", workerID, err)
							} else {
								resp.Body.Close()
							}
							if workload == "burst" && j%10 == 0 {
								// Simulate a short burst by sleeping
								time.Sleep(100 * time.Millisecond)
							}
						}
					}(i)
				}
				wg.Wait()

				// Stop profiling and print memory stats
				stopCPUProfiling(f)
				runtime.ReadMemStats(&memStats)

				fmt.Printf("Test completed in %v\n", time.Since(start))
				fmt.Printf("Memory Allocations: %v KB\n", memStats.Alloc/1024)
				fmt.Printf("Total Allocations: %v KB\n", memStats.TotalAlloc/1024)
				fmt.Printf("Heap Allocations: %v KB\n", memStats.HeapAlloc/1024)
				fmt.Printf("Number of Goroutines: %v\n", runtime.NumGoroutine())

				// Write memory profile to file
				err = writeMemProfile(memProfile)
				if err != nil {
					b.Errorf("Failed to write memory profile: %v", err)
				}
			})
		}
	}
}

func startCPUProfiling(profilePath string) (*bytes.Buffer, error) {
	var buf bytes.Buffer
	err := pprof.StartCPUProfile(&buf)
	if err != nil {
		return nil, fmt.Errorf("error starting CPU profile: %v", err)
	}
	return &buf, nil
}

func stopCPUProfiling(f *bytes.Buffer) {
	pprof.StopCPUProfile()
	fmt.Println("CPU profile written successfully.")
}

func writeMemProfile(profilePath string) error {
	err := pprof.Lookup("heap").WriteTo(&bytes.Buffer{}, 0)
	if err != nil {
		return fmt.Errorf("error writing memory profile: %v", err)
	}
	fmt.Println("Memory profile written successfully.")
	return nil
}

func main() {
	// Execute the benchmark
	bench := testing.Benchmark(BenchmarkAppHighConcurrency)
	fmt.Println(bench)
}
