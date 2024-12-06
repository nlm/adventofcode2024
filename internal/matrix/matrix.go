package matrix

import (
	"bufio"
	"fmt"
	"io"
	"iter"
	"strings"
)

type Matrix[T comparable] struct {
	Data []T
	Len  Coord
}

func (m *Matrix[T]) Clone() *Matrix[T] {
	data := make([]T, len(m.Data))
	copy(data, m.Data)
	return &Matrix[T]{
		Data: data,
		Len:  m.Len,
	}
}

var ErrInconsistentGeometry = fmt.Errorf("inconsistent geometry")

func New[T comparable](x, y int) *Matrix[T] {
	return &Matrix[T]{
		Data: make([]T, x*y),
		Len:  Coord{x, y},
	}
}

func NewFromReader(input io.Reader) (*Matrix[byte], error) {
	matrix := &Matrix[byte]{}
	s := bufio.NewScanner(input)
	cols := -1
	rows := 0
	for s.Scan() {
		if cols != -1 {
			if len(s.Bytes()) != cols {
				return nil, ErrInconsistentGeometry
			}
		} else {
			cols = len(s.Bytes())
		}
		matrix.Data = append(matrix.Data, s.Bytes()...)
		rows++
	}
	if s.Err() != nil {
		return nil, s.Err()
	}
	matrix.Len.X = cols
	matrix.Len.Y = rows
	return matrix, nil
}

// 1111
// 2222
// 3333
// 4444
//
// 1111222233334444
func (m *Matrix[T]) Find(value T) (Coord, bool) {
	for i := 0; i < len(m.Data); i++ {
		if m.Data[i] == value {
			return Coord{i % m.Len.X, i / m.Len.X}, true
		}
	}
	return Coord{}, false
}

func (m *Matrix[T]) Count(value T) int {
	count := 0
	for _, v := range m.Data {
		if v == value {
			count++
		}
	}
	return count
}

func (m *Matrix[T]) Fill(value T) {
	for i := 0; i < len(m.Data); i++ {
		m.Data[i] = value
	}
}

func (m *Matrix[T]) IterCoords() iter.Seq[Coord] {
	return func(yield func(Coord) bool) {
		for y := 0; y < m.Len.Y; y++ {
			for x := 0; x < m.Len.X; x++ {
				if !yield(Coord{X: x, Y: y}) {
					return
				}
			}
		}
	}
}

// func (m *Matrix[T]) InsertLineBefore(y int, value T) {
// 	yLen := m.Len.Y
// 	// m.Data = append(m.Data, []byte{})
// 	copy(m.Data[1:2], m.Data[2:3])
// 	for j := yLen; j > y; j-- {
// 		m.Data[j] = m.Data[j-1]
// 	}
// 	m.Len.Y++
// }

// func (m *Matrix[T]) InsertColumnBefore(x int, value T) {
// 	xLen := m.Len.X
// 	for y := 0; y < m.Len.Y; y++ {
// 		m.Data[y] = append(m.Data[y], byte(0))
// 		for i := xLen; i > x; i-- {
// 			m.Data[y][i] = m.Data[y][i-1]
// 		}
// 		m.Data[y][x] = b
// 	}
// 	m.Len.X++
// }

func (m *Matrix[T]) AtCoord(c Coord) T {
	return m.At(c.X, c.Y)
}

func (m *Matrix[T]) At(x, y int) T {
	return m.Data[y*m.Len.X+x]
}

func (m *Matrix[T]) SetAt(x, y int, value T) {
	m.Data[y*m.Len.X+x] = value
}

func (m *Matrix[T]) SetAtCoord(c Coord, value T) {
	m.SetAt(c.X, c.Y, value)
}

func (m *Matrix[T]) In(x, y int) bool {
	return x >= 0 && x <= m.Len.X-1 && y >= 0 && y <= m.Len.Y-1
}

func (m *Matrix[T]) InCoord(c Coord) bool {
	return m.In(c.X, c.Y)
}

func SMatrix(m *Matrix[byte]) string {
	sb := strings.Builder{}
	for y := 0; y < m.Len.Y; y++ {
		sb.Write(m.Data[y*m.Len.X : (y+1)*m.Len.X])
		sb.WriteByte('\n')
	}
	return sb.String()
}

func (m *Matrix[T]) String() string {
	sb := strings.Builder{}
	for y := 0; y < m.Len.Y; y++ {
		fmt.Fprint(&sb, m.Data[y*m.Len.X:(y+1)*m.Len.X])
		sb.WriteByte('\n')
	}
	return sb.String()
}
