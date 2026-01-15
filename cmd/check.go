package cmd

import (
	"fmt"

	"github.com/nidhalelbkir/ddt/detector"
	"github.com/nidhalelbkir/ddt/storage"
	"github.com/spf13/cobra"
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check for outdated decisions",
	Long:  "Run the debt detection engine to find decisions that may be outdated",
	RunE:  runCheck,
}

func init() {
	rootCmd.AddCommand(checkCmd)
}

func runCheck(cmd *cobra.Command, args []string) error {
	root, err := getWorkspaceRoot()
	if err != nil {
		return fmt.Errorf("failed to get workspace root: %w", err)
	}

	store, err := storage.NewStorage(root)
	if err != nil {
		return fmt.Errorf("failed to initialize storage: %w", err)
	}
	defer store.Close()

	fmt.Println("\n🔍 Checking decisions for expiration triggers...")

	det := detector.NewDetector(store, root)
	alerts, err := det.CheckAll()
	if err != nil {
		return fmt.Errorf("failed to check decisions: %w", err)
	}

	if len(alerts) == 0 {
		printSuccess("All decisions are up to date! ✨")
		return nil
	}

	fmt.Printf("\n⚠️  Found %d potential issues:\n\n", len(alerts))

	for _, alert := range alerts {
		decision, err := store.GetDecision(alert.DecisionID)
		if err != nil {
			continue
		}

		severityIcon := "⚠️"
		if alert.Severity == "critical" {
			severityIcon = "🚨"
		}

		fmt.Printf("%s Decision [%s] may be outdated\n", severityIcon, alert.DecisionID)
		fmt.Printf("   Title: %s\n", decision.Title)
		fmt.Printf("   Trigger: %s\n", alert.Trigger)
		if alert.Evidence != "" {
			fmt.Printf("   Evidence: %s\n", alert.Evidence)
		}
		fmt.Printf("   Detected: %s\n", alert.DetectedAt.Format("2006-01-02 15:04"))
		fmt.Println()
	}

	fmt.Printf("💡 Review these decisions with: ddt show [decision-id]\n")
	fmt.Printf("💡 Update a decision with: ddt update [decision-id]\n")
	fmt.Printf("💡 Retire a decision with: ddt retire [decision-id]\n")

	return nil
}
