package src

import (
	"fmt"
	"testing"
	"unsafe"
)

func TestMapRange(t *testing.T) {
	ids := map[int]string{}
	ids[1] = "a"
	ids[2] = "b"
	for k := range ids {
		t.Log(k)
	}

	for k, v := range ids {
		t.Log(k, v)
	}
}

func TestMapNilLength(t *testing.T) {
	var a = make(map[int32][]string)
	t.Log(len(a))
	t.Log(len(a[1]))

	a[1] = append(a[1], "a")
	t.Log(a)
}

// TestMapDistinct 利用map去重
func TestMapDistinct(t *testing.T) {
	ids := []int{1, 2, 3, 4, 4, 3, 2, 1}
	// 去重
	m := make(map[int]struct{}, len(ids)/2)
	var result []int
	for _, v := range ids {
		if _, ok := m[v]; ok {
			continue
		}
		m[v] = struct{}{}
		result = append(result, v)
	}

	t.Log(result)
}

// TestMapExpend map初始化与扩容
func TestMapExpend(t *testing.T) {
	m := make(map[int]struct{}, 1)
	t.Logf("m len: %d", len(m))
	m[1] = struct{}{}
	m[2] = struct{}{}
	t.Log("----------")
	t.Logf("m len: %d", len(m))
	t.Log(m)
}

// TestMapOK
func TestMapOK(t *testing.T) {
	// var m map[int]struct{}
	m := make(map[int]int, 1)
	m[0] = 0
	m[1] = 1
	v, ok := m[1]
	t.Logf("m[1] value: %d", v)
	t.Logf("m[1] ok: %t", ok) // %t bool类型占位符

	t.Log("----------")
	v2, ok2 := m[2]
	t.Logf("m[2] value: %d", v2)
	t.Logf("m[2] ok: %t", ok2) // %t bool类型占位符

	t.Log("----------")
	v3 := m[3]
	t.Logf("m[2] value: %d", v3)
}

func TestLenMapByPointer(t *testing.T) {
	m := make(map[int32]any, 2)
	m[1] = "1"
	m[2] = "2"

	fmt.Println(len(m))
	fmt.Printf("%p\n", &m)
	// NOTE[Dokiy] 2022/11/9: ??为啥取不出来长度？
	fmt.Println(*(*int)(unsafe.Pointer(&m)))
}

func TestLenChanByPointer(t *testing.T) {
	c := make(chan struct{}, 2)
	c <- struct{}{}
	c <- struct{}{}

	fmt.Println(len(c))
	fmt.Printf("%p\n", &c)
	// NOTE[Dokiy] 2022/11/9: ??为啥取不出来长度？
	fmt.Println(*(*int)(unsafe.Pointer(&c)))
}
