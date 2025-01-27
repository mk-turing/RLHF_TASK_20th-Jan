package _28204

import (
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

const (
	minConcurrency     = 100
	maxConcurrency     = 1000
	concurrencyStep    = 100
	maxRequestsPerSec  = 100
	requestBurstSize   = 100
	maxIdleConnections = 200
)

var client = &http.Client{
	Transport: &http.Transport{
		MaxIdleConnsPerHost: maxIdleConnections,
	},
}

func simulateLoad(wg *sync.WaitGroup, b *testing.B, concurrency int, rate float64, btime time.Duration) {
	defer wg.Done()
	ticker := time.NewTicker(time.Duration(float64(time.Second) / rate))
	defer ticker.Stop()
	var respCount uint64

	for range ticker.C {
		for i := 0; i < requestBurstSize; i++ {
			go func() {
				resp, err := client.Get("http://your-cloud-application-url.com/health")
				if err != nil {
					b.Fatalf("Failed to reach the application: %v", err)
				}
				resp.Body.Close()
				respCount++
			}()
		}
	}

	time.Sleep(btime)
	atomic.AddUint64(&totalRequests, respCount)
}

func BenchmarkAppPerformanceConcurrency(b *testing.B) {
	var totalRequests uint64
	wg := &sync.WaitGroup{}

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			concurrency := rand.Intn(maxConcurrency-minConcurrency+1) + minConcurrency
			rate := float64(concurrency) / maxRequestsPerSec
			btime := time.Duration(float64(time.Second) / rate)

			wg.Add(1)
			go simulateLoad(wg, b, concurrency, rate, btime)

			// Wait for a randomized amount of time to introduce various concurrency levels during the benchmark.
			time.Sleep(time.Duration(rand.Intn(1000) * int(time.Millisecond)))
		}
	})

	wg.Wait()
	b.ReportMetric(float64(totalRequests)/float64(b.N), "requests_per_second")
}

func main() {
	// Hook up testing.B to execute the benchmark as a normal test run.
	bench := testing.Benchmark(BenchmarkAppPerformanceConcurrency)
	fmt.Println(bench)
}
