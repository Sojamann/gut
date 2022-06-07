package main

import (
	"strings"
	"testing"
)

func TestAccess(t *testing.T) {
	items := []string{"T", "e", "s", "t"}

	check := func(a span, expected string) {
		actual := strings.Join(access(a, items), "")
		if actual != expected {
			t.Errorf("%v: Expected (%s) Actual (%s)", a, expected, actual)
		}
	}

	check(span{}, "Test")

	check(span{left: 1}, "Test")
	check(span{left: 2}, "est")
	check(span{left: 3}, "st")
	check(span{left: 4}, "t")
	check(span{left: 5}, "")

	check(span{left: -1}, "t")
	check(span{left: -2}, "st")
	check(span{left: -3}, "est")
	check(span{left: -4}, "Test")
	check(span{left: -5}, "Test")

	check(span{right: 1}, "T")
	check(span{right: 2}, "Te")
	check(span{right: 3}, "Tes")
	check(span{right: 4}, "Test")
	check(span{right: 5}, "Test")
	check(span{right: 100}, "Test")

	check(span{right: -1}, "Test")
	check(span{right: -2}, "Tes")
	check(span{right: -3}, "Te")
	check(span{right: -4}, "T")
	check(span{right: -5}, "")

	check(span{left: 1, right: 1}, "T")
	check(span{left: 1, right: 2}, "Te")
	check(span{left: 1, right: 3}, "Tes")
	check(span{left: 1, right: 4}, "Test")
	check(span{left: 1, right: 5}, "Test")
	check(span{left: 1, right: 100}, "Test")

	check(span{left: 1, right: -1}, "Test")
	check(span{left: 1, right: -2}, "Tes")
	check(span{left: 1, right: -3}, "Te")
	check(span{left: 1, right: -4}, "T")
	check(span{left: 1, right: -5}, "")
	check(span{left: 1, right: -100}, "")

	check(span{left: -4, right: 4}, "Test")
	check(span{left: -3, right: 4}, "est")
	check(span{left: -2, right: 4}, "st")
	check(span{left: -1, right: 4}, "t")
	check(span{left: 0, right: 4}, "Test")

	check(span{left: -2, right: 3}, "s")

	check(span{left: 1, right: 1}, "T")
	check(span{left: 2, right: 2}, "e")
	check(span{left: 3, right: 3}, "s")
	check(span{left: 4, right: 4}, "t")
	check(span{left: -4, right: -4}, "T")
	check(span{left: -3, right: -3}, "e")
	check(span{left: -2, right: -2}, "s")
	check(span{left: -1, right: -1}, "t")

	check(span{left: 1, right: 3}, "Tes")
	check(span{left: 2, right: 3}, "es")
	check(span{left: 3, right: 5}, "st")
	check(span{left: -2, right: 5}, "st")
	check(span{left: -2, right: -1}, "st")
	check(span{left: -100, right: 100}, "Test")
}

func TestFlagToSpans(t *testing.T) {
	testOk := func(flag, sep string, expectedSpans []span) {
		actualsSpans, err := flagToSpans(flag, sep)

		if err != nil {
			t.Errorf("There should not be an error converting from (%s) but got (%v)", flag, err)
		}

		if len(actualsSpans) != len(expectedSpans) {
			t.Errorf("Expected %d spans but got %d From (%s)", len(expectedSpans), len(actualsSpans), flag)
		}

		for i := 0; i < len(expectedSpans); i++ {
			if actualsSpans[i] != expectedSpans[i] {
				t.Errorf("Expected (%v) Got (%v) From (%s)", expectedSpans[i], actualsSpans[i], flag)
			}
		}
	}

	testFailed := func(flag, sep string) {
		_, err := flagToSpans(flag, sep)

		if err == nil {
			t.Errorf("There should be an error converting from (%s) but got (%v)", flag, err)
		}
	}

	// single
	testOk("", ",", []span{{}})

	testOk("1", ",", []span{{left: 1, right: 1}})
	testOk("10", ",", []span{{left: 10, right: 10}})
	testOk("-1", ",", []span{{left: -1, right: -1}})

	testOk("1:", ",", []span{{left: 1}})
	testOk("-1:", ",", []span{{left: -1}})
	testOk(":2", ",", []span{{right: 2}})
	testOk(":-2", ",", []span{{right: -2}})
	testOk("1:2", ",", []span{{left: 1, right: 2}})
	testOk("-3:-2", ",", []span{{left: -3, right: -2}})

	testOk("1:,1:", ",", []span{{left: 1}, {left: 1}})
	testOk("-1:,1:", ",", []span{{left: -1}, {left: 1}})
	testOk(":2,1:", ",", []span{{right: 2}, {left: 1}})
	testOk(":-2,1:", ",", []span{{right: -2}, {left: 1}})
	testOk("1:2,1:", ",", []span{{left: 1, right: 2}, {left: 1}})
	testOk("1:-2,1:", ",", []span{{left: 1, right: -2}, {left: 1}})

	testFailed(",", ",")
	testFailed("1:2,", ",")
	testFailed(",1:2", ",")
	testFailed("1|2,", ",")
	testFailed("::2,", ",")
	testFailed("1::2,", ",")
	testFailed("4:2,", ",")
	testFailed("-2:-4,", ",")
}
