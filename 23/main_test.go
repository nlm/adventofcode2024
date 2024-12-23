package main

import (
	"testing"

	"github.com/nlm/adventofcode2024/internal/stage"
)

var Stage1TestCases = []stage.TestCase{
	{
		Name:   "example",
		Result: 7,
	},
	{
		Name:   "input",
		Result: 1323,
	},
}

var Stage2TestCases = []stage.TestCase{
	{
		Name:   "example",
		Result: "co,de,ka,ta",
	},
	{
		Name:   "input",
		Result: "er,fh,fi,ir,kk,lo,lp,qi,ti,vb,xf,ys,yu",
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
