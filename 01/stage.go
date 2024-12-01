package main

import (
	"io"
	"sort"
	"strings"

	"github.com/nlm/adventofcode2024/internal/iterators"
	"github.com/nlm/adventofcode2024/internal/math"
	"github.com/nlm/adventofcode2024/internal/utils"
)

func readLists(input io.Reader) ([]int, []int) {
	var (
		list1 = make([]int, 0, 1024)
		list2 = make([]int, 0, 1024)
	)
	for line := range iterators.MustLines(input) {
		items := strings.Fields(line)
		list1 = append(list1, utils.MustAtoi(items[0]))
		list2 = append(list2, utils.MustAtoi(items[1]))
	}
	return list1, list2
}

func Stage1(input io.Reader) (any, error) {
	list1, list2 := readLists(input)
	sort.IntSlice(list1).Sort()
	sort.IntSlice(list2).Sort()
	total := 0
	for i := range len(list1) {
		total += math.Abs(list1[i] - list2[i])
	}
	return total, nil
}

func occurrences(list []int) map[int]int {
	occur := make(map[int]int, len(list))
	for _, v := range list {
		occur[v]++
	}
	return occur
}

func Stage2(input io.Reader) (any, error) {
	list1, list2 := readLists(input)
	occur := occurrences(list2)
	total := 0
	for i := range len(list1) {
		total += list1[i] * occur[list1[i]]
	}
	return total, nil
}
