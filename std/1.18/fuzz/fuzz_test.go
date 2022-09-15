package fuzz

import (
	"testing"
)

func add(a, b int) int {
	return a + b
}

// go test -fuzz=Fuzz
func FuzzAdd(f *testing.F) {
	f.Fuzz(func(t *testing.T, a int, b int) {
		if add(a, b) == 120 {
			t.Error("Error 120")
		}
	})
}
