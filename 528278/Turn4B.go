package main

import (
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"runtime/pprof"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

// Config holds the configuration parameters for the benchmark.
type Config struct {
	NumNodes        int   // Number of nodes in the distributed system.
	NodeID          string // Unique ID for each node.
	WorkloadSize   int   // Size of the data processed by each node.
	LoopComplexity  int   // Complexity factor for the range loop.
	SyncInterval   time.Duration // Interval for synchronization.
	MaxSimultaneous int   // Maximum number of simultaneous range loops.
	ServerURL       string // URL of the central aggregation server.
	RedisAddr      string // Address of the Redis server for synchronization.
	ProfilesDir     string // Directory to store performance profiles.
}

// startCPUProfile begins CPU profiling and returns a function to stop it.
func startCPUProfile(filename string) func() {
	f, err := os.Create(filename)
	if err != nil {
		fmt.Println("Could not create CPU profile:", err)
		return func() {}
	}
	if err := pprof.StartCPUProfile(f); err != nil {
		fmt.Println("Could not start CPU profile:", err)
		return func() {}
	}
	return func() {
		pprof.StopCPUProfile()
		f.Close()
	}
}

// simulateWorkload simulates the computational work performed within a range loop.
func simulateWorkload(c *Config) {
	time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
}

// runSingleNodeBenchmark performs the benchmark on a single node.
func runSingleNodeBenchmark(c *Config) {
	stopCPUProfile := startCPUProfile(fmt.Sprintf("%s/single_node_cpu.prof", c.ProfilesDir))
	defer stopCPUProfile()

	for i := 0; i < c.WorkloadSize; i++ {
		simulateWorkload(c)
	}
}

// runDistributedBenchmark performs the benchmark in a distributed system.
func runDistributedBenchmark(c *Config) {
	if c.NumNodes <= 1 {
		fmt.Println("Distributed benchmark requires more than one node.")
		return
	}

	var wg sync.WaitGroup
	rdb := redis.NewClient(&redis.Options{Addr: c.RedisAddr})

	// Start benchmarking on each node
	for i := 0; i < c.NumNodes; i++ {
		nodeID := fmt.Sprintf("node-%d", i+1)
		wg.Add(1)
		go func(nodeID string) {
			defer wg.Done()
			runDistributedNodeBenchmark(c, nodeID, rdb)
		}(nodeID)
	}

	wg.Wait()
}

// runDistributedNodeBenchmark performs the benchmark on a single node in a distributed system.
func runDistributedNodeBenchmark(c *Config, nodeID string, rdb *redis.Client) {
	stopCPUProfile := startCPUProfile(fmt.Sprintf("%s/distributed_node_%s_cpu.prof", c.ProfilesDir, nodeID))
	defer stopCPUProfile()

	lockKey := "range_loop_lock"
	for {
		if acquireLock(context.Background(), rdb, lockKey, 5*time.Second) {
			defer releaseLock(context.Background(), rdb, lockKey)

			// Simulate workload within the critical section
			for i := 0; i < c.WorkloadSize; i++ {
				simulateWorkload(c)