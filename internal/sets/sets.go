package sets

import (
	"iter"
	"maps"
)

type Set[T comparable] map[T]struct{}

func (s Set[T]) Contains(value T) bool {
	_, ok := s[value]
	return ok
}

func (s Set[T]) Add(values ...T) {
	for _, v := range values {
		s[v] = struct{}{}
	}
}

func Append[T comparable](s Set[T], values ...T) Set[T] {
	if s == nil {
		s = Set[T]{}
	}
	s.Add(values...)
	return s
}

func Values[T comparable](s Set[T]) iter.Seq[T] {
	return maps.Keys(s)
}

func (s Set[T]) Remove(v T) {
	delete(s, v)
}
