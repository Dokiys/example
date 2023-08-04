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
			TableName: "people",
			PK:        "id",
			// TODO[Dokiy] 2023/8/3: 验证是否需要转移
			TableComment: "测试表model渲染\\n 123",
		}
		// TODO[Dokiy] 2023/8/3: 补全
		cols := []*basic.Column{
			{
				ColumnName:      "id",
				OrdinalPosition: 1,
				DataType:        "int",
				ColumnType:      "int unsigned",
				ColumnComment:   "ID",
			}, {
				ColumnName:      "name",
				OrdinalPosition: 2,
				DataType:        "varchar",
				ColumnType:      "varchar(255)",
				ColumnComment:   "姓名",
				SrsId:           0,
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
