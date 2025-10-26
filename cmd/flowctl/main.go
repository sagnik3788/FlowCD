package main

import (
    "fmt"
    "os"
)

func main() {
    fmt.Println("Welcome to FlowCTL!")
    // Here you can add command handling logic for FlowCTL
    if len(os.Args) < 2 {
        fmt.Println("Please provide a command.")
        os.Exit(1)
    }

    command := os.Args[1]
    switch command {
    case "version":
        fmt.Println("FlowCTL version 0.1.0")
    case "help":
        fmt.Println("Available commands: version, help")
    default:
        fmt.Printf("Unknown command: %s\n", command)
        os.Exit(1)
    }
}