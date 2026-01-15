package models

import (
	"time"
)

// Decision represents a technical decision logged in the system
type Decision struct {
	ID                 string    `json:"id" yaml:"id"`
	Title              string    `json:"title" yaml:"title"`
	Reasoning          string    `json:"reasoning" yaml:"reasoning"`
	DecidedBy          string    `json:"decided_by" yaml:"decided_by"`
	DecidedAt          time.Time `json:"decided_at" yaml:"decided_at"`
	Assumptions        []string  `json:"assumptions" yaml:"assumptions"`
	ExpirationTriggers []string  `json:"expiration_triggers" yaml:"expiration_triggers"`
	LinkedResources    []string  `json:"linked_resources" yaml:"linked_resources"`
	Tags               []string  `json:"tags" yaml:"tags"`
	Status             string    `json:"status" yaml:"status"` // active, updated, retired
	UpdatedAt          time.Time `json:"updated_at" yaml:"updated_at"`
	RetiredAt          *time.Time `json:"retired_at,omitempty" yaml:"retired_at,omitempty"`
	RetirementReason   string    `json:"retirement_reason,omitempty" yaml:"retirement_reason,omitempty"`
}

// Alert represents a triggered warning about a potentially outdated decision
type Alert struct {
	ID            string    `json:"id" yaml:"id"`
	DecisionID    string    `json:"decision_id" yaml:"decision_id"`
	Trigger       string    `json:"trigger" yaml:"trigger"`
	Evidence      string    `json:"evidence" yaml:"evidence"`
	Severity      string    `json:"severity" yaml:"severity"` // warning, critical
	DetectedAt    time.Time `json:"detected_at" yaml:"detected_at"`
	Acknowledged  bool      `json:"acknowledged" yaml:"acknowledged"`
}

// DecisionUpdate represents a revision or update to a decision
type DecisionUpdate struct {
	ID         string    `json:"id" yaml:"id"`
	DecisionID string    `json:"decision_id" yaml:"decision_id"`
	UpdatedBy  string    `json:"updated_by" yaml:"updated_by"`
	UpdatedAt  time.Time `json:"updated_at" yaml:"updated_at"`
	Changes    string    `json:"changes" yaml:"changes"`
	Reason     string    `json:"reason" yaml:"reason"`
}
