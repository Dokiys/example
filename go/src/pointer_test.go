package src

import (
	"encoding/json"
	"fmt"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

// TestPointer 指向指针的指针
func TestPointer(t *testing.T) {
	type A struct{}

	a := &A{}
	b := &a
	t.Logf("type b:%T", b)

	c := &b
	t.Logf("type c:%T", c)

	t.Log("-----")
	t.Logf("pointer a:%p", a)
	t.Logf("pointer b:%p", b)
	t.Logf("pointer c:%p", c)
}

type A struct {
	b *B
}
type B struct {
	i int32
}

// TestPointerPassByPointer 测试指针传递
func TestPointerPassByPointer(t *testing.T) {
	b := &B{i: 1}
	a := &A{b: b}
	fmt.Printf("before pass by pointer \n")
	fmt.Printf("\tfunc a is %p\n", &a)
	fmt.Printf("\tfunc a.b is %v\n", a.b)
	passByPointer(a)
	fmt.Printf("after pass by pointer \n")
	fmt.Printf("\tfunc a is %p\n", &a)
	fmt.Printf("\tfunc a.b is %v\n", a.b)
}

// TestPointerPassByValue 测试值传递
func TestPointerPassByValue(t *testing.T) {
	b := &B{i: 1}
	a := &A{b: b}
	fmt.Printf("before pass by value \n")
	fmt.Printf("\tfunc a is %p\n", &a)
	fmt.Printf("\tfunc a.b is %v\n", a.b)
	passByValue(*a)
	fmt.Printf("after pass by value \n")
	fmt.Printf("\tfunc a is %p\n", &a)
	fmt.Printf("\tfunc a.b is %v\n", a.b)
}

func passByPointer(a *A) {
	b := B{i: 2}
	a.b = &b
}

func passByValue(a A) {
	b := B{i: 2}
	a.b = &b
}

func TestPointerCopy(t *testing.T) {
	b := &B{i: 1}
	a := &A{b: b}
	var copy_a = new(A)
	*copy_a = *a
	assert.True(t, copy_a.b == a.b)
	copy_a = &A{b: nil}
	copy_a.b = &*a.b
	assert.True(t, copy_a != a)
	assert.True(t, copy_a.b == a.b)

	bb := *a.b
	copy_a.b = &bb
	assert.True(t, copy_a != a)
	assert.True(t, copy_a.b != a.b)
}

func TestPointerMultiUnmarshall(t *testing.T) {

	type A struct {
		Name string
	}

	str := "{\"ColumnName\":\"123\"}"

	f := func() (result *A) {
		t.Log(result)
		t.Log(&result)
		var i interface{}
		i = &result
		t.Log(&i)

		err := json.Unmarshal([]byte(str), &i)
		if err != nil {
			t.Fatal(err)
		}

		t.Log(**i.(**A))
		return
	}

	f()
}

func TestPointer_Offset(t *testing.T) {
	var a, b []int
	a = []int{1, 2, 3}
	b = a
	a = a[:2]

	fmt.Printf("%v\n", unsafe.Pointer(&a))
	aLenPtr := uintptr(unsafe.Pointer(&a)) + uintptr(8)
	aCapPtr := uintptr(unsafe.Pointer(&a)) + uintptr(16)
	fmt.Printf("a len: %v\n", (*(*int)(unsafe.Pointer(aLenPtr))))
	fmt.Printf("a cap: %v\n", (*(*int)(unsafe.Pointer(aCapPtr))))

	fmt.Printf("%v\n", unsafe.Pointer(&b))
	bLenPtr := uintptr(unsafe.Pointer(&b)) + uintptr(8)
	bCapPtr := uintptr(unsafe.Pointer(&b)) + uintptr(16)
	fmt.Printf("b len: %v\n", (*(*int)(unsafe.Pointer(bLenPtr))))
	fmt.Printf("b cap: %v\n", (*(*int)(unsafe.Pointer(bCapPtr))))
}
