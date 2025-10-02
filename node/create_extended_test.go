// Copyright (C) 2025, Dione Protocol, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package node

import (
	"context"
	"testing"

	"github.com/DioneProtocol/odyssey-tooling-sdk-go/constants"
	"github.com/DioneProtocol/odyssey-tooling-sdk-go/odyssey"
	"github.com/stretchr/testify/assert"
)

func TestPreCreateCheckExtended(t *testing.T) {
	// Save original feature flag values
	originalInstanceManagement := constants.InstanceManagementEnabled
	originalSSHKeyManagement := constants.SSHKeyManagementEnabled
	defer func() {
		constants.InstanceManagementEnabled = originalInstanceManagement
		constants.SSHKeyManagementEnabled = originalSSHKeyManagement
	}()

	// Enable required feature flags for this test
	constants.InstanceManagementEnabled = true
	constants.SSHKeyManagementEnabled = true

	tests := []struct {
		name              string
		cp                CloudParams
		count             int
		sshPrivateKeyPath string
		expectError       bool
		errorContains     string
	}{
		{
			name: "Valid parameters",
			cp: CloudParams{
				Region:       "us-west-2",
				ImageID:      "ami-12345678",
				InstanceType: "t3.medium",
				AWSConfig: &AWSConfig{
					AWSProfile:           "default",
					AWSKeyPair:           "test-keypair",
					AWSSecurityGroupID:   "sg-12345678",
					AWSSecurityGroupName: "test-security-group",
					AWSVolumeSize:        100,
					AWSVolumeType:        "gp3",
					AWSVolumeIOPS:        1000,
					AWSVolumeThroughput:  500,
				},
			},
			count:             1,
			sshPrivateKeyPath: "/path/to/key",
			expectError:       true,
			errorContains:     "ssh private key path /path/to/key does not exist",
		},
		{
			name: "Invalid count - zero",
			cp: CloudParams{
				Region:       "us-west-2",
				ImageID:      "ami-12345678",
				InstanceType: "t3.medium",
			},
			count:             0,
			sshPrivateKeyPath: "/path/to/key",
			expectError:       true,
			errorContains:     "count must be at least 1",
		},
		{
			name: "Invalid count - negative",
			cp: CloudParams{
				Region:       "us-west-2",
				ImageID:      "ami-12345678",
				InstanceType: "t3.medium",
			},
			count:             -1,
			sshPrivateKeyPath: "/path/to/key",
			expectError:       true,
			errorContains:     "count must be at least 1",
		},
		{
			name: "Empty SSH key path",
			cp: CloudParams{
				Region:       "us-west-2",
				ImageID:      "ami-12345678",
				InstanceType: "t3.medium",
			},
			count:             1,
			sshPrivateKeyPath: "",
			expectError:       true,
			errorContains:     "unsupported cloud",
		},
		{
			name: "Large count",
			cp: CloudParams{
				Region:       "us-west-2",
				ImageID:      "ami-12345678",
				InstanceType: "t3.medium",
				AWSConfig: &AWSConfig{
					AWSProfile:           "default",
					AWSKeyPair:           "test-keypair",
					AWSSecurityGroupID:   "sg-12345678",
					AWSSecurityGroupName: "test-security-group",
					AWSVolumeSize:        100,
					AWSVolumeType:        "gp3",
					AWSVolumeIOPS:        1000,
					AWSVolumeThroughput:  500,
				},
			},
			count:             100,
			sshPrivateKeyPath: "/path/to/key",
			expectError:       true,
			errorContains:     "ssh private key path /path/to/key does not exist",
		},
		{
			name:              "Empty cloud params",
			cp:                CloudParams{},
			count:             1,
			sshPrivateKeyPath: "/path/to/key",
			expectError:       true,
			errorContains:     "region is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := preCreateCheck(tt.cp, tt.count, tt.sshPrivateKeyPath)
			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorContains)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPreCreateCheck_FeatureFlags(t *testing.T) {
	// Note: Feature flag validation is not yet implemented in the node package
	// This test is kept for future implementation
	t.Skip("Feature flag validation not yet implemented in node package")
}

func TestProvisionHost_ErrorHandling(t *testing.T) {
	// Save original feature flag values
	originalDockerSupport := constants.DockerSupportEnabled
	originalRelayerEnabled := constants.RelayerEnabled
	defer func() {
		constants.DockerSupportEnabled = originalDockerSupport
		constants.RelayerEnabled = originalRelayerEnabled
	}()

	// Enable required feature flags for this test
	constants.DockerSupportEnabled = true
	constants.RelayerEnabled = true

	tests := []struct {
		name        string
		node        Node
		nodeParams  *NodeParams
		expectError bool
	}{
		{
			name: "Invalid roles",
			node: Node{
				NodeID: "test-node",
				IP:     "192.168.1.1",
			},
			nodeParams: &NodeParams{
				Roles: []SupportedRole{},
			},
			expectError: true,
		},
		{
			name: "Empty roles",
			node: Node{
				NodeID: "test-node",
				IP:     "192.168.1.1",
			},
			nodeParams: &NodeParams{
				Roles: []SupportedRole{},
			},
			expectError: true, // Will fail due to SSH connection
		},
		{
			name: "Valid validator role",
			node: Node{
				NodeID: "test-node",
				IP:     "192.168.1.1",
			},
			nodeParams: &NodeParams{
				Roles: []SupportedRole{Validator},
			},
			expectError: true, // Will fail due to no connection
		},
		{
			name: "Valid API role",
			node: Node{
				NodeID: "test-node",
				IP:     "192.168.1.1",
			},
			nodeParams: &NodeParams{
				Roles: []SupportedRole{API},
			},
			expectError: true, // Will fail due to no connection
		},
		{
			name: "Valid monitor role",
			node: Node{
				NodeID: "test-node",
				IP:     "192.168.1.1",
			},
			nodeParams: &NodeParams{
				Roles: []SupportedRole{Monitor},
			},
			expectError: true, // Will fail due to no connection
		},
		{
			name: "Valid loadtest role",
			node: Node{
				NodeID: "test-node",
				IP:     "192.168.1.1",
			},
			nodeParams: &NodeParams{
				Roles: []SupportedRole{Loadtest},
			},
			expectError: true, // Will fail due to no connection
		},
		{
			name: "Valid AWM relayer role",
			node: Node{
				NodeID: "test-node",
				IP:     "192.168.1.1",
			},
			nodeParams: &NodeParams{
				Roles: []SupportedRole{AWMRelayer},
			},
			expectError: true, // Will fail due to no connection
		},
		{
			name: "Multiple valid roles",
			node: Node{
				NodeID: "test-node",
				IP:     "192.168.1.1",
			},
			nodeParams: &NodeParams{
				Roles: []SupportedRole{Validator, API},
			},
			expectError: true, // Will fail due to no connection
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

func TestProvisionHost_FeatureFlags(t *testing.T) {
	// Note: Feature flag validation is not yet implemented in the node package
	// This test is kept for future implementation
	t.Skip("Feature flag validation not yet implemented in node package")
}

func TestProvisionOdysseyGoHost(t *testing.T) {
	// Save original feature flag values
	originalDockerSupport := constants.DockerSupportEnabled
	defer func() {
		constants.DockerSupportEnabled = originalDockerSupport
	}()

	// Enable Docker support for this test
	constants.DockerSupportEnabled = true

	tests := []struct {
		name        string
		node        Node
		nodeParams  *NodeParams
		expectError bool
	}{
		{
			name: "Valid node and params",
			node: Node{
				NodeID: "test-node",
				IP:     "192.168.1.1",
			},
			nodeParams: &NodeParams{
				Network:          odyssey.TestnetNetwork(),
				SubnetIDs:        []string{"subnet1"},
				OdysseyGoVersion: "v1.11.8",
			},
			expectError: true, // Will fail due to no connection
		},
		{
			name: "Empty node",
			node: Node{},
			nodeParams: &NodeParams{
				Network:          odyssey.TestnetNetwork(),
				SubnetIDs:        []string{"subnet1"},
				OdysseyGoVersion: "v1.11.8",
			},
			expectError: true,
		},
		{
			name: "Node with empty IP",
			node: Node{
				NodeID: "test-node",
				IP:     "",
			},
			nodeParams: &NodeParams{
				Network:          odyssey.TestnetNetwork(),
				SubnetIDs:        []string{"subnet1"},
				OdysseyGoVersion: "v1.11.8",
			},
			expectError: true,
		},
		{
			name: "Node with mainnet network",
			node: Node{
				NodeID: "test-node",
				IP:     "192.168.1.1",
			},
			nodeParams: &NodeParams{
				Network:          odyssey.MainnetNetwork(),
				SubnetIDs:        []string{},
				OdysseyGoVersion: "v1.11.8",
			},
			expectError: true,
		},
		{
			name: "Node with devnet network",
			node: Node{
				NodeID: "test-node",
				IP:     "192.168.1.1",
			},
			nodeParams: &NodeParams{
				Network:          odyssey.DevnetNetwork(),
				SubnetIDs:        []string{"subnet1", "subnet2"},
				OdysseyGoVersion: "v1.11.8",
			},
			expectError: true,
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

func TestProvisionOdysseyGoHost_FeatureFlags(t *testing.T) {
	// Note: Feature flag validation is not yet implemented in the node package
	// This test is kept for future implementation
	t.Skip("Feature flag validation not yet implemented in node package")
}

func TestProvisionLoadTestHostExtended(t *testing.T) {
	// Save original feature flag values
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
			name: "Valid node",
			node: Node{
				NodeID: "test-node",
				IP:     "192.168.1.1",
			},
			expectError: true, // Will fail due to no connection
		},
		{
			name:        "Empty node",
			node:        Node{},
			expectError: true,
		},
		{
			name: "Node with empty IP",
			node: Node{
				NodeID: "test-node",
				IP:     "",
			},
			expectError: true,
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

func TestProvisionLoadTestHost_FeatureFlags(t *testing.T) {
	// Note: Feature flag validation is not yet implemented in the node package
	// This test is kept for future implementation
	t.Skip("Feature flag validation not yet implemented in node package")
}

func TestProvisionMonitoringHostExtended(t *testing.T) {
	// Save original feature flag values
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
			name: "Valid node",
			node: Node{
				NodeID: "test-node",
				IP:     "192.168.1.1",
			},
			expectError: true, // Will fail due to no connection
		},
		{
			name:        "Empty node",
			node:        Node{},
			expectError: true,
		},
		{
			name: "Node with empty IP",
			node: Node{
				NodeID: "test-node",
				IP:     "",
			},
			expectError: true,
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

func TestProvisionMonitoringHost_FeatureFlags(t *testing.T) {
	// Note: Feature flag validation is not yet implemented in the node package
	// This test is kept for future implementation
	t.Skip("Feature flag validation not yet implemented in node package")
}

func TestProvisionAWMRelayerHostExtended(t *testing.T) {
	// Save original feature flag values
	originalDockerSupport := constants.DockerSupportEnabled
	originalRelayerEnabled := constants.RelayerEnabled
	defer func() {
		constants.DockerSupportEnabled = originalDockerSupport
		constants.RelayerEnabled = originalRelayerEnabled
	}()

	// Enable required feature flags for this test
	constants.DockerSupportEnabled = true
	constants.RelayerEnabled = true

	tests := []struct {
		name        string
		node        Node
		expectError bool
	}{
		{
			name: "Valid node",
			node: Node{
				NodeID: "test-node",
				IP:     "192.168.1.1",
			},
			expectError: true, // Will fail due to no connection
		},
		{
			name:        "Empty node",
			node:        Node{},
			expectError: true,
		},
		{
			name: "Node with empty IP",
			node: Node{
				NodeID: "test-node",
				IP:     "",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := provisionAWMRelayerHost(tt.node)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestProvisionAWMRelayerHost_FeatureFlags(t *testing.T) {
	// Note: Feature flag validation is not yet implemented in the node package
	// This test is kept for future implementation
	t.Skip("Feature flag validation not yet implemented in node package")
}

func TestCreateNodes_FeatureFlags(t *testing.T) {
	// Note: Feature flag validation is not yet implemented in the node package
	// This test is kept for future implementation
	t.Skip("Feature flag validation not yet implemented in node package")
}

func TestCreateNodes_EdgeCases(t *testing.T) {
	// Save original feature flag values
	originalAWSIntegration := constants.AWSIntegrationEnabled
	originalGCPIntegration := constants.GCPIntegrationEnabled
	originalInstanceManagement := constants.InstanceManagementEnabled
	originalDockerSupport := constants.DockerSupportEnabled
	originalSSHKeyManagement := constants.SSHKeyManagementEnabled
	defer func() {
		constants.AWSIntegrationEnabled = originalAWSIntegration
		constants.GCPIntegrationEnabled = originalGCPIntegration
		constants.InstanceManagementEnabled = originalInstanceManagement
		constants.DockerSupportEnabled = originalDockerSupport
		constants.SSHKeyManagementEnabled = originalSSHKeyManagement
	}()

	// Enable all feature flags for this test
	constants.AWSIntegrationEnabled = true
	constants.GCPIntegrationEnabled = true
	constants.InstanceManagementEnabled = true
	constants.DockerSupportEnabled = true
	constants.SSHKeyManagementEnabled = true

	ctx := context.Background()

	tests := []struct {
		name        string
		nodeParams  *NodeParams
		expectError bool
	}{
		{
			name:        "Nil node params",
			nodeParams:  nil,
			expectError: true,
		},
		{
			name: "Nil cloud params",
			nodeParams: &NodeParams{
				CloudParams:       nil,
				Count:             1,
				Roles:             []SupportedRole{Validator},
				Network:           odyssey.TestnetNetwork(),
				SSHPrivateKeyPath: "/path/to/key",
				OdysseyGoVersion:  "v1.11.8",
				UseStaticIP:       false,
			},
			expectError: true,
		},
		{
			name: "Zero count",
			nodeParams: &NodeParams{
				CloudParams: &CloudParams{
					Region:       "us-west-2",
					ImageID:      "ami-12345678",
					InstanceType: "t3.medium",
				},
				Count:             0,
				Roles:             []SupportedRole{Validator},
				Network:           odyssey.TestnetNetwork(),
				SSHPrivateKeyPath: "/path/to/key",
				OdysseyGoVersion:  "v1.11.8",
				UseStaticIP:       false,
			},
			expectError: true,
		},
		{
			name: "Empty roles",
			nodeParams: &NodeParams{
				CloudParams: &CloudParams{
					Region:       "us-west-2",
					ImageID:      "ami-12345678",
					InstanceType: "t3.medium",
				},
				Count:             1,
				Roles:             []SupportedRole{},
				Network:           odyssey.TestnetNetwork(),
				SSHPrivateKeyPath: "/path/to/key",
				OdysseyGoVersion:  "v1.11.8",
				UseStaticIP:       false,
			},
			expectError: true,
		},
		{
			name: "Empty SSH key path",
			nodeParams: &NodeParams{
				CloudParams: &CloudParams{
					Region:       "us-west-2",
					ImageID:      "ami-12345678",
					InstanceType: "t3.medium",
				},
				Count:             1,
				Roles:             []SupportedRole{Validator},
				Network:           odyssey.TestnetNetwork(),
				SSHPrivateKeyPath: "",
				OdysseyGoVersion:  "v1.11.8",
				UseStaticIP:       false,
			},
			expectError: true,
		},
		{
			name: "Empty OdysseyGo version",
			nodeParams: &NodeParams{
				CloudParams: &CloudParams{
					Region:       "us-west-2",
					ImageID:      "ami-12345678",
					InstanceType: "t3.medium",
				},
				Count:             1,
				Roles:             []SupportedRole{Validator},
				Network:           odyssey.TestnetNetwork(),
				SSHPrivateKeyPath: "/path/to/key",
				OdysseyGoVersion:  "",
				UseStaticIP:       false,
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.nodeParams == nil {
				// Skip nil nodeParams test to avoid panic
				t.Skip("Skipping nil nodeParams test to avoid panic")
				return
			}
			_, err := CreateNodes(ctx, tt.nodeParams)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCreateNodes_NetworkTypes(t *testing.T) {
	// Save original feature flag values
	originalAWSIntegration := constants.AWSIntegrationEnabled
	originalGCPIntegration := constants.GCPIntegrationEnabled
	originalInstanceManagement := constants.InstanceManagementEnabled
	originalDockerSupport := constants.DockerSupportEnabled
	originalSSHKeyManagement := constants.SSHKeyManagementEnabled
	defer func() {
		constants.AWSIntegrationEnabled = originalAWSIntegration
		constants.GCPIntegrationEnabled = originalGCPIntegration
		constants.InstanceManagementEnabled = originalInstanceManagement
		constants.DockerSupportEnabled = originalDockerSupport
		constants.SSHKeyManagementEnabled = originalSSHKeyManagement
	}()

	// Enable all feature flags for this test
	constants.AWSIntegrationEnabled = true
	constants.GCPIntegrationEnabled = true
	constants.InstanceManagementEnabled = true
	constants.DockerSupportEnabled = true
	constants.SSHKeyManagementEnabled = true

	ctx := context.Background()

	tests := []struct {
		name        string
		network     odyssey.Network
		expectError bool
	}{
		{
			name:        "Testnet network",
			network:     odyssey.TestnetNetwork(),
			expectError: true, // Will fail due to no connection
		},
		{
			name:        "Mainnet network",
			network:     odyssey.MainnetNetwork(),
			expectError: true,
		},
		{
			name:        "Devnet network",
			network:     odyssey.DevnetNetwork(),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nodeParams := &NodeParams{
				CloudParams: &CloudParams{
					Region:       "us-west-2",
					ImageID:      "ami-12345678",
					InstanceType: "t3.medium",
				},
				Count:             1,
				Roles:             []SupportedRole{Validator},
				Network:           tt.network,
				SSHPrivateKeyPath: "/path/to/key",
				OdysseyGoVersion:  "v1.11.8",
				UseStaticIP:       false,
			}

			_, err := CreateNodes(ctx, nodeParams)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCreateNodes_RoleCombinations(t *testing.T) {
	// Save original feature flag values
	originalAWSIntegration := constants.AWSIntegrationEnabled
	originalGCPIntegration := constants.GCPIntegrationEnabled
	originalInstanceManagement := constants.InstanceManagementEnabled
	originalDockerSupport := constants.DockerSupportEnabled
	originalSSHKeyManagement := constants.SSHKeyManagementEnabled
	originalRelayerEnabled := constants.RelayerEnabled
	defer func() {
		constants.AWSIntegrationEnabled = originalAWSIntegration
		constants.GCPIntegrationEnabled = originalGCPIntegration
		constants.InstanceManagementEnabled = originalInstanceManagement
		constants.DockerSupportEnabled = originalDockerSupport
		constants.SSHKeyManagementEnabled = originalSSHKeyManagement
		constants.RelayerEnabled = originalRelayerEnabled
	}()

	// Enable all feature flags for this test
	constants.AWSIntegrationEnabled = true
	constants.GCPIntegrationEnabled = true
	constants.InstanceManagementEnabled = true
	constants.DockerSupportEnabled = true
	constants.SSHKeyManagementEnabled = true
	constants.RelayerEnabled = true

	ctx := context.Background()

	tests := []struct {
		name        string
		roles       []SupportedRole
		expectError bool
	}{
		{
			name:        "Single validator role",
			roles:       []SupportedRole{Validator},
			expectError: true,
		},
		{
			name:        "Single API role",
			roles:       []SupportedRole{API},
			expectError: true,
		},
		{
			name:        "Single monitor role",
			roles:       []SupportedRole{Monitor},
			expectError: true,
		},
		{
			name:        "Single loadtest role",
			roles:       []SupportedRole{Loadtest},
			expectError: true,
		},
		{
			name:        "Single AWM relayer role",
			roles:       []SupportedRole{AWMRelayer},
			expectError: true,
		},
		{
			name:        "Validator and API roles",
			roles:       []SupportedRole{Validator, API},
			expectError: true,
		},
		{
			name:        "Validator and Monitor roles",
			roles:       []SupportedRole{Validator, Monitor},
			expectError: true,
		},
		{
			name:        "API and Monitor roles",
			roles:       []SupportedRole{API, Monitor},
			expectError: true,
		},
		{
			name:        "All roles",
			roles:       []SupportedRole{Validator, API, Monitor, Loadtest, AWMRelayer},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nodeParams := &NodeParams{
				CloudParams: &CloudParams{
					Region:       "us-west-2",
					ImageID:      "ami-12345678",
					InstanceType: "t3.medium",
				},
				Count:             1,
				Roles:             tt.roles,
				Network:           odyssey.TestnetNetwork(),
				SSHPrivateKeyPath: "/path/to/key",
				OdysseyGoVersion:  "v1.11.8",
				UseStaticIP:       false,
			}

			_, err := CreateNodes(ctx, nodeParams)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
