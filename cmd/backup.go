package cmd

import (
	"fmt"

	"github.com/nidhalelbkir/ddt/storage"
	"github.com/spf13/cobra"
)

var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Create a backup of all decisions",
	Long:  "Create a timestamped backup file of all decisions",
	RunE:  runBackup,
}

func init() {
	rootCmd.AddCommand(backupCmd)
}

func runBackup(cmd *cobra.Command, args []string) error {
	root, err := getWorkspaceRoot()
	if err != nil {
		return fmt.Errorf("failed to get workspace root: %w", err)
	}

	store, err := storage.NewStorage(root)
	if err != nil {
		return fmt.Errorf("failed to initialize storage: %w", err)
	}
	defer store.Close()

	fmt.Println("\n💾 Creating backup...")

	if err := store.BackupDatabase(); err != nil {
		return fmt.Errorf("failed to create backup: %w", err)
	}

	printSuccess("Backup created in %s", store.GetBackupPath())

	return nil
}
