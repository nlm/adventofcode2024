package sets

import (
	"fmt"
	"iter"
	"maps"
	"strings"
)

// Set is a collection of unique items.
//
// To create a new set, use make(Set[T], cap).
type Set[T comparable] map[T]struct{}

// Contains returns true if the Set contains the value.
// This operation has O(1) complexity.
func (s Set[T]) Contains(value T) bool {
	_, ok := s[value]
	return ok
}

// Add adds values to an existing set.
func (s Set[T]) Add(values ...T) {
	for _, v := range values {
		s[v] = struct{}{}
	}
}

func (s Set[T]) String() string {
	b := strings.Builder{}
	b.Grow(len(s) * 4)
	first := true
	b.WriteString("set{")
	for k := range s {
		if !first {
			b.WriteString(" ")
		} else {
			first = false
		}
		b.WriteString(fmt.Sprint(k))
	}
	b.WriteString("}")
	return b.String()
}

// Append is the same as Add, but it creates a new set if s is nil.
func Append[T comparable](s Set[T], values ...T) Set[T] {
	if s == nil {
		s = make(Set[T], len(values))
	}
	s.Add(values...)
	return s
}

// Values returns an interator over the Set values.
func Values[T comparable](s Set[T]) iter.Seq[T] {
	return maps.Keys(s)
}

// Remove removes an item from an existing Set.
func (s Set[T]) Remove(v T) {
	delete(s, v)
}

// Clone returns a copy of the Set.
func (s Set[T]) Clone() Set[T] {
	return maps.Clone(s)
}

// Copy copies the items from src to dst.
func Copy[T comparable](dst, src Set[T]) {
	maps.Copy(dst, src)
}

// Equal reports whether two sets contain the same values.
// Values are compared using ==.
func Equal[S Set[T], T comparable](s1, s2 Set[T]) bool {
	return maps.Equal(s1, s2)
}

// Insert adds the items from seq to the Set.
func Insert[T comparable](s Set[T], seq iter.Seq[T]) {
	for v := range seq {
		s[v] = struct{}{}
	}
}

// Collects builds a slice containing all items in the set,
// in no particular order.
func Collect[T comparable](s Set[T]) []T {
	l := make([]T, len(s))
	for v := range s {
		l = append(l, v)
	}
	return l
}
