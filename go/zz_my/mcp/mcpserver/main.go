package main

import (
	"database/sql"
	"log"
	"os"

	"example/go/zz_my/mcp/mcpserver/tools"
	"github.com/go-sql-driver/mysql"
	"github.com/mark3labs/mcp-go/server"
)

func initDb() (*sql.DB, func()) {
	// 获取数据库连接
	dsn := os.Getenv("MCP_SQL_DSN")

	var err error
	dsnConf, err := mysql.ParseDSN(dsn)
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
	db, cancel := initDb()
	defer cancel()

	s := NewMCPServer()
	// Add Tools
	s.RegisterTool(tools.NewSqlExecutor(db))

	httpServer := server.NewStreamableHTTPServer(s.MCPServer)
	log.Printf("HTTP server listening on :8080/mcp")
	if err := httpServer.Start(":8080"); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

type MCPServer struct {
	*server.MCPServer
}

func NewMCPServer() MCPServer {
	mcpServer := server.NewMCPServer(
		"MCP服务助手",
		"1.0.0",
		server.WithPromptCapabilities(true),
		server.WithToolCapabilities(true),
		server.WithLogging(),
	)

	return MCPServer{
		mcpServer,
	}
}
func (m *MCPServer) RegisterTool(tool tools.Tool) {
	m.MCPServer.AddTool(tool.ToolInfo(), tool.Handle)
}
