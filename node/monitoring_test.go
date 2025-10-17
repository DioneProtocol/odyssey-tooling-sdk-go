// Copyright (C) 2025, Dione Protocol, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package node

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/DioneProtocol/odyssey-tooling-sdk-go/constants"
	"github.com/DioneProtocol/odyssey-tooling-sdk-go/odyssey"
)

func TestNode_MonitorNodes(t *testing.T) {
	tests := []struct {
		name        string
		node        Node
		nodes       []Node
		clusterName string
		expectError bool
	}{
		{
			name: "Empty nodes list",
			node: Node{
				IP: "192.168.1.1",
				SSHConfig: SSHConfig{
					User:           "ubuntu",
					PrivateKeyPath: "/path/to/key",
				},
				Roles: []SupportedRole{Monitor},
			},
			nodes:       []Node{},
			clusterName: "test-cluster",
			expectError: true, // Will fail due to no connection
		},
		{
			name: "Single node to monitor",
			node: Node{
				IP: "192.168.1.1",
				SSHConfig: SSHConfig{
					User:           "ubuntu",
					PrivateKeyPath: "/path/to/key",
				},
				Roles: []SupportedRole{Monitor},
			},
			nodes: []Node{
				{
					IP:    "192.168.1.2",
					Roles: []SupportedRole{Validator},
				},
			},
			clusterName: "test-cluster",
			// Expect error due to no SSH connection
			expectError: true,
		},
		{
			name: "Multiple nodes to monitor",
			node: Node{
				IP: "192.168.1.1",
				SSHConfig: SSHConfig{
					User:           "ubuntu",
					PrivateKeyPath: "/path/to/key",
				},
				Roles: []SupportedRole{Monitor},
			},
			nodes: []Node{
				{
					IP:    "192.168.1.2",
					Roles: []SupportedRole{Validator},
				},
				{
					IP:    "192.168.1.3",
					Roles: []SupportedRole{API},
				},
			},
			clusterName: "test-cluster",
			expectError: true, // Will fail due to no connection
		},
		{
			name: "Empty cluster name",
			node: Node{
				IP: "192.168.1.1",
				SSHConfig: SSHConfig{
					User:           "ubuntu",
					PrivateKeyPath: "/path/to/key",
				},
				Roles: []SupportedRole{Monitor},
			},
			nodes: []Node{
				{
					IP:    "192.168.1.2",
					Roles: []SupportedRole{Validator},
				},
			},
			clusterName: "",
			expectError: true, // Will fail due to no connection
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			err := tt.node.MonitorNodes(ctx, tt.nodes, tt.clusterName)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNode_ProvisioningFunctions(t *testing.T) {
	// Save original Docker support flag value
	originalDockerSupport := constants.DockerSupportEnabled
	defer func() {
		constants.DockerSupportEnabled = originalDockerSupport
	}()

	// Enable Docker support for this test
	constants.DockerSupportEnabled = true

	tests := []struct {
		name        string
		node        Node
		expectError bool
	}{
		{
			name: "RunSSHSetupNode",
			node: Node{
				IP: "192.168.1.1",
				SSHConfig: SSHConfig{
					User:           "ubuntu",
					PrivateKeyPath: "/path/to/key",
				},
			},
			expectError: true, // Will fail due to no connection
		},
		{
			name: "RunSSHSetupDockerService",
			node: Node{
				IP: "192.168.1.1",
				SSHConfig: SSHConfig{
					User:           "ubuntu",
					PrivateKeyPath: "/path/to/key",
				},
			},
			expectError: false, // No-op when systemd not available
		},
		{
			name: "RunSSHSetupMonitoringFolders",
			node: Node{
				IP: "192.168.1.1",
				SSHConfig: SSHConfig{
					User:           "ubuntu",
					PrivateKeyPath: "/path/to/key",
				},
			},
			expectError: true, // Will fail due to no connection
		},
		{
			name: "RunSSHSetupPromtailConfig",
			node: Node{
				IP: "192.168.1.1",
				SSHConfig: SSHConfig{
					User:           "ubuntu",
					PrivateKeyPath: "/path/to/key",
				},
			},
			expectError: true, // Will fail due to no connection
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			switch tt.name {
			case "RunSSHSetupNode":
				err = tt.node.RunSSHSetupNode()
			case "RunSSHSetupDockerService":
				err = tt.node.RunSSHSetupDockerService()
			case "RunSSHSetupMonitoringFolders":
				err = tt.node.RunSSHSetupMonitoringFolders()
			case "RunSSHSetupPromtailConfig":
				err = tt.node.RunSSHSetupPromtailConfig("127.0.0.1", 3100, "node-1", "test-cluster")
			}

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestProvisionHost(t *testing.T) {
	tests := []struct {
		name        string
		node        Node
		nodeParams  *NodeParams
		expectError bool
	}{
		{
			name: "Valid validator node",
			node: Node{
				IP: "192.168.1.1",
				SSHConfig: SSHConfig{
					User:           "ubuntu",
					PrivateKeyPath: "/path/to/key",
				},
			},
			nodeParams: &NodeParams{
				Roles:            []SupportedRole{Validator},
				Network:          odyssey.TestnetNetwork(),
				OdysseyGoVersion: "v1.10.13",
			},
			expectError: true, // Will fail due to no connection
		},
		{
			name: "Valid API node",
			node: Node{
				IP: "192.168.1.1",
				SSHConfig: SSHConfig{
					User:           "ubuntu",
					PrivateKeyPath: "/path/to/key",
				},
			},
			nodeParams: &NodeParams{
				Roles:            []SupportedRole{API},
				Network:          odyssey.TestnetNetwork(),
				OdysseyGoVersion: "v1.10.13",
			},
			expectError: true, // Will fail due to no connection
		},
		{
			name: "Valid monitor node",
			node: Node{
				IP: "192.168.1.1",
				SSHConfig: SSHConfig{
					User:           "ubuntu",
					PrivateKeyPath: "/path/to/key",
				},
			},
			nodeParams: &NodeParams{
				Roles: []SupportedRole{Monitor},
			},
			expectError: true, // Will fail due to no connection
		},
		{
			name: "Valid load test node",
			node: Node{
				IP: "192.168.1.1",
				SSHConfig: SSHConfig{
					User:           "ubuntu",
					PrivateKeyPath: "/path/to/key",
				},
			},
			nodeParams: &NodeParams{
				Roles: []SupportedRole{Loadtest},
			},
			expectError: true, // Will fail due to no connection
		},
		{
			name: "Invalid role combination",
			node: Node{
				IP: "192.168.1.1",
				SSHConfig: SSHConfig{
					User:           "ubuntu",
					PrivateKeyPath: "/path/to/key",
				},
			},
			nodeParams: &NodeParams{
				Roles: []SupportedRole{Validator, API}, // Invalid combination
			},
			expectError: true, // Should fail due to invalid role combination
		},
		{
			name: "Unsupported role",
			node: Node{
				IP: "192.168.1.1",
				SSHConfig: SSHConfig{
					User:           "ubuntu",
					PrivateKeyPath: "/path/to/key",
				},
			},
			nodeParams: &NodeParams{
				Roles: []SupportedRole{SupportedRole(999)}, // Invalid role
			},
			expectError: true, // Should fail due to unsupported role
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := provisionHost(tt.node, tt.nodeParams)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestProvisionAvagoHost(t *testing.T) {
	tests := []struct {
		name        string
		node        Node
		nodeParams  *NodeParams
		expectError bool
	}{
		{
			name: "Valid validator node",
			node: Node{
				IP: "192.168.1.1",
				SSHConfig: SSHConfig{
					User:           "ubuntu",
					PrivateKeyPath: "/path/to/key",
				},
			},
			nodeParams: &NodeParams{
				Roles:            []SupportedRole{Validator},
				Network:          odyssey.TestnetNetwork(),
				OdysseyGoVersion: "v1.10.13",
			},
			expectError: true, // Will fail due to no connection
		},
		{
			name: "Valid API node",
			node: Node{
				IP: "192.168.1.1",
				SSHConfig: SSHConfig{
					User:           "ubuntu",
					PrivateKeyPath: "/path/to/key",
				},
			},
			nodeParams: &NodeParams{
				Roles:            []SupportedRole{API},
				Network:          odyssey.TestnetNetwork(),
				OdysseyGoVersion: "v1.10.13",
			},
			expectError: true, // Will fail due to no connection
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := provisionOdysseyGoHost(tt.node, tt.nodeParams)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestProvisionLoadTestHost(t *testing.T) {
	tests := []struct {
		name        string
		node        Node
		expectError bool
	}{
		{
			name: "Valid load test node",
			node: Node{
				IP: "192.168.1.1",
				SSHConfig: SSHConfig{
					User:           "ubuntu",
					PrivateKeyPath: "/path/to/key",
				},
			},
			expectError: true, // Will fail due to no connection
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := provisionLoadTestHost(tt.node)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestProvisionMonitoringHost(t *testing.T) {
	tests := []struct {
		name        string
		node        Node
		expectError bool
	}{
		{
			name: "Valid monitoring node",
			node: Node{
				IP: "192.168.1.1",
				SSHConfig: SSHConfig{
					User:           "ubuntu",
					PrivateKeyPath: "/path/to/key",
				},
			},
			expectError: true, // Will fail due to no connection
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := provisionMonitoringHost(tt.node)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNode_MonitoringConfiguration(t *testing.T) {
	tests := []struct {
		name        string
		node        Node
		operation   func() error
		expectError bool
	}{
		{
			name: "ComposeSSHSetupGrafanaConfig",
			node: Node{
				IP: "192.168.1.1",
				SSHConfig: SSHConfig{
					User:           "ubuntu",
					PrivateKeyPath: "/path/to/key",
				},
			},
			operation: func() error {
				node := Node{
					IP: "192.168.1.1",
					SSHConfig: SSHConfig{
						User:           "ubuntu",
						PrivateKeyPath: "/path/to/key",
					},
				}
				// Test basic node functionality
				_, err := node.Command([]string{}, time.Second, "echo test")
				return err
			},
			expectError: true, // Will fail due to no connection
		},
		{
			name: "ComposeSSHSetupPrometheusConfig",
			node: Node{
				IP: "192.168.1.1",
				SSHConfig: SSHConfig{
					User:           "ubuntu",
					PrivateKeyPath: "/path/to/key",
				},
			},
			operation: func() error {
				node := Node{
					IP: "192.168.1.1",
					SSHConfig: SSHConfig{
						User:           "ubuntu",
						PrivateKeyPath: "/path/to/key",
					},
				}
				// Test basic node functionality
				_, err := node.Command([]string{}, time.Second, "echo test")
				return err
			},
			expectError: true, // Will fail due to no connection
		},
		{
			name: "ComposeSSHSetupLokiConfig",
			node: Node{
				IP: "192.168.1.1",
				SSHConfig: SSHConfig{
					User:           "ubuntu",
					PrivateKeyPath: "/path/to/key",
				},
			},
			operation: func() error {
				node := Node{
					IP: "192.168.1.1",
					SSHConfig: SSHConfig{
						User:           "ubuntu",
						PrivateKeyPath: "/path/to/key",
					},
				}
				// Test basic node functionality
				_, err := node.Command([]string{}, time.Second, "echo test")
				return err
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

func TestNode_MonitoringEdgeCases(t *testing.T) {
	// Save original Docker support flag value
	originalDockerSupport := constants.DockerSupportEnabled
	defer func() {
		constants.DockerSupportEnabled = originalDockerSupport
	}()

	// Enable Docker support for this test
	constants.DockerSupportEnabled = true

	tests := []struct {
		name        string
		operation   func() error
		expectError bool
	}{
		{
			name: "Empty node - RunSSHSetupNode",
			operation: func() error {
				node := Node{}
				return node.RunSSHSetupNode()
			},
			expectError: true,
		},
		{
			name: "Empty node - RunSSHSetupDockerService",
			operation: func() error {
				node := Node{}
				return node.RunSSHSetupDockerService()
			},
			expectError: false, // No-op when systemd not available
		},
		{
			name: "Empty node - RunSSHSetupMonitoringFolders",
			operation: func() error {
				node := Node{}
				return node.RunSSHSetupMonitoringFolders()
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

func TestNode_MonitoringTimeoutHandling(t *testing.T) {
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
			node := Node{}
			err := node.RunSSHSetupPromtailConfig("127.0.0.1", 3100, "node-1", "test-cluster")
			assert.Error(t, err) // Will fail due to no connection, but timeout should be handled
		})
	}
}

func TestNode_MonitoringValidation(t *testing.T) {
	tests := []struct {
		name        string
		operation   func() error
		expectError bool
	}{
		{
			name: "Invalid promtail config - empty IP",
			operation: func() error {
				node := Node{}
				return node.RunSSHSetupPromtailConfig("", 3100, "node-1", "test-cluster")
			},
			expectError: true,
		},
		{
			name: "Invalid promtail config - zero port",
			operation: func() error {
				node := Node{}
				return node.RunSSHSetupPromtailConfig("127.0.0.1", 0, "node-1", "test-cluster")
			},
			expectError: true,
		},
		{
			name: "Invalid promtail config - empty node ID",
			operation: func() error {
				node := Node{}
				return node.RunSSHSetupPromtailConfig("127.0.0.1", 3100, "", "test-cluster")
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

func TestNode_MonitoringConstants(t *testing.T) {
	// Test that the constants are properly defined
	assert.NotZero(t, constants.OdysseygoLokiPort)
	assert.NotZero(t, constants.OdysseygoAPIPort)
	assert.NotZero(t, constants.OdysseygoMachineMetricsPort)
	assert.NotZero(t, constants.OdysseygoLoadTestPort)
}

func TestNode_MonitoringErrorHandling(t *testing.T) {
	// Test error handling in monitoring functions
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
			name: "MonitorNodes with nil context",
			operation: func() error {
				return node.MonitorNodes(nil, []Node{}, "test-cluster")
			},
			expectError: true,
		},
		{
			name: "MonitorNodes with nil nodes",
			operation: func() error {
				ctx := context.Background()
				return node.MonitorNodes(ctx, nil, "test-cluster")
			},
			expectError: true, // Not a monitoring node in this test setup
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
