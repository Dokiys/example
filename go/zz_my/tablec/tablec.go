package tablec

import (
	"fmt"
	"io"

	"github.com/Dokiys/go_test/go/zz_my/tablec/basic"
	"github.com/Dokiys/go_test/go/zz_my/tablec/model"
)

const (
	TypModel = "model"
)

type TableC struct {
	conf    *DBConf
	tables  map[string]*basic.Table
	columns map[string][]*basic.Column
}

func NewTableC(conf *DBConf) (*TableC, error) {
	t := &TableC{conf: conf}
	if err := t.loadData(); err != nil {
		return nil, err
	}

	return t, nil
}

func (t *TableC) Do(wr io.Writer, typ, tableName string) error {
	switch typ {
	case TypModel:
		modelGenerator := model.NewModel(t.tables[tableName], t.columns[tableName], "module")
		if err := modelGenerator.Gen(wr); err != nil {
			return fmt.Errorf("generate model failed: %w", err)
		}
	default:
		return fmt.Errorf("not support type: %s", typ)
	}

	return nil
}

func (t *TableC) loadData() error {
	db, err := connectionDB(t.conf)
	if err != nil {
		return fmt.Errorf("connection db failed: %w", err)
	}
	defer db.Close()

	t.tables, err = loadTables(db, t.conf.Schema)
	if err != nil {
		return fmt.Errorf("load tables failed: %w", err)
	}

	t.columns, err = loadColumns(db, t.conf.Schema)
	if err != nil {
		return fmt.Errorf("load columns failed: %w", err)
	}

	return nil
}
