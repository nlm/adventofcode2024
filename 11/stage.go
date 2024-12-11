package main

import (
	"fmt"
	"io"
	"slices"
	"strings"

	"github.com/nlm/adventofcode2024/internal/iterators"
	"github.com/nlm/adventofcode2024/internal/utils"
)

func SplitStone(stone int) []int {
	sStone := fmt.Sprint(stone)
	digits := len(sStone)
	return []int{
		utils.MustAtoi(sStone[:digits/2]),
		utils.MustAtoi(sStone[digits/2:]),
	}
}

// func UpdateLine(stones []int) []int {
// 	newStones := make([]int, 0, len(stones))
// 	for _, stone := range stones {
// 		if stone == 0 {
// 			newStones = append(newStones, 1)
// 			continue
// 		}
// 		if len(fmt.Sprint(stone))%2 == 0 {
// 			newStones = append(newStones, SplitStones[stone]...)
// 			continue
// 		}
// 		newStones = append(newStones, stone*2024)
// 	}
// 	return newStones
// }

func SolveOne(stone int) []int {
	if stone == 0 {
		return []int{1}
	}
	if len(fmt.Sprint(stone))%2 == 0 {
		return SplitStone(stone)
	}
	return []int{stone * 2024}
}

var Cache = make(map[[2]int]int)

func GetLengthIn(stone int, iterations int) int {
	if iterations == 0 {
		return 1
	}

	// check cache
	key := [2]int{stone, iterations}
	if v, ok := Cache[key]; ok {
		return v
	}

	// new result
	newStones := SolveOne(stone)
	total := 0
	for _, s := range newStones {
		total += GetLengthIn(s, iterations-1)
	}

	// put result in cache
	Cache[key] = total
	return total
}

func Stage1(input io.Reader) (any, error) {
	line := iterators.MapSlice(strings.Fields(slices.Collect(iterators.MustLines(input))[0]), utils.MustAtoi)
	total := 0
	for _, item := range line {
		total += GetLengthIn(item, 25)
	}
	return total, nil
}

func Stage2(input io.Reader) (any, error) {
	line := iterators.MapSlice(strings.Fields(slices.Collect(iterators.MustLines(input))[0]), utils.MustAtoi)
	total := 0
	for _, item := range line {
		total += GetLengthIn(item, 75)
	}
	return total, nil
}
