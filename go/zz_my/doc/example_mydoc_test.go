package doc_test

import (
	"fmt"
)

func ExampleVersion() {
	Version()
	// Output:
	// v1.0.0
}

func ExampleStu_Study() {
	stu := Stu{Name: "zhangsan"}
	fmt.Println(stu.Code())
	stu.Study()
	fmt.Println(stu.Code())
	// Output:
	// 0
	// 1
}

func ExampleMeet() {
	stu := Stu{Name: "zhangsan"}
	teacher := Teacher{Name: "lisi"}

	Meet(teacher, stu)
	// Output:
	// lisi: Hello zhangsan
	// zhangsan: Hello Mr.lisi
}
