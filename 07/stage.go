package main

import (
	"flag"
	"fmt"
	"io"
	"iter"
	"math"
	"strings"

	"github.com/nlm/adventofcode2024/internal/iterators"
	"github.com/nlm/adventofcode2024/internal/stage"
	"github.com/nlm/adventofcode2024/internal/utils"
)

// Equation represents an equation to be solved.
type Equation struct {
	Result int
	Parts  []int
}

// ParseEquation parses an equation string (b: a1 a2 a3 ...)
// and returns an Equation.
func ParseEquation(line string) Equation {
	elts := strings.Fields(line)
	return Equation{
		Result: utils.MustAtoi(strings.Replace(elts[0], ":", "", 1)),
		Parts:  iterators.Map(elts[1:], utils.MustAtoi),
	}
}

// Op is an operation on 2 int that returns an int.
type Op func(a, b int) int

func Add(a, b int) int {
	return a + b
}

func Mul(a, b int) int {
	return a * b
}

func StringConcat(a, b int) int {
	return utils.MustAtoi(fmt.Sprintf("%d%d", a, b))
}

func Concat(a, b int) int {
	if a == 0 {
		return b
	}
	if b == 0 {
		return a * 10
	}
	return a*int(math.Pow10(int(math.Log10(float64(b)))+1)) + b

}

func CheckEquation(eq Equation, ops []Op) bool {
	if len(eq.Parts) != len(ops)+1 {
		panic(fmt.Sprint("equation and ops len mismatch: ", eq, ops))
	}
	res := eq.Parts[0]
	for i, op := range ops {
		// optim that works only if all operation increase res
		// if res > eq.Result {
		// 	return false
		// }
		res = op(res, eq.Parts[i+1])
	}
	return res == eq.Result
}

func EquationHasSolutions(eq Equation, possibleOps []Op) bool {
	opsCombinations := CartesianProduct(possibleOps, len(eq.Parts)-1)
	for opsCombination := range opsCombinations {
		if CheckEquation(eq, opsCombination) {
			return true
		}
	}
	return false
}

// EquationHasSolutionsRecursive is an alternative implementation that
// uses a recursive algorithm based on optimizations that rely on the fact
// that operations are only increasing the value of the result.
//
// It uses a tree-like algoritm, where all sub-branches that already have
// a result superior to the target result are ignored.
func EquationHasSolutionsRecursive(eq Equation, possibleOps []Op) bool {
	return Solver{Operators: possibleOps}.Solve(eq)
}

// Solver is an equation solver.
type Solver struct {
	Operators []Op
}

// Solve is a recursive solver for equations.
func (s Solver) Solve(eq Equation) bool {
	return s.checkResult(eq.Result, eq.Parts[0], eq.Parts[1:])
}

func (s Solver) checkResult(targetResult int, currentValue int, remainingParts []int) bool {
	for i := range len(s.Operators) {
		newValue := s.Operators[i](currentValue, remainingParts[0])
		// result is too high, continue
		if newValue > targetResult {
			continue
		}
		// if we have more work to do, go deeper in the tree
		if len(remainingParts[1:]) > 0 {
			if s.checkResult(targetResult, newValue, remainingParts[1:]) {
				return true
			}
			continue
		}
		// we're at a leaf, check if we're good
		if newValue == targetResult {
			return true
		}
	}
	return false
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

var flagRecurse = flag.Bool("no-recursive", false, "use original non-recursive algorithm")

func Stage1(input io.Reader) (any, error) {
	fn := EquationHasSolutionsRecursive
	if *flagRecurse {
		stage.Println("disabling recursive algorithm")
		fn = EquationHasSolutions
	}

	total := 0
	for line := range iterators.MustLines(input) {
		eq := ParseEquation(line)
		if fn(eq, []Op{Add, Mul}) {
			total += eq.Result
		}
	}
	return total, nil
}

func Stage2(input io.Reader) (any, error) {
	fn := EquationHasSolutionsRecursive
	if *flagRecurse {
		stage.Println("disabling recursive algorithm")
		fn = EquationHasSolutions
	}

	total := 0
	for line := range iterators.MustLines(input) {
		eq := ParseEquation(line)
		if fn(eq, []Op{Add, Mul, Concat}) {
			total += eq.Result
		}
	}
	return total, nil
}
