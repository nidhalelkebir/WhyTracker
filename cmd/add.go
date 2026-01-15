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

var addCmd = &cobra.Command{
	Use:   "add [decision title]",
	Short: "Add a new decision",
	Long:  "Add a new technical decision with context and assumptions",
	Args:  cobra.MinimumNArgs(1),
	RunE:  runAdd,
}

func init() {
	rootCmd.AddCommand(addCmd)
}

func runAdd(cmd *cobra.Command, args []string) error {
	title := strings.Join(args, " ")

	// Initialize storage
	root, err := getWorkspaceRoot()
	if err != nil {
		return fmt.Errorf("failed to get workspace root: %w", err)
	}

	store, err := storage.NewStorage(root)
	if err != nil {
		return fmt.Errorf("failed to initialize storage: %w", err)
	}
	defer store.Close()

	// Interactive prompts
	fmt.Printf("\n📝 Adding decision: %s\n\n", title)

	reasoning, err := promptInput("Why was this decision made?", "")
	if err != nil {
		return err
	}

	decidedBy, err := promptInput("Who decided? (person/team)", "")
	if err != nil {
		return err
	}

	// Assumptions (multi-line)
	fmt.Println("\nAssumptions (enter each assumption, empty line to finish):")
	assumptions := promptMultiLine()

	// Expiration triggers
	fmt.Println("\nExpiration triggers (when should this be reconsidered?):")
	triggers := promptMultiLine()

	// Linked resources
	fmt.Println("\nLinked resources (commit hash, URL, doc path - optional):")
	resources := promptMultiLine()

	// Tags
	fmt.Println("\nTags (optional, comma-separated):")
	tagsInput, err := promptInput("Tags", "")
	if err != nil {
		return err
	}
	
	var tags []string
	if tagsInput != "" {
		tags = strings.Split(tagsInput, ",")
		for i := range tags {
			tags[i] = strings.TrimSpace(tags[i])
		}
	}

	// Generate ID
	now := time.Now()
	id := fmt.Sprintf("ddt-%s-%03d", now.Format("2006-01-02"), now.Unix()%1000)

	decision := &models.Decision{
		ID:                 id,
		Title:              title,
		Reasoning:          reasoning,
		DecidedBy:          decidedBy,
		DecidedAt:          now,
		Assumptions:        assumptions,
		ExpirationTriggers: triggers,
		LinkedResources:    resources,
		Tags:               tags,
		Status:             "active",
		UpdatedAt:          now,
	}

	// Save to database
	if err := store.SaveDecision(decision); err != nil {
		return fmt.Errorf("failed to save decision: %w", err)
	}

	// Export to YAML backup
	if err := store.ExportToYAML(decision); err != nil {
		printWarning("Failed to export to YAML: %v", err)
	}

	printSuccess("Decision saved [ID: %s]", id)
	fmt.Printf("\n📄 Decision Details:\n")
	fmt.Printf("   Title: %s\n", title)
	fmt.Printf("   ID: %s\n", id)
	fmt.Printf("   Decided by: %s\n", decidedBy)
	fmt.Printf("   Assumptions: %d\n", len(assumptions))
	fmt.Printf("   Expiration triggers: %d\n", len(triggers))
	
	return nil
}

func promptInput(label, defaultValue string) (string, error) {
	prompt := promptui.Prompt{
		Label:   label,
		Default: defaultValue,
	}

	result, err := prompt.Run()
	if err != nil {
		return "", err
	}

	return result, nil
}

func promptMultiLine() []string {
	var lines []string
	lineNum := 1

	for {
		prompt := promptui.Prompt{
			Label: fmt.Sprintf("  %d", lineNum),
		}

		result, err := prompt.Run()
		if err != nil {
			break
		}

		if result == "" {
			break
		}

		lines = append(lines, result)
		lineNum++
	}

	return lines
}
