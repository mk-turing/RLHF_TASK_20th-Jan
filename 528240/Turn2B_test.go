package _28240

import (
	"fmt"
	"github.com/cucumber/godog"
	"os"
	"testing"
	"time"
)

var (
	firstNumber, secondNumber, result int
	asyncResult                       chan int
)

func iHaveEnteredIntoTheCalculator(arg1 int) error {
	if firstNumber == 0 {
		firstNumber = arg1
	} else {
		secondNumber = arg1
	}
	return nil
}

func iPressAddAsynchronously() error {
	asyncResult = make(chan int)
	go func() {
		time.Sleep(2 * time.Second) // Simulate some asynchronous processing
		result = firstNumber + secondNumber
		asyncResult <- result
	}()
	return nil
}

func iShouldReceiveTheResultEventually(expectedResult int) error {
	select {
	case res := <-asyncResult:
		if res != expectedResult {
			return fmt.Errorf("expected %d but got %d", expectedResult, res)
		}
		return nil
	case <-time.After(5 * time.Second): // Timeout after 5 seconds
		return fmt.Errorf("timed out waiting for result")
	}
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Step(`^I have entered (\d+) into the calculator$`, iHaveEnteredIntoTheCalculator)
	ctx.Step(`^I press add asynchronously$`, iPressAddAsynchronously)
	ctx.Step(`^I should receive the result (\d+) on the screen eventually$`, iShouldReceiveTheResultEventually)
}

func TestMain(m *testing.M) {
	opts := godog.Options{
		Format: "pretty",             // Output format
		Paths:  []string{"features"}, // Path to the feature files
	}
	status := godog.TestSuite{
		Name:                "calculator",
		ScenarioInitializer: InitializeScenario,
		Options:             &opts,
	}.Run()

	if st := m.Run(); st > status {
		status = st
	}
	os.Exit(status)
}
