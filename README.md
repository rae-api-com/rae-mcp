# RAE Model Context Protocol (MCP) Server

This repository contains an implementation of a Model Context Protocol (MCP) server for the Royal Spanish Academy (RAE) API. It allows language models to interact with RAE's dictionary and linguistic resources.

## Requirements

- Go 1.21+

## Installation

```bash
git clone https://github.com/rae-api-com/rae-mpc.git
cd rae-mpc
go build
```

## Usage

### Command Line Arguments

Run the server with stdio transport (for integration with LLMs):

```bash
./rae-mpc --transport stdio
```

Or run it as an SSE server:

```bash
./rae-mpc --transport sse --port 8080
```

### Available Tools

The MCP server exposes the following tools to LLMs:

1. `search` - Search RAE API for information
   - Parameters:
     - `query` (required): The search query
     - `lang` (optional): Language code (default: "es")

2. `get_word_info` - Get detailed information about a word
   - Parameters:
     - `word` (required): The word to look up
     - `lang` (optional): Language code (default: "es")

## Integration with LLMs

This MCP server can be integrated with language models that support the Model Context Protocol, allowing them to access RAE's linguistic resources for improved Spanish language capabilities.
