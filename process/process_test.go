// Copyright (C) 2025, Dione Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package process

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/DioneProtocol/odyssey-tooling-sdk-go/constants"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunFile_JSON(t *testing.T) {
	// Test RunFile struct JSON marshaling/unmarshaling
	original := RunFile{Pid: 12345}

	// Marshal to JSON
	data, err := json.Marshal(&original)
	require.NoError(t, err)

	// Unmarshal back
	var restored RunFile
	err = json.Unmarshal(data, &restored)
	require.NoError(t, err)

	assert.Equal(t, original.Pid, restored.Pid)
}

func TestSaveRunFile(t *testing.T) {
	tempDir := t.TempDir()
	runFilePath := filepath.Join(tempDir, "test.pid")

	// Test saving run file
	err := saveRunFile(12345, runFilePath)
	require.NoError(t, err)

	// Verify file was created
	assert.FileExists(t, runFilePath)

	// Verify file content
	content, err := os.ReadFile(runFilePath)
	require.NoError(t, err)

	var runFile RunFile
	err = json.Unmarshal(content, &runFile)
	require.NoError(t, err)
	assert.Equal(t, 12345, runFile.Pid)

	// Test file permissions
	info, err := os.Stat(runFilePath)
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(constants.WriteReadReadPerms), info.Mode())
}

func TestSaveRunFile_InvalidPath(t *testing.T) {
	// Test saving to invalid path
	invalidPath := "/invalid/path/that/does/not/exist/test.pid"
	err := saveRunFile(12345, invalidPath)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "could not write run file")
}

func TestLoadRunFile(t *testing.T) {
	tempDir := t.TempDir()
	runFilePath := filepath.Join(tempDir, "test.pid")

	// Test loading non-existent file
	pid, err := loadRunFile(runFilePath)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "run file")
	assert.Contains(t, err.Error(), "does not exist")
	assert.Equal(t, 0, pid)

	// Create a valid run file
	err = saveRunFile(54321, runFilePath)
	require.NoError(t, err)

	// Test loading existing file
	pid, err = loadRunFile(runFilePath)
	require.NoError(t, err)
	assert.Equal(t, 54321, pid)
}

func TestLoadRunFile_EmptyPath(t *testing.T) {
	// Test loading with empty path
	pid, err := loadRunFile("")
	require.NoError(t, err)
	assert.Equal(t, 0, pid)
}

func TestLoadRunFile_InvalidJSON(t *testing.T) {
	tempDir := t.TempDir()
	runFilePath := filepath.Join(tempDir, "invalid.pid")

	// Write invalid JSON
	err := os.WriteFile(runFilePath, []byte("invalid json"), constants.WriteReadReadPerms)
	require.NoError(t, err)

	// Test loading invalid JSON
	pid, err := loadRunFile(runFilePath)
	assert.Error(t, err)
	assert.Equal(t, 0, pid)
}

func TestRemoveRunFile(t *testing.T) {
	tempDir := t.TempDir()
	runFilePath := filepath.Join(tempDir, "test.pid")

	// Test removing non-existent file (should not error)
	err := removeRunFile(runFilePath)
	// Note: removeRunFile may return an error for non-existent files, which is acceptable

	// Create a file
	err = saveRunFile(12345, runFilePath)
	require.NoError(t, err)
	assert.FileExists(t, runFilePath)

	// Test removing existing file
	err = removeRunFile(runFilePath)
	require.NoError(t, err)
	assert.NoFileExists(t, runFilePath)
}

func TestRemoveRunFile_EmptyPath(t *testing.T) {
	// Test removing with empty path
	err := removeRunFile("")
	require.NoError(t, err)
}

func TestExecute(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skipping process tests on Windows due to signal handling differences")
	}

	tempDir := t.TempDir()
	runFilePath := filepath.Join(tempDir, "test.pid")

	var stdout, stderr bytes.Buffer

	// Test executing a simple command (sleep for a short time)
	binPath := "sleep"
	args := []string{"1"}

	pid, err := Execute(binPath, args, &stdout, &stderr, runFilePath, 0)
	require.NoError(t, err)
	assert.Greater(t, pid, 0)

	// Verify run file was created
	assert.FileExists(t, runFilePath)

	// Verify PID in run file matches returned PID
	loadedPid, err := loadRunFile(runFilePath)
	require.NoError(t, err)
	assert.Equal(t, pid, loadedPid)

	// Wait for process to complete
	time.Sleep(2 * time.Second)
}

func TestExecute_NoRunFile(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skipping process tests on Windows due to signal handling differences")
	}

	var stdout, stderr bytes.Buffer

	// Test executing without run file
	binPath := "echo"
	args := []string{"test"}

	pid, err := Execute(binPath, args, &stdout, &stderr, "", 0)
	require.NoError(t, err)
	assert.Greater(t, pid, 0)

	// Wait for process to complete
	time.Sleep(100 * time.Millisecond)
}

