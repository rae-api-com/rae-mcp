package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/rs/zerolog"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger()

	var transport string
	var port int

	flag.StringVar(&transport, "t", "stdio", "Transport type (stdio or sse)")
	flag.StringVar(&transport, "transport", "stdio", "Transport type (stdio or sse)")
	flag.IntVar(&port, "port", 8080, "Port for SSE server")
	flag.Parse()

	logger.Info().Str("transport", transport).Int("port", port).Msg("Starting server")

	s := NewRaeMCPServer()

	switch transport {
	case "stdio":
		logger.Info().Msg("Serving via stdio")
		if err := s.ServeStdio(); err != nil {
			logger.Fatal().Err(err).Msg("Server error")
		}
	case "sse":
		addr := fmt.Sprintf("localhost:%d", port)
		sseServer := s.ServeSSE(addr)
		logger.Info().Int("port", port).Msg("SSE server listening")
		if err := sseServer.Start(fmt.Sprintf(":%d", port)); err != nil {
			logger.Fatal().Err(err).Msg("Server error")
		}
	default:
		logger.Fatal().
			Str("transport", transport).
			Msg("Invalid transport type: must be 'stdio' or 'sse'")
	}
}
