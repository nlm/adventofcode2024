package main

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"regexp"
	"slices"

	"github.com/nlm/adventofcode2023/internal/utils"
)

var lineRe = regexp.MustCompile(`(\d+)\s+(\d+)`)

func readAndSortLists(input io.Reader) ([]int, []int, error) {
	var (
		list1 = make([]int, 0)
		list2 = make([]int, 0)
	)
	s := bufio.NewScanner(input)
	for s.Scan() {
		items := lineRe.FindSubmatch(s.Bytes())
		if items == nil {
			return nil, nil, fmt.Errorf("no match")
		}
		list1 = append(list1, utils.MustAtoi(string(items[1])))
		list2 = append(list2, utils.MustAtoi(string(items[2])))
	}
	list1 = slices.Sorted(slices.Values(list1))
	list2 = slices.Sorted(slices.Values(list2))
	return list1, list2, nil
}

func Stage1(input io.Reader) (any, error) {
	list1, list2, err := readAndSortLists(input)
	if err != nil {
		return nil, err
	}
	total := 0
	for i := range len(list1) {
		total += int(math.Abs(float64(list1[i] - list2[i])))
	}
	return total, nil
}

func countAppearance(list []int, value int) int {
	appearances := 0
	for _, v := range list {
		if v == value {
			appearances++
		}
		if v > value {
			break
		}
	}
	return appearances
}

func Stage2(input io.Reader) (any, error) {
	list1, list2, err := readAndSortLists(input)
	if err != nil {
		return nil, err
	}
	total := 0
	for i := range len(list1) {
		total += list1[i] * countAppearance(list2, list1[i])
	}
	return total, nil
}
