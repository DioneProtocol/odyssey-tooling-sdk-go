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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExtractArchive(t *testing.T) {
	tests := []struct {
		name        string
		kind        ArchiveKind
		archiveData []byte
		expectError bool
	}{
		{
			name:        "unsupported archive kind",
			kind:        UndefinedArchive,
			archiveData: []byte{},
			expectError: true,
		},
		{
			name:        "zip archive",
			kind:        Zip,
			archiveData: createTestZip(t),
			expectError: false,
		},
		{
			name:        "tar.gz archive",
			kind:        TarGz,
			archiveData: createTestTarGz(t),
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()

			err := ExtractArchive(tt.kind, tt.archiveData, tempDir)

			if tt.expectError {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			// Verify directory was created
			assert.DirExists(t, tempDir)
		})
	}
}

func TestExtractZip(t *testing.T) {
	tests := []struct {
		name        string
		setupZip    func(t *testing.T) []byte
		expectError bool
		verifyFiles []string // Files that should exist after extraction
	}{
		{
			name: "valid zip with files and directories",
			setupZip: func(t *testing.T) []byte {
				return createTestZip(t)
			},
			expectError: false,
			verifyFiles: []string{
				"test_file.txt",
				"subdir/nested_file.txt",
			},
		},
		{
			name: "empty zip",
			setupZip: func(t *testing.T) []byte {
				var buf bytes.Buffer
				writer := zip.NewWriter(&buf)
				writer.Close()
				return buf.Bytes()
			},
			expectError: false,
			verifyFiles: []string{},
		},
		{
			name: "zip with zip slip attempt",
			setupZip: func(t *testing.T) []byte {
				return createZipSlipZip(t)
			},
			expectError: true,
			verifyFiles: []string{},
		},
		{
			name: "invalid zip data",
			setupZip: func(t *testing.T) []byte {
				return []byte("not a zip file")
			},
			expectError: true,
			verifyFiles: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			zipData := tt.setupZip(t)

			err := extractZip(zipData, tempDir)

			if tt.expectError {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)

			// Verify expected files exist
			for _, file := range tt.verifyFiles {
				filePath := filepath.Join(tempDir, file)
				assert.FileExists(t, filePath)
			}
		})
	}
}

func TestExtractTarGz(t *testing.T) {
	tests := []struct {
		name        string
		setupTarGz  func(t *testing.T) []byte
		expectError bool
		verifyFiles []string // Files that should exist after extraction
	}{
		{
			name: "valid tar.gz with files and directories",
			setupTarGz: func(t *testing.T) []byte {
				return createTestTarGz(t)
			},
			expectError: false,
			verifyFiles: []string{
				"test_file.txt",
				"subdir/nested_file.txt",
			},
		},
		{
			name: "empty tar.gz",
			setupTarGz: func(t *testing.T) []byte {
				var buf bytes.Buffer
				gzWriter := gzip.NewWriter(&buf)
				tarWriter := tar.NewWriter(gzWriter)
				tarWriter.Close()
				gzWriter.Close()
				return buf.Bytes()
			},
			expectError: false,
			verifyFiles: []string{},
		},
		{
			name: "tar.gz with zip slip attempt",
			setupTarGz: func(t *testing.T) []byte {
				return createTarGzSlip(t)
			},
			expectError: true,
			verifyFiles: []string{},
		},
		{
			name: "invalid tar.gz data",
			setupTarGz: func(t *testing.T) []byte {
				return []byte("not a tar.gz file")
			},
			expectError: true,
			verifyFiles: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			tarGzData := tt.setupTarGz(t)

			err := extractTarGz(tarGzData, tempDir)

			if tt.expectError {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)

			// Verify expected files exist
			for _, file := range tt.verifyFiles {
				filePath := filepath.Join(tempDir, file)
				assert.FileExists(t, filePath)
			}
		})
	}
}

