package main

import (
	"bytes"
	"io"
	"iter"
	"regexp"
	"strings"

	"github.com/nlm/adventofcode2024/internal/iterators"
	"github.com/nlm/adventofcode2024/internal/matrix"
	"github.com/nlm/adventofcode2024/internal/maze"
	"github.com/nlm/adventofcode2024/internal/stage"
	"github.com/nlm/adventofcode2024/internal/utils"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/path"
	"gonum.org/v1/gonum/graph/simple"
)

func ManhattanNeighbors[T comparable](m *matrix.Matrix[T], c matrix.Coord) iter.Seq[matrix.Coord] {
	return func(yield func(matrix.Coord) bool) {
		for _, dir := range []matrix.Vec{matrix.Up, matrix.Right, matrix.Down, matrix.Left} {
			nc := c.Add(dir)
			if m.InCoord(nc) {
				if !yield(nc) {
					return
				}
			}
		}
	}
}

type Pad struct {
	Matrix   *matrix.Matrix[byte]
	Origin   matrix.Coord
	Pointer  matrix.Coord
	Graph    *simple.UndirectedGraph
	KeyIndex map[byte]matrix.Coord
}

func (p *Pad) Reset() {
	p.Pointer = p.Origin
}

func (p *Pad) Current() byte {
	return p.Matrix.AtCoord(p.Pointer)
}

func (p Pad) Clone() *Pad {
	return &p
}

func (p *Pad) SetCurrent(b byte) {
	p.Pointer = p.KeyIndex[b]
}

func (p *Pad) FindPath(from, to byte) []matrix.Coord {
	fromCoord, ok := p.KeyIndex[from]
	if !ok {
		panic("'from' key not found")
	}
	toCoord, ok := p.KeyIndex[to]
	if !ok {
		panic("'to' key not found")
	}
	paths := path.DijkstraFrom(p.Graph.Node(maze.CoordToId(p.Matrix, fromCoord)), p.Graph)
	sp, _ := paths.To(maze.CoordToId(p.Matrix, toCoord))
	return iterators.MapSlice(sp, func(node graph.Node) matrix.Coord { return maze.IdToCoord(p.Matrix, node.ID()) })
}

func (p *Pad) FindAllPath(from, to byte) [][]matrix.Coord {
	fromCoord, ok := p.KeyIndex[from]
	if !ok {
		panic("'from' key not found")
	}
	toCoord, ok := p.KeyIndex[to]
	if !ok {
		panic("'to' key not found")
	}
	// paths := path.DijkstraFrom(p.Graph.Node(maze.CoordToId(p.Matrix, fromCoord)), p.Graph)
	sp, _ := path.DijkstraAllFrom(p.Graph.Node(maze.CoordToId(p.Matrix, fromCoord)), p.Graph).
		AllTo(maze.CoordToId(p.Matrix, toCoord))

	return iterators.MapSlice(sp, func(nodes []graph.Node) []matrix.Coord {
		return iterators.MapSlice(nodes, func(node graph.Node) matrix.Coord {
			return maze.IdToCoord(p.Matrix, node.ID())
		})
	})
}

func NewPad(m *matrix.Matrix[byte], origin matrix.Coord) Pad {
	p := Pad{
		Matrix:   m,
		Pointer:  origin,
		KeyIndex: make(map[byte]matrix.Coord, len(m.Data)),
	}
	g := simple.NewUndirectedGraph()
	for c := range p.Matrix.Coords() {
		// Index
		p.KeyIndex[m.AtCoord(c)] = c
		// Graph
		node1, isNew := g.NodeWithID(maze.CoordToId(p.Matrix, c))
		if isNew {
			g.AddNode(node1)
		}
		for nc := range ManhattanNeighbors(p.Matrix, c) {
			node2, isNew := g.NodeWithID(maze.CoordToId(p.Matrix, nc))
			if isNew {
				g.AddNode(node2)
			}
			g.SetEdge(g.NewEdge(node1, node2))
		}
	}
	p.Graph = g
	return p
}

