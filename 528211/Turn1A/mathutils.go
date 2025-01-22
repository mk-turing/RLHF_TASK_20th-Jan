package Turn1A

// Add adds two integers.
func Add(a, b int) int {
	return a + b
}

// Multiply multiplies two integers.
func Multiply(a, b int) int {
	return a * b
}

// Divide divides two integers. Returns 0 if dividing by zero.
func Divide(a, b int) int {
	if b == 0 {
		return 0
	}
	return a / b
}
