package main

import (
	"fmt"
	"io"
	"maps"
	"math"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/nlm/adventofcode2024/internal/iterators"
	"github.com/nlm/adventofcode2024/internal/stage"
	"github.com/nlm/adventofcode2024/internal/utils"
)

type Computer struct {
	IP     int
	InitRg map[byte]int
	Rg     map[byte]int
	Prog   []int
	output []string
	outInt []int
}

func (cmp *Computer) Reset() {
	cmp.IP = 0
	maps.Copy(cmp.Rg, cmp.InitRg)
	// for _, b := range []byte{'A', 'B', 'C'} {
	// 	cmp.Rg[b] = cmp.InitRg[b]
	// }
	cmp.outInt = cmp.outInt[:0]
	cmp.output = cmp.output[:0]
}

var registerRe = regexp.MustCompile(`^Register (\w): (\d+)`)
var programRe = regexp.MustCompile(`Program: ([\d,]+)`)

func ParseInput(input io.Reader) Computer {
	cmp := Computer{
		Rg:     make(map[byte]int, 3),
		InitRg: make(map[byte]int, 3),
	}
	for line := range iterators.MustLines(input) {
		reg := registerRe.FindStringSubmatch(line)
		if reg != nil {
			cmp.InitRg[byte(reg[1][0])] = utils.MustAtoi(reg[2])
			cmp.Rg[byte(reg[1][0])] = utils.MustAtoi(reg[2])
			continue
		}
		program := programRe.FindStringSubmatch(line)
		if program != nil {
			for _, n := range strings.Split(program[1], ",") {
				cmp.Prog = append(cmp.Prog, utils.MustAtoi(n))
			}
		}
	}
	return cmp
}

func (cp *Computer) ComboOp(x int) int {
	switch x {
	case 0, 1, 2, 3:
		return x
	case 4:
		return cp.Rg['A']
	case 5:
		return cp.Rg['B']
	case 6:
		return cp.Rg['C']
	default:
		panic("operand error")
	}
}

func (cp *Computer) Step() bool {
	if cp.IP > len(cp.Prog)-2 {
		return false
	}
	op, arg := cp.Prog[cp.IP], cp.Prog[cp.IP+1]
	// stage.Println(">>", "op", op, "arg", arg)
	switch op {
	case 0: // adv
		cp.Rg['A'] = int(float64(cp.Rg['A']) / math.Pow(float64(2), float64(cp.ComboOp(arg))))
		cp.IP += 2
	case 1: // bxl
		cp.Rg['B'] = cp.Rg['B'] ^ arg
		cp.IP += 2
	case 2: // bst
		cp.Rg['B'] = cp.ComboOp(arg) % 8
		cp.IP += 2
	case 3: // jnz
		if cp.Rg['A'] == 0 {
			cp.IP += 2
			break
		}
		cp.IP = arg
	case 4: // bxc
		cp.Rg['B'] = cp.Rg['B'] ^ cp.Rg['C']
		cp.IP += 2
	case 5: // out
		cp.Output(cp.ComboOp(arg) % 8)
		cp.IP += 2
	case 6: // bdv
		cp.Rg['B'] = int(float64(cp.Rg['A']) / math.Pow(float64(2), float64(cp.ComboOp(arg))))
		cp.IP += 2
	case 7: // cdv
		cp.Rg['C'] = int(float64(cp.Rg['A']) / math.Pow(float64(2), float64(cp.ComboOp(arg))))
		cp.IP += 2
	}
	return true
}

func (cp *Computer) Output(v int) {
	cp.output = append(cp.output, fmt.Sprint(v))
	cp.outInt = append(cp.outInt, v)
	// stage.Println("Out:", v)
}

func Stage1(input io.Reader) (any, error) {
	cmp := ParseInput(input)
	stage.Println(cmp)
	for cmp.Step() {
	}
	stage.Println(cmp)
	return strings.Join(cmp.output, ","), nil
}

const MaxUint = ^uint(0)
const MinUint = 0
const MaxInt = int(MaxUint >> 1)
const MinInt = -MaxInt - 1

func Stage2(input io.Reader) (any, error) {
	cmp := ParseInput(input)
	var i = 0
	go func() {
		for {
			fmt.Println(i)
			time.Sleep(1 * time.Second)
		}
	}()
	// Try to match size quickly
	for ; i < MaxInt; i++ {
		cmp.Reset()
		cmp.Rg['A'] = i
		for cmp.Step() {
		}
		stage.Println("prog", cmp.Prog, "out", cmp.outInt)
		if len(cmp.outInt) < len(cmp.Prog) {
			i *= 2
			continue
		}
		i /= 2
		break
	}
	// Precision search
	accuracy := 1
	step := 1000000000
	lastI := i
	for ; i < MaxInt; i += step {
		cmp.Reset()
		cmp.Rg['A'] = i
		for cmp.Step() {
		}
		stage.Println("prog", cmp.Prog, "out", cmp.outInt, accuracy-1, "/", len(cmp.Prog), i)
		if len(cmp.outInt) > len(cmp.Prog) {
			i = lastI
			if step > 1 {
				step /= 2
			}
			continue
		}
		if slices.Compare(cmp.Prog[len(cmp.Prog)-accuracy:], cmp.outInt[len(cmp.outInt)-accuracy:]) == 0 {
			if accuracy == len(cmp.Prog) {
				return i, nil
			}
			accuracy++
			i = lastI
			if step > 1 {
				step /= 10
			}
			continue
		}
		// Comment to make it work on the example
		lastI = i
	}
	return 0, nil
}
