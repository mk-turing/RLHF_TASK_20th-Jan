package _94218

import "testing"

// ... (Same code as before)

func BenchmarkDistributedCacheRandomAccess(b *testing.B) {
	// Create two in-memory cache nodes for the distributed system
	node1 := &InMemoryCache{store: make(map[string][]byte)}
	node2 := &InMemoryCache{store: make(map[string][]byte)}
	distributedCache := NewDistributedCache(map[string]CacheStrategy{"node1": node1, "node2": node2})

	benchmarkCacheStrategy(b, distributedCache, 1024, randomAccessPattern)
}

func BenchmarkDistributedCacheSequentialAccess(b *testing.B) {
	// Create two in-memory cache nodes for the distributed system
	node1 := &InMemoryCache{store: make(map[string][]byte)}
	node2 := &InMemoryCache{store: make(map[string][]byte)}
	distributedCache := NewDistributedCache(map[string]CacheStrategy{"node1": node1, "node2": node2})
	benchmarkCacheStrategy(b, distributedCache, 1024, sequentialAccessPattern)
}

func TestDistributedCacheNodeFailure(t *testing.T) {
	// Create two in-memory cache nodes for the distributed system
	node1 := &InMemoryCache{store: make(map[string][]byte)}
	node2 := &InMemoryCache{store: make(map[string][]byte)}
	distributedCache := NewDistributedCache(map[string]CacheStrategy{"node1": node1, "node2": node2})

	// Simulate node failure by removing it from the cache
	distributedCache.mu.Lock()
	delete(distributedCache.nodes, "node2")
	distributedCache.mu.Unlock()