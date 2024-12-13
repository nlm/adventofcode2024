package main

import (
	"testing"

	"github.com/nlm/adventofcode2024/internal/stage"
)

var Stage1TestCases = []stage.TestCase{
	{
		Name:   "example",
		Result: int64(480),
	},
	{
		Name:   "input",
		Result: int64(33427),
	},
}

var Stage2TestCases = []stage.TestCase{
	{
		Name:   "example",
		Result: int64(875318608908),
	},
	{
		Name:   "input",
		Result: int64(91649162972270),
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
