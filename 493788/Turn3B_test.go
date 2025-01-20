package main

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestStress(t *testing.T) {
	const url = "http://localhost:8080/myapp/component1"

	count := 100 // number of requests
	meanTime := make(chan time.Duration)
	var wg sync.WaitGroup

	wg.Add(count)
	start := time.Now()

	ticker := time.NewTicker(time.Millisecond * 50) // Rapid firing requests
	defer ticker.Stop()

	go func() {
		for _ = range ticker.C {
			go func() {
				defer wg.Done()
				checkTime := pingdom.CheckTime(url)
				meanTime <- checkTime
			}()
		}
	}()

	wg.Wait()
	close(meanTime)

	wg.Add(count)

	for i := 0; i < count; i++ {
		go func() {
			defer wg.Done()
			checkTime := pingdom.CheckTime(url)
			meanTime <- checkTime
		}()
	}

	wg.Wait()
	close(meanTime)

	var totalTime time.Duration
	for result := range meanTime {
		totalTime += result
	}

	//Calculate average time
	averageTime := totalTime / time.Duration(count)
	fps := int(1 / float64(averageTime.Seconds()))

	fmt.Printf("Average Response time: %s, Frame rate: %dFPS \n", averageTime, fps)
}
