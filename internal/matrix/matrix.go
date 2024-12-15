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
	Size Vec
}

func (m *Matrix[T]) Clone() *Matrix[T] {
	data := make([]T, len(m.Data))
	copy(data, m.Data)
	return &Matrix[T]{
		Data: data,
		Size: m.Size,
	}
}

var ErrInconsistentGeometry = fmt.Errorf("inconsistent geometry")

// New allocates a new Matrix of size x * y.
func New[T comparable](x, y int) *Matrix[T] {
	return &Matrix[T]{
		Data: make([]T, x*y),
		Size: Vec{x, y},
	}
}

// NewFromReader reads lines from the Reader and builds a new Matrix[byte]
// where each line is a row of the matrix, starting from the top.
// Each line must be of the same length, else an error will be returned.
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
	matrix.Size.X = cols
	matrix.Size.Y = rows
	return matrix, nil
}

// 1111
// 2222
// 3333
// 4444
//
// 1111222233334444

// Find searches for a value in the Matrix and returns
// a coordinate of the first match and a boolean indicating
// if the value was found.
func (m *Matrix[T]) Find(value T) (Coord, bool) {
	for i := 0; i < len(m.Data); i++ {
		if m.Data[i] == value {
			return Coord{i % m.Size.X, i / m.Size.X}, true
		}
	}
	return Coord{}, false
}

// Count counts the number of occurences of a value in the Matrix.
func (m *Matrix[T]) Count(value T) int {
	count := 0
	for _, v := range m.Data {
		if v == value {
			count++
		}
	}
	return count
}

// Fill will fill the matrix with a given value.
func (m *Matrix[T]) Fill(value T) {
	for i := 0; i < len(m.Data); i++ {
		m.Data[i] = value
	}
}

// Copy copies data from a matrix into the current matrix.
// If geometries are different, an error will be returned.
func (m *Matrix[T]) Copy(src *Matrix[T]) error {
	if src.Size != m.Size {
		return ErrInconsistentGeometry
	}
	copy(m.Data, src.Data)
	return nil
}

// Coords will return an iterator over all the coordinates
// that exist in the matrix.
func (m *Matrix[T]) Coords() iter.Seq[Coord] {
	return func(yield func(Coord) bool) {
		for y := 0; y < m.Size.Y; y++ {
			for x := 0; x < m.Size.X; x++ {
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

// AtCoord returns the value present at a coordinate.
// It's the responsibility of the user to check that
// the coordinate exists within the Matrix with InCoord.
func (m *Matrix[T]) AtCoord(c Coord) T {
	return m.At(c.X, c.Y)
}

// At returns the value present at a coordinate.
// It's the responsibility of the user to check that
// the coordinate exists within the Matrix with In.
func (m *Matrix[T]) At(x, y int) T {
	return m.Data[y*m.Size.X+x]
}

// SetAt sets the value present at a coordinate.
// It's the responsibility of the user to check that
// the coordinate exists within the Matrix with In.
func (m *Matrix[T]) SetAt(x, y int, value T) {
	m.Data[y*m.Size.X+x] = value
}

// SetAtCoord sets the value present at a coordinate.
// It's the responsibility of the user to check that
// the coordinate exists within the Matrix with InCoord.
func (m *Matrix[T]) SetAtCoord(c Coord, value T) {
	m.SetAt(c.X, c.Y, value)
}

// In checks that a coordinate exists within the Matrix.
func (m *Matrix[T]) In(x, y int) bool {
	return x >= 0 && x <= m.Size.X-1 && y >= 0 && y <= m.Size.Y-1
}

// InCoord checks that a coordinate exists within the Matrix.
func (m *Matrix[T]) InCoord(c Coord) bool {
	return m.In(c.X, c.Y)
}

func SMatrix(m *Matrix[byte]) string {
	sb := strings.Builder{}
	for y := 0; y < m.Size.Y; y++ {
		sb.Write(m.Data[y*m.Size.X : (y+1)*m.Size.X])
		sb.WriteByte('\n')
	}
	return sb.String()
}

func (m *Matrix[T]) String() string {
	sb := strings.Builder{}
	for y := 0; y < m.Size.Y; y++ {
		fmt.Fprint(&sb, m.Data[y*m.Size.X:(y+1)*m.Size.X])
		sb.WriteByte('\n')
	}
	return sb.String()
}
