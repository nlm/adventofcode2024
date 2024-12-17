package sets

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContains(t *testing.T) {
	s := Set[int]{}
	assert.False(t, s.Contains(0))
	assert.NotContains(t, s, 0)
	s.Add(1)
	assert.False(t, s.Contains(0))
	assert.NotContains(t, s, 0)
	assert.True(t, s.Contains(1))
	assert.Contains(t, s, 1)
}

func TestLen(t *testing.T) {
	s := Set[int]{}
	// empty set
	assert.Equal(t, 0, len(s))

	// add 1 item sets length to 1
	s.Add(1)
	assert.Equal(t, 1, len(s))

	// adding same item does not change length
	s.Add(1)
	assert.Equal(t, 1, len(s))
}

func TestAppend(t *testing.T) {
	s := Set[int]{}
	assert.Equal(t, 0, len(s))

	s = Append(s, 1, 2, 3)
	assert.Equal(t, 3, len(s))
	assert.NotContains(t, s, 0)
	assert.True(t, s.Contains(1))
	assert.True(t, s.Contains(2))
	assert.True(t, s.Contains(3))

	s.Remove(3)
	assert.Equal(t, 2, len(s))
	assert.False(t, s.Contains(0))
	assert.True(t, s.Contains(1))
	assert.True(t, s.Contains(2))
	assert.False(t, s.Contains(3))
}

func TestNilSet(t *testing.T) {
	s := Set[int](nil)
	assert.NotContains(t, s, 0)
	assert.False(t, s.Contains(0))
	s.Remove(0)
	assert.NotContains(t, s, 0)
	assert.False(t, s.Contains(0))
	s = Append(s, 0)
	assert.Contains(t, s, 0)
	assert.True(t, s.Contains(0))
}

func TestIterSet(t *testing.T) {
	values := []int{1, 2, 4, 5, 3}
	s := Append(nil, values...)
	sl := slices.Collect(Values(s))
	assert.Len(t, sl, len(values))
	for v := range slices.Values(values) {
		assert.Contains(t, s, v)
	}
}

func TestClone(t *testing.T) {
	a := make(Set[int])
	a.Add(42)
	b := a.Clone()
	b.Add(43)
	assert.True(t, b.Contains(42))
	assert.False(t, a.Contains(43))
	assert.True(t, b.Contains(43))
}
