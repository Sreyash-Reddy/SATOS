package verifier

import (
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// Result represents a verification result
type Result struct {
	Success   bool
	ExitCode  int
	Output    string
	Duration  time.Duration
	Timestamp time.Time
}

// Verifier manages verification of task outputs
type Verifier struct {
	Path       string
	Command    string
	FailedRuns int
	Flagged    bool
}

// NewVerifier creates a new verifier
func NewVerifier(command string) *Verifier {
	return &Verifier{
		Command: command,
	}
}

// Run executes the verification command
func (v *Verifier) Run() (*Result, error) {
	start := time.Now()

	cmd := exec.Command("bash", "-c", v.Command)
	output, err := cmd.CombinedOutput()

	duration := time.Since(start)
	result := &Result{
		Duration:  duration,
		Timestamp: time.Now(),
		Output:    string(output),
	}

	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitErr.ExitCode()
		} else {
			result.ExitCode = -1
		}
		result.Success = result.ExitCode == 0
	} else {
		result.ExitCode = 0
		result.Success = true
	}

	v.FailedRuns = 0
	if !result.Success {
		v.FailedRuns = 1
	}

	// Check if verifier should be flagged (3+ consecutive failures handled externally)
	v.Flagged = v.FailedRuns >= 3

	return result, nil
}

// ParseTaskFile parses a task markdown file
func ParseTaskFile(content string) (*TaskFile, error) {
	tf := &TaskFile{}

	lines := strings.Split(content, "\n")
	var currentSection string

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Parse section header
		if strings.HasPrefix(trimmed, "##") {
			// Get section content after ##
			sectionPart := strings.TrimPrefix(trimmed, "##")

			// Check if value is on same line: ## Status: queued
			if colonIdx := strings.Index(sectionPart, ":"); colonIdx > 0 {
				sectionName := strings.TrimSpace(sectionPart[:colonIdx])
				sectionValue := strings.TrimSpace(sectionPart[colonIdx+1:])
				currentSection = strings.ToLower(sectionName)
				// Set value if on same line
				switch currentSection {
				case "status":
					tf.Status = sectionValue
				case "output":
					tf.Output = sectionValue
				case "depends-on":
					tf.DependsOn = sectionValue
				case "sprint":
					tf.Sprint = sectionValue
				case "priority":
					tf.Priority = sectionValue
				case "assignee":
					tf.Assignee = sectionValue
				}
			} else {
				// No colon - just section name, value on next line
				currentSection = strings.ToLower(strings.TrimSpace(sectionPart))
			}
			continue
		}

		// Track code blocks
		if strings.HasPrefix(trimmed, "```") {
			continue
		}

		// Skip empty lines and non-section lines
		if trimmed == "" || currentSection == "" {
			continue
		}

		// Command and Verify are multi-line (code blocks)
		// For now, collect lines between section headers
		if currentSection == "command" {
			tf.Command += trimmed + "\n"
		} else if currentSection == "verify" {
			tf.Verify += trimmed + "\n"
		} else {
			// Single value sections
			switch currentSection {
			case "status":
				tf.Status = trimmed
			case "output":
				tf.Output = trimmed
			case "depends-on":
				tf.DependsOn = trimmed
			case "sprint":
				tf.Sprint = trimmed
			case "priority":
				tf.Priority = trimmed
			case "assignee":
				tf.Assignee = trimmed
			}
		}
	}

	// Clean up
	tf.Command = strings.Trim(tf.Command, "\n ")
	tf.Verify = strings.Trim(tf.Verify, "\n ")

	return tf, nil
}

// TaskFile represents a parsed task file
type TaskFile struct {
	Status    string
	Command   string
	Output    string
	Verify    string
	DependsOn string
	Sprint    string
	Priority  string
	Assignee  string
}

// Validate validates the task file
func (t *TaskFile) Validate() error {
	if t.Command == "" {
		return fmt.Errorf("command is required")
	}
	if t.Output == "" {
		return fmt.Errorf("output is required")
	}
	return nil
}

// ToTask converts TaskFile to a state Task
func (t *TaskFile) ToTask(taskID string) map[string]interface{} {
	return map[string]interface{}{
		"task_id":     taskID,
		"status":      t.Status,
		"command":     t.Command,
		"output":      t.Output,
		"verify":      t.Verify,
		"depends_on":  t.DependsOn,
		"sprint":      t.Sprint,
		"priority":    t.Priority,
		"assignee":    t.Assignee,
	}
}

// ReliabilityChecker tracks verifier reliability
type ReliabilityChecker struct {
	Failures map[string]int
	Total    map[string]int
}

// NewReliabilityChecker creates a new reliability checker
func NewReliabilityChecker() *ReliabilityChecker {
	return &ReliabilityChecker{
		Failures: make(map[string]int),
		Total:    make(map[string]int),
	}
}

// Record records a verification result
func (rc *ReliabilityChecker) Record(verifierPath string, success bool) {
	rc.Total[verifierPath]++
	if !success {
		rc.Failures[verifierPath]++
	}
}

// ShouldFlag returns true if verifier should be flagged (3+ consecutive failures)
func (rc *ReliabilityChecker) ShouldFlag(verifierPath string) bool {
	if rc.Failures[verifierPath] >= 3 {
		return true
	}
	return false
}

// GetReliability returns reliability percentage
func (rc *ReliabilityChecker) GetReliability(verifierPath string) float64 {
	total := rc.Total[verifierPath]
	if total == 0 {
		return 100.0
	}
	failures := rc.Failures[verifierPath]
	return float64(total-failures) / float64(total) * 100
}