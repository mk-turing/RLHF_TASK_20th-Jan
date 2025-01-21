package _93898

import (
	"flag"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"sync"
	"testing"
)

var (
	configFile   string
	funcSize     int
	concurrency  int
	iterations   int
	useTestify   bool
)

func init() {
	flag.StringVar(&configFile, "config", "", "Configuration file path")
	flag.IntVar(&funcSize, "size", 1000, "Function size (number of elements)")
	flag.IntVar(&concurrency, "concurrency", 4, "Concurrency level")
	flag.IntVar(&iterations, "iterations", 100000, "Number of iterations")
	flag.BoolVar(&useTestify, "testify", false, "Use testify framework for benchmarking")
	flag.Parse()

	// Load configuration from file if provided
	if configFile != "" {
		viper.SetConfigFile(configFile)
		if err := viper.ReadInConfig(); err != nil {
			log.Fatalf("Error reading config file: %s", err)
		}
		viper.Unmarshal(&funcSize)
		viper.Unmarshal(&concurrency)
		viper.Unmarshal(&iterations)
		viper.Unmarshal(&useTestify)
	}
}

func SampleFunction(n int) int {
	sum := 0
	for i := 0; i < n; i++ {
		sum += i
	}
	return sum
}

func ConcurrentSum(numbers []int) int {
	const partSize = 100
	results := make(chan int, concurrency)
	var wg sync.WaitGroup

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(start int) {
			defer wg.Done()
			end := start + partSize
			if end > len(numbers) {
				end = len(numbers)
			}
			sum := 0
			for _, num := range numbers[start:end] {
				sum += num
			}
			results <- sum
		}(i * partSize)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	finalSum := 0
	for sum := range results {
		finalSum += sum
	}

	return finalSum
}

func BenchmarkSampleFunction(b *testing.B) {
	numbers := make([]int, funcSize)
	for i := range numbers {
		numbers[i] = i
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		SampleFunction(funcSize)
	}
}

func BenchmarkConcurrentSum(b *testing.B) {
	numbers := make([]int, funcSize)
	for i := range numbers {
		numbers[i] = i
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ConcurrentSum(numbers)
	}
}

func main() {
	fmt.Printf("Benchmarking with: Function size = %d, Concurrency = %d, Iterations = %d, UseTestify = %t\n",
		funcSize, concurrency, iterations, useTestify)

	if useTestify {
		// Run benchmarks using testify