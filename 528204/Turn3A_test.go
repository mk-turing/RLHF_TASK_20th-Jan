package _28204

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

const (
	benchmarkFile = "./cloud-benchmark-test/benchmark-file.txt"
	dataSize      = 10485760 // 10 MB
)

func BenchmarkStorageIntensive(b *testing.B) {
	// Make sure the directory exists
	if err := os.MkdirAll("./cloud-benchmark-test", 0755); err != nil {
		b.Fatalf("Failed to create benchmark directory: %v", err)
	}

	// Open file to write to
	file, err := os.Create(benchmarkFile)
	if err != nil {
		b.Fatalf("Writing to disk failed: %v", err)
	}
	defer file.Close()

	// Write data to the file only once
	data := make([]byte, dataSize)
	_, err = file.Write(data)
	if err != nil {
		b.Fatalf("Error writing to file: %v", err)
	}

	// Reading from the file to verify
	readData, err := ioutil.ReadFile(benchmarkFile)
	if err != nil {
		b.Fatalf("Reading from disk failed: %v", err)
	}

	// Check the size of the file
	if len(readData) != dataSize {
		b.Fatalf("Read data mismatch: expected %d, got %d", dataSize, len(readData))
	}

	// Clean up
	err = os.Remove(benchmarkFile)
	if err != nil {
		b.Fatalf("Failed to remove temporary file: %v", err)
	}
}

func main() {
	// Run benchmark manually for testing purposes (for demonstration here)
	result := testing.Benchmark(BenchmarkStorageIntensive)
	fmt.Println(result)
}
