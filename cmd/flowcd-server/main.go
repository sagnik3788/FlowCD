package main

import (
    "log"
    "net/http"

    "github.com/yourusername/flowcd/pkg/server/http"
)

func main() {
    // Initialize the HTTP server
    server := http.NewServer()

    // Start listening for incoming requests
    log.Println("Starting FlowCD server on :8080")
    if err := http.ListenAndServe(":8080", server); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}