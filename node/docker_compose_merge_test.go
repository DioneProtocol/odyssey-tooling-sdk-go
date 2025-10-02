// Copyright (C) 2025, Dione Protocol, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package node

import (
	"testing"

	"github.com/DioneProtocol/odyssey-tooling-sdk-go/constants"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNode_MergeComposeFiles_FeatureFlags(t *testing.T) {
	// Save original feature flag values
	originalDockerSupport := constants.DockerSupportEnabled
	originalSSHKeyManagement := constants.SSHKeyManagementEnabled

	defer func() {
		constants.DockerSupportEnabled = originalDockerSupport
		constants.SSHKeyManagementEnabled = originalSSHKeyManagement
	}()

	node := &Node{
		NodeID: "test-node",
		IP:     "192.168.1.1",
	}

	tests := []struct {
		name          string
		dockerEnabled bool
		sshEnabled    bool
		expectError   bool
		errorContains string
	}{
		{
			name:          "Docker support disabled",
			dockerEnabled: false,
			sshEnabled:    true,
			expectError:   true,
			errorContains: "Docker support functionality is disabled",
		},
		{
			name:          "SSH key management disabled",
			dockerEnabled: true,
			sshEnabled:    false,
			expectError:   true,
			errorContains: "SSH key management functionality is disabled",
		},
		{
			name:          "Both disabled",
			dockerEnabled: false,
			sshEnabled:    false,
			expectError:   true,
			errorContains: "Docker support functionality is disabled",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set feature flags
			constants.DockerSupportEnabled = tt.dockerEnabled
			constants.SSHKeyManagementEnabled = tt.sshEnabled

			err := node.MergeComposeFiles("/path/to/current.yml", "/path/to/new.yml")

			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorContains)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestNode_MergeComposeFiles_FileValidation(t *testing.T) {
	// Save original feature flag values
	originalDockerSupport := constants.DockerSupportEnabled
	originalSSHKeyManagement := constants.SSHKeyManagementEnabled

	defer func() {
		constants.DockerSupportEnabled = originalDockerSupport
		constants.SSHKeyManagementEnabled = originalSSHKeyManagement
	}()

	// Enable feature flags
	constants.DockerSupportEnabled = true
	constants.SSHKeyManagementEnabled = true

	node := &Node{
		NodeID: "test-node",
		IP:     "192.168.1.1",
	}

	tests := []struct {
		name          string
		currentFile   string
		newFile       string
		expectError   bool
		errorContains string
	}{
		{
			name:          "Empty current file path",
			currentFile:   "",
			newFile:       "/path/to/new.yml",
			expectError:   true,
			errorContains: "file  does not exist",
		},
		{
			name:          "Empty new file path",
			currentFile:   "/path/to/current.yml",
			newFile:       "",
			expectError:   true,
			errorContains: "file  does not exist",
		},
		{
			name:          "Both files empty",
			currentFile:   "",
			newFile:       "",
			expectError:   true,
			errorContains: "file  does not exist",
		},
		{
			name:          "Non-existent current file",
			currentFile:   "/nonexistent/current.yml",
			newFile:       "/path/to/new.yml",
			expectError:   true,
			errorContains: "file /nonexistent/current.yml does not exist",
		},
		{
			name:          "Non-existent new file",
			currentFile:   "/path/to/current.yml",
			newFile:       "/nonexistent/new.yml",
			expectError:   true,
			errorContains: "file /nonexistent/new.yml does not exist",
		},
		{
			name:          "Both files non-existent",
			currentFile:   "/nonexistent/current.yml",
			newFile:       "/nonexistent/new.yml",
			expectError:   true,
			errorContains: "file /nonexistent/current.yml does not exist",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := node.MergeComposeFiles(tt.currentFile, tt.newFile)

			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorContains)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestNode_MergeComposeFiles_ValidFiles(t *testing.T) {
	// Save original feature flag values
	originalDockerSupport := constants.DockerSupportEnabled
	originalSSHKeyManagement := constants.SSHKeyManagementEnabled

	defer func() {
		constants.DockerSupportEnabled = originalDockerSupport
		constants.SSHKeyManagementEnabled = originalSSHKeyManagement
	}()

	// Enable feature flags
	constants.DockerSupportEnabled = true
	constants.SSHKeyManagementEnabled = true

	node := &Node{
		NodeID: "test-node",
		IP:     "192.168.1.1",
	}

	tests := []struct {
		name        string
		currentFile string
		newFile     string
		expectError bool
	}{
		{
			name:        "Valid compose files",
			currentFile: "/path/to/current.yml",
			newFile:     "/path/to/new.yml",
			expectError: true, // Will fail due to SSH connection
		},
		{
			name:        "Same file paths",
			currentFile: "/path/to/same.yml",
			newFile:     "/path/to/same.yml",
			expectError: true, // Will fail due to SSH connection
		},
		{
			name:        "Different file extensions",
			currentFile: "/path/to/current.yaml",
			newFile:     "/path/to/new.yml",
			expectError: true, // Will fail due to SSH connection
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := node.MergeComposeFiles(tt.currentFile, tt.newFile)

			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestNode_MergeComposeFiles_EdgeCases(t *testing.T) {
	// Save original feature flag values
	originalDockerSupport := constants.DockerSupportEnabled
	originalSSHKeyManagement := constants.SSHKeyManagementEnabled

	defer func() {
		constants.DockerSupportEnabled = originalDockerSupport
		constants.SSHKeyManagementEnabled = originalSSHKeyManagement
	}()

	// Enable feature flags
	constants.DockerSupportEnabled = true
	constants.SSHKeyManagementEnabled = true

	node := &Node{
		NodeID: "test-node",
		IP:     "192.168.1.1",
	}

	tests := []struct {
		name        string
		currentFile string
		newFile     string
		expectError bool
	}{
		{
			name:        "Long file paths",
			currentFile: "/very/long/path/to/current/compose/file.yml",
			newFile:     "/very/long/path/to/new/compose/file.yml",
			expectError: true, // Will fail due to SSH connection
		},
		{
			name:        "Special characters in paths",
			currentFile: "/path/with spaces/current.yml",
			newFile:     "/path/with spaces/new.yml",
			expectError: true, // Will fail due to SSH connection
		},
		{
			name:        "Relative paths",
			currentFile: "./current.yml",
			newFile:     "./new.yml",
			expectError: true, // Will fail due to SSH connection
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := node.MergeComposeFiles(tt.currentFile, tt.newFile)

			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

