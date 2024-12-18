package main

import (
	"fmt"
	"io"
	"strings"

	"github.com/nlm/adventofcode2024/internal/iterators"
	"github.com/nlm/adventofcode2024/internal/matrix"
	"github.com/nlm/adventofcode2024/internal/stage"
	"github.com/nlm/adventofcode2024/internal/utils"
	"gonum.org/v1/gonum/graph/path"
	"gonum.org/v1/gonum/graph/simple"
)

func ParseInput(input io.Reader, maxX, maxY int) (*matrix.Matrix[byte], []matrix.Coord) {
	coords := make([]matrix.Coord, 0)
	for line := range iterators.MustLines(input) {
		fields := strings.Split(line, ",")
		coords = append(coords, matrix.Coord{X: utils.MustAtoi(fields[0]), Y: utils.MustAtoi(fields[1])})
	}
	m := matrix.New[byte](maxX+1, maxY+1)
	m.Fill('.')
	return m, coords
}

func Corrupt(m *matrix.Matrix[byte], coords []matrix.Coord, nbytes int) {
	for i := range nbytes {
		m.SetAtCoord(coords[i], '#')
	}
}

func CoordToId(c matrix.Coord) int64 {
	return int64(c.X*1000 + c.Y)
}

func IdToCoord(id int64) matrix.Coord {
	return matrix.Coord{X: int(id / 1000), Y: int(id % 1000)}
}

func Stage1(input io.Reader) (any, error) {
	var example = false
	var maxX, maxY = 70, 70
	var nCorr = 1024

	if example {
		maxX, maxY = 6, 6
		nCorr = 12
	}

	m, coords := ParseInput(input, maxX, maxY)
	Corrupt(m, coords, nCorr)
	stage.Println(matrix.SMatrix(m))
	stage.Println(coords)

	start := matrix.Coord{X: 0, Y: 0}
	end := matrix.Coord{X: maxX, Y: maxY}
	return FindPath(m, start, end), nil
}

func FindPath(m *matrix.Matrix[byte], from, to matrix.Coord) int {
	g := simple.NewWeightedDirectedGraph(0, 0)
	for c := range m.Coords() {
		currNode, isNew := g.NodeWithID(CoordToId(c))
		if isNew {
			g.AddNode(currNode)
		}
		for _, dir := range []matrix.Vec{matrix.Up, matrix.Down, matrix.Left, matrix.Right} {
			next := c.Add(dir)
			if !m.InCoord(next) || m.AtCoord(next) != '.' {
				continue
			}
			nextNode, isNew := g.NodeWithID(CoordToId(next))
			if isNew {
				g.AddNode(nextNode)
			}
			g.SetWeightedEdge(g.NewWeightedEdge(currNode, nextNode, 1))
		}
	}

	paths := path.DijkstraFrom(g.Node(CoordToId(from)), g)
	if sp, weight := paths.To(CoordToId(to)); sp != nil {
		if stage.Verbose() {
			stage.Println("found:")
			vm := m.Clone()
			for _, n := range sp {
				c := IdToCoord(n.ID())
				vm.SetAtCoord(c, 'O')
				stage.Printf("(%d, %d) -> ", c.X, c.Y)
			}
			stage.Printf("end\ntotal : %.2f\n", weight)
			stage.Println(matrix.SMatrix(vm))
		}
		return int(weight)
	} else {
		stage.Println("not found")
		return -1
	}
}

func Stage2(input io.Reader) (any, error) {
	example := false

	var maxX, maxY = 70, 70
	var nCorr = 1024
	if example {
		maxX, maxY = 6, 6
		nCorr = 12
	}

	m, coords := ParseInput(input, maxX, maxY)
	stage.Println(matrix.SMatrix(m))
	stage.Println(coords)

	for nC := nCorr; nC < len(coords); nC++ {
		m2 := m.Clone()
		Corrupt(m2, coords, nC)

		start := matrix.Coord{X: 0, Y: 0}
		end := matrix.Coord{X: maxX, Y: maxY}
		w := FindPath(m2, start, end)
		stage.Println("->", nC, coords[nC-1])
		if w < 0 {
			c := coords[nC-1]
			return fmt.Sprintf("%d,%d", c.X, c.Y), nil
		}
	}

	return 0, nil
}
