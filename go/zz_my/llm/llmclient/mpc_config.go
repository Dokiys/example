package main

type ServerConfig struct {
	Type   string            `json:"type"`
	URL    string            `json:"url"`
	Header map[string]string `json:"header"`
}

type McpConfig struct {
	MCPServers map[string]ServerConfig `json:"mcpServers"`
}
