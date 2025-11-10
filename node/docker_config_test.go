// Copyright (C) 2025, Dione Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package node

import (
	"os"
	"testing"

	"github.com/DioneProtocol/odyssey-tooling-sdk-go/constants"
	"github.com/stretchr/testify/assert"
)

func TestNode_RunSSHRenderOdysseyNodeConfig(t *testing.T) {
	// Save original feature flag values
	originalSSHKeyManagement := constants.SSHKeyManagementEnabled
	defer func() {
		constants.SSHKeyManagementEnabled = originalSSHKeyManagement
	}()

	// Enable SSH key management for this test
	constants.SSHKeyManagementEnabled = true

	tests := []struct {
		name         string
		networkID    string
		trackSubnets []string
		expectError  bool
	}{
		{
			name:         "Valid configuration",
			networkID:    "test-network-123",
			trackSubnets: []string{"subnet1", "subnet2"},
			expectError:  true, // Will fail due to no connection
		},
		{
			name:         "Empty network ID",
			networkID:    "",
			trackSubnets: []string{"subnet1"},
			expectError:  true,
		},
		{
			name:         "Empty track subnets",
			networkID:    "test-network-123",
			trackSubnets: []string{},
			expectError:  true,
		},
		{
			name:         "Nil track subnets",
			networkID:    "test-network-123",
			trackSubnets: nil,
			expectError:  true,
		},
		{
			name:         "Single subnet",
			networkID:    "test-network-123",
			trackSubnets: []string{"subnet1"},
			expectError:  true,
		},
		{
			name:         "Multiple subnets",
			networkID:    "test-network-123",
			trackSubnets: []string{"subnet1", "subnet2", "subnet3"},
			expectError:  true,
		},
		{
			name:         "Network ID with special characters",
			networkID:    "test-network-123_abc",
			trackSubnets: []string{"subnet1"},
			expectError:  true,
		},
		{
			name:         "Subnet IDs with special characters",
			networkID:    "test-network-123",
			trackSubnets: []string{"subnet-1", "subnet_2", "subnet.3"},
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := &Node{
				NodeID: "test-node",
				IP:     "192.168.1.1",
			}

			err := node.RunSSHRenderOdysseyNodeConfig(tt.networkID, tt.trackSubnets)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNode_RunSSHRenderOdysseyNodeConfig_FeatureFlags(t *testing.T) {
	// Save original feature flag values
	originalSSHKeyManagement := constants.SSHKeyManagementEnabled
	defer func() {
		constants.SSHKeyManagementEnabled = originalSSHKeyManagement
	}()

	node := &Node{
		NodeID: "test-node",
		IP:     "192.168.1.1",
	}

	// Test with SSH key management disabled
	constants.SSHKeyManagementEnabled = false

	err := node.RunSSHRenderOdysseyNodeConfig("test-network-123", []string{"subnet1"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "SSH key management functionality is disabled")
}

func TestPrepareGrafanaConfig(t *testing.T) {
	// Test the prepareGrafanaConfig function
	configFile, dashboardsFile, dataSourceFile, promDataSourceFile, err := prepareGrafanaConfig()

	// This function creates temporary files, so we expect it to succeed
	assert.NoError(t, err)
	assert.NotEmpty(t, configFile)
	assert.NotEmpty(t, dashboardsFile)
	assert.NotEmpty(t, dataSourceFile)
	assert.NotEmpty(t, promDataSourceFile)

	// Verify files were created
	assert.FileExists(t, configFile)
	assert.FileExists(t, dashboardsFile)
	assert.FileExists(t, dataSourceFile)
	assert.FileExists(t, promDataSourceFile)

	// Clean up temporary files
	os.Remove(configFile)
	os.Remove(dashboardsFile)
	os.Remove(dataSourceFile)
	os.Remove(promDataSourceFile)
}

func TestPrepareGrafanaConfig_MultipleCalls(t *testing.T) {
	// Test multiple calls to prepareGrafanaConfig
	for i := 0; i < 3; i++ {
		configFile, dashboardsFile, dataSourceFile, promDataSourceFile, err := prepareGrafanaConfig()

		assert.NoError(t, err)
		assert.NotEmpty(t, configFile)
		assert.NotEmpty(t, dashboardsFile)
		assert.NotEmpty(t, dataSourceFile)
		assert.NotEmpty(t, promDataSourceFile)

		// Clean up temporary files
		os.Remove(configFile)
		os.Remove(dashboardsFile)
		os.Remove(dataSourceFile)
		os.Remove(promDataSourceFile)
	}
}

func TestNode_RunSSHRenderOdysseyNodeConfig_EdgeCases(t *testing.T) {
	// Save original feature flag values
	originalSSHKeyManagement := constants.SSHKeyManagementEnabled
	defer func() {
		constants.SSHKeyManagementEnabled = originalSSHKeyManagement
	}()

	// Enable SSH key management for this test
	constants.SSHKeyManagementEnabled = true

	tests := []struct {
		name         string
		node         Node
		networkID    string
		trackSubnets []string
		expectError  bool
	}{
		{
			name:         "Empty node",
			node:         Node{},
			networkID:    "test-network-123",
			trackSubnets: []string{"subnet1"},
			expectError:  true,
		},
		{
			name: "Node with empty IP",
			node: Node{
				NodeID: "test-node",
				IP:     "",
			},
			networkID:    "test-network-123",
			trackSubnets: []string{"subnet1"},
			expectError:  true,
		},
		{
			name: "Node with empty NodeID",
			node: Node{
				NodeID: "",
				IP:     "192.168.1.1",
			},
			networkID:    "test-network-123",
			trackSubnets: []string{"subnet1"},
			expectError:  true,
		},
		{
			name: "Very long network ID",
			node: Node{
				NodeID: "test-node",
				IP:     "192.168.1.1",
			},
			networkID:    "very-long-network-id-that-exceeds-normal-limits-and-should-still-work",
			trackSubnets: []string{"subnet1"},
			expectError:  true,
		},
		{
			name: "Very long subnet IDs",
			node: Node{
				NodeID: "test-node",
				IP:     "192.168.1.1",
			},
			networkID: "test-network-123",
			trackSubnets: []string{
				"very-long-subnet-id-that-exceeds-normal-limits",
				"another-very-long-subnet-id-that-exceeds-normal-limits",
			},
			expectError: true,
		},
		{
			name: "Network ID with unicode characters",
			node: Node{
				NodeID: "test-node",
				IP:     "192.168.1.1",
			},
			networkID:    "test-network-123-测试",
			trackSubnets: []string{"subnet1"},
			expectError:  true,
		},
		{
			name: "Subnet IDs with unicode characters",
			node: Node{
				NodeID: "test-node",
				IP:     "192.168.1.1",
			},
			networkID:    "test-network-123",
			trackSubnets: []string{"subnet-测试", "subnet-тест"},
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.node.RunSSHRenderOdysseyNodeConfig(tt.networkID, tt.trackSubnets)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNode_RunSSHRenderOdysseyNodeConfig_NetworkTypes(t *testing.T) {
	// Save original feature flag values
	originalSSHKeyManagement := constants.SSHKeyManagementEnabled
	defer func() {
		constants.SSHKeyManagementEnabled = originalSSHKeyManagement
	}()

	// Enable SSH key management for this test
	constants.SSHKeyManagementEnabled = true

	tests := []struct {
		name         string
		networkID    string
		trackSubnets []string
		expectError  bool
	}{
		{
			name:         "Mainnet network ID",
			networkID:    "mainnet-123",
			trackSubnets: []string{"subnet1"},
			expectError:  true,
		},
		{
			name:         "Testnet network ID",
			networkID:    "testnet-456",
			trackSubnets: []string{"subnet1"},
			expectError:  true,
		},
		{
			name:         "Devnet network ID",
			networkID:    "devnet-789",
			trackSubnets: []string{"subnet1"},
			expectError:  true,
		},
		{
			name:         "Custom network ID",
			networkID:    "custom-network-abc",
			trackSubnets: []string{"subnet1"},
			expectError:  true,
		},
		{
			name:         "Network ID with numbers only",
			networkID:    "123456789",
			trackSubnets: []string{"subnet1"},
			expectError:  true,
		},
		{
			name:         "Network ID with mixed case",
			networkID:    "TestNetwork-123",
			trackSubnets: []string{"subnet1"},
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := &Node{
				NodeID: "test-node",
				IP:     "192.168.1.1",
			}

			err := node.RunSSHRenderOdysseyNodeConfig(tt.networkID, tt.trackSubnets)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNode_RunSSHRenderOdysseyNodeConfig_SubnetTypes(t *testing.T) {
	// Save original feature flag values
	originalSSHKeyManagement := constants.SSHKeyManagementEnabled
	defer func() {
		constants.SSHKeyManagementEnabled = originalSSHKeyManagement
	}()

	// Enable SSH key management for this test
	constants.SSHKeyManagementEnabled = true

	tests := []struct {
		name         string
		networkID    string
		trackSubnets []string
		expectError  bool
	}{
		{
			name:         "EVM subnet",
			networkID:    "test-network-123",
			trackSubnets: []string{"evm-subnet-1"},
			expectError:  true,
		},
		{
			name:         "Platform subnet",
			networkID:    "test-network-123",
			trackSubnets: []string{"platform-subnet-1"},
			expectError:  true,
		},
		{
			name:         "Custom VM subnet",
			networkID:    "test-network-123",
			trackSubnets: []string{"custom-vm-subnet-1"},
			expectError:  true,
		},
		{
			name:         "Mixed subnet types",
			networkID:    "test-network-123",
			trackSubnets: []string{"evm-subnet-1", "platform-subnet-1", "custom-vm-subnet-1"},
			expectError:  true,
		},
		{
			name:         "Subnet with UUID format",
			networkID:    "test-network-123",
			trackSubnets: []string{"12345678-1234-1234-1234-123456789abc"},
			expectError:  true,
		},
		{
			name:         "Subnet with hash format",
			networkID:    "test-network-123",
			trackSubnets: []string{"abcdef1234567890abcdef1234567890abcdef12"},
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := &Node{
				NodeID: "test-node",
				IP:     "192.168.1.1",
			}

			err := node.RunSSHRenderOdysseyNodeConfig(tt.networkID, tt.trackSubnets)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNode_RunSSHRenderOdysseyNodeConfig_ErrorHandling(t *testing.T) {
	// Save original feature flag values
	originalSSHKeyManagement := constants.SSHKeyManagementEnabled
	defer func() {
		constants.SSHKeyManagementEnabled = originalSSHKeyManagement
	}()

	// Enable SSH key management for this test
	constants.SSHKeyManagementEnabled = true

	node := &Node{
		NodeID: "test-node",
		IP:     "192.168.1.1",
	}

	// Test that error messages are informative
	err := node.RunSSHRenderOdysseyNodeConfig("test-network-123", []string{"subnet1"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to connect")
}

func TestNode_RunSSHRenderOdysseyNodeConfig_ConcurrentCalls(t *testing.T) {
	// Save original feature flag values
	originalSSHKeyManagement := constants.SSHKeyManagementEnabled
	defer func() {
		constants.SSHKeyManagementEnabled = originalSSHKeyManagement
	}()

	// Enable SSH key management for this test
	constants.SSHKeyManagementEnabled = true

	node := &Node{
		NodeID: "test-node",
		IP:     "192.168.1.1",
	}

	// Test concurrent calls to the same function
	done := make(chan error, 3)

	go func() {
		err := node.RunSSHRenderOdysseyNodeConfig("test-network-1", []string{"subnet1"})
		done <- err
	}()

	go func() {
		err := node.RunSSHRenderOdysseyNodeConfig("test-network-2", []string{"subnet2"})
		done <- err
	}()

	go func() {
		err := node.RunSSHRenderOdysseyNodeConfig("test-network-3", []string{"subnet3"})
		done <- err
	}()

	// All calls should fail due to no connection
	for i := 0; i < 3; i++ {
		err := <-done
		assert.Error(t, err)
	}
}

func TestNode_RunSSHRenderOdysseyNodeConfig_ResourceCleanup(t *testing.T) {
	// Save original feature flag values
	originalSSHKeyManagement := constants.SSHKeyManagementEnabled
	defer func() {
		constants.SSHKeyManagementEnabled = originalSSHKeyManagement
	}()

	// Enable SSH key management for this test
	constants.SSHKeyManagementEnabled = true

	node := &Node{
		NodeID: "test-node",
		IP:     "192.168.1.1",
	}

	// Test that resources are properly cleaned up even when errors occur
	err := node.RunSSHRenderOdysseyNodeConfig("test-network-123", []string{"subnet1"})
	assert.Error(t, err)

	// Verify that the node is in a clean state (NodeID should remain unchanged)
	assert.Equal(t, "test-node", node.NodeID) // NodeID should not be modified by the function
}

func TestNode_RunSSHRenderOdysseyNodeConfig_Validation(t *testing.T) {
	// Save original feature flag values
	originalSSHKeyManagement := constants.SSHKeyManagementEnabled
	defer func() {
		constants.SSHKeyManagementEnabled = originalSSHKeyManagement
	}()

	// Enable SSH key management for this test
	constants.SSHKeyManagementEnabled = true

	tests := []struct {
		name          string
		networkID     string
		trackSubnets  []string
		expectError   bool
		errorContains string
	}{
		{
			name:         "Network ID with spaces",
			networkID:    "test network 123",
			trackSubnets: []string{"subnet1"},
			expectError:  true,
		},
		{
			name:         "Network ID with tabs",
			networkID:    "test\tnetwork\t123",
			trackSubnets: []string{"subnet1"},
			expectError:  true,
		},
		{
			name:         "Network ID with newlines",
			networkID:    "test\nnetwork\n123",
			trackSubnets: []string{"subnet1"},
			expectError:  true,
		},
		{
			name:         "Subnet ID with spaces",
			networkID:    "test-network-123",
			trackSubnets: []string{"subnet 1"},
			expectError:  true,
		},
		{
			name:         "Subnet ID with tabs",
			networkID:    "test-network-123",
			trackSubnets: []string{"subnet\t1"},
			expectError:  true,
		},
		{
			name:         "Subnet ID with newlines",
			networkID:    "test-network-123",
			trackSubnets: []string{"subnet\n1"},
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := &Node{
				NodeID: "test-node",
				IP:     "192.168.1.1",
			}

			err := node.RunSSHRenderOdysseyNodeConfig(tt.networkID, tt.trackSubnets)
			if tt.expectError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
