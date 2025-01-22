package _28233

import (
	"reflect"
	"testing"
)

// Assume we have a ProcessData function defined elsewhere
func ProcessData(input interface{}) interface{} {
	// Sample implementation
	return input // Just a placeholder
}

// Helper function to compare if two slices are deeply equal
func slicesEqual(a, b interface{}) bool {
	va := reflect.ValueOf(a)
	vb := reflect.ValueOf(b)

	if va.Len() != vb.Len() {
		return false
	}

	for i := 0; i < va.Len(); i++ {
		if !reflect.DeepEqual(va.Index(i).Interface(), vb.Index(i).Interface()) {
			return false
		}
	}
	return true
}

// Test function
func TestProcessData(t *testing.T) {
	// Creating a cache for types
	typeCache := make(map[reflect.Type]struct{})

	// Define test cases
	tests := []struct {
		input    interface{}
		expected interface{}
	}{
		{123, 123},
		{"string", "string"},
		{[]int{1, 2, 3}, []int{1, 2, 3}},
		{map[string]int{"one": 1}, map[string]int{"one": 1}},
	}

	// Loop through each test case
	for _, test := range tests {
		// Check the cache first to see if we’ve processed this type before
		inputType := reflect.TypeOf(test.input)
		if _, present := typeCache[inputType]; !present {
			// If type is not in cache, process and add to cache
			typeCache[inputType] = struct{}{}
		} else {
			continue
		}

		// Process the data
		result := ProcessData(test.input)

		// Use reflection for comparison if needed
		resultValue := reflect.ValueOf(result)
		expectedValue := reflect.ValueOf(test.expected)

		// If types are different, report error immediately
		if resultValue.Type() != expectedValue.Type() {
			t.Errorf("Type mismatch for input %v: got %s, expected %s", test.input, resultValue.Type(), expectedValue.Type())
			continue
		}

		// For slices and maps, use reflect.DeepEqual or a custom comparison
		switch resultValue.Kind() {
		case reflect.Slice, reflect.Map:
			// Avoid redundant reflection comparisons by using a custom slice/map comparison
			if !slicesEqual(result, test.expected) {
				t.Errorf("Mismatch for input %v: got %v, expected %v", test.input, result, test.expected)
			}
		default:
			// For simple cases, direct equality is faster
			if result != test.expected {
				t.Errorf("Mismatch for input %v: got %v, expected %v", test.input, result, test.expected)
			}
		}
	}
}

func main() {
	// Normally testing is done via 'go test', but you can call the test directly:
	TestProcessData(nil)
}
