package main

import (
	"bytes"
	"fmt"
	"log"
	"runtime"
	"strings"
	"sync"
	"time"

	"golang.org/x/exp/mmap"
)

const chunkSize = 1024 * 1024 // 1MB chunk size

// openAndFindPattern opens a file, finds a pattern in it using memory mapping, and performs concurrent processing.
func openAndFindPattern(filePath, pattern string) ([]int, error) {
	if pattern == "" {
		return nil, fmt.Errorf("empty pattern")
	}

	patternBytes := []byte(pattern)

	reader, err := mmap.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer func() {
		if err := reader.Close(); err != nil {
			log.Printf("error closing file: %v", err)
		}
	}()

	fileSize, err := reader.Size()
	if err != nil {
		return nil, fmt.Errorf("failed to get file size: %w", err)
	}

	var wg sync.WaitGroup
	offset := int64(0)
	patternIndices := make([]int, 0, 1024)
	numChunks := (fileSize + chunkSize - 1) / chunkSize
	wg.Add(int(numChunks))

	for i := 0; i < int(numChunks); i++ {
		go func(chunkOffset int64) {
			defer wg.Done()
			chunkLength := int64(chunkSize)
			if chunkOffset+chunkLength > fileSize {
				chunkLength = fileSize - chunkOffset
			}
			chunk, err := reader.Slice(chunkOffset, chunkOffset+chunkLength)
			if err != nil {
				log.Printf("error reading chunk: %v", err)
				return
			}
			for idx := bytes.Index(chunk, patternBytes); idx != -1; idx = bytes.Index(chunk[idx+len(patternBytes):], patternBytes) {
				patternIndices = append(patternIndices, int(chunkOffset+int64(idx)))
			}
		}(offset)
		offset += chunkSize
	}

	wg.Wait()

	return patternIndices, nil
}

func main() {
	filePath := "largefile.dat"
	pattern := "hello"

	startTime := time.Now()
	indices, err := openAndFindPattern(filePath, pattern)
	elapsedTime := time.Since(startTime)
	if err != nil {
		log.Fatalf("error finding pattern: %v", err)
	}

	fmt.Printf("Found %d occurrences of '%s' in %s\n", len(indices), pattern, elapsedTime)

	fmt.Println("Performance Comparison:")
	runtime.GC() // force GC to get stable memory usage
	startTime = time.Now()
	data, err := mmap.Open(filePath)
	if err != nil {
		log.Fatalf("error opening file for non-mapped approach: %v", err)
	}
	defer func() {
		if err := data.Close(); err != nil {
			log.Printf("error closing file for non-mapped approach: %v", err)
		}
	}()

	fileContent := string(data.Data())
	nonMappedIndices := make([]int, 0, 1024)
	for idx := strings.Index(fileContent, pattern); idx != -1; idx = strings.Index(fileContent[idx+len(pattern):], pattern) {
		nonMappedIndices = append(nonMappedIndices, idx)
	}
	elapsedTime = time.Since(startTime)
	fmt.Printf("Found %d occurrences of '%s' in %s (non-mapped approach)\n", len(nonMappedIndices), pattern, elapsedTime)
}
