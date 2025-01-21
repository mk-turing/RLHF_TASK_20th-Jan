package _28216

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"
)

// LoadConfig loads a configuration file and returns its content as a string.
// It measures and logs the time taken to load the file.
func LoadConfig(filename string) (string, error) {
	start := time.Now()

	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	content, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}

	duration := time.Since(start)
	log.Printf("Loaded file: %s, Size: %d bytes, Time taken: %v\n", filename, len(content), duration)

	return string(content), nil
}

// BenchmarkLoadConfig benchmarks the LoadConfig function with various file sizes.
func BenchmarkLoadConfig(b *testing.B) {
	sizes := map[string]int{
		"small":  1024,     // 1 KB
		"medium": 1048576,  // 1 MB
		"large":  10485760, // 10 MB
	}

	for name, size := range sizes {
		// Create a temporary file of the specified size.
		tempFile, err := ioutil.TempFile("", "config_"+name+"_*.tmp")
		if err != nil {
			b.Fatalf("Failed to create temp file: %s", err)
		}
		defer os.Remove(tempFile.Name())

		// Write random data to the file.
		data := make([]byte, size)
		if _, err := tempFile.Write(data); err != nil {
			b.Fatalf("Failed to write to temp file: %s", err)
		}
		if err := tempFile.Close(); err != nil {
			b.Fatalf("Failed to close temp file: %s", err)
		}

		b.Run(name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, err := LoadConfig(tempFile.Name())
				if err != nil {
					b.Fatalf("Failed to load config: %s", err)
				}
			}
		})
	}
}
