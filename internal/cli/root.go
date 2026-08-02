package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "golens",
	Short: "GoLens - Comprehensive health analysis for Go backend services",
	Long: `GoLens (golens) performs a comprehensive health analysis of Go backend services.
It helps developers identify correctness, performance, reliability, concurrency, security, observability, and maintainability issues.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
