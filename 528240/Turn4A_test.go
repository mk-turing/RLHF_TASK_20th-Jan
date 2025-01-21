package _28240

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type UserStory struct {
	ID        string   `json:"id"`
	Title     string   `json:"title"`
	Feature   string   `json:"feature"`
	Scenarios []string `json:"scenarios"`
}

func fetchUserStories(apiURL string) ([]UserStory, error) {
	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var stories []UserStory
	if err = json.Unmarshal(body, &stories); err != nil {
		return nil, err
	}

	return stories, nil
}

func generateFeatureFiles(stories []UserStory) error {
	for _, story := range stories {
		fileName := fmt.Sprintf("%s.feature", story.ID)
		file, err := os.Create(fileName)
		if err != nil {
			return err
		}
		defer file.Close()

		featureContent := fmt.Sprintf("Feature: %s\n\n%s\n", story.Title, story.Feature)
		for _, scenario := range story.Scenarios {
			featureContent += fmt.Sprintf("\nScenario: %s\n%s\n", scenario, "  # Scenario steps here")
		}

		if _, err := file.WriteString(featureContent); err != nil {
			return err
		}
	}
	return nil
}

func startPeriodicSync(apiURL string, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		stories, err := fetchUserStories(apiURL)
		if err != nil {
			log.Printf("Error fetching user stories: %v", err)
			continue
		}

		if err := generateFeatureFiles(stories); err != nil {
			log.Printf("Error generating feature files: %v", err)
		} else {
			log.Println("Feature files updated successfully.")
		}
	}
}
