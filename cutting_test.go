package main

import (
	"strings"
	"testing"
)

func TestCutterFromSeperator(t *testing.T) {
	test := func(s, expectedL, expectedR string, findable bool, c StringCutter) {
		actualL, actualR, found := c(s)
		if expectedL != actualL || expectedR != actualR || findable != found {
			t.Errorf("Expected (%s,%s,%v) Got (%s,%s,%v)", expectedL, expectedR, findable, actualL, actualR, found)
		}
	}

	testFindable := func(expectedL, sep, expectedR string) {
		s := expectedL + sep + expectedR
		c := cutterFromSeperator(sep)
		test(s, expectedL, expectedR, true, c)
	}

	testUnfindable := func(s, sep string) {
		c := cutterFromSeperator(sep)
		test(s, s, "", false, c)
	}

	// findable with simple cut
	testFindable("l", " ", "r")
	testFindable("l", "\t", "r")
	testFindable("l", "_", "r")
	testFindable("l", "\n", "right")
	testFindable("l", " ", "r")
	testFindable("l", "  ", "r")
	testFindable("l", "\ts", "r")

	// findable with recurring sepereator
	testFindable("l", " ", " r")
	testFindable("l", " ", "r ")
	testFindable("l", " ", "  ")

	// not findable seperators
	testUnfindable("lr", " ")
	testUnfindable("l r", "\t")
}

func TestMultiWsCutter(t *testing.T) {
	test := func(s, expectedL, expectedR string, findable bool) {
		actualL, actualR, found := multiWsCutter(s)
		if expectedL != actualL || expectedR != actualR || findable != found {
			t.Errorf("Expected (%s,%s,%v) Got (%s,%s,%v) For (%s)", expectedL, expectedR, findable, actualL, actualR, found, s)
		}
	}

	// test if different combination of whitespace result in the same
	test("l  r", "l", "r", true)
	test("l\t r", "l", "r", true)
	test("l \tr", "l", "r", true)
	test("l\t\tr", "l", "r", true)
	test("l   \t   r", "l", "r", true)
	test("longerleft  longerright", "longerleft", "longerright", true)

	// test that normal whitespace is no problem
	test("no double ws\t\tbut now", "no double ws", "but now", true)
	test("a b  c", "a b", "c", true)
	test("a  b c", "a", "b c", true)

	// seperator is at the start or at the end
	test("a b\tc", "a b", "c", true)
	test("  b", "", "b", true)
	test("a  ", "a", "", true)

	// test that it reports if there is no double ws
	test("no double ws", "no double ws", "", false)

	if tabCountsAsMultiWs {
		test("l\tr", "l", "r", true)
		test("a b\tc", "a b", "c", true)
		test("a\tb c", "a", "b c", true)
		test("\tb", "", "b", true)
		test("a\t", "a", "", true)
	}

}

func TestWsSingleCutter(t *testing.T) {
	test := func(s, expectedL, expectedR string, findable bool) {
		actualL, actualR, found := singleWsCutter(s)
		if expectedL != actualL || expectedR != actualR || findable != found {
			t.Errorf("Expected (%s,%s,%v) Got (%s,%s,%v) For (%s)", expectedL, expectedR, findable, actualL, actualR, found, s)
		}
	}

	// test if different combination of whitespace result in the same
	test("l r", "l", "r", true)
	test("l\tr", "l", "r", true)
	test("l  r", "l", "r", true)
	test("l\t r", "l", "r", true)
	test("l \tr", "l", "r", true)
	test("l\t\tr", "l", "r", true)
	test("l   \t   r", "l", "r", true)
	test("longerleft longerright", "longerleft", "longerright", true)

	// seperator is at the start or at the end
	test("a b\tc", "a", "b\tc", true)
	test(" b", "", "b", true)
	test("a  ", "a", "", true)
}

