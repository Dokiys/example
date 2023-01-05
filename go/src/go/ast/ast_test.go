package ast

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"testing"
)

func TestAstPrint(t *testing.T) {
	src := `
package hello

import "fmt"

type A struct {
	s string
	s1 string
}
func greet() {
	a := &A{
		s: "123",

		s1: "123",
	}
	fmt.Println(a)
}
`
	// Create the AST by parsing src.
	fset := token.NewFileSet() // positions are relative to fset
	f, err := parser.ParseFile(fset, "", src, 0)
	if err != nil {
		panic(err)
	}

	// Print the AST.
	ast.Print(fset, f)
}

func TestFormatSrc(t *testing.T) {
	src := `a := 		1
`
	got, err := format.Source([]byte(src))
	if err != nil {
		panic(err)
	}

	fmt.Println(string(got))
}
