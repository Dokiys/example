package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"example/go/zz_my/llm/mcpserver/tools"
	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
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
	godotenv.Load()
	db, cancel := initDb()
	defer cancel()

	s := NewMCPServer()
	// Add Tools
	s.RegisterTool(tools.NewSqlExecutor(db))
	httpServer := server.NewStreamableHTTPServer(s.MCPServer)

	log.Printf("HTTP server listening on :8081/mcp")
	err := http.ListenAndServe(":8081", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")

		if auth != "123" {
			w.WriteHeader(http.StatusForbidden)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"error":"forbidden"}`))
			return
		}
		httpServer.ServeHTTP(w, r)
	}))
	if err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

type MCPServer struct {
	*server.MCPServer
}

func NewMCPServer() MCPServer {
	hooks := &server.Hooks{}
	hooks.AddOnRequestInitialization(func(ctx context.Context, id any, message any) error {
		fmt.Printf("AddOnRequestInitialization: %v, %v\n", id, message)
		return nil
	})
	mcpServer := server.NewMCPServer(
		"MCP服务助手",
		"1.0.0",
		server.WithPromptCapabilities(true),
		server.WithToolCapabilities(true),
		server.WithLogging(),
		server.WithHooks(hooks),
	)

	return MCPServer{
		mcpServer,
	}
}
func (m *MCPServer) RegisterTool(tool tools.Tool) {
	m.MCPServer.AddTool(tool.ToolInfo(), tool.Handle)
}
