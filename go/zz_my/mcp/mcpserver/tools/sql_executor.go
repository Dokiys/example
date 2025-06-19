package tools

import (
	"context"
	"database/sql"
	_ "embed"
	"encoding/json"

	"github.com/mark3labs/mcp-go/mcp"
)

var _ Tool = (*SqlExecutor)(nil)

type SqlExecutor struct {
	db *sql.DB
}

func NewSqlExecutor(db *sql.DB) *SqlExecutor {
	return &SqlExecutor{
		db: db,
	}
}
func (s *SqlExecutor) ToolInfo() mcp.Tool {
	return mcp.NewTool("SQL_Executor",
		mcp.WithDescription("该工具能且仅能帮助执行数据查询语句并返回信息"),
		mcp.WithString("sql",
			mcp.Required(),
			mcp.Description("需要执行的SQL语句"),
		),
	)
}
func (s *SqlExecutor) Handle(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	sqlStr := req.GetString("sql", "")
	if sqlStr == "" {
		return mcp.NewToolResultError("请输入SQL语句"), nil
	}

	resultJson, err := queryRowsAsJSON(ctx, s.db, sqlStr)
	if err != nil {
		return mcp.NewToolResultError("执行SQL错误：" + err.Error()), nil
	}

	return mcp.NewToolResultText(resultJson), nil
}
func queryRowsAsJSON(ctx context.Context, db *sql.DB, query string, args ...interface{}) (string, error) {
	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return "", err
	}

	var results []json.RawMessage

	for rows.Next() {
		// 创建一个 []interface{} 用于接值
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))

		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return "", err
		}

		// 构造 map[string]interface{}
		rowMap := make(map[string]interface{})
		for i, col := range columns {
			val := values[i]
			if b, ok := val.([]byte); ok {
				rowMap[col] = string(b)
			} else {
				rowMap[col] = val
			}
		}

		// 将 map 转为 JSON 并包裹成 RawMessage
		jsonBytes, err := json.Marshal(rowMap)
		if err != nil {
			return "", err
		}

		results = append(results, jsonBytes)
	}

	if err := rows.Err(); err != nil {
		return "", err
	}
	// 转为最终的 JSON 数组字符串
	finalJSON, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return "", err
	}

	return string(finalJSON), nil
}
