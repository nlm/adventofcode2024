package main

import (
	"io"
	"math"

	"github.com/nlm/adventofcode2024/internal/iterators"
	"github.com/nlm/adventofcode2024/internal/stage"
	"github.com/nlm/adventofcode2024/internal/utils"
)

func Mix(secnum, v int64) int64 {
	return secnum ^ v
}

func Prune(secnum int64) int64 {
	return ((secnum % 16777216) + 16777216) % 16777216
}

func GenerateNext(secnum int64) int64 {
	secnum = Prune(Mix(secnum, secnum*64))
	secnum = Prune(Mix(secnum, int64(math.Floor(float64(secnum)/32))))
	secnum = Prune(Mix(secnum, secnum*2048))
	return secnum
}

func Stage1(input io.Reader) (any, error) {
	stage.Println(Mix(42, 15) == 37)
	stage.Println(Prune(100000000) == 16113920)
	stage.Println(GenerateNext(123) == 15887950)
	total := int64(0)
	for line := range iterators.MustLines(input) {
		number := int64(utils.MustAtoi(line))
		initial := number
		for range 2000 {
			number = GenerateNext(number)
		}
		stage.Println(initial, "->", number)
		total += number
	}
	return total, nil
}

func CalcOneDigit(n int64) int64 {
	return n % 10
}

func Handle1(number int64) map[[4]int64]int64 {
	diffs := make([]int64, 0, 2000)
	wins := make(map[[4]int64]int64)
	for range 2000 {
		next := GenerateNext(number)
		diff := CalcOneDigit(next) - CalcOneDigit(number)
		if len(diffs) >= 4 {
			key := [4]int64(diffs[len(diffs)-4:])
			if _, ok := wins[key]; !ok {
				wins[key] = CalcOneDigit(number)
				stage.Println("wins", key, "->", wins[key])
			}
		}
		diffs = append(diffs, diff)
		number = next
	}
	return wins
}

func Stage2(input io.Reader) (any, error) {
	stage.Println(CalcOneDigit(15887950) == 0)
	stage.Println(CalcOneDigit(16495136) == 6)

	allWins := make(map[[4]int64]int64)
	for line := range iterators.MustLines(input) {
		number := int64(utils.MustAtoi(line))
		wins := Handle1(number)
		stage.Println(number, "->", wins[[4]int64{-2, 1, -1, 3}])
		for k, v := range wins {
			allWins[k] += v
		}
	}

	maxBananas := int64(0)
	for k, v := range allWins {
		if v > maxBananas {
			stage.Println("new max", k, v)
			maxBananas = v
		}
	}

	return maxBananas, nil
}
