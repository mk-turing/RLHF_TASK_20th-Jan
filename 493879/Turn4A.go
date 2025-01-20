package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"golang.org/x/exp/mmap"
)

type Result struct {
	FilePath  string
	Positions []int64
	Err       error
}

func findPatternInFileMMAP(ctx context.Context, filePath, pattern string, chunkSize int64) ([]int64, error) {
	if pattern == "" {
		return nil, fmt.Errorf("pattern cannot be empty")
	}

	reader, err := mmap.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer func() {
		if cErr := reader.Close(); cErr != nil {
			log.Printf("error closing file: %v", cErr)
		}
	}()

	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}
	fileSize := fileInfo.Size()

	var positions []int64
	buffer := make([]byte, chunkSize)
	overlap := int64(len(pattern) - 1)

	for offset := int64(0); offset < fileSize; offset += chunkSize - overlap {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
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
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			default:
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

func processFilesConcurrently(files []string, pattern string, chunkSize int64, timeout time.Duration) []Result {
	results := make([]Result, len(files))
	var wg sync.WaitGroup
	var mu sync.Mutex

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	resultChan := make(chan Result, len(files))

	for _, file := range files {
		wg.Add(1)
		// Launch a Goroutine for each file
		go func(filePath string) {
			defer wg.Done()
			fileCtx, fileCancel := context.WithCancel(ctx)
			defer fileCancel()

			positions, err := findPatternInFileMMAP(fileCtx, filePath, pattern, chunkSize)
			resultChan <- Result{
				FilePath:  filePath,
				Positions: positions,
				Err:       err,
			}
		}(file)
	}

	// Wait for all Goroutines to complete and close the result channel
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	i := 0
	for result := range resultChan {
		mu.Lock()
		results[i] = result
		i++
		mu.Unlock()
	}

	return results
}

func main() {
	files := []string{"largefile1.dat", "largefile2.dat", "largefile3.dat"}
	pattern := "keyword"
	chunkSize := int64(1024 * 1024) // 1 MB chunks
	timeout := 10 * time.Second

	results := processFilesConcurrently(files, pattern, chunkSize, timeout)

	for _, result := range results {
		if result.Err != nil {
			if result.Err == context.DeadlineExceeded {
				fmt.Printf("Search in file '%s' canceled due to timeout.\n", result.FilePath)
			} else {
				fmt.Printf("Error processing file '%s': %v\n", result.FilePath, result.Err)
			}
		} else {
			fmt.Printf("Pattern found in file '%s' at positions: %v\n", result.FilePath, result.Positions)
		}
	}
}
