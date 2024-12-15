package main

import (
	"testing"

	"github.com/nlm/adventofcode2024/internal/stage"
)

var Stage1TestCases = []stage.TestCase{
	{
		Name:   "example",
		Result: 10092,
	},
	{
		Name:   "example2",
		Result: 2028,
	},
	{
		Name:   "input",
		Result: 1412971,
	},
}

var Stage2TestCases = []stage.TestCase{
	{
		Name:   "example",
		Result: 9021,
	},
	{
		Name:   "input",
		Result: 1429299,
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
