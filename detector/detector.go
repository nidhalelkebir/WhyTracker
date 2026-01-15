package detector

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/nidhalelbkir/ddt/models"
	"github.com/nidhalelbkir/ddt/storage"
)

// Detector checks for outdated decisions based on expiration triggers
type Detector struct {
	store    *storage.Storage
	rootPath string
}

// NewDetector creates a new detector instance
func NewDetector(store *storage.Storage, rootPath string) *Detector {
	return &Detector{
		store:    store,
		rootPath: rootPath,
	}
}

// CheckAll checks all active decisions for expiration
func (d *Detector) CheckAll() ([]*models.Alert, error) {
	decisions, err := d.store.ListDecisions("active", 0)
	if err != nil {
		return nil, err
	}

	var alerts []*models.Alert

	for _, decision := range decisions {
		decisionAlerts := d.checkDecision(decision)
		alerts = append(alerts, decisionAlerts...)
	}

	// Save alerts to database
	for _, alert := range alerts {
		if err := d.store.SaveAlert(alert); err != nil {
			return nil, fmt.Errorf("failed to save alert: %w", err)
		}
	}

	return alerts, nil
}

// checkDecision checks a single decision for expiration
func (d *Detector) checkDecision(decision *models.Decision) []*models.Alert {
	var alerts []*models.Alert

	for _, trigger := range decision.ExpirationTriggers {
		if alert := d.evaluateTrigger(decision, trigger); alert != nil {
			alerts = append(alerts, alert)
		}
	}

	return alerts
}

// evaluateTrigger evaluates a single expiration trigger
func (d *Detector) evaluateTrigger(decision *models.Decision, trigger string) *models.Alert {
	triggerLower := strings.ToLower(trigger)

	// Check for file size triggers
	if strings.Contains(triggerLower, "size") || strings.Contains(triggerLower, "exceed") {
		if alert := d.checkFileSizeTrigger(decision, trigger); alert != nil {
			return alert
		}
	}

	// Check for file existence triggers
	if strings.Contains(triggerLower, "when") && (strings.Contains(triggerLower, "add") || 
		strings.Contains(triggerLower, "create") || strings.Contains(triggerLower, "exist")) {
		if alert := d.checkFileExistenceTrigger(decision, trigger); alert != nil {
			return alert
		}
	}

	// Check for count/number triggers
	if strings.Contains(triggerLower, "more than") || strings.Contains(triggerLower, ">") {
		if alert := d.checkCountTrigger(decision, trigger); alert != nil {
			return alert
		}
	}

	// Check for time-based triggers
	if strings.Contains(triggerLower, "month") || strings.Contains(triggerLower, "year") || 
		strings.Contains(triggerLower, "days") {
		if alert := d.checkTimeTrigger(decision, trigger); alert != nil {
			return alert
		}
	}

	// Check for keyword presence in codebase
	if strings.Contains(triggerLower, "when") && (strings.Contains(triggerLower, "use") || 
		strings.Contains(triggerLower, "implement") || strings.Contains(triggerLower, "require")) {
		if alert := d.checkKeywordTrigger(decision, trigger); alert != nil {
			return alert
		}
	}

	return nil
}

// checkFileSizeTrigger checks if file size exceeds threshold
func (d *Detector) checkFileSizeTrigger(decision *models.Decision, trigger string) *models.Alert {
	// Extract size threshold from trigger (e.g., "Data size > 1GB")
	re := regexp.MustCompile(`(\d+)\s*(GB|MB|KB)`)
	matches := re.FindStringSubmatch(trigger)
	
	if len(matches) < 3 {
		return nil
	}

	threshold, _ := strconv.ParseInt(matches[1], 10, 64)
	unit := matches[2]

	// Convert to bytes
	var thresholdBytes int64
	switch unit {
	case "GB":
		thresholdBytes = threshold * 1024 * 1024 * 1024
	case "MB":
		thresholdBytes = threshold * 1024 * 1024
	case "KB":
		thresholdBytes = threshold * 1024
	}

	// Check database files, config files, etc.
	var totalSize int64
	filepath.Walk(d.rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() && (strings.HasSuffix(path, ".db") || 
			strings.HasSuffix(path, ".sqlite") ||
			strings.HasSuffix(path, ".json") ||
			strings.HasSuffix(path, ".yaml")) {
			totalSize += info.Size()
		}
		return nil
	})

	if totalSize > thresholdBytes {
		return &models.Alert{
			ID:           fmt.Sprintf("alert-%d", time.Now().Unix()),
			DecisionID:   decision.ID,
			Trigger:      trigger,
			Evidence:     fmt.Sprintf("Current data size: %.2f %s", float64(totalSize)/float64(thresholdBytes)*float64(threshold), unit),
			Severity:     "warning",
			DetectedAt:   time.Now(),
			Acknowledged: false,
		}
	}

	return nil
}

// checkFileExistenceTrigger checks if certain files exist
func (d *Detector) checkFileExistenceTrigger(decision *models.Decision, trigger string) *models.Alert {
	// Extract keywords from trigger
	keywords := extractKeywords(trigger)

	for _, keyword := range keywords {
		// Search for files matching the keyword
		var foundFiles []string
		filepath.Walk(d.rootPath, func(path string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() {
				return nil
			}
			if strings.Contains(strings.ToLower(filepath.Base(path)), keyword) ||
				strings.Contains(strings.ToLower(path), keyword) {
				foundFiles = append(foundFiles, path)
			}
			return nil
		})

		if len(foundFiles) > 0 {
			return &models.Alert{
				ID:           fmt.Sprintf("alert-%d", time.Now().Unix()),
				DecisionID:   decision.ID,
				Trigger:      trigger,
				Evidence:     fmt.Sprintf("Found files: %s", strings.Join(foundFiles, ", ")),
				Severity:     "warning",
				DetectedAt:   time.Now(),
				Acknowledged: false,
			}
		}
	}

	return nil
}

