package runtime

import (
	"runtime"
	"testing"
)

func TestRuntimeCaller(t *testing.T) {
	callerWrap(t)
}

func caller(t *testing.T) {
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

func callerWrap(t *testing.T) {
	caller(t)
}

func TestRuntimeCallers(t *testing.T) {
	callersWrap(t)
}

func callersWrap(t *testing.T) {
	callers(t)
}

func callers(t *testing.T) {
	pc := make([]uintptr, 10)
	n := runtime.Callers(0, pc)
	for i := 0; i < n; i++ {
		funcPc := runtime.FuncForPC(pc[i])
		t.Logf("func: [%s]\n", funcPc.Name())
		file, line := funcPc.FileLine(pc[i])
		t.Logf("stack: [%s:%d]\n", file, line)
	}
}
