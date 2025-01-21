package _28216

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"sync"
	"testing"
)

// LoadConfigConcurrently loads configuration files concurrently using a worker pool pattern.
func LoadConfigConcurrently(filenames []string, numWorkers int) ([]interface{}, error) {
	var wg sync.WaitGroup
	results := make([]interface{}, len(filenames))
	errors := make(chan error, len(filenames))

	wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go func(workerID int) {
			defer wg.Done()
			for j := workerID; j < len(filenames); j += numWorkers {
				result, err := LoadConfig(filenames[j])
				if err != nil {
					errors <- fmt.Errorf("worker %d: failed to load config %s: %w", workerID, filenames[j], err)
					continue
				}
				results[j] = result
			}
		}(i)
	}

	wg.Wait()
	close(errors)

	for err := range errors {
		return nil, err
	}

	return results, nil
}

// BenchmarkLoadConfigConcurrently benchmarks the LoadConfigConcurrently function with varying file sizes and formats.
func BenchmarkLoadConfigConcurrently(b *testing.B) {
	formats := []string{"json", "yaml", "xml"}
	sizes := map[string]int{
		"small":  1024,     // 1 KB
		"medium": 1048576,  // 1 MB
		"large":  10485760, // 10 MB
	}
	numWorkers := 10

	for _, format := range formats {
		for name, size := range sizes {
			b.Run(fmt.Sprintf("%s-%s-sequential", format, name), func(b *testing.B) {
				filenames := createTempFiles(b, format, name, size, 1)
				defer cleanTempFiles(filenames)

				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					_, err := LoadConfigSequentially(filenames)
					if err != nil {
						b.Fatalf("Failed to load config: %s", err)
					}
				}
			})

			b.Run(fmt.Sprintf("%s-%s-concurrent", format, name), func(b *testing.B) {
				filenames := createTempFiles(b, format, name, size, numWorkers)
				defer cleanTempFiles(filenames)

				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					_, err := LoadConfigConcurrently(filenames, numWorkers)
					if err != nil {
						b.Fatalf("Failed to load config: %s", err)
					}
				}
			})
		}
	}
}

// createTempFiles creates temporary files in the specified format and size for benchmarking.
func createTempFiles(b *testing.B, format, name string, size int, numWorkers int) []string {
	var filenames []string
	for i := 0; i < numWorkers; i++ {
		tempFile, err := ioutil.TempFile("", fmt.Sprintf("config_%s_%s_%d.*.%s", format, name, i, format))
		if err != nil {
			b.Fatalf("Failed to create temp file: %s", err)
		}
		filenames = append(filenames, tempFile.Name())

		// Generate content based on format.
		var data []byte
		switch format {
		case "json":
			data, _ = json.Marshal(generateDummyData(size))
		case "yaml":
			data, _ = yaml.Marshal(generateDummyData(size))
		case "xml":
			data, _ = xml.Marshal(generateXMLDummyData(size))
		}

		if _, err := tempFile.Write(data); err != nil {
			b.Fatalf("Failed to write to temp file: %s", err)
		}
		if err := tempFile.Close(); err != nil {
			b.Fatalf("Failed to close temp file: %s", err)
		}
	}
	return filenames
}

// cleanTempFiles removes the temporary files created for benchmarking.
func cleanTempFiles(filenames []string) {
	for _, filename := range filenames {
		os.Remove(filename)
	}
}

// LoadConfigSequentially loads configuration files sequentially.
func LoadConfigSequentially(filenames []string) ([]interface{}, error) {
	var results []interface{}
	for _, filename := range filenames {
		result, err := LoadConfig(filename)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	return results, nil
}
