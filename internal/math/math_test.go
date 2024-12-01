package math

import (
	"math"
	"testing"
)

func BenchmarkStdAbs(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = int(math.Abs(float64(i)))
		_ = int(math.Abs(float64(-i)))
	}
}

func BenchmarkAbs(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = Abs(i)
		_ = Abs(-i)
	}
}
