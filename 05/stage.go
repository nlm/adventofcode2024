package main

import (
	"bufio"
	"fmt"
	"io"
	"slices"
	"strings"

	"github.com/nlm/adventofcode2024/internal/sets"
	"github.com/nlm/adventofcode2024/internal/stage"
	"github.com/nlm/adventofcode2024/internal/utils"
)

func ReadRequirements(s *bufio.Scanner) map[string]sets.Set[string] {
	requirements := make(map[string]sets.Set[string], 0)
	for s.Scan() {
		line := s.Text()

		// Detect section separator
		if strings.TrimSpace(line) == "" {
			break
		}
		// Read ordering rules
		parts := strings.Split(line, "|")
		if len(parts) != 2 {
			panic(fmt.Errorf("malformed line: %v", line))
		}
		// requiredments[PAGE] = all pages that must already be printed
		requirements[parts[1]] = sets.Append(requirements[parts[1]], parts[0])
	}
	stage.Println(requirements)
	return requirements
}

func ReadUpdates(s *bufio.Scanner) [][]string {
	updates := make([][]string, 0)
	for s.Scan() {
		line := s.Text()
		// Read pages to print
		updates = append(updates, strings.Split(line, ","))
	}
	// stage.Println(updates)
	return updates
}

func ReportIsOrdered(pages []string, requirements map[string]sets.Set[string]) bool {
	// stage.Println("== printing update:", pages, "==")
	allPages := sets.Append(nil, pages...)
	alreadyPrinted := make(sets.Set[string], len(pages))
	updateOk := true
	for _, page := range pages {
		// stage.Println("printing page", page)

		// this page is now printed
		alreadyPrinted.Add(page)

		// checkint requirements
		pageOk := true
		requirements, ok := requirements[page]
		if !ok {
			// no requirements found
			// stage.Println("no requirements")
			continue
		}
		for requiredPage := range sets.Values(requirements) {
			// stage.Print("  page ", page, " requires ", requiredPage)
			// requirement not part of the update
			if !allPages.Contains(requiredPage) {
				// stage.Println(" which is not part of the update: ignore")
				continue
			}
			// requirement not yet printed
			if !alreadyPrinted.Contains(requiredPage) {
				// stage.Println(" which is not yet printed: error")
				pageOk = false
				break
			}
			// stage.Println(" which is printed: ok")
		}
		if pageOk {
			// stage.Println("  page is OK")
		} else {
			// stage.Println("  page is BAD")
			updateOk = false
			break
		}
	}
	// if updateOk {
	// 	stage.Println("=> update is OK")
	// } else {
	// 	stage.Println("=> update is BAD")
	// }
	return updateOk
}

func Stage1(input io.Reader) (any, error) {
	s := bufio.NewScanner(input)

	// scan ordering rules
	requirements := ReadRequirements(s)

	// scan pages to print
	updates := ReadUpdates(s)

	// check order
	total := 0
	for _, pages := range updates {
		if ReportIsOrdered(pages, requirements) {
			middlePage := pages[len(pages)/2]
			total += utils.MustAtoi(middlePage)
		}
	}
	return total, nil
}

func ReorderPages(pages []string, requirements map[string]sets.Set[string]) []string {
	return slices.SortedStableFunc(slices.Values(pages), func(a, b string) int {
		if requirements[b].Contains(a) {
			return -1
		}
		return 0
	})
}

func Stage2(input io.Reader) (any, error) {
	s := bufio.NewScanner(input)

	// scan ordering rules
	requirements := ReadRequirements(s)

	// scan pages to print
	updates := ReadUpdates(s)

	// check order
	total := 0
	for _, pages := range updates {
		if ReportIsOrdered(pages, requirements) {
			// report already ok, skip
			continue
		}
		// stage.Println("before:", pages)
		pages = ReorderPages(pages, requirements)
		// stage.Println("after:", pages)
		middlePage := pages[len(pages)/2]
		total += utils.MustAtoi(middlePage)
	}
	return total, nil
}
