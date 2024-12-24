package main

import (
	"bytes"
	"fmt"
	"io"
	"maps"
	"os"
	"regexp"
	"slices"
	"sort"
	"strings"

	"github.com/nlm/adventofcode2024/internal/iterators"
	"github.com/nlm/adventofcode2024/internal/stage"
	"github.com/nlm/adventofcode2024/internal/utils"
)

var (
	regRe  = regexp.MustCompile(`(\w{3}): (\d+)`)
	instRe = regexp.MustCompile(`(\w{3}) (AND|OR|XOR) (\w{3}) -> (\w{3})`)
)

type Inst struct {
	Reg1 string
	Op   string
	Reg2 string
	Res  string
}

type InputData struct {
	Regs map[string]bool
	Inst []Inst
}

func ItoB(n int) bool {
	return n != 0
}

func BtoI(b bool) int {
	if b {
		return 1
	}
	return 0
}

func ParseInput(input io.Reader) *InputData {
	parts := bytes.SplitN(utils.Must(io.ReadAll(input)), []byte("\n\n"), 2)
	// stage.Println(string(parts[0]))
	// stage.Println(string(parts[1]))

	inputData := InputData{
		Regs: make(map[string]bool),
	}

	for line := range iterators.MustLines(bytes.NewReader(parts[0])) {
		// stage.Println(">>>", line)
		sparts := regRe.FindStringSubmatch(line)
		if sparts == nil || len(sparts) != 3 {
			panic("parse error")
		}
		// stage.Println("sub", sparts)
		inputData.Regs[sparts[1]] = ItoB(utils.MustAtoi(sparts[2]))
	}

	for line := range iterators.MustLines(bytes.NewReader(parts[1])) {
		// stage.Println(">>>", line)
		sparts := instRe.FindStringSubmatch(line)
		if sparts == nil || len(sparts) != 5 {
			panic("parse error")
		}
		// stage.Println("sub", sparts)
		inputData.Inst = append(inputData.Inst, Inst{
			Reg1: sparts[1],
			Op:   sparts[2],
			Reg2: sparts[3],
			Res:  sparts[4],
		})
	}

	return &inputData
}

func Contains[T1 comparable, T2 any](m map[T1]T2, k T1) bool {
	_, ok := m[k]
	return ok
}

func Remove[T any](slice []T, s int) []T {
	return append(slice[:s], slice[s+1:]...)
}

func Stage1(input io.Reader) (any, error) {
	inputData := ParseInput(input)
	stage.Println(inputData)
	regs := maps.Clone(inputData.Regs)
	insts := slices.Clone(inputData.Inst)

	done := false
	for !done {
		done = true
		for i, inst := range insts {
			if Contains(regs, inst.Reg1) && Contains(regs, inst.Reg2) {
				switch inst.Op {
				case "OR":
					regs[inst.Res] = regs[inst.Reg1] || regs[inst.Reg2]
				case "AND":
					regs[inst.Res] = regs[inst.Reg1] && regs[inst.Reg2]
				case "XOR":
					regs[inst.Res] = regs[inst.Reg1] != regs[inst.Reg2]
				default:
					panic(fmt.Sprintln("unknown instruction", inst.Op))
				}
				stage.Printf("solvable: %+v -> %d\n", inst, BtoI(regs[inst.Res]))
				done = false
				insts = Remove(insts, i)
				break
			}
		}
	}
	stage.Println(regs, insts)
	for _, k := range slices.Sorted(maps.Keys(regs)) {
		stage.Printf("%s: %d\n", k, BtoI(regs[k]))
	}

	res := RegToI("z", regs)

	return res, nil
}

func RegToI(b string, regs map[string]bool) int64 {
	var res int64 = 0
	for i := 63; i >= 0; i-- {
		res = (res << 1) + int64(BtoI(regs[fmt.Sprintf("%s%02d", b, i)]))
	}
	return res
}

func RegToBin(s string, regs map[string]bool) string {
	b := strings.Builder{}
	b.Grow(64)
	isZero := true
	for i := 63; i >= 0; i-- {
		if regs[fmt.Sprintf("%s%02d", s, i)] {
			b.WriteByte('1')
			isZero = false
		} else {
			if isZero {
				b.WriteByte(' ')

			} else {
				b.WriteByte('0')
			}
		}
	}
	return b.String()
}

func MakeGraph(file string, insts []Inst) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()

	opColor := map[string]string{
		"AND": "blue",
		"OR":  "red",
		"XOR": "green",
	}
	// tsw XOR wwm -> z05
	// hdt, z05

	f.WriteString("digraph demo {\n")
	for _, inst := range insts {
		f.WriteString(fmt.Sprint(inst.Reg1, " -> ", inst.Res, " [color=\"", opColor[inst.Op], "\"];\n"))
		f.WriteString(fmt.Sprint(inst.Reg2, " -> ", inst.Res, " [color=\"", opColor[inst.Op], "\"];\n"))
		if strings.HasPrefix(inst.Res, "z") && inst.Op != "XOR" {
			fmt.Println("suspect:", inst.Res, "is not of XOR from", inst.Reg1, inst.Reg2)
		}
	}
	f.WriteString("}")
	return nil
}

