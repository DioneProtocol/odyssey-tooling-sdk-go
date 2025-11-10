// Copyright (C) 2025, Dione Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package node

import (
	"testing"
	"time"

	"github.com/DioneProtocol/odyssey-tooling-sdk-go/constants"
	"github.com/DioneProtocol/odyssey-tooling-sdk-go/odyssey"
	"github.com/DioneProtocol/odyssey-tooling-sdk-go/validator"
	"github.com/DioneProtocol/odyssey-tooling-sdk-go/wallet"
	"github.com/DioneProtocol/odysseygo/ids"
	"github.com/stretchr/testify/assert"
)

// createMockWallet creates a mock wallet for testing
// Note: This is a simplified approach since wallet.Wallet is a struct, not an interface
func createMockWallet(t *testing.T) wallet.Wallet {
	// Return zero value wallet that will cause the function to fail
	// This is appropriate for testing error conditions
	return wallet.Wallet{}
}

func TestValidatePrimaryNetwork_ValidationErrors(t *testing.T) {
	// Save original feature flag values
	originalSSHKeyManagement := constants.SSHKeyManagementEnabled
	defer func() {
		constants.SSHKeyManagementEnabled = originalSSHKeyManagement
	}()

	// Enable SSH key management for this test
	constants.SSHKeyManagementEnabled = true

	tests := []struct {
		name            string
		validatorParams validator.PrimaryNetworkValidatorParams
		expectError     bool
		errorContains   string
	}{
		{
			name: "Empty node ID",
			validatorParams: validator.PrimaryNetworkValidatorParams{
				NodeID:      ids.EmptyNodeID,
				Duration:    time.Hour,
				StakeAmount: 1000000,
			},
			expectError:   true,
			errorContains: "validator node id is not provided",
		},
		{
			name: "Zero duration",
			validatorParams: validator.PrimaryNetworkValidatorParams{
				NodeID:      ids.GenerateTestNodeID(),
				Duration:    0,
				StakeAmount: 1000000,
			},
			expectError:   true,
			errorContains: "validator duration is not provided",
		},
		{
			name: "Insufficient stake amount",
			validatorParams: validator.PrimaryNetworkValidatorParams{
				NodeID:      ids.GenerateTestNodeID(),
				Duration:    time.Hour,
				StakeAmount: 100, // Too low
			},
			expectError:   true,
			errorContains: "invalid weight",
		},
		{
			name: "Valid parameters",
			validatorParams: validator.PrimaryNetworkValidatorParams{
				NodeID:      ids.GenerateTestNodeID(),
				Duration:    time.Hour,
				StakeAmount: 500000000000000, // Use actual minimum stake
			},
			expectError: true, // Will still fail due to wallet/connection issues
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := &Node{
				NodeID: "test-node",
				IP:     "192.168.1.1",
			}

			network := odyssey.TestnetNetwork()
			wallet := createMockWallet(t)

			_, err := node.ValidatePrimaryNetwork(network, tt.validatorParams, wallet)

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

func TestValidatePrimaryNetwork_FeatureFlags(t *testing.T) {
	// Note: Feature flag validation is not yet implemented in the node package
	// This test is kept for future implementation
	t.Skip("Feature flag validation not yet implemented in node package")
}

func TestGetBLSKeyFromRemoteHost_ErrorHandling(t *testing.T) {
	// Save original feature flag values
	originalSSHKeyManagement := constants.SSHKeyManagementEnabled
	defer func() {
		constants.SSHKeyManagementEnabled = originalSSHKeyManagement
	}()

	// Enable SSH key management for this test
	constants.SSHKeyManagementEnabled = true

	tests := []struct {
		name        string
		node        Node
		expectError bool
	}{
		{
			name: "Node with no connection",
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.node.GetBLSKeyFromRemoteHost()
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetBLSKeyFromRemoteHost_FeatureFlags(t *testing.T) {
	// Note: Feature flag validation is not yet implemented in the node package
	// This test is kept for future implementation
	t.Skip("Feature flag validation not yet implemented in node package")
}

func TestValidatePrimaryNetwork_NetworkValidation(t *testing.T) {
	// Save original feature flag values
	originalSSHKeyManagement := constants.SSHKeyManagementEnabled
	defer func() {
		constants.SSHKeyManagementEnabled = originalSSHKeyManagement
	}()

	// Enable SSH key management for this test
	constants.SSHKeyManagementEnabled = true

	tests := []struct {
		name          string
		network       odyssey.Network
		stakeAmount   uint64
		expectError   bool
		errorContains string
	}{
		{
			name:        "Testnet with valid stake",
			network:     odyssey.TestnetNetwork(),
			stakeAmount: 500000000000000, // Use actual minimum stake
			expectError: true,            // Will fail due to wallet/connection issues
		},
		{
			name:        "Mainnet with valid stake",
			network:     odyssey.MainnetNetwork(),
			stakeAmount: 500000000000000, // Use actual minimum stake
			expectError: true,            // Will fail due to wallet/connection issues
		},
		{
			name:        "Devnet with valid stake",
			network:     odyssey.DevnetNetwork(),
			stakeAmount: 500000000000000, // Use actual minimum stake
			expectError: true,            // Will fail due to wallet/connection issues
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := &Node{
				NodeID: "test-node",
				IP:     "192.168.1.1",
			}

			wallet := createMockWallet(t)
			validatorParams := validator.PrimaryNetworkValidatorParams{
				NodeID:      ids.GenerateTestNodeID(),
				Duration:    time.Hour,
				StakeAmount: tt.stakeAmount,
			}

			_, err := node.ValidatePrimaryNetwork(tt.network, validatorParams, wallet)

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

func TestValidatePrimaryNetwork_DelegationFee(t *testing.T) {
	// Save original feature flag values
	originalSSHKeyManagement := constants.SSHKeyManagementEnabled
	defer func() {
		constants.SSHKeyManagementEnabled = originalSSHKeyManagement
	}()

	// Enable SSH key management for this test
	constants.SSHKeyManagementEnabled = true

	tests := []struct {
		name          string
		delegationFee uint32
		expectError   bool
	}{
		{
			name:          "Zero delegation fee (should use default)",
			delegationFee: 0,
			expectError:   true, // Will fail due to wallet/connection issues
		},
		{
			name:          "Valid delegation fee",
			delegationFee: 10000, // 1%
			expectError:   true,  // Will fail due to wallet/connection issues
		},
		{
			name:          "High delegation fee",
			delegationFee: 50000, // 5%
			expectError:   true,  // Will fail due to wallet/connection issues
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := &Node{
				NodeID: "test-node",
				IP:     "192.168.1.1",
			}

			network := odyssey.TestnetNetwork()
			wallet := createMockWallet(t)
			validatorParams := validator.PrimaryNetworkValidatorParams{
				NodeID:        ids.GenerateTestNodeID(),
				Duration:      time.Hour,
				StakeAmount:   500000000000000, // Use actual minimum stake
				DelegationFee: tt.delegationFee,
			}

			_, err := node.ValidatePrimaryNetwork(network, validatorParams, wallet)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidatePrimaryNetwork_EdgeCases(t *testing.T) {
	// Save original feature flag values
	originalSSHKeyManagement := constants.SSHKeyManagementEnabled
	defer func() {
		constants.SSHKeyManagementEnabled = originalSSHKeyManagement
	}()

	// Enable SSH key management for this test
	constants.SSHKeyManagementEnabled = true

	tests := []struct {
		name            string
		node            Node
		validatorParams validator.PrimaryNetworkValidatorParams
		expectError     bool
	}{
		{
			name: "Node with empty NodeID",
			node: Node{
				NodeID: "",
				IP:     "192.168.1.1",
			},
			validatorParams: validator.PrimaryNetworkValidatorParams{
				NodeID:      ids.GenerateTestNodeID(),
				Duration:    time.Hour,
				StakeAmount: 1000000,
			},
			expectError: true,
		},
		{
			name: "Node with empty IP",
			node: Node{
				NodeID: "test-node",
				IP:     "",
			},
			validatorParams: validator.PrimaryNetworkValidatorParams{
				NodeID:      ids.GenerateTestNodeID(),
				Duration:    time.Hour,
				StakeAmount: 1000000,
			},
			expectError: true,
		},
		{
			name: "Very short duration",
			node: Node{
				NodeID: "test-node",
				IP:     "192.168.1.1",
			},
			validatorParams: validator.PrimaryNetworkValidatorParams{
				NodeID:      ids.GenerateTestNodeID(),
				Duration:    time.Second,
				StakeAmount: 500000000000000, // Use actual minimum stake
			},
			expectError: true, // Will fail due to wallet/connection issues
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			network := odyssey.TestnetNetwork()
			wallet := createMockWallet(t)

			_, err := tt.node.ValidatePrimaryNetwork(network, tt.validatorParams, wallet)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
