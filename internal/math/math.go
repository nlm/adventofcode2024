package math

import (
	"math"

	"golang.org/x/exp/constraints"
)

type Number interface {
	constraints.Integer | constraints.Float
}

//go:inline
func Abs[T Number](n T) T {
	return T(math.Abs(float64(n)))
}