func Stage2(input io.Reader) (any, error) {
	inputData := ParseInput(input)
	// stage.Println(inputData)
	regs := maps.Clone(inputData.Regs)

	insts := slices.Clone(inputData.Inst)

	// Swaps
	insts = SwapOutputs(insts, Inst{"tsw", "XOR", "wwm", ""}, Inst{"rnk", "OR", "mkq", ""})
	insts = SwapOutputs(insts, Inst{"x09", "AND", "y09", ""}, Inst{"vkd", "XOR", "wqr", ""})
	insts = SwapOutputs(insts, Inst{"dpr", "XOR", "nvv", ""}, Inst{"dpr", "AND", "nvv", ""})
	insts = SwapOutputs(insts, Inst{"y15", "AND", "x15", ""}, Inst{"y15", "XOR", "x15", ""})

	swapped := []string{"hdt", "z05", "z09", "gbf", "z30", "nbf", "mht", "jgt"}
	sort.StringSlice(swapped).Sort()

	// MakeGraph("graph2.dot", insts)

	// Method: using the graph generation, identify patterns that are not normal.
	// All stairs of the graph have the same behavior:
	// xN XOR yN -> aN
	// xN AND yN -> bN
	// bN-1 OR dN-1 -> cN
	// aN XOR cN -> zN
	// aN AND result of last OR -> dN
	// bN OR dN -> cN+1

	done := false
	for !done {
		done = true
		for i, inst := range insts {
			if Contains(regs, inst.Reg1) && Contains(regs, inst.Reg2) {
				switch inst.Op {
				case "OR":
					regs[inst.Res] = regs[inst.Reg1] || regs[inst.Reg2]
				case "AND":
					regs[inst.Res] = regs[inst.Reg1] && regs[inst.Reg2]
				case "XOR":
					regs[inst.Res] = regs[inst.Reg1] != regs[inst.Reg2]
				default:
					panic(fmt.Sprintln("unknown instruction", inst.Op))
				}
				// stage.Printf("solvable: %+v -> %d\n", inst, BtoI(regs[inst.Res]))
				done = false
				insts = Remove(insts, i)
				break
			}
		}
	}
	if RegToI("x", regs)+RegToI("y", regs) == RegToI("z", regs) {
		fmt.Println("OK")
	}

	stage.Println(RegToI("x", regs), "+", RegToI("y", regs), "=", RegToI("z", regs))
	stage.Println(" ", RegToBin("x", regs), "\n+", RegToBin("y", regs), "\n\n=", RegToBin("z", regs))
	DoTheMath("x", "y", "Z", regs)
	stage.Println(">", RegToBin("Z", regs))

	return strings.Join(swapped, ","), nil
}

func SwapOutputs(insts []Inst, inst1, inst2 Inst) []Inst {
	var toSwap []int
	// stage.Println(insts)
	for i, inst := range insts {
		if inst.Reg1 == inst1.Reg1 && inst.Reg2 == inst1.Reg2 && inst.Op == inst1.Op {
			toSwap = append(toSwap, i)
			continue
		}
		if inst.Reg1 == inst2.Reg1 && inst.Reg2 == inst2.Reg2 && inst.Op == inst2.Op {
			toSwap = append(toSwap, i)
			continue
		}
	}
	if len(toSwap) != 2 {
		panic(fmt.Sprintln("swap error", insts, inst1, inst2))
	}
	fmt.Println("SWAPPED:", insts[toSwap[0]].Res, insts[toSwap[1]].Res)
	insts[toSwap[0]].Res, insts[toSwap[1]].Res = insts[toSwap[1]].Res, insts[toSwap[0]].Res
	return insts
}

func DoTheMath(x, y, z string, regs map[string]bool) {
	carry := false
	for i := 0; i < 64; i++ {
		xbit := regs[fmt.Sprintf("%s%02d", x, i)]
		ybit := regs[fmt.Sprintf("%s%02d", y, i)]
		if xbit && ybit {
			// 1 + 1 + c0 -> 0 c1
			// 1 + 1 + c1 -> 1 c1
			regs[fmt.Sprintf("%s%02d", z, i)] = carry
			carry = true
			continue
		}
		if xbit || ybit {
			// 1 + 0 + c0 -> 1 c0
			// 1 + 0 + c1 -> 0 c1
			regs[fmt.Sprintf("%s%02d", z, i)] = !carry
			continue
		}
		// 0 + 0 + c0 -> 0 c0
		// 0 + 0 + c1 -> 1 c0
		regs[fmt.Sprintf("%s%02d", z, i)] = carry
		carry = false
		continue
	}
}
