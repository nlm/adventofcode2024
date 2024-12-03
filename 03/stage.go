package main

import (
	"io"
	"regexp"

	"github.com/nlm/adventofcode2024/internal/iterators"
	"github.com/nlm/adventofcode2024/internal/stage"
	"github.com/nlm/adventofcode2024/internal/utils"
)

func Stage1(input io.Reader) (any, error) {
	re := regexp.MustCompile(`mul\((\d{1,3}),(\d{1,3})\)`)
	total := 0
	for line := range iterators.MustLines(input) {
		for _, match := range re.FindAllString(line, -1) {
			stage.Println("match:", match)
			n := iterators.Map(re.FindStringSubmatch(match)[1:], utils.MustAtoi)
			total += n[0] * n[1]
		}
	}
	return total, nil
}

func Stage2(input io.Reader) (any, error) {
	re := regexp.MustCompile(`mul\((\d{1,3}),(\d{1,3})\)|do\(\)|don't\(\)`)
	total := 0
	enabled := true
	for line := range iterators.MustLines(input) {
		for _, match := range re.FindAllString(line, -1) {
			stage.Println("match:", match)
			switch match {
			case `do()`:
				enabled = true
			case `don't()`:
				enabled = false
			default:
				if enabled {
					n := iterators.Map(re.FindStringSubmatch(match)[1:], utils.MustAtoi)
					total += n[0] * n[1]
				}
			}
		}
	}
	return total, nil
}
