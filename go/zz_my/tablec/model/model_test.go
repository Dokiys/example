package model

import (
	"bytes"
	"testing"

	"github.com/Dokiys/go_test/go/zz_my/tablec/basic"
)

func TestRenderModelTemp(t *testing.T) {
	t.Run("测试渲染模版", func(t *testing.T) {
		buf := &bytes.Buffer{}
		table := &basic.Table{
			TableName:    pointVar("people"),
			PK:           pointVar("id"),
			TableComment: pointVar("测试表model渲染\\n 123"),
		}
		cols := []*basic.Column{
			{
				ColumnName:      pointVar("id"),
				OrdinalPosition: pointVar(1),
				DataType:        pointVar("int"),
				ColumnType:      pointVar("int unsigned"),
				ColumnComment:   pointVar("ID"),
			}, {
				ColumnName:      pointVar("name"),
				OrdinalPosition: pointVar(2),
				DataType:        pointVar("varchar"),
				ColumnType:      pointVar("varchar(255)"),
				ColumnComment:   pointVar("姓名"),
				SrsId:           pointVar(0),
			},
		}
		data := NewModel(table, cols, "module")
		err := data.Gen(buf)
		if err != nil {
			t.Fatal(err)
		}

		t.Log(buf.String())
	})
}

func pointVar[T comparable](t T) *T {
	var zero T
	if t == zero {
		return &zero
	}

	return &t
}
