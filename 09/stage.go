package main

import (
	"fmt"
	"io"
	"slices"
	"strings"

	"github.com/nlm/adventofcode2024/internal/iterators"
	"github.com/nlm/adventofcode2024/internal/stage"
	"github.com/nlm/adventofcode2024/internal/utils"
)

func ExpandIntLineToDiskMap(intLine []int) []int {
	expandedLine := make([]int, 0)
	id := 0
	for i := range intLine {
		isSpace := i%2 == 1
		if !isSpace {
			for range intLine[i] {
				expandedLine = append(expandedLine, id)
			}
			id++
		} else {
			for range intLine[i] {
				expandedLine = append(expandedLine, -1)
			}
		}
	}
	return expandedLine
}

func DiskMapString(diskMap []int) string {
	var nums = []byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}
	b := strings.Builder{}
	for _, v := range diskMap {
		if IsFreeBlock(v) {
			b.WriteByte('.')
		} else {
			b.WriteByte(nums[v%len(nums)])
		}
	}
	return b.String()
}

func IsFreeBlock(n int) bool {
	return n < 0
}

func CompactDiskMap(diskMap []int) {
	var minIdx, maxIdx = 0, len(diskMap) - 1
	fwdCursor, revCursor := minIdx, maxIdx
	// for !IsFreeBlock(diskMap[fwdCursor]) && fwdCursor <= maxIdx {
	// 	fwdCursor++
	// }
	// for IsFreeBlock(diskMap[revCursor]) && revCursor >= maxIdx {
	// 	revCursor--
	// }
	for {
		// place fwdCursor on a free block
		for fwdCursor <= maxIdx && !IsFreeBlock(diskMap[fwdCursor]) {
			fwdCursor++
		}
		// place revCursor on a data block
		for revCursor >= minIdx && IsFreeBlock(diskMap[revCursor]) {
			revCursor--
		}
		if fwdCursor >= revCursor {
			break
		}
		// stage.Println("fwd", fwdCursor, "->", diskMap[fwdCursor], "rev", revCursor, "->", diskMap[revCursor])
		diskMap[fwdCursor] = diskMap[revCursor]
		diskMap[revCursor] = -1
		// stage.Println(DiskMapString(diskMap))
		fwdCursor++
		revCursor--
	}
}

func CheckSum(diskMap []int) int {
	ckSum := 0
	for i, n := range diskMap {
		if n < 0 {
			break
		}
		ckSum += i * n
	}
	return ckSum
}

// func FindLastBlock(diskMap []int) (int, int) {
// 	i := len(diskMap) - 1
// 	for i > 0 && diskMap[i] < 0 {
// 		i--
// 	}
// 	value := diskMap[i]
// 	length := 0
// 	for i > 0 && diskMap[i] == value {
// 		i--
// 		length++
// 	}
// 	return i + 1, length
// }

// func FindFreeSpace(diskMap []int, size int) (int, int) {
// 	for {

// 	}
// }

func Stage1(input io.Reader) (any, error) {
	// Parse Input
	line := slices.Collect(iterators.MustLinesBytes(input))[0]
	intLine := make([]int, 0, len(line))
	for _, c := range line {
		intLine = append(intLine, utils.MustAtoi(string(c)))
	}
	diskMap := ExpandIntLineToDiskMap(intLine)
	// stage.Println(diskMap)
	// stage.Println(DiskMapString(diskMap))
	CompactDiskMap(diskMap)
	// stage.Println(diskMap)
	// stage.Println(DiskMapString(diskMap))
	return CheckSum(diskMap), nil
}

func FindBlock(diskMap []int, value int) (int, int, bool) {
	i := len(diskMap) - 1
	for ; i >= 0; i-- {
		if diskMap[i] > value || diskMap[i] < 0 {
			continue
		}
		break
	}
	length := 0
	for ; i >= 0 && diskMap[i] == value; i-- {
		length++
	}
	if length > 0 {
		return i + 1, length, true
	}
	return 0, 0, false
}

func FindFreeSpace(diskMap []int, size int, boundary int) (int, bool) {
	length := 0
	for i := 0; i < boundary; i++ {
		if diskMap[i] >= 0 {
			length = 0
			continue
		}
		length++
		if length >= size {
			return i - length + 1, true
		}
	}
	return 0, false
}

func FindLastId(diskMap []int) int {
	for i := len(diskMap) - 1; i >= 0; i-- {
		if diskMap[i] > 0 {
			return diskMap[i]
		}
	}
	return -1
}

func CheckSum2(diskMap []int) int {
	ckSum := 0
	for i, n := range diskMap {
		if n < 0 {
			continue
		}
		ckSum += i * n
	}
	return ckSum
}

func Stage2(input io.Reader) (any, error) {
	// Parse Input
	line := slices.Collect(iterators.MustLinesBytes(input))[0]
	intLine := make([]int, 0, len(line))
	for _, c := range line {
		intLine = append(intLine, utils.MustAtoi(string(c)))
	}
	diskMap := ExpandIntLineToDiskMap(intLine)
	// stage.Println(DiskMapString(diskMap))
	lastId := FindLastId(diskMap)
	for id := lastId; id >= 0; id-- {
		stage.Println("---", id, "---")
		i, l, ok := FindBlock(diskMap, id)
		if !ok {
			return 0, fmt.Errorf("block %d not found", id)
		}
		// stage.Println("block", diskMap[i : i+l])
		j, ok := FindFreeSpace(diskMap, l, i) //len(diskMap)-1)
		if !ok {
			stage.Println("no free space")
			continue
		}
		// stage.Println(DiskMapString(diskMap))
		copy(diskMap[j:j+l], diskMap[i:i+l])
		// stage.Println(DiskMapString(diskMap))
		for k := i; k < i+l; k++ {
			diskMap[k] = -1
		}
		// stage.Println(DiskMapString(diskMap))
	}
	return CheckSum2(diskMap), nil
}
