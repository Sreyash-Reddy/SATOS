package state

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/Sreyash-Reddy/SATOS/internal/verifier"
	_ "github.com/mattn/go-sqlite3"
)

// Manager handles all SATOS state via SQLite
type Manager struct {
	db *sql.DB
}

// Config holds SATOS configuration
type Config struct {
	BaseDir     string
	DBPath      string
	LogPath     string
	QueueDir    string
	DoneDir     string
	FailedDir   string
	VerifiersDir string
}

// NewManager creates a new state manager
func NewManager(cfg Config) (*Manager, error) {
	db, err := sql.Open("sqlite3", cfg.DBPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Enable foreign keys
	_, err = db.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	m := &Manager{db: db}
	if err := m.initSchema(); err != nil {
		return nil, fmt.Errorf("failed to init schema: %w", err)
	}

	return m, nil
}

// Close closes the database connection
func (m *Manager) Close() error {
	return m.db.Close()
}

// initSchema creates all necessary tables
func (m *Manager) initSchema() error {
	schema := `
	CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		task_id TEXT UNIQUE NOT NULL,
		title TEXT NOT NULL,
		description TEXT,
		status TEXT DEFAULT 'queued' CHECK(status IN ('queued', 'in_progress', 'done', 'failed')),
		team TEXT,
		assignee TEXT,
		story_points INTEGER,
		sprint TEXT,
		priority TEXT,
		depends_on TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		completed_at DATETIME,
		verify_command TEXT,
		output_path TEXT
	);

	CREATE TABLE IF NOT EXISTS verification_logs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		task_id TEXT NOT NULL,
		command TEXT NOT NULL,
		exit_code INTEGER,
		output TEXT,
		verified_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS verifier_reliability (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		verifier_path TEXT NOT NULL UNIQUE,
		total_runs INTEGER DEFAULT 0,
		failures INTEGER DEFAULT 0,
		last_run DATETIME,
		flagged BOOLEAN DEFAULT FALSE
	);

	CREATE TABLE IF NOT EXISTS recovery_log (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		task_id TEXT NOT NULL,
		event_type TEXT NOT NULL,
		details TEXT,
		recovered_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS director_meta (
		key TEXT PRIMARY KEY,
		value TEXT,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_tasks_status ON tasks(status);
	CREATE INDEX IF NOT EXISTS idx_tasks_team ON tasks(team);
	CREATE INDEX IF NOT EXISTS idx_tasks_sprint ON tasks(sprint);
	CREATE INDEX IF NOT EXISTS idx_verification_task ON verification_logs(task_id);
	`

	_, err := m.db.Exec(schema)
	return err
}

// Task represents a SATOS task
type Task struct {
	ID            int64
	TaskID        string
	Title         string
	Description   string
	Status        string
	Team          string
	Assignee      string
	StoryPoints   int
	Sprint        string
	Priority      string
	DependsOn     string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	CompletedAt   *time.Time
	VerifyCommand string
	OutputPath    string
}

// CreateTask creates a new task
func (m *Manager) CreateTask(t *Task) error {
	query := `
	INSERT INTO tasks (task_id, title, description, status, team, assignee, story_points, sprint, priority, depends_on, verify_command, output_path)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	result, err := m.db.Exec(query, t.TaskID, t.Title, t.Description, t.Status, t.Team, t.Assignee, t.StoryPoints, t.Sprint, t.Priority, t.DependsOn, t.VerifyCommand, t.OutputPath)
	if err != nil {
		return fmt.Errorf("failed to create task: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert ID: %w", err)
	}
	t.ID = id
	return nil
}

// GetTask retrieves a task by ID
func (m *Manager) GetTask(taskID string) (*Task, error) {
	query := `SELECT id, task_id, title, description, status, team, assignee, story_points, sprint, priority, depends_on, created_at, updated_at, completed_at, verify_command, output_path FROM tasks WHERE task_id = ?`
	row := m.db.QueryRow(query, taskID)

	var t Task
	var completedAt sql.NullTime
	err := row.Scan(&t.ID, &t.TaskID, &t.Title, &t.Description, &t.Status, &t.Team, &t.Assignee, &t.StoryPoints, &t.Sprint, &t.Priority, &t.DependsOn, &t.CreatedAt, &t.UpdatedAt, &completedAt, &t.VerifyCommand, &t.OutputPath)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("task not found: %s", taskID)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	if completedAt.Valid {
		t.CompletedAt = &completedAt.Time
	}

	return &t, nil
}

// UpdateTaskStatus updates a task's status
func (m *Manager) UpdateTaskStatus(taskID, status string) error {
	query := `UPDATE tasks SET status = ?, updated_at = CURRENT_TIMESTAMP WHERE task_id = ?`
	if status == "done" || status == "failed" {
		query = `UPDATE tasks SET status = ?, updated_at = CURRENT_TIMESTAMP, completed_at = CURRENT_TIMESTAMP WHERE task_id = ?`
	}

	_, err := m.db.Exec(query, status, taskID)
	return err
}

// GetTasksByStatus returns all tasks with a given status
func (m *Manager) GetTasksByStatus(status string) ([]*Task, error) {
	query := `SELECT id, task_id, title, description, status, team, assignee, story_points, sprint, priority, depends_on, created_at, updated_at, completed_at, verify_command, output_path FROM tasks WHERE status = ?`

	rows, err := m.db.Query(query, status)
	if err != nil {
		return nil, fmt.Errorf("failed to query tasks: %w", err)
	}
	defer rows.Close()

	var tasks []*Task
	for rows.Next() {
		var t Task
		var completedAt sql.NullTime
		err := rows.Scan(&t.ID, &t.TaskID, &t.Title, &t.Description, &t.Status, &t.Team, &t.Assignee, &t.StoryPoints, &t.Sprint, &t.Priority, &t.DependsOn, &t.CreatedAt, &t.UpdatedAt, &completedAt, &t.VerifyCommand, &t.OutputPath)
		if err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}
		if completedAt.Valid {
			t.CompletedAt = &completedAt.Time
		}
		tasks = append(tasks, &t)
	}

	return tasks, nil
}

// LogVerification logs a verification event
func (m *Manager) LogVerification(taskID, command string, exitCode int, output string) error {
	query := `INSERT INTO verification_logs (task_id, command, exit_code, output) VALUES (?, ?, ?, ?)`
	_, err := m.db.Exec(query, taskID, command, exitCode, output)
	return err
}

// RecordRecovery records a recovery event
func (m *Manager) RecordRecovery(taskID, eventType, details string) error {
	query := `INSERT INTO recovery_log (task_id, event_type, details) VALUES (?, ?, ?)`
	_, err := m.db.Exec(query, taskID, eventType, details)
	return err
}

// UpdateVerifierReliability updates verifier reliability stats
func (m *Manager) UpdateVerifierReliability(v *verifier.Verifier) error {
	query := `
	INSERT INTO verifier_reliability (verifier_path, total_runs, failures, last_run, flagged)
	VALUES (?, ?, ?, CURRENT_TIMESTAMP, ?)
	ON CONFLICT(verifier_path) DO UPDATE SET
		total_runs = total_runs + 1,
		failures = failures + ?,
		last_run = CURRENT_TIMESTAMP,
		flagged = CASE WHEN failures + ? >= 3 THEN TRUE ELSE flagged END
	`
	_, err := m.db.Exec(query, v.Path, 1, v.FailedRuns, v.Flagged, v.FailedRuns, v.FailedRuns)
	return err
}

// SetMeta sets a director metadata value
func (m *Manager) SetMeta(key, value string) error {
	query := `INSERT OR REPLACE INTO director_meta (key, value, updated_at) VALUES (?, ?, CURRENT_TIMESTAMP)`
	_, err := m.db.Exec(query, key, value)
	return err
}

// GetMeta gets a director metadata value
func (m *Manager) GetMeta(key string) (string, error) {
	query := `SELECT value FROM director_meta WHERE key = ?`
	row := m.db.QueryRow(query, key)
	var value string
	err := row.Scan(&value)
	if err == sql.ErrNoRows {
		return "", nil
	}
	return value, err
}