package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

const (
	apiURL = "http://localhost:2525/imposter/weather"
	normalDelay = 20 * time.Millisecond  // Normal response time
)

func simulateNetworkIssue(t *testing.T, delay time.Duration) {
	// Start Mountebank imposter to simulate network issue
	resp, err := http.Post("http://localhost:2525/imposters", "application/json", `
   {
      "protocol": "http",
      "port": 8000,
      "name": "weather",
      "stubs": [
         {
            "predicates": [{ "equals": { "path": "/data/2.5/weather?q=London&appid=YOUR_API_KEY" } }],
            "responses": [
               {
                  "is": {
                     "statusCode": 200,
                     "headers": { "Content-Type": "application/json" },
                     "body": "Sunny"
                  },
                  "delay": { "time": delay }
               }
            ]
         }
      ]
   }`)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	// Simulate network issue by making HTTP requests to the imposter
	for i := 0; i < t.N; i++ {
		http.Get(apiURL)
	}
}

func BenchmarkNetworkIssue(b *testing.B) {
	b.Run("NormalDelay", func(b *testing.B) {
		simulateNetworkIssue(b, normalDelay)
	})

	b.Run("HighDelay", func(b *testing.B) {
		simulateNetworkIssue(b, normalDelay * 5)  // Simulate 5 times the normal delay
	})

	b.Run("ConnectionTimeout", func(b *testing.B) {
		// Simulate connection timeout
		http.DefaultClient.Timeout = 10 * time.Millisecond
		simulateNetworkIssue(b, normalDelay)
	})
}

func main() {
	// Start Mountebank server
	go func() {
		if err := http.ListenAndServe(":2525", nil); err != nil {
			fmt.Println(err)
		}
	}()

	// Test scenarios
	testing.RunBenchmarks(nil, "BenchmarkNetworkIssue")