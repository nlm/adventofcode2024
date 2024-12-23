package iterators

import (
	"bufio"
	"io"
	"iter"
)

// MustLines returns an iteator over the lines read from the io.Reader.
// This version returns slices of bytes.
func MustLinesBytes(r io.Reader) iter.Seq[[]byte] {
	return func(yield func([]byte) bool) {
		s := bufio.NewScanner(r)
		for s.Scan() {
			if !yield(s.Bytes()) {
				return
			}
		}
		if s.Err() != nil {
			panic(s.Err())
		}
	}
}

// MustLines returns an iteator over the lines read from the io.Reader.
// This version returns strings.
func MustLines(r io.Reader) iter.Seq[string] {
	return func(yield func(string) bool) {
		s := bufio.NewScanner(r)
		for s.Scan() {
			if !yield(s.Text()) {
				return
			}
		}
		if s.Err() != nil {
			panic(s.Err())
		}
	}
}

// Map returns an iterator that outputs the result of the function f applied
// to each item of the given Sequence.
func Map[T1, T2 any](items iter.Seq[T1], f func(T1) T2) iter.Seq[T2] {
	return func(yield func(T2) bool) {
		for item := range items {
			if !yield(f(item)) {
				return
			}
		}
	}
}

// Map returns a slice of which each items contains the result
// of the function f applied to each item from the provided slice.
func MapSlice[T1, T2 any](items []T1, f func(T1) T2) []T2 {
	res := make([]T2, len(items))
	for i := range len(items) {
		res[i] = f(items[i])
	}
	return res
}

// Filter returns an iteartor that outputs each item from the given
// sequence that satisfy the condition "f(item) == true".
func Filter[T1 any](items iter.Seq[T1], f func(T1) bool) iter.Seq[T1] {
	return func(yield func(T1) bool) {
		for item := range items {
			if f(item) {
				if !yield(item) {
					return
				}
			}
		}
	}
}

// Filter returns a slice that contains all items from the given
// slice that satisfy the condition "f(item) == true".
func FilterSlice[T1 any](items []T1, f func(T1) bool) []T1 {
	res := make([]T1, 0, len(items))
	for _, item := range items {
		if f(item) {
			res = append(res, item)
		}
	}
	return res
}

// All returns true if fn returns true for all of
// the items contained in the slice.
func All[T any](seq iter.Seq[T], fn func(T) bool) bool {
	for elt := range seq {
		if !fn(elt) {
			return false
		}
	}
	return true
}

// All returns true if fn returns true for all of
// the items contained in the slice.
func AllSlice[T any](slice []T, fn func(T) bool) bool {
	for i := 0; i < len(slice); i++ {
		if !fn(slice[i]) {
			return false
		}
	}
	return true
}

// Any returns true if fn returns true for any of
// the items produced by the iterator.
func Any[T any](seq iter.Seq[T], fn func(T) bool) bool {
	for elt := range seq {
		if fn(elt) {
			return true
		}
	}
	return false
}

// Any returns true if fn returns true for any of
// the items contained in the slice.
func AnySlice[T any](slice []T, fn func(T) bool) bool {
	for i := 0; i < len(slice); i++ {
		if fn(slice[i]) {
			return true
		}
	}
	return false
}
