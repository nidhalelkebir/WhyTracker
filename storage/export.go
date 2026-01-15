package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
	"github.com/nidhalelbkir/ddt/models"
)

// ExportToYAML exports a decision to a YAML file in the backup directory
func (s *Storage) ExportToYAML(decision *models.Decision) error {
	backupPath := s.GetBackupPath()
	filename := fmt.Sprintf("%s.yaml", decision.ID)
	filePath := filepath.Join(backupPath, filename)

	data, err := yaml.Marshal(decision)
	if err != nil {
		return fmt.Errorf("failed to marshal decision to YAML: %w", err)
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write YAML file: %w", err)
	}

	return nil
}

// ExportAllToYAML exports all active decisions to YAML files
func (s *Storage) ExportAllToYAML() error {
	decisions, err := s.ListDecisions("active", 0)
	if err != nil {
		return err
	}

	for _, decision := range decisions {
		if err := s.ExportToYAML(decision); err != nil {
			return err
		}
	}

	return nil
}

// ImportFromYAML imports a decision from a YAML file
func (s *Storage) ImportFromYAML(filePath string) (*models.Decision, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read YAML file: %w", err)
	}

	var decision models.Decision
	if err := yaml.Unmarshal(data, &decision); err != nil {
		return nil, fmt.Errorf("failed to unmarshal YAML: %w", err)
	}

	return &decision, nil
}

// BackupDatabase creates a timestamped backup of all decisions
func (s *Storage) BackupDatabase() error {
	timestamp := time.Now().Format("20060102-150405")
	backupFile := filepath.Join(s.GetBackupPath(), fmt.Sprintf("backup-%s.yaml", timestamp))

	decisions, err := s.ListDecisions("", 0)
	if err != nil {
		return err
	}

	type Backup struct {
		Timestamp time.Time          `yaml:"timestamp"`
		Decisions []*models.Decision `yaml:"decisions"`
	}

	backup := Backup{
		Timestamp: time.Now(),
		Decisions: decisions,
	}

	data, err := yaml.Marshal(backup)
	if err != nil {
		return fmt.Errorf("failed to marshal backup: %w", err)
	}

	if err := os.WriteFile(backupFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write backup file: %w", err)
	}

	return nil
}
