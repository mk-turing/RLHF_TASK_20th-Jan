package main

import (
	"log"
	"time"
)

// RenderComponent simulates rendering a UI component
func RenderComponent() {
	startTime := time.Now()

	// Simulate component rendering with a sleep
	time.Sleep(50 * time.Millisecond) // Simulated render time

	elapsedTime := time.Since(startTime)
	log.Printf("RenderComponent executed in %s", elapsedTime)
}

func main() {
	RenderComponent()
}
