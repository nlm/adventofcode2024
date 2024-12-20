package main

import (
	"fmt"
	"io"

	"github.com/nlm/adventofcode2024/internal/math"
	"github.com/nlm/adventofcode2024/internal/matrix"
	"github.com/nlm/adventofcode2024/internal/maze"
	"github.com/nlm/adventofcode2024/internal/stage"
	"github.com/nlm/adventofcode2024/internal/utils"
)

func Stage(input io.Reader, minPicoSaved, shortcutLength int) (any, error) {
	m := utils.Must(matrix.NewFromReader(input))
	stage.Println(matrix.SMatrix(m))
	start, ok := m.Find('S')
	if !ok {
		return nil, fmt.Errorf("start not found")
	}
	end, ok := m.Find('E')
	if !ok {
		return nil, fmt.Errorf("end not found")
	}

	// find shortest path
	pf := maze.NewSimplePathFinder(m)
	pf.AddSpecialNode(m, end, true)
	spath, _ := pf.FindDijkstra(start, end)

	// index coordinates in shortest path for fast lookup
	pathIdx := make(map[matrix.Coord]int)
	for i, v := range spath {
		pathIdx[v] = i
	}

	// find shortcuts
	total := 0
	for i, c1 := range spath {
		for _, c2 := range spath[i+1:] {
			if IsShortcut(m, shortcutLength, c1, c2) {
				timeSaved := GetTimeSaved(pathIdx, c1, c2)
				if timeSaved >= minPicoSaved {
					stage.Println("shortcut", c1, "->", c2, "cost", ManhattanDistance(c1, c2), "saved", timeSaved)
					total++
				}
			}
		}
	}
	return total, nil
}

func GetTimeSaved(pathIdx map[matrix.Coord]int, c1, c2 matrix.Coord) int {
	i1, ok := pathIdx[c1]
	if !ok {
		return -1
	}
	i2, ok := pathIdx[c2]
	if !ok {
		return -1
	}
	return max(i2-i1-ManhattanDistance(c1, c2), -1)
}

func IsShortcut(m *matrix.Matrix[byte], length int, c1, c2 matrix.Coord) bool {
	dist := ManhattanDistance(c1, c2)
	return dist > 0 && dist <= length
}

func ManhattanDistance(c1, c2 matrix.Coord) int {
	return math.Abs(c2.Y-c1.Y) + math.Abs(c2.X-c1.X)
}

func Stage1(input io.Reader) (any, error) {
	// Stage(example.txt, 1, 2) -> 44
	// Stage(input.txt, 100, 2) -> 1351
	return Stage(input, 100, 2)
}

func Stage2(input io.Reader) (any, error) {
	// Stage(example.txt, 50, 20) -> 44
	// Stage(input.txt, 100, 20) -> 966130
	return Stage(input, 100, 20)
}
