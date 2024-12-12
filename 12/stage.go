package main

import (
	"io"

	"github.com/nlm/adventofcode2024/internal/matrix"
	"github.com/nlm/adventofcode2024/internal/stage"
	"github.com/nlm/adventofcode2024/internal/utils"
)

var Directions = []matrix.Vec{
	matrix.Up,
	matrix.Right,
	matrix.Down,
	matrix.Left,
}

// Visit returns area and perimeter
func Visit(m *matrix.Matrix[byte], c matrix.Coord, visited *matrix.Matrix[bool]) (int, int) {
	visited.SetAtCoord(c, true)
	area, perimeter := 1, 0
	for _, dir := range Directions {
		newCoord := c.Add(dir)
		if !SameValue(m, c, newCoord) {
			perimeter++
		}
		if SameValue(m, c, newCoord) && !visited.AtCoord(newCoord) {
			a, p := Visit(m, newCoord, visited)
			area += a
			perimeter += p
		}
	}
	return area, perimeter
}

func Stage1(input io.Reader) (any, error) {
	m := utils.Must(matrix.NewFromReader(input))
	visited := matrix.New[bool](m.Len.X, m.Len.Y)
	total := 0
	for c := range m.Coords() {
		if visited.AtCoord(c) {
			continue
		}
		area, perimeter := Visit(m, c, visited)
		stage.Println(string(m.AtCoord(c)), "| area", area, "perim", perimeter)
		total += area * perimeter
	}
	return total, nil
}

// Cells of same value at c1 and c2
func SameValue(m *matrix.Matrix[byte], c1, c2 matrix.Coord) bool {
	return m.InCoord(c1) && m.InCoord(c2) && m.AtCoord(c1) == m.AtCoord(c2)
}

var NextDir = map[matrix.Vec]matrix.Vec{
	matrix.Up:    matrix.Right,
	matrix.Right: matrix.Down,
	matrix.Down:  matrix.Left,
	matrix.Left:  matrix.Up,
}

func CountCorners(m *matrix.Matrix[byte], c matrix.Coord) int {
	corners := 0
	for _, dir := range Directions {
		// Convex
		if !SameValue(m, c, c.Add(dir)) && !SameValue(m, c, c.Add(NextDir[dir])) {
			corners++
		}
		// Concave
		if SameValue(m, c, c.Add(dir)) && SameValue(m, c, c.Add(NextDir[dir])) && !SameValue(m, c, c.Add(dir).Add(NextDir[dir])) {
			corners++
		}
	}
	return corners
}

// Visit returns area and sides
func VisitWithCorners(m *matrix.Matrix[byte], c matrix.Coord, visited *matrix.Matrix[bool]) (int, int) {
	visited.SetAtCoord(c, true)
	area, sides := 1, CountCorners(m, c)
	stage.Println(c, "has", area, "area and", sides, "corners")
	for _, dir := range Directions {
		newCoord := c.Add(dir)
		if SameValue(m, c, newCoord) && !visited.AtCoord(newCoord) {
			a, s := VisitWithCorners(m, newCoord, visited)
			area += a
			sides += s
		}
	}
	return area, sides
}

func Stage2(input io.Reader) (any, error) {
	m := utils.Must(matrix.NewFromReader(input))
	visited := matrix.New[bool](m.Len.X, m.Len.Y)
	total := 0
	stage.Println(matrix.SMatrix(m))
	for c := range m.Coords() {
		if visited.AtCoord(c) {
			continue
		}
		area, sides := VisitWithCorners(m, c, visited)
		stage.Println(string(m.AtCoord(c)), "| area", area, "sides", sides)
		total += area * sides
	}
	return total, nil
}
