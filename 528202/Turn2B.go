package main
import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"
)

// Define the functions you want to benchmark
func FunctionOne() {
	time.Sleep(50 * time.Millisecond) // Simulate some work
}

func FunctionTwo() {
	time.Sleep(100 * time.Millisecond) // Simulate some work
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

func runGoVersionBenchmark(goVersion string) error {
	// Ensure the specified Go version is installed
	cmd := exec.Command("go", "version")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to get Go version: %v", err)
	}

	if !containsVersion(string(output), goVersion) {
		return fmt.Errorf("Go %s is not installed or not in use: %s", goVersion, string(output))
	}

	// Run benchmarks for the specified Go version
	for _, bm := range benchmarks {
		result := testing.Benchmark(func(b *testing.B) { benchmarkFunc(b, bm) })
		fmt.Printf("Go %s - %s: %s\n", goVersion, bm.Name, result)
	}

	// Run code coverage analysis for the specified Go version
	coverageOutput, err := runCoverageAnalysis(goVersion)
	if err != nil {
		return fmt.Errorf("failed to run coverage analysis for Go %s: %v", goVersion, err)
	}

	// Parse and print the coverage percentage
	coveragePercentage, err := parseCoveragePercentage(coverageOutput)
	if err != nil {
		return fmt.Errorf("failed to parse coverage output: %v", err)
	}
	fmt.Printf("Go %s - Coverage: %.2f%%\n", goVersion, coveragePercentage)

	return nil
}

func runCoverageAnalysis(goVersion string) (string, error) {
	// Set GO version to the specified version for coverage analysis
	os.Setenv("GO", goVersion)

	// Run 'go test -cover' to generate coverage profile
	cmd := exec.Command("go", "test", "-cover")
	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("failed to run coverage analysis: %v", err)
	}

	// Generate coverage profile in temporary file
	tempFile, err := ioutil.TempFile("", "coverage-profile-*.out")
	if err != nil {
		return "", fmt.Errorf("failed to create temporary file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	// Use 'go tool cover' to convert profile to text format
	coverCmd := exec.Command("go", "tool", "cover", "-html", tempFile.Name(), "-o", "coverage.html")
	err = coverCmd.Run()
	if err != nil {
		return "", fmt.Errorf("failed to convert coverage profile: %v", err)
	}

	// Read the coverage output
	coverageOutput, err := ioutil.ReadFile(tempFile.Name())
	if err != nil {
		return "", fmt.Errorf("failed to read coverage output: %v", err)