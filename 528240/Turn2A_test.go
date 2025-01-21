package _28240

import (
	"fmt"
	"github.com/cucumber/godog"
	"os"
	"sync"
	"testing"
	"time"
)

var taskQueue chan string
var taskResults map[string]bool
var mu sync.Mutex

func iHaveATaskToProcess() error {
	taskQueue = make(chan string, 1) // Channel for task queue
	taskResults = make(map[string]bool)
	taskQueue <- "task1" // Add a task to the queue
	return nil
}

func iStartTheTaskProcessor() error {
	go func() {
		for task := range taskQueue {
			time.Sleep(1 * time.Second) // Simulate task processing time
			mu.Lock()
			taskResults[task] = true // Mark task as processed
			mu.Unlock()
		}
	}()
	return nil
}

func theTaskShouldBeProcessedSuccessfully() error {
	time.Sleep(2 * time.Second) // Wait for task to be processed
	mu.Lock()
	defer mu.Unlock()
	if !taskResults["task1"] {
		return fmt.Errorf("task1 was not processed successfully")
	}
	return nil
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Step(`^I have a task to process$`, iHaveATaskToProcess)
	ctx.Step(`^I start the task processor$`, iStartTheTaskProcessor)
	ctx.Step(`^the task should be processed successfully$`, theTaskShouldBeProcessedSuccessfully)
}

func TestMain(m *testing.M) {
	opts := godog.Options{
		Format: "pretty",             // Output format
		Paths:  []string{"features"}, // Path to the feature files
	}
	status := godog.TestSuite{
		Name:                "async",
		ScenarioInitializer: InitializeScenario,
		Options:             &opts,
	}.Run()

	if st := m.Run(); st > status {
		status = st
	}
	os.Exit(status)
}
