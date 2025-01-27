package main

import (
	"net/http"
	"testing"
	"time"
)

func BenchmarkRequestLatency(b *testing.B) {
	for n := 0; n < b.N; n++ {
		start := time.Now()
		resp, err := http.Get("http://your-api-endpoint.com/resource") // Replace with your real application endpoint
		if err != nil {
			b.Fatal(err)
		}
		_ = resp.Body.Close()
		elapsed := time.Since(start)
		b.ReportMetric(float64(elapsed.Nanoseconds()), "ns/op")
	}
}

func BenchmarkThroughput(b *testing.B) {
	start := time.Now()
	for n := 0; n < b.N; n++ {
		resp, err := http.Get("http://your-api-endpoint.com/resource") // Replace with your real application endpoint
		if err != nil {
			b.Fatal(err)
		}
		_ = resp.Body.Close()
	}
	elapsed := time.Since(start)
	b.ReportMetric(float64(b.N)/elapsed.Seconds(), "ops/s")
}
