package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"

	"github.com/nlm/adventofcode2024/internal/matrix"
	"github.com/nlm/adventofcode2024/internal/stage"
	"github.com/nlm/adventofcode2024/internal/utils"
)

func ParseInput(input io.Reader) (*matrix.Matrix[byte], []byte) {
	// bytes.Split(\n\n)
	s := bufio.NewScanner(input)
	mbuffer := bytes.NewBuffer(nil)
	for s.Scan() {
		if len(s.Text()) == 0 {
			break
		}
		mbuffer.Write(s.Bytes())
		mbuffer.WriteByte('\n')
	}
	instructions := bytes.NewBuffer(nil)
	for s.Scan() {
		instructions.Write(s.Bytes())
	}
	m := utils.Must(matrix.NewFromReader(bytes.NewReader(mbuffer.Bytes())))
	return m, instructions.Bytes()
}

var Dir = map[byte]matrix.Vec{
	'v': matrix.Down,
	'^': matrix.Up,
	'>': matrix.Right,
	'<': matrix.Left,
}

func MoveBox(m *matrix.Matrix[byte], coord matrix.Coord, vec matrix.Vec) bool {
	nextCoord := coord.Add(vec)
	switch m.AtCoord(nextCoord) {
	case '#':
		stage.Println("box hit wall")
		return false
	case '.':
		stage.Println("box move ok")
		m.SetAtCoord(nextCoord, 'O')
		m.SetAtCoord(coord, '.')
		return true
	case 'O':
		if MoveBox(m, nextCoord, vec) {
			stage.Println("box move recuse ok")
			m.SetAtCoord(nextCoord, 'O')
			m.SetAtCoord(coord, '.')
			return true
		}
		stage.Println("box move recurse fail")
		return false
	default:
		panic("unknown move box")
	}
}

func Stage1(input io.Reader) (any, error) {
	m, instructions := ParseInput(input)
	stage.Println(matrix.SMatrix(m))
	stage.Println(string(instructions))
	robot, ok := m.Find('@')
	if !ok {
		panic("bot not found")
	}
	for _, inst := range instructions {
		vec := Dir[inst]
		// stage.Println(vec)
		newRobot := robot.Add(vec)
		stage.Println("move", string(inst))
		switch m.AtCoord(newRobot) {
		case '#':
			// Wall
			stage.Println("hit wall")
		case '.':
			stage.Println("move ok")
			m.SetAtCoord(robot, '.')
			m.SetAtCoord(newRobot, '@')
			robot = newRobot
		case 'O':
			if MoveBox(m, newRobot, vec) {
				m.SetAtCoord(robot, '.')
				m.SetAtCoord(newRobot, '@')
				robot = newRobot
				stage.Println("move box ok")
			} else {
				stage.Println("move box blocked")
			}
		}
		stage.Println(matrix.SMatrix(m))
	}
	total := 0
	for c := range m.Coords() {
		if m.AtCoord(c) == 'O' {
			total += c.X + c.Y*100
		}
	}
	return total, nil
}

func ExpandMatrix(m *matrix.Matrix[byte]) *matrix.Matrix[byte] {
	m2 := matrix.New[byte](m.Size.X*2, m.Size.Y)
	m2.Data = m2.Data[:0]
	for _, b := range m.Data {
		switch b {
		case '#':
			m2.Data = append(m2.Data, '#', '#')
		case '.':
			m2.Data = append(m2.Data, '.', '.')
		case '@':
			m2.Data = append(m2.Data, '@', '.')
		case 'O':
			m2.Data = append(m2.Data, '[', ']')
		default:
			panic("expand error")
		}
	}
	return m2
}

func AddVec2(coords [2]matrix.Coord, vec matrix.Vec) [2]matrix.Coord {
	var nextCoords [2]matrix.Coord
	nextCoords[0] = coords[0].Add(vec)
	nextCoords[1] = coords[1].Add(vec)
	return nextCoords
}

func CanMoveBigBox(m *matrix.Matrix[byte], coords [2]matrix.Coord, vec matrix.Vec) bool {
	nextCoords := AddVec2(coords, vec)
	stage.Println("CanMove?", coords, "+", vec, "->", nextCoords)
	switch vec {
	case matrix.Left:
		// look right of next
		// [][]
		switch m.AtCoord(coords[0].Add(vec)) {
		case '.':
			// ok
			return true
		case '#':
			// wall
			return false
		case ']':
			return CanMoveBigBox(m, AddVec2(coords, vec.Mul(2)), vec)
		// case '[':
		// 	// panic("[")
		default:
			panic("XL")
		}
	case matrix.Right:
		// look right of next
		// [][]
		//   ^
		switch m.AtCoord(coords[1].Add(vec)) {
		case '.':
			// ok
			return true
		case '#':
			// wall
			return false
		case '[':
			return CanMoveBigBox(m, AddVec2(coords, vec.Mul(2)), vec)
		// case ']':
		default:
			panic("XR")
		}
	case matrix.Up, matrix.Down:
		// wall
		for _, nc := range nextCoords {
			if m.AtCoord(nc) == '#' {
				return false
			}
		}
		if m.AtCoord(nextCoords[0]) == '.' && m.AtCoord(nextCoords[1]) == '.' {
			stage.Println("empty")
			return true
		}
		// []  []  []  [][]
		//  [] [] []    []
		if m.AtCoord(nextCoords[0]) == m.AtCoord(coords[0]) {
			// box is X-aligned
			stage.Println("box is x-aligned")
			return CanMoveBigBox(m, nextCoords, vec)
		}
		if m.AtCoord(nextCoords[0]) == ']' {
			// box is at left
			if !CanMoveBigBox(m, AddVec2(nextCoords, matrix.Left), vec) {
				return false
			}
		}
		if m.AtCoord(nextCoords[1]) == '[' {
			// box is at right
			if !CanMoveBigBox(m, AddVec2(nextCoords, matrix.Right), vec) {
				return false
			}
		}
		return true
	default:
		panic("x")
	}
}

