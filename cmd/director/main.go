package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/Sreyash-Reddy/SATOS/internal/state"
	"github.com/spf13/cobra"
)

var (
	baseDir     string
	dbPath      string
	logPath     string
	queueDir    string
	doneDir     string
	failedDir   string
	verifiersDir string
)

var rootCmd = &cobra.Command{
	Use:   "satos",
	Short: "SATOS - Swarm Agentic Team Operating System",
	Long:  `Director V0 - Manages AI agent teams with structured task execution`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(`
╔══════════════════════════════════════════════════╗
║   SATOS Director V0                              ║
║   Swarm Agentic Team Operating System            ║
╚══════════════════════════════════════════════════╝
`)
		cmd.Help()
	},
}

// initCmd initializes SATOS state
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize SATOS state directory and database",
	Run: func(cmd *cobra.Command, args []string) {
		if err := initializeSATOS(); err != nil {
			log.Fatalf("Initialization failed: %v", err)
		}
		fmt.Println("\n✓ SATOS initialized successfully")
	},
}

// statusCmd shows SATOS status
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show SATOS status and task statistics",
	Run: func(cmd *cobra.Command, args []string) {
		showStatus()
	},
}

// tasksCmd shows queued tasks
var tasksCmd = &cobra.Command{
	Use:   "tasks",
	Short: "List all tasks",
	Run: func(cmd *cobra.Command, args []string) {
		listTasks()
	},
}

// Create the directory structure
func createDirectories() error {
	dirs := []string{queueDir, doneDir, failedDir, verifiersDir}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create %s: %w", dir, err)
		}
		fmt.Printf("  ✓ %s\n", dir)
	}
	return nil
}

func initializeSATOS() error {
	fmt.Println("\nInitializing SATOS...")
	fmt.Println()

	// Check if this is a git repository
	if _, err := os.Stat(".git"); os.IsNotExist(err) {
		fmt.Println("  ✗ Not a git repository")
		return fmt.Errorf("not a git repository. SATOS requires a git-tracked project")
	}
	fmt.Println("  ✓ Git repository verified")

	// Create directory structure
	fmt.Println()
	fmt.Println("Creating directory structure...")
	if err := createDirectories(); err != nil {
		return err
	}

	// Initialize database
	fmt.Println()
	fmt.Println("Initializing database...")
	cfg := state.Config{
		DBPath: dbPath,
	}
	mgr, err := state.NewManager(cfg)
	if err != nil {
		return fmt.Errorf("failed to create state manager: %w", err)
	}
	defer mgr.Close()

	// Set metadata
	mgr.SetMeta("version", "0.1.0")
	mgr.SetMeta("initialized", "true")
	fmt.Printf("  ✓ Database: %s\n", dbPath)

	// Create director log
	if err := os.MkdirAll(filepath.Dir(logPath), 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}
	f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to create log file: %w", err)
	}
	f.Close()
	fmt.Printf("  ✓ Log file: %s\n", logPath)

	// Create task template
	fmt.Println()
	fmt.Println("Creating task template...")
	templatePath := filepath.Join(baseDir, "TASK-TEMPLATE.md")
	if err := os.WriteFile(templatePath, []byte(taskTemplate), 0644); err != nil {
		return fmt.Errorf("failed to create task template: %w", err)
	}
	fmt.Printf("  ✓ Template: %s\n", templatePath)

	return nil
}

func showStatus() {
	cfg := state.Config{DBPath: dbPath}
	mgr, err := state.NewManager(cfg)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer mgr.Close()

	fmt.Println()
	fmt.Println("═══ SATOS Status ═══")
	fmt.Println()

	// Get version
	version, _ := mgr.GetMeta("version")
	fmt.Printf("  Version:    %s\n", version)

	// Count tasks by status
	queued, _ := mgr.GetTasksByStatus("queued")
	inProgress, _ := mgr.GetTasksByStatus("in_progress")
	done, _ := mgr.GetTasksByStatus("done")
	failed, _ := mgr.GetTasksByStatus("failed")

	fmt.Printf("  Queued:     %d\n", len(queued))
	fmt.Printf("  In Progress: %d\n", len(inProgress))
	fmt.Printf("  Done:       %d\n", len(done))
	fmt.Printf("  Failed:     %d\n", len(failed))
	fmt.Println()
}

func listTasks() {
	cfg := state.Config{DBPath: dbPath}
	mgr, err := state.NewManager(cfg)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer mgr.Close()

	fmt.Println()
	fmt.Println("═══ Tasks ═══")
	fmt.Println()

	statuses := []string{"queued", "in_progress", "done", "failed"}
	for _, status := range statuses {
		tasks, _ := mgr.GetTasksByStatus(status)
		if len(tasks) > 0 {
			fmt.Printf("  [%s]\n", status)
			for _, t := range tasks {
				fmt.Printf("    • %s - %s\n", t.TaskID, t.Title)
			}
			fmt.Println()
		}
	}
}

const taskTemplate = `# Task Template

## Status
queued | in_progress | done | failed

## Command
{concrete action: edit X, implement Y, test Z}

## Output
{file path or stdout expectation}

## Verify
{shell command that exits 0 = success}

## Depends-On
{TASK-id, optional — omit if N/A}

## Example
# TASK-001
## Status: queued

## Command
Create a hello world Python file at src/hello.py

## Output
src/hello.py

## Verify
python3 src/hello.py && echo "Hello, SATOS!"

## Depends-On
null
`

func init() {
	// Set default paths
	home, _ := os.UserHomeDir()

	rootCmd.PersistentFlags().StringVar(&baseDir, "dir", ".satos", "SATOS state directory")
	rootCmd.PersistentFlags().StringVar(&dbPath, "db", ".satos/state.db", "Database path")
	rootCmd.PersistentFlags().StringVar(&logPath, "log", ".satos/director.log", "Log file path")
	rootCmd.PersistentFlags().StringVar(&queueDir, "queue", ".satos/tasks/queue", "Queue directory")
	rootCmd.PersistentFlags().StringVar(&doneDir, "done", ".satos/tasks/done", "Done directory")
	rootCmd.PersistentFlags().StringVar(&failedDir, "failed", ".satos/tasks/failed", "Failed directory")
	rootCmd.PersistentFlags().StringVar(&verifiersDir, "verifiers", ".satos/verifiers", "Verifiers directory")

	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(tasksCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}