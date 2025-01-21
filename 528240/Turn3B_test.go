package _28240

import (
	"fmt"
	"github.com/cucumber/godog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"os"
	"testing"
)

var apiClient *MockAPIClient
var database *MockDatabase

func iHaveAMockedAPIClient() error {
	apiClient = new(MockAPIClient)
	return nil
}

func theAPIClientReturnsUserDataForUser(userID string) error {
	apiClient.On("GetUserData", userID).Return(map[string]interface{}{"user_id": userID, "name": "John Doe"}, nil)
	return nil
}

func iHaveAMockedDatabase() error {
	database = new(MockDatabase)
	return nil
}

func iRetrieveUserDataForUser(userID string) error {
	_, err := apiClient.GetUserData(userID)
	return err
}

func iSaveTheUserDataToTheDatabase() error {
	data := map[string]interface{}{"user_id": "user123", "name": "John Doe"}
	return database.InsertRecord(data)
}

func theDatabaseShouldContainUserDataForUser(userID string) error {
	database.AssertCalled(t, "InsertRecord", map[string]interface{}{"user_id": userID, "name": "John Doe"})
	return nil
}

func InitializeScenario(ctx *godog.ScenarioContext) {