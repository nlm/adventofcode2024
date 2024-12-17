package matrix

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRotate(t *testing.T) {
	for _, tc := range []struct {
		Src Vec
		Ang int
		Dst Vec
	}{
		{Left, 90, Up},
		{Up, 90, Right},
		{Right, 90, Down},
		{Down, 90, Left},
		{Up, 360, Up},
		{Up, 720, Up},
		{Left, -90, Down},
		{Up, 45, UpRight},
		{UpRight, 45, Right},
	} {
		t.Run(fmt.Sprint(VecName[tc.Src], "+", tc.Ang, "=", VecName[tc.Dst]), func(t *testing.T) {
			assert.Equal(t, tc.Dst, tc.Src.Rotate(tc.Ang))
		})
	}
}
