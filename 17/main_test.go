package main

import (
	"testing"

	"github.com/nlm/adventofcode2024/internal/stage"
)

var Stage1TestCases = []stage.TestCase{
	{
		Name:   "example",
		Result: "4,6,3,5,6,3,5,2,1,0",
	},
	{
		Name:   "input",
		Result: "7,6,1,5,3,1,4,2,6",
	},
}

var Stage2TestCases = []stage.TestCase{
	// {
	// 	Name:   "example5",
	// 	Result: 117440,
	// },
	// {
	// 	Name:   "input",
	// 	Result: 164541017976509,
	// },
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
