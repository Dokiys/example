package mcp

import (
	"context"
	"fmt"
	"testing"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/client/transport"
	"github.com/mark3labs/mcp-go/mcp"
)

func TestSqlMcpClient(t *testing.T) {
	var ctx = context.Background()
	var httpURL = "http://localhost:8080/mcp"

	httpTransport, err := transport.NewStreamableHTTP(httpURL)
	if err != nil {
		t.Fatalf("Failed to create HTTP transport: %v", err)
	}

	c := client.NewClient(httpTransport)
	if err := c.Start(ctx); err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	// Set up notification handler
	c.OnNotification(func(notification mcp.JSONRPCNotification) {
		fmt.Printf("Received notification: %s\n", notification.Method)
	})

	// Initialize the client
	fmt.Println("Initializing client...")
	initRequest := mcp.InitializeRequest{}
	initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initRequest.Params.ClientInfo = mcp.Implementation{
		Name:    "McpList-Go Simple Client Example",
		Version: "1.0.0",
	}
	initRequest.Params.Capabilities = mcp.ClientCapabilities{}

	serverInfo, err := c.Initialize(ctx, mcp.InitializeRequest{})
	if err != nil {
		t.Fatalf("Failed to initialize: %v", err)
	}
	// Display server information
	fmt.Printf("Connected to server: %s (version %s)\n", serverInfo.ServerInfo.Name, serverInfo.ServerInfo.Version)
	fmt.Printf("Server capabilities: %+v\n", serverInfo.Capabilities)

	// List available tools if the server supports them
	if serverInfo.Capabilities.Tools != nil {
		fmt.Println("Fetching available tools...")
		toolsRequest := mcp.ListToolsRequest{}
		toolsResult, err := c.ListTools(ctx, toolsRequest)
		if err != nil {
			fmt.Printf("Failed to list tools: %v", err)
		} else {
			fmt.Printf("Server has %d tools available\n", len(toolsResult.Tools))
			for i, tool := range toolsResult.Tools {
				fmt.Printf("  %d. %s - %s\n", i+1, tool.Name, tool.Description)
			}
		}
	}
}
