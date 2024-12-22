package maze

import (
	"testing"

	"github.com/nlm/adventofcode2024/internal/matrix"
	"github.com/nlm/adventofcode2024/internal/sets"
	"github.com/stretchr/testify/assert"
)

func TestCoordsNotEqual(t *testing.T) {
	m := matrix.New[byte](5, 9)
	assert.NotEqual(t,
		CoordToId(m, matrix.Coord{X: 2, Y: 5}),
		CoordToId(m, matrix.Coord{X: 5, Y: 2}),
	)
}

func TestCoordsInvert(t *testing.T) {
	m := matrix.New[byte](42, 42)
	ids := make(sets.Set[int64])
	for c := range m.Coords() {
		id := CoordToId(m, c)
		if ids.Contains(id) {
			t.Error("non-unique id:", id)
		}
		assert.Equal(t, c, IdToCoord(m, id))
		ids.Add(id)
	}
}
