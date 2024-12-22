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
	Matrix *matrix.Matrix[byte]
	// Origin        matrix.Coord
	// Pointer       matrix.Coord
	Graph         *simple.UndirectedGraph
	KeyIndex      map[byte]matrix.Coord
	ShortestPaths map[[2]byte][][]byte
}

// func (p *Pad) Reset() {
// 	p.Pointer = p.Origin
// }

// func (p *Pad) Current() byte {
// 	return p.Matrix.AtCoord(p.Pointer)
// }

// func (p Pad) Clone() *Pad {
// 	return &p
// }

// func (p *Pad) SetCurrent(b byte) {
// 	p.Pointer = p.KeyIndex[b]
// }

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

func (p *Pad) GetTypedPaths(from, to byte) [][]byte {
	shortestPaths, ok := p.ShortestPaths[[2]byte{from, to}]
	if !ok {
		panic("invalid pair for path")
	}
	return shortestPaths
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
		Matrix: m,
		// Pointer:  origin,
		// Origin:   origin,
		KeyIndex: make(map[byte]matrix.Coord, len(m.Data)),
	}
	// Precompute graph
	g := simple.NewUndirectedGraph()
	for c := range p.Matrix.Coords() {
		// skip empty slots
		if m.AtCoord(c) == ' ' {
			continue
		}
		// Index
		p.KeyIndex[m.AtCoord(c)] = c
		// Graph
		node1, isNew := g.NodeWithID(maze.CoordToId(p.Matrix, c))
		if isNew {
			g.AddNode(node1)
		}
		for nc := range ManhattanNeighbors(p.Matrix, c) {
			// skip empty slots
			if m.AtCoord(nc) == ' ' {
				continue
			}
			node2, isNew := g.NodeWithID(maze.CoordToId(p.Matrix, nc))
			if isNew {
				g.AddNode(node2)
			}
			g.SetEdge(g.NewEdge(node1, node2))
		}
	}
	p.Graph = g
	// Precompute shortests paths
	p.ShortestPaths = make(map[[2]byte][][]byte, len(m.Data))
	data := iterators.FilterSlice(m.Data, func(b byte) bool { return b != ' ' })
	for _, d1 := range data {
		for _, d2 := range data {
			key := [2]byte{d1, d2}
			for _, spath := range p.FindAllPath(d1, d2) {
				p.ShortestPaths[key] = append(p.ShortestPaths[key], TranslateVecs(PathToVecs(spath)))
			}
		}
	}
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

// func (p *Pad) Type(seq []byte) []byte {
// 	buf := bytes.Buffer{}
// 	for _, b := range seq {
// 		buf.Write(TranslateVecs(PathToVecs(p.FindPath(p.Current(), b))))
// 		buf.WriteByte('A')
// 		p.SetCurrent(b)
// 	}
// 	return buf.Bytes()
// }

// Not used
func (p *Pad) TypeRec(from byte, seq []byte) iter.Seq[[]byte] {
	return func(yield func([]byte) bool) {
		if len(seq) == 0 {
			return
		}
		allPath := p.FindAllPath(from, seq[0])

		buf := bytes.NewBuffer(nil)
		for _, path := range allPath {
			// stage.Println("path", path)
			prefix := TranslateVecs(PathToVecs(path))
			allSuffixes := p.TypeRec(seq[0], seq[1:])
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
			}
			if !iterated {
				buf.Reset()
				buf.Write(prefix)
				buf.WriteByte('A')
				if !yield(buf.Bytes()) {
					return
				}
			}
		}
	}
}

// Numeric(NPad) <- Robot1(DPad) <- Robot2(Dpad) <-Robot3(Dpad)

var numRe = regexp.MustCompile(`\d+`)

// func (p *Pad) AllTypeRec(seqs iter.Seq[[]byte]) iter.Seq[[]byte] {
// 	return func(yield func([]byte) bool) {
// 		for seq := range seqs {
// 			for input := range p.Clone().TypeRec(seq) {
// 				if !yield(input) {
// 					return
// 				}
// 			}
// 		}
// 	}
// }

// func CompareLen(a, b []byte) int {
// 	if len(a) < len(b) {
// 		return -1
// 	}
// 	if len(a) > len(b) {
// 		return 1
// 	}
// 	return 0
// }

// func First[T any](seq iter.Seq[T]) T {
// 	for v := range seq {
// 		return v
// 	}
// 	panic("empty seq")
// }

// 02
// *
func (p *Pad) BuildSeq(fromKey byte, toKeys []byte, currpath []byte) iter.Seq[[]byte] {
	// stage.Println("buildseq", "from", string(fromKey), "to", string(toKeys), "curr", string(currpath))
	return func(yield func([]byte) bool) {
		if len(toKeys) == 0 {
			yield(currpath)
			return
		}
		nextKey := toKeys[0]
		// stage.Println("BSEQ2:", string(p.Current()), string(nextKey))
		// for _, prefix := range p.ShortestPaths[[2]byte{p.Current(), nextKey}] {
		for _, prefix := range p.GetTypedPaths(fromKey, nextKey) {
			// stage.Println("bseq3:", string(currpath), "prefix:", string(prefix))
			nextPath := append([]byte(nil), currpath...)
			nextPath = append(nextPath, prefix...)
			nextPath = append(nextPath, 'A')
			for suffix := range p.BuildSeq(nextKey, toKeys[1:], nextPath) {
				yield(suffix)
			}
			// here
		}
	}
}

type CacheKey struct {
	String string
	Depth  int
}

func (p *Pad) FindShortestSeq(keys []byte, depth int, cache map[CacheKey]int) int {
	if depth == 0 {
		return len(keys)
	}
	if v, ok := cache[CacheKey{string(keys), depth}]; ok {
		return v
	}
	total := 0
	for _, subkeys := range strings.SplitAfter(string(keys), "A") {
		minSeqLen := MaxInt
		for seq := range p.BuildSeq('A', []byte(subkeys), nil) {
			v := p.FindShortestSeq(seq, depth-1, cache)
			if v < minSeqLen {
				minSeqLen = v
			}
		}
		total += minSeqLen
	}
	cache[CacheKey{string(keys), depth}] = total
	return total
}

const MaxUint = ^uint(0)
const MinUint = 0
const MaxInt = int(MaxUint >> 1)
const MinInt = -MaxInt - 1

func Solve(numericInputs iter.Seq[[]byte], maxDepth int) int64 {
	npad := NewNPad()
	dpad := NewDPad()

	total := int64(0)
	for numericInput := range numericInputs {
		lowest := MaxInt
		for seq := range npad.BuildSeq('A', numericInput, nil) {
			v := dpad.FindShortestSeq(seq, maxDepth, make(map[CacheKey]int))
			if v < lowest {
				lowest = v
			}
		}
		num := utils.MustAtoi(string(numRe.Find(numericInput)))
		total += int64(lowest) * int64(num)
	}
	return total
}

func Stage1(input io.Reader) (any, error) {
	return Solve(iterators.MustLinesBytes(input), 2), nil

}

func Stage2(input io.Reader) (any, error) {
	return Solve(iterators.MustLinesBytes(input), 25), nil
}
