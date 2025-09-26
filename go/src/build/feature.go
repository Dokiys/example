//go:build !mymodule_coolfeature

package main

import "fmt"

func Greet() {
	fmt.Println("Hello, old world!")
}
