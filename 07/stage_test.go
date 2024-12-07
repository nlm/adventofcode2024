package main

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConcat(t *testing.T) {
	for _, tc := range [][3]int{
		{12, 1, 2},
		{10, 1, 0},
		{1, 0, 1},
		{11, 1, 1},
		{0, 0, 0},
		{1234, 123, 4},
		{1234, 12, 34},
		{1234, 1, 234},
	} {
		assert.Equal(t, tc[0], Concat(tc[1], tc[2]))
	}
}

func TestStringConcat(t *testing.T) {
	for _, tc := range [][3]int{
		{12, 1, 2},
		{10, 1, 0},
		{1, 0, 1},
		{1234, 12, 34},
		{1234, 123, 4},
		{1234, 1, 234},
	} {
		assert.Equal(t, tc[0], StringConcat(tc[1], tc[2]))
	}
}

func BenchmarkConcat(b *testing.B) {
	b.Run("Concat", func(b *testing.B) {
		for range b.N {
			Concat(123, 456)
		}
	})
	b.Run("StringConcat", func(b *testing.B) {
		for range b.N {
			StringConcat(123, 456)
		}
	})
}

func BenchmarkProduct1(b *testing.B) {
	for range b.N {
		slices.Collect(CarthesianProduct([]Op{Add, Mul, Concat}, 10))
	}
}

func TestCarthesianProduct(t *testing.T) {
	result := slices.Collect(CarthesianProduct([]string{"A", "B", "C"}, 2))
	expected := [][]string{
		{"A", "A"},
		{"A", "B"},
		{"A", "C"},
		{"B", "A"},
		{"B", "B"},
		{"B", "C"},
		{"C", "A"},
		{"C", "B"},
		{"C", "C"},
	}
	assert.Equal(t, result, expected)
}
