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
	Error     error
}

func findPatternInFileMMAP(ctx context.Context, filePath, pattern string, chunkSize int64, result chan<- Result) {
	var res Result
	defer func() {
		result <- res
	}()

	if pattern == "" {
		res.Error = fmt.Errorf("pattern cannot be empty")
		return
	}

	reader, err := mmap.Open(filePath)
	if err != nil {
		res.Error = fmt.Errorf("failed to open file %s: %w", filePath, err)
		return
	}
	defer func() {
		if err := reader.Close(); err != nil {
			log.Printf("error closing file: %v", err)
		}
	}()

	fileInfo, err := os.Stat(filePath)
	if err != nil {
		res.Error = fmt.Errorf("failed to get file info: %w", err)
		return
	}
	fileSize := fileInfo.Size()

	var positions []int64
	var buffer = make([]byte, chunkSize)
	var overlap = int64(len(pattern) - 1)

	for offset := int64(0); offset < fileSize; offset += chunkSize - overlap {
		select {
		case <-ctx.Done():
			res.Error = ctx.Err()
			return
		default:
		}

		length := chunkSize
		if offset+length > fileSize {
			length = fileSize - offset
		}

		n, err := reader.ReadAt(buffer[:length], offset)
		if err != nil {
			res.Error = fmt.Errorf("failed to read file segment: %w", err)
			return
		}
		bufferSegment := buffer[:n]

		for {
			select {
			case <-ctx.Done():
				res.Error = ctx.Err()
				return
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

	res.FilePath = filePath
	res.Positions = positions
}

func main() {
	filePaths := []string{"largefile1.dat", "largefile2.dat", "largefile3.dat"}
	pattern := "keyword"
	chunkSize := int64(1024 * 1024) // 1 MB chunks
	numWorkers := len(filePaths)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	results := make(chan Result, numWorkers)
	var wg sync.WaitGroup

	wg.Add(numWorkers)
	for _, filePath := range filePaths {
		go func(filePath string) {
			defer wg.Done()
			findPatternInFileMMAP(ctx, filePath, pattern, chunkSize, results)
		}(filePath)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for res := range results {
		if res.Error != nil {
			if res.Error == context.DeadlineExceeded {
				fmt.Printf("Search in %s canceled due to timeout.\n", res.FilePath)
			} else {
				log.Fatalf("error finding pattern in %s: %v", res.FilePath, res.Error)
			}
		} else {
			fmt.Printf("Pattern found in %s at positions: %v\n", res.FilePath, res.Positions)
		}
	}
}
