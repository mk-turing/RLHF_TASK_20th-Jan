package Turn3A

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestReadFile(t *testing.T) {
	// Setup temporary file with content.
	file, err := ioutil.TempFile("", "example")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())

	if _, err := file.Write([]byte("hello world")); err != nil {
		t.Fatal(err)
	}
	file.Close()

	// Test reading the file correctly.
	content, err := ReadFile(file.Name())
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if content != "hello world" {
		t.Errorf("Expected 'hello world', got %s", content)
	}

	// Test reading a non-existent file.
	_, err = ReadFile("non-existent-file")
	if err == nil {
		t.Error("Expected an error, got nil")
	}

	// Test reading an empty file.
	emptyFile, _ := ioutil.TempFile("", "empty")
	defer os.Remove(emptyFile.Name())
	emptyContent, err := ReadFile(emptyFile.Name())
	if err == nil || emptyContent != "" {
		t.Error("Expected an error for empty file")
	}
}
