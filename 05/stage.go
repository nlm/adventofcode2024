package main

import (
	"bufio"
	"fmt"
	"io"
	"slices"
	"strings"

	"github.com/nlm/adventofcode2024/internal/stage"
	"github.com/nlm/adventofcode2024/internal/utils"
)

func ReadRequirements(s *bufio.Scanner) map[string][]string {
	requirements := make(map[string][]string, 0)
	for s.Scan() {
		line := s.Text()
		// Section separator
		if strings.TrimSpace(line) == "" {
			stage.Println("break")
			break
		}
		// Read ordering rules
		parts := strings.Split(line, "|")
		if len(parts) != 2 {
			panic(fmt.Errorf("malformed line: %v", line))
		}
		// requiredments[PAGE] = all pages that must already be printed
		requirements[parts[1]] = append(requirements[parts[1]], parts[0])
	}
	stage.Println(requirements)
	return requirements
}

func ReadUpdates(s *bufio.Scanner) [][]string {
	updates := make([][]string, 0)
	for s.Scan() {
		line := s.Text()
		stage.Println(line)
		// Read pages to print
		updates = append(updates, strings.Split(line, ","))
	}
	stage.Println(updates)
	return updates
}

func ReportIsOrdered(pages []string, requirements map[string][]string) bool {
	stage.Println("== printing update:", pages, "==")
	alreadyPrinted := make([]string, 0, len(pages))
	updateOk := true
	for _, page := range pages {
		stage.Println("printing page", page)

		// this page is now printed
		alreadyPrinted = append(alreadyPrinted, page)

		// checkint requirements
		pageOk := true
		requirements, ok := requirements[page]
		if !ok {
			// no requirements found
			stage.Println("no requirements")
			continue
		}
		for _, requiredPage := range requirements {
			stage.Println(" page", page, "requires", requiredPage)
			// requiredment not part of the update
			if !slices.Contains(pages, requiredPage) {
				stage.Println("  which is not part of the update")
				continue
			}
			if !slices.Contains(alreadyPrinted, requiredPage) {
				stage.Println("  which is not yet printed: error")
				pageOk = false
				break
			}
			stage.Println("  requirement ok", requiredPage)
		}
		if pageOk {
			stage.Println("=> page is OK")
		} else {
			updateOk = false
			stage.Println("=> page is BAD")
			break
		}
	}
	if updateOk {
		stage.Println("=> report is OK")
	} else {
		stage.Println("=> report is BAD")
	}
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

func ReorderPages(pages []string, requirements map[string][]string) []string {
	return slices.SortedStableFunc(slices.Values(pages), func(a, b string) int {
		if slices.Contains(requirements[b], a) {
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
		stage.Println("before:", pages)
		pages = ReorderPages(pages, requirements)
		stage.Println("after:", pages)
		middlePage := pages[len(pages)/2]
		total += utils.MustAtoi(middlePage)
	}
	return total, nil
}
