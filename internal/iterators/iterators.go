package iterators

import (
	"bufio"
	"io"
	"iter"
)

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

func Map[T1, T2 any](items iter.Seq[T1], f func(T1) T2) iter.Seq[T2] {
	return func(yield func(T2) bool) {
		for item := range items {
			if !yield(f(item)) {
				return
			}
		}
	}
}

func MapSlice[T1, T2 any](items []T1, f func(T1) T2) []T2 {
	res := make([]T2, len(items))
	for i := range len(items) {
		res[i] = f(items[i])
	}
	return res
}

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

func FilterSlice[T1 any](items []T1, f func(T1) bool) []T1 {
	res := make([]T1, 0, len(items))
	for _, item := range items {
		if f(item) {
			res = append(res, item)
		}
	}
	return res
}
