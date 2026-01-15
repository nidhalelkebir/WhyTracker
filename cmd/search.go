package cmd

import (
	"fmt"
	"strings"

	"github.com/nidhalelbkir/ddt/storage"
	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "Search for decisions",
	Long:  "Search decisions by keywords in title, reasoning, or tags",
	Args:  cobra.MinimumNArgs(1),
	RunE:  runSearch,
}

func init() {
	rootCmd.AddCommand(searchCmd)
}

func runSearch(cmd *cobra.Command, args []string) error {
	query := strings.Join(args, " ")

	root, err := getWorkspaceRoot()
	if err != nil {
		return fmt.Errorf("failed to get workspace root: %w", err)
	}

	store, err := storage.NewStorage(root)
	if err != nil {
		return fmt.Errorf("failed to initialize storage: %w", err)
	}
	defer store.Close()

	decisions, err := store.SearchDecisions(query)
	if err != nil {
		return fmt.Errorf("failed to search decisions: %w", err)
	}

	if len(decisions) == 0 {
		printInfo("No decisions found matching '%s'", query)
		return nil
	}

	fmt.Printf("\n🔍 Search Results for '%s' (%d found)\n\n", query, len(decisions))

	for _, decision := range decisions {
		printDecisionSummary(decision)
		fmt.Println()
	}

	return nil
}
