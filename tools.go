package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/rs/zerolog"
)

// handleSearchTool is a tool that searches the RAE API.
func handleSearchTool(
	logger zerolog.Logger,
) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		query, ok := request.GetArguments()["word"].(string)
		if !ok || query == "" {
			logger.Error().Msg("missing or empty query in handleSearchTool")
			return nil, fmt.Errorf("missing or empty query")
		}

		logger.Info().Str("query", query).Msg("Searching RAE API")
		results, err := client.Word(ctx, query)
		if err != nil {
			logger.Error().Err(err).Str("query", query).Msg("RAE API search error")
			return nil, fmt.Errorf("RAE API search error: %v", err)
		}

		// Format the results for the LLM
		output, err := formatSearchResults(logger, results)
		if err != nil {
			logger.Error().Err(err).Msg("error formatting search results")
			return nil, fmt.Errorf("error formatting results: %v", err)
		}

		logger.Info().Str("query", query).Msg("Search successful")
		return mcp.NewToolResultText(output), nil
	}
}

// newHandleGetWordInfoTool returns a tool that gets detailed information about a word.
func newHandleGetWordInfoTool(
	logger zerolog.Logger,
) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		word, ok := args["word"].(string)
		if !ok || word == "" {
			logger.Error().Msg("missing or empty word in handleGetWordInfoTool")
			return nil, fmt.Errorf("missing or empty word")
		}

		lang, ok := args["lang"].(string)
		if !ok || lang == "" {
			lang = "es" // Default to Spanish
		}

		logger.Info().Str("word", word).Str("lang", lang).Msg("Getting word info from RAE API")
		wordInfo, err := client.Word(ctx, word)
		if err != nil {
			if len(wordInfo.Suggestions) > 0 {
				logger.Warn().Str("word", word).Msg("Word not found, suggesting alternatives")
				output, err := formatSuggestions(wordInfo.Suggestions)
				if err != nil {
					logger.Error().Err(err).Str("word", word).Msg("error formatting suggestions")
					return nil, fmt.Errorf("error formatting suggestions: %v", err)
				}

				return mcp.NewToolResultText(output), nil
			}

			logger.Error().Err(err).Str("word", word).Msg("RAE API word info error")
			return nil, fmt.Errorf("RAE API word info error: %v", err)
		}

		// Format the word information for the LLM
		output, err := formatWordInfo(logger, wordInfo)
		if err != nil {
			logger.Error().Err(err).Str("word", word).Msg("error formatting word info")
			return nil, fmt.Errorf("error formatting word info: %v", err)
		}

		logger.Info().Str("word", word).Msg("GetWordInfo successful")
		return mcp.NewToolResultText(output), nil
	}
}

// formatSuggestions formats the suggestions into a readable format for the LLM
func formatSuggestions(res []string) (string, error) {
	type suggestions struct {
		Suggestions []string `json:"suggestions"`
		Msg         string   `json:"msg"`
	}
	if len(res) == 0 {
		return "No suggestions available", nil
	}

	bts, err := json.MarshalIndent(suggestions{Suggestions: res, Msg: "Did you mean one of these words?"}, "", "  ")
	if err != nil {
		return "", err
	}
	return string(bts), nil
}

// formatSearchResults formats the search results into a readable format for the LLM
func formatSearchResults(logger zerolog.Logger, results any) (string, error) {
	// Convert results to a JSON string for now
	// In a production environment, you would create a more structured and readable format
	jsonData, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		logger.Error().Err(err).Msg("error marshaling search results")
		return "", err
	}
	return string(jsonData), nil
}

// formatWordInfo formats word information into a readable format for the LLM
func formatWordInfo(logger zerolog.Logger, wordInfo any) (string, error) {
	// Convert word info to a JSON string for now
	jsonData, err := json.MarshalIndent(wordInfo, "", "  ")
	if err != nil {
		logger.Error().Err(err).Msg("error marshaling word info")
		return "", err
	}
	return string(jsonData), nil
}

// newHandleGetDailyWordTool returns a tool that gets the word of the day.
func newHandleGetDailyWordTool(
	logger zerolog.Logger,
) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logger.Info().Msg("Getting daily word from RAE API")

		dailyWord, err := client.Daily(ctx)
		if err != nil {
			logger.Error().Err(err).Msg("RAE API daily word error")
			return nil, fmt.Errorf("RAE API daily word error: %v", err)
		}

		// Format the daily word information for the LLM
		output, err := formatDailyWord(logger, dailyWord)
		if err != nil {
			logger.Error().Err(err).Msg("error formatting daily word")
			return nil, fmt.Errorf("error formatting daily word: %v", err)
		}

		logger.Info().Msg("GetDailyWord successful")
		return mcp.NewToolResultText(output), nil
	}
}

// newHandleGetRandomWordTool returns a tool that gets a random word.
func newHandleGetRandomWordTool(
	logger zerolog.Logger,
) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()

		// Parse optional min_length parameter
		var minLength *int
		if minLengthStr, ok := args["min_length"].(string); ok && minLengthStr != "" {
			if val, err := strconv.Atoi(minLengthStr); err == nil && val > 0 {
				minLength = &val
			}
		} else if minLengthFloat, ok := args["min_length"].(float64); ok && minLengthFloat > 0 {
			val := int(minLengthFloat)
			minLength = &val
		}

		// Parse optional max_length parameter
		var maxLength *int
		if maxLengthStr, ok := args["max_length"].(string); ok && maxLengthStr != "" {
			if val, err := strconv.Atoi(maxLengthStr); err == nil && val > 0 {
				maxLength = &val
			}
		} else if maxLengthFloat, ok := args["max_length"].(float64); ok && maxLengthFloat > 0 {
			val := int(maxLengthFloat)
			maxLength = &val
		}

		logger.Info().
			Interface("min_length", minLength).
			Interface("max_length", maxLength).
			Msg("Getting random word from RAE API")

		randomWord, err := client.Random(ctx)
		if err != nil {
			logger.Error().Err(err).Msg("RAE API random word error")
			return nil, fmt.Errorf("RAE API random word error: %v", err)
		}

		// Format the random word information for the LLM
		output, err := formatRandomWord(logger, randomWord)
		if err != nil {
			logger.Error().Err(err).Msg("error formatting random word")
			return nil, fmt.Errorf("error formatting random word: %v", err)
		}

		logger.Info().Msg("GetRandomWord successful")
		return mcp.NewToolResultText(output), nil
	}
}

// formatDailyWord formats daily word information into a readable format for the LLM
func formatDailyWord(logger zerolog.Logger, dailyWord any) (string, error) {
	jsonData, err := json.MarshalIndent(dailyWord, "", "  ")
	if err != nil {
		logger.Error().Err(err).Msg("error marshaling daily word")
		return "", err
	}
	return fmt.Sprintf("Word of the Day:\n%s", string(jsonData)), nil
}

// formatRandomWord formats random word information into a readable format for the LLM
func formatRandomWord(logger zerolog.Logger, randomWord any) (string, error) {
	jsonData, err := json.MarshalIndent(randomWord, "", "  ")
	if err != nil {
		logger.Error().Err(err).Msg("error marshaling random word")
		return "", err
	}
	return fmt.Sprintf("Random Word:\n%s", string(jsonData)), nil
}
