package Turn2B

import (
	"sync"
	"testing"
)

func TestConcurrentAdd(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(10)

	for i := 0; i < 5; i++ {
		go ConcurrentAdd(i, 2, &wg)
	}
	for i := 5; i < 10; i++ {
		go ConcurrentAdd(i, -2, &wg)
	}
	wg.Wait()
}

func TestConcurrentMultiply(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(10)

	for i := 0; i < 5; i++ {
		go ConcurrentMultiply(i, 3, &wg)
	}
	for i := 5; i < 10; i++ {
		go ConcurrentMultiply(i, -3, &wg)
	}
	wg.Wait()
}

func TestConcurrentDivide(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(10)

	for i := 0; i < 5; i++ {
		go ConcurrentDivide(i, 2, &wg)
	}
	for i := 5; i < 10; i++ {
		go ConcurrentDivide(i, 0, &wg)
	}
	wg.Wait()
}

func TestSafeDivide(t *testing.T) {
	var wg sync.WaitGroup
	var mutex sync.Mutex
	wg.Add(10)

	for i := 0; i < 5; i++ {
		go SafeDivide(i, 2, &wg, &mutex)
	}
	for i := 5; i < 10; i++ {
		go SafeDivide(i, 0, &wg, &mutex)
	}
	wg.Wait()
}
