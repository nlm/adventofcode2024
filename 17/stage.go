package main

import (
	"fmt"
	"io"
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

func (cmp *Computer) Clear() {
	cmp.IP = 0
	for _, b := range []byte{'A', 'B', 'C'} {
		cmp.Rg[b] = cmp.InitRg[b]
	}
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
	// fmt.Print(v, ",")
}

func Stage1(input io.Reader) (any, error) {
	cmp := ParseInput(input)
	stage.Println(cmp)
	for cmp.Step() {

	}
	stage.Println(cmp)
	stage.Println(strings.Join(cmp.output, ","))
	return 0, nil
}

const MaxUint = ^uint(0)
const MinUint = 0
const MaxInt = int(MaxUint >> 1)
const MinInt = -MaxInt - 1

func Stage2(input io.Reader) (any, error) {
	// 164541002416128
	// 164545925312189
	// 164541017976509
	cmp := ParseInput(input)
	// stage.Println(cmp)
	var i = 0
	go func() {
		for {
			fmt.Println(i)
			time.Sleep(1 * time.Second)
		}
	}()
	n := 16
	// for i = 40000000000000; i < MaxInt; i += 10000000 {
	// for i = 164541002416100; i < MaxInt; i += 10 {
	for i = 164541017972735; i < MaxInt; i += 1 {
		// fmt.Println(i)
		cmp.Clear()
		cmp.Rg['A'] = i
		// stage.Println(cmp)
		for cmp.Step() {
		}
		stage.Println("prog", cmp.Prog, "out", cmp.outInt)
		if len(cmp.outInt) != len(cmp.Prog) {
			i += 10000000
			continue
		}
		if slices.Compare(cmp.Prog[len(cmp.Prog)-n:], cmp.outInt[len(cmp.outInt)-n:]) == 0 {
			return i, nil
		}
		if slices.Compare(cmp.Prog, cmp.outInt) == 0 {
			fmt.Println("OK")
			return i, nil
		}
	}
	// stage.Println(cmp)
	// stage.Println(strings.Join(cmp.output, ","))
	return 0, nil
}
