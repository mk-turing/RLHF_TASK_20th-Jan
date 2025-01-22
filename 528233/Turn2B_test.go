package _28233

import (
	"testing"
)

// Assume we have a ProcessData function defined elsewhere
func ProcessData(input Data) Data {
	// Sample implementation
	return input // Just a placeholder
}

// Define a custom testing interface
type Data interface {
	Equals(other Data) bool
	GetValue() interface{}
}

// Implement the Data interface for various types
type IntData struct {
	Value int
}

func (d IntData) Equals(other Data) bool {
	if od, ok := other.(IntData); ok {
		return d.Value == od.Value
	}
	return false
}

func (d IntData) GetValue() interface{} {
	return d.Value
}

// Similar implementations for StringData, SliceData, and MapData

// Test function
func TestProcessData(t *testing.T) {
	// Define test cases using custom Data types
	tests := []struct {
		input    Data
		expected Data
	}{
		{IntData{123}, IntData{123}},
		{StringData{"string"}, StringData{"string"}},
		{SliceData{[]int{1, 2, 3}}, SliceData{[]int{1, 2, 3}}},
		{MapData{map[string]int{"one": 1}}, MapData{map[string]int{"one": 1}}},
	}

	// Loop through each test case
	for _, test := range tests {
		// Process the data using the custom Data interface
		result := ProcessData(test.input)

		// Use the Equals method for comparison
		if !result.Equals(test.expected) {
			t.Errorf("Mismatch for input %v: got %v, expected %v", test.input.GetValue(), result.GetValue(), test.expected.GetValue())
		}
	}
}
