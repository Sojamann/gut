package main

import "io"

// struct that enables to use the flag
// package with different aliases while
// ensuring that all flags use the same
// value and default value
type arg[T any] struct {
	value        T
	aliases      []string
	defaultValue T
}

// stores indecies which which a slice
// can be accessed. Both left and right
// inclusive and 1-index based, so that 0
// the default value, stands for unspecified
type span struct {
	left, right int
}

// A StringCutter cuts a string into left, right and found.
// It behaves like strings.Cut with the
// seperation token being encapulated in the function
// itself.
type StringCutter func(string) (string, string, bool)

// A StringSplitter cuts a string into multiple pieces.
// Unlike the StringCutter that only performs one cut
// a StringSplitter can perform multiple cuts.
type StringSplitter func(string) []string

// An AutoCloseReader is a reader that encapsulates
// a ReadCloser. The AutoCloseReader closes the
// underlying ReadCloser when a read from it
// returns an error.
type AutoCloseReader struct {
	r io.ReadCloser
}

func (c AutoCloseReader) Read(b []byte) (int, error) {
	n, err := c.r.Read(b)
	if err != nil {
		c.r.Close()
	}

	return n, err
}
