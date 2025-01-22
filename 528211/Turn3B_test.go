package _28211

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouter_GET_WithParams(t *testing.T) {
	r := Default()
	r.GET("/user/:name", func(c *Context) {
		c.String(http.StatusOK, "Hello %s", c.Param("name"))
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/user/John", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", w.Code)
	}
	expected := "Hello John"
	if w.Body.String() != expected {
		t.Errorf("Expected body '%s', got '%s'", expected, w.Body.String())
	}

	// Add more test cases with different parameter values
	// ...
}
