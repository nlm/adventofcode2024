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
		curr := orig.Add(dir.Mul(i))
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
	for coord := range m.IterCoords() {
		if m.AtCoord(coord) != 'X' {
			continue
		}
		stage.Println("Search", coord)
		for _, vec := range []matrix.Vec{
			// Clock 12
			matrix.Up,
			// Clock 1.5
			matrix.UpRight,
			// Clock 3
			matrix.Right,
			// Clock 4.5
			matrix.DownRight,
			// CLock 6
			matrix.Down,
			// Clock 7.5
			matrix.DownLeft,
			// Clock 9
			matrix.Left,
			// Clock 10.5
			matrix.UpLeft,
		} {
			total += SearchXmas(m, coord, vec)
		}
	}
	return total, nil
}

func SearchCrossMass(m *matrix.Matrix[byte], orig matrix.Coord) int {
	vecs := []matrix.Vec{
		matrix.UpLeft,
		matrix.UpRight,
		matrix.DownLeft,
		matrix.DownRight,
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
	if m.AtCoord(orig.Add(matrix.UpLeft)) == m.AtCoord(orig.Add(matrix.DownRight)) {
		return 0
	}
	return 1
}

func Stage2(input io.Reader) (any, error) {
	total := 0
	m := utils.Must(matrix.NewFromReader(input))
	for coord := range m.IterCoords() {
		if m.AtCoord(coord) != 'A' {
			continue
		}
		stage.Println("Search", coord)
		total += SearchCrossMass(m, coord)
	}
	return total, nil
}
