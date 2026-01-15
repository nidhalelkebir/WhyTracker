package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/nidhalelbkir/ddt/storage"
	"github.com/spf13/cobra"
)

var (
	listAll    bool
	listLimit  int
	listStatus string
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List decisions",
	Long:  "List all logged decisions with their status and basic info",
	RunE:  runList,
}

func init() {
	listCmd.Flags().BoolVarP(&listAll, "all", "a", false, "Show all decisions including retired")
	listCmd.Flags().IntVarP(&listLimit, "limit", "l", 10, "Limit number of results (0 for all)")
	listCmd.Flags().StringVarP(&listStatus, "status", "s", "active", "Filter by status (active, updated, retired)")
	rootCmd.AddCommand(listCmd)
}

func runList(cmd *cobra.Command, args []string) error {
	root, err := getWorkspaceRoot()
	if err != nil {
		return fmt.Errorf("failed to get workspace root: %w", err)
	}

	store, err := storage.NewStorage(root)
	if err != nil {
		return fmt.Errorf("failed to initialize storage: %w", err)
	}
	defer store.Close()

	status := listStatus
	if listAll {
		status = ""
	}

	decisions, err := store.ListDecisions(status, listLimit)
	if err != nil {
		return fmt.Errorf("failed to list decisions: %w", err)
	}

	if len(decisions) == 0 {
		printInfo("No decisions found")
		return nil
	}

	fmt.Printf("\n📋 Recent Decisions (%d)\n\n", len(decisions))

	for _, decision := range decisions {
		printDecisionSummary(decision)
		fmt.Println()
	}

	return nil
}

func printDecisionSummary(d interface{}) {
	switch decision := d.(type) {
	case *storage.Decision:
		// Handle storage.Decision if needed
	default:
		// Handle models.Decision
		if dec, ok := d.(*models.Decision); ok {
			statusIcon := getStatusIcon(dec.Status)
			fmt.Printf("%s [%s] %s\n", statusIcon, dec.ID, dec.Title)
			fmt.Printf("   Decided: %s", dec.DecidedAt.Format("2006-01-02"))
			if dec.DecidedBy != "" {
				fmt.Printf(" by %s", dec.DecidedBy)
			}
			fmt.Println()
			
			if dec.Reasoning != "" {
				reasoning := dec.Reasoning
				if len(reasoning) > 80 {
					reasoning = reasoning[:77] + "..."
				}
				fmt.Printf("   Why: %s\n", reasoning)
			}
			
			if len(dec.Assumptions) > 0 {
				fmt.Printf("   Assumptions: %d | Triggers: %d\n", 
					len(dec.Assumptions), len(dec.ExpirationTriggers))
			}
			
			if len(dec.Tags) > 0 {
				fmt.Printf("   Tags: %s\n", strings.Join(dec.Tags, ", "))
			}
		}
	}
}

func getStatusIcon(status string) string {
	switch status {
	case "active":
		return "🟢"
	case "updated":
		return "🟡"
	case "retired":
		return "⚫"
	default:
		return "⚪"
	}
}

func formatDuration(d time.Duration) string {
	days := int(d.Hours() / 24)
	if days == 0 {
		return "today"
	} else if days == 1 {
		return "1 day ago"
	} else if days < 7 {
		return fmt.Sprintf("%d days ago", days)
	} else if days < 30 {
		weeks := days / 7
		if weeks == 1 {
			return "1 week ago"
		}
		return fmt.Sprintf("%d weeks ago", weeks)
	} else if days < 365 {
		months := days / 30
		if months == 1 {
			return "1 month ago"
		}
		return fmt.Sprintf("%d months ago", months)
	} else {
		years := days / 365
		if years == 1 {
			return "1 year ago"
		}
		return fmt.Sprintf("%d years ago", years)
	}
}
