package main

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEach(t *testing.T) {
	results := slices.Collect(eachReportMinusOne([]int{0, 1, 2, 3, 4}))
	assert.Equal(t, [][]int{
		{1, 2, 3, 4},
		{0, 2, 3, 4},
		{0, 1, 3, 4},
		{0, 1, 2, 4},
		{0, 1, 2, 3},
	}, results)
}
