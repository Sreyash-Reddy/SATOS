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
	var inCodeBlock bool

	for _, line := range lines {
		// Skip markdown headers for now
		if strings.HasPrefix(line, "##") {
			parts := strings.SplitN(strings.TrimSpace(line), " ", 2)
			if len(parts) >= 2 {
				currentSection = strings.ToLower(strings.TrimSpace(parts[1]))
				continue
			}
		}

		// Track code blocks
		if strings.HasPrefix(line, "```") {
			inCodeBlock = !inCodeBlock
			continue
		}

		if inCodeBlock {
			switch currentSection {
			case "command":
				tf.Command += line + "\n"
			case "verify":
				tf.Verify += line + "\n"
			}
			continue
		}

		switch currentSection {
		case "status":
			tf.Status = strings.TrimSpace(line)
		case "output":
			tf.Output = strings.TrimSpace(line)
		case "depends-on":
			tf.DependsOn = strings.TrimSpace(line)
		case "sprint":
			tf.Sprint = strings.TrimSpace(line)
		case "priority":
			tf.Priority = strings.TrimSpace(line)
		case "assignee":
			tf.Assignee = strings.TrimSpace(line)
		}
	}

	// Clean up command and verify (remove trailing newlines and leading/trailing whitespace)
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