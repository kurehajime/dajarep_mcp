package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/kurehajime/dajarep"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// Create a new MCP server
	s := server.NewMCPServer(
		"Dajarep",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithLogging(),
	)

	// Add a dajare tool
	dajareTool := mcp.NewTool("dajarep",
		mcp.WithDescription("Determines if a sentence contains a pun. 文章に駄洒落が含まれているか判定します"),
		mcp.WithString("sentence",
			mcp.Required(),
			mcp.Description("This argument is the sentence you want to determine if it is a pun."),
		),
	)

	// Add the dajare handler
	s.AddTool(dajareTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		sentence := request.Params.Arguments["sentence"].(string)
		s, _ := dajarep.Dajarep(sentence, 2, false)
		result := map[string]interface{}{
			"result":  len(s) > 0,
			"matches": s,
		}
		jsonResult, err := json.Marshal(result)
		if err != nil {
			return nil, fmt.Errorf("JSON変換エラー: %v", err)
		}
		return mcp.NewToolResultText(string(jsonResult)), nil
	})

	// Start the server
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
