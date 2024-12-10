package main

import (
	"io"
	"iter"

	"github.com/nlm/adventofcode2024/internal/iterators"
	"github.com/nlm/adventofcode2024/internal/matrix"
	"github.com/nlm/adventofcode2024/internal/sets"
	"github.com/nlm/adventofcode2024/internal/stage"
	"github.com/nlm/adventofcode2024/internal/utils"
)

var IntToByte = map[int]byte{
	0: '0',
	1: '1',
	2: '2',
	3: '3',
	4: '4',
	5: '5',
	6: '6',
	7: '7',
	8: '8',
	9: '9',
}

func FindOrigins(m *matrix.Matrix[byte]) iter.Seq[matrix.Coord] {
	return iterators.Filter(m.Coords(), func(c matrix.Coord) bool {
		return m.AtCoord(c) == '0'
	})
}

var Directions = []matrix.Vec{
	matrix.Up,
	matrix.Right,
	matrix.Down,
	matrix.Left,
}

func SearchTrails(m *matrix.Matrix[byte], c matrix.Coord, v int) []matrix.Coord {
	// stage.Println("considering", c)
	// check we're in the matrix
	if !m.InCoord(c) {
		// stage.Println(" not in coord")
		return nil
	}
	// check current value is ok
	if m.AtCoord(c) != IntToByte[v] {
		// stage.Printf("% s != %s\n", string(m.AtCoord(c)), string(IntToByte[v]))
		return nil
	}
	// stage.Println(" value ok ", c, "->", v)
	// check end of trail
	if m.AtCoord(c) == '9' {
		// stage.Println(" found 9 at", c)
		return []matrix.Coord{c}
	}
	var total []matrix.Coord
	for _, dir := range Directions {
		next := c.Add(dir)
		// stage.Println(" next is", "from", v, "at", c, "to", v+1, "at", next, "via", dir)
		total = append(total, SearchTrails(m, next, v+1)...)
	}
	return total
}

func Stage1(input io.Reader) (any, error) {
	m := utils.Must(matrix.NewFromReader(input))
	stage.Println(matrix.SMatrix(m))
	total := 0
	for origin := range FindOrigins(m) {
		stage.Println("origin at", origin)
		trails := SearchTrails(m, origin, 0)
		nines := sets.Append(nil, trails...)
		stage.Println("->", len(nines))
		total += len(nines)
	}
	return total, nil
}

func Stage2(input io.Reader) (any, error) {
	m := utils.Must(matrix.NewFromReader(input))
	stage.Println(matrix.SMatrix(m))
	total := 0
	for origin := range FindOrigins(m) {
		stage.Println("origin at", origin)
		trails := SearchTrails(m, origin, 0)
		stage.Println("->", len(trails))
		total += len(trails)
	}
	return total, nil
}
