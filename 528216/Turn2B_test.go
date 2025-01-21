package _28216

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Key   string `json:"key" yaml:"key" xml:"key"`
	Value string `json:"value" yaml:"value" xml:"value"`
}

func LoadConfig(filename string) (*Config, error) {
	start := time.Now()

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	content, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var config Config
	var format string
	switch ext := filepath.Ext(filename); ext {
	case ".json":
		format = "JSON"
		err = json.Unmarshal(content, &config)
	case ".yaml", ".yml":
		format = "YAML"
		err = yaml.Unmarshal(content, &config)
	case ".xml":
		format = "XML"
		err = xml.Unmarshal(content, &config)
	default:
		return nil, fmt.Errorf("unsupported file format: %s", ext)
	}

	duration := time.Since(start)
	log.Printf("Loaded file: %s, Format: %s, Size: %d bytes, Time taken: %v\n", filename, format, len(content), duration)

	return &config, err
}

func BenchmarkLoadConfig(b *testing.B) {
	sizes := map[string]int{
		"small":  1024,     // 1 KB
		"medium": 1048576,  // 1 MB
		"large":  10485760, // 10 MB
	}

	formats := []string{"json", "yaml", "xml"}

	for _, format := range formats {
		for name, size := range sizes {
			tempFile, err := ioutil.TempFile("", fmt.Sprintf("config_%s_%s_*.tmp", name, format))
			if err != nil {
				b.Fatalf("Failed to create temp file: %s", err)
			}
			defer os.Remove(tempFile.Name())

			data := make([]byte, size)
			if _, err := tempFile.Write(data); err != nil {
				b.Fatalf("Failed to write to temp file: %s", err)
			}
			if err := tempFile.Close(); err != nil {
				b.Fatalf("Failed to close temp file: %s", err)
			}

			b.Run(fmt.Sprintf("%s_%s", format, name), func(b *testing.B) {
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
