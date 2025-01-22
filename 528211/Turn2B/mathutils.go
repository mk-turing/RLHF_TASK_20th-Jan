package Turn2B

import (
	"sync"
	"time"
)

// ConcurrentAdd adds two integers concurrently.
func ConcurrentAdd(a, b int, wg *sync.WaitGroup) {
	defer wg.Done()
	time.Sleep(time.Duration(a) * time.Millisecond) // Introduce delay to simulate concurrency
	result := a + b
	_ = result // Avoid compiler optimization
}

// ConcurrentMultiply multiplies two integers concurrently.
func ConcurrentMultiply(a, b int, wg *sync.WaitGroup) {
	defer wg.Done()
	result := a * b
	_ = result // Avoid compiler optimization
}

// ConcurrentDivide divides two integers concurrently. Returns 0 if dividing by zero.
func ConcurrentDivide(a, b int, wg *sync.WaitGroup) {
	defer wg.Done()
	if b == 0 {
		return
	}
	result := a / b
	_ = result // Avoid compiler optimization
}

// SafeDivide divides two integers concurrently using a mutex to prevent race conditions.
func SafeDivide(a, b int, wg *sync.WaitGroup, mutex *sync.Mutex) {
	defer wg.Done()
	mutex.Lock()
	defer mutex.Unlock()
	if b == 0 {
		return
	}
	result := a / b
	_ = result // Avoid compiler optimization
}
