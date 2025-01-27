package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"testing"
	"time"
)

// Predefined benchmark functions
func FunctionOne() {
	time.Sleep(50 * time.Millisecond) // Simulate work
}

func FunctionTwo() {
	time.Sleep(100 * time.Millisecond) // Simulate work
}

type BenchmarkFunc struct {
	Name string
	Fn   func()
}

var benchmarks = []BenchmarkFunc{
	{"FunctionOne", FunctionOne},
	{"FunctionTwo", FunctionTwo},
}

func benchmarkFunc(b *testing.B, bf BenchmarkFunc) {
	for n := 0; n < b.N; n++ {
		bf.Fn()
	}
}

func captureCoverageOutput() (string, error) {
	// Run tests with coverage to capture the coverage output
	cmd := exec.Command("go", "test", "-coverprofile=coverage.out", "./...")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("failed to run coverage analysis: %v", err)
	}

	// Read the coverage file to capture the details
	data, err := os.ReadFile("coverage.out")
	if err != nil {
		return "", fmt.Errorf("failed to read coverage file: %v", err)
	}

	return string(data), nil
}

func parseCoverage(coverageOutput string) (float64, error) {
	// Parse the coverage file for the overall coverage percentage
	lines := strings.Split(coverageOutput, "\n")

	var coveredLines, totalLines int

	for _, line := range lines {
		if strings.HasPrefix(line, "main/") { // File lines start with "main/"
			parts := strings.Fields(line)
			if len(parts) > 2 {
				// The first number indicates the number of times a line was covered
				covered, err := strconv.Atoi(parts[1])
				if err != nil {
					return 0, fmt.Errorf("failed to parse covered count: %v", err)
				}

				// The second number indicates the number of times a line was not covered
				notCovered, err := strconv.Atoi(parts[2])
				if err != nil {
					return 0, fmt.Errorf("failed to parse not-covered count: %v", err)
				}

				coveredLines += covered
				totalLines += covered + notCovered
			}
		}
	}

	if totalLines == 0 {
		return 0, fmt.Errorf("no coverage data available")
	}

	// Calculate coverage percentage
	coverage := (float64(coveredLines) / float64(totalLines)) * 100
	return coverage, nil
}

func runBenchmarksWithCoverage() error {
	// Run benchmarks for each function
	for _, bm := range benchmarks {
		result := testing.Benchmark(func(b *testing.B) { benchmarkFunc(b, bm) })
		fmt.Printf("Benchmark - %s: %s\n", bm.Name, result)

		// Capture and parse coverage
		coverageOutput, err := captureCoverageOutput()
		if err != nil {
			return fmt.Errorf("failed to capture coverage output: %v", err)
		}

		coverage, err := parseCoverage(coverageOutput)
		if err != nil {
			return fmt.Errorf("failed to parse coverage output: %v", err)
		}

		// Print out combined benchmarking and coverage results
		fmt.Printf("Function: %s, Benchmark: %s, Coverage: %.2f%%\n", bm.Name, result, coverage)
	}
	return nil
}

func main() {
	err := runBenchmarksWithCoverage()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
