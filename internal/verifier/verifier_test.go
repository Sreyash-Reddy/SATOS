package verifier

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewVerifier(t *testing.T) {
	v := NewVerifier("echo test")
	assert.Equal(t, "echo test", v.Command)
	assert.Equal(t, 0, v.FailedRuns)
	assert.False(t, v.Flagged)
}

func TestVerifierRunSuccess(t *testing.T) {
	v := NewVerifier("echo 'hello world' | grep hello")
	result, err := v.Run()

	require.NoError(t, err)
	assert.True(t, result.Success)
	assert.Equal(t, 0, result.ExitCode)
	assert.Contains(t, result.Output, "hello world")
}

func TestVerifierRunFailure(t *testing.T) {
	v := NewVerifier("exit 1")
	result, err := v.Run()

	require.NoError(t, err)
	assert.False(t, result.Success)
	assert.Equal(t, 1, result.ExitCode)
}

func TestParseTaskFile(t *testing.T) {
	content := `# TASK-001
## Status: queued

## Command
echo "Hello, SATOS!"

## Output
stdout

## Verify
echo "Verified!"

## Depends-On
null
`

	tf, err := ParseTaskFile(content)
	require.NoError(t, err)

	assert.Equal(t, "queued", tf.Status)
	assert.Equal(t, "echo \"Hello, SATOS!\"", tf.Command)
	assert.Equal(t, "echo \"Verified!\"", tf.Verify)
	assert.Equal(t, "stdout", tf.Output)
}

func TestTaskFileValidation(t *testing.T) {
	valid := &TaskFile{
		Command: "echo test",
		Output:  "file.txt",
	}
	assert.NoError(t, valid.Validate())

	invalid := &TaskFile{
		Command: "",
		Output:  "file.txt",
	}
	assert.Error(t, invalid.Validate())
}

func TestReliabilityChecker(t *testing.T) {
	rc := NewReliabilityChecker()

	// Record successes
	for i := 0; i < 5; i++ {
		rc.Record("test-verifier", true)
	}

	// Should not be flagged
	assert.False(t, rc.ShouldFlag("test-verifier"))
	assert.Equal(t, 100.0, rc.GetReliability("test-verifier"))

	// Add failures
	rc.Record("test-verifier", false)
	rc.Record("test-verifier", false)
	rc.Record("test-verifier", false)

	// Should be flagged
	assert.True(t, rc.ShouldFlag("test-verifier"))
	assert.Equal(t, 62.5, rc.GetReliability("test-verifier"))
}

func TestVerifyRealScript(t *testing.T) {
	// Create a test script
	script := t.TempDir() + "/test-script.sh"
	err := os.WriteFile(script, []byte("#!/bin/bash\necho 'Hello'\n"), 0755)
	require.NoError(t, err)

	v := NewVerifier(script)
	result, err := v.Run()

	require.NoError(t, err)
	assert.True(t, result.Success)
}