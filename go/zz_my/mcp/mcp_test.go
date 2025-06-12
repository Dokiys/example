package mcp

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func TestMCP(t *testing.T) {
	// Create a new MCP server
	s := server.NewMCPServer(
		"Current Time",
		"1.0.0",
		server.WithToolCapabilities(false),
	)

	// Add tool
	tool := mcp.NewTool("Current Time",
		mcp.WithDescription("获取时间助手"),
		mcp.WithString("location",
			mcp.Required(),
			mcp.Description("所要获取时间的地区"),
		),
	)

	// Add tool handler
	s.AddTool(tool, currentHandler)

	// Start the stdio server
	// go build -o current_time && PATH="$PATH:`pwd`"
	//
	// 客户的请求：
	// {"method":"initialize","params":{"protocolVersion":"2025-03-26","capabilities":{},"clientInfo":{"name":"Cherry Studio","version":"1.4.1"}},"jsonrpc":"2.0","id":0}
	// {"jsonrpc":"2.0","method":"notifications/cancelled","params":{"requestId":0,"reason":"McpError: MCP error -32001: Request timed out"}}
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}

	// Start the sse server
	// sseServer := server.NewSSEServer(s)
	// if err := sseServer.Start("localhost:8080"); err != nil {
	// 	fmt.Printf("Server error: %v\n", err)
	// }

}
func currentHandler(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	_, err := request.RequireString("location")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText(time.Now().Add(-2 * time.Hour).Format("2006-01-02 15:04:05")), nil
}
