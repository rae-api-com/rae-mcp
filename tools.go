package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
)

// handleSearchTool is a tool that searches the RAE API.
func handleSearchTool(
	ctx context.Context,
	request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	query, ok := request.Params.Arguments["word"].(string)
	if !ok || query == "" {
		return nil, fmt.Errorf("missing or empty query")
	}

	results, err := client.Word(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("RAE API search error: %v", err)
	}

	// Format the results for the LLM
	output, err := formatSearchResults(results)
	if err != nil {
		return nil, fmt.Errorf("error formatting results: %v", err)
	}

	return mcp.NewToolResultText(output), nil
}

// handleGetWordInfoTool is a tool that gets detailed information about a word.
func handleGetWordInfoTool(
	ctx context.Context,
	request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	word, ok := request.Params.Arguments["word"].(string)
	if !ok || word == "" {
		return nil, fmt.Errorf("missing or empty word")
	}

	lang, ok := request.Params.Arguments["lang"].(string)
	if !ok || lang == "" {
		lang = "es" // Default to Spanish
	}

	wordInfo, err := client.Word(ctx, word)
	if err != nil {
		return nil, fmt.Errorf("RAE API word info error: %v", err)
	}

	// Format the word information for the LLM
	output, err := formatWordInfo(wordInfo)
	if err != nil {
		return nil, fmt.Errorf("error formatting word info: %v", err)
	}

	return mcp.NewToolResultText(output), nil
}

// formatSearchResults formats the search results into a readable format for the LLM
func formatSearchResults(results any) (string, error) {
	// Convert results to a JSON string for now
	// In a production environment, you would create a more structured and readable format
	jsonData, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

// formatWordInfo formats word information into a readable format for the LLM
func formatWordInfo(wordInfo any) (string, error) {
	// Convert word info to a JSON string for now
	jsonData, err := json.MarshalIndent(wordInfo, "", "  ")
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}
