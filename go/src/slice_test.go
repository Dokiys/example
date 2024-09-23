package src

import (
	"slices"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestSliceRange range两个参数遍历
func TestSliceRange(t *testing.T) {
	s := []string{"a", "b", "c"}

	for i, v := range s { // index, value
		t.Log(i)
		t.Log(v)
	}

	for i := range s { // index
		t.Log(i) // 0, 1, 2
	}
}

// TestSliceRangePointer 遍历切片时，v 是同一个指针
func TestSliceRangePointer(t *testing.T) {
	type A struct{ Id int }
	slice := []A{{Id: 1}, {Id: 2}}

	for _, v := range slice {
		t.Logf("%p", &v)
	}
}

// TestSliceRangePointer2 遍历切片时，v 是同一个指针
func TestSliceRangePointer2(t *testing.T) {
	slice := []int{0, 1, 2, 3}
	myMap := make(map[int]*int)

	var index, value int
	for index, value = range slice {
		myMap[index] = &value
	}
	t.Log(myMap)
}

// TestSliceExpend 切片扩容
func TestSliceExpend(t *testing.T) {
	nums := make([]int32, 0, 2)

	t.Logf("len: %d", len(nums))
	t.Logf("cap: %d", cap(nums))
	t.Logf("--------")

	nums = append(nums, 1)
	t.Logf("len: %d", len(nums))
	t.Logf("cap: %d", cap(nums))
	t.Logf("--------")

	nums = append(nums, 2)
	t.Logf("len: %d", len(nums))
	t.Logf("cap: %d", cap(nums))
	t.Logf("--------")

	nums = append(nums, 3)
	t.Logf("len: %d", len(nums))
	t.Logf("cap: %d", cap(nums))
	t.Logf("--------")
}

// TestSortSlice 切片排序
func TestSortSlice(t *testing.T) {
	nums := []int{1, 2, 3, 4, 3, 2, 1}
	sort.Slice(nums, func(i, j int) bool {
		return nums[i] < nums[j]
	})

	t.Log(nums) // [1 1 2 2 3 3 4]
}

// TestSortSlice 切片排序
func TestSliceSort(t *testing.T) {
	type A struct {
		id   int
		name string
	}

	aList := []A{{id: 1, name: "a1"}, {id: 3, name: "a3"}, {id: 2, name: "a2"}}
	var expected string
	for i := 0; i < 3; i++ {
		// sort.Slice(aList, func(i, j int) bool {
		// 	return aList[i].id == 1
		// })
		slices.SortStableFunc(aList, func(a, b A) int {
			switch {
			case a.id == 1 && b.id != 1:
				return 1
			case a.id == 1 && b.id == 1:
				return 0
			case a.id != 1 && b.id == 1:
				return -1
			case a.id != 1 && b.id != 1:
				return 0
			}
			return 0
		})
		var actual string
		for _, a := range aList {
			actual += a.name
		}
		if expected == "" {
			expected = actual
		}
		if expected != actual {
			t.Errorf("expected: %s, actual: %s", expected, actual)
		}
		expected = actual
	}
}

// Insert向 index 左边插入元素, index从0开始
func insert[T any](arr []T, e T, index int) []T {
	if index >= len(arr) {
		panic("array.Insert index out of range")
	}

	r := make([]T, len(arr)+1)
	copy(r[:index], arr[:index])
	r[index] = e
	copy(r[index+1:], arr[index:])

	return r
}

func TestSliceInsert(t *testing.T) {
	arr := []string{"1", "2", "3"}
	arr = insert(arr, "a", 2)
	assert.Equal(t, []string{"1", "2", "a", "3"}, arr)
}
