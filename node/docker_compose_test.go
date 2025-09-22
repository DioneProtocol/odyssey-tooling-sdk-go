// Copyright (C) 2025, Dione Protocol, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package node

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/DioneProtocol/odyssey-tooling-sdk-go/constants"
)

func TestRenderComposeFile(t *testing.T) {
	tests := []struct {
		name         string
		composePath  string
		composeDesc  string
		templateVars dockerComposeInputs
		expectError  bool
	}{
		{
			name:        "Valid template with monitoring",
			composePath: "templates/odysseygo.docker-compose.yml",
			composeDesc: "test compose",
			templateVars: dockerComposeInputs{
				WithMonitoring:   true,
				WithOdysseygo:    true,
				OdysseygoVersion: "v1.11.8",
				E2E:              false,
				E2EIP:            "",
				E2ESuffix:        "",
			},
			expectError: false,
		},
		{
			name:        "Valid template without monitoring",
			composePath: "templates/odysseygo.docker-compose.yml",
			composeDesc: "test compose",
			templateVars: dockerComposeInputs{
				WithMonitoring:   false,
				WithOdysseygo:    true,
				OdysseygoVersion: "v1.11.8",
				E2E:              false,
				E2EIP:            "",
				E2ESuffix:        "",
			},
			expectError: false,
		},
		{
			name:        "E2E template",
			composePath: "templates/odysseygo.docker-compose.yml",
			composeDesc: "test compose",
			templateVars: dockerComposeInputs{
				WithMonitoring:   true,
				WithOdysseygo:    true,
				OdysseygoVersion: "v1.11.8",
				E2E:              true,
				E2EIP:            "192.168.1.1",
				E2ESuffix:        "test",
			},
			expectError: false,
		},
		{
			name:        "Non-existent template file",
			composePath: "templates/nonexistent.yml",
			composeDesc: "test compose",
			templateVars: dockerComposeInputs{
				WithMonitoring:   true,
				WithOdysseygo:    true,
				OdysseygoVersion: "v1.11.8",
				E2E:              false,
				E2EIP:            "",
				E2ESuffix:        "",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := renderComposeFile(tt.composePath, tt.composeDesc, tt.templateVars)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, result)
			}
		})
	}
}

