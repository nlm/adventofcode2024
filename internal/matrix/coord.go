package matrix

import (
	"fmt"
)

type Coord struct {
	X int
	Y int
}

func (c Coord) Left() Coord {
	return c.Add(Left)
}

func (c Coord) Right() Coord {
	return c.Add(Right)
}

func (c Coord) Up() Coord {
	return c.Add(Up)
}

func (c Coord) Down() Coord {
	return c.Add(Down)
}

func (c Coord) String() string {
	return fmt.Sprintf("{X: %d, Y: %d}", c.X, c.Y)
}

func (c Coord) Add(v Vec) Coord {
	return Coord{c.X + v.X, c.Y + v.Y}
}

func (c Coord) Clone(v Vec) Coord {
	return c
}

func (c *Coord) Move(v Vec) {
	c.X += v.X
	c.Y += v.Y
}

type Vec Coord

var (
	Left      = Vec{X: -1, Y: 0}
	Right     = Vec{X: 1, Y: 0}
	Up        = Vec{X: 0, Y: -1}
	Down      = Vec{X: 0, Y: 1}
	UpLeft    = Up.Add(Left)
	UpRight   = Up.Add(Right)
	DownLeft  = Down.Add(Left)
	DownRight = Down.Add(Right)
)

// Add adds a vector to another vector.
func (v Vec) Add(v2 Vec) Vec {
	return Vec{X: v.X + v2.X, Y: v.Y + v2.Y}
}

// Mul multiplies a vector by a factor of n.
func (v Vec) Mul(n int) Vec {
	return Vec{X: v.X * n, Y: v.Y * n}
}

// String returns a string representation of Vec.
func (v Vec) String() string {
	return Coord(v).String()
}
