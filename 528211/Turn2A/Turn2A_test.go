package Turn2A

import (
	"sync"
	"testing"
)

func TestConcurrentAccumulator(t *testing.T) {
	acc := &ConcurrentAccumulator{}

	// Edge case: Adding zero should not affect the total.
	acc.Add(0)
	if acc.Total() != 0 {
		t.Errorf("Expected total 0, got %d", acc.Total())
	}

	// Concurrency scenario: Simulate concurrent additions.
	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			acc.Add(1)
		}()
	}
	wg.Wait()

	// Validate final total matches expected value.
	if acc.Total() != 1000 {
		t.Errorf("Expected total 1000, got %d", acc.Total())
	}
}
