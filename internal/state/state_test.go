package state

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewManager(t *testing.T) {
	// Create temp directory
	tmpDir := t.TempDir()
	dbPath := tmpDir + "/test.db"

	cfg := Config{DBPath: dbPath}
	mgr, err := NewManager(cfg)
	require.NoError(t, err)
	defer mgr.Close()

	// Verify database was created
	_, err = os.Stat(dbPath)
	assert.NoError(t, err)
}

func TestCreateAndGetTask(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := tmpDir + "/test.db"

	cfg := Config{DBPath: dbPath}
	mgr, err := NewManager(cfg)
	require.NoError(t, err)
	defer mgr.Close()

	// Create a task
	task := &Task{
		TaskID:  "TEST-001",
		Title:   "Test Task",
		Status:  "queued",
		Sprint:  "Sprint 1",
		Priority: "High",
	}

	err = mgr.CreateTask(task)
	require.NoError(t, err)
	assert.Greater(t, task.ID, int64(0))

	// Get the task back
	retrieved, err := mgr.GetTask("TEST-001")
	require.NoError(t, err)
	assert.Equal(t, "Test Task", retrieved.Title)
	assert.Equal(t, "queued", retrieved.Status)
}

func TestUpdateTaskStatus(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := tmpDir + "/test.db"

	cfg := Config{DBPath: dbPath}
	mgr, err := NewManager(cfg)
	require.NoError(t, err)
	defer mgr.Close()

	// Create a task
	task := &Task{TaskID: "TEST-002", Title: "Test", Status: "queued"}
	err = mgr.CreateTask(task)
	require.NoError(t, err)

	// Update status
	err = mgr.UpdateTaskStatus("TEST-002", "in_progress")
	require.NoError(t, err)

	// Verify
	retrieved, err := mgr.GetTask("TEST-002")
	require.NoError(t, err)
	assert.Equal(t, "in_progress", retrieved.Status)
}

func TestGetTasksByStatus(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := tmpDir + "/test.db"

	cfg := Config{DBPath: dbPath}
	mgr, err := NewManager(cfg)
	require.NoError(t, err)
	defer mgr.Close()

	// Create multiple tasks
	for i := 1; i <= 3; i++ {
		task := &Task{
			TaskID: "TEST-" + string(rune('0'+i)),
			Title:  "Task " + string(rune('0'+i)),
			Status: "queued",
		}
		err = mgr.CreateTask(task)
		require.NoError(t, err)
	}

	// Get queued tasks
	tasks, err := mgr.GetTasksByStatus("queued")
	require.NoError(t, err)
	assert.Len(t, tasks, 3)

	// Get done tasks
	doneTasks, err := mgr.GetTasksByStatus("done")
	require.NoError(t, err)
	assert.Len(t, doneTasks, 0)
}

func TestSetAndGetMeta(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := tmpDir + "/test.db"

	cfg := Config{DBPath: dbPath}
	mgr, err := NewManager(cfg)
	require.NoError(t, err)
	defer mgr.Close()

	// Set metadata
	err = mgr.SetMeta("version", "0.1.0")
	require.NoError(t, err)

	// Get metadata
	value, err := mgr.GetMeta("version")
	require.NoError(t, err)
	assert.Equal(t, "0.1.0", value)
}

func TestLogVerification(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := tmpDir + "/test.db"

	cfg := Config{DBPath: dbPath}
	mgr, err := NewManager(cfg)
	require.NoError(t, err)
	defer mgr.Close()

	// Create a task
	task := &Task{TaskID: "TEST-VERIFY", Title: "Verify Test", Status: "queued"}
	err = mgr.CreateTask(task)
	require.NoError(t, err)

	// Log verification
	err = mgr.LogVerification("TEST-VERIFY", "echo test", 0, "test output")
	require.NoError(t, err)
}

func TestRecordRecovery(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := tmpDir + "/test.db"

	cfg := Config{DBPath: dbPath}
	mgr, err := NewManager(cfg)
	require.NoError(t, err)
	defer mgr.Close()

	// Create a task
	task := &Task{TaskID: "TEST-RECOV", Title: "Recovery Test", Status: "queued"}
	err = mgr.CreateTask(task)
	require.NoError(t, err)

	// Record recovery
	err = mgr.RecordRecovery("TEST-RECOV", "crash", "Agent crashed mid-execution")
	require.NoError(t, err)
}