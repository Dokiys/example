package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

var dsn string
var db *sql.DB
var dsnConf *mysql.Config

func initDb() (*sql.DB, func()) {
	var err error
	dsn = os.Getenv("MCP_SQL_DSN")
	dsnConf, err = mysql.ParseDSN(dsn)
	if err != nil {
		log.Fatalf("[db] failed parse dsn: %s\n", err)
	}
	log.Printf("[db] connect: %s %s\n", dsnConf.DBName, dsnConf.Addr)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("[db] failed to connect: %s %s\n", dsnConf.DBName, err)
	}

	if err := db.Ping(); err != nil {
		_ = db.Close()
		log.Fatalf("[db] failed to connect: %s %s\n", dsnConf.DBName, err)
	}

	db.SetMaxOpenConns(2)
	db.SetMaxIdleConns(2)

	return db, func() {
		db.Close()
	}
}
func main() {
	var cancel func()
	db, cancel = initDb()
	defer cancel()

	s := server.NewMCPServer(
		"Database Helper",
		"1.0.0",
		server.WithPromptCapabilities(true),
		server.WithToolCapabilities(true),
		server.WithLogging(),
	)

	// Add tool handler
	s.AddTool(mcp.NewTool("SQL_Helper",
		mcp.WithDescription("For any database query-related tasks, you can first seek help from this tool. The tool provides business-oriented database query logic, generates reference SQL statements, and offers processing workflows."),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return mcp.NewToolResultText(`
你可以按照以下流程去分析应该如何完成一次业务相关的查询：
1.生成SQL调用工具查询数据库中的所有表及其备注。
2.识别可能包含所需数据的相关表。
3.生成SQL调用工具查看相关表的结构以及表字段的备注信息。
4.编写需求相关的查询语句。生成SQL调用工具进行查询
处理查询结果，整理结果后回复。

你需要注意：
禁止去猜想有关数据库的任何信息，包括database、需要查询的表名称、表字段名称等，这些信息你都需要生成SQL调用工具去查询数据库中的备注信息获得。
必须根据表名和字段名和Comment信息综合考虑需要执行的操作。
每一个步骤的操作只能包含单条可执行的SQL语句。
`), nil
	})

	s.AddTool(mcp.NewTool("SQL_Executor",
		mcp.WithDescription("This tool can help execute SQL queries and return information, with support for single SQL statement execution only. Under all circumstances, this tool is restricted to executing SQL query operations only and is prohibited from performing any data modification operations."),
		mcp.WithString("sql",
			mcp.Required(),
			mcp.Description("SQL statement to execute"),
		),
	), SqlExecuteHandler)

	if err := server.ServeStdio(s); err != nil {
		return
	}
}

func SqlExecuteHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	sqlStr := req.GetString("sql", "")
	if sqlStr == "" {
		return mcp.NewToolResultError("请输入SQL语句"), nil
	}

	resultJson, err := queryRowsAsJSON(ctx, db, sqlStr)
	if err != nil {
		return mcp.NewToolResultError("执行SQL错误：" + err.Error()), nil
	}

	return mcp.NewToolResultText(resultJson), nil
}

// queryRowsAsJSON 执行查询，将所有结果行映射为 json.RawMessage
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
		log.Fatal(err)
	}

	return string(finalJSON), nil
}
