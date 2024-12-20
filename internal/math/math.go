package math

import (
	"golang.org/x/exp/constraints"
)

type Number interface {
	constraints.Integer | constraints.Float
}

func Abs[T Number](n T) T {
	if n < 0 {
		return -n
	}
	return n
}
