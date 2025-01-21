package _28240

import (
	"fmt"
	"github.com/cucumber/godog"
	"os"
	"testing"
)

var firstNumber, secondNumber, result int

func iHaveEnteredIntoTheCalculator(arg1 int) error {
	if firstNumber == 0 {
		firstNumber = arg1
	} else {
		secondNumber = arg1
	}
	return nil
}

func iPressAdd() error {
	result = firstNumber + secondNumber
	return nil
}

func theResultShouldBeOnTheScreen(expectedResult int) error {
	if result != expectedResult {
		return fmt.Errorf("expected %d but got %d", expectedResult, result)
	}
	return nil
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Step(`^I have entered (\d+) into the calculator$`, iHaveEnteredIntoTheCalculator)
	ctx.Step(`^I press add$`, iPressAdd)
	ctx.Step(`^the result should be (\d+) on the screen$`, theResultShouldBeOnTheScreen)
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
