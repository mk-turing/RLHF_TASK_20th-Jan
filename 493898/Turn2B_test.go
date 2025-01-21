package _93898

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"runtime"
	"sync"
	"testing"
)

// Sample function for concurrent computation
func ConcurrentSampleFunction(n int, wg *sync.WaitGroup) {
	defer wg.Done()
	sum := 0
	for i := 0; i < n; i++ {
		sum += i
	}
}

// Benchmarking the concurrent sample function using the 'testing' package
func BenchmarkConcurrentSampleFunction(b *testing.B) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			var wg sync.WaitGroup
			wg.Add(b.N)
			for i := 0; i < b.N; i++ {
				go ConcurrentSampleFunction(1000, &wg)
			}
			wg.Wait()
		}
	})
}

// Benchmarking the concurrent sample function using the 'testify' package
func BenchmarkConcurrentSampleFunctionWithTestify(b *testing.B) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			var wg sync.WaitGroup
			wg.Add(b.N)
			for i := 0; i < b.N; i++ {
				go func() {
					defer wg.Done()
					result := SampleFunction(1000)
					assert.NotNil(b, result)
					assert.Greater(b, result, 0)
				}()
			}
			wg.Wait()
		}
	})
}

func main() {
	fmt.Println("Running main function...")
	// Sample function call for non-benchmarking purposes
	fmt.Println(SampleFunction(1000))
}