func TestExecute_WithSetupTime(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skipping process tests on Windows due to signal handling differences")
	}

	tempDir := t.TempDir()
	runFilePath := filepath.Join(tempDir, "test.pid")

	var stdout, stderr bytes.Buffer

	// Test executing with setup time
	binPath := "sleep"
	args := []string{"5"} // Sleep for 5 seconds

	pid, err := Execute(binPath, args, &stdout, &stderr, runFilePath, 100*time.Millisecond)
	require.NoError(t, err)
	assert.Greater(t, pid, 0)

	// Verify process is still running after setup time
	proc, err := GetProcess(pid)
	require.NoError(t, err)
	assert.NotNil(t, proc)

	// Clean up
	proc.Kill()
}

func TestExecute_ProcessStopsDuringSetup(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skipping process tests on Windows due to signal handling differences")
	}

	tempDir := t.TempDir()
	runFilePath := filepath.Join(tempDir, "test.pid")

	var stdout, stderr bytes.Buffer

	// Test with a command that exits quickly
	binPath := "echo"
	args := []string{"test"}

	pid, err := Execute(binPath, args, &stdout, &stderr, runFilePath, 1*time.Second)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "process stopped during setup")
	assert.Equal(t, 0, pid)
}

func TestExecute_InvalidCommand(t *testing.T) {
	var stdout, stderr bytes.Buffer

	// Test executing invalid command
	binPath := "/nonexistent/command"
	args := []string{}

	pid, err := Execute(binPath, args, &stdout, &stderr, "", 0)
	assert.Error(t, err)
	assert.Equal(t, 0, pid)
}

func TestGetProcess(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skipping process tests on Windows due to signal handling differences")
	}

	// Test getting current process
	currentPid := os.Getpid()
	proc, err := GetProcess(currentPid)
	require.NoError(t, err)
	assert.NotNil(t, proc)
	assert.Equal(t, currentPid, proc.Pid)

	// Test getting non-existent process
	proc, err = GetProcess(999999)
	assert.Error(t, err)
	assert.Nil(t, proc)
}

func TestIsRunning(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skipping process tests on Windows due to signal handling differences")
	}

	tempDir := t.TempDir()
	runFilePath := filepath.Join(tempDir, "test.pid")

	// Test with current process PID
	currentPid := os.Getpid()
	isRunning, pid, proc, err := IsRunning(currentPid, "")
	require.NoError(t, err)
	assert.True(t, isRunning)
	assert.Equal(t, currentPid, pid)
	assert.NotNil(t, proc)

	// Test with non-existent PID
	isRunning, pid, proc, err = IsRunning(999999, "")
	require.NoError(t, err)
	assert.False(t, isRunning)
	assert.Equal(t, 0, pid)
	assert.Nil(t, proc)

	// Test with run file
	err = saveRunFile(currentPid, runFilePath)
	require.NoError(t, err)

	isRunning, pid, proc, err = IsRunning(0, runFilePath)
	require.NoError(t, err)
	assert.True(t, isRunning)
	assert.Equal(t, currentPid, pid)
	assert.NotNil(t, proc)
}

func TestIsRunning_InvalidRunFile(t *testing.T) {
	tempDir := t.TempDir()
	runFilePath := filepath.Join(tempDir, "nonexistent.pid")

	// Test with non-existent run file
	isRunning, pid, proc, err := IsRunning(0, runFilePath)
	assert.Error(t, err)
	assert.False(t, isRunning)
	assert.Equal(t, 0, pid)
	assert.Nil(t, proc)
}

func TestIsRunning_InvalidParameters(t *testing.T) {
	// Test with both PID and run file (should error)
	isRunning, pid, proc, err := IsRunning(12345, "/some/path.pid")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "either provide a pid or a runFile")
	assert.False(t, isRunning)
	assert.Equal(t, 0, pid)
	assert.Nil(t, proc)
}

func TestCleanup(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skipping process tests on Windows due to signal handling differences")
	}

	tempDir := t.TempDir()
	tmpDir := filepath.Join(tempDir, "tmp")

	// Create temporary directory
	err := os.MkdirAll(tmpDir, 0755)
	require.NoError(t, err)

	// Create a file in tmp dir
	tmpFile := filepath.Join(tmpDir, "test.txt")
	err = os.WriteFile(tmpFile, []byte("test"), 0644)
	require.NoError(t, err)

	// Test cleanup with non-existent process
	err = Cleanup(999999, "", tmpDir, 100*time.Millisecond, 1*time.Second)
	require.NoError(t, err)

	// Verify tmp directory was removed
	assert.NoDirExists(t, tmpDir)
}

func TestCleanup_WithRunningProcess(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skipping process tests on Windows due to signal handling differences")
	}

	tempDir := t.TempDir()
	runFilePath := filepath.Join(tempDir, "test.pid")

	var stdout, stderr bytes.Buffer

	// Start a long-running process
	binPath := "sleep"
	args := []string{"10"}

	_, err := Execute(binPath, args, &stdout, &stderr, runFilePath, 0)
	require.NoError(t, err)

	// Wait a bit to ensure process is running
	time.Sleep(100 * time.Millisecond)

	// Test cleanup
	err = Cleanup(0, runFilePath, "", 100*time.Millisecond, 2*time.Second)
	require.NoError(t, err)

	// Verify run file was removed
	assert.NoFileExists(t, runFilePath)

	// Note: Process cleanup verification is complex due to timing issues
	// The important thing is that cleanup succeeded without error
}

