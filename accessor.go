package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const spanRangeIndicator = ':'

// Access returns the through the span selected
// items from a given slice.
// =0 indicates that there is no bound -> all preceeding/following
// <0 [1 2 3 4][-1] == -4 == [1 2 3 4][4]
// >0 [1 2 3 4][1] == 1
func access[T any](a span, items []T) []T {
	if a.left > len(items) {
		return []T{}
	}

	if a.right > len(items) {
		a.right = len(items)
	}

	if a.left < 0 {
		a.left = (len(items) + a.left) + 1
		if a.left < 0 {
			a.left = 0
		}
	}
	if a.right < 0 {
		a.right = (len(items) + a.right) + 1
		if a.right <= 0 {
			return []T{}
		}
	}

	switch {
	case a.left == 0 && a.right == 0:
		return items
	case a.left == 0 && a.right != 0:
		return items[:a.right]
	case a.left != 0 && a.right == 0:
		return items[a.left-1:]
	case a.left != 0 && a.right != 0:
		return items[a.left-1 : a.right]
	}

	panic("Should never happen")
}

// flagToSpans returns the spans specified by the flag.
// The flag argument can contain multiple spans seperated
// by the specified seperator.
func flagToSpans(flag string, sep string) ([]span, error) {
	// empty means select all, this is the default for empty flag
	if len(flag) == 0 {
		return []span{{}}, nil
	}

	items := strings.Split(flag, sep)

	result := make([]span, 0, 3)
	for _, item := range items {
		if len(item) == 0 {
			return nil, errors.New("input seems to be seperated badly")
		}

		var idx, lDigit, rDigit int
		var err error

		count := strings.Count(item, string(spanRangeIndicator))
		if count > 1 {
			return nil, fmt.Errorf("there are too many '%c' in the '%s'", spanRangeIndicator, item)
		}

		idx = strings.IndexRune(item, spanRangeIndicator)

		if count == 0 {
			lDigit, err = strconv.Atoi(item)
			rDigit = lDigit
			if err != nil {
				return nil, fmt.Errorf("the left part of '%s' does not seem to be an integer", item)
			}
		} else {
			if idx == 0 {
				lDigit = 0
			} else {
				lDigit, err = strconv.Atoi(item[:idx])
				if err != nil {
					return nil, fmt.Errorf("the left part of '%s' does not seem to be an integer", item)
				}
			}

			if idx == len(item)-1 {
				rDigit = 0
			} else {
				rDigit, err = strconv.Atoi(item[idx+1:])
				if err != nil {
					return nil, fmt.Errorf("the right part of '%s' does not seem to be an integer", item)
				}
			}
		}

		// we can only check for the left <= right logic if both have the same sign
		if ((lDigit > 0 && rDigit > 0) || (lDigit < 0 && rDigit < 0)) && lDigit > rDigit {
			return nil, fmt.Errorf("the left number cannot be greater than the right number in '%s'", item)
		}
		result = append(result, span{left: lDigit, right: rDigit})
	}

	return result, nil
}
