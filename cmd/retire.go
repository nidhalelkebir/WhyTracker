package cmd

import (
	"fmt"
	"time"

	"github.com/manifoldco/promptui"
	"github.com/nidhalelbkir/ddt/storage"
	"github.com/spf13/cobra"
)

var retireCmd = &cobra.Command{
	Use:   "retire [decision-id]",
	Short: "Retire a decision",
	Long:  "Mark a decision as retired/obsolete with a reason",
	Args:  cobra.ExactArgs(1),
	RunE:  runRetire,
}

func init() {
	rootCmd.AddCommand(retireCmd)
}

func runRetire(cmd *cobra.Command, args []string) error {
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

	if decision.Status == "retired" {
		printWarning("Decision %s is already retired", id)
		return nil
	}

	// Prompt for retirement reason
	prompt := promptui.Prompt{
		Label: "Retirement reason",
	}

	reason, err := prompt.Run()
	if err != nil {
		return err
	}

	// Update decision
	now := time.Now()
	decision.Status = "retired"
	decision.RetiredAt = &now
	decision.RetirementReason = reason
	decision.UpdatedAt = now

	if err := store.UpdateDecision(decision); err != nil {
		return fmt.Errorf("failed to retire decision: %w", err)
	}

	// Export updated decision
	if err := store.ExportToYAML(decision); err != nil {
		printWarning("Failed to export to YAML: %v", err)
	}

	printSuccess("Decision %s retired", id)
	fmt.Printf("   Reason: %s\n", reason)

	return nil
}