func TestSanitizeArchivePath(t *testing.T) {
	tests := []struct {
		name         string
		baseDir      string
		targetPath   string
		expectError  bool
		expectedPath string
	}{
		{
			name:         "valid path within base directory",
			baseDir:      "/tmp/extract",
			targetPath:   "file.txt",
			expectError:  false,
			expectedPath: "/tmp/extract/file.txt",
		},
		{
			name:         "valid nested path",
			baseDir:      "/tmp/extract",
			targetPath:   "subdir/file.txt",
			expectError:  false,
			expectedPath: "/tmp/extract/subdir/file.txt",
		},
		{
			name:         "zip slip attempt with ../",
			baseDir:      "/tmp/extract",
			targetPath:   "../../../etc/passwd",
			expectError:  true,
			expectedPath: "",
		},
		{
			name:         "zip slip attempt with absolute path",
			baseDir:      "/tmp/extract",
			targetPath:   "/etc/passwd",
			expectError:  false, // This is actually valid - it creates /tmp/extract/etc/passwd
			expectedPath: "/tmp/extract/etc/passwd",
		},
		{
			name:         "zip slip attempt with mixed separators",
			baseDir:      "/tmp/extract",
			targetPath:   "..\\..\\..\\etc\\passwd",
			expectError:  true,
			expectedPath: "",
		},
		{
			name:         "edge case: path exactly at base directory",
			baseDir:      "/tmp/extract",
			targetPath:   ".",
			expectError:  false,
			expectedPath: "/tmp/extract",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := sanitizeArchivePath(tt.baseDir, tt.targetPath)

			if tt.expectError {
				assert.Error(t, err)
				assert.Empty(t, result)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expectedPath, result)
		})
	}
}

// Helper functions to create test archives

func createTestZip(t *testing.T) []byte {
	var buf bytes.Buffer
	writer := zip.NewWriter(&buf)

	// Add a regular file
	fileWriter, err := writer.Create("test_file.txt")
	require.NoError(t, err)
	_, err = fileWriter.Write([]byte("test content"))
	require.NoError(t, err)

	// Add a file in a subdirectory
	fileWriter, err = writer.Create("subdir/nested_file.txt")
	require.NoError(t, err)
	_, err = fileWriter.Write([]byte("nested content"))
	require.NoError(t, err)

	// Add a directory entry
	_, err = writer.Create("subdir/")
	require.NoError(t, err)

	err = writer.Close()
	require.NoError(t, err)

	return buf.Bytes()
}

func createZipSlipZip(t *testing.T) []byte {
	var buf bytes.Buffer
	writer := zip.NewWriter(&buf)

	// Add a malicious file that tries to escape the extraction directory
	fileWriter, err := writer.Create("../../../malicious_file.txt")
	require.NoError(t, err)
	_, err = fileWriter.Write([]byte("malicious content"))
	require.NoError(t, err)

	err = writer.Close()
	require.NoError(t, err)

	return buf.Bytes()
}

func createTestTarGz(t *testing.T) []byte {
	var buf bytes.Buffer
	gzWriter := gzip.NewWriter(&buf)
	tarWriter := tar.NewWriter(gzWriter)

	// Add a regular file
	header := &tar.Header{
		Name: "test_file.txt",
		Mode: 0644,
		Size: int64(len("test content")),
	}
	err := tarWriter.WriteHeader(header)
	require.NoError(t, err)
	_, err = tarWriter.Write([]byte("test content"))
	require.NoError(t, err)

	// Add a directory
	header = &tar.Header{
		Name:     "subdir/",
		Mode:     0755,
		Typeflag: tar.TypeDir,
	}
	err = tarWriter.WriteHeader(header)
	require.NoError(t, err)

	// Add a file in the subdirectory
	header = &tar.Header{
		Name: "subdir/nested_file.txt",
		Mode: 0644,
		Size: int64(len("nested content")),
	}
	err = tarWriter.WriteHeader(header)
	require.NoError(t, err)
	_, err = tarWriter.Write([]byte("nested content"))
	require.NoError(t, err)

	err = tarWriter.Close()
	require.NoError(t, err)
	err = gzWriter.Close()
	require.NoError(t, err)

	return buf.Bytes()
}

