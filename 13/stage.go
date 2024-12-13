package main

import (
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/nlm/adventofcode2024/internal/iterators"
	"github.com/nlm/adventofcode2024/internal/stage"
	"github.com/nlm/adventofcode2024/internal/utils"
)

type Vec struct {
	X int64
	Y int64
}

type Machine struct {
	A Vec
	B Vec
	T Vec
}

var (
	buttonRe = regexp.MustCompile(`Button (\S): X\+(\d+), Y\+(\d+)`)
	prizeRe  = regexp.MustCompile(`Prize: X=(\d+), Y=(\d+)`)
)

func ParseInput(input io.Reader) ([]Machine, error) {
	machines := make([]Machine, 0)
	currentMachine := Machine{}
	for line := range iterators.MustLines(input) {
		switch {
		case strings.HasPrefix(line, "Button A"):
			sm := buttonRe.FindStringSubmatch(line)
			if sm == nil {
				return nil, fmt.Errorf("match error: %s", line)
			}
			currentMachine.A.X = int64(utils.MustAtoi(sm[2]))
			currentMachine.A.Y = int64(utils.MustAtoi(sm[3]))
		case strings.HasPrefix(line, "Button B"):
			sm := buttonRe.FindStringSubmatch(line)
			if sm == nil {
				return nil, fmt.Errorf("match error: %s", line)
			}
			currentMachine.B.X = int64(utils.MustAtoi(sm[2]))
			currentMachine.B.Y = int64(utils.MustAtoi(sm[3]))
		case strings.HasPrefix(line, "Prize"):
			sm := prizeRe.FindStringSubmatch(line)
			if sm == nil {
				return nil, fmt.Errorf("match error: %s", line)
			}
			currentMachine.T.X = int64(utils.MustAtoi(sm[1]))
			currentMachine.T.Y = int64(utils.MustAtoi(sm[2]))
			machines = append(machines, currentMachine)
			currentMachine = Machine{}
		}
	}
	return machines, nil
}

func (m *Machine) Check(a, b int64) bool {
	return m.A.X*a+m.B.X*b == m.T.X && m.A.Y*a+m.B.Y*b == m.T.Y
}

func Solve1(m Machine) (bool, int64) {
	found := false
	cheapest := int64(0)
	for a := int64(0); a <= 100; a++ {
		for b := int64(0); b < 100; b++ {
			if m.Check(a, b) {
				price := int64(a*3 + b)
				if !found {
					found = true
					cheapest = price
				} else {
					if cheapest > price {
						cheapest = price
					}
				}
			}
		}
	}
	return found, cheapest
}

func Stage1(input io.Reader) (any, error) {
	machines, err := ParseInput(input)
	if err != nil {
		return nil, err
	}
	total := int64(0)
	for _, machine := range machines {
		found, price := Solve1(machine)
		if found {
			total += price
		}
	}
	return total, nil
}

func Solve2(m Machine) (bool, int64) {
	a := (m.B.Y*m.T.X - m.B.X*m.T.Y) / (m.B.Y*m.A.X - m.B.X*m.A.Y)
	b := (m.A.X*m.T.Y - m.A.Y*m.T.X) / (m.A.X*m.B.Y - m.A.Y*m.B.X)
	if !m.Check(a, b) {
		stage.Println("error")
		return false, 0
	}
	return true, a*3 + b
}

// Button A: X+94, Y+34
// Button B: X+22, Y+67
// Prize: X=10000000008400, Y=10000000005400
// (Ax * a) + (Bx * b) = Px
// (Ay * a) + (By * b) = Py
//
// (Ax * a) + (Bx * b) = Px

func Stage2(input io.Reader) (any, error) {
	machines, err := ParseInput(input)
	if err != nil {
		return nil, err
	}
	total := int64(0)
	for _, machine := range machines {
		machine.T.X = machine.T.X + 10000000000000
		machine.T.Y = machine.T.Y + 10000000000000
		stage.Println(machine)
		found, price := Solve2(machine)
		if found {
			stage.Println("found", price)
			total += price
		}
	}
	return total, nil
}
