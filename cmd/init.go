package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize DDT in current directory",
	Long:  "Create .ddt directory and initialize the decision tracking system",
	RunE:  runInit,
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func runInit(cmd *cobra.Command, args []string) error {
	root, err := getWorkspaceRoot()
	if err != nil {
		return fmt.Errorf("failed to get workspace root: %w", err)
	}

	store, err := storage.NewStorage(root)
	if err != nil {
		return fmt.Errorf("failed to initialize storage: %w", err)
	}
	defer store.Close()

	printSuccess("DDT initialized in %s", root)
	fmt.Printf("   Database: %s/.ddt/decisions.db\n", root)
	fmt.Printf("   Backups: %s/.ddt/backups/\n", root)
	fmt.Println()
	fmt.Println("📝 Get started:")
	fmt.Println("   ddt add \"Your first decision\" - Log a new decision")
	fmt.Println("   ddt list                      - View all decisions")
	fmt.Println("   ddt check                     - Check for outdated decisions")

	return nil
}
