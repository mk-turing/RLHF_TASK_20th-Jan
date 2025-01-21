package _28216

import (
	"bufio"
	"bytes"
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

// StreamConfigLoader reads and parses config files in blocks to reduce memory usage.
func StreamConfigLoader(filename string) (interface{}, error) {
	start := time.Now()

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	var result interface{}

	switch {
	case hasExtension(filename, ".json"):
		err = json.NewDecoder(reader).Decode(&result)
	case hasExtension(filename, ".yaml", ".yml"):
		// YAML streaming in Go typically loads everything; you'd ideally use a streaming YAML parser here.
		buf := new(bytes.Buffer)
		_, err = buf.ReadFrom(reader)
		if err == nil {
			err = yaml.Unmarshal(buf.Bytes(), &result)
		}
	case hasExtension(filename, ".xml"):
		err = xml.NewDecoder(reader).Decode(&result)
	default:
		return nil, fmt.Errorf("unsupported file format for file: %s", filename)
	}

	if err != nil {
		return nil, err
	}

	duration := time.Since(start)
	log.Printf("Loaded and parsed file: %s, Size: N/A (stream), Time taken: %v\n", filename, duration)

	return result, nil
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

// LoadConfigsConcurrentlyWithStream loads multiple configuration files concurrently using streaming.
func LoadConfigsConcurrentlyWithStream(filenames []string) ([]interface{}, error) {
	var wg sync.WaitGroup
	results := make([]interface{}, len(filenames))
	errors := make(chan error, len(filenames))

	for i, filename := range filenames {
		wg.Add(1)
		go func(i int, filename string) {
			defer wg.Done()
			result, err := StreamConfigLoader(filename)
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

// BenchmarkLoadConfigsConcurrentlyWithStream benchmarks concurrent loading using streaming.
func BenchmarkLoadConfigsConcurrentlyWithStream(b *testing.B) {
	formats := []string{"json", "yaml", "xml"}
	sizes := map[string]int{
		"small":  1024,
		"medium": 1048576,
		"large":  10485760,
	}
	numFiles := 5

	for _, format := range formats {
		for name, size := range sizes {
			b.Run(fmt.Sprintf("ConcurrentStream-%s-%s", format, name), func(b *testing.B) {
				var tempFiles []string
				for i := 0; i < numFiles; i++ {
					tempFile := createTempFile(b, format, fmt.Sprintf("%s%d", name, i), size)
					defer os.Remove(tempFile.Name())
					tempFiles = append(tempFiles, tempFile.Name())
				}

				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					_, err := LoadConfigsConcurrentlyWithStream(tempFiles)
					if err != nil {
						b.Fatalf("Failed to load configs concurrently with stream: %s", err)
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
