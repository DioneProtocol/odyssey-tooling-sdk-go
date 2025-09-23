// Copyright (C) 2025, Dione Protocol, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package install

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestArchiveBuilder helps build test archives for testing
type TestArchiveBuilder struct {
	t *testing.T
}

// NewTestArchiveBuilder creates a new test archive builder
func NewTestArchiveBuilder(t *testing.T) *TestArchiveBuilder {
	return &TestArchiveBuilder{t: t}
}

// CreateZipWithFiles creates a zip archive with the specified files
func (b *TestArchiveBuilder) CreateZipWithFiles(files map[string]string) []byte {
	var buf bytes.Buffer
	writer := zip.NewWriter(&buf)

	for filePath, content := range files {
		fileWriter, err := writer.Create(filePath)
		require.NoError(b.t, err)
		_, err = fileWriter.Write([]byte(content))
		require.NoError(b.t, err)
	}

	err := writer.Close()
	require.NoError(b.t, err)

	return buf.Bytes()
}

// CreateTarGzWithFiles creates a tar.gz archive with the specified files
func (b *TestArchiveBuilder) CreateTarGzWithFiles(files map[string]string) []byte {
	var buf bytes.Buffer
	gzWriter := gzip.NewWriter(&buf)
	tarWriter := tar.NewWriter(gzWriter)

	for filePath, content := range files {
		header := &tar.Header{
			Name: filePath,
			Mode: 0644,
			Size: int64(len(content)),
		}
		err := tarWriter.WriteHeader(header)
		require.NoError(b.t, err)
		_, err = tarWriter.Write([]byte(content))
		require.NoError(b.t, err)
	}

	err := tarWriter.Close()
	require.NoError(b.t, err)
	err = gzWriter.Close()
	require.NoError(b.t, err)

	return buf.Bytes()
}

// CreateZipWithZipSlip creates a zip archive that attempts zip slip attacks
func (b *TestArchiveBuilder) CreateZipWithZipSlip() []byte {
	maliciousFiles := map[string]string{
		"../../../etc/passwd":           "malicious content 1",
		"..\\..\\..\\windows\\system32": "malicious content 2",
		"normal_file.txt":               "normal content",
		"subdir/../malicious.txt":       "malicious content 3",
	}
	return b.CreateZipWithFiles(maliciousFiles)
}

// CreateTarGzWithZipSlip creates a tar.gz archive that attempts zip slip attacks
func (b *TestArchiveBuilder) CreateTarGzWithZipSlip() []byte {
	maliciousFiles := map[string]string{
		"../../../etc/passwd":           "malicious content 1",
		"..\\..\\..\\windows\\system32": "malicious content 2",
		"normal_file.txt":               "normal content",
		"subdir/../malicious.txt":       "malicious content 3",
	}
	return b.CreateTarGzWithFiles(maliciousFiles)
}

// MockHTTPClient provides a mock implementation for HTTP requests
type MockHTTPClient struct {
	responses map[string][]byte
	errors    map[string]error
}

// NewMockHTTPClient creates a new mock HTTP client
func NewMockHTTPClient() *MockHTTPClient {
	return &MockHTTPClient{
		responses: make(map[string][]byte),
		errors:    make(map[string]error),
	}
}

// SetResponse sets a mock response for a URL
func (m *MockHTTPClient) SetResponse(url string, data []byte) {
	m.responses[url] = data
}

// SetError sets a mock error for a URL
func (m *MockHTTPClient) SetError(url string, err error) {
	m.errors[url] = err
}

// Get simulates an HTTP GET request
func (m *MockHTTPClient) Get(url string) ([]byte, error) {
	if err, exists := m.errors[url]; exists {
		return nil, err
	}
	if data, exists := m.responses[url]; exists {
		return data, nil
	}
	return nil, &MockHTTPError{URL: url, Message: "no mock response set"}
}

// MockHTTPError represents an error from the mock HTTP client
type MockHTTPError struct {
	URL     string
	Message string
}

func (e *MockHTTPError) Error() string {
	return e.Message + " for URL: " + e.URL
}

// TestFileSystem provides utilities for testing file operations
type TestFileSystem struct {
	t       *testing.T
	baseDir string
	created []string // Track created files for cleanup
}

// NewTestFileSystem creates a new test file system
func NewTestFileSystem(t *testing.T) *TestFileSystem {
	baseDir := t.TempDir()
	return &TestFileSystem{
		t:       t,
		baseDir: baseDir,
		created: make([]string, 0),
	}
}

