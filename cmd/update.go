package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/manifoldco/promptui"
	"github.com/nidhalelbkir/ddt/models"
	"github.com/nidhalelbkir/ddt/storage"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update [decision-id]",
	Short: "Update a decision",
	Long:  "Log a change or revision to an existing decision",
	Args:  cobra.ExactArgs(1),
	RunE:  runUpdate,
}

func init() {
	rootCmd.AddCommand(updateCmd)
}

func runUpdate(cmd *cobra.Command, args []string) error {
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

	fmt.Printf("\n📝 Updating decision: %s\n\n", decision.Title)

	// Select what to update
	selectPrompt := promptui.Select{
		Label: "What would you like to update?",
		Items: []string{
			"Reasoning",
			"Assumptions",
			"Expiration Triggers",
			"Linked Resources",
			"Tags",
			"All fields",
		},
	}

	_, selection, err := selectPrompt.Run()
	if err != nil {
		return err
	}

	updateReason, err := promptInput("Reason for update", "")
	if err != nil {
		return err
	}

	changes := []string{}

	switch selection {
	case "Reasoning":
		newReasoning, err := promptInput("New reasoning", decision.Reasoning)
		if err != nil {
			return err
		}
		if newReasoning != decision.Reasoning {
			decision.Reasoning = newReasoning
			changes = append(changes, "Updated reasoning")
		}

	case "Assumptions":
		fmt.Println("\nCurrent assumptions:")
		for i, a := range decision.Assumptions {
			fmt.Printf("  %d. %s\n", i+1, a)
		}
		fmt.Println("\nNew assumptions (empty line to finish):")
		newAssumptions := promptMultiLine()
		if len(newAssumptions) > 0 {
			decision.Assumptions = newAssumptions
			changes = append(changes, "Updated assumptions")
		}

	case "Expiration Triggers":
		fmt.Println("\nCurrent triggers:")
		for i, t := range decision.ExpirationTriggers {
			fmt.Printf("  %d. %s\n", i+1, t)
		}
		fmt.Println("\nNew expiration triggers (empty line to finish):")
		newTriggers := promptMultiLine()
		if len(newTriggers) > 0 {
			decision.ExpirationTriggers = newTriggers
			changes = append(changes, "Updated expiration triggers")
		}

	case "Linked Resources":
		fmt.Println("\nCurrent resources:")
		for i, r := range decision.LinkedResources {
			fmt.Printf("  %d. %s\n", i+1, r)
		}
		fmt.Println("\nNew linked resources (empty line to finish):")
		newResources := promptMultiLine()
		if len(newResources) > 0 {
			decision.LinkedResources = newResources
			changes = append(changes, "Updated linked resources")
		}

	case "Tags":
		currentTags := strings.Join(decision.Tags, ", ")
		newTagsStr, err := promptInput("Tags (comma-separated)", currentTags)
		if err != nil {
			return err
		}
		if newTagsStr != currentTags {
			var newTags []string
			if newTagsStr != "" {
				newTags = strings.Split(newTagsStr, ",")
				for i := range newTags {
					newTags[i] = strings.TrimSpace(newTags[i])
				}
			}
			decision.Tags = newTags
			changes = append(changes, "Updated tags")
		}

	case "All fields":
		// Update all fields
		newReasoning, _ := promptInput("Reasoning", decision.Reasoning)
		if newReasoning != decision.Reasoning {
			decision.Reasoning = newReasoning
			changes = append(changes, "reasoning")
		}

		fmt.Println("\nAssumptions (empty line to finish):")
		newAssumptions := promptMultiLine()
		if len(newAssumptions) > 0 {
			decision.Assumptions = newAssumptions
			changes = append(changes, "assumptions")
		}

		fmt.Println("\nExpiration triggers (empty line to finish):")
		newTriggers := promptMultiLine()
		if len(newTriggers) > 0 {
			decision.ExpirationTriggers = newTriggers
			changes = append(changes, "triggers")
		}
	}

	if len(changes) == 0 {
		printInfo("No changes made")
		return nil
	}

	// Update decision
	now := time.Now()
	decision.Status = "updated"
	decision.UpdatedAt = now

	if err := store.UpdateDecision(decision); err != nil {
		return fmt.Errorf("failed to update decision: %w", err)
	}

	// Log the update
	update := &models.DecisionUpdate{
		ID:         fmt.Sprintf("upd-%d", now.Unix()),
		DecisionID: decision.ID,
		UpdatedAt:  now,
		Changes:    strings.Join(changes, ", "),
		Reason:     updateReason,
	}

	// Export updated decision
	if err := store.ExportToYAML(decision); err != nil {
		printWarning("Failed to export to YAML: %v", err)
	}

	printSuccess("Decision %s updated", id)
	fmt.Printf("   Changes: %s\n", update.Changes)
	fmt.Printf("   Reason: %s\n", updateReason)

	return nil
}
