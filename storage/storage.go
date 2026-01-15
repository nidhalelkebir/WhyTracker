package storage

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	"github.com/nidhalelbkir/ddt/models"
)

const (
	ddtDir     = ".ddt"
	dbFile     = "decisions.db"
	backupDir  = "backups"
)

type Storage struct {
	db       *sql.DB
	rootPath string
}

// NewStorage initializes the storage layer
func NewStorage(rootPath string) (*Storage, error) {
	ddtPath := filepath.Join(rootPath, ddtDir)
	
	// Create .ddt directory if it doesn't exist
	if err := os.MkdirAll(ddtPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create .ddt directory: %w", err)
	}

	// Create backups directory
	backupPath := filepath.Join(ddtPath, backupDir)
	if err := os.MkdirAll(backupPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create backups directory: %w", err)
	}

	dbPath := filepath.Join(ddtPath, dbFile)
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	s := &Storage{
		db:       db,
		rootPath: rootPath,
	}

	if err := s.initSchema(); err != nil {
		return nil, err
	}

	return s, nil
}

// initSchema creates the necessary database tables
func (s *Storage) initSchema() error {
	schema := `
	CREATE TABLE IF NOT EXISTS decisions (
		id TEXT PRIMARY KEY,
		title TEXT NOT NULL,
		reasoning TEXT,
		decided_by TEXT,
		decided_at DATETIME NOT NULL,
		assumptions TEXT, -- JSON array
		expiration_triggers TEXT, -- JSON array
		linked_resources TEXT, -- JSON array
		tags TEXT, -- JSON array
		status TEXT DEFAULT 'active',
		updated_at DATETIME,
		retired_at DATETIME,
		retirement_reason TEXT
	);

	CREATE TABLE IF NOT EXISTS alerts (
		id TEXT PRIMARY KEY,
		decision_id TEXT NOT NULL,
		trigger_text TEXT NOT NULL,
		evidence TEXT,
		severity TEXT DEFAULT 'warning',
		detected_at DATETIME NOT NULL,
		acknowledged BOOLEAN DEFAULT 0,
		FOREIGN KEY (decision_id) REFERENCES decisions(id)
	);

	CREATE TABLE IF NOT EXISTS decision_updates (
		id TEXT PRIMARY KEY,
		decision_id TEXT NOT NULL,
		updated_by TEXT,
		updated_at DATETIME NOT NULL,
		changes TEXT,
		reason TEXT,
		FOREIGN KEY (decision_id) REFERENCES decisions(id)
	);

	CREATE INDEX IF NOT EXISTS idx_decisions_status ON decisions(status);
	CREATE INDEX IF NOT EXISTS idx_decisions_decided_at ON decisions(decided_at);
	CREATE INDEX IF NOT EXISTS idx_alerts_decision_id ON alerts(decision_id);
	CREATE INDEX IF NOT EXISTS idx_alerts_acknowledged ON alerts(acknowledged);
	`

	_, err := s.db.Exec(schema)
	return err
}

