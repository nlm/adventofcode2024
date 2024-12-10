package main

import (
	"testing"

	"github.com/nlm/adventofcode2024/internal/stage"
)

var Stage1TestCases = []stage.TestCase{
	{
		Name:   "example",
		Result: 36,
	},
	{
		Name:   "example2",
		Result: 2,
	},
	{
		Name:   "example3",
		Result: 4,
	},
	{
		Name:   "example4",
		Result: 1,
	},
	{
		Name:   "input",
		Result: 694,
	},
}

var Stage2TestCases = []stage.TestCase{
	{
		Name:   "example",
		Result: 81,
	},
	{
		Name:   "example2",
		Result: 2,
	},
	{
		Name:   "example3",
		Result: 13,
	},
	{
		Name:   "example4",
		Result: 3,
	},
	{
		Name:   "input",
		Result: 1497,
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
