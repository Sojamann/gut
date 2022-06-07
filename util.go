package main

import (
	"fmt"
	"os"
)

func die(format string, stuff ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", stuff...)
	os.Exit(1)
}

func countValue[T comparable](value T, items ...T) int {
	counter := 0
	for _, item := range items {
		if item == value {
			counter += 1
		}
	}
	return counter
}
