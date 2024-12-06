package main

import (
	"fmt"
	"io"
	"slices"

	"github.com/nlm/adventofcode2024/internal/matrix"
	"github.com/nlm/adventofcode2024/internal/utils"
)

var NextDirection = map[matrix.Vec]matrix.Vec{
	matrix.Up:    matrix.Right,
	matrix.Right: matrix.Down,
	matrix.Down:  matrix.Left,
	matrix.Left:  matrix.Up,
}

func RunGuard(m *matrix.Matrix[byte], orig matrix.Coord, dir matrix.Vec) error {
	curr := orig
	for {
		m.SetAtCoord(curr, 'X')
		// fmt.Println(matrix.SMatrix(m))
		next := curr.Add(dir)
		if !m.InCoord(next) {
			return fmt.Errorf("out of bounds")
		}
		if m.AtCoord(next) == '#' {
			dir = NextDirection[dir]
			continue
		}
		curr = next
	}
}

func Stage1(input io.Reader) (any, error) {
	m := utils.Must(matrix.NewFromReader(input))
	origin, _ := m.Find('^')
	dir := matrix.Up
	_ = RunGuard(m, origin, dir)
	total := m.Count('X')
	return total, nil
}

func DetectLoop(m *matrix.Matrix[byte], orig matrix.Coord, dir matrix.Vec) bool {
	visits := make(map[matrix.Coord][]matrix.Vec)
	curr := orig
	for {
		m.SetAtCoord(curr, 'X')
		// fmt.Println(matrix.SMatrix(m))
		next := curr.Add(dir)
		if !m.InCoord(next) {
			return false
		}
		if m.AtCoord(next) == '#' || m.AtCoord(next) == 'O' {
			dir = NextDirection[dir]
			if slices.Contains(visits[curr], dir) {
				return true
			}
			visits[curr] = append(visits[curr], dir)
			continue
		}
		curr = next
	}
}

func Stage2(input io.Reader) (any, error) {
	m := utils.Must(matrix.NewFromReader(input))
	origin, _ := m.Find('^')
	dir := matrix.Up
	total := 0
	for coord := range m.IterCoords() {
		if m.AtCoord(coord) == '#' {
			// can't place obstacle on a wall
			continue
		}
		m2 := m.Clone()
		// Set obstacle
		m2.SetAtCoord(coord, 'O')
		if DetectLoop(m2, origin, dir) {
			total++
		}
	}
	return total, nil
}
