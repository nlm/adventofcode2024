package main

import (
	"fmt"
	"io"
	"iter"
	"math"
	"strings"

	"github.com/nlm/adventofcode2024/internal/iterators"
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
	line = strings.Replace(line, ":", "", 1)
	elts := strings.Fields(line)
	return Equation{
		Result: utils.MustAtoi(elts[0]),
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
		panic(fmt.Sprint("malformed equation: ", eq))
	}
	res := eq.Parts[0]
	for i, op := range ops {
		res = op(res, eq.Parts[i+1])
	}
	return res == eq.Result
}

func EquationHasSolutions(eq Equation, possibleOps []Op) bool {
	opsCombinations := CarthesianProduct(possibleOps, len(eq.Parts)-1)
	for opsCombination := range opsCombinations {
		if CheckEquation(eq, opsCombination) {
			return true
		}
	}
	return false
}

func CarthesianProduct[T any](elts []T, n int) iter.Seq[[]T] {
	return func(yield func([]T) bool) {
		if len(elts) == 0 || n <= 0 {
			return
		}
		indexes := make([]int, n)
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
	total := 0
	for line := range iterators.MustLines(input) {
		eq := ParseEquation(line)
		if EquationHasSolutions(eq, []Op{Add, Mul}) {
			total += eq.Result
		}
	}
	return total, nil
}

func Stage2(input io.Reader) (any, error) {
	total := 0
	for line := range iterators.MustLines(input) {
		eq := ParseEquation(line)
		if EquationHasSolutions(eq, []Op{Add, Mul, Concat}) {
			total += eq.Result
		}
	}
	return total, nil
}
