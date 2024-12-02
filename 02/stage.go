package main

import (
	"io"
	"iter"
	"strings"

	"github.com/nlm/adventofcode2024/internal/iterators"
	"github.com/nlm/adventofcode2024/internal/stage"
	"github.com/nlm/adventofcode2024/internal/utils"
)

func reportIsSafe(items []int) bool {
	stage.Println(items)
	last := items[0]
	incr := items[1] > items[0]
	stage.Println("increasing ?", incr)
	for _, item := range items[1:] {
		stage.Println("next", item)
		switch incr {
		case true:
			if item-last < 1 || item-last > 3 {
				stage.Println("not incr", last, "->", item)
				return false
			}
		case false:
			if item-last < -3 || item-last > -1 {
				stage.Println("not decr", last, "->", item)
				return false
			}
		}
		last = item
	}
	return true
}

// func countReportUnsafeness(items []int) int {
// 	unsafe := 0
// 	stage.Println(items)
// 	last := items[0]
// 	incr := items[1] > items[0]
// 	stage.Println("increasing ?", incr)
// 	for _, item := range items[1:] {
// 		stage.Println("next", item)
// 		switch incr {
// 		case true:
// 			if item-last < 1 || item-last > 3 {
// 				stage.Println("not incr", last, "->", item)
// 				unsafe++
// 				continue
// 			}
// 		case false:
// 			if item-last < -3 || item-last > -1 {
// 				stage.Println("not decr", last, "->", item)
// 				unsafe++
// 				continue
// 			}
// 		}
// 		last = item
// 	}
// 	stage.Println("Count:", unsafe)
// 	return unsafe
// }

func Stage1(input io.Reader) (any, error) {
	total := 0
	for line := range iterators.MustLines(input) {
		items := iterators.Map(strings.Fields(line), utils.MustAtoi)
		if reportIsSafe(items) {
			total++
		}
	}
	return total, nil
}

func eachReportMinusOne(items []int) iter.Seq[[]int] {
	return func(yield func([]int) bool) {
		for i := 0; i < len(items); i++ {
			newItems := make([]int, 0, len(items))
			newItems = append(newItems, items[:i]...)
			newItems = append(newItems, items[i+1:]...)
			if !yield(newItems) {
				return
			}
		}
	}
}

func Stage2(input io.Reader) (any, error) {
	total := 0
	for line := range iterators.MustLines(input) {
		items := iterators.Map(strings.Fields(line), utils.MustAtoi)
		for report := range eachReportMinusOne(items) {
			if reportIsSafe(report) {
				total++
				break
			}
		}
	}
	return total, nil
}
