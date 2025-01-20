
package _94218_test

import (
	"bytes"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/benchmark/metrics"
	"github.com/patrickmn/go-cache"
	"go.etcd.io/bbolt"
)

// Define the different caching strategies
type CachingStrategy interface {
	ReadFile(string) ([]byte, error)
	WriteFile(string, []byte) error
}

// NoCachingStrategy: Reads/writes directly to the disk
type NoCachingStrategy struct{}

func (NoCachingStrategy) ReadFile(filename string) ([]byte, error) {
	return ioutil.ReadFile(filename)
}

func (NoCachingStrategy) WriteFile(filename string, data []byte) error {
	return ioutil.WriteFile(filename, data, 0644)
}

// InMemoryCache: Uses a simple map to cache file data
type InMemoryCache struct {
	cache map[string][]byte
}

func NewInMemoryCache() *InMemoryCache {
	return &InMemoryCache{cache: make(map[string][]byte)}
}

func (c *InMemoryCache) ReadFile(filename string) ([]byte, error) {
	data, ok := c.cache[filename]
	if !ok {
		data, err := ioutil.ReadFile(filename)
		if err != nil {
			return nil, err
		}
		c.cache[filename] = data
	}
	return data, nil
}

func (c *InMemoryCache) WriteFile(filename string, data []byte) error {
	c.cache[filename] = data
	return ioutil.WriteFile(filename, data, 0644)
}

// LRUCache: Uses a least-recently-used cache implementation
type LRUCache struct {
	cache *cache.Cache
}

func NewLRUCache(size int) *LRUCache {
	return &LRUCache{cache: cache.New(size, 0)}
}

func (c *LRUCache) ReadFile(filename string) ([]byte, error) {
	data, found := c.cache.Get(filename)
	if !found {
		data, err := ioutil.ReadFile(filename)
		if err != nil {
			return nil, err
		}
		c.cache.Set(filename, data, cache.DefaultExpiration)
	}
	return data.([]byte), nil
}

func (c *LRUCache) WriteFile(filename string, data []byte) error {
	c.cache.Set(filename, data, cache.DefaultExpiration)
	return ioutil.WriteFile(filename, data, 0644)
}

// DiskCache: Uses a disk-based cache (BoltDB)
type DiskCache struct {
	db *bolt.DB
}

func NewDiskCache(filename string) (*DiskCache, error) {
	db, err := bolt.Open(filename, 0644, nil)
	if err != nil {
		return nil, err
	}
	return &DiskCache{db: db}, nil
}

func (c *DiskCache) ReadFile(filename string) ([]byte, error) {
	var data []byte
	err := c.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("files"))
		data = b.Get([]byte(filename))
		return nil