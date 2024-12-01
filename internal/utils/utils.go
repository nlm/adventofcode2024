package utils

import (
	"strconv"
	"strings"
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

// Map returns a new slice of the same size as t1,
// containing the result of fn applied on each value of t1.
func Map[T1, T2 any](t1 []T1, fn func(T1) T2) []T2 {
	var t2 = make([]T2, len(t1))
	for i := 0; i < len(t1); i++ {
		t2[i] = fn(t1[i])
	}
	return t2
}

// Filter returns a new slice containing elements of s1
// if fn returns true for that element.
func Filter[T any](s1 []T, fn func(T) bool) []T {
	var s2 = make([]T, 0, len(s1))
	for i := 0; i < len(s1); i++ {
		if fn(s1[i]) {
			s2 = append(s2, s1[i])
		}
	}
	return s2
}

// All returns true if fn returns true for all of
// the items contained in the slice.
func All[T any](slice []T, fn func(T) bool) bool {
	for i := 0; i < len(slice); i++ {
		if !fn(slice[i]) {
			return false
		}
	}
	return true
}

// Any returns true if fn returns true for any of
// the items contained in the slice.
func Any[T any](slice []T, fn func(T) bool) bool {
	for i := 0; i < len(slice); i++ {
		if fn(slice[i]) {
			return true
		}
	}
	return false
}
