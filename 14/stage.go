package main

import (
	"flag"
	"fmt"
	"io"
	"iter"
	"regexp"

	"github.com/nlm/adventofcode2024/internal/iterators"
	"github.com/nlm/adventofcode2024/internal/matrix"
	"github.com/nlm/adventofcode2024/internal/sets"
	"github.com/nlm/adventofcode2024/internal/stage"
	"github.com/nlm/adventofcode2024/internal/utils"
)

var lineRe = regexp.MustCompile(`p=(-?\d+),(-?\d+) v=(-?\d+),(-?\d+)`)

type Robot struct {
	Pos matrix.Coord
	Vel matrix.Vec
}

func ParseInput(input io.Reader) []*Robot {
	robots := make([]*Robot, 0)
	for line := range iterators.MustLines(input) {
		match := lineRe.FindStringSubmatch(line)
		if len(match) != 5 {
			panic(fmt.Sprintf("parse error: %v", line))
		}
		robots = append(robots, &Robot{
			Pos: matrix.Coord{X: utils.MustAtoi(match[1]), Y: utils.MustAtoi(match[2])},
			Vel: matrix.Vec{X: utils.MustAtoi(match[3]), Y: utils.MustAtoi(match[4])},
		})
		// stage.Printf("%s -> %v\n", line, robots[len(robots)-1:])
	}
	return robots
}

var flagExample = flag.Bool("stage-is-example", false, "run on example surface")

func Stage1(input io.Reader) (any, error) {
	robots := ParseInput(input)
	// Handle
	lenX := 101
	lenY := 103
	if *flagExample {
		lenX = 11
		lenY = 7
	}
	for range 100 { // run for 100 seconds
		MoveRobots(robots, lenX, lenY)
	}
	ShowRobots(lenX, lenY, robots)
	product := 1
	for qTL, qDR := range Quadrants(lenX, lenY) {
		cnt := CountRobotsInArea(robots, qTL, qDR)
		stage.Println(qTL, qDR, "->", cnt)
		product *= cnt
	}
	return product, nil
}

func Quadrants(lenX, lenY int) iter.Seq2[matrix.Coord, matrix.Coord] {
	vecX := matrix.Vec{X: lenX/2 - 1, Y: 0}
	vecY := matrix.Vec{X: 0, Y: lenY/2 - 1}
	orig := matrix.Coord{X: 0, Y: 0}
	return func(yield func(matrix.Coord, matrix.Coord) bool) {
		for _, v := range [][2]matrix.Coord{
			// UpLeft
			{orig, orig.Add(vecX).Add(vecY)},
			// UpRight
			{orig.Add(vecX).Add(matrix.Vec{X: 2}), orig.Add(vecX).Add(matrix.Vec{X: 2}).Add(vecX).Add(vecY)},
			// DownLeft
			{orig.Add(vecY).Add(matrix.Vec{Y: 2}), orig.Add(vecY).Add(matrix.Vec{Y: 2}).Add(vecX).Add(vecY)},
			// DownRight
			{orig.Add(vecX).Add(vecY).Add(matrix.Vec{X: 2, Y: 2}), orig.Add(vecX).Add(vecY).Add(matrix.Vec{X: 2, Y: 2}).Add(vecX).Add(vecY)},
		} {
			if !yield(v[0], v[1]) {
				return
			}
		}
	}
}

func CountRobotsInArea(robots []*Robot, orig, end matrix.Coord) int {
	count := 0
	for _, r := range robots {
		if r.Pos.X >= orig.X &&
			r.Pos.X <= end.X &&
			r.Pos.Y >= orig.Y &&
			r.Pos.Y <= end.Y {
			count++
		}
	}
	return count
}

func HasCollisions(robots []*Robot) bool {
	s := make(sets.Set[matrix.Coord], len(robots))
	for _, r := range robots {
		if s.Contains(r.Pos) {
			return false
		}
		s.Add(r.Pos)
	}
	return true
}

func MoveRobots(robots []*Robot, lenX, lenY int) {
	for _, r := range robots {
		MoveRobot(r, lenX, lenY)
	}
}

func MoveRobot(r *Robot, lenX, lenY int) {
	r.Pos = r.Pos.Add(r.Vel)
	r.Pos.X = mod(r.Pos.X, lenX)
	r.Pos.Y = mod(r.Pos.Y, lenY)
}

func mod(a, b int) int {
	return (a%b + b) % b
}

func ShowRobots(x, y int, robots []*Robot) {
	m := matrix.New[byte](x, y)
	m.Fill('.')
	for _, r := range robots {
		cur := m.AtCoord(r.Pos)
		if cur == '.' {
			m.SetAtCoord(r.Pos, '1')
		} else {
			m.SetAtCoord(r.Pos, cur+1)
		}
	}
	stage.Println(matrix.SMatrix(m))
}

func RobotsToMatrix(robots []*Robot, lenX, lenY int) *matrix.Matrix[bool] {
	m := matrix.New[bool](lenX, lenY)
	for _, r := range robots {
		m.SetAtCoord(r.Pos, true)
	}
	return m
}

func DetectLine(m *matrix.Matrix[bool], cnt int) bool {
	inLine := 0
	for c := range m.Coords() {
		if m.AtCoord(c) {
			inLine++
		} else {
			inLine = 0
		}
		if inLine >= cnt {
			return true
		}
	}
	return false
}

func Stage2(input io.Reader) (any, error) {
	robots := ParseInput(input)
	lenX := 101
	lenY := 103
	// parts := make([]int, 0, 4)
	for i := range 100000 { // 100 seconds
		// parts = parts[:0]
		// stage.Println("===== after", i, "=====")
		// ShowRobots(lenX, lenY, robots)
		// for qTL, qDR := range Quadrants(lenX, lenY) {
		// 	cnt := CountRobotsInArea(robots, qTL, qDR)
		// 	parts = append(parts, cnt)
		// }
		// tlZoneTL := matrix.Coord{0, 0}
		// tlZoneDR := matrix.Coord{10, 10}
		// cntTl := CountRobotsInArea(robots, tlZoneTL, tlZoneDR)

		// trZoneTL := matrix.Coord{90, 0}
		// trZoneDR := matrix.Coord{100, 10}
		// cntTr := CountRobotsInArea(robots, trZoneTL, trZoneDR)

		// if cntTl == 0 && cntTr == 0 {
		// 	fmt.Println("Yippee")
		// 	time.Sleep(2 * time.Second)
		// }

		m := RobotsToMatrix(robots, lenX, lenY)
		if DetectLine(m, 8) {
			stage.Println("===== after", i, "=====")
			ShowRobots(lenX, lenY, robots)
			return i, nil
		}

		// if parts[0] == parts[1] && parts[2] == parts[3] {
		// 	fmt.Println("Yippee")
		// 	time.Sleep(2 * time.Second)
		// }
		MoveRobots(robots, lenX, lenY)
	}
	return 0, nil
}
