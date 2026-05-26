package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/Sreyash-Reddy/SATOS/internal/state"
	"github.com/Sreyash-Reddy/SATOS/internal/verifier"
	"github.com/spf13/cobra"
)

var (
	dbPathTL     string
	queueDirTL   string
	doneDirTL    string
	failedDirTL  string
	verifiersDirTL string
	timeout      int
)

var rootCmdTL = &cobra.Command{
	Use:   "tl",
	Short: "SATOS Team Leader - Dispatches and manages tasks",
	Long:  `TL manages task execution, spawns agents, and verifies results`,
}

// dispatchCmd dispatches next queued task
var dispatchCmd = &cobra.Command{
	Use:   "dispatch",
	Short: "Dispatch the next queued task to an agent",
	Run: func(cmd *cobra.Command, args []string) {
		dispatchTask()
	},
}

// runCmd runs the TL in continuous mode
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run TL continuously, processing tasks as they appear",
	Run: func(cmd *cobra.Command, args []string) {
		runTL()
	},
}

// TaskFile represents a task file
type TaskFile struct {
	Status    string
	Command   string
	Output    string
	Verify    string
	DependsOn string
}

func parseTaskFile(path string) (*TaskFile, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read task file: %w", err)
	}

	tf := &TaskFile{}
	lines := strings.Split(string(content), "\n")

	for i, line := range lines {
		line = strings.TrimSpace(line)

		switch {
		case strings.HasPrefix(line, "## Status:"):
			tf.Status = strings.TrimSpace(strings.TrimPrefix(line, "## Status:"))
		case strings.HasPrefix(line, "## Command") && i+1 < len(lines):
			// Collect until next ##
			var cmdLines []string
			for j := i + 1; j < len(lines) && !strings.HasPrefix(strings.TrimSpace(lines[j]), "##"); j++ {
				cmdLines = append(cmdLines, lines[j])
			}
			tf.Command = strings.Trim(strings.Join(cmdLines, "\n"), "\n ")
		case strings.HasPrefix(line, "## Verify") && i+1 < len(lines):
			var verifyLines []string
			for j := i + 1; j < len(lines) && !strings.HasPrefix(strings.TrimSpace(lines[j]), "##"); j++ {
				verifyLines = append(verifyLines, lines[j])
			}
			tf.Verify = strings.Trim(strings.Join(verifyLines, "\n"), "\n ")
		case strings.HasPrefix(line, "## Output:"):
			tf.Output = strings.TrimSpace(strings.TrimPrefix(line, "## Output:"))
		case strings.HasPrefix(line, "## Depends-On:"):
			tf.DependsOn = strings.TrimSpace(strings.TrimPrefix(line, "## Depends-On:"))
		}
	}

	return tf, nil
}

func dispatchTask() {
	cfg := state.Config{DBPath: dbPathTL}
	mgr, err := state.NewManager(cfg)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer mgr.Close()

	// Get next queued task from database
	tasks, err := mgr.GetTasksByStatus("queued")
	if err != nil {
		log.Fatalf("Failed to get tasks: %v", err)
	}

	if len(tasks) == 0 {
		// Check queue directory for task files
		files, err := os.ReadDir(queueDirTL)
		if err != nil {
			log.Fatalf("Failed to read queue directory: %v", err)
		}

		for _, f := range files {
			if filepath.Ext(f.Name()) == ".md" {
				taskID := strings.TrimSuffix(f.Name(), ".md")
				fmt.Printf("Processing task file: %s\n", taskID)
				processTaskFile(filepath.Join(queueDirTL, f.Name()), taskID, mgr)
				return
			}
		}

		fmt.Println("No tasks in queue")
		return
	}

	// Process from database
	task := tasks[0]
	processTask(task.TaskID, mgr)
}

