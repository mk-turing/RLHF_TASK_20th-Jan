package _28240

import (
	"testing"
)

// UserStory represents a user story in the BDD framework
type UserStory struct {
	Title       string
	Scenarios   []Scenario
	TestResults []TestResult
}

// Scenario represents a scenario in a user story
type Scenario struct {
	Title      string
	Given      []string
	When       string
	Then       []string
	TestFunc   func(*testing.T)
	TestResult TestResult
}

// TestResult represents the outcome of a test
type TestResult struct {
	Passed bool
	Error  error
}

// RunTests executes the tests for a given user story
func (us *UserStory) RunTests(t *testing.T) {
	for _, sc := range us.Scenarios {
		t.Run(sc.Title, sc.TestFunc)
		us.TestResults = append(us.TestResults, TestResult{Passed: t.Failed()})
	}
}

// Given adds a "Given" step to a scenario
func Given(step string) func(sc *Scenario) {
	return func(sc *Scenario) {
		sc.Given = append(sc.Given, step)
	}
}

// When adds a "When" step to a scenario
func When(step string) func(sc *Scenario) {
	return func(sc *Scenario) {
		sc.When = step
	}
}

// Then adds a "Then" step to a scenario
func Then(step string) func(sc *Scenario) {
	return func(sc *Scenario) {
		sc.Then = append(sc.Then, step)
	}
}

// And is a utility function to chain Given, When, and Then steps
func And(steps ...func(sc *Scenario)) func(sc *Scenario) {
	return func(sc *Scenario) {
		for _, step := range steps {
			step(sc)
		}
	}
}

// UserStoryWithScenario defines a user story with one or more scenarios
func UserStoryWithScenario(title string, scenarios ...Scenario) UserStory {
	return UserStory{Title: title, Scenarios: scenarios}
}

func TestUserStory(t *testing.T) {
	additionUS := UserStoryWithScenario(
		"Adding two numbers",
		Scenario{
			Title: "Adding two positive numbers",
			Given: []string{"a = 5", "b = 3"},
			When:  "I add a and b",
			Then:  []string{"The result should be 8"},
			TestFunc: func(t *testing.T) {
				a, b := 5, 3
				result := a + b
				if result != 8 {
					t.Errorf("Expected 8, got %d", result)
				}
			},
		},
		Scenario{
			Title: "Adding two negative numbers",
			Given: []string{"a = -5", "b = -3"},
			When:  "I add a and b",
			Then:  []string{"The result should be -8"},
			TestFunc: func(t *testing.T) {
				a, b := -5, -3
				result := a + b
				if result != -8 {
					t.Errorf("Expected -8, got %d", result)
				}
			},
		},
	)

	t.Run("Addition User Story", func(t *testing.T) {
		additionUS.RunTests(t)
		for i, sc := range additionUS.Scenarios {