package _28204

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"testing"
	"time"
)

const ephemeralStorageDir = "/tmp/ephemeral_storage" // Temporary directory for ephemeral storage
const fileSize int = 1024 * 1024 // 1MB file size

func init() {
	os.MkdirAll(ephemeralStorageDir, 0777) // Create ephemeral storage directory
}

func BenchmarkAppHighConcurrencyWithStorage(b *testing.B) {
	client := &http.Client{}
	url := "https://example.com/data"

	concurrencyLevels := []int{10, 50, 100} // Number of goroutines to simulate

	for _, concurrency := range concurrencyLevels {
		b.Run(fmt.Sprintf("Concurrency-%d", concurrency), func(b *testing.B) {
			var wg sync.WaitGroup
			var memStats runtime.MemStats

			runtime.GC() // Run garbage collection before the benchmark
			runtime.ReadMemStats(&memStats)
			fmt.Printf("Initial Memory Usage: %v KB\n", memStats.Alloc/1024)

			// Start a go routine to simulate frequent file reads and writes
			go func() {
				for {
					buf := bytes.Buffer{}
					// Simulate writing a large file
					if _, err := buf.WriteString(string(make([]byte, fileSize))); err != nil {
						b.Fatalf("Error writing to buffer: %v", err)
					}

					// Choose a random file name for persistence
					filename := filepath.Join(ephemeralStorageDir, fmt.Sprintf("%d.txt", rand.Int63()))
					if err := ioutil.WriteFile(filename, buf.Bytes(), 0666); err != nil {
						b.Fatalf("Error writing file to disk: %v", err)
					}

					// Simulate reading the file from ephemeral storage
					data, err := ioutil.ReadFile(filename)
					if err != nil {
						b.Fatalf("Error reading file from disk: %v", err)
					}
					if bytes.Compare(data, buf.Bytes()) != 0 {
						b.Fatalf("Data read from file does not match written data")
					}