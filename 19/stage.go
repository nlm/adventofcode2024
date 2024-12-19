package main

import (
	"io"
	"strings"

	"github.com/nlm/adventofcode2024/internal/stage"
)

type Data struct {
	Patterns []string
	Designs  []string
}

func ParseInput(input io.Reader) Data {
	d := Data{}
	lines, _ := io.ReadAll(input)
	parts := strings.SplitN(string(lines), "\n\n", 2)
	stage.Println(parts)
	for _, p := range strings.Split(string(parts[0]), ", ") {
		if strings.TrimSpace(p) != "" {
			d.Patterns = append(d.Patterns, strings.TrimSpace(p))
		}
	}
	for _, ds := range strings.Split(string(parts[1]), "\n") {
		if strings.TrimSpace(ds) != "" {
			d.Designs = append(d.Designs, strings.TrimSpace(ds))
		}
	}
	return d
}

func (dc *DesignCounter) CountPossibleDesigns(design string) int {
	if v, ok := dc.Cache[design]; ok {
		return v
	}
	if len(design) == 0 {
		return 1
	}
	count := 0
	for _, p := range dc.Patterns {
		if len(p) > len(design) {
			continue
		}
		if strings.HasPrefix(design, p) {
			count += dc.CountPossibleDesigns(design[len(p):])
		}
	}
	dc.Cache[design] = count
	return count
}

type DesignCounter struct {
	Cache    map[string]int
	Patterns []string
}

func Stage1(input io.Reader) (any, error) {
	data := ParseInput(input)
	dc := DesignCounter{
		Cache:    make(map[string]int),
		Patterns: data.Patterns,
	}
	total := 0
	for _, ds := range data.Designs {
		if dc.CountPossibleDesigns(ds) > 0 {
			total++
		}
	}
	return total, nil
}

func Stage2(input io.Reader) (any, error) {
	data := ParseInput(input)
	dc := DesignCounter{
		Cache:    make(map[string]int),
		Patterns: data.Patterns,
	}
	total := 0
	for _, ds := range data.Designs {
		total += dc.CountPossibleDesigns(ds)
	}
	return total, nil
}
