package _28216

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"gopkg.in/yaml.v3"
)

// LoadConfig supports loading and parsing configuration files in JSON, YAML, and XML formats.
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

// BenchmarkLoadConfig benchmarks the LoadConfig function with various formats and sizes.
func BenchmarkLoadConfig(b *testing.B) {
	formats := []string{"json", "yaml", "xml"}
	sizes := map[string]int{
		"small":  1024,     // 1 KB
		"medium": 1048576,  // 1 MB
		"large":  10485760, // 10 MB
	}

	for _, format := range formats {
		for name, size := range sizes {
			b.Run(fmt.Sprintf("%s-%s", format, name), func(b *testing.B) {
				// Create a temporary file in the specified format and size.
				tempFile, err := ioutil.TempFile("", fmt.Sprintf("config_%s_%s.*.%s", format, name, format))
				if err != nil {
					b.Fatalf("Failed to create temp file: %s", err)
				}
				defer os.Remove(tempFile.Name())

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

				// Run the benchmark.
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

// generateDummyData creates dummy data with a specific size.
func generateDummyData(size int) map[string]interface{} {
	// Reduce the complexity by adjusting the filler size based on the requested size.
	data := make(map[string]interface{})
	filler := []byte("x")
	numKeys := size / 100 // Create approximately 100 bytes per key.
	if numKeys == 0 {
		numKeys = 1 // Ensure at least one key is generated for very small sizes.
	}
	for i := 0; i < numKeys; i++ {
		data[fmt.Sprintf("key%d", i)] = filler
	}
	return data
}

// generateXMLDummyData creates dummy XML data with a specific size.
func generateXMLDummyData(size int) interface{} {
	// Ensure the generated XML has a root element and well-formed structure.
	type Item struct {
		XMLName xml.Name `xml:"item"`
		Value   string   `xml:",chardata"`
	}

	// Create XML data with a root element
	root := struct {
		XMLName xml.Name `xml:"root"`
		Items   []Item   `xml:"item"`
	}{
		Items: make([]Item, size/100), // Adjust the number of items based on size
	}

	for i := 0; i < len(root.Items); i++ {
		root.Items[i] = Item{
			Value: fmt.Sprintf("value%d", i),
		}
	}

	return root
}
