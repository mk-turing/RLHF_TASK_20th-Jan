package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"golang.org/x/exp/mmap"
)

// findPatternInFileMMAP searches for occurrences of a pattern in a file using memory-mapped I/O.
// It processes the file in chunks to handle large files efficiently and supports context cancellation.
func findPatternInFileMMAP(ctx context.Context, filePath, pattern string, chunkSize int64) ([]int64, error) {
	if pattern == "" {
		return nil, fmt.Errorf("pattern cannot be empty")
	}

	// Open the file using mmap.ReaderAt
	reader, err := mmap.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer func() {
		if err := reader.Close(); err != nil {
			log.Printf("error closing file: %v", err)
		}
	}()

	// Get the size of the file
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}
	fileSize := fileInfo.Size()

	var positions []int64
	buffer := make([]byte, chunkSize)
	overlap := int64(len(pattern) - 1)

	for offset := int64(0); offset < fileSize; offset += chunkSize - overlap {
		// Check for context cancellation
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			// Continue processing
		}

		length := chunkSize
		if offset+length > fileSize {
			length = fileSize - offset
		}

		n, err := reader.ReadAt(buffer[:length], offset)
		if err != nil {
			return nil, fmt.Errorf("failed to read file segment: %w", err)
		}
		bufferSegment := buffer[:n]

		// Search for the pattern in the current chunk
		for {
			index := bytes.Index(bufferSegment, []byte(pattern))
			if index == -1 {
				break
			}
			positions = append(positions, offset+int64(index))
			bufferSegment = bufferSegment[index+1:]
			offset += int64(index) + 1
		}
	}

	return positions, nil
}

func main() {
	filePath := "largefile.dat"
	pattern := "keyword"
	chunkSize := int64(1024 * 1024) // 1 MB chunks

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Use a goroutine to simulate an external cancellation trigger
	go func() {
		time.Sleep(2 * time.Second) // Simulate some processing time
		cancel()                    // Trigger cancellation
	}()

	// Find pattern occurrences using memory-mapped file I/O with context for cancellation
	positions, err := findPatternInFileMMAP(ctx, filePath, pattern, chunkSize)
	if err != nil {
		log.Printf("operation canceled or error occurred: %v", err)
		return
	}
	fmt.Printf("Pattern found at positions: %v\n", positions)
}
