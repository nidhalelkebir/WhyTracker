package cmd

import (
	"fmt"
	"strings"

	"github.com/nidhalelbkir/ddt/storage"
	"github.com/spf13/cobra"
)

var showCmd = &cobra.Command{
	Use:   "show [decision-id]",
	Short: "Show full details of a decision",
	Long:  "Display complete information about a specific decision",
	Args:  cobra.ExactArgs(1),
	RunE:  runShow,
}

func init() {
	rootCmd.AddCommand(showCmd)
}

func runShow(cmd *cobra.Command, args []string) error {
	id := args[0]

	root, err := getWorkspaceRoot()
	if err != nil {
		return fmt.Errorf("failed to get workspace root: %w", err)
	}

	store, err := storage.NewStorage(root)
	if err != nil {
		return fmt.Errorf("failed to initialize storage: %w", err)
	}
	defer store.Close()

	decision, err := store.GetDecision(id)
	if err != nil {
		return fmt.Errorf("decision not found: %w", err)
	}

	// Print detailed decision info
	fmt.Printf("\n" + strings.Repeat("=", 60) + "\n")
	fmt.Printf("Decision: %s\n", decision.Title)
	fmt.Printf(strings.Repeat("=", 60) + "\n\n")

	fmt.Printf("ID: %s\n", decision.ID)
	fmt.Printf("Status: %s %s\n", getStatusIcon(decision.Status), decision.Status)
	fmt.Printf("Decided: %s", decision.DecidedAt.Format("2006-01-02 15:04:05"))
	if decision.DecidedBy != "" {
		fmt.Printf(" by %s", decision.DecidedBy)
	}
	fmt.Println()

	if decision.UpdatedAt.After(decision.DecidedAt) {
		fmt.Printf("Updated: %s\n", decision.UpdatedAt.Format("2006-01-02 15:04:05"))
	}

	if decision.RetiredAt != nil {
		fmt.Printf("Retired: %s\n", decision.RetiredAt.Format("2006-01-02 15:04:05"))
		if decision.RetirementReason != "" {
			fmt.Printf("Retirement Reason: %s\n", decision.RetirementReason)
		}
	}

	if decision.Reasoning != "" {
		fmt.Printf("\n📝 Reasoning:\n")
		fmt.Printf("%s\n", wrapText(decision.Reasoning, 60, "   "))
	}

	if len(decision.Assumptions) > 0 {
		fmt.Printf("\n💡 Assumptions:\n")
		for i, assumption := range decision.Assumptions {
			fmt.Printf("   %d. %s\n", i+1, assumption)
		}
	}

	if len(decision.ExpirationTriggers) > 0 {
		fmt.Printf("\n⏰ Expiration Triggers:\n")
		for i, trigger := range decision.ExpirationTriggers {
			fmt.Printf("   %d. %s\n", i+1, trigger)
		}
	}

	if len(decision.LinkedResources) > 0 {
		fmt.Printf("\n🔗 Linked Resources:\n")
		for _, resource := range decision.LinkedResources {
			fmt.Printf("   • %s\n", resource)
		}
	}

	if len(decision.Tags) > 0 {
		fmt.Printf("\n🏷️  Tags: %s\n", strings.Join(decision.Tags, ", "))
	}

	// Show alerts for this decision
	alerts, err := store.GetAlerts(id, false)
	if err == nil && len(alerts) > 0 {
		fmt.Printf("\n⚠️  Alerts (%d):\n", len(alerts))
		for _, alert := range alerts {
			ackStatus := ""
			if alert.Acknowledged {
				ackStatus = " [✓]"
			}
			fmt.Printf("   • [%s] %s%s\n", alert.Severity, alert.Trigger, ackStatus)
			if alert.Evidence != "" {
				fmt.Printf("     Evidence: %s\n", alert.Evidence)
			}
		}
	}

	fmt.Println()

	return nil
}

func wrapText(text string, width int, prefix string) string {
	if len(text) <= width {
		return prefix + text
	}

	var result strings.Builder
	words := strings.Fields(text)
	lineLen := 0

	result.WriteString(prefix)

	for _, word := range words {
		if lineLen+len(word)+1 > width {
			result.WriteString("\n")
			result.WriteString(prefix)
			lineLen = 0
		}

		if lineLen > 0 {
			result.WriteString(" ")
			lineLen++
		}

		result.WriteString(word)
		lineLen += len(word)
	}

	return result.String()
}
