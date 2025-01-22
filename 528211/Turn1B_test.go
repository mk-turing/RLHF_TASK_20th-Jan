package _28211

import (
	"math"
	"math/rand"
	"testing"
)

func sqrt(x float64) float64 { if x < 0 { return 0 } z := x / 2 for i := 0; i < 100; i++ { z -= (z*z - x) / (2 * z) } return z }

func TestSqrt(t *testing.T) { cases := []struct { in, want float64 }{ {4, 2}, {9, 3}, {2, 1.4142135623730951}, {-2, 0}, } for _, c := range cases { got := sqrt(c.in) if got != c.want { t.Errorf("sqrt(%g) == %g, want %g", c.in, got, c.want) } } }

func TestSqrtRandom(t *testing.T) { for i := 0; i < 1000; i++ { x := rand.Float64() * 100 got := sqrt(x) if got != math.Sqrt(x) { t.Errorf("sqrt(%g) == %g, want %g", x, got, math.Sqrt(x)) } } }