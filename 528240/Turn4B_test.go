package _28240

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/cucumber/godog"
)

type UserStory struct {
	Feature   string `json:"feature"`
	Scenarios []struct {
		Name  string   `json:"name"`
		Steps []string `json:"steps"`
	} `json:"scenarios"`
}

func generateTests(userStories []UserStory) {
	featureDir := "features"
	for _, story := range userStories {
		featureFilePath := filepath.Join(featureDir, fmt.Sprintf("%s.feature", story.Feature))
		content := fmt.Sprintf("Feature: %s\n", story.Feature)
		for _, scenario := range story.Scenarios {
			content += fmt.Sprintf("  Scenario: %s\n", scenario.Name)
			for _, step := range scenario.Steps {
				content += fmt.Sprintf("    %s\n", step)
			}
		}
		_ = ioutil.WriteFile(featureFilePath, []byte(content), 0644)
	}
}

func syncWithCentralizedSystem() ([]UserStory, error) {
	// In a real scenario, replace this with an API call
	jsonFilePath := "centralized_system.json"
	data, err := ioutil.ReadFile(jsonFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read centralized system data: %v", err)
	}
	var userStories []UserStory
	err = json.Unmarshal(data, &userStories)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal centralized system data: %v", err)
	}
	return userStories, nil
}

func periodicallySync() {
	for {
		userStories, err := syncWithCentralizedSystem()
		if err != nil {
			fmt.Printf("Error syncing: %v\n", err)
		} else {
			generateTests(userStories)
			fmt.Println("Tests synchronized successfully.")
		}
		time.Sleep(10 * time.Second) // Sleep for 10 seconds before next sync
	}
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	// Start the periodic sync goroutine in the background
	go periodicallySync()
}

func TestMain(m *testing.M) {
	opts := godog.Options{
		Format: "pretty",             // Output format
		Paths:  []string{"features"}, // Path to the feature files
	}
	status := godog.TestSuite{
		Name:                "dynamic",
		ScenarioInitializer: InitializeScenario,
		Options:             &opts,
	}.Run()

	if st := m.Run(); st > status {
		status = st
	}
	os.Exit(status)
}
