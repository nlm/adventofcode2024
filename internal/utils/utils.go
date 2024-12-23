package utils

import (
	"strconv"
	"strings"

	"golang.org/x/exp/constraints"
)

// Must takes any type of result and error.
// It panics if the error is non-nil.
// Otherwise, it returns the result.
func Must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

// MustAtoi converts a string to an integer.
// It panics if it fails.
func MustAtoi(s string) int {
	return Must(strconv.Atoi(strings.TrimSpace(s)))
}

// MustAtoi converts a string to an integer.
// It panics if it fails.
func MustAtoInt[T constraints.Integer](s string) T {
	return T(Must(strconv.Atoi(strings.TrimSpace(s))))
}

// NoError panics if an error is present
func MustNoErr(err error) {
	if err != nil {
		panic(err)
	}
}