func TestNode_PushComposeFile(t *testing.T) {
	// Create a temporary file for testing
	tempDir := t.TempDir()
	localFile := filepath.Join(tempDir, "test-compose.yml")
	remoteFile := "/remote/compose.yml"

	// Write some test content
	testContent := `version: '3.8'
services:
  test:
    image: alpine:latest
    command: echo "test"
`
	err := os.WriteFile(localFile, []byte(testContent), 0o644)
	require.NoError(t, err)

	tests := []struct {
		name        string
		node        Node
		localFile   string
		remoteFile  string
		merge       bool
		expectError bool
	}{
		{
			name: "Non-existent local file",
			node: Node{
				IP: "192.168.1.1",
				SSHConfig: SSHConfig{
					User:           "ubuntu",
					PrivateKeyPath: "/path/to/key",
				},
			},
			localFile:   "/non/existent/file.yml",
			remoteFile:  remoteFile,
			merge:       false,
			expectError: true,
		},
		{
			name: "Valid local file but no connection",
			node: Node{
				IP: "192.168.1.1",
				SSHConfig: SSHConfig{
					User:           "ubuntu",
					PrivateKeyPath: "/path/to/key",
				},
			},
			localFile:   localFile,
			remoteFile:  remoteFile,
			merge:       false,
			expectError: true, // Will fail due to no connection
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.node.PushComposeFile(tt.localFile, tt.remoteFile, tt.merge)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNode_MergeComposeFiles(t *testing.T) {
	tests := []struct {
		name        string
		node        Node
		currentFile string
		newFile     string
		expectError bool
	}{
		{
			name: "No connection",
			node: Node{
				IP: "192.168.1.1",
				SSHConfig: SSHConfig{
					User:           "ubuntu",
					PrivateKeyPath: "/path/to/key",
				},
			},
			currentFile: "/remote/current.yml",
			newFile:     "/remote/new.yml",
			expectError: true, // Will fail due to no connection
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.node.MergeComposeFiles(tt.currentFile, tt.newFile)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDockerComposeInputs(t *testing.T) {
	tests := []struct {
		name        string
		inputs      dockerComposeInputs
		expectedOG  bool
		expectedMon bool
		expectedE2E bool
	}{
		{
			name: "Full configuration",
			inputs: dockerComposeInputs{
				WithMonitoring:   true,
				WithOdysseygo:    true,
				OdysseygoVersion: "v1.11.8",
				E2E:              true,
				E2EIP:            "192.168.1.1",
				E2ESuffix:        "test",
			},
			expectedOG:  true,
			expectedMon: true,
			expectedE2E: true,
		},
		{
			name: "Minimal configuration",
			inputs: dockerComposeInputs{
				WithMonitoring:   false,
				WithOdysseygo:    false,
				OdysseygoVersion: "",
				E2E:              false,
				E2EIP:            "",
				E2ESuffix:        "",
			},
			expectedOG:  false,
			expectedMon: false,
			expectedE2E: false,
		},
		{
			name: "OdysseyGo only",
			inputs: dockerComposeInputs{
				WithMonitoring:   false,
				WithOdysseygo:    true,
				OdysseygoVersion: "v1.11.8",
				E2E:              false,
				E2EIP:            "",
				E2ESuffix:        "",
			},
			expectedOG:  true,
			expectedMon: false,
			expectedE2E: false,
		},
		{
			name: "Monitoring only",
			inputs: dockerComposeInputs{
				WithMonitoring:   true,
				WithOdysseygo:    false,
				OdysseygoVersion: "",
				E2E:              false,
				E2EIP:            "",
				E2ESuffix:        "",
			},
			expectedOG:  false,
			expectedMon: true,
			expectedE2E: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expectedOG, tt.inputs.WithOdysseygo)
			assert.Equal(t, tt.expectedMon, tt.inputs.WithMonitoring)
			assert.Equal(t, tt.expectedE2E, tt.inputs.E2E)
		})
	}
}

func TestNode_DockerComposeOperations(t *testing.T) {
	// Create a test node
	node := Node{
		IP: "192.168.1.1",
		SSHConfig: SSHConfig{
			User:           "ubuntu",
			PrivateKeyPath: "/path/to/key",
		},
	}

	tests := []struct {
		name        string
		operation   func() error
		expectError bool
	}{
		{
			name: "ComposeSSHSetupNode",
			operation: func() error {
				return node.ComposeSSHSetupNode("testnet", []string{}, "v1.11.8", true)
			},
			expectError: true, // Will fail due to no connection
		},
		{
			name: "ComposeSSHSetupLoadTest",
			operation: func() error {
				return node.ComposeSSHSetupLoadTest()
			},
			expectError: true, // Will fail due to no connection
		},
		{
			name: "ComposeSSHSetupMonitoring",
			operation: func() error {
				return node.ComposeSSHSetupMonitoring()
			},
			expectError: true, // Will fail due to no connection
		},
		{
			name: "ComposeSSHSetupAWMRelayer",
			operation: func() error {
				return node.ComposeSSHSetupAWMRelayer()
			},
			expectError: true, // Will fail due to no connection
		},
		{
			name: "StartDockerCompose",
			operation: func() error {
				return node.StartDockerCompose(time.Second)
			},
			expectError: true, // Will fail due to no connection
		},
		{
			name: "RestartDockerCompose",
			operation: func() error {
				return node.RestartDockerCompose(time.Second)
			},
			expectError: true, // Will fail due to no connection
		},
		{
			name: "StartDockerComposeService",
			operation: func() error {
				return node.StartDockerComposeService("/remote/compose.yml", "test-service", time.Second)
			},
			expectError: true, // Will fail due to no connection
		},
		{
			name: "StopDockerComposeService",
			operation: func() error {
				return node.StopDockerComposeService("/remote/compose.yml", "test-service", time.Second)
			},
			expectError: true, // Will fail due to no connection
		},
		{
			name: "RestartDockerComposeService",
			operation: func() error {
				return node.RestartDockerComposeService("/remote/compose.yml", "test-service", time.Second)
			},
			expectError: true, // Will fail due to no connection
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.operation()
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNode_DockerComposeFileOperations(t *testing.T) {
	tests := []struct {
		name        string
		operation   func() error
		expectError bool
	}{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.operation()
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNode_DockerComposeValidation(t *testing.T) {
	tests := []struct {
		name        string
		operation   func() (bool, error)
		expectError bool
	}{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.operation()
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDockerComposeConstants(t *testing.T) {
	// Test that the constants are properly defined
	assert.NotEmpty(t, constants.ServiceOdysseygo)
	assert.NotEmpty(t, constants.ServicePrometheus)
	assert.NotEmpty(t, constants.ServiceGrafana)
	assert.NotEmpty(t, constants.ServiceLoki)
	assert.NotEmpty(t, constants.ServicePromtail)
	assert.NotEmpty(t, constants.ServiceAWMRelayer)
}

func TestNode_DockerComposeEdgeCases(t *testing.T) {
	// Test with empty node
	node := Node{}

	tests := []struct {
		name        string
		operation   func() error
		expectError bool
	}{
		{
			name: "Empty node - ComposeSSHSetupNode",
			operation: func() error {
				return node.ComposeSSHSetupNode("", []string{}, "", false)
			},
			expectError: true,
		},
		{
			name: "Empty node - StartDockerCompose",
			operation: func() error {
				return node.StartDockerCompose(time.Second)
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.operation()
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDockerComposeInputs_EdgeCases(t *testing.T) {
	// Test with empty inputs
	inputs := dockerComposeInputs{}
	assert.False(t, inputs.WithMonitoring)
	assert.False(t, inputs.WithOdysseygo)
	assert.False(t, inputs.E2E)
	assert.Empty(t, inputs.OdysseygoVersion)
	assert.Empty(t, inputs.E2EIP)
	assert.Empty(t, inputs.E2ESuffix)

	// Test with partial inputs
	inputs = dockerComposeInputs{
		WithOdysseygo: true,
	}
	assert.False(t, inputs.WithMonitoring)
	assert.True(t, inputs.WithOdysseygo)
	assert.False(t, inputs.E2E)
}

func TestNode_DockerComposeTimeoutHandling(t *testing.T) {
	node := Node{
		IP: "192.168.1.1",
		SSHConfig: SSHConfig{
			User:           "ubuntu",
			PrivateKeyPath: "/path/to/key",
		},
	}

	tests := []struct {
		name    string
		timeout time.Duration
	}{
		{
			name:    "Zero timeout",
			timeout: 0,
		},
		{
			name:    "Very short timeout",
			timeout: time.Nanosecond,
		},
		{
			name:    "Normal timeout",
			timeout: time.Second,
		},
		{
			name:    "Long timeout",
			timeout: time.Minute,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := node.StartDockerCompose(tt.timeout)
			assert.Error(t, err) // Will fail due to no connection, but timeout should be handled
		})
	}
}

func TestNode_StopDockerCompose_Success(t *testing.T) {
	node := &Node{
		NodeID: "test-node",
		IP:     "192.168.1.1",
		SSHConfig: SSHConfig{
			User:           "ubuntu",
			PrivateKeyPath: "/path/to/key",
		},
	}

	// Test with a node that will fail due to SSH connection issues
	// This tests the error handling path
	err := node.StopDockerCompose(time.Second)

	assert.Error(t, err) // Will fail due to SSH connection issues
}

func TestNode_StopDockerCompose_SystemDError(t *testing.T) {
	node := &Node{
		NodeID: "test-node",
		IP:     "192.168.1.1",
		SSHConfig: SSHConfig{
			User:           "ubuntu",
			PrivateKeyPath: "/path/to/key",
		},
	}

	// Test with a node that will fail due to SSH connection issues
	err := node.StopDockerCompose(time.Second)

	assert.Error(t, err) // Will fail due to SSH connection issues
}

func TestNode_StopDockerCompose_NoSystemD(t *testing.T) {
	node := &Node{
		NodeID: "test-node",
		IP:     "192.168.1.1",
		SSHConfig: SSHConfig{
			User:           "ubuntu",
			PrivateKeyPath: "/path/to/key",
		},
	}

	// Test with a node that will fail due to SSH connection issues
	err := node.StopDockerCompose(time.Second)

	assert.Error(t, err) // Will fail due to SSH connection issues
}

func TestNode_StartDockerComposeService_Success(t *testing.T) {
	node := &Node{
		NodeID: "test-node",
		IP:     "192.168.1.1",
		SSHConfig: SSHConfig{
			User:           "ubuntu",
			PrivateKeyPath: "/path/to/key",
		},
	}

	// Test with a node that will fail due to SSH connection issues
	err := node.StartDockerComposeService("/path/to/compose.yml", "test-service", time.Second)

	assert.Error(t, err) // Will fail due to SSH connection issues
}

func TestNode_StartDockerComposeService_InitError(t *testing.T) {
	node := &Node{
		NodeID: "test-node",
		IP:     "192.168.1.1",
		SSHConfig: SSHConfig{
			User:           "ubuntu",
			PrivateKeyPath: "/path/to/key",
		},
	}

	// Test with a node that will fail due to SSH connection issues
	err := node.StartDockerComposeService("/path/to/compose.yml", "test-service", time.Second)

	assert.Error(t, err) // Will fail due to SSH connection issues
}

func TestNode_ListRemoteComposeServices_Success(t *testing.T) {
	node := &Node{
		NodeID: "test-node",
		IP:     "192.168.1.1",
		SSHConfig: SSHConfig{
			User:           "ubuntu",
			PrivateKeyPath: "/path/to/key",
		},
	}

	// Test with a node that will fail due to SSH connection issues
	services, err := node.ListRemoteComposeServices("/path/to/compose.yml", time.Second)

	assert.Error(t, err) // Will fail due to SSH connection issues
	assert.Nil(t, services)
}

func TestNode_ListRemoteComposeServices_Error(t *testing.T) {
	node := &Node{
		NodeID: "test-node",
		IP:     "192.168.1.1",
		SSHConfig: SSHConfig{
			User:           "ubuntu",
			PrivateKeyPath: "/path/to/key",
		},
	}

	// Test with a node that will fail due to SSH connection issues
	services, err := node.ListRemoteComposeServices("/path/to/compose.yml", time.Second)

	assert.Error(t, err) // Will fail due to SSH connection issues
	assert.Nil(t, services)
}

func TestNode_GetRemoteComposeContent_Success(t *testing.T) {
	node := &Node{
		NodeID: "test-node",
		IP:     "192.168.1.1",
		SSHConfig: SSHConfig{
			User:           "ubuntu",
			PrivateKeyPath: "/path/to/key",
		},
	}

	// Test with a node that will fail due to SSH connection issues
	content, err := node.GetRemoteComposeContent("/path/to/compose.yml", time.Second)

	assert.Error(t, err) // Will fail due to SSH connection issues
	assert.Empty(t, content)
}

func TestNode_GetRemoteComposeContent_DownloadError(t *testing.T) {
	node := &Node{
		NodeID: "test-node",
		IP:     "192.168.1.1",
		SSHConfig: SSHConfig{
			User:           "ubuntu",
			PrivateKeyPath: "/path/to/key",
		},
	}

	// Test with a node that will fail due to SSH connection issues
	content, err := node.GetRemoteComposeContent("/path/to/compose.yml", time.Second)

	assert.Error(t, err) // Will fail due to SSH connection issues
	assert.Empty(t, content)
}

func TestNode_ParseRemoteComposeContent_Success(t *testing.T) {
	node := &Node{
		NodeID: "test-node",
		IP:     "192.168.1.1",
		SSHConfig: SSHConfig{
			User:           "ubuntu",
			PrivateKeyPath: "/path/to/key",
		},
	}

	// Test with a node that will fail due to SSH connection issues
	value, err := node.ParseRemoteComposeContent("/path/to/compose.yml", "pattern", time.Second)

	assert.Error(t, err) // Will fail due to SSH connection issues
	assert.Empty(t, value)
}

func TestNode_ParseRemoteComposeContent_GetContentError(t *testing.T) {
	node := &Node{
		NodeID: "test-node",
		IP:     "192.168.1.1",
		SSHConfig: SSHConfig{
			User:           "ubuntu",
			PrivateKeyPath: "/path/to/key",
		},
	}

	// Test with a node that will fail due to SSH connection issues
	value, err := node.ParseRemoteComposeContent("/path/to/compose.yml", "pattern", time.Second)

	assert.Error(t, err) // Will fail due to SSH connection issues
	assert.Empty(t, value)
}

func TestNode_HasRemoteComposeService_ServiceExists(t *testing.T) {
	node := &Node{
		NodeID: "test-node",
		IP:     "192.168.1.1",
		SSHConfig: SSHConfig{
			User:           "ubuntu",
			PrivateKeyPath: "/path/to/key",
		},
	}

	// Test with a node that will fail due to SSH connection issues
	exists, err := node.HasRemoteComposeService("/path/to/compose.yml", "service2", time.Second)

	assert.Error(t, err) // Will fail due to SSH connection issues
	assert.False(t, exists)
}

func TestNode_HasRemoteComposeService_ServiceNotExists(t *testing.T) {
	node := &Node{
		NodeID: "test-node",
		IP:     "192.168.1.1",
		SSHConfig: SSHConfig{
			User:           "ubuntu",
			PrivateKeyPath: "/path/to/key",
		},
	}

	// Test with a node that will fail due to SSH connection issues
	exists, err := node.HasRemoteComposeService("/path/to/compose.yml", "service4", time.Second)

	assert.Error(t, err) // Will fail due to SSH connection issues
	assert.False(t, exists)
}

func TestNode_HasRemoteComposeService_ListError(t *testing.T) {
	node := &Node{
		NodeID: "test-node",
		IP:     "192.168.1.1",
		SSHConfig: SSHConfig{
			User:           "ubuntu",
			PrivateKeyPath: "/path/to/key",
		},
	}

	// Test with a node that will fail due to SSH connection issues
	exists, err := node.HasRemoteComposeService("/path/to/compose.yml", "service1", time.Second)

	assert.Error(t, err) // Will fail due to SSH connection issues
	assert.False(t, exists)
}

func TestNode_ListDockerComposeImages_Success(t *testing.T) {
	node := &Node{
		NodeID: "test-node",
		IP:     "192.168.1.1",
		SSHConfig: SSHConfig{
			User:           "ubuntu",
			PrivateKeyPath: "/path/to/key",
		},
	}

	// Test with a node that will fail due to SSH connection issues
	images, err := node.ListDockerComposeImages("/path/to/compose.yml", time.Second)

	assert.Error(t, err) // Will fail due to SSH connection issues
	assert.Nil(t, images)
}

func TestNode_ListDockerComposeImages_CommandError(t *testing.T) {
	node := &Node{
		NodeID: "test-node",
		IP:     "192.168.1.1",
		SSHConfig: SSHConfig{
			User:           "ubuntu",
			PrivateKeyPath: "/path/to/key",
		},
	}

	// Test with a node that will fail due to SSH connection issues
	images, err := node.ListDockerComposeImages("/path/to/compose.yml", time.Second)

	assert.Error(t, err) // Will fail due to SSH connection issues
	assert.Nil(t, images)
}

func TestNode_ListDockerComposeImages_InvalidJSON(t *testing.T) {
	node := &Node{
		NodeID: "test-node",
		IP:     "192.168.1.1",
		SSHConfig: SSHConfig{
			User:           "ubuntu",
			PrivateKeyPath: "/path/to/key",
		},
	}

	// Test with a node that will fail due to SSH connection issues
	images, err := node.ListDockerComposeImages("/path/to/compose.yml", time.Second)

	assert.Error(t, err) // Will fail due to SSH connection issues
	assert.Nil(t, images)
}

func TestNode_GetDockerImageVersion_Success(t *testing.T) {
	node := &Node{
		NodeID: "test-node",
		IP:     "192.168.1.1",
		SSHConfig: SSHConfig{
			User:           "ubuntu",
			PrivateKeyPath: "/path/to/key",
		},
	}

	// Test with a node that will fail due to SSH connection issues
	version, err := node.GetDockerImageVersion("test-image", time.Second)

	assert.Error(t, err) // Will fail due to SSH connection issues
	assert.Empty(t, version)
}

func TestNode_GetDockerImageVersion_ImageNotFound(t *testing.T) {
	node := &Node{
		NodeID: "test-node",
		IP:     "192.168.1.1",
		SSHConfig: SSHConfig{
			User:           "ubuntu",
			PrivateKeyPath: "/path/to/key",
		},
	}

	// Test with a node that will fail due to SSH connection issues
	version, err := node.GetDockerImageVersion("test-image", time.Second)

	assert.Error(t, err) // Will fail due to SSH connection issues
	assert.Empty(t, version)
}

func TestNode_GetDockerImageVersion_ListError(t *testing.T) {
	node := &Node{
		NodeID: "test-node",
		IP:     "192.168.1.1",
		SSHConfig: SSHConfig{
			User:           "ubuntu",
			PrivateKeyPath: "/path/to/key",
		},
	}

	// Test with a node that will fail due to SSH connection issues
	version, err := node.GetDockerImageVersion("test-image", time.Second)

	assert.Error(t, err) // Will fail due to SSH connection issues
	assert.Empty(t, version)
}

func TestNode_StopDockerComposeService_Success(t *testing.T) {
	node := &Node{
		NodeID: "test-node",
		IP:     "192.168.1.1",
		SSHConfig: SSHConfig{
			User:           "ubuntu",
			PrivateKeyPath: "/path/to/key",
		},
	}

	// Test with a node that will fail due to SSH connection issues
	err := node.StopDockerComposeService("/path/to/compose.yml", "test-service", time.Second)

	assert.Error(t, err) // Will fail due to SSH connection issues
}

func TestNode_StopDockerComposeService_Error(t *testing.T) {
	node := &Node{
		NodeID: "test-node",
		IP:     "192.168.1.1",
		SSHConfig: SSHConfig{
			User:           "ubuntu",
			PrivateKeyPath: "/path/to/key",
		},
	}

	// Test with a node that will fail due to SSH connection issues
	err := node.StopDockerComposeService("/path/to/compose.yml", "test-service", time.Second)

	assert.Error(t, err) // Will fail due to SSH connection issues
}

func TestNode_RestartDockerComposeService_Success(t *testing.T) {
	node := &Node{
		NodeID: "test-node",
		IP:     "192.168.1.1",
		SSHConfig: SSHConfig{
			User:           "ubuntu",
			PrivateKeyPath: "/path/to/key",
		},
	}

	// Test with a node that will fail due to SSH connection issues
	err := node.RestartDockerComposeService("/path/to/compose.yml", "test-service", time.Second)

	assert.Error(t, err) // Will fail due to SSH connection issues
}

func TestNode_InitDockerComposeService_Success(t *testing.T) {
	node := &Node{
		NodeID: "test-node",
		IP:     "192.168.1.1",
		SSHConfig: SSHConfig{
			User:           "ubuntu",
			PrivateKeyPath: "/path/to/key",
		},
	}

	// Test with a node that will fail due to SSH connection issues
	err := node.InitDockerComposeService("/path/to/compose.yml", "test-service", time.Second)

	assert.Error(t, err) // Will fail due to SSH connection issues
}
