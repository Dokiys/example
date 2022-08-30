package main

import (
	"fmt"
	"strings"
	"unicode"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"protoc-gen-go-myerrors/v2/proto/errors"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
)

const (
	sqlPackage     = protogen.GoImportPath("database/sql")
	fmtPackage     = protogen.GoImportPath("fmt")
	runtimePackage = protogen.GoImportPath("runtime")
	errorsPackage  = protogen.GoImportPath("github.com/go-kratos/kratos/v2/errors")
)

var enCases = cases.Title(language.AmericanEnglish, cases.NoLower)

func generate(gen *protogen.Plugin, file *protogen.File) *protogen.GeneratedFile {
	if len(file.Enums) == 0 {
		return nil
	}
	filename := file.GeneratedFilenamePrefix + "_replace_errors.go"
	g := gen.NewGeneratedFile(filename, file.GoImportPath)
	g.P("package ", file.GoPackageName)
	g.P()
	g.QualifiedGoIdent(sqlPackage.Ident(""))
	g.QualifiedGoIdent(fmtPackage.Ident(""))
	g.QualifiedGoIdent(runtimePackage.Ident(""))
	g.QualifiedGoIdent(errorsPackage.Ident(""))
	g.P()

	g.P("// This is a compile-time assertion to ensure that this generated file")
	g.P("// is compatible with the kratos package it is being compiled against.")
	g.P("const _ = ", errorsPackage.Ident("SupportPackageIsVersion1"))
	g.P()
	index := 0
	for _, enum := range file.Enums {
		if !genErrorsReason(g, enum) {
			index++
		}
	}
	// If all enums do not contain 'errors.code', the current file is skipped
	if index == 0 {
		g.Skip()
	}
	return g
}

func genErrorsReason(g *protogen.GeneratedFile, enum *protogen.Enum) bool {
	defaultCode := proto.GetExtension(enum.Desc.Options(), errors.E_DefaultCode)
	code := 0
	if ok := defaultCode.(int32); ok != 0 {
		code = int(ok)
	}
	if code > 600 || code < 0 {
		panic(fmt.Sprintf("Enum '%s' range must be greater than 0 and less than or equal to 600", string(enum.Desc.Name())))
	}
	var ew errorWrapper
	for _, v := range enum.Values {
		eCode := proto.GetExtension(v.Desc.Options(), errors.E_Code)
		if ok := eCode.(int32); ok != 0 {
			code = int(ok)
		}
		// If the current enumeration does not contain 'errors.code'
		// or the code value exceeds the range, the current enum will be skipped
		if code > 600 || code < 0 {
			panic(fmt.Sprintf("Enum '%s' range must be greater than 0 and less than or equal to 600", string(v.Desc.Name())))
		}
		if code == 0 {
			continue
		}

		if !strings.HasSuffix(case2Camel(string(v.Desc.Name())), "NotFound") {
			continue
		}
		comment := v.Comments.Leading.String()
		if comment == "" {
			comment = v.Comments.Trailing.String()
		}

		err := &errorInfo{
			Name:       string(enum.Desc.Name()),
			Value:      string(v.Desc.Name()),
			CamelValue: case2Camel(string(v.Desc.Name())),
			HTTPCode:   code,
			Comment:    comment,
			HasComment: len(comment) > 0,
		}
		ew.Errors = append(ew.Errors, err)
	}
	if len(ew.Errors) == 0 {
		return true
	}
	g.P(ew.execute())

	return false
}

func case2Camel(name string) string {
	if !strings.Contains(name, "_") {
		if name == strings.ToUpper(name) {
			name = strings.ToLower(name)
		}
		return enCases.String(name)
	}
	strs := strings.Split(name, "_")
	words := make([]string, 0, len(strs))
	for _, w := range strs {
		hasLower := false
		for _, r := range w {
			if unicode.IsLower(r) {
				hasLower = true
				break
			}
		}
		if !hasLower {
			w = strings.ToLower(w)
		}
		w = enCases.String(w)
		words = append(words, w)
	}

	return strings.Join(words, "")
}
