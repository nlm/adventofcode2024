package main

import (
	"testing"

	"github.com/nlm/adventofcode2024/internal/stage"
)

var Stage1TestCases = []stage.TestCase{
	{
		Name:   "example0",
		Result: 4,
	},
	{
		Name:   "example1",
		Result: 18,
	},
	{
		Name:   "example1bis",
		Result: 18,
	},
	{
		Name:   "input",
		Result: 2517,
	},
}

var Stage2TestCases = []stage.TestCase{
	{
		Name:   "example2",
		Result: 9,
	},
	{
		Name:   "input",
		Result: 1960,
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
