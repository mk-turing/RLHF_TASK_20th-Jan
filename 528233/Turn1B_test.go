package _28233

import (
	"reflect"
	"testing"
)

func ProcessData(data interface{}) interface{} {
	switch data.(type) {
	case int:
		return data.(int) + 1
	case string:
		return data.(string) + "!"
	case []int:
		result := make([]int, len(data.([]int)))
		copy(result, data.([]int))
		for i := range result {
			result[i]++
		}
		return result
	case []string:
		result := make([]string, len(data.([]string)))
		copy(result, data.([]string))
		for i := range result {
			result[i] += "!"
		}
		return result
	default:
		return data
	}
}

func TestProcessData(t *testing.T) {
	testCases := []struct {
		input    interface{}
		expected interface{}
	}{
		{1, 2},
		{"hello", "hello!"},
		{[]int{1, 2, 3}, []int{2, 3, 4}},
		{[]string{"hello", "world"}, []string{"hello!", "world!"}},
	}

	for _, tc := range testCases {
		result := ProcessData(tc.input)
		if !reflect.DeepEqual(result, tc.expected) {
			t.Errorf("ProcessData(%v) returned %v, expected %v", tc.input, result, tc.expected)
		}
	}
}
