package generics

import "testing"

type ModelBasicType interface {
	~int | ~int32 | ~int64 | float64 | ~string
}

func SetOrInitBasic[T ModelBasicType](a *T, b T) {
	var s T
	if s == b {
		return
	}

	*a = b
}

func SetOrInitArr[T ModelBasicType](a *[]T, b []T) {
	*a = make([]T, 0)
	if len(b) <= 0 {
		return
	}

	*a = b
}

func TestSetOrInit(t *testing.T) {
	var a, b int
	b = 2
	SetOrInitBasic[int](&a, b)
	t.Log(a)
}

func TestSetOrInitArr(t *testing.T) {
	var a, b []int
	b = []int{1, 2, 3}
	SetOrInitArr[int](&a, b)
	t.Log(a)
}
