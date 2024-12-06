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
