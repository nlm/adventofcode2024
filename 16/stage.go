package main

import (
	"fmt"
	"io"
	"sort"
	"time"

	"github.com/nlm/adventofcode2024/internal/matrix"
	"github.com/nlm/adventofcode2024/internal/sets"
	"github.com/nlm/adventofcode2024/internal/stage"
	"github.com/nlm/adventofcode2024/internal/utils"
)

var Dirs = []matrix.Vec{
	matrix.Up,
	matrix.Right,
	matrix.Down,
	matrix.Left,
}

var Opposite = map[matrix.Vec]matrix.Vec{
	matrix.Up:    matrix.Down,
	matrix.Down:  matrix.Up,
	matrix.Left:  matrix.Right,
	matrix.Right: matrix.Left,
}

func CalcInc(dir matrix.Vec, lastDir matrix.Vec) int {
	if dir == lastDir {
		return 1
	}
	if dir == Opposite[lastDir] {
		panic("opposite")
	}
	return 1001
}

type WeightedVec struct {
	Vec  matrix.Vec
	Cost int
	Dist int
}

func Distance(a, b matrix.Coord) int {
	x := a.X - b.X
	y := a.Y - b.Y
	if x < 0 {
		x = -x
	}
	if y < 0 {
		y = -y
	}
	return x + y
}

type CoordVec struct {
	Coord matrix.Coord
	Last  matrix.Vec
	Next  matrix.Vec
}

var Ops = 0
var BestScore = 0
var BestCosts = make(map[CoordVec]int)
var BestScores = make(map[matrix.Coord]int)

func FindPath(m *matrix.Matrix[byte], visited sets.Set[matrix.Coord], lastDir matrix.Vec, score int, start, end matrix.Coord) (int, bool) {
	Ops++
	curr := start
	visited.Add(curr)
	bestScore := 0
	if stage.Verbose() {
		// uViz
		vm := m.Clone()
		for c := range visited {
			vm.SetAtCoord(c, 'X')
		}
		// for c := range BestCosts {
		// 	vm.SetAtCoord(c, 'i')
		// }
		time.Sleep(100 * time.Millisecond)
		stage.Println()
		stage.Println("curr", curr)
		stage.Println("score", score, "/", BestScore)
		// stage.Println("best", BestScores)
		stage.Println(matrix.SMatrix(vm))
	}
	potentialPaths := make([]WeightedVec, 0, 4)
	for _, dir := range Dirs {
		next := curr.Add(dir)
		switch m.AtCoord(next) {
		case '#':
			continue
		case '.':
			if !visited.Contains(next) {
				potentialPaths = append(potentialPaths, WeightedVec{dir, CalcInc(dir, lastDir), Distance(next, end)})
			}
		case 'E':
			inc := CalcInc(dir, lastDir)
			if bestScore == 0 || score+inc < BestScore {
				BestScore = score + inc
				BestScores[next] = BestScore
			}
			return score + inc, true
		}
	}
	// priority sort best paths
	sort.Slice(potentialPaths, func(i, j int) bool {
		if potentialPaths[i].Cost != potentialPaths[j].Cost {
			return potentialPaths[i].Cost < potentialPaths[j].Cost
		}
		return potentialPaths[i].Dist < potentialPaths[j].Dist
	})
	stage.Println("pot:", potentialPaths)
	var found bool
	for _, p := range potentialPaths {
		inc := CalcInc(p.Vec, lastDir)
		next := curr.Add(p.Vec)
		// Path has previous lowest score
		if BestScore > 0 && score+inc > BestScore {
			continue
		}
		cvec := CoordVec{next, lastDir, p.Vec}
		if BestCosts[cvec] != 0 && score+inc > BestCosts[cvec] {
			// panic("optim")
			continue
		}
		BestCosts[cvec] = score + inc
		newVisited := visited.Clone()
		newScore, ok := FindPath(m, newVisited, p.Vec, score+inc, next, end)
		if !ok {
			continue
		}
		if BestScores[next] == 0 || BestScores[next] > newScore {
			BestScores[next] = newScore
		}
		if !found || newScore < bestScore {
			found = true
			bestScore = newScore
			if bestScore < BestScore {
				BestScore = bestScore
			}
			// if BestScores[curr] == 0 || bestScore < BestScores[curr] {
			// 	BestScores[curr] = bestScore
			// }
			// TEST
			// return bestScore, true
		}
		// if bestScore == BestScore {
		// 	BestScores[curr] = bestScore
		// }
	}
	if found {
		return bestScore, true
	}
	return 0, false
}

func Stage1(input io.Reader) (any, error) {
	m := utils.Must(matrix.NewFromReader(input))
	start, ok := m.Find('S')
	if !ok {
		return nil, fmt.Errorf("not found")
	}
	end, ok := m.Find('E')
	if !ok {
		return nil, fmt.Errorf("not found")
	}
	visited := make(sets.Set[matrix.Coord])
	go func() {
		for range time.Tick(1 * time.Second) {
			fmt.Println(BestScore, Ops)
		}
	}()
	res, ok := FindPath(m, visited, matrix.Right, 0, start, end)
	fmt.Println(Ops)
	if ok {
		// if res <= 149516 {
		// 	panic("false")
		// }
		return res, nil
	}
	return 0, nil
}

func Stage2(input io.Reader) (any, error) {
	m := utils.Must(matrix.NewFromReader(input))
	start, ok := m.Find('S')
	if !ok {
		return nil, fmt.Errorf("not found")
	}
	end, ok := m.Find('E')
	if !ok {
		return nil, fmt.Errorf("not found")
	}
	visited := make(sets.Set[matrix.Coord])
	go func() {
		for range time.Tick(1 * time.Second) {
			fmt.Println(BestScore, Ops)
		}
	}()
	res, ok := FindPath(m, visited, matrix.Right, 0, start, end)
	total := 0
	mv := m.Clone()
	for c, v := range BestScores {
		if res == v {
			total++
		}
		mv.SetAtCoord(c, 'O')
	}
	stage.Println(matrix.SMatrix(mv))
	// fmt.Println(Ops)
	// > 451
	return total + 1, nil
}
