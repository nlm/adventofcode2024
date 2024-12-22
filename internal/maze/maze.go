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

// PathFinder helps find paths within a Matrix[byte].
type PathFinder struct {
	m *matrix.Matrix[byte]
	g *simple.WeightedDirectedGraph
}

// CoordToId translates coordinates inside the matrix to a unique id.
func CoordToId[T comparable](m *matrix.Matrix[T], c matrix.Coord) int64 {
	return int64(c.Y*m.Size.X + c.X)
}

// IdToCoord translates an previously generated id back to coordinates inside the matrix.
func IdToCoord[T comparable](m *matrix.Matrix[T], id int64) matrix.Coord {
	return matrix.Coord{X: int(id % int64(m.Size.X)), Y: int(id / int64(m.Size.X))}
}

// NewSimplePathFinder creates a PathFinder
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

// AddSpecialNode adds a new node in the graph at coordinate c.
// It will search all the reachable neightbors from this point and create path FROM it to these.
// If invert is true, it will create path from all reachable neighbors TO this point.
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

// FindDijkstra uses the Dijkstra algorithm to find a shortest path from a coordinate to another.
func (pf *PathFinder) FindDijkstra(from, to matrix.Coord) ([]matrix.Coord, int64) {
	paths := path.DijkstraFrom(pf.g.Node(CoordToId(pf.m, from)), pf.g)
	sp, w := paths.To(CoordToId(pf.m, to))
	spres := iterators.MapSlice(sp, func(node graph.Node) matrix.Coord { return IdToCoord(pf.m, node.ID()) })
	return spres, int64(w)
}

// FindDijkstra uses the Dijkstra algorithm to find all shortest paths from a coordinate to another.
func (pf *PathFinder) FindAllDijkstra(from, to matrix.Coord) ([][]matrix.Coord, int64) {
	paths := path.DijkstraAllFrom(pf.g.Node(CoordToId(pf.m, from)), pf.g)
	sp, w := paths.AllTo(CoordToId(pf.m, to))
	spres := iterators.MapSlice(sp, func(nodes []graph.Node) []matrix.Coord {
		return iterators.MapSlice(nodes, func(node graph.Node) matrix.Coord {
			return IdToCoord(pf.m, node.ID())
		})
	})
	return spres, int64(w)
}

// Graph returns the underlying graph contained in the PathFinder.
func (pf *PathFinder) Graph() *simple.WeightedDirectedGraph {
	return pf.g
}
