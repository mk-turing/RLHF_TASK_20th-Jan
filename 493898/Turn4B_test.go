package _93898

import (
	"fmt"
	"runtime"
	"sync"
	"testing"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

type BenchmarkFunc func(int) (int, error)

type BenchmarkResult struct {
	Name      string
	Time      testing.BenchmarkResult
	CPUUsage  float64
	MemUsage  uint64
	Error     error
}

func BenchmarkFunctions(b *testing.B, funcs []BenchmarkFunc, n int) []BenchmarkResult {
	results := make([]BenchmarkResult, len(funcs))
	var wg sync.WaitGroup

	for i, fn := range funcs {
		wg.Add(1)
		go func(idx int, f BenchmarkFunc) {
			defer wg.Done()
			// Measure CPU and memory usage before and after the benchmark
			cpuPercent, _ := cpu.Percent(0, false)
			v, _ := mem.VirtualMemory()
			startMem := v.Used

			res := testing.Benchmark(func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					_, err := f(n)
					if err != nil {
						b.Fatalf("Error in benchmark: %v", err)
					}
				}
			})

			// Measure CPU and memory usage again after the benchmark
			cpuPercentAfter, _ := cpu.Percent(0, false)
			vAfter, _ := mem.VirtualMemory()
			endMem := vAfter.Used

			results[idx] = BenchmarkResult{
				Name:      fmt.Sprintf("%T", f),
				Time:      res,
				CPUUsage:  cpuPercentAfter[0] - cpuPercent[0],