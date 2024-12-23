package main

import (
	"io"
	"iter"
	"maps"
	"slices"
	"sort"
	"strings"

	"github.com/nlm/adventofcode2024/internal/iterators"
	"github.com/nlm/adventofcode2024/internal/sets"
	"github.com/nlm/adventofcode2024/internal/stage"
)

// func StrToId(s string) int64 {
// 	b := []byte(s)
// 	if len(b) != 2 {
// 		panic("X2")
// 	}
// 	return int64(b[0])<<8 + int64(b[1])
// }

// func IdToStr(id int64) string {
// 	return string(byte(id>>8)) + string(byte(id&0xff))
// }

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

func BuildGroups(input io.Reader) sets.Set[[3]string] {
	linksMap := make(map[string][]string)

	for line := range iterators.MustLines(input) {
		parts := strings.SplitN(line, "-", 2)
		if len(parts) != 2 {
			panic("X")
		}
		sort.StringSlice(parts).Sort()
		linksMap[parts[0]] = append(linksMap[parts[0]], parts[1])
		linksMap[parts[1]] = append(linksMap[parts[1]], parts[0])
	}

	groups := make(sets.Set[[3]string], 0)
	for _, node1 := range slices.Sorted(maps.Keys(linksMap)) {
		nodes := linksMap[node1]
		for neighs := range CartesianProduct(nodes, 2) {
			// stage.Println("consider", node1, "with", neighs)
			if slices.Contains(linksMap[neighs[0]], neighs[1]) {
				// stage.Println("OK", node1, neighs)
				groups.Add([3]string(slices.Sorted(slices.Values([]string{node1, neighs[0], neighs[1]}))))
			}
		}
	}
	return groups
}

func HasHistorian(group [3]string) bool {
	return strings.HasPrefix(group[0], "t") || strings.HasPrefix(group[1], "t") || strings.HasPrefix(group[2], "t")
}

func Stage1(input io.Reader) (any, error) {
	historians := 0
	groups := BuildGroups(input)
	for group := range maps.Keys(groups) {
		if HasHistorian(group) {
			historians++
			stage.Println("GROUP", group)
		}
	}
	return historians, nil
}

func Stage2(input io.Reader) (any, error) {
	linksMap := make(map[string][]string)
	directLinks := make(sets.Set[[2]string])
	for line := range iterators.MustLines(input) {
		parts := strings.SplitN(line, "-", 2)
		if len(parts) != 2 {
			panic("X")
		}
		sort.StringSlice(parts).Sort()
		// list all possible direct links
		linksMap[parts[0]] = append(linksMap[parts[0]], parts[1])
		// cache direct links
		directLinks.Add([2]string{parts[0], parts[1]})
		directLinks.Add([2]string{parts[1], parts[0]})
	}

	var longest sets.Set[string]
	for k := range linksMap {
		s := make(sets.Set[string])
		for seen := range (Finder{Links: linksMap}).FindChains(k, s) {
			// skip short sequences
			if longest != nil && len(seen) <= len(longest) {
				continue
			}
			// check that every item in the chain is connected to others
			ok := true
			elts := slices.Collect(maps.Keys(seen))
			for duo := range CartesianProduct(elts, 2) {
				if duo[0] == duo[1] {
					continue
				}
				if !directLinks.Contains([2]string(duo)) {
					ok = false
					break
				}
			}
			if ok {
				longest = seen
			}
		}
	}

	codeword := slices.Collect(maps.Keys(longest))
	sort.StringSlice(codeword).Sort()
	return strings.Join(codeword, ","), nil
}

type Finder struct {
	Links map[string][]string
}

func (f Finder) FindChains(from string, seen sets.Set[string]) iter.Seq[sets.Set[string]] {
	return func(yield func(sets.Set[string]) bool) {
		if seen.Contains(from) {
			return
		}
		seen.Add(from)
		foundNeighs := false
		for _, neigh := range f.Links[from] {
			foundNeighs = true
			for elt := range f.FindChains(neigh, seen.Clone()) {
				if !yield(elt) {
					return
				}
			}
		}
		if !foundNeighs {
			yield(seen)
		}
	}
}
