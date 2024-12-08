package main

import (
	"testing"

	"github.com/nlm/adventofcode2024/internal/stage"
)

var Stage1TestCases = []stage.TestCase{
	{
		Name:   "example",
		Result: 14,
	},
	{
		Name:   "example2",
		Result: 2,
	},
	{
		Name:   "example4",
		Result: 4,
	},
	{
		Name:   "input",
		Result: 261,
	},
}

var Stage2TestCases = []stage.TestCase{
	{
		Name:   "example",
		Result: 34,
	},
	{
		Name:   "example5",
		Result: 9,
	},
	{
		Name:   "input",
		Result: 898,
	},
}

// Do not edit below

func TestStage1(t *testing.T) {
	stage.Test(t, Stage1, Stage1TestCases)
}

func TestStage2(t *testing.T) {
	stage.Test(t, Stage2, Stage2TestCases)
}

func BenchmarkStage1(b *testing.B) {
	stage.Benchmark(b, Stage1, Stage1TestCases)
}

func BenchmarkStage2(b *testing.B) {
	stage.Benchmark(b, Stage2, Stage2TestCases)
}
