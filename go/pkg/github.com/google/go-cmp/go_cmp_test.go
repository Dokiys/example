package go_cmp

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCmp(t *testing.T) {
	type Student struct {
		Id   int
		Name string
		Age  int
	}
	s1 := Student{
		Id:   1,
		Name: "zhangsan",
		Age:  1,
	}

	s2 := Student{
		Id:   2,
		Name: "lisi",
		Age:  2,
	}
	diff := cmp.Diff(s1, s2)
	t.Log(diff)
}
