// Copyright (C) 2025, Dione Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package node

import (
	"testing"

	"github.com/DioneProtocol/odyssey-tooling-sdk-go/constants"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProvisionHost_FeatureFlagsExtended(t *testing.T) {
	// Save original feature flag values
	originalSSHKeyManagement := constants.SSHKeyManagementEnabled
	originalDockerSupport := constants.DockerSupportEnabled

	defer func() {
		constants.SSHKeyManagementEnabled = originalSSHKeyManagement
		constants.DockerSupportEnabled = originalDockerSupport
	}()

	node := Node{
		NodeID: "test-node",
		IP:     "192.168.1.1",
	}

	nodeParams := &NodeParams{
		Roles: []SupportedRole{Validator},
	}

	tests := []struct {
		name          string
		sshEnabled    bool
		dockerEnabled bool
		expectError   bool
		errorContains string
	}{
		{
			name:          "SSH key management disabled",
			sshEnabled:    false,
			dockerEnabled: true,
			expectError:   true,
			errorContains: "SSH key management functionality is disabled",
		},
		{
			name:          "Docker support disabled",
			sshEnabled:    true,
			dockerEnabled: false,
			expectError:   true,
			errorContains: "Docker support functionality is disabled",
		},
		{
			name:          "Both disabled",
			sshEnabled:    false,
			dockerEnabled: false,
			expectError:   true,
			errorContains: "SSH key management functionality is disabled",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set feature flags
			constants.SSHKeyManagementEnabled = tt.sshEnabled
			constants.DockerSupportEnabled = tt.dockerEnabled

			err := provisionHost(node, nodeParams)

			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorContains)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestProvisionHost_RoleValidation(t *testing.T) {
	// Save original feature flag values
	originalSSHKeyManagement := constants.SSHKeyManagementEnabled
	originalDockerSupport := constants.DockerSupportEnabled

	defer func() {
		constants.SSHKeyManagementEnabled = originalSSHKeyManagement
		constants.DockerSupportEnabled = originalDockerSupport
	}()

	// Enable feature flags
	constants.SSHKeyManagementEnabled = true
	constants.DockerSupportEnabled = true

	node := Node{
		NodeID: "test-node",
		IP:     "192.168.1.1",
	}

	tests := []struct {
		name          string
		roles         []SupportedRole
		expectError   bool
		errorContains string
	}{
		{
			name:          "Empty roles",
			roles:         []SupportedRole{},
			expectError:   true,
			errorContains: "roles cannot be empty",
		},
		{
			name:          "Nil roles",
			roles:         nil,
			expectError:   true,
			errorContains: "roles cannot be empty",
		},
		{
			name:          "Invalid role",
			roles:         []SupportedRole{SupportedRole(99)},
			expectError:   true,
			errorContains: "unsupported role",
		},
		{
			name:          "Mixed valid and invalid roles",
			roles:         []SupportedRole{Validator, SupportedRole(99)},
			expectError:   true,
			errorContains: "unsupported role",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nodeParams := &NodeParams{
				Roles: tt.roles,
			}

			err := provisionHost(node, nodeParams)

			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorContains)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestProvisionHost_ValidRoles(t *testing.T) {
	// Save original feature flag values
	originalSSHKeyManagement := constants.SSHKeyManagementEnabled
	originalDockerSupport := constants.DockerSupportEnabled

	defer func() {
		constants.SSHKeyManagementEnabled = originalSSHKeyManagement
		constants.DockerSupportEnabled = originalDockerSupport
	}()

	// Enable feature flags
	constants.SSHKeyManagementEnabled = true
	constants.DockerSupportEnabled = true

	node := Node{
		NodeID: "test-node",
		IP:     "192.168.1.1",
	}

	tests := []struct {
		name        string
		roles       []SupportedRole
		expectError bool
	}{
		{
			name:        "Single validator role",
			roles:       []SupportedRole{Validator},
			expectError: true, // Will fail due to SSH connection
		},
		{
			name:        "Single API role",
			roles:       []SupportedRole{API},
			expectError: true, // Will fail due to SSH connection
		},
		{
			name:        "Single loadtest role",
			roles:       []SupportedRole{Loadtest},
			expectError: true, // Will fail due to SSH connection
		},
		{
			name:        "Single monitor role",
			roles:       []SupportedRole{Monitor},
			expectError: true, // Will fail due to SSH connection
		},
		{
			name:        "Multiple roles",
			roles:       []SupportedRole{Validator, API},
			expectError: true, // Will fail due to SSH connection
		},
		{
			name:        "All roles",
			roles:       []SupportedRole{Validator, API, Loadtest, Monitor},
			expectError: true, // Will fail due to SSH connection
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nodeParams := &NodeParams{
				Roles: tt.roles,
			}

			err := provisionHost(node, nodeParams)

			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestProvisionHost_EdgeCases(t *testing.T) {
	// Save original feature flag values
	originalSSHKeyManagement := constants.SSHKeyManagementEnabled
	originalDockerSupport := constants.DockerSupportEnabled

	defer func() {
		constants.SSHKeyManagementEnabled = originalSSHKeyManagement
		constants.DockerSupportEnabled = originalDockerSupport
	}()

	// Enable feature flags
	constants.SSHKeyManagementEnabled = true
	constants.DockerSupportEnabled = true

	tests := []struct {
		name          string
		node          Node
		nodeParams    *NodeParams
		expectError   bool
		errorContains string
	}{
		{
			name: "Nil node params",
			node: Node{
				NodeID: "test-node",
				IP:     "192.168.1.1",
			},
			nodeParams:    nil,
			expectError:   true,
			errorContains: "nodeParams cannot be nil",
		},
		{
			name: "Empty node ID",
			node: Node{
				NodeID: "",
				IP:     "192.168.1.1",
			},
			nodeParams: &NodeParams{
				Roles: []SupportedRole{Validator},
			},
			expectError:   true,
			errorContains: "node ID is required",
		},
		{
			name: "Empty IP",
			node: Node{
				NodeID: "test-node",
				IP:     "",
			},
			nodeParams: &NodeParams{
				Roles: []SupportedRole{Validator},
			},
			expectError:   true,
			errorContains: "IP address is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := provisionHost(tt.node, tt.nodeParams)

			if tt.expectError {
				require.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
}
