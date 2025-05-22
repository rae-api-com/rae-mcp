package main

import (
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rae-api-com/go-rae"
	"github.com/rs/zerolog"
)

var (
	client = rae.New()
)

// RaeMCPServer represents the MCP server for RAE API
type RaeMCPServer struct {
	server *server.MCPServer
}

// Updated NewRaeMCPServer function with all tools
func NewRaeMCPServer(logger zerolog.Logger) *RaeMCPServer {
	mcpServer := server.NewMCPServer(
		"rae-api-mcp",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithPromptCapabilities(true),
		server.WithToolCapabilities(true),
	)

	// Register existing RAE API tool
	mcpServer.AddTool(mcp.NewTool("get_word_info",
		mcp.WithDescription("Get detailed information about a word from RAE API"),
		mcp.WithString("word",
			mcp.Description("The word to look up"),
			mcp.Required(),
		),
	), newHandleGetWordInfoTool(logger))

	// Register new Daily Word tool
	mcpServer.AddTool(mcp.NewTool("get_daily_word",
		mcp.WithDescription("Get the word of the day from RAE dictionary"),
	), newHandleGetDailyWordTool(logger))

	// Register new Random Word tool
	mcpServer.AddTool(mcp.NewTool("get_random_word",
		mcp.WithDescription("Get a random word from RAE dictionary"),
		mcp.WithNumber("min_length",
			mcp.Description("Minimum length of the random word"),
		),
		mcp.WithNumber("max_length",
			mcp.Description("Maximum length of the random word"),
		),
	), newHandleGetRandomWordTool(logger))

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
