package main

import (
	"testing"

	"github.com/nlm/adventofcode2024/internal/stage"
)

var Stage1TestCases = []stage.TestCase{
	{
		Name:   "example",
		Result: 1930,
	},
	{
		Name:   "input",
		Result: 1450422,
	},
}

var Stage2TestCases = []stage.TestCase{
	{
		Name:   "example",
		Result: 1206,
	},
	{
		Name:   "example3",
		Result: 80,
	},
	{
		Name:   "example4",
		Result: 236,
	},
	{
		Name:   "example7",
		Result: 368,
	},
	{
		Name:   "input",
		Result: 906606,
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