// SaveDecision saves a decision to the database
func (s *Storage) SaveDecision(decision *models.Decision) error {
	query := `
	INSERT INTO decisions (
		id, title, reasoning, decided_by, decided_at,
		assumptions, expiration_triggers, linked_resources, tags,
		status, updated_at, retired_at, retirement_reason
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := s.db.Exec(
		query,
		decision.ID,
		decision.Title,
		decision.Reasoning,
		decision.DecidedBy,
		decision.DecidedAt,
		stringSliceToJSON(decision.Assumptions),
		stringSliceToJSON(decision.ExpirationTriggers),
		stringSliceToJSON(decision.LinkedResources),
		stringSliceToJSON(decision.Tags),
		decision.Status,
		decision.UpdatedAt,
		decision.RetiredAt,
		decision.RetirementReason,
	)

	return err
}

// GetDecision retrieves a decision by ID
func (s *Storage) GetDecision(id string) (*models.Decision, error) {
	query := `
	SELECT id, title, reasoning, decided_by, decided_at,
		   assumptions, expiration_triggers, linked_resources, tags,
		   status, updated_at, retired_at, retirement_reason
	FROM decisions
	WHERE id = ?
	`

	var decision models.Decision
	var assumptions, triggers, resources, tags string
	var retiredAt sql.NullTime
	var retirementReason sql.NullString

	err := s.db.QueryRow(query, id).Scan(
		&decision.ID,
		&decision.Title,
		&decision.Reasoning,
		&decision.DecidedBy,
		&decision.DecidedAt,
		&assumptions,
		&triggers,
		&resources,
		&tags,
		&decision.Status,
		&decision.UpdatedAt,
		&retiredAt,
		&retirementReason,
	)

	if err != nil {
		return nil, err
	}

	decision.Assumptions = jsonToStringSlice(assumptions)
	decision.ExpirationTriggers = jsonToStringSlice(triggers)
	decision.LinkedResources = jsonToStringSlice(resources)
	decision.Tags = jsonToStringSlice(tags)
	
	if retiredAt.Valid {
		decision.RetiredAt = &retiredAt.Time
	}
	if retirementReason.Valid {
		decision.RetirementReason = retirementReason.String
	}

	return &decision, nil
}

// ListDecisions retrieves all decisions with optional status filter
func (s *Storage) ListDecisions(status string, limit int) ([]*models.Decision, error) {
	query := `
	SELECT id, title, reasoning, decided_by, decided_at,
		   assumptions, expiration_triggers, linked_resources, tags,
		   status, updated_at, retired_at, retirement_reason
	FROM decisions
	`
	args := []interface{}{}

	if status != "" {
		query += " WHERE status = ?"
		args = append(args, status)
	}

	query += " ORDER BY decided_at DESC"

	if limit > 0 {
		query += " LIMIT ?"
		args = append(args, limit)
	}

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var decisions []*models.Decision
	for rows.Next() {
		var decision models.Decision
		var assumptions, triggers, resources, tags string
		var retiredAt sql.NullTime
		var retirementReason sql.NullString

		err := rows.Scan(
			&decision.ID,
			&decision.Title,
			&decision.Reasoning,
			&decision.DecidedBy,
			&decision.DecidedAt,
			&assumptions,
			&triggers,
			&resources,
			&tags,
			&decision.Status,
			&decision.UpdatedAt,
			&retiredAt,
			&retirementReason,
		)
		if err != nil {
			return nil, err
		}

		decision.Assumptions = jsonToStringSlice(assumptions)
		decision.ExpirationTriggers = jsonToStringSlice(triggers)
		decision.LinkedResources = jsonToStringSlice(resources)
		decision.Tags = jsonToStringSlice(tags)
		
		if retiredAt.Valid {
			decision.RetiredAt = &retiredAt.Time
		}
		if retirementReason.Valid {
			decision.RetirementReason = retirementReason.String
		}

		decisions = append(decisions, &decision)
	}

	return decisions, nil
}

// SearchDecisions searches for decisions by query string
func (s *Storage) SearchDecisions(query string) ([]*models.Decision, error) {
	searchQuery := `
	SELECT id, title, reasoning, decided_by, decided_at,
		   assumptions, expiration_triggers, linked_resources, tags,
		   status, updated_at, retired_at, retirement_reason
	FROM decisions
	WHERE title LIKE ? OR reasoning LIKE ? OR tags LIKE ?
	ORDER BY decided_at DESC
	`

	pattern := "%" + query + "%"
	rows, err := s.db.Query(searchQuery, pattern, pattern, pattern)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var decisions []*models.Decision
	for rows.Next() {
		var decision models.Decision
		var assumptions, triggers, resources, tags string
		var retiredAt sql.NullTime
		var retirementReason sql.NullString

		err := rows.Scan(
			&decision.ID,
			&decision.Title,
			&decision.Reasoning,
			&decision.DecidedBy,
			&decision.DecidedAt,
			&assumptions,
			&triggers,
			&resources,
			&tags,
			&decision.Status,
			&decision.UpdatedAt,
			&retiredAt,
			&retirementReason,
		)
		if err != nil {
			return nil, err
		}

		decision.Assumptions = jsonToStringSlice(assumptions)
		decision.ExpirationTriggers = jsonToStringSlice(triggers)
		decision.LinkedResources = jsonToStringSlice(resources)
		decision.Tags = jsonToStringSlice(tags)
		
		if retiredAt.Valid {
			decision.RetiredAt = &retiredAt.Time
		}
		if retirementReason.Valid {
			decision.RetirementReason = retirementReason.String
		}

		decisions = append(decisions, &decision)
	}

	return decisions, nil
}

// UpdateDecision updates an existing decision
func (s *Storage) UpdateDecision(decision *models.Decision) error {
	query := `
	UPDATE decisions SET
		title = ?, reasoning = ?, decided_by = ?,
		assumptions = ?, expiration_triggers = ?, linked_resources = ?, tags = ?,
		status = ?, updated_at = ?, retired_at = ?, retirement_reason = ?
	WHERE id = ?
	`

	_, err := s.db.Exec(
		query,
		decision.Title,
		decision.Reasoning,
		decision.DecidedBy,
		stringSliceToJSON(decision.Assumptions),
		stringSliceToJSON(decision.ExpirationTriggers),
		stringSliceToJSON(decision.LinkedResources),
		stringSliceToJSON(decision.Tags),
		decision.Status,
		decision.UpdatedAt,
		decision.RetiredAt,
		decision.RetirementReason,
		decision.ID,
	)

	return err
}

// SaveAlert saves an alert to the database
func (s *Storage) SaveAlert(alert *models.Alert) error {
	query := `
	INSERT INTO alerts (
		id, decision_id, trigger_text, evidence, severity, detected_at, acknowledged
	) VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	_, err := s.db.Exec(
		query,
		alert.ID,
		alert.DecisionID,
		alert.Trigger,
		alert.Evidence,
		alert.Severity,
		alert.DetectedAt,
		alert.Acknowledged,
	)

	return err
}

// GetAlerts retrieves alerts, optionally filtered by decision ID
func (s *Storage) GetAlerts(decisionID string, unacknowledgedOnly bool) ([]*models.Alert, error) {
	query := `
	SELECT id, decision_id, trigger_text, evidence, severity, detected_at, acknowledged
	FROM alerts
	WHERE 1=1
	`
	args := []interface{}{}

	if decisionID != "" {
		query += " AND decision_id = ?"
		args = append(args, decisionID)
	}

	if unacknowledgedOnly {
		query += " AND acknowledged = 0"
	}

	query += " ORDER BY detected_at DESC"

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var alerts []*models.Alert
	for rows.Next() {
		var alert models.Alert
		err := rows.Scan(
			&alert.ID,
			&alert.DecisionID,
			&alert.Trigger,
			&alert.Evidence,
			&alert.Severity,
			&alert.DetectedAt,
			&alert.Acknowledged,
		)
		if err != nil {
			return nil, err
		}
		alerts = append(alerts, &alert)
	}

	return alerts, nil
}

// Close closes the database connection
func (s *Storage) Close() error {
	return s.db.Close()
}

// GetDDTPath returns the .ddt directory path
func (s *Storage) GetDDTPath() string {
	return filepath.Join(s.rootPath, ddtDir)
}

// GetBackupPath returns the backups directory path
func (s *Storage) GetBackupPath() string {
	return filepath.Join(s.rootPath, ddtDir, backupDir)
}