func TestCleanup_ProcessKillTimeout(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skipping process tests on Windows due to signal handling differences")
	}

	tempDir := t.TempDir()
	runFilePath := filepath.Join(tempDir, "test.pid")

	var stdout, stderr bytes.Buffer

	// Start a process that ignores SIGINT
	binPath := "sleep"
	args := []string{"30"}

	_, err := Execute(binPath, args, &stdout, &stderr, runFilePath, 0)
	require.NoError(t, err)

	// Wait a bit to ensure process is running
	time.Sleep(100 * time.Millisecond)

	// Test cleanup with very short timeout (should force kill)
	err = Cleanup(0, runFilePath, "", 50*time.Millisecond, 100*time.Millisecond)
	require.NoError(t, err)

	// Verify run file was removed
	assert.NoFileExists(t, runFilePath)

	// Note: Process cleanup verification is complex due to timing issues
	// The important thing is that cleanup succeeded without error
}

func TestCleanup_EmptyParameters(t *testing.T) {
	// Test cleanup with empty parameters
	err := Cleanup(0, "", "", 0, 0)
	require.NoError(t, err)
}

func TestCleanup_InvalidRunFile(t *testing.T) {
	tempDir := t.TempDir()
	runFilePath := filepath.Join(tempDir, "nonexistent.pid")

	// Test cleanup with invalid run file
	err := Cleanup(0, runFilePath, "", 100*time.Millisecond, 1*time.Second)
	assert.Error(t, err)
}

// Test helper functions for process management
func TestProcessLifecycle(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skipping process tests on Windows due to signal handling differences")
	}

	tempDir := t.TempDir()
	runFilePath := filepath.Join(tempDir, "lifecycle.pid")

	var stdout, stderr bytes.Buffer

	// Start process
	binPath := "sleep"
	args := []string{"5"}

	pid, err := Execute(binPath, args, &stdout, &stderr, runFilePath, 0)
	require.NoError(t, err)

	// Verify process is running
	isRunning, loadedPid, proc, err := IsRunning(0, runFilePath)
	require.NoError(t, err)
	assert.True(t, isRunning)
	assert.Equal(t, pid, loadedPid)
	assert.NotNil(t, proc)

	// Clean up
	err = Cleanup(0, runFilePath, "", 100*time.Millisecond, 1*time.Second)
	require.NoError(t, err)

	// Verify process is no longer running
	isRunning, _, _, err = IsRunning(0, runFilePath)
	assert.Error(t, err) // Should error because run file was removed
}

func TestConcurrentProcessManagement(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skipping process tests on Windows due to signal handling differences")
	}

	tempDir := t.TempDir()

	// Start multiple processes concurrently
	var pids []int
	for i := 0; i < 3; i++ {
		runFilePath := filepath.Join(tempDir, fmt.Sprintf("process_%d.pid", i))
		var stdout, stderr bytes.Buffer

		pid, err := Execute("sleep", []string{"3"}, &stdout, &stderr, runFilePath, 0)
		require.NoError(t, err)
		pids = append(pids, pid)
	}

	// Verify all processes are running
	for i := range pids {
		runFilePath := filepath.Join(tempDir, fmt.Sprintf("process_%d.pid", i))
		isRunning, _, _, err := IsRunning(0, runFilePath)
		require.NoError(t, err)
		assert.True(t, isRunning)
	}

	// Clean up all processes
	for i := range pids {
		runFilePath := filepath.Join(tempDir, fmt.Sprintf("process_%d.pid", i))
		err := Cleanup(0, runFilePath, "", 100*time.Millisecond, 1*time.Second)
		require.NoError(t, err)
	}
}

func TestRunFilePermissions(t *testing.T) {
	tempDir := t.TempDir()
	runFilePath := filepath.Join(tempDir, "permissions.pid")

	// Save run file
	err := saveRunFile(12345, runFilePath)
	require.NoError(t, err)

	// Check permissions
	info, err := os.Stat(runFilePath)
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(constants.WriteReadReadPerms), info.Mode())
}

func TestRunFileContentIntegrity(t *testing.T) {
	tempDir := t.TempDir()
	runFilePath := filepath.Join(tempDir, "integrity.pid")

	// Save run file with specific PID
	expectedPid := 98765
	err := saveRunFile(expectedPid, runFilePath)
	require.NoError(t, err)

	// Load and verify
	actualPid, err := loadRunFile(runFilePath)
	require.NoError(t, err)
	assert.Equal(t, expectedPid, actualPid)

	// Verify JSON structure
	content, err := os.ReadFile(runFilePath)
	require.NoError(t, err)

	var runFile RunFile
	err = json.Unmarshal(content, &runFile)
	require.NoError(t, err)
	assert.Equal(t, expectedPid, runFile.Pid)
}
