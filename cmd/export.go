package cmd

import (
	"fmt"

	"github.com/nidhalelbkir/ddt/storage"
	"github.com/spf13/cobra"
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export decisions to YAML",
	Long:  "Export all active decisions to YAML backup files",
	RunE:  runExport,
}

func init() {
	rootCmd.AddCommand(exportCmd)
}

func runExport(cmd *cobra.Command, args []string) error {
	root, err := getWorkspaceRoot()
	if err != nil {
		return fmt.Errorf("failed to get workspace root: %w", err)
	}

	store, err := storage.NewStorage(root)
	if err != nil {
		return fmt.Errorf("failed to initialize storage: %w", err)
	}
	defer store.Close()

	fmt.Println("\n📦 Exporting decisions to YAML...")

	if err := store.ExportAllToYAML(); err != nil {
		return fmt.Errorf("failed to export decisions: %w", err)
	}

	printSuccess("All decisions exported to %s", store.GetBackupPath())

	return nil
}
