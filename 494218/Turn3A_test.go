package _94218

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

type DistributedCacheStrategy interface {
	Write(nodeID string, data []byte, filename string) error
	Read(nodeID string, filename string) ([]byte, error)
	SimulateFailure(nodeID string) error
	RecoverNode(nodeID string) error
}

type Node struct {
	ID    string
	Cache map[string][]byte
}

type SimpleDistributedCache struct {
	nodes        map[string]*Node
	mu           sync.RWMutex
	networkDelay time.Duration
}

func NewSimpleDistributedCache(nodeIDs []string, delay time.Duration) *SimpleDistributedCache {
	nodes := make(map[string]*Node)
	for _, id := range nodeIDs {
		nodes[id] = &Node{ID: id, Cache: make(map[string][]byte)}
	}
	return &SimpleDistributedCache{nodes: nodes, networkDelay: delay}
}

func (s *SimpleDistributedCache) Write(nodeID string, data []byte, filename string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	time.Sleep(s.networkDelay) // Simulating network delay

	if node, exists := s.nodes[nodeID]; exists && node != nil {
		node.Cache[filename] = data
		return nil
	}
	return fmt.Errorf("node %s not found or down", nodeID)
}

func (s *SimpleDistributedCache) Read(nodeID string, filename string) ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	time.Sleep(s.networkDelay) // Simulating network delay

	if node, exists := s.nodes[nodeID]; exists && node != nil {
		data, exists := node.Cache[filename]
		if !exists {
			return nil, fmt.Errorf("file %s not found", filename)
		}
		return data, nil
	}
	return nil, fmt.Errorf("node %s not found or down", nodeID)
}

func (s *SimpleDistributedCache) SimulateFailure(nodeID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.nodes[nodeID]; exists {
		s.nodes[nodeID] = nil
		return nil
	}
	return fmt.Errorf("node %s not found", nodeID)
}

func (s *SimpleDistributedCache) RecoverNode(nodeID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.nodes[nodeID]; exists {
		s.nodes[nodeID] = &Node{ID: nodeID, Cache: make(map[string][]byte)}
		return nil
	}
	return fmt.Errorf("node %s not found", nodeID)
}

func TestDistributedCache(t *testing.T) {
	nodeIDs := []string{"node1", "node2", "node3"}
	cache := NewSimpleDistributedCache(nodeIDs, 100*time.Millisecond)

	// Write data to all nodes
	for _, nodeID := range nodeIDs {
		if err := cache.Write(nodeID, []byte("Hello"), "greeting.txt"); err != nil {
			t.Fatal(err)
		}
	}

	// Test data retrieval
	for _, nodeID := range nodeIDs {
		data, err := cache.Read(nodeID, "greeting.txt")
		if err != nil || string(data) != "Hello" {
			t.Fatalf("failed to read from %s", nodeID)
		}
	}

	// Simulate node failure
	if err := cache.SimulateFailure("node2"); err != nil {
		t.Fatal(err)
	}

	// Test retrieval from failed node
	_, err := cache.Read("node2", "greeting.txt")
	if err == nil {
		t.Fatalf("expected error when reading from down node")
	}

	// Recover node and verify it's empty
	if err := cache.RecoverNode("node2"); err != nil {
		t.Fatal(err)
	}
	data, err := cache.Read("node2", "greeting.txt")
	if err == nil {
		t.Fatalf("expected error after recovery, got data: %s", string(data))
	}
}