func processTaskFile(path, taskID string, mgr *state.Manager) {
	tf, err := parseTaskFile(path)
	if err != nil {
		log.Printf("Failed to parse task file %s: %v", taskID, err)
		return
	}

	// Create task in database
	task := &state.Task{
		TaskID:        taskID,
		Title:         taskID,
		Status:        "in_progress",
		VerifyCommand: tf.Verify,
		OutputPath:    tf.Output,
	}
	if err := mgr.CreateTask(task); err != nil {
		log.Printf("Failed to create task: %v", err)
	}

	fmt.Printf("  Executing: %s\n", tf.Command)
	fmt.Printf("  Output: %s\n", tf.Output)

	// Execute the command
	cmd := exec.Command("bash", "-c", tf.Command)
	cmd.Dir, _ = os.Getwd()

	// Set timeout
	done := make(chan error, 1)
	go func() {
		done <- cmd.Run()
	}()

	select {
	case err := <-done:
		if err != nil {
			fmt.Printf("  ✗ Execution failed: %v\n", err)
			mgr.UpdateTaskStatus(taskID, "failed")
			moveToFailed(path)
			return
		}
	case <-time.After(time.Duration(timeout) * time.Second):
		cmd.Process.Kill()
		fmt.Printf("  ✗ Execution timed out after %d seconds\n", timeout)
		mgr.UpdateTaskStatus(taskID, "failed")
		moveToFailed(path)
		return
	}

	fmt.Println("  ✓ Command executed")

	// Verify if verification command is provided
	if tf.Verify != "" {
		fmt.Printf("  Verifying: %s\n", tf.Verify)

		v := verifier.NewVerifier(tf.Verify)
		result, err := v.Run()
		if err != nil {
			fmt.Printf("  ✗ Verification error: %v\n", err)
			mgr.RecordRecovery(taskID, "verification_error", err.Error())
		}

		// Log verification
		mgr.LogVerification(taskID, tf.Verify, result.ExitCode, result.Output)

		if !result.Success {
			fmt.Printf("  ✗ Verification failed (exit code: %d)\n", result.ExitCode)
			mgr.UpdateTaskStatus(taskID, "failed")
			moveToFailed(path)
			return
		}

		fmt.Println("  ✓ Verification passed")
	}

	// Mark as done
	mgr.UpdateTaskStatus(taskID, "done")
	moveToDone(path)

	fmt.Printf("  ✓ Task %s completed successfully\n", taskID)
}

func processTask(taskID string, mgr *state.Manager) {
	task, err := mgr.GetTask(taskID)
	if err != nil {
		log.Printf("Failed to get task %s: %v", taskID, err)
		return
	}

	fmt.Printf("Processing task: %s\n", taskID)

	// Execute command from database
	fmt.Printf("  Executing: %s\n", task.Description)

	cmd := exec.Command("bash", "-c", task.Description)
	cmd.Dir, _ = os.Getwd()

	done := make(chan error, 1)
	go func() {
		done <- cmd.Run()
	}()

	select {
	case err := <-done:
		if err != nil {
			fmt.Printf("  ✗ Execution failed: %v\n", err)
			mgr.UpdateTaskStatus(taskID, "failed")
			return
		}
	case <-time.After(time.Duration(timeout) * time.Second):
		cmd.Process.Kill()
		fmt.Printf("  ✗ Execution timed out after %d seconds\n", timeout)
		mgr.UpdateTaskStatus(taskID, "failed")
		return
	}

	// Run verification if command exists
	if task.VerifyCommand != "" {
		fmt.Printf("  Verifying: %s\n", task.VerifyCommand)

		v := verifier.NewVerifier(task.VerifyCommand)
		result, _ := v.Run()
		mgr.LogVerification(taskID, task.VerifyCommand, result.ExitCode, result.Output)

		if !result.Success {
			fmt.Printf("  ✗ Verification failed\n")
			mgr.UpdateTaskStatus(taskID, "failed")
			return
		}
	}

	mgr.UpdateTaskStatus(taskID, "done")
	fmt.Printf("  ✓ Task %s completed\n", taskID)
}

func moveToDone(path string) {
	dest := filepath.Join(doneDirTL, filepath.Base(path))
	os.MkdirAll(doneDirTL, 0755)
	os.Rename(path, dest)
}

func moveToFailed(path string) {
	dest := filepath.Join(failedDirTL, filepath.Base(path))
	os.MkdirAll(failedDirTL, 0755)
	os.Rename(path, dest)
}

func runTL() {
	cfg := state.Config{DBPath: dbPathTL}
	mgr, err := state.NewManager(cfg)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer mgr.Close()

	fmt.Println("SATOS TL running continuously...")
	fmt.Println("Press Ctrl+C to stop")
	fmt.Println()

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			dispatchTask()
		}
	}
}

func init() {
	rootCmdTL.PersistentFlags().StringVar(&dbPathTL, "db", ".satos/state.db", "Database path")
	rootCmdTL.PersistentFlags().StringVar(&queueDirTL, "queue", ".satos/tasks/queue", "Queue directory")
	rootCmdTL.PersistentFlags().StringVar(&doneDirTL, "done", ".satos/tasks/done", "Done directory")
	rootCmdTL.PersistentFlags().StringVar(&failedDirTL, "failed", ".satos/tasks/failed", "Failed directory")
	rootCmdTL.PersistentFlags().IntVar(&timeout, "timeout", 300, "Task timeout in seconds")

	rootCmdTL.AddCommand(dispatchCmd)
	rootCmdTL.AddCommand(runCmd)
}

func main() {
	if err := rootCmdTL.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}