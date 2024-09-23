package __23

import (
	"iter"
	"testing"
)

func integers(data []int) iter.Seq[int] {
	return func(yield func(int) bool) {
		for i := range data {
			if !yield(i) {
				return
			}
		}
	}
}
func even(seq iter.Seq[int]) iter.Seq[int] {
	return func(yield func(int) bool) {
		for n := range seq {
			if n%2 == 0 {
				if !yield(n) {
					return
				}
			}
		}
	}
}
func singleDigit(seq iter.Seq[int]) iter.Seq[int] {
	return func(yield func(int) bool) {
		for n := range seq {
			if n < 10 {
				if !yield(n) {
					return
				}
			}
		}
	}
}

func TestIterSeq(t *testing.T) {
	var data = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
	for p := range singleDigit(even(integers(data))) {
		t.Log(p)
	}
}
