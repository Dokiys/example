package main

import (
	"fmt"
	"os"
	"runtime/pprof"
	"testing"
)

func fib(n int) int {
	if n <= 1 {
		return 1
	}

	return fib(n-1) + fib(n-2)
}

func TestPprof(t *testing.T) {
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
