package main

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"golang.org/x/exp/mmap"
)

// findPatternInFile searches for occurrences of a pattern in a file using memory-mapped I/O.
// It processes the file in chunks to handle large files efficiently.
func findPatternInFileMMAP(filePath, pattern string, chunkSize int64) ([]int64, error) {
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
	var buffer = make([]byte, chunkSize)
	var overlap = int64(len(pattern) - 1)

	for offset := int64(0); offset < fileSize; offset += chunkSize - overlap {
		length := chunkSize
		if offset+length > fileSize {
			length = fileSize - offset
		}

		n, err := reader.ReadAt(buffer[:length], offset)
		if err != nil {
			return nil, fmt.Errorf("failed to read file segment: %w", err)
		}
		bufferSegment := buffer[:n]

		// Search for pattern in the current chunk
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

	// Find pattern occurrences using memory-mapped file I/O
	positions, err := findPatternInFileMMAP(filePath, pattern, chunkSize)
	if err != nil {
		log.Fatalf("error finding pattern: %v", err)
	}
	fmt.Printf("Pattern found at positions: %v\n", positions)
}