func DoMoveBigBox(m *matrix.Matrix[byte], coords [2]matrix.Coord, vec matrix.Vec) {
	nextCoords := AddVec2(coords, vec)
	switch vec {
	case matrix.Left:
		switch m.AtCoord(coords[0].Add(vec)) {
		case ']':
			DoMoveBigBox(m, AddVec2(coords, vec.Mul(2)), vec)
		}
	case matrix.Right:
		// look right of next
		switch m.AtCoord(coords[1].Add(vec)) {
		case '[':
			DoMoveBigBox(m, AddVec2(coords, vec.Mul(2)), vec)
		}
	case matrix.Up, matrix.Down:
		// []  []  []  [][]
		//  [] [] []    []
		if m.AtCoord(nextCoords[0]) == m.AtCoord(coords[0]) {
			// box is X-aligned
			DoMoveBigBox(m, nextCoords, vec)
		}
		if m.AtCoord(nextCoords[0]) == ']' {
			// box is at left
			DoMoveBigBox(m, AddVec2(nextCoords, matrix.Left), vec)
		}
		if m.AtCoord(nextCoords[1]) == '[' {
			// box is at right
			DoMoveBigBox(m, AddVec2(nextCoords, matrix.Right), vec)
		}
	default:
		panic("x")
	}
	stage.Println("DoMove", coords, "+", vec)
	m.SetAtCoord(coords[0], '.')
	m.SetAtCoord(coords[1], '.')
	m.SetAtCoord(nextCoords[0], '[')
	m.SetAtCoord(nextCoords[1], ']')
}

// func DoMoveBigBox2(m *matrix.Matrix[byte], coords [2]matrix.Coord, vec matrix.Vec) {
// 	stage.Println("DoMove", coords, "+", vec)
// 	var nextCoords [2]matrix.Coord
// 	nextCoords[0] = coords[0].Add(vec)
// 	nextCoords[1] = coords[1].Add(vec)
// 	for _, nc := range nextCoords {
// 		switch m.AtCoord(nc) {
// 		case '[':
// 			if vec == matrix.Left {
// 				// it's ourselves
// 				continue
// 			}
// 			DoMoveBigBox(m, [2]matrix.Coord{nextCoords[0].Add(vec), nextCoords[1].Add(vec)}, vec)
// 		case ']':
// 			if vec == matrix.Right {
// 				// it's ourselves
// 				continue
// 			}
// 			DoMoveBigBox(m, [2]matrix.Coord{nextCoords[0].Add(vec), nextCoords[1].Add(vec)}, vec)
// 		}
// 	}
// 	m.SetAtCoord(coords[0], '.')
// 	m.SetAtCoord(coords[1], '.')
// 	m.SetAtCoord(nextCoords[0], '[')
// 	m.SetAtCoord(nextCoords[1], ']')
// }

func Stage2(input io.Reader) (any, error) {
	m, instructions := ParseInput(input)
	stage.Println(matrix.SMatrix(m))
	stage.Println(string(instructions))
	m = ExpandMatrix(m)
	stage.Println(matrix.SMatrix(m))
	robot, ok := m.Find('@')
	if !ok {
		return nil, fmt.Errorf("robot not found")
	}
	for _, inst := range instructions {
		vec := Dir[inst]
		// stage.Println(vec)
		newRobot := robot.Add(vec)
		stage.Println("move", string(inst))
		switch m.AtCoord(newRobot) {
		case '#':
			// Wall
			stage.Println("hit wall")
		case '.':
			stage.Println("move ok")
			m.SetAtCoord(robot, '.')
			m.SetAtCoord(newRobot, '@')
			robot = newRobot
		case '[', ']':
			// find box current coords
			var boxCoords [2]matrix.Coord
			b := m.AtCoord(newRobot)
			switch vec {
			case matrix.Left:
				if b == '[' {
					panic("unexp left")
				}
				boxCoords = [2]matrix.Coord{
					newRobot.Add(matrix.Left),
					newRobot,
				}
			case matrix.Right:
				if b == ']' {
					panic("unexp right")
				}
				boxCoords = [2]matrix.Coord{
					newRobot,
					newRobot.Add(matrix.Right),
				}
			case matrix.Up, matrix.Down:
				if b == '[' {
					boxCoords = [2]matrix.Coord{
						newRobot,
						newRobot.Add(matrix.Right),
					}
				} else if b == ']' {
					boxCoords = [2]matrix.Coord{
						newRobot.Add(matrix.Left),
						newRobot,
					}
				} else {
					panic("unex abc")
				}
			default:
				panic("unexpected boxcoord")
			}
			// move box
			if CanMoveBigBox(m, boxCoords, vec) {
				DoMoveBigBox(m, boxCoords, vec)
				m.SetAtCoord(robot, '.')
				m.SetAtCoord(newRobot, '@')
				robot = newRobot
				stage.Println("move box ok")
			} else {
				stage.Println("move box blocked")
			}
		default:
			return nil, fmt.Errorf("unexpected move")
		}
		stage.Println(matrix.SMatrix(m))
	}
	total := 0
	for c := range m.Coords() {
		if m.AtCoord(c) == '[' {
			total += c.X + c.Y*100
		}
	}
	return total, nil
}
