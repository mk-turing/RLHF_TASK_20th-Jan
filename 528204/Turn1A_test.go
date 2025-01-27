package _28204

import (
	"fmt"
	"net/http"
	"testing"
)

func BenchmarkAppPerformance(b *testing.B) {
	client := &http.Client{}
	url := "https://randomuser.me/"

	b.Run("LoadTest", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			resp, err := client.Get(url)
			if err != nil {
				b.Fatalf("Failed to reach the application: %v", err)
			}
			resp.Body.Close()
			b.ReportMetric(float64(resp.StatusCode), "status_code")
		}
	})
}

func main() {
	// Hook up testing.B to execute the benchmark as a normal test run.
	bench := testing.Benchmark(BenchmarkAppPerformance)
	fmt.Println(bench)
}
