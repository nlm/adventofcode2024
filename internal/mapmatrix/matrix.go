package mapmatrix

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/nlm/adventofcode2023/internal/matrix"
)

type Matrix[T comparable] struct {
	Data      map[matrix.Coord]T
	UpLeft    matrix.Coord
	DownRight matrix.Coord
	// AutoExpand bool
}

func (m *Matrix[T]) Clone() *Matrix[T] {
	data := make(map[matrix.Coord]T, len(m.Data))
	for k, v := range m.Data {
		data[k] = v
	}
	return &Matrix[T]{
		Data:      data,
		UpLeft:    m.UpLeft,
		DownRight: m.DownRight,
	}
}

var ErrInconsistentGeometry = fmt.Errorf("inconsistent geometry")

func New[T comparable](x, y int) *Matrix[T] {
	return &Matrix[T]{
		Data:      make(map[matrix.Coord]T, x*y),
		UpLeft:    matrix.Coord{X: 0, Y: 0},
		DownRight: matrix.Coord{X: x, Y: y},
	}
}

func NewFromReader(input io.Reader) (*Matrix[byte], error) {
	m := New[byte](0, 0)
	s := bufio.NewScanner(input)
	cols := 0
	rows := 0
	for s.Scan() {
		if cols != 0 {
			if len(s.Bytes()) != cols {
				return nil, ErrInconsistentGeometry
			}
		} else {
			cols = len(s.Bytes())
		}
		for x, b := range s.Bytes() {
			m.Data[matrix.Coord{
				X: x,
				Y: rows,
			}] = b
		}
		rows++
	}
	if s.Err() != nil {
		return nil, s.Err()
	}
	m.DownRight = matrix.Coord{
		X: cols,
		Y: rows,
	}
	return m, nil
}

// 1111
// 2222
// 3333
// 4444
//
// 1111222233334444
func (m *Matrix[T]) Find(value T) (matrix.Coord, bool) {
	for k, v := range m.Data {
		if v == value {
			return k, true
		}
	}
	return matrix.Coord{}, false
}

func (m *Matrix[T]) Fill(value T) {
	for y := m.UpLeft.Y; y < m.DownRight.Y; y++ {
		for x := m.UpLeft.X; x < m.DownRight.X; x++ {
			m.Data[matrix.Coord{X: x, Y: y}] = value
		}
	}
}

func (m *Matrix[T]) AtCoord(c matrix.Coord) T {
	return m.Data[c]
}

func (m *Matrix[T]) At(x, y int) T {
	return m.AtCoord(matrix.Coord{X: x, Y: y})
}

func (m *Matrix[T]) SetAt(x, y int, value T) {
	m.SetAtCoord(matrix.Coord{X: x, Y: y}, value)
}

func (m *Matrix[T]) SetAtCoord(c matrix.Coord, value T) {
	m.Data[c] = value
}

func (m *Matrix[T]) In(x, y int) bool {
	return x >= m.UpLeft.X && x <= m.DownRight.X-1 && y >= m.UpLeft.Y && y <= m.DownRight.Y-1
}

func (m *Matrix[T]) InCoord(c matrix.Coord) bool {
	return m.In(c.X, c.Y)
}

func SMatrix(m *Matrix[byte]) string {
	sb := strings.Builder{}
	for y := m.UpLeft.Y; y < m.DownRight.Y; y++ {
		for x := m.UpLeft.X; x < m.DownRight.X; x++ {
			sb.WriteByte(m.At(x, y))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func (m *Matrix[T]) String() string {
	sb := strings.Builder{}
	for y := m.UpLeft.Y; y < m.DownRight.Y; y++ {
		for x := m.UpLeft.X; x < m.DownRight.X; x++ {
			fmt.Fprint(&sb, m.At(x, y))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}