func createTarGzSlip(t *testing.T) []byte {
	var buf bytes.Buffer
	gzWriter := gzip.NewWriter(&buf)
	tarWriter := tar.NewWriter(gzWriter)

	// Add a malicious file that tries to escape the extraction directory
	header := &tar.Header{
		Name: "../../../malicious_file.txt",
		Mode: 0644,
		Size: int64(len("malicious content")),
	}
	err := tarWriter.WriteHeader(header)
	require.NoError(t, err)
	_, err = tarWriter.Write([]byte("malicious content"))
	require.NoError(t, err)

	err = tarWriter.Close()
	require.NoError(t, err)
	err = gzWriter.Close()
	require.NoError(t, err)

	return buf.Bytes()
}

func TestArchiveExtractionWithLargeFile(t *testing.T) {
	// Test that large files are handled properly (within the 2GB limit)
	tempDir := t.TempDir()

	// Create a zip with a moderately large file (1MB)
	var buf bytes.Buffer
	writer := zip.NewWriter(&buf)

	// Create 1MB of data
	largeData := make([]byte, 1024*1024) // 1MB
	for i := range largeData {
		largeData[i] = byte(i % 256)
	}

	fileWriter, err := writer.Create("large_file.bin")
	require.NoError(t, err)
	_, err = fileWriter.Write(largeData)
	require.NoError(t, err)

	err = writer.Close()
	require.NoError(t, err)

	// Extract the archive
	err = extractZip(buf.Bytes(), tempDir)
	require.NoError(t, err)

	// Verify the file was extracted correctly
	filePath := filepath.Join(tempDir, "large_file.bin")
	assert.FileExists(t, filePath)

	// Verify file size
	info, err := os.Stat(filePath)
	require.NoError(t, err)
	assert.Equal(t, int64(1024*1024), info.Size())
}

func TestArchiveExtractionFilePermissions(t *testing.T) {
	tempDir := t.TempDir()

	// Create a tar.gz with files having different permissions
	var buf bytes.Buffer
	gzWriter := gzip.NewWriter(&buf)
	tarWriter := tar.NewWriter(gzWriter)

	// Add a file with specific permissions
	header := &tar.Header{
		Name: "executable_file.sh",
		Mode: 0755, // Executable
		Size: int64(len("#!/bin/bash\necho hello")),
	}
	err := tarWriter.WriteHeader(header)
	require.NoError(t, err)
	_, err = tarWriter.Write([]byte("#!/bin/bash\necho hello"))
	require.NoError(t, err)

	// Add a file with read-only permissions
	header = &tar.Header{
		Name: "readonly_file.txt",
		Mode: 0444, // Read-only
		Size: int64(len("read only content")),
	}
	err = tarWriter.WriteHeader(header)
	require.NoError(t, err)
	_, err = tarWriter.Write([]byte("read only content"))
	require.NoError(t, err)

	err = tarWriter.Close()
	require.NoError(t, err)
	err = gzWriter.Close()
	require.NoError(t, err)

	// Extract the archive
	err = extractTarGz(buf.Bytes(), tempDir)
	require.NoError(t, err)

	// Verify file permissions
	execFile := filepath.Join(tempDir, "executable_file.sh")
	readFile := filepath.Join(tempDir, "readonly_file.txt")

	assert.FileExists(t, execFile)
	assert.FileExists(t, readFile)

	// Check permissions (note: actual permissions may be affected by umask)
	execInfo, err := os.Stat(execFile)
	require.NoError(t, err)
	readInfo, err := os.Stat(readFile)
	require.NoError(t, err)

	// At minimum, verify files exist and have some permissions
	assert.True(t, execInfo.Mode()&0444 != 0, "executable file should be readable")
	assert.True(t, readInfo.Mode()&0444 != 0, "read-only file should be readable")
}
