package _28243

import (
	"encoding/json"
	"errors"
	"os"
	"runtime"
	"testing"
	"time"
)

type TestResult struct {
	TestName        string        `json:"test_name"`
	ExecutionTime   time.Duration `json:"execution_time"`
	MemoryAllocated uint64        `json:"memory_allocated"`
}

func FunctionToTest(input int) int {
	// Example function logic
	return input * 2
}

func readBaseline() ([]TestResult, error) {
	file, err := os.ReadFile("baseline_metrics.json")
	if err != nil {
		return nil, err
	}
	var baseline []TestResult
	err = json.Unmarshal(file, &baseline)
	if err != nil {
		return nil, errors.New("failed to unmarshal baseline metrics")
	}
	return baseline, nil
}

func updateBaseline(results []TestResult) error {
	file, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile("baseline_metrics.json", file, 0644)
}

const executionTimePercentThreshold = 50 // Allow a 50% difference in execution time
const memoryPercentThreshold = 50        // Allow a 50% difference in memory usage

func compareMetrics(current, baseline TestResult) bool {
	// Calculate the percentage deviation
	exTimePercentDeviation := float64(current.ExecutionTime) / float64(baseline.ExecutionTime) * 100
	memPercentDeviation := float64(current.MemoryAllocated) / float64(baseline.MemoryAllocated) * 100

	// Allow deviations within a specified percentage
	return exTimePercentDeviation > (100+executionTimePercentThreshold) || memPercentDeviation > (100+memoryPercentThreshold)
}

func saveReport(results []TestResult) error {
	file, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile("test_report.json", file, 0644)
}

func monitorTest(testName string, testFunc func(t *testing.T), t *testing.T) TestResult {
	startTime := time.Now()
	var memStatsStart, memStatsEnd runtime.MemStats

	// Capture memory usage before running the test
	runtime.ReadMemStats(&memStatsStart)

	// Run the actual test function
	testFunc(t)

	// Capture memory usage after running the test
	runtime.ReadMemStats(&memStatsEnd)

	// Calculate memory usage and execution time
	executionTime := time.Since(startTime)
	memoryAllocated := memStatsEnd.Alloc - memStatsStart.Alloc

	return TestResult{
		TestName:        testName,
		ExecutionTime:   executionTime,
		MemoryAllocated: memoryAllocated,
	}
}

func TestFunctionToTest(t *testing.T) {
	// Read the baseline metrics
	baseline, err := readBaseline()
	if err != nil {
		t.Fatalf("Error reading baseline metrics: %v", err)
	}

	// Monitor the current test result
	currentResults := []TestResult{
		monitorTest("TestFunctionToTest", func(t *testing.T) {
			result := FunctionToTest(2)
			if result != 4 {
				t.Errorf("Expected 4, but got %d", result)
			}
		}, t),
	}

	// Validate the number of baseline tests matches current tests
	if len(baseline) != len(currentResults) {
		t.Fatalf("Mismatch between baseline (%d) and current results (%d) length", len(baseline), len(currentResults))
	}

	// Compare each result and detect performance regression
	for i, result := range currentResults {
		if compareMetrics(result, baseline[i]) {
			t.Errorf("Performance regression detected in %s: ExecutionTime: %v (baseline: %v), MemoryAllocated: %d (baseline: %d)",
				result.TestName, result.ExecutionTime, baseline[i].ExecutionTime, result.MemoryAllocated, baseline[i].MemoryAllocated)
		}
	}

	// Save the current results for reporting
	err = saveReport(currentResults)
	if err != nil {
		t.Errorf("Error saving test report: %v", err)
	}
}
