package generics

import (
	"testing"
)

type Int int
type INT interface {
	int | int32 | int64 | Int
}

type A[N INT] struct {
	n N
}

func add[T INT](a, b T) T {
	return a + b
}

func TestGenerics(t *testing.T) {
	var i1, i2 int = 1, 2
	var i32_1, i32_2 int32 = 1, 2
	var int1, int2 Int = 1, 2
	t.Log(add(i1, i2))
	t.Log(add(int1, int2))
	t.Log(add(i32_1, i32_2))
	// t.Log(add(i1, i32_2)) invalid

	var a A[int32]
	t.Logf("a.N type: %T", a.n)
}
