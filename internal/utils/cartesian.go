package utils

import "iter"

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
