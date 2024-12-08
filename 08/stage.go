package main

import (
	"io"
	"iter"

	"github.com/nlm/adventofcode2024/internal/matrix"
	"github.com/nlm/adventofcode2024/internal/sets"
	"github.com/nlm/adventofcode2024/internal/stage"
	"github.com/nlm/adventofcode2024/internal/utils"
)

func ParseLocations(m *matrix.Matrix[byte]) map[byte][]matrix.Coord {
	result := make(map[byte][]matrix.Coord)
	for c := range m.Coords() {
		v := m.AtCoord(c)
		if v != '.' && v != '#' {
			result[v] = append(result[v], c)
		}
	}
	return result
}

// Pairs generate unique pairs from the elements:
// [a, b, c, d] => a [b, c, d] + b [c, d] + c [d]
// func Pairs[T any](elts []T) iter.Seq[[2]T] {
// 	return func(yield func([2]T) bool) {
// 		if len(elts) < 2 {
// 			return
// 		}
// 		if len(elts) == 2 {
// 			if !yield([2]T{elts[0], elts[1]}) {
// 				return
// 			}
// 			return
// 		}
// 		for _, elt := range elts[1:] {
// 			if !yield([2]T{elts[0], elt}) {
// 				return
// 			}
// 		}
// 		for pair := range Pairs(elts[1:]) {
// 			if !yield(pair) {
// 				return
// 			}
// 		}
// 	}
// }

// func Swaps[T any](elts [2]T) iter.Seq[[2]T] {
// 	return func(yield func([2]T) bool) {
// 		if !yield([2]T{elts[0], elts[1]}) {
// 			return
// 		}
// 		if !yield([2]T{elts[1], elts[0]}) {
// 			return
// 		}
// 	}
// }

func CountAntinodes(m *matrix.Matrix[byte]) int {
	locs := ParseLocations(m)
	antinodes := make(sets.Set[matrix.Coord])
	// m2 is only used for visualization
	m2 := m.Clone()
	m2.Fill('.')
	for name, antennas := range locs {
		stage.Println("considering antennas of type", string(name), "->", antennas)
		// shamelessly reuse last day's code
		for pair := range CartesianProduct(antennas, 2) {
			pair := [2]matrix.Coord(pair)
			// skip combinations with 2 identical items
			if pair[0] == pair[1] {
				continue
			}
			// --- visualization
			m2.SetAtCoord(pair[0], name)
			m2.SetAtCoord(pair[1], name)
			stage.Println(" testing pair", pair)
			// checking only one way, because CartesianProduct generates [a, b] and [b, a]
			c := pair[1].Add(pair[1].Sub(pair[0]))
			stage.Println("  antinode at", c, "->", m.InCoord(c))
			if m.InCoord(c) {
				antinodes.Add(c)
				m2.SetAtCoord(c, '#')
			}
			stage.Println(matrix.SMatrix(m2))
		}
	}
	return len(antinodes)
}

func CountAntinodes2(m *matrix.Matrix[byte]) int {
	locs := ParseLocations(m)
	antinodes := make(sets.Set[matrix.Coord])
	// m2 is only used for visualization
	m2 := m.Clone()
	m2.Fill('.')
	for name, antennas := range locs {
		stage.Println("considering", string(name), antennas)
		// shamelessly reuse last day's code
		for pair := range CartesianProduct(antennas, 2) {
			pair := [2]matrix.Coord(pair)
			// skip combinations with 2 identical items
			if pair[0] == pair[1] {
				continue
			}
			// visualization
			m2.SetAtCoord(pair[0], name)
			m2.SetAtCoord(pair[1], name)
			stage.Println(" testing pair", pair)
			// add itself for this stage
			antinodes.Add(pair[1])
			// checking only one way, because CartesianProduct generates [a, b] and [b, a]
			v := pair[1].Sub(pair[0])
			for c := pair[1].Add(v); m.InCoord(c); c = c.Add(v) {
				stage.Println("  antinode at", c)
				antinodes.Add(c)
				m2.SetAtCoord(c, '#')
			}
			stage.Println(matrix.SMatrix(m2))
		}
	}
	return len(antinodes)
}

// CartesianProduct iterates over every combination of n items from the elts list.
// It will allocate "1 + len(elts) to the power of n" slices.
//
// Internally it's using an array of indexes to iterate over every possible
// combination of positions in the elts list. It's using basic math to pass to
// the next possiblity by adding 1 to the lowest weight index and propagating
// the carry to a higher index, until all solutions have been interated over.
func CartesianProduct[T any](elts []T, n int) iter.Seq[[]T] {
	return func(yield func([]T) bool) {
		if len(elts) == 0 || n <= 0 {
			return
		}
		indexes := make([]int, n)
		// faster version with buffer reuse
		// curr := make([]T, n)
		for {
			// return current combination
			curr := make([]T, n)
			for i := range n {
				curr[i] = elts[indexes[i]]
			}
			if !yield(curr) {
				return
			}
			// calculate indexes for next combination
			carry := 0
			for i := n - 1; i >= 0; i-- {
				indexes[i]++
				if indexes[i] >= len(elts) {
					carry++
					indexes[i] = 0
				} else {
					break
				}
			}
			// return if every bit triggered carry
			if carry == n {
				return
			}
		}
	}
}

func Stage1(input io.Reader) (any, error) {
	m := utils.Must(matrix.NewFromReader(input))
	return CountAntinodes(m), nil
}

func Stage2(input io.Reader) (any, error) {
	m := utils.Must(matrix.NewFromReader(input))
	return CountAntinodes2(m), nil
}
