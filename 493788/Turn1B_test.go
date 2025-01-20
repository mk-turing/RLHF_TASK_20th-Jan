package main

import (
	"testing"
	"time"
)

func RenderUIComponent() {
	// Simulate the process of rendering a UI component
	time.Sleep(50 * time.Millisecond) // Simulated delay
}

func BenchmarkRenderUI(b *testing.B) {
	// Start benchmarking the RenderUIComponent
	for i := 0; i < b.N; i++ {
		RenderUIComponent()
	}
}
