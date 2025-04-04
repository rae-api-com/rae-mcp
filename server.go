package main

import (
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rae-api-com/go-rae"
)

var (
	client = rae.New()
)

// RaeMCPServer represents the MCP server for RAE API
type RaeMCPServer struct {
	server *server.MCPServer
}

// NewRaeMCPServer creates a new MCP server with RAE API tools
func NewRaeMCPServer() *RaeMCPServer {
	mcpServer := server.NewMCPServer(
		"rae-api-mcp",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithPromptCapabilities(true),
		server.WithToolCapabilities(true),
	)

	// Register RAE API tools
	mcpServer.AddTool(mcp.NewTool("get_word_info",
		mcp.WithDescription("Get detailed information about a word from RAE API"),
		mcp.WithString("word",
			mcp.Description("The word to look up"),
			mcp.Required(),
		),
	), handleGetWordInfoTool)

	return &RaeMCPServer{
		server: mcpServer,
	}
}

// ServeSSE creates and returns an SSE server
func (s *RaeMCPServer) ServeSSE(addr string) *server.SSEServer {
	return server.NewSSEServer(s.server,
		server.WithBaseURL(fmt.Sprintf("http://%s", addr)),
	)
}

// ServeStdio serves the MCP server over stdio
func (s *RaeMCPServer) ServeStdio() error {
	return server.ServeStdio(s.server)
}
