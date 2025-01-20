package main

import (
	"fmt"
	"sync"
	"time"
)

// Mock any UI component that needs rendering.
type Component struct {
	ID          int
	NeedsRender bool
}

// Simulate the rendering of a UI component
func (c *Component) Render() {
	time.Sleep(20 * time.Millisecond) // Simulated delay
	fmt.Printf("Rendered component %d\n", c.ID)
}

// Optimize the rendering of multiple components using batching and concurrency.
func BatchRender(components []Component) {
	var wg sync.WaitGroup
	batchSize := 10 // Number of components per batch

	for i := 0; i < len(components); i += batchSize {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			end := i + batchSize
			if end > len(components) {
				end = len(components)
			}
			batch := components[i:end]
			for _, component := range batch {
				if component.NeedsRender {
					component.Render()
				}
			}
		}(i)
	}
	wg.Wait()
}

func main() {
	// Example set of components that need rendering.
	components := make([]Component, 50)
	for i := 0; i < 50; i++ {
		components[i] = Component{ID: i, NeedsRender: (i%2 == 0)} // Example condition
	}

	start := time.Now()
	BatchRender(components)
	elapsed := time.Since(start)

	fmt.Printf("Batch rendering completed in %s\n", elapsed)
}
