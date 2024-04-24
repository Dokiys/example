package generics

import (
	"reflect"
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

func conv[S comparable, T any](tList []T, dbTag string) map[S]T {
	var result = map[S]T{}
	if len(tList) <= 0 {
		return result
	}

	var idx = -1
	for _, t := range tList {
		tValue := reflect.ValueOf(t)
		if idx == -1 {
			if tValue.Type().Kind() != reflect.Struct {
				panic("LocalCache getKey: provided value is not a struct")
			}

			for j := 0; j < tValue.NumField(); j++ {
				if tag, ok := tValue.Type().Field(j).Tag.Lookup("db"); ok && tag == dbTag {
					if tValue.Field(j).Type().AssignableTo(reflect.TypeOf(*new(S))) {
						idx = j
						break
					} else {
						panic("LocalCache getKey: cannot assign field to type S")
					}
				}
			}
			panic("LocalCache getKey: no field with db tag '" + dbTag + "' found")
		}
		result[tValue.Field(idx).Interface().(S)] = t
	}

	return result
}
