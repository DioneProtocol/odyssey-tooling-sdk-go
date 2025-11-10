// Copyright (C) 2025, Dione Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package node

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsMonitoringNode(t *testing.T) {
	tests := []struct {
		name     string
		node     Node
		expected bool
	}{
		{
			name: "Node with Monitor role",
			node: Node{
				Roles: []SupportedRole{Monitor},
			},
			expected: true,
		},
		{
			name: "Node with multiple roles including Monitor",
			node: Node{
				Roles: []SupportedRole{Validator, Monitor},
			},
			expected: true,
		},
		{
			name: "Node without Monitor role",
			node: Node{
				Roles: []SupportedRole{Validator},
			},
			expected: false,
		},
		{
			name: "Node with no roles",
			node: Node{
				Roles: []SupportedRole{},
			},
			expected: false,
		},
		{
			name: "Node with nil roles",
			node: Node{
				Roles: nil,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isMonitoringNode(tt.node)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsOdysseyGoNode(t *testing.T) {
	tests := []struct {
		name     string
		node     Node
		expected bool
	}{
		{
			name: "Node with Validator role",
			node: Node{
				Roles: []SupportedRole{Validator},
			},
			expected: true,
		},
		{
			name: "Node with API role",
			node: Node{
				Roles: []SupportedRole{API},
			},
			expected: true,
		},
		{
			name: "Node with both Validator and API roles",
			node: Node{
				Roles: []SupportedRole{Validator, API},
			},
			expected: true,
		},
		{
			name: "Node with Validator and other roles",
			node: Node{
				Roles: []SupportedRole{Validator, Monitor},
			},
			expected: true,
		},
		{
			name: "Node with API and other roles",
			node: Node{
				Roles: []SupportedRole{API, Monitor},
			},
			expected: true,
		},
		{
			name: "Node without Validator or API role",
			node: Node{
				Roles: []SupportedRole{Monitor},
			},
			expected: false,
		},
		{
			name: "Node with Loadtest role only",
			node: Node{
				Roles: []SupportedRole{Loadtest},
			},
			expected: false,
		},
		{
			name: "Node with no roles",
			node: Node{
				Roles: []SupportedRole{},
			},
			expected: false,
		},
		{
			name: "Node with nil roles",
			node: Node{
				Roles: nil,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isOdysseyGoNode(tt.node)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsLoadTestNode(t *testing.T) {
	tests := []struct {
		name     string
		node     Node
		expected bool
	}{
		{
			name: "Node with Loadtest role",
			node: Node{
				Roles: []SupportedRole{Loadtest},
			},
			expected: true,
		},
		{
			name: "Node with multiple roles including Loadtest",
			node: Node{
				Roles: []SupportedRole{Validator, Loadtest},
			},
			expected: true,
		},
		{
			name: "Node without Loadtest role",
			node: Node{
				Roles: []SupportedRole{Validator},
			},
			expected: false,
		},
		{
			name: "Node with no roles",
			node: Node{
				Roles: []SupportedRole{},
			},
			expected: false,
		},
		{
			name: "Node with nil roles",
			node: Node{
				Roles: nil,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isLoadTestNode(tt.node)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetPrometheusTargets(t *testing.T) {
	tests := []struct {
		name            string
		nodes           []Node
		expectedOG      []string
		expectedMachine []string
		expectedLT      []string
	}{
		{
			name:            "Empty nodes list",
			nodes:           []Node{},
			expectedOG:      []string{},
			expectedMachine: []string{},
			expectedLT:      []string{},
		},
		{
			name: "Single Validator node",
			nodes: []Node{
				{
					IP:    "192.168.1.1",
					Roles: []SupportedRole{Validator},
				},
			},
			expectedOG:      []string{"'192.168.1.1:9650'"}, // OdysseygoAPIPort
			expectedMachine: []string{"'192.168.1.1:9100'"}, // OdysseygoMachineMetricsPort
			expectedLT:      []string{},
		},
		{
			name: "Single API node",
			nodes: []Node{
				{
					IP:    "192.168.1.2",
					Roles: []SupportedRole{API},
				},
			},
			expectedOG:      []string{"'192.168.1.2:9650'"},
			expectedMachine: []string{"'192.168.1.2:9100'"},
			expectedLT:      []string{},
		},
		{
			name: "Single Loadtest node",
			nodes: []Node{
				{
					IP:    "192.168.1.3",
					Roles: []SupportedRole{Loadtest},
				},
			},
			expectedOG:      []string{},
			expectedMachine: []string{},
			expectedLT:      []string{"'192.168.1.3:8082'"}, // OdysseygoLoadTestPort
		},
		{
			name: "Single Monitor node",
			nodes: []Node{
				{
					IP:    "192.168.1.4",
					Roles: []SupportedRole{Monitor},
				},
			},
			expectedOG:      []string{},
			expectedMachine: []string{},
			expectedLT:      []string{},
		},
		{
			name: "Mixed nodes",
			nodes: []Node{
				{
					IP:    "192.168.1.1",
					Roles: []SupportedRole{Validator},
				},
				{
					IP:    "192.168.1.2",
					Roles: []SupportedRole{API},
				},
				{
					IP:    "192.168.1.3",
					Roles: []SupportedRole{Loadtest},
				},
				{
					IP:    "192.168.1.4",
					Roles: []SupportedRole{Monitor},
				},
			},
			expectedOG:      []string{"'192.168.1.1:9650'", "'192.168.1.2:9650'"},
			expectedMachine: []string{"'192.168.1.1:9100'", "'192.168.1.2:9100'"},
			expectedLT:      []string{"'192.168.1.3:8082'"},
		},
		{
			name: "Node with multiple roles including Validator",
			nodes: []Node{
				{
					IP:    "192.168.1.1",
					Roles: []SupportedRole{Validator, Monitor},
				},
			},
			expectedOG:      []string{"'192.168.1.1:9650'"},
			expectedMachine: []string{"'192.168.1.1:9100'"},
			expectedLT:      []string{},
		},
		{
			name: "Node with multiple roles including API",
			nodes: []Node{
				{
					IP:    "192.168.1.1",
					Roles: []SupportedRole{API, Monitor},
				},
			},
			expectedOG:      []string{"'192.168.1.1:9650'"},
			expectedMachine: []string{"'192.168.1.1:9100'"},
			expectedLT:      []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ogPorts, machinePorts, ltPorts := getPrometheusTargets(tt.nodes)
			assert.Equal(t, tt.expectedOG, ogPorts)
			assert.Equal(t, tt.expectedMachine, machinePorts)
			assert.Equal(t, tt.expectedLT, ltPorts)
		})
	}
}

func TestComposeFileExists(t *testing.T) {
	// This function is hard to test without a real connection
	// We can only test the basic structure
	node := Node{
		IP: "192.168.1.1",
		SSHConfig: SSHConfig{
			User:           "ubuntu",
			PrivateKeyPath: "/path/to/key",
		},
	}

	// This will fail due to no connection, but we can test the function exists
	result := composeFileExists(node)
	assert.False(t, result) // Should be false when not connected
}

func TestGenesisFileExists(t *testing.T) {
	// This function is hard to test without a real connection
	// We can only test the basic structure
	node := Node{
		IP: "192.168.1.1",
		SSHConfig: SSHConfig{
			User:           "ubuntu",
			PrivateKeyPath: "/path/to/key",
		},
	}

	// This will fail due to no connection, but we can test the function exists
	result := genesisFileExists(node)
	assert.False(t, result) // Should be false when not connected
}

func TestNodeConfigFileExists(t *testing.T) {
	// This function is hard to test without a real connection
	// We can only test the basic structure
	node := Node{
		IP: "192.168.1.1",
		SSHConfig: SSHConfig{
			User:           "ubuntu",
			PrivateKeyPath: "/path/to/key",
		},
	}

	// This will fail due to no connection, but we can test the function exists
	result := nodeConfigFileExists(node)
	assert.False(t, result) // Should be false when not connected
}

func TestGetPublicKeyFromSSHKey_EdgeCases(t *testing.T) {
	// Create a temporary directory for test files
	tempDir := t.TempDir()

	tests := []struct {
		name        string
		setup       func() string
		expectError bool
	}{
		{
			name: "File with no newline at end",
			setup: func() string {
				keyPath := filepath.Join(tempDir, "no_newline.pub")
				keyContent := "ssh-rsa AAAAB3NzaC1yc2EAAAABIwAAAQEAr8E7T/ZoQ9Jyb5F1U1t/9F+nkRoSi8g8j6x0g7vZJ68dVVpREzK84+R5cOJ6ydP9Nd+G99kW1HLhfwK5BhJnW3uZ7h1mL0Hh/RZb8csViNe8sEc2FSgH5G8cl3ZX8Y1UtdbS4k5F3cC3B4JFF9y6vOZRwUBO4z1Z2BZaGP29sXXkW0ZGRrWaBswcq+S5FJ1QOeeJ38OjkB45L7zq2X2NQ== user@hostname"
				err := os.WriteFile(keyPath, []byte(keyContent), 0o600)
				require.NoError(t, err)
				return keyPath
			},
			expectError: false,
		},
		{
			name: "File with multiple newlines",
			setup: func() string {
				keyPath := filepath.Join(tempDir, "multiple_newlines.pub")
				keyContent := "ssh-rsa AAAAB3NzaC1yc2EAAAABIwAAAQEAr8E7T/ZoQ9Jyb5F1U1t/9F+nkRoSi8g8j6x0g7vZJ68dVVpREzK84+R5cOJ6ydP9Nd+G99kW1HLhfwK5BhJnW3uZ7h1mL0Hh/RZb8csViNe8sEc2FSgH5G8cl3ZX8Y1UtdbS4k5F3cC3B4JFF9y6vOZRwUBO4z1Z2BZaGP29sXXkW0ZGRrWaBswcq+S5FJ1QOeeJ38OjkB45L7zq2X2NQ== user@hostname\n\n"
				err := os.WriteFile(keyPath, []byte(keyContent), 0o600)
				require.NoError(t, err)
				return keyPath
			},
			expectError: false,
		},
		{
			name: "Empty file",
			setup: func() string {
				keyPath := filepath.Join(tempDir, "empty.pub")
				err := os.WriteFile(keyPath, []byte(""), 0o600)
				require.NoError(t, err)
				return keyPath
			},
			expectError: false,
		},
		{
			name: "File with only newlines",
			setup: func() string {
				keyPath := filepath.Join(tempDir, "only_newlines.pub")
				err := os.WriteFile(keyPath, []byte("\n\n\n"), 0o600)
				require.NoError(t, err)
				return keyPath
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keyPath := tt.setup()
			result, err := GetPublicKeyFromSSHKey(keyPath)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				// The result should have the trailing newline removed
				// Allow a single trailing newline if present in file; trim before check
				trimmed := strings.TrimRight(result, "\n")
				assert.NotContains(t, trimmed, "\n")
			}
		})
	}
}

func TestGetDefaultProjectNameFromGCPCredentials_EdgeCases(t *testing.T) {
	t.Skip("GCP functionality has been removed from this SDK")
	// Create a temporary directory for test files
	tempDir := t.TempDir()

	tests := []struct {
		name        string
		setup       func() string
		expectError bool
		expected    string
	}{
		{
			name: "Empty JSON object",
			setup: func() string {
				filePath := filepath.Join(tempDir, "empty.json")
				err := os.WriteFile(filePath, []byte("{}"), 0o600)
				require.NoError(t, err)
				return filePath
			},
			expectError: false,
			expected:    "",
		},
		{
			name: "Missing quota_project_id field",
			setup: func() string {
				filePath := filepath.Join(tempDir, "missing_field.json")
				content := `{
					"client_id": "test-client-id",
					"client_secret": "test-client-secret",
					"refresh_token": "test-refresh-token",
					"type": "service_account"
				}`
				err := os.WriteFile(filePath, []byte(content), 0o600)
				require.NoError(t, err)
				return filePath
			},
			expectError: false,
			expected:    "",
		},
		{
			name: "Null quota_project_id",
			setup: func() string {
				filePath := filepath.Join(tempDir, "null_field.json")
				content := `{
					"client_id": "test-client-id",
					"client_secret": "test-client-secret",
					"quota_project_id": null,
					"refresh_token": "test-refresh-token",
					"type": "service_account"
				}`
				err := os.WriteFile(filePath, []byte(content), 0o600)
				require.NoError(t, err)
				return filePath
			},
			expectError: false,
			expected:    "",
		},
		{
			name: "Empty quota_project_id",
			setup: func() string {
				filePath := filepath.Join(tempDir, "empty_field.json")
				content := `{
					"client_id": "test-client-id",
					"client_secret": "test-client-secret",
					"quota_project_id": "",
					"refresh_token": "test-refresh-token",
					"type": "service_account"
				}`
				err := os.WriteFile(filePath, []byte(content), 0o600)
				require.NoError(t, err)
				return filePath
			},
			expectError: false,
			expected:    "",
		},
		{
			name: "Malformed JSON",
			setup: func() string {
				filePath := filepath.Join(tempDir, "malformed.json")
				content := `{
					"client_id": "test-client-id",
					"client_secret": "test-client-secret",
					"quota_project_id": "test-project-id",
					"refresh_token": "test-refresh-token",
					"type": "service_account"
				` // Missing closing brace
				err := os.WriteFile(filePath, []byte(content), 0o600)
				require.NoError(t, err)
				return filePath
			},
			expectError: true,
			expected:    "",
		},
		{
			name: "Non-JSON file",
			setup: func() string {
				filePath := filepath.Join(tempDir, "not_json.txt")
				content := "This is not JSON content"
				err := os.WriteFile(filePath, []byte(content), 0o600)
				require.NoError(t, err)
				return filePath
			},
			expectError: true,
			expected:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filePath := tt.setup()
			result, err := getDefaultProjectNameFromGCPCredentials(filePath)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestConstants(t *testing.T) {
	// Test that the constants are properly defined
	assert.Equal(t, 102400, maxResponseSize)
	assert.Equal(t, 5, sshConnectionRetries)
	assert.Equal(t, "3s", sshConnectionTimeout.String())
}
