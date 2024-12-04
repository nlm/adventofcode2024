package main

import (
	"io"

	"github.com/nlm/adventofcode2024/internal/matrix"
	"github.com/nlm/adventofcode2024/internal/stage"
	"github.com/nlm/adventofcode2024/internal/utils"
)

func SearchXmas(m *matrix.Matrix[byte], orig matrix.Coord, dir matrix.Vec) int {
	word := []byte{'X', 'M', 'A', 'S'}
	for i := range len(word) {
		curr := orig.Add(matrix.Vec{X: dir.X * i, Y: dir.Y * i})
		if !m.InCoord(curr) {
			// stage.Println("out of bounds")
			return 0
		}
		if m.AtCoord(curr) != word[i] {
			return 0
		}
	}
	stage.Println("found", orig, "->", dir)
	return 1
}

func Stage1(input io.Reader) (any, error) {
	total := 0
	m := utils.Must(matrix.NewFromReader(input))
	for y := 0; y < m.Len.Y; y++ {
		for x := 0; x < m.Len.X; x++ {
			if m.At(x, y) != 'X' {
				continue
			}
			stage.Println("Search", x, y)
			coord := matrix.Coord{X: x, Y: y}
			// Clock 12
			total += SearchXmas(m, coord, matrix.Vec{X: 0, Y: 1})
			// Clock 1.5
			total += SearchXmas(m, coord, matrix.Vec{X: 1, Y: 1})
			// Clock 3
			total += SearchXmas(m, coord, matrix.Vec{X: 1, Y: 0})
			// Clock 4.5
			total += SearchXmas(m, coord, matrix.Vec{X: 1, Y: -1})
			// CLock 6
			total += SearchXmas(m, coord, matrix.Vec{X: 0, Y: -1})
			// Clock 7.5
			total += SearchXmas(m, coord, matrix.Vec{X: -1, Y: -1})
			// Clock 9
			total += SearchXmas(m, coord, matrix.Vec{X: -1, Y: 0})
			// Clock 10.5
			total += SearchXmas(m, coord, matrix.Vec{X: -1, Y: 1})
		}
	}
	return total, nil
}

func SearchCrossMass(m *matrix.Matrix[byte], orig matrix.Coord) int {
	vecs := []matrix.Vec{
		{X: -1, Y: -1},
		{X: -1, Y: 1},
		{X: 1, Y: 1},
		{X: 1, Y: -1},
	}
	xs := map[byte]int{'M': 0, 'S': 0}
	for _, v := range vecs {
		curr := orig.Add(v)
		if !m.InCoord(curr) {
			return 0
		}
		xs[m.AtCoord(curr)]++
	}
	// MAS are there
	if xs['M'] != 2 || xs['S'] != 2 {
		return 0
	}
	// MAM / SAS
	if m.AtCoord(orig.Add(matrix.Vec{X: -1, Y: -1})) == m.AtCoord(orig.Add(matrix.Vec{X: 1, Y: 1})) {
		return 0
	}
	return 1
}

func Stage2(input io.Reader) (any, error) {
	total := 0
	m := utils.Must(matrix.NewFromReader(input))
	for y := 0; y < m.Len.Y; y++ {
		for x := 0; x < m.Len.X; x++ {
			if m.At(x, y) != 'A' {
				continue
			}
			stage.Println("Search", x, y)
			coord := matrix.Coord{X: x, Y: y}
			total += SearchCrossMass(m, coord)
		}
	}
	return total, nil
}