// CreateFile creates a file with the given content
func (fs *TestFileSystem) CreateFile(relativePath, content string) string {
	fullPath := filepath.Join(fs.baseDir, relativePath)

	// Create parent directories
	err := os.MkdirAll(filepath.Dir(fullPath), 0755)
	require.NoError(fs.t, err)

	// Create the file
	err = os.WriteFile(fullPath, []byte(content), 0644)
	require.NoError(fs.t, err)

	fs.created = append(fs.created, fullPath)
	return fullPath
}

// CreateExecutableFile creates an executable file
func (fs *TestFileSystem) CreateExecutableFile(relativePath, content string) string {
	fullPath := fs.CreateFile(relativePath, content)

	// Make it executable
	err := os.Chmod(fullPath, 0755)
	require.NoError(fs.t, err)

	return fullPath
}

// CreateDirectory creates a directory
func (fs *TestFileSystem) CreateDirectory(relativePath string) string {
	fullPath := filepath.Join(fs.baseDir, relativePath)
	err := os.MkdirAll(fullPath, 0755)
	require.NoError(fs.t, err)

	fs.created = append(fs.created, fullPath)
	return fullPath
}

// GetPath returns the full path for a relative path
func (fs *TestFileSystem) GetPath(relativePath string) string {
	return filepath.Join(fs.baseDir, relativePath)
}

// AssertFileExists checks if a file exists
func (fs *TestFileSystem) AssertFileExists(relativePath string) {
	fullPath := filepath.Join(fs.baseDir, relativePath)
	require.FileExists(fs.t, fullPath)
}

// AssertFileNotExists checks if a file does not exist
func (fs *TestFileSystem) AssertFileNotExists(relativePath string) {
	fullPath := filepath.Join(fs.baseDir, relativePath)
	require.NoFileExists(fs.t, fullPath)
}

// AssertFileContent checks if a file has the expected content
func (fs *TestFileSystem) AssertFileContent(relativePath, expectedContent string) {
	fullPath := filepath.Join(fs.baseDir, relativePath)
	content, err := os.ReadFile(fullPath)
	require.NoError(fs.t, err)
	require.Equal(fs.t, expectedContent, string(content))
}

// AssertFileExecutable checks if a file is executable
func (fs *TestFileSystem) AssertFileExecutable(relativePath string) {
	fullPath := filepath.Join(fs.baseDir, relativePath)
	info, err := os.Stat(fullPath)
	require.NoError(fs.t, err)
	require.True(fs.t, info.Mode()&0111 != 0, "file should be executable")
}

// Cleanup removes all created files (usually called automatically by t.TempDir())
func (fs *TestFileSystem) Cleanup() {
	for _, path := range fs.created {
		os.RemoveAll(path)
	}
}

// TestArchiveValidator helps validate archive extraction results
type TestArchiveValidator struct {
	t       *testing.T
	baseDir string
}

// NewTestArchiveValidator creates a new archive validator
func NewTestArchiveValidator(t *testing.T, baseDir string) *TestArchiveValidator {
	return &TestArchiveValidator{
		t:       t,
		baseDir: baseDir,
	}
}

// AssertFilesExist checks that the specified files exist after extraction
func (v *TestArchiveValidator) AssertFilesExist(files []string) {
	for _, file := range files {
		fullPath := filepath.Join(v.baseDir, file)
		require.FileExists(v.t, fullPath, "file should exist after extraction: %s", file)
	}
}

// AssertFilesNotExist checks that the specified files do not exist after extraction
func (v *TestArchiveValidator) AssertFilesNotExist(files []string) {
	for _, file := range files {
		fullPath := filepath.Join(v.baseDir, file)
		require.NoFileExists(v.t, fullPath, "file should not exist after extraction: %s", file)
	}
}

// AssertFileContent checks that a file has the expected content
func (v *TestArchiveValidator) AssertFileContent(relativePath, expectedContent string) {
	fullPath := filepath.Join(v.baseDir, relativePath)
	content, err := os.ReadFile(fullPath)
	require.NoError(v.t, err)
	require.Equal(v.t, expectedContent, string(content), "file content should match expected")
}

// AssertNoZipSlipFiles checks that no zip slip attack files were created
func (v *TestArchiveValidator) AssertNoZipSlipFiles() {
	// Check for common zip slip attack patterns
	zipSlipPatterns := []string{
		"../",
		"..\\",
		"/etc/passwd",
		"\\windows\\system32",
	}

	// Walk the directory and check for any files with zip slip patterns
	err := filepath.Walk(v.baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(v.baseDir, path)
		if err != nil {
			return err
		}

		for _, slipPattern := range zipSlipPatterns {
			if filepath.Base(relPath) == slipPattern ||
				filepath.Dir(relPath) == slipPattern {
				require.Fail(v.t, "zip slip attack detected",
					"found potentially malicious file: %s", relPath)
			}
		}

		return nil
	})
	require.NoError(v.t, err)
}
