package matrix

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	m := New[byte](4, 2)
	assert.Len(t, m.Data, 4*2)
	assert.Equal(t, m.Len, Coord{4, 2})
}

func TestAt(t *testing.T) {
	m := New[byte](4, 2)
	for i := 0; i < len(m.Data); i++ {
		m.Data[i] = byte(i)
	}
	for _, tc := range []struct {
		Coord Coord
		Value byte
	}{
		{Coord{0, 0}, 0},
		{Coord{0, 1}, 4},
		{Coord{1, 0}, 1},
		{Coord{1, 1}, 5},
		{Coord{2, 0}, 2},
		{Coord{3, 0}, 3},
		{Coord{3, 1}, 7},
	} {
		t.Run("AtCoord"+tc.Coord.String(), func(t *testing.T) {
			assert.Equal(t, tc.Value, m.AtCoord(tc.Coord))
		})
		t.Run("At"+tc.Coord.String(), func(t *testing.T) {
			assert.Equal(t, tc.Value, m.At(tc.Coord.X, tc.Coord.Y))
		})
	}
}

func TestCopy(t *testing.T) {
	m1 := New[int](2, 2)
	m2 := New[int](2, 2)
	assert.Equal(t, []int{0, 0, 0, 0}, m1.Data)
	assert.Equal(t, []int{0, 0, 0, 0}, m2.Data)
	m1.Fill(1)
	assert.Equal(t, []int{1, 1, 1, 1}, m1.Data)
	m2.Copy(m1)
	assert.Equal(t, []int{1, 1, 1, 1}, m2.Data)
}

func BenchmarkMatrix(b *testing.B) {
	for range b.N {
		m := New[byte](100, 100)
		m.Fill('X')
		for c := range m.Coords() {
			assert.Equal(b, byte('X'), m.AtCoord(c))
		}
	}
}
