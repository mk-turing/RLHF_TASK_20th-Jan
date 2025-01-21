package _93898

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func SampleFunction(n int) int {
	sum := 0
	for i := 0; i < n; i++ {
		sum += i
	}
	return sum
}

func BenchmarkSampleFunction(b *testing.B) {
	for i := 0; i < b.N; i++ {
		SampleFunction(1000)
	}
}

func BenchmarkSampleFunctionWithTestify(b *testing.B) {
	for i := 0; i < b.N; i++ {
		result := SampleFunction(1000)
		assert.NotNil(b, result)
		assert.Greater(b, result, 0)
	}
}
