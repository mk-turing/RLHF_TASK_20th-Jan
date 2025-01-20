package main

import (
	"bytes"
	"fmt"
	"golang.org/x/exp/mmap"
	"io"
	"log"
	"os"
)

// findPatternInFile reads the file using standard I/O and searches for a pattern.
func findPatternInFile(filePath, pattern string, chunkSize int64) ([]int64, error) {
	if pattern == "" {
		return nil, fmt.Errorf("pattern cannot be empty")
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var positions []int64
	var buffer = make([]byte, chunkSize)
	var overlap = int64(len(pattern) - 1)
	var offset int64 = 0

	for {
		n, err := file.ReadAt(buffer[:], offset)
		if err != nil && err != io.EOF {
			return nil, fmt.Errorf("failed to read file segment: %w", err)
		}
		if n == 0 {
			break
		}
		bufferSegment := buffer[:n]

		for {
			index := bytes.Index(bufferSegment, []byte(pattern))
			if index == -1 {
				break
			}
			positions = append(positions, offset+int64(index))
			bufferSegment = bufferSegment[index+1:]
			offset += int64(index) + 1
		}
		offset += int64(n) - overlap
	}

	return positions, nil
}

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

// Compare performance between memory-mapped and standard I/O
func main() {
	filePath := "largefile.dat"
	pattern := "keyword"
	chunkSize := int64(1024 * 1024) // 1 MB chunks

	// Memory-mapped approach
	positionsMMAP, err := findPatternInFileMMAP(filePath, pattern, chunkSize)
	if err != nil {
		log.Fatalf("error finding pattern with mmap: %v", err)
	}
	fmt.Printf("MMAP Pattern found at positions: %v\n", positionsMMAP)

	// Standard I/O approach
	positionsIO, err := findPatternInFile(filePath, pattern, chunkSize)
	if err != nil {
		log.Fatalf("error finding pattern with stdio: %v", err)
	}
	fmt.Printf("Standard I/O Pattern found at positions: %v\n", positionsIO)
}
