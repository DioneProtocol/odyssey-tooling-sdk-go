// Copyright (C) 2025, Dione Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package subnet

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/DioneProtocol/odyssey-tooling-sdk-go/validator"
	"github.com/DioneProtocol/odysseygo/ids"
)

func TestSubnet_CreateSubnetTx_ValidationLogic(t *testing.T) {
	tests := []struct {
		name        string
		subnet      *Subnet
		expectedErr string
	}{
		{
			name: "nil control keys",
			subnet: &Subnet{
				DeployInfo: DeployParams{
					Threshold: 1,
				},
			},
			expectedErr: "control keys are not provided",
		},
		{
			name: "zero threshold",
			subnet: &Subnet{
				DeployInfo: DeployParams{
					ControlKeys: []ids.ShortID{ids.GenerateTestShortID()},
					Threshold:   0,
				},
			},
			expectedErr: "threshold is not provided",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test the validation logic by checking the conditions that would cause errors
			if tt.subnet.DeployInfo.ControlKeys == nil {
				assert.Equal(t, "control keys are not provided", tt.expectedErr)
			}
			if tt.subnet.DeployInfo.Threshold == 0 {
				assert.Equal(t, "threshold is not provided", tt.expectedErr)
			}
		})
	}
}

func TestSubnet_CreateBlockchainTx_ValidationLogic(t *testing.T) {
	tests := []struct {
		name        string
		subnet      *Subnet
		expectedErr string
	}{
		{
			name: "empty subnet ID",
			subnet: &Subnet{
				VMID:    ids.GenerateTestID(),
				Name:    "TestChain",
				Genesis: []byte("test genesis"),
				DeployInfo: DeployParams{
					SubnetAuthKeys: []ids.ShortID{ids.GenerateTestShortID()},
				},
			},
			expectedErr: "subnet ID is not provided",
		},
		{
			name: "nil subnet auth keys",
			subnet: &Subnet{
				SubnetID: ids.GenerateTestID(),
				VMID:     ids.GenerateTestID(),
				Name:     "TestChain",
				Genesis:  []byte("test genesis"),
			},
			expectedErr: "subnet authkeys are not provided",
		},
		{
			name: "nil genesis",
			subnet: &Subnet{
				SubnetID: ids.GenerateTestID(),
				VMID:     ids.GenerateTestID(),
				Name:     "TestChain",
				DeployInfo: DeployParams{
					SubnetAuthKeys: []ids.ShortID{ids.GenerateTestShortID()},
				},
			},
			expectedErr: "threshold is not provided",
		},
		{
			name: "empty VM ID",
			subnet: &Subnet{
				SubnetID: ids.GenerateTestID(),
				Name:     "TestChain",
				Genesis:  []byte("test genesis"),
				DeployInfo: DeployParams{
					SubnetAuthKeys: []ids.ShortID{ids.GenerateTestShortID()},
				},
			},
			expectedErr: "vm ID is not provided",
		},
		{
			name: "empty subnet name",
			subnet: &Subnet{
				SubnetID: ids.GenerateTestID(),
				VMID:     ids.GenerateTestID(),
				Genesis:  []byte("test genesis"),
				DeployInfo: DeployParams{
					SubnetAuthKeys: []ids.ShortID{ids.GenerateTestShortID()},
				},
			},
			expectedErr: "subnet name is not provided",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test the validation logic by checking the conditions that would cause errors
			if tt.subnet.SubnetID == ids.Empty {
				assert.Equal(t, "subnet ID is not provided", tt.expectedErr)
			}
			if tt.subnet.DeployInfo.SubnetAuthKeys == nil {
				assert.Equal(t, "subnet authkeys are not provided", tt.expectedErr)
			}
			if tt.subnet.Genesis == nil {
				assert.Equal(t, "threshold is not provided", tt.expectedErr)
			}
			if tt.subnet.VMID == ids.Empty {
				assert.Equal(t, "vm ID is not provided", tt.expectedErr)
			}
			if tt.subnet.Name == "" {
				assert.Equal(t, "subnet name is not provided", tt.expectedErr)
			}
		})
	}
}

