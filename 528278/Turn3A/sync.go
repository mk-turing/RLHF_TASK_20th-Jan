package main

import (
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"golang.org/x/net/context"
)

func acquireLock(ctx context.Context, rdb *redis.Client, lockKey string, ttl time.Duration) bool {
	ok, err := rdb.SetNX(ctx, lockKey, "locked", ttl).Result()
	if err != nil {
		fmt.Println("Error acquiring lock:", err)
		return false
	}
	return ok
}

func releaseLock(ctx context.Context, rdb *redis.Client, lockKey string) {
	_, err := rdb.Del(ctx, lockKey).Result()
	if err != nil {
		fmt.Println("Error releasing lock:", err)
	}
}

func synchronizedRangeLoop(ctx context.Context, rdb *redis.Client, nodeID string) {
	lockKey := "range_loop_lock"
	if acquireLock(ctx, rdb, lockKey, 5*time.Second) {
		defer releaseLock(ctx, rdb, lockKey)
		fmt.Printf("Node %s executing range loop\n", nodeID)
		for i := 0; i < 100; i++ {
			time.Sleep(10 * time.Millisecond) // Simulate work
		}
	} else {
		fmt.Printf("Node %s could not acquire lock\n", nodeID)
	}
}

func main() {
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	nodeID := "node-1"
	for {
		synchronizedRangeLoop(ctx, rdb, nodeID)
		time.Sleep(1 * time.Second)
	}
}
