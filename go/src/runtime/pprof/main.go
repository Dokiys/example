package main

import (
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
)

func fib(n int) int {
	if n <= 1 {
		return 1
	}

	return fib(n-1) + fib(n-2)
}

func main() {
	if os.Args[1] == "cup" {
		PprofCup()
	} else if os.Args[1] == "heap" {
		PprofHeap()
	}
}

func PprofCup() {
	f, _ := os.OpenFile("fib.profile", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	defer f.Close()
	if err := pprof.StartCPUProfile(f); err != nil {
		fmt.Printf("could not start CPU profile: %s", err)
	}
	defer pprof.StopCPUProfile()

	n := 10
	for i := 1; i <= 5; i++ {
		fmt.Printf("fib(%d)=%d\n", n, fib(n))
		n += 3 * i
	}
}

var s = make([]int, 0)
var s1 []int

func fib2(n int) int {
	s = append(s, n)
	s1 = append(s1, n)
	if n <= 1 {
		return n
	}

	return fib2(n-1) + fib2(n-2)
}

func PprofHeap() {
	s1 = make([]int, 0, 29954739)

	n := 5
	for i := 1; i <= 5; i++ {
		fmt.Printf("fib2(%d)=%d\n", n, fib2(n))
		n += 3 * i
	}

	println("cap(s): \t", cap(s))
	println("cap(s1): \t", cap(s1))
	runtime.GC() // get up-to-date statistics
	f, _ := os.OpenFile("fib_heap.profile", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	defer f.Close()
	if err := pprof.WriteHeapProfile(f); err != nil {
		fmt.Printf("could not write heap profile: %s", err)
	}
}
