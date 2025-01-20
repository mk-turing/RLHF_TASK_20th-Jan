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

type FileSystemCache struct {
	directory string
}

func (f *FileSystemCache) Write(data []byte, filename string) error {
	return ioutil.WriteFile(f.directory+"/"+filename, data, 0644)
}

func (f *FileSystemCache) Read(filename string) ([]byte, error) {
	return ioutil.ReadFile(f.directory + "/" + filename)
}

// Test suite demonstrating the performance evaluation
func benchmarkCacheStrategy(b *testing.B, strategy CacheStrategy, fileSize int, accessPattern func(b *testing.B, files []string, strategy CacheStrategy)) {
	files := setupFiles(b, fileSize)
	b.ResetTimer()
	accessPattern(b, files, strategy)
}

func setupFiles(b *testing.B, fileSize int) []string {
	files := make([]string, 0, b.N)
	for i := 0; i < b.N; i++ {
		filename := time.Now().Format("20060102150405") + "-" + strconv.Itoa(i)
		files = append(files, filename)
	}
	return files
}

func randomAccessPattern(b *testing.B, files []string, strategy CacheStrategy) {
	for i := 0; i < b.N; i++ {
		if err := strategy.Write(make([]byte, 1024), files[i%len(files)]); err != nil {
			b.Fatal(err)
		}
		if _, err := strategy.Read(files[i%len(files)]); err != nil && err != os.ErrNotExist {
			b.Fatal(err)
		}
	}
}

func sequentialAccessPattern(b *testing.B, files []string, strategy CacheStrategy) {
	for i := 0; i < b.N; i++ {
		filename := files[i%len(files)]
		if err := strategy.Write(make([]byte, 1024), filename); err != nil {
			b.Fatal(err)
		}
		if _, err := strategy.Read(filename); err != nil && err != os.ErrNotExist {
			b.Fatal(err)
		}
	}
}

func BenchmarkInMemoryCacheRandomAccess(b *testing.B) {
	cache := &InMemoryCache{store: make(map[string][]byte)}
	benchmarkCacheStrategy(b, cache, 1024, randomAccessPattern)
}

func BenchmarkFileSystemCacheRandomAccess(b *testing.B) {
	cache := &FileSystemCache{directory: os.TempDir()}
	benchmarkCacheStrategy(b, cache, 1024, randomAccessPattern)
}

func BenchmarkInMemoryCacheSequentialAccess(b *testing.B) {
	cache := &InMemoryCache{store: make(map[string][]byte)}
	benchmarkCacheStrategy(b, cache, 1024, sequentialAccessPattern)
}

func BenchmarkFileSystemCacheSequentialAccess(b *testing.B) {
	cache := &FileSystemCache{directory: os.TempDir()}
	benchmarkCacheStrategy(b, cache, 1024, sequentialAccessPattern)
}

const (
	// Maximum number of goroutines for concurrent testing
	maxConcurrency = 100
)

func benchmarkCacheStrategyConcurrently(b *testing.B, strategy CacheStrategy, fileSize int, accessPattern func(b *testing.B, files []string, strategy CacheStrategy), numGoroutines int) {
	files := setupFiles(b, fileSize)
	var wg sync.WaitGroup
	var cacheHits int64

	// Start the concurrent test
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func() {
			for j := 0; j < b.N/numGoroutines; j++ {
				if err := strategy.Write(make([]byte, 1024), files[j%len(files)]); err != nil {
					b.Fatal(err)
				}
				data, err := strategy.Read(files[j%len(files)])
				if err != nil && err != os.ErrNotExist {
					b.Fatal(err)
				}
				if len(data) > 0 {
					atomic.AddInt64(&cacheHits, 1)
				}
			}
			wg.Done()
		}()
	}
	// Wait for all goroutines to complete
	wg.Wait()

	// Calculate cache hit rate
	cacheHitRate := float64(cacheHits) / float64(b.N) * 100

	// Report the result
	b.ReportMetric(cacheHitRate, "cache_hit_rate")
}

func BenchmarkInMemoryCacheRandomAccessConcurrently(b *testing.B) {
	cache := &InMemoryCache{store: make(map[string][]byte)}
	benchmarkCacheStrategyConcurrently(b, cache, 1024, randomAccessPattern, maxConcurrency)
}

func BenchmarkFileSystemCacheRandomAccessConcurrently(b *testing.B) {
	cache := &FileSystemCache{directory: os.TempDir()}
	benchmarkCacheStrategyConcurrently(b, cache, 1024, randomAccessPattern, maxConcurrency)
}

func BenchmarkInMemoryCacheSequentialAccessConcurrently(b *testing.B) {
	cache := &InMemoryCache{store: make(map[string][]byte)}
	benchmarkCacheStrategyConcurrently(b, cache, 1024, sequentialAccessPattern, maxConcurrency)
}
func BenchmarkFileSystemCacheSequentialAccessConcurrently(b *testing.B) {
	cache := &FileSystemCache{directory: os.TempDir()}
	benchmarkCacheStrategyConcurrently(b, cache, 1024, sequentialAccessPattern, maxConcurrency)
}
