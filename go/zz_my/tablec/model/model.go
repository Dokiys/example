package model

import (
	"bytes"
	_ "embed"
	"fmt"
	"go/format"
	"io"
	"sort"
	"strings"
	"text/template"
	"unicode"

	"github.com/Dokiys/go_test/go/zz_my/tablec/basic"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var _ basic.Generator = (*Model)(nil)

type Model struct {
	*basic.Table
	Imports []string       // 需要导入的包
	PkgName string         // 包名
	Columns []*modelColumn // 表字段
}

type modelColumn struct {
	*basic.Column
	CamelName   string // 字段驼峰名称
	MappingType string // 映射类型，初始为空，需要各自实现
}

func NewModel(table *basic.Table, cols []*basic.Column, pkgName string) *Model {
	return &Model{
		Table:   table,
		Imports: importPkgs(cols),
		PkgName: pkgName,
		Columns: newModelColumns(cols),
	}
}

//go:embed model.tmpl
var modelTmpl string

func (t *Model) Gen(wr io.Writer) error {
	buf := &bytes.Buffer{}
	tmpl, err := template.New("modelTemp").Parse(strings.TrimSpace(modelTmpl))
	if err != nil {
		return fmt.Errorf("parse model temp failed: %w", err)
	}

	err = tmpl.Execute(buf, t)
	if err != nil {
		return fmt.Errorf("execute model temp failed: %w", err)
	}

	src, err := format.Source(buf.Bytes())
	if err != nil {
		return fmt.Errorf("format source failed: %w", err)
	}

	_, err = wr.Write(src)
	if err != nil {
		return fmt.Errorf("write source failed: %w", err)
	}

	return nil
}

func (t *Model) ModelName() string {
	return snake2Camel(t.GetTableName())
}

func (t *Model) ReceiverName() string {
	if len(t.GetTableName()) == 0 {
		return ""
	}

	return strings.ToLower(t.GetTableName()[0:1])
}

func importPkgs(cols []*basic.Column) (result []string) {
	var m = make(map[string]struct{})
	for _, c := range cols {
		if c.GetColumnType() == "date" || c.GetColumnType() == "datetime" || c.GetColumnType() == "timestamp" {
			m["time"] = struct{}{}
		}
	}

	for k := range m {
		result = append(result, k)
	}

	return
}

func newModelColumns(col []*basic.Column) []*modelColumn {
	cols := make([]*modelColumn, 0, len(col))
	for _, c := range col {
		cols = append(cols, &modelColumn{
			Column:      c,
			CamelName:   snake2Camel(c.GetColumnName()),
			MappingType: mapType(c.GetDataType()),
		})
	}
	sort.Slice(cols, func(i, j int) bool {
		return cols[i].GetOrdinalPosition() < cols[j].GetOrdinalPosition()
	})
	return cols
}

func mapType(dataType string) string {
	switch dataType {
	case "int":
		return "int32"
	case "tinyint":
		return "int32"
	case "smallint":
		return "int32"
	case "mediumint":
		return "int32"
	case "enum":
		return "int32"
	case "bigint":
		return "int64"
	case "char":
		return "string"
	case "varchar":
		return "string"
	case "json":
		return "string"
	case "timestamp":
		return "time.Time"
	case "date":
		return "time.Time"
	case "datetime":
		return "time.Time"
	case "text":
		return "string"
	case "mediumtext":
		return "string"
	case "longtext":
		return "string"
	case "double":
		return "float64"
	case "decimal":
		return "float64"
	case "float":
		return "float64"
	}

	panic("Unknown type: " + dataType)
}

// 将大小写的蛇形字符串转换为驼峰字符串
func snake2Camel(name string) string {
	var enCases = cases.Title(language.AmericanEnglish, cases.NoLower)

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