func TestCutWithCutters(t *testing.T) {
	// it is getting joined on '|' as it is easier to compare
	test := func(s, expectedJoined string, cutters []StringCutter) {
		actualJoined := strings.Join(cutWithCutters(s, cutters), "|")
		if expectedJoined != actualJoined {
			t.Errorf("Expected (%s)  Actual (%s)", expectedJoined, actualJoined)
		}
	}

	// tests should be completely independend
	// on the actual cutter used as they should have their own tests

	sCutter := cutterFromSeperator(" ")
	tCutter := cutterFromSeperator("\t")

	test("a b", "a b", nil)                     // empty
	test("a b", "a b", []StringCutter{})        // empty
	test("a b", "a b", []StringCutter{tCutter}) // not found
	test("a b", "a|b", []StringCutter{sCutter}) // found
	test("a b c", "a|b c", []StringCutter{sCutter})
	test("a b c", "a|b|c", []StringCutter{sCutter, sCutter})
	test("a b\tc", "a|b|c", []StringCutter{sCutter, tCutter})
	test("a\tb c", "a|b|c", []StringCutter{tCutter, sCutter})

}

func TestCutAllWithCutter(t *testing.T) {
	test := func(s, expected string, c StringCutter) {
		actual := strings.Join(cutAllWithCutter(s, c), "|")
		if expected != actual {
			t.Errorf("Expected (%s) Got (%s)", expected, actual)
		}
	}

	test("a b", "a b", multiWsCutter)
	test("a  b", "a|b", multiWsCutter)
	test("a  b c", "a|b c", multiWsCutter)
	test("a c\t\tb c", "a c|b c", multiWsCutter)

	test("a b", "a|b", cutterFromSeperator(" "))
	test("a b", "a|b", cutterFromSeperator(" "))
	test("a b", "a b", cutterFromSeperator("\t"))
}

func TestFlagToCutters(t *testing.T) {
	testOk := func(flag, delim string, expectedCutters []StringCutter) {
		actualCutters, err := flagToCutters(flag, delim)

		if err != nil {
			t.Errorf("Did not expect to get an error for flag '%s'", flag)
		}

		testcases := []string{
			"a",
			"a b",
			"a\tb",
			"a  b",
			"a \t b",
			"a\t\tb",
			"abc",
			" a",
			"  a",
			"\ta",
			"\t a",
			"\t\ta",
			"a ",
			"a  ",
			"a\t",
			"a\t\t",
			"a\t b c",
			"a\t b  c",
			"a\t\tb\t\tc",
			"a\tb\tb   c    d     e",
			"ab ab   ab",
		}

		for _, testcase := range testcases {

			// we use cut here which has it's own tests
			actualJoined := strings.Join(cutWithCutters(testcase, actualCutters), "|")
			expectedJoined := strings.Join(cutWithCutters(testcase, expectedCutters), "|")

			if expectedJoined != actualJoined {
				t.Errorf("Expected (%s)  Got (%s)  On (%s) For Flag (%s)", expectedJoined, actualJoined, testcase, flag)
			}
		}
	}

	testFailed := func(flag, sep string) {
		_, err := flagToCutters(flag, sep)
		if err == nil {
			t.Errorf("Expected to get an error from '%s'", flag)
		}
	}

	// single so that we know that the function can get the right cutters
	testOk("s", "|", []StringCutter{cutterFromSeperator(" ")})
	testOk("t", "|", []StringCutter{cutterFromSeperator("\t")})
	testOk("m", "|", []StringCutter{multiWsCutter})
	testOk("<a>", "|", []StringCutter{cutterFromSeperator("a")})

	// combination
	testOk("s|s", "|", []StringCutter{cutterFromSeperator(" "), cutterFromSeperator(" ")})
	testOk("s|t", "|", []StringCutter{cutterFromSeperator(" "), cutterFromSeperator("\t")})
	testOk("s|<a>", "|", []StringCutter{cutterFromSeperator(" "), cutterFromSeperator("a")})

	// current error hanlding on invalid input
	testFailed("s|t", "-")
	testFailed("", "|")
	testFailed("|", "|")
	testFailed("a|", "|")

}
