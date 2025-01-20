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

func findPatternInFileMMAP(ctx context.Context, filePath, pattern string, chunkSize int64) ([]int64, error) {
	if pattern == "" {
		return nil, fmt.Errorf("pattern cannot be empty")
	}

	reader, err := mmap.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer func() {
		if err := reader.Close(); err != nil {
			log.Printf("error closing file: %v", err)
		}
	}()

	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}
	fileSize := fileInfo.Size()

	var positions []int64
	var buffer = make([]byte, chunkSize)
	var overlap = int64(len(pattern) - 1)

	for offset := int64(0); offset < fileSize; offset += chunkSize - overlap {
		if err := ctx.Err(); err != nil {
			return nil, err // Cancel the search if context is canceled
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

		for {
			if err := ctx.Err(); err != nil {
				return nil, err // Cancel the search if context is canceled
			}

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

	// Create a context with a timeout of 5 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	positions, err := findPatternInFileMMAP(ctx, filePath, pattern, chunkSize)
	if err != nil {
		if err == context.DeadlineExceeded {
			fmt.Println("Search canceled due to timeout.")
		} else {
			log.Fatalf("error finding pattern: %v", err)
		}
	} else {
		fmt.Printf("Pattern found at positions: %v\n", positions)
	}
}
