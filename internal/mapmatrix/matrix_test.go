package mapmatrix

import (
	"testing"

	"github.com/nlm/adventofcode2023/internal/matrix"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	m := New[byte](4, 2)
	assert.Len(t, m.Data, 0)
	m.Fill('.')
	assert.Len(t, m.Data, 4*2)
	assert.Equal(t, m.DownRight, matrix.Coord{X: 0, Y: 0})
	assert.Equal(t, m.DownRight, matrix.Coord{X: 4, Y: 2})
}

func TestAt(t *testing.T) {
	m := New[byte](4, 2)
	for y := m.UpLeft.Y; y < m.DownRight.Y; y++ {
		for x := m.UpLeft.X; x < m.DownRight.X; x++ {
			m.SetAt(x, y, byte(y*x+x))
		}
	}
	for _, tc := range []struct {
		Coord matrix.Coord
		Value byte
	}{
		{matrix.Coord{X: 0, Y: 0}, 0},
		{matrix.Coord{X: 0, Y: 1}, 4},
		{matrix.Coord{X: 1, Y: 0}, 1},
		{matrix.Coord{X: 1, Y: 1}, 5},
		{matrix.Coord{X: 2, Y: 0}, 2},
		{matrix.Coord{X: 3, Y: 0}, 3},
		{matrix.Coord{X: 3, Y: 1}, 7},
	} {
		t.Run("AtCoord"+tc.Coord.String(), func(t *testing.T) {
			assert.Equal(t, tc.Value, m.AtCoord(tc.Coord))
		})
		t.Run("At"+tc.Coord.String(), func(t *testing.T) {
			assert.Equal(t, tc.Value, m.At(tc.Coord.X, tc.Coord.Y))
		})
	}
}
