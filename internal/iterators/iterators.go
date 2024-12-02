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

func Map[T1, T2 any](items []T1, f func(T1) T2) []T2 {
	result := make([]T2, len(items))
	for i := range len(items) {
		result[i] = f(items[i])
	}
	return result
}
