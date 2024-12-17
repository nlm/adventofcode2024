package matrix

import "fmt"

type Vec XY

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

var VecName = map[Vec]string{
	Left:      "Left",
	Right:     "Right",
	Up:        "Up",
	Down:      "Down",
	UpLeft:    "UpLeft",
	UpRight:   "UpRight",
	DownLeft:  "DownLeft",
	DownRight: "DownRight",
}

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
	return fmt.Sprintf("{X: %d, Y: %d}", v.X, v.Y)
}

// Inv returns the inverse of this vector
func (v Vec) Inv() Vec {
	return Vec{X: -v.X, Y: -v.Y}
}
