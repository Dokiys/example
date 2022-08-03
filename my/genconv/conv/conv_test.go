package conv

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"testing"

	"go_test/my/genconv/data"
	"golang.org/x/tools/go/packages"
)

func TestName(t *testing.T) {
	src, _ := ioutil.ReadFile("../temp/template.go")
	fset := token.NewFileSet()
	astf, err := parser.ParseFile(fset, "", string(src), 0)
	if err != nil {
		t.Error(err)
	}
	t.Log(astf)
}

func TestName2(t *testing.T) {
	src, _ := ioutil.ReadFile("../temp/template2.go")
	fset := token.NewFileSet()
	astf, err := parser.ParseFile(fset, "", string(src), 0)
	if err != nil {
		t.Error(err)
	}
	t.Log(astf)
	t.Log(data.A{})
}

// TODO[Dokiy] 2022/8/3: 加载到pacakge中获取类型，并反射创建实例
func TestPackage(t *testing.T) {
	cfg := &packages.Config{
		Mode: packages.LoadImports | packages.LoadSyntax,
		// TODO: Need to think about constants in test files. Maybe write type_string_test.go
		// in a separate pass? For later.
		//Tests: false,
		//BuildFlags: []string{fmt.Sprintf("-tags=%s", strings.Join(tags, " "))},
	}
	pkgs, err := packages.Load(cfg, "go_test/my/genconv/temp")
	if err != nil {
		t.Error(err)
	}

	//ast.Inspect(pkgs[0].Syntax[0], inspect)

	//var obj types.Object
	//obj.Type().Underlying().(*types.Struct).NumFields()

	// NOTE[Dokiy] 2022/8/3: 获取入参类型
	//pkgs[0].Syntax[1].Decls[0].Type.Params.List.Type...
	// NOTE[Dokiy] 2022/8/3: 去获取参数
	// pkgs[0].Imports["go_test/my/gencov/data"].Syntax.Tok == Type
	// pkgs[0].Imports["go_test/my/gencov/data"].Syntax.Speces.Type.(StructType).Fields
	defs := pkgs[0].TypesInfo.Defs
	t.Log(defs)
}

func inspect(node ast.Node) bool {
	decl, ok := node.(*ast.GenDecl)
	if !ok || decl.Tok != token.IMPORT {
		return true
	}

	for _, spec := range decl.Specs {
		fmt.Println(spec)
	}
	return false
}
