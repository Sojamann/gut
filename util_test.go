package main

import (
	"errors"
	"testing"
)

func TestCountValue(t *testing.T) {
	if actual := countValue(true, true, false); actual != 1 {
		t.Errorf("Expected count to be 1. Actual : %v", actual)
	}

	if actual := countValue(2, 0, 0); actual != 0 {
		t.Errorf("Expected count to be 0. Actual : %v", actual)
	}

	if actual := countValue("A", "A", "A", "AB"); actual != 2 {
		t.Errorf("Expected count to be 2. Actual : %v", actual)
	}
}

type stubFile struct {
	open bool
}

func (s *stubFile) Read(b []byte) (int, error) {
	return 0, errors.New("")
}

func (s *stubFile) Close() error {
	s.open = false
	return nil
}

func TestAutoCloseReader(t *testing.T) {
	f := &stubFile{}
	r := AutoCloseReader{f}

	// this will error but just to make sure
	if _, err := r.Read([]byte{}); err == nil {
		t.Errorf("Expected the stubFile to return an error!!!!!")
	}

	if f.open {
		t.Errorf("The file should have been closed")
	}
}
