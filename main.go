package main

import (
	"flag"
	"fmt"
	"log"
)

func main() {
	var transport string
	var port int

	flag.StringVar(&transport, "t", "stdio", "Transport type (stdio or sse)")
	flag.StringVar(&transport, "transport", "stdio", "Transport type (stdio or sse)")
	flag.IntVar(&port, "port", 8080, "Port for SSE server")
	flag.Parse()

	s := NewRaeMCPServer()

	switch transport {
	case "stdio":
		if err := s.ServeStdio(); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	case "sse":
		addr := fmt.Sprintf("localhost:%d", port)
		sseServer := s.ServeSSE(addr)
		log.Printf("SSE server listening on :%d", port)
		if err := sseServer.Start(fmt.Sprintf(":%d", port)); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	default:
		log.Fatalf(
			"Invalid transport type: %s. Must be 'stdio' or 'sse'",
			transport,
		)
	}
}
