package main

import (
    "fmt"
    "os"

    "github.com/spf13/cobra"
)

const (
    version = "0.1.0"
    banner = `
    ███████╗██╗      ██████╗ ██╗    ██╗ ██████╗██████╗ 
    ██╔════╝██║     ██╔═══██╗██║    ██║██╔════╝██╔══██╗
    █████╗  ██║     ██║   ██║██║ █╗ ██║██║     ██║  ██║
    ██╔══╝  ██║     ██║   ██║██║███╗██║██║     ██║  ██║
    ██║     ███████╗╚██████╔╝╚███╔███╔╝╚██████╗██████╔╝
    ╚═╝     ╚══════╝ ╚═════╝  ╚══╝╚══╝  ╚═════╝╚═════╝ 
                                                        
    Continuous Deployment with Flow - v%s
    `
)

var rootCmd = &cobra.Command{
    Use:   "flowctl",
    Short: "FlowCD CLI - Continuous Deployment with Flow",
    Long: fmt.Sprintf(banner, version) + `
FlowCD is a Kubernetes-native continuous deployment tool that enables 
GitOps workflows with advanced deployment strategies.

Use flowctl to manage your FlowCD resources, pipelines, and deployments.`,
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Printf(banner, version)
        fmt.Println("\nUse 'flowctl --help' to see available commands")
    },
}

var versionCmd = &cobra.Command{
    Use:   "version",
    Short: "Print the version number of flowctl",
    Long:  "All software has versions. This is FlowCD's",
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Printf("FlowCD CLI v%s\n", version)
    },
}

func init() {
    rootCmd.AddCommand(versionCmd)
}

func main() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
}