package main

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestMcpConfigLoad(t *testing.T) {
	// 示例 json 字符串（替换 ... 为具体内容）
	jsonData := []byte(`{
		"mcpServers": {
			"xxx1": {
				"type": "streamable-http-example",
				"url": "http://localhost:3001",
				"header": {
					"key": "value"
				}
			},
			"xxx2": {
				"type": "sse",
				"url": "http://localhost:8080",
				"header": {
					"key": "value"
				}
			}
		}
	}`)

	var config McpConfig
	err := json.Unmarshal(jsonData, &config)
	if err != nil {
		panic(err)
	}

	// 打印反序列化结果
	for name, server := range config.MCPServers {
		fmt.Printf("服务名：%s\n类型：%s\nURL：%s\nHeader：%v\n\n", name, server.Type, server.URL, server.Header)
	}
}
