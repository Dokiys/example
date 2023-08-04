package tablec

import (
	"database/sql"
	"fmt"
	"regexp"

	"github.com/Dokiys/go_test/go/zz_my/tablec/basic"
	_ "github.com/go-sql-driver/mysql"
)

const (
	queryTableCommentSQL = `
SELECT t.TABLE_NAME, t.TABLE_COMMENT, c.COLUMN_NAME AS PK
FROM information_schema.TABLES t LEFT JOIN information_schema.COLUMNS c ON t.TABLE_NAME = c.TABLE_NAME
WHERE t.TABLE_SCHEMA = ? AND c.COLUMN_KEY = 'PRI';
`
	queryTableColumnSQL = `
SELECT c.*
FROM INFORMATION_SCHEMA.columns c
WHERE c.TABLE_SCHEMA = ?;
`
)

type DBConf struct {
	Host     string `yaml:"host" json:"host"`
	Port     int    `yaml:"port" json:"port"`
	Schema   string `yaml:"schema" json:"schema"`
	Username string `yaml:"username" json:"username"`
	Password string `yaml:"password" json:"password"`
}

var lineBreakRegex = regexp.MustCompile("[\r\n\t]")

const driverName = "mysql"

func connectionDB(conf *DBConf) (*sql.DB, error) {
	var dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true", conf.Username, conf.Password, conf.Host, conf.Port, conf.Schema)

	db, err := sql.Open(driverName, dsn)
	if err != nil {
		return nil, fmt.Errorf("connect to database error: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping database error: %w", err)
	}
	return db, nil
}

func loadTables(db *sql.DB, schema string) (map[string]*basic.Table, error) {
	tables := make(map[string]*basic.Table)
	rows, err := db.Query(queryTableCommentSQL, schema)
	if err != nil {
		return nil, fmt.Errorf("query table comment error: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var table = &basic.Table{}
		if err := rows.Scan(&table.TableName, &table.TableComment, &table.PK); err != nil {
			return nil, fmt.Errorf("scan table comment error: %w", err)
		}

		*table.TableComment = lineBreakRegex.ReplaceAllLiteralString(*table.TableComment, " ")
		tables[table.GetTableName()] = table
	}
	return tables, nil
}

func loadColumns(db *sql.DB, schema string) (map[string][]*basic.Column, error) {
	crows, err := db.Query(queryTableColumnSQL, schema)
	if err != nil {
		return nil, fmt.Errorf("query table column error: %w", err)
	}
	defer crows.Close()

	var columns = make(map[string][]*basic.Column)
	for crows.Next() {
		var column = &basic.Column{}
		if err := crows.Scan(
			&column.TableCatalog,
			&column.TableSchema,
			&column.TableName,
			&column.ColumnName,
			&column.OrdinalPosition,
			&column.ColumnDefault,
			&column.IsNullable,
			&column.DataType,
			&column.CharacterMaximumLength,
			&column.CharacterOctetLength,
			&column.NumericPrecision,
			&column.NumericScale,
			&column.DatetimePrecision,
			&column.CharacterSetName,
			&column.CollationName,
			&column.ColumnType,
			&column.ColumnKey,
			&column.Extra,
			&column.Privileges,
			&column.ColumnComment,
			&column.GenerationExpression,
			&column.SrsId,
		); err != nil {
			return nil, fmt.Errorf("scan table column error: %w", err)
		}

		*column.ColumnComment = lineBreakRegex.ReplaceAllLiteralString(column.GetColumnComment(), " ")
		columns[column.GetTableName()] = append(columns[column.GetTableName()], column)
	}

	return columns, nil
}
