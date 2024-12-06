package main

import (
	"io"

	"github.com/nlm/adventofcode2024/internal/matrix"
	"github.com/nlm/adventofcode2024/internal/stage"
	"github.com/nlm/adventofcode2024/internal/utils"
)

func SearchXmas(m *matrix.Matrix[byte], orig matrix.Coord, dir matrix.Vec) bool {
	word := []byte{'X', 'M', 'A', 'S'}
	for i := range len(word) {
		curr := orig.Add(dir.Mul(i))
		if !m.InCoord(curr) {
			// stage.Println("out of bounds")
			return false
		}
		if m.AtCoord(curr) != word[i] {
			return false
		}
	}
	stage.Println("found", orig, "->", dir)
	return true
}

func Stage1(input io.Reader) (any, error) {
	total := 0
	m := utils.Must(matrix.NewFromReader(input))
	for coord := range m.IterCoords() {
		// XMAS always starts with an X
		if m.AtCoord(coord) != 'X' {
			continue
		}
		stage.Println("Search", coord)
		// Lookup XMAS in every direction
		for _, vec := range []matrix.Vec{
			// Clock 12:00
			matrix.Up,
			// Clock 1:30
			matrix.UpRight,
			// Clock 3:00
			matrix.Right,
			// Clock 4:30
			matrix.DownRight,
			// Clock 6:00
			matrix.Down,
			// Clock 7:30
			matrix.DownLeft,
			// Clock 9:00
			matrix.Left,
			// Clock 10:30
			matrix.UpLeft,
		} {
			if SearchXmas(m, coord, vec) {
				total++
			}
		}
	}
	return total, nil
}

func SearchCrossMass(m *matrix.Matrix[byte], orig matrix.Coord) bool {
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
			return false
		}
		xs[m.AtCoord(curr)]++
	}
	// check that the two M & S are there
	if xs['M'] != 2 || xs['S'] != 2 {
		return false
	}
	// check that it's not MAM / SAS
	if m.AtCoord(orig.Add(matrix.UpLeft)) == m.AtCoord(orig.Add(matrix.DownRight)) {
		return false
	}
	return true
}

func Stage2(input io.Reader) (any, error) {
	total := 0
	m := utils.Must(matrix.NewFromReader(input))
	for coord := range m.IterCoords() {
		// We want to look around A
		if m.AtCoord(coord) != 'A' {
			continue
		}
		stage.Println("Search", coord)
		// Look around this A for matches
		if SearchCrossMass(m, coord) {
			total++
		}
	}
	return total, nil
}
