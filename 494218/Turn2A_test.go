package _94218

import (
	"io/ioutil"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

type CacheStrategy interface {
	Write(data []byte, filename string) error
	Read(filename string) ([]byte, error)
}

type InMemoryCache struct {
	sync.RWMutex
	store map[string][]byte
}

func (c *InMemoryCache) Write(data []byte, filename string) error {
	c.Lock()
	defer c.Unlock()
	c.store[filename] = data
	return nil
}

func (c *InMemoryCache) Read(filename string) ([]byte, error) {
	c.RLock()
	defer c.RUnlock()
	data, exists := c.store[filename]
	if !exists {
		return nil, os.ErrNotExist
	}
	return data, nil
}

// Adding concurrency support and metrics measurement
func concurrentAccessPattern(b *testing.B, files []string, strategy CacheStrategy) {
	var wg sync.WaitGroup
	var cacheHits int32

	// Simulate multiple goroutines
	numGoroutines := 10
	for g := 0; g < numGoroutines; g++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < b.N; i++ {
				filename := files[i%len(files)]
				data := make([]byte, 1024)

				if err := strategy.Write(data, filename); err != nil {
					b.Error(err)
				}

				readData, err := strategy.Read(filename)
				if err == nil && len(readData) == len(data) {
					atomic.AddInt32(&cacheHits, 1)
				} else if err != os.ErrNotExist {
					b.Error(err)
				}
			}
		}()
	}

	wg.Wait()
	hitRate := float64(cacheHits) / float64(b.N*numGoroutines)
	b.ReportMetric(hitRate, "cache-hit-rate")
}

func setupFiles(b *testing.B, fileSize int) []string {
	files := make([]string, 0, b.N)
	for i := 0; i < b.N; i++ {
		// Use counter and time to ensure unique filenames
		filename := time.Now().Format("20060102150405") + "-" + strconv.Itoa(i)
		files = append(files, filename)
	}
	return files
}

type FileSystemCache struct {
	directory string
}

func (f *FileSystemCache) Write(data []byte, filename string) error {
	return ioutil.WriteFile(f.directory+"/"+filename, data, 0644)
}

func (f *FileSystemCache) Read(filename string) ([]byte, error) {
	return ioutil.ReadFile(f.directory + "/" + filename)
}

func BenchmarkInMemoryCacheConcurrentAccess(b *testing.B) {
	cache := &InMemoryCache{store: make(map[string][]byte)}
	files := setupFiles(b, 1024)
	concurrentAccessPattern(b, files, cache)
}

func BenchmarkFileSystemCacheConcurrentAccess(b *testing.B) {
	// Create a temporary directory for filesystem cache
	dir, err := ioutil.TempDir("", "fs-benchmark")
	if err != nil {
		b.Fatal(err)
	}
	defer os.RemoveAll(dir)

	cache := &FileSystemCache{directory: dir}
	files := setupFiles(b, 1024)
	concurrentAccessPattern(b, files, cache)
}
