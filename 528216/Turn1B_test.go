package _28216

import (
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func LoadConfig(filename string) ([]byte, error) {
	start := time.Now()
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	end := time.Now()
	duration := end.Sub(start)
	return data, nil
}

func BenchmarkLoadConfigSmall(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = LoadConfig("small_config.txt")
	}
}

func BenchmarkLoadConfigMedium(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = LoadConfig("medium_config.txt")
	}
}

func BenchmarkLoadConfigLarge(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = LoadConfig("large_config.txt")
	}
}

func main() {
	filenames := []string{"small_config.txt", "medium_config.txt", "large_config.txt"}

	for _, filename := range filenames {
		data, err := ioutil.ReadFile(filename)
		if err != nil {
			panic(err)
		}
		start := time.Now()
		for i := 0; i < 1000; i++ {
			_, _ = LoadConfig(filename)
		}
		end := time.Now()
		duration := end.Sub(start)
		println("Loaded", filename, len(data), "bytes in", duration)
	}

	os.Exit(0)
}
