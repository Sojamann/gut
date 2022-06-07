package main

import (
	"fmt"
	"strings"
)

// wsChars contains all runes/characters that
// are considered whitespace characters to split on
const wsChars = "\t "

// the TESTS need to know about this
const tabCountsAsMultiWs = true

// cutterToSplitter takes a StringCutter and turns it into
// a StringSplitter by returning a new function where the
// StringCutter is applied as long as possible on a provided string.
func cutterToSplitter(c StringCutter) StringSplitter {
	return func(s string) []string {
		return cutAllWithCutter(s, c)
	}
}

// A StringCutter is created from a seperator as long
// as the seperator is not empty.
func cutterFromSeperator(sep string) StringCutter {
	if len(sep) == 0 {
		panic("cutterFromSeperator should never be used with empty seperators")
	}
	return func(s string) (string, string, bool) {
		return strings.Cut(s, sep)
	}
}

// A StringCutter that cuts the string into before and after when finding
// more then one consecutive whitespace characters.
// Note that all consecutive whitespaces are consumed and not only two.
func multiWsCutter(s string) (string, string, bool) {
	var (
		offset, firstWsIdx, secondWsIdx int
	)
	for offset < len(s) {
		firstWsIdx = strings.IndexAny(s[offset:], wsChars)
		secondWsIdx = strings.IndexAny(s[offset+firstWsIdx+1:], wsChars)

		// these things need to be checked in order

		if firstWsIdx < 0 {
			goto notFound
		}

		if tabCountsAsMultiWs && s[offset+firstWsIdx] == '\t' {
			goto found
		}

		if secondWsIdx < 0 {
			goto notFound
		}

		// the secondWsIdx can only be non-negative if the firstWsIdx is also non-negative
		// as the search for the second ws starts right after the first ws was found
		// a 0 idx indicates that the item was found right afterwards.
		// If the first one is a tab and this behavior is allowed, then we can also end it there
		if secondWsIdx == 0 {
			goto found
		}

		// advance search marker to the location of the second whitespace
		// and consider this the first ws as the next one might be right
		// after it
		offset += firstWsIdx + secondWsIdx
	}
	// this needs to be the first case after the loop when it ends
notFound:
	return s, "", false

found:
	//fmt.Println(s, offset, offset+firstWsIdx, offset+firstWsIdx+secondWsIdx+1)
	// Trim any additional leading ws from the right string as there could be more than two.
	return s[:offset+firstWsIdx], strings.TrimLeft(s[offset+firstWsIdx+1:], wsChars), true
}

// A StringCutter that cuts on any whitespace character.
// Note that all consecutive whitespace it consumed.
func singleWsCutter(s string) (string, string, bool) {
	idx := strings.IndexAny(s, wsChars)
	if idx < 0 {
		return s, "", false
	}

	return s[:idx], strings.TrimLeft(s[idx+1:], wsChars), true
}

// Cuts string using one cutter after another.
// The result is like strings.Split using different
// seperators after each cut. The result of this
// are the parts that the cutters cut the string into.
func cutWithCutters(s string, cutters []StringCutter) []string {
	var (
		match  string
		found  bool
		result = make([]string, 0, 2)
	)

	for _, c := range cutters {
		match, s, found = c(s)

		result = append(result, match)

		if !found {
			return result
		}
	}

	return append(result, s)
}

// Apply a cutter on a given string until the cutter is done.
// Return the parts that the cutter cut the string into.
func cutAllWithCutter(s string, c StringCutter) []string {
	var (
		match  string
		found  bool = true
		result      = make([]string, 0, 3)
	)
	for found {
		match, s, found = c(s)
		result = append(result, match)
	}

	return result
}

// predefinedCutters defines special cutters which
// can be used in the cut on flag command
var predefinedCutters = map[string]StringCutter{
	"t": cutterFromSeperator("\t"),
	"s": cutterFromSeperator(" "),
	"a": singleWsCutter,
	"m": multiWsCutter,
}

// Converts the user supplied cut specification into StringCutters
func flagToCutters(s string, sep string) ([]StringCutter, error) {
	result := make([]StringCutter, 0, 3)
	for _, item := range strings.Split(s, sep) {
		item = strings.Trim(item, wsChars)

		if strings.HasPrefix(item, "<") && strings.HasSuffix(item, ">") && len(item) > 2 {
			result = append(result, cutterFromSeperator(item[1:len(item)-1]))
			continue
		}

		if c, found := predefinedCutters[item]; found {
			result = append(result, c)
			continue
		}

		return nil, fmt.Errorf("unknown delimiter '%s'", item)
	}

	return result, nil
}
