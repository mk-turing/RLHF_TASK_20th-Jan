package main

import (
	"fmt"
	"golang.org/x/exp/mmap"
	"log"
)

// openAndReadSegment opens a file and reads a segment from it using memory mapping.
// start and length define the segment of the file you want to read.
func openAndReadSegment(filePath string, start, length int64) ([]byte, error) {
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

	// Create a buffer to hold the segment data
	buffer := make([]byte, length)

	// Read the segment from the file
	n, err := reader.ReadAt(buffer, start)
	if err != nil {
		return nil, fmt.Errorf("failed to read from file: %w", err)
	}
	if int64(n) < length {
		return nil, fmt.Errorf("failed to read full segment, got %d bytes", n)
	}

	return buffer, nil
}

func main() {
	// Example usage of openAndReadSegment
	filePath := "largefile.dat"
	start := int64(0)   // Start reading from the beginning of the file
	length := int64(64) // Read 64 bytes (for demonstration purposes)

	// Read the specified segment of the file
	data, err := openAndReadSegment(filePath, start, length)
	if err != nil {
		log.Fatalf("error reading segment: %v", err)
	}

	// Print the data read from the file
	fmt.Printf("Segment data: %s\n", data)
}
