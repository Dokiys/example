package main

import (
	"fmt"
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/packages"
)

func main() {
	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedTypes | packages.NeedTypesInfo | packages.NeedSyntax | packages.NeedImports | packages.NeedDeps,
		// TODO: Need to think about constants in test files. Maybe write type_string_test.go
		// in a separate pass? For later.
		Tests: false,
		//BuildFlags: []string{fmt.Sprintf("-tags=%s", strings.Join(tags, " "))},
	}
	pkgs, err := packages.Load(cfg, ".")
	if err != nil {
		panic(err)
	}

	ast.Inspect(pkgs[0].Syntax[0], inspect)
}

func inspect(node ast.Node) bool {
	decl, ok := node.(*ast.GenDecl)
	fmt.Println(node)
	if !ok || decl.Tok != token.STRUCT {
		return true
	}

	for _, spec := range decl.Specs {
		fmt.Println(spec)
	}
	return false
}
