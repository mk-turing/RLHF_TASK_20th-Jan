package Turn2A

import "sync"

// ConcurrentAccumulator accumulates values concurrently.
type ConcurrentAccumulator struct {
	total int
	mu    sync.Mutex
}

// Add adds a value to the total.
func (acc *ConcurrentAccumulator) Add(value int) {
	acc.mu.Lock()
	defer acc.mu.Unlock()
	acc.total += value
}

// Total returns the accumulated total.
func (acc *ConcurrentAccumulator) Total() int {
	acc.mu.Lock()
	defer acc.mu.Unlock()
	return acc.total
}
