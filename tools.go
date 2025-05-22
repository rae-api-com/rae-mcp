package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/rs/zerolog/log"
)

// handleSearchTool is a tool that searches the RAE API.
func handleSearchTool(
	ctx context.Context,
	request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	query, ok := request.GetArguments()["word"].(string)
	if !ok || query == "" {
		log.Error().Msg("missing or empty query in handleSearchTool")
		return nil, fmt.Errorf("missing or empty query")
	}

	log.Info().Str("query", query).Msg("Searching RAE API")
	results, err := client.Word(ctx, query)
	if err != nil {
		log.Error().Err(err).Str("query", query).Msg("RAE API search error")
		return nil, fmt.Errorf("RAE API search error: %v", err)
	}

	// Format the results for the LLM
	output, err := formatSearchResults(results)
	if err != nil {
		log.Error().Err(err).Msg("error formatting search results")
		return nil, fmt.Errorf("error formatting results: %v", err)
	}

	log.Info().Str("query", query).Msg("Search successful")
	return mcp.NewToolResultText(output), nil
}

// handleGetWordInfoTool is a tool that gets detailed information about a word.
func handleGetWordInfoTool(
	ctx context.Context,
	request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	args := request.GetArguments()
	word, ok := args["word"].(string)
	if !ok || word == "" {
		log.Error().Msg("missing or empty word in handleGetWordInfoTool")
		return nil, fmt.Errorf("missing or empty word")
	}

	lang, ok := args["lang"].(string)
	if !ok || lang == "" {
		lang = "es" // Default to Spanish
	}

	log.Info().Str("word", word).Str("lang", lang).Msg("Getting word info from RAE API")
	wordInfo, err := client.Word(ctx, word)
	if err != nil {
		log.Error().Err(err).Str("word", word).Msg("RAE API word info error")
		return nil, fmt.Errorf("RAE API word info error: %v", err)
	}

	// Format the word information for the LLM
	output, err := formatWordInfo(wordInfo)
	if err != nil {
		log.Error().Err(err).Str("word", word).Msg("error formatting word info")
		return nil, fmt.Errorf("error formatting word info: %v", err)
	}

	log.Info().Str("word", word).Msg("GetWordInfo successful")
	return mcp.NewToolResultText(output), nil
}

// formatSearchResults formats the search results into a readable format for the LLM
func formatSearchResults(results any) (string, error) {
	// Convert results to a JSON string for now
	// In a production environment, you would create a more structured and readable format
	jsonData, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		log.Error().Err(err).Msg("error marshaling search results")
		return "", err
	}
	return string(jsonData), nil
}

// formatWordInfo formats word information into a readable format for the LLM
func formatWordInfo(wordInfo any) (string, error) {
	// Convert word info to a JSON string for now
	jsonData, err := json.MarshalIndent(wordInfo, "", "  ")
	if err != nil {
		log.Error().Err(err).Msg("error marshaling word info")
		return "", err
	}
	return string(jsonData), nil
}
