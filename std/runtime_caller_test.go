package std

import (
	"runtime"
	"testing"
)

func s1(t *testing.T) {
	pc1, file1, line1, ok1 := runtime.Caller(0)
	if !ok1 {
		t.Logf("Out of range！")
		return
	}
	t.Logf("stack: [%s:%d]\n", file1, line1)
	funcPc1 := runtime.FuncForPC(pc1)
	t.Logf("func: [%s]\n", funcPc1.Name())

	//
	pc, file, line, ok := runtime.Caller(2)
	if !ok {
		t.Logf("Out of range！")
		return
	}
	t.Logf("stack: [%s:%d]\n", file, line)
	funcPc := runtime.FuncForPC(pc)
	t.Logf("func: [%s]\n", funcPc.Name())

	// common parent path
	var split [2]int
	for i := 0; i < len(file); i++ {
		if file[i] == '/' {
			split[0], split[1] = i, split[0]
		}
		if file[i] != file1[i] {
			break
		}
	}

	t.Logf("%s:%d\n", file[split[1]:], line)
}

func s2(t *testing.T) {
	s1(t)
}
func TestRuntimeCaller(t *testing.T) {
	s2(t)
}
