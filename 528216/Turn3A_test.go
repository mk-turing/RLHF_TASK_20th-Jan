package _28216

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"gopkg.in/yaml.v3"
)

// LoadConfig supports loading and parsing configuration files in JSON, YAML, and XML formats.
// It returns the loaded configuration or an error.
func LoadConfig(filename string) (interface{}, error) {
	start := time.Now()

	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var result interface{}
	switch {
	case hasExtension(filename, ".json"):
		err = json.Unmarshal(content, &result)
	case hasExtension(filename, ".yaml", ".yml"):
		err = yaml.Unmarshal(content, &result)
	case hasExtension(filename, ".xml"):
		err = xml.Unmarshal(content, &result)
	default:
		return nil, fmt.Errorf("unsupported file format for file: %s", filename)
	}
	if err != nil {
		return nil, err
	}

	duration := time.Since(start)
	log.Printf("Loaded and parsed file: %s, Size: %d bytes, Time taken: %v\n", filename, len(content), duration)

	return result, nil
}

// LoadConfigsConcurrently loads multiple configuration files concurrently.
func LoadConfigsConcurrently(filenames []string) ([]interface{}, error) {
	var wg sync.WaitGroup
	results := make([]interface{}, len(filenames))
	errors := make(chan error, len(filenames))

	for i, filename := range filenames {
		wg.Add(1)
		go func(i int, filename string) {
			defer wg.Done()
			result, err := LoadConfig(filename)
			if err != nil {
				errors <- err
				return
			}
			results[i] = result
		}(i, filename)
	}

	wg.Wait()
	close(errors)

	if len(errors) > 0 {
		return nil, <-errors
	}
	return results, nil
}

// hasExtension checks if a file has any of the provided extensions.
func hasExtension(filename string, extensions ...string) bool {
	ext := strings.ToLower(filename)
	for _, e := range extensions {
		if len(ext) >= len(e) && ext[len(ext)-len(e):] == e {
			return true
		}
	}
	return false
}

// BenchmarkLoadConfig benchmarks loading configurations sequentially.
func BenchmarkLoadConfig(b *testing.B) {
	formats := []string{"json", "yaml", "xml"}
	sizes := map[string]int{
		"small":  1024,     // 1 KB
		"medium": 1048576,  // 1 MB
		"large":  10485760, // 10 MB
	}

	for _, format := range formats {
		for name, size := range sizes {
			b.Run(fmt.Sprintf("Sequential-%s-%s", format, name), func(b *testing.B) {
				tempFile := createTempFile(b, format, name, size)
				defer os.Remove(tempFile.Name())

				// Run the benchmark sequentially.
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					_, err := LoadConfig(tempFile.Name())
					if err != nil {
						b.Fatalf("Failed to load config: %s", err)
					}
				}
			})
		}
	}
}

// BenchmarkLoadConfigsConcurrently benchmarks loading configurations concurrently.
func BenchmarkLoadConfigsConcurrently(b *testing.B) {
	formats := []string{"json", "yaml", "xml"}
	sizes := map[string]int{
		"small":  1024,     // 1 KB
		"medium": 1048576,  // 1 MB
		"large":  10485760, // 10 MB
	}
	numFiles := 5 // Number of files to load concurrently.

	for _, format := range formats {
		for name, size := range sizes {
			b.Run(fmt.Sprintf("Concurrent-%s-%s", format, name), func(b *testing.B) {
				// Create temporary files.
				var tempFiles []string
				for i := 0; i < numFiles; i++ {
					tempFile := createTempFile(b, format, fmt.Sprintf("%s%d", name, i), size)
					defer os.Remove(tempFile.Name())
					tempFiles = append(tempFiles, tempFile.Name())
				}

				// Run the benchmark concurrently.
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					_, err := LoadConfigsConcurrently(tempFiles)
					if err != nil {
						b.Fatalf("Failed to load configs concurrently: %s", err)
					}
				}
			})
		}
	}
}

// createTempFile creates a temporary file filled with dummy data of a specific format and size.
func createTempFile(b *testing.B, format, name string, size int) *os.File {
	tempFile, err := ioutil.TempFile("", fmt.Sprintf("config_%s_%s.*.%s", format, name, format))
	if err != nil {
		b.Fatalf("Failed to create temp file: %s", err)
	}

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

	return tempFile
}

// generateDummyData creates dummy data with a specific size for JSON and YAML.
func generateDummyData(size int) map[string]interface{} {
	data := make(map[string]interface{})
	filler := []byte("x")
	numKeys := size / 100 // approx 100 bytes per key.
	if numKeys == 0 {
		numKeys = 1 // ensure at least one key.
	}
	for i := 0; i < numKeys; i++ {
		data[fmt.Sprintf("key%d", i)] = filler
	}
	return data
}

// generateXMLDummyData creates dummy XML data with a specific size.
func generateXMLDummyData(size int) interface{} {
	type Item struct {
		XMLName xml.Name `xml:"item"`
		Value   string   `xml:",chardata"`
	}

	root := struct {
		XMLName xml.Name `xml:"root"`
		Items   []Item   `xml:"item"`
	}{
		Items: make([]Item, size/100), // adjust items based on size
	}

	for i := 0; i < len(root.Items); i++ {
		root.Items[i] = Item{
			Value: fmt.Sprintf("value%d", i),
		}
	}

	return root
}
