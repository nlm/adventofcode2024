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
