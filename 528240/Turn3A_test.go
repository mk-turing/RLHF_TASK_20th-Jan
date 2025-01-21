package _28240

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/cucumber/godog"
	"os"
	"testing"
)

var server *httptest.Server
var apiAvailable bool
var apiDelay time.Duration

func theAPIIsAvailable() error {
	apiAvailable = true
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if apiDelay > 0 {
			time.Sleep(apiDelay)
		}
		if apiAvailable {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"status":"success"}`))
		} else {
			w.WriteHeader(http.StatusServiceUnavailable)
		}
	}))
	return nil
}

func theAPIIsDown() error {
	apiAvailable = false
	return nil
}

func theAPIHasADelay() error {
	apiDelay = 2 * time.Second
	apiAvailable = true
	return nil
}

func iMakeARequestToTheExternalAPI() error {
	resp, err := http.Get(server.URL)
	if err != nil {
		return fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()
	if apiAvailable && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("expected status 200 but got %d", resp.StatusCode)
	}
	if !apiAvailable && resp.StatusCode != http.StatusServiceUnavailable {
		return fmt.Errorf("expected status 503 but got %d", resp.StatusCode)
	}
	return nil
}

func iShouldReceiveASuccessfulResponse() error {
	// Already checked in the previous step, so we just acknowledge here
	return nil
}

func iShouldReceiveAnErrorResponse() error {
	// Already checked in the previous step, so we just acknowledge here
	return nil
}

func iShouldReceiveAResponseAfterADelay() error {
	start := time.Now()
	resp, err := http.Get(server.URL)
	elapsed := time.Since(start)

	if err != nil {
		return fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if elapsed < 2*time.Second {
		return fmt.Errorf("expected delay of 2 seconds but got %v", elapsed)
	}

	return nil
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Step(`^the API is available$`, theAPIIsAvailable)
	ctx.Step(`^the API is down$`, theAPIIsDown)
	ctx.Step(`^the API has a delay$`, theAPIHasADelay)
	ctx.Step(`^I make a request to the external API$`, iMakeARequestToTheExternalAPI)
	ctx.Step(`^I should receive a successful response$`, iShouldReceiveASuccessfulResponse)
	ctx.Step(`^I should receive an error response$`, iShouldReceiveAnErrorResponse)
	ctx.Step(`^I should receive a response after a delay$`, iShouldReceiveAResponseAfterADelay)
}

func TestMain(m *testing.M) {
	opts := godog.Options{
		Format: "pretty",             // Output format
		Paths:  []string{"features"}, // Path to the feature files
	}
	status := godog.TestSuite{
		Name:                "externalAPI",
		ScenarioInitializer: InitializeScenario,
		Options:             &opts,
	}.Run()

	if st := m.Run(); st > status {
		status = st
	}
	os.Exit(status)
}
