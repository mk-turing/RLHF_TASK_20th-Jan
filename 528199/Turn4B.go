package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	_ "github.com/lib/pq"
)

const (
	dsn               = "user=postgres password=example dbname=postgres sslmode=disable"
	batchSize         = 10
	streamSize        = 500
	maxQueueSize      = 10
	totalBatches      = 5
	numStreams        = 3
	maxFailedAttempts = 3
	cooldownDuration  = 5 * time.Second // Circuit breaker cooldown
)

type circuitBreaker struct {
	mutex         sync.Mutex
	failureCount  int
	lastFailure   time.Time
	open          bool
	maxAttempts   int
	cooldown      time.Duration
}

func (cb *circuitBreaker) Allow() bool {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	if cb.open {
		if time.Since(cb.lastFailure) > cb.cooldown {
			cb.failureCount = 0
			cb.open = false
		}
	}
	return !cb.open
}

func (cb *circuitBreaker) Trip() {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	cb.open = true
	cb.failureCount++
	cb.lastFailure = time.Now()
}

// Simulate failure for demonstration purposes
func processDataBatch(data []int, streamId int, batchId int, wg *sync.WaitGroup, results chan<- int, failureProbability float64) {
	defer wg.Done()

	sum := 0
	for _, value := range data {
		sum += value
	}

	// Introduce random failures with the given probability
	if rand.Float64() < failureProbability {
		log.Printf("Stream %d - Batch %d processing failed.", streamId, batchId)
		return
	}

	fmt.Printf("Stream %d - Batch %d processed, sum: %d\n", streamId, batchId, sum)
	results <- sum
}

func handleStream(streamId int, results chan<- int, wg *sync.WaitGroup, failureProbability float64) {
	defer wg.Done()
	cb := circuitBreaker{maxAttempts: maxFailedAttempts, cooldown: cooldownDuration}

	dataStream := make(chan int, maxQueueSize)
	go realTimeDataGenerator(dataStream, streamId)

	var localWg sync.WaitGroup
	batchCounter := 0
	currentBatch := make([]int, 0, batchSize)

	for batchCounter < totalBatches {
		if !cb.Allow() {
			log.Printf("Stream %d: Circuit breaker open. Waiting for cooldown period...", streamId)
			time.Sleep(cb.cooldown)