func NewNPad() Pad {
	return NewPad(
		utils.Must(matrix.NewFromReader(strings.NewReader("789\n456\n123\n 0A\n"))),
		matrix.Coord{X: 2, Y: 3},
	)
}

func NewDPad() Pad {
	return NewPad(
		utils.Must(matrix.NewFromReader(strings.NewReader(" ^A\n<v>\n"))),
		matrix.Coord{X: 2, Y: 0},
	)
}

func PathToVecs(coords []matrix.Coord) []matrix.Vec {
	if len(coords) < 2 {
		return nil
	}
	vecs := make([]matrix.Vec, 0, len(coords)-1)
	for i := 0; i+1 < len(coords); i++ {
		vecs = append(vecs, coords[i+1].Sub(coords[i]))
	}
	return vecs
}

func TranslateVecs(vecs []matrix.Vec) []byte {
	b := bytes.NewBuffer(make([]byte, 0, len(vecs)))
	for _, v := range vecs {
		switch v {
		case matrix.Left:
			b.WriteByte('<')
		case matrix.Right:
			b.WriteByte('>')
		case matrix.Up:
			b.WriteByte('^')
		case matrix.Down:
			b.WriteByte('v')
		default:
			panic("unrecognized vec")
		}
	}
	return b.Bytes()
}

func (p *Pad) Type(seq []byte) []byte {
	buf := bytes.Buffer{}
	for _, b := range seq {
		buf.Write(TranslateVecs(PathToVecs(p.FindPath(p.Current(), b))))
		buf.WriteByte('A')
		p.SetCurrent(b)
	}
	return buf.Bytes()
}

// func (p *Pad) AllType(seq []byte) [][]byte {
// 	res := make([][]byte, 0)

// 	for _, b := range seq {
// 		p2 := p.Clone()
// 		for _, path := range p.FindAllPath(p.Current(), b) {
// 			buf := bytes.Buffer{}
// 			buf.Write(TranslateVecs(PathToVecs(path)))
// 			buf.WriteByte('A')
// 			p2.SetCurrent(b)
// 		}
// 	}

// 	return res
// }

func (p *Pad) TypeRec(seq []byte) iter.Seq[[]byte] {
	return func(yield func([]byte) bool) {
		if len(seq) == 0 {
			return
		}
		// res := make([][]byte, 0)
		allPath := p.FindAllPath(p.Current(), seq[0])
		p.SetCurrent(seq[0])

		buf := bytes.NewBuffer(nil)
		for _, path := range allPath {
			// stage.Println("path", path)
			prefix := TranslateVecs(PathToVecs(path))
			allSuffixes := p.Clone().TypeRec(seq[1:])
			iterated := false
			for suffix := range allSuffixes {
				iterated = true
				buf.Reset()
				buf.Write(prefix)
				buf.WriteByte('A')
				buf.Write(suffix)
				if !yield(buf.Bytes()) {
					return
				}
				// res = append(res, buf.Bytes())
			}
			if !iterated {
				buf.Reset()
				buf.Write(prefix)
				buf.WriteByte('A')
				if !yield(buf.Bytes()) {
					return
				}
				// res = append(res, buf.Bytes())
			}
			// if len(allSuffixes) == 0 {
			// 	buf := bytes.Buffer{}
			// 	buf.Write(prefix)
			// 	buf.WriteByte('A')
			// 	res = append(res, buf.Bytes())
			// 	continue
			// }
			// for _, suffix := range allSuffixes {
			// 	// stage.Println("suffix", suffix)
			// 	buf := bytes.Buffer{}
			// 	buf.Write(prefix)
			// 	buf.WriteByte('A')
			// 	buf.Write(suffix)
			// 	res = append(res, buf.Bytes())
			// }
		}
		// stage.Println("res", res)
		// return res
	}
}