// checkCountTrigger checks if counts exceed thresholds
func (d *Detector) checkCountTrigger(decision *models.Decision, trigger string) *models.Alert {
	// Extract number from trigger
	re := regexp.MustCompile(`more than (\d+)|> ?(\d+)`)
	matches := re.FindStringSubmatch(trigger)
	
	if len(matches) < 2 {
		return nil
	}

	var threshold int
	if matches[1] != "" {
		threshold, _ = strconv.Atoi(matches[1])
	} else if matches[2] != "" {
		threshold, _ = strconv.Atoi(matches[2])
	}

	// Count relevant items based on trigger context
	keywords := extractKeywords(trigger)
	for _, keyword := range keywords {
		count := d.countOccurrences(keyword)
		if count > threshold {
			return &models.Alert{
				ID:           fmt.Sprintf("alert-%d", time.Now().Unix()),
				DecisionID:   decision.ID,
				Trigger:      trigger,
				Evidence:     fmt.Sprintf("Found %d occurrences of '%s' (threshold: %d)", count, keyword, threshold),
				Severity:     "warning",
				DetectedAt:   time.Now(),
				Acknowledged: false,
			}
		}
	}

	return nil
}

// checkTimeTrigger checks time-based expiration
func (d *Detector) checkTimeTrigger(decision *models.Decision, trigger string) *models.Alert {
	// Extract duration from trigger
	re := regexp.MustCompile(`(\d+)\s*(month|year|day)s?`)
	matches := re.FindStringSubmatch(strings.ToLower(trigger))
	
	if len(matches) < 3 {
		return nil
	}

	amount, _ := strconv.Atoi(matches[1])
	unit := matches[2]

	var duration time.Duration
	switch unit {
	case "day":
		duration = time.Duration(amount) * 24 * time.Hour
	case "month":
		duration = time.Duration(amount*30) * 24 * time.Hour
	case "year":
		duration = time.Duration(amount*365) * 24 * time.Hour
	}

	if time.Since(decision.DecidedAt) > duration {
		return &models.Alert{
			ID:           fmt.Sprintf("alert-%d", time.Now().Unix()),
			DecisionID:   decision.ID,
			Trigger:      trigger,
			Evidence:     fmt.Sprintf("Decision is %.0f days old (threshold: %.0f days)", 
				time.Since(decision.DecidedAt).Hours()/24, duration.Hours()/24),
			Severity:     "warning",
			DetectedAt:   time.Now(),
			Acknowledged: false,
		}
	}

	return nil
}

// checkKeywordTrigger checks for keyword presence in codebase
func (d *Detector) checkKeywordTrigger(decision *models.Decision, trigger string) *models.Alert {
	keywords := extractKeywords(trigger)

	for _, keyword := range keywords {
		if d.countOccurrences(keyword) > 0 {
			return &models.Alert{
				ID:           fmt.Sprintf("alert-%d", time.Now().Unix()),
				DecisionID:   decision.ID,
				Trigger:      trigger,
				Evidence:     fmt.Sprintf("Found keyword '%s' in codebase", keyword),
				Severity:     "warning",
				DetectedAt:   time.Now(),
				Acknowledged: false,
			}
		}
	}

	return nil
}

// countOccurrences counts occurrences of a keyword in the codebase
func (d *Detector) countOccurrences(keyword string) int {
	count := 0
	keywordLower := strings.ToLower(keyword)

	filepath.Walk(d.rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}

		// Skip binary files and hidden directories
		if strings.Contains(path, "/.") || strings.Contains(path, "\\.") {
			return nil
		}

		// Check common source code extensions
		ext := filepath.Ext(path)
		validExts := []string{".go", ".js", ".ts", ".py", ".java", ".rb", ".rs", 
			".c", ".cpp", ".h", ".hpp", ".cs", ".php", ".swift", ".kt"}
		
		isValid := false
		for _, validExt := range validExts {
			if ext == validExt {
				isValid = true
				break
			}
		}

		if !isValid {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return nil
		}

		contentLower := strings.ToLower(string(content))
		count += strings.Count(contentLower, keywordLower)

		return nil
	})

	return count
}

// extractKeywords extracts meaningful keywords from a trigger string
func extractKeywords(trigger string) []string {
	// Remove common words
	stopWords := []string{"when", "if", "the", "a", "an", "is", "are", "was", "were", 
		"add", "adds", "added", "create", "creates", "created", "require", "requires", 
		"use", "uses", "implement", "implements", "need", "needs"}

	words := strings.Fields(strings.ToLower(trigger))
	var keywords []string

	for _, word := range words {
		// Remove punctuation
		word = strings.Trim(word, ".,!?;:'\"()[]{}><")
		
		// Skip stop words and short words
		if len(word) < 3 {
			continue
		}

		isStopWord := false
		for _, sw := range stopWords {
			if word == sw {
				isStopWord = true
				break
			}
		}

		if !isStopWord {
			keywords = append(keywords, word)
		}
	}

	return keywords
}
