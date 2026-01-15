package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ddt",
	Short: "Decision Debt Tracker - Track technical decisions and their assumptions",
	Long: `DDT is a CLI tool for logging and monitoring technical decisions.
	
It helps teams track:
- The decisions themselves (what was chosen)
- Context (why, who, when, constraints)
- Assumptions and expiration conditions
- Links to code, tickets, or docs
- Alerts when reality no longer matches assumptions`,
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}

// getWorkspaceRoot finds the repository root or current directory
func getWorkspaceRoot() (string, error) {
	// Try to find .git directory
	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// For now, just use current directory
	// In production, you might want to walk up to find .git
	return currentDir, nil
}

// printError prints an error message to stderr
func printError(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "Error: "+msg+"\n", args...)
}

// printSuccess prints a success message
func printSuccess(msg string, args ...interface{}) {
	fmt.Printf("✅ "+msg+"\n", args...)
}

// printWarning prints a warning message
func printWarning(msg string, args ...interface{}) {
	fmt.Printf("⚠️  "+msg+"\n", args...)
}

// printInfo prints an info message
func printInfo(msg string, args ...interface{}) {
	fmt.Printf("ℹ️  "+msg+"\n", args...)
}
