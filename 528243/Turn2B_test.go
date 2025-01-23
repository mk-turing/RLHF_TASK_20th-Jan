package _28243

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

type TestResult struct {
	TestName        string `json:"test_name"`
	ExecutionTime   time.Duration `json:"execution_time"`
	MemoryAllocated uint64 `json:"memory_allocated"`
}

func monitorTest(testName string, testFunc func(*testing.T), t *testing.T) TestResult {
	start := time.Now()
	var memStatsStart, memStatsEnd runtime.MemStats
	runtime.ReadMemStats(&memStatsStart)
	testFunc(t)
	runtime.ReadMemStats(&memStatsEnd)
	elapsed := time.Since(start)
	return TestResult{
		TestName:        testName,
		ExecutionTime:   elapsed,
		MemoryAllocated: memStatsEnd.Alloc - memStatsStart.Alloc,
	}
}

func runTestsAndRecordResults() []TestResult {
	results := make([]TestResult, 0)
	results = append(results, monitorTest("TestFunctionToTest1", testFunctionToTest1, testing.New()))
	results = append(results, monitorTest("TestFunctionToTest2", testFunctionToTest2, testing.New()))
	return results
}

func getBaselinePath() string {
	dir, _ := os.Getwd()
	return filepath.Join(dir, "test_baseline.json")
}

func loadBaselineData() ([]TestResult, error) {
	file, err := os.Open(getBaselinePath())
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var baseline []TestResult
	if err := json.NewDecoder(file).Decode(&baseline); err != nil {
		return nil, err
	}
	return baseline, nil
}

func saveBaselineData(results []TestResult) error {
	file, err := os.Create(getBaselinePath())
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", " ")
	if err := encoder.Encode(results); err != nil {
		return err
	}
	return nil
}

func compareResultsWithBaseline(results []TestResult) error {
	baseline, err := loadBaselineData()
	if err != nil {
		return fmt.Errorf("failed to load baseline data: %w", err)
	}

	for _, result := range results {
		var baselineResult *TestResult
		for _, br := range baseline {
			if br.TestName == result.TestName {
				baselineResult = &br
				break
			}
		}
		if baselineResult == nil {
			fmt.Printf("Baseline not found for test: %s\n", result.TestName)
			continue
		}

		// Adjust comparison logic as needed
		timeDiff := result.ExecutionTime.Sub(baselineResult.ExecutionTime)
		memoryDiff := float64(result.MemoryAllocated) - float64(baselineResult.MemoryAllocated)
		timeThreshold := time.Duration(float64(baselineResult.ExecutionTime) * 0.1) // 10% threshold
		memoryThreshold := float64(baselineResult.MemoryAllocated) * 0.1 // 10% threshold

		if timeDiff > timeThreshold {
			fmt.Printf("Performance Regression in Test: %s\n", result.TestName)
			fmt.Printf("  Baseline Time: %s\n", baselineResult.ExecutionTime)
			fmt.Printf("  Current Time: %s\n", result.ExecutionTime)
			fmt.Printf("  Time Difference: %s\n", timeDiff)
		}