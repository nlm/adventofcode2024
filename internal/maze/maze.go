package maze

import (
	"github.com/nlm/adventofcode2024/internal/iterators"
	"github.com/nlm/adventofcode2024/internal/matrix"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/path"
	"gonum.org/v1/gonum/graph/simple"
)

const (
	SymbolWall  = '#'
	SymbolEmpty = '.'
)

type PathFinder struct {
	m *matrix.Matrix[byte]
	g *simple.WeightedDirectedGraph
}

func CoordToId[T comparable](m *matrix.Matrix[T], c matrix.Coord) int64 {
	// return int64(c.Y*m.Size.X + c.X)
	return int64(1000*c.X + c.Y)
}

func IdToCoord[T comparable](m *matrix.Matrix[T], id int64) matrix.Coord {
	// return matrix.Coord{X: int(id % int64(m.Size.X)), Y: int(id / int64(m.Size.X))}
	return matrix.Coord{
		X: int(id / 1000),
		Y: int(id % 1000),
	}
}

func NewSimplePathFinder(m *matrix.Matrix[byte]) *PathFinder {
	g := simple.NewWeightedDirectedGraph(0, 0)
	for c := range m.Coords() {
		currNode, isNew := g.NodeWithID(CoordToId(m, c))
		if isNew {
			g.AddNode(currNode)
		}
		for _, dir := range []matrix.Vec{matrix.Up, matrix.Down, matrix.Left, matrix.Right} {
			next := c.Add(dir)
			if !m.InCoord(next) || m.AtCoord(next) != SymbolEmpty {
				continue
			}
			nextNode, isNew := g.NodeWithID(CoordToId(m, next))
			if isNew {
				g.AddNode(nextNode)
			}
			g.SetWeightedEdge(g.NewWeightedEdge(currNode, nextNode, 1))
		}
	}
	return &PathFinder{
		m: m,
		g: g,
	}
}

func (pf *PathFinder) AddSpecialNode(m *matrix.Matrix[byte], c matrix.Coord, invert bool) {
	currNode, isNew := pf.g.NodeWithID(CoordToId(m, c))
	if isNew {
		pf.g.AddNode(currNode)
	}
	for _, dir := range []matrix.Vec{matrix.Up, matrix.Down, matrix.Left, matrix.Right} {
		next := c.Add(dir)
		if !m.InCoord(next) || m.AtCoord(next) != SymbolEmpty {
			continue
		}
		nextNode, isNew := pf.g.NodeWithID(CoordToId(m, next))
		if isNew {
			pf.g.AddNode(nextNode)
		}
		if invert {
			currNode, nextNode = nextNode, currNode
		}
		pf.g.SetWeightedEdge(pf.g.NewWeightedEdge(currNode, nextNode, 1))
	}
}

func (pf *PathFinder) FindDijkstra(from, to matrix.Coord) ([]matrix.Coord, int64) {
	paths := path.DijkstraFrom(pf.g.Node(CoordToId(pf.m, from)), pf.g)
	sp, w := paths.To(CoordToId(pf.m, to))
	spres := iterators.MapSlice(sp, func(node graph.Node) matrix.Coord { return IdToCoord(pf.m, node.ID()) })
	return spres, int64(w)
}

func (pf *PathFinder) Graph() *simple.WeightedDirectedGraph {
	return pf.g
}
