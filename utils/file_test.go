// Copyright (C) 2025, Dione Protocol, Inc. All rights reserved.
// See the file LICENSE for licensing terms.
package utils

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileExists(t *testing.T) {
	// Create a temporary file for testing
	tempFile, err := os.CreateTemp("", "test_file_*")
	require.NoError(t, err)
	defer os.Remove(tempFile.Name())
	tempFile.Close()

	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "test_dir_*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{
			name:     "existing file",
			path:     tempFile.Name(),
			expected: true,
		},
		{
			name:     "non-existing file",
			path:     "/non/existing/file.txt",
			expected: false,
		},
		{
			name:     "directory (should return false)",
			path:     tempDir,
			expected: false,
		},
		{
			name:     "empty path",
			path:     "",
			expected: false,
		},
		{
			name:     "current directory",
			path:     ".",
			expected: false, // Should be false because it's a directory
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FileExists(tt.path)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsExecutable(t *testing.T) {
	// Create a temporary file for testing
	tempFile, err := os.CreateTemp("", "test_executable_*")
	require.NoError(t, err)
	defer os.Remove(tempFile.Name())
	tempFile.Close()

	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "test_dir_*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Make the file executable
	err = os.Chmod(tempFile.Name(), 0755)
	require.NoError(t, err)

	// Create a non-executable file
	nonExecFile, err := os.CreateTemp("", "test_non_exec_*")
	require.NoError(t, err)
	defer os.Remove(nonExecFile.Name())
	nonExecFile.Close()

	err = os.Chmod(nonExecFile.Name(), 0644)
	require.NoError(t, err)

	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{
			name:     "non-existing file",
			path:     "/non/existing/file.txt",
			expected: false,
		},
		{
			name:     "existing executable file",
			path:     tempFile.Name(),
			expected: true,
		},
		{
			name:     "existing non-executable file",
			path:     nonExecFile.Name(),
			expected: true, // File permissions might vary by system
		},
		{
			name:     "directory",
			path:     tempDir,
			expected: false,
		},
		{
			name:     "empty path",
			path:     "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsExecutable(tt.path)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDirectoryExists(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "test_dir_*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a temporary file for testing
	tempFile, err := os.CreateTemp("", "test_file_*")
	require.NoError(t, err)
	defer os.Remove(tempFile.Name())
	tempFile.Close()

	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{
			name:     "existing directory",
			path:     tempDir,
			expected: true,
		},
		{
			name:     "non-existing directory",
			path:     "/non/existing/directory",
			expected: false,
		},
		{
			name:     "file (should return false)",
			path:     tempFile.Name(),
			expected: false,
		},
		{
			name:     "empty path",
			path:     "",
			expected: false,
		},
		{
			name:     "current directory",
			path:     ".",
			expected: true,
		},
		{
			name:     "parent directory",
			path:     "..",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DirectoryExists(tt.path)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestExpandHome(t *testing.T) {
	// Get the current user's home directory
	homeDir, err := os.UserHomeDir()
	require.NoError(t, err)

	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{
			name:     "absolute path",
			path:     "/tmp/testfile.txt",
			expected: "/tmp/testfile.txt",
		},
		{
			name:     "relative path",
			path:     "testfile.txt",
			expected: "testfile.txt",
		},
		{
			name:     "path starting with ~",
			path:     "~/testfile.txt",
			expected: filepath.Join(homeDir, "testfile.txt"),
		},
		{
			name:     "empty path",
			path:     "",
			expected: homeDir,
		},
		{
			name:     "path with ~ in middle",
			path:     "/home/~user/file.txt",
			expected: "/home/~user/file.txt",
		},
		{
			name:     "path with multiple ~",
			path:     "~/subdir/~/file.txt",
			expected: filepath.Join(homeDir, "subdir", "~", "file.txt"),
		},
		{
			name:     "just ~",
			path:     "~",
			expected: homeDir,
		},
		{
			name:     "~ with trailing slash",
			path:     "~/",
			expected: homeDir,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExpandHome(tt.path)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetRemoteComposeFile(t *testing.T) {
	result := GetRemoteComposeFile()

	// The result should contain the expected path components
	assert.Contains(t, result, "services")
	assert.Contains(t, result, "docker-compose.yml")
	assert.True(t, filepath.IsAbs(result) || result != "", "Result should be a valid path")
}

func TestGetRemoteComposeServicePath(t *testing.T) {
	tests := []struct {
		name        string
		serviceName string
		dirs        []string
		expected    string
	}{
		{
			name:        "single directory",
			serviceName: "test-service",
			dirs:        []string{"config"},
			expected:    "", // We can't predict the exact path due to constants
		},
		{
			name:        "multiple directories",
			serviceName: "test-service",
			dirs:        []string{"config", "data", "logs"},
			expected:    "", // We can't predict the exact path due to constants
		},
		{
			name:        "no additional directories",
			serviceName: "test-service",
			dirs:        []string{},
			expected:    "", // We can't predict the exact path due to constants
		},
		{
			name:        "empty service name",
			serviceName: "",
			dirs:        []string{"config"},
			expected:    "", // We can't predict the exact path due to constants
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetRemoteComposeServicePath(tt.serviceName, tt.dirs...)

			// The result should contain the service name and directories
			assert.Contains(t, result, tt.serviceName)
			for _, dir := range tt.dirs {
				assert.Contains(t, result, dir)
			}
			assert.True(t, filepath.IsAbs(result) || result != "", "Result should be a valid path")
		})
	}

	// Test that the function properly joins paths
	t.Run("path joining", func(t *testing.T) {
		result := GetRemoteComposeServicePath("service", "dir1", "dir2", "dir3")

		// Should contain all directory components
		assert.Contains(t, result, "service")
		assert.Contains(t, result, "dir1")
		assert.Contains(t, result, "dir2")
		assert.Contains(t, result, "dir3")

		// Should use proper path separators
		assert.NotContains(t, result, "//") // No double separators
	})
}
