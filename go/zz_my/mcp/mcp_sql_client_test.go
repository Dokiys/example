package mcp

import (
	"context"
	"fmt"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func TestSqlMcpClient(t *testing.T) {
	var cancel func()
	db, cancel = initDb()
	defer cancel()

	// Create a new MCP server
	s := server.NewMCPServer(
		"数据库查询助手",
		"1.0.0",
		server.WithToolCapabilities(false),
	)

	// Add tool handler
	s.AddTool(mcp.NewTool("数据库查询逻辑助手",
		mcp.WithDescription("任何涉及数据库查询的内容都可以先请求该工具获取帮助，该工具提供业务有关的数据库查询思维逻辑，提供可供参考的SQL生成以及处理流程。"),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return mcp.NewToolResultText(`
你可以按照以下流程去分析应该如何完成一次业务相关的查询：
1.生成SQL调用工具查询数据库中的所有表及其备注。
2.识别可能包含所需数据的相关表。
3.生成SQL调用工具查看相关表的结构以及表字段的备注信息。
4.编写需求相关的查询语句。生成SQL调用工具进行查询
处理查询结果，整理结果后回复。

你需要注意：
请不要自主去猜想有关数据库的任何信息，包括database、需要查询的表名称、表字段名称等，这些信息你都需要生成SQL调用工具去查询数据库中的备注信息获得。
你需要根据表和字段的name和Comment信息综合考虑需要执行的操作。
每一个步骤的操作只能包含单条可执行的SQL语句，可以反复多次的调用工具来执行每条SQL以完成需求。
`), nil
	})

	s.AddTool(mcp.NewTool("执行SQL工具",
		mcp.WithDescription("可以帮助执行查询类的SQL并返回信息，每次仅支持单条SQL执行。"),
		mcp.WithString("sql",
			mcp.Required(),
			mcp.Description("所要执行的sql语句"),
		),
	), SqlExecuteHandler)

	// Start the sse server
	sseServer := server.NewSSEServer(s)
	if err := sseServer.Start("localhost:8080"); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