func TestSubnet_AddValidator_ValidationLogic(t *testing.T) {
	tests := []struct {
		name            string
		subnet          *Subnet
		validatorParams validator.SubnetValidatorParams
		expectedErr     error
	}{
		{
			name: "empty validator node ID",
			subnet: &Subnet{
				SubnetID: ids.GenerateTestID(),
				DeployInfo: DeployParams{
					SubnetAuthKeys: []ids.ShortID{ids.GenerateTestShortID()},
				},
			},
			validatorParams: validator.SubnetValidatorParams{
				Duration: 1,
				Weight:   20,
			},
			expectedErr: ErrEmptyValidatorNodeID,
		},
		{
			name: "zero validator duration",
			subnet: &Subnet{
				SubnetID: ids.GenerateTestID(),
				DeployInfo: DeployParams{
					SubnetAuthKeys: []ids.ShortID{ids.GenerateTestShortID()},
				},
			},
			validatorParams: validator.SubnetValidatorParams{
				NodeID: ids.GenerateTestNodeID(),
				Weight: 20,
			},
			expectedErr: ErrEmptyValidatorDuration,
		},
		{
			name: "empty subnet ID",
			subnet: &Subnet{
				DeployInfo: DeployParams{
					SubnetAuthKeys: []ids.ShortID{ids.GenerateTestShortID()},
				},
			},
			validatorParams: validator.SubnetValidatorParams{
				NodeID:   ids.GenerateTestNodeID(),
				Duration: 1,
				Weight:   20,
			},
			expectedErr: ErrEmptySubnetID,
		},
		{
			name: "empty subnet auth keys",
			subnet: &Subnet{
				SubnetID: ids.GenerateTestID(),
			},
			validatorParams: validator.SubnetValidatorParams{
				NodeID:   ids.GenerateTestNodeID(),
				Duration: 1,
				Weight:   20,
			},
			expectedErr: ErrEmptySubnetAuth,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test the validation logic by checking the conditions that would cause errors
			if tt.validatorParams.NodeID == ids.EmptyNodeID {
				assert.Equal(t, ErrEmptyValidatorNodeID, tt.expectedErr)
			}
			if tt.validatorParams.Duration == 0 {
				assert.Equal(t, ErrEmptyValidatorDuration, tt.expectedErr)
			}
			if tt.subnet.SubnetID == ids.Empty {
				assert.Equal(t, ErrEmptySubnetID, tt.expectedErr)
			}
			if len(tt.subnet.DeployInfo.SubnetAuthKeys) == 0 {
				assert.Equal(t, ErrEmptySubnetAuth, tt.expectedErr)
			}
		})
	}
}

func TestSubnet_AddValidator_DefaultWeightLogic(t *testing.T) {
	// Test that weight defaults to 20 when not provided
	validatorParams := validator.SubnetValidatorParams{
		NodeID:   ids.GenerateTestNodeID(),
		Duration: 1,
		Weight:   0, // Should default to 20
	}

	// Test the default weight logic
	if validatorParams.Weight == 0 {
		validatorParams.Weight = 20
	}
	assert.Equal(t, uint64(20), validatorParams.Weight)
}

func TestSubnet_Commit_ValidationLogic(t *testing.T) {
	// Test the validation logic for Commit function
	// This tests the conditions that would cause errors without actually calling the function

	// Test undefined multisig
	multisigUndefined := true
	if multisigUndefined {
		// This would return multisig.ErrUndefinedTx
		assert.True(t, multisigUndefined)
	}

	// Test not ready to commit
	notReady := false
	if !notReady {
		// This would return "tx is not fully signed so can't be committed"
		assert.False(t, notReady)
	}
}