// Numeric(NPad) <- Robot1(DPad) <- Robot2(Dpad) <-Robot3(Dpad)

var numRe = regexp.MustCompile(`\d+`)

func (p *Pad) AllTypeRec(seqs iter.Seq[[]byte]) iter.Seq[[]byte] {
	return func(yield func([]byte) bool) {
		for seq := range seqs {
			for input := range p.Clone().TypeRec(seq) {
				if !yield(input) {
					return
				}
			}
		}
	}
}

func CompareLen(a, b []byte) int {
	if len(a) < len(b) {
		return -1
	}
	if len(a) > len(b) {
		return 1
	}
	return 0
}

func First[T any](seq iter.Seq[T]) T {
	for v := range seq {
		return v
	}
	panic("empty seq")
}

func Stage1(input io.Reader) (any, error) {
	// npad := NewNPad()  // One numeric keypad (on a door) that a robot is using.
	// dpad1 := NewDPad() // Two directional keypads that robots are using.
	// dpad2 := NewDPad() // Two directional keypads that robots are using.
	// dpad3 := NewDPad() // One directional keypad that you are using.

	total := 0
	for line := range iterators.MustLinesBytes(input) {
		npad := NewNPad()  // One numeric keypad (on a door) that a robot is using.
		dpad1 := NewDPad() // Two directional keypads that robots are using.
		dpad2 := NewDPad() // Two directional keypads that robots are using.

		minLen := 1000000
		for input := range dpad2.AllTypeRec(dpad1.AllTypeRec(npad.TypeRec(line))) {
			// for input := range npad.TypeRec(line) {
			// stage.Println(string(input))
			if len(input) < minLen {
				minLen = len(input)
			}
		}
		num := utils.MustAtoi(string(numRe.Find(line)))
		stage.Println("num:", num, "len:", minLen, "->", num*minLen)
		total += num * minLen

		// allInputs := npad.AllTypeRec([][]byte{line, line, line})
		// a := slices.Values(slices.SortedFunc(npad.TypeRec(line), CompareLen))
		// b := slices.Values(slices.SortedFunc(dpad1.AllTypeRec(a), CompareLen))
		// c := slices.Values(slices.SortedFunc(dpad2.AllTypeRec(b), CompareLen))
		// dpad2Ks := First(dpad2.AllTypeRec(b))
		// allInputs := dpad2.AllTypeRec(dpad1.AllTypeRec(npad.TypeRec(line)))
		// minLen := 10000
		// checked := 0
		// go func() {
		// 	for range time.Tick(1 * time.Second) {
		// 		stage.Println(minLen, checked)
		// 	}
		// }()
		// for v := range c {
		// 	checked++
		// 	if len(v) < minLen {
		// 		minLen = len(v)
		// 	}
		// 	// stage.Println(string(v))
		// }
		// stage.Println(minLen)
		// stage.Println(string(dpad2.Type(dpad1.Type(npad.Type([]byte("029A"))))))
		// break
		// dpad1 := NewDPad() // Two directional keypads that robots are using.
		// dpad2 := NewDPad() // Two directional keypads that robots are using.

		// code := []byte(line)
		// stage.Println("code:", string(code))
		// npadKs := npad.Type(code)
		// stage.Println("numeric pad:", string(npadKs))
		// dpad1Ks := dpad1.Type(npadKs)
		// stage.Println("directional pad 1:", string(dpad1Ks))
		// dpad2Ks := dpad2.Type(dpad1Ks)
		// stage.Println("directional pad 2:", string(dpad2Ks))
		// // stage.Println(string(dpad2.Type(dpad1.Type(npad.Type([]byte("029A"))))))

		// num := utils.MustAtoi(numRe.FindString(line))
		// stage.Println("num:", num)
		// stage.Println("len:", len(dpad2Ks))
		// total += num * len(dpad2Ks)
		// stage.Println()
	}

	return total, nil
}

func Stage2(input io.Reader) (any, error) {
	return nil, stage.ErrUnimplemented
}
