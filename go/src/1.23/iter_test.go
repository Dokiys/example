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

func fibonacci(yield func(int) bool) {
	i, j := 1, 0
	for {
		isBreak := yield(i)
		if !isBreak {
			break
		}
		i, j = i+j, i
	}
}

func TestIterFibPush(t *testing.T) {
	for n := range fibonacci {
		if n > 100 {
			break
		}
		println(n)
	}
}

func TestIterFibPull(t *testing.T) {
	next, stop := iter.Pull(fibonacci)
	i, b := next()
	t.Log(i, b)
	i2, b2 := next()
	t.Log(i2, b2)
	i3, b3 := next()
	t.Log(i3, b3)
	i4, b4 := next()
	t.Log(i4, b4)
	for n, ok := next(); ok; n, ok = next() {
		if n > 100 {
			stop()
			break
		}
		println(n)
	}
}
