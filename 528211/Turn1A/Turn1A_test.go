package Turn1A

import "testing"

func TestAdd(t *testing.T) {
	tests := []struct {
		name     string
		a, b     int
		expected int
	}{
		{"Add positives", 1, 2, 3},
		{"Add negatives", -1, -2, -3},
		{"Add zero", 0, 5, 5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Add(tt.a, tt.b); got != tt.expected {
				t.Errorf("Add(%d, %d) = %d, want %d", tt.a, tt.b, got, tt.expected)
			}
		})
	}
}

func TestMultiply(t *testing.T) {
	tests := []struct {
		name     string
		a, b     int
		expected int
	}{
		{"Multiply positives", 2, 3, 6},
		{"Multiply by zero", 0, 4, 0},
		{"Multiply negatives", -2, -3, 6},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Multiply(tt.a, tt.b); got != tt.expected {
				t.Errorf("Multiply(%d, %d) = %d, want %d", tt.a, tt.b, got, tt.expected)
			}
		})
	}
}

func TestDivide(t *testing.T) {
	tests := []struct {
		name     string
		a, b     int
		expected int
	}{
		{"Divide positives", 6, 3, 2},
		{"Divide by zero", 6, 0, 0},
		{"Divide negatives", -6, -3, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Divide(tt.a, tt.b); got != tt.expected {
				t.Errorf("Divide(%d, %d) = %d, want %d", tt.a, tt.b, got, tt.expected)
			}
		})
	}
}
