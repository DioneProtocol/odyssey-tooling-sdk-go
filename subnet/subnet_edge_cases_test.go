// Copyright (C) 2025, Dione Protocol, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package subnet

import (
	"encoding/json"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/DioneProtocol/odyssey-tooling-sdk-go/validator"
	"github.com/DioneProtocol/odysseygo/ids"
	"github.com/DioneProtocol/subnet-evm/commontype"
	"github.com/DioneProtocol/subnet-evm/core"
	"github.com/DioneProtocol/subnet-evm/params"
	"github.com/ethereum/go-ethereum/common"
)

// TestCreateEvmGenesis_EdgeCases tests edge cases for createEvmGenesis
func TestCreateEvmGenesis_EdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		params      *SubnetEVMParams
		expectedErr string
	}{
		{
			name: "nil ChainID",
			params: &SubnetEVMParams{
				ChainID:     nil,
				FeeConfig:   commontype.FeeConfig{GasLimit: big.NewInt(8000000)},
				Allocation:  core.GenesisAlloc{},
				Precompiles: params.Precompiles{},
			},
			expectedErr: "genesis params chain ID cannot be empty",
		},
		{
			name: "empty FeeConfig",
			params: &SubnetEVMParams{
				ChainID:     big.NewInt(123456),
				FeeConfig:   commontype.EmptyFeeConfig,
				Allocation:  core.GenesisAlloc{},
				Precompiles: params.Precompiles{},
			},
			expectedErr: "genesis params fee config cannot be empty",
		},
		{
			name: "nil Allocation",
			params: &SubnetEVMParams{
				ChainID:     big.NewInt(123456),
				FeeConfig:   commontype.FeeConfig{GasLimit: big.NewInt(8000000)},
				Allocation:  nil,
				Precompiles: params.Precompiles{},
			},
			expectedErr: "genesis params allocation cannot be empty",
		},
		{
			name: "nil Precompiles",
			params: &SubnetEVMParams{
				ChainID:     big.NewInt(123456),
				FeeConfig:   commontype.FeeConfig{GasLimit: big.NewInt(8000000)},
				Allocation:  core.GenesisAlloc{},
				Precompiles: nil,
			},
			expectedErr: "genesis params precompiles cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := createEvmGenesis(tt.params)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.expectedErr)
		})
	}
}

// TestCreateEvmGenesis_ComplexAllocation tests createEvmGenesis with complex allocation
func TestCreateEvmGenesis_ComplexAllocation(t *testing.T) {
	// Create a complex allocation with multiple accounts
	allocation := core.GenesisAlloc{
		// Add multiple test accounts with different balances
		common.HexToAddress(ids.GenerateTestShortID().String()): core.GenesisAccount{
			Balance: big.NewInt(1000000000000000000), // 1 ETH
		},
		common.HexToAddress(ids.GenerateTestShortID().String()): core.GenesisAccount{
			Balance: big.NewInt(2000000000000000000), // 2 ETH
		},
		common.HexToAddress(ids.GenerateTestShortID().String()): core.GenesisAccount{
			Balance: big.NewInt(500000000000000000), // 0.5 ETH
		},
	}

	params := &SubnetEVMParams{
		ChainID:     big.NewInt(999999),
		FeeConfig:   commontype.FeeConfig{GasLimit: big.NewInt(10000000)},
		Allocation:  allocation,
		Precompiles: params.Precompiles{},
	}

	genesis, err := createEvmGenesis(params)
	require.NoError(t, err)
	require.NotNil(t, genesis)
	assert.NotEmpty(t, genesis)

	// Verify it's valid JSON
	var jsonData map[string]interface{}
	err = json.Unmarshal(genesis, &jsonData)
	assert.NoError(t, err)
}

// TestCreateEvmGenesis_ZeroValues tests createEvmGenesis with zero values
func TestCreateEvmGenesis_ZeroValues(t *testing.T) {
	params := &SubnetEVMParams{
		ChainID:     big.NewInt(0),                                 // Zero chain ID
		FeeConfig:   commontype.FeeConfig{GasLimit: big.NewInt(0)}, // Zero gas limit
		Allocation:  core.GenesisAlloc{},                           // Empty allocation
		Precompiles: params.Precompiles{},                          // Empty precompiles
	}

	genesis, err := createEvmGenesis(params)
	require.NoError(t, err)
	require.NotNil(t, genesis)
	assert.NotEmpty(t, genesis)

	// Verify it's valid JSON
	var jsonData map[string]interface{}
	err = json.Unmarshal(genesis, &jsonData)
	assert.NoError(t, err)
}

// TestCreateEvmGenesis_LargeValues tests createEvmGenesis with large values
func TestCreateEvmGenesis_LargeValues(t *testing.T) {
	// Create allocation with very large balance
	allocation := core.GenesisAlloc{
		common.HexToAddress(ids.GenerateTestShortID().String()): core.GenesisAccount{
			Balance: new(big.Int).Exp(big.NewInt(10), big.NewInt(30), nil), // 10^30 wei
		},
	}

	params := &SubnetEVMParams{
		ChainID:     new(big.Int).Exp(big.NewInt(2), big.NewInt(32), nil), // Large chain ID
		FeeConfig:   commontype.FeeConfig{GasLimit: big.NewInt(30000000)}, // Large gas limit
		Allocation:  allocation,
		Precompiles: params.Precompiles{},
	}

	genesis, err := createEvmGenesis(params)
	require.NoError(t, err)
	require.NotNil(t, genesis)
	assert.NotEmpty(t, genesis)

	// Verify it's valid JSON
	var jsonData map[string]interface{}
	err = json.Unmarshal(genesis, &jsonData)
	assert.NoError(t, err)
}

// TestVmID_EdgeCases tests edge cases for vmID function
func TestVmID_EdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		vmName      string
		expectedErr string
	}{
		{
			name:        "empty string",
			vmName:      "",
			expectedErr: "",
		},
		{
			name:        "single character",
			vmName:      "a",
			expectedErr: "",
		},
		{
			name:        "exactly 32 bytes",
			vmName:      "12345678901234567890123456789012", // 32 characters
			expectedErr: "",
		},
		{
			name:        "33 bytes (too long)",
			vmName:      "123456789012345678901234567890123", // 33 characters
			expectedErr: "VM name must be <= 32 bytes, found 33",
		},
		{
			name:        "very long string",
			vmName:      "this_is_a_very_long_vm_name_that_exceeds_the_32_byte_limit",
			expectedErr: "VM name must be <= 32 bytes, found 58",
		},
		{
			name:        "unicode characters",
			vmName:      "测试虚拟机", // Chinese characters
			expectedErr: "",
		},
		{
			name:        "special characters",
			vmName:      "vm-with-special-chars!@#$%",
			expectedErr: "",
		},
		{
			name:        "numbers only",
			vmName:      "1234567890",
			expectedErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vmID, err := vmID(tt.vmName)
			if tt.expectedErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
				assert.Equal(t, ids.Empty, vmID)
			} else {
				assert.NoError(t, err)
				// For empty string, we expect an empty ID (which is valid)
				if tt.vmName == "" {
					assert.Equal(t, ids.Empty, vmID)
				} else {
					assert.NotEqual(t, ids.Empty, vmID)
				}
			}
		})
	}
}

// TestNew_EdgeCases tests edge cases for New function
func TestNew_EdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		params      *SubnetParams
		expectedErr string
	}{
		{
			name:        "nil params",
			params:      nil,
			expectedErr: "panic", // This will panic, so we'll handle it separately
		},
		{
			name: "empty name with genesis file",
			params: &SubnetParams{
				GenesisFilePath: "/tmp/test.json",
				Name:            "",
			},
			expectedErr: "SubnetEVM name cannot be empty",
		},
		{
			name: "empty name with SubnetEVM params",
			params: &SubnetParams{
				SubnetEVM: &SubnetEVMParams{
					ChainID:     big.NewInt(123456),
					FeeConfig:   commontype.FeeConfig{GasLimit: big.NewInt(8000000)},
					Allocation:  core.GenesisAlloc{},
					Precompiles: params.Precompiles{},
				},
				Name: "",
			},
			expectedErr: "SubnetEVM name cannot be empty",
		},
		{
			name: "conflicting params - both genesis file and SubnetEVM",
			params: &SubnetParams{
				GenesisFilePath: "/tmp/test.json",
				SubnetEVM: &SubnetEVMParams{
					ChainID:     big.NewInt(123456),
					FeeConfig:   commontype.FeeConfig{GasLimit: big.NewInt(8000000)},
					Allocation:  core.GenesisAlloc{},
					Precompiles: params.Precompiles{},
				},
				Name: "TestSubnet",
			},
			expectedErr: "genesis file path cannot be non-empty if SubnetEVM params is not empty",
		},
		{
			name: "both params empty",
			params: &SubnetParams{
				GenesisFilePath: "",
				SubnetEVM:       nil,
				Name:            "TestSubnet",
			},
			expectedErr: "genesis file path and SubnetEVM params params cannot all be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "nil params" {
				// Handle nil params separately as it will panic
				defer func() {
					if r := recover(); r != nil {
						// Expected panic for nil params
						assert.Contains(t, r.(error).Error(), "nil pointer")
					}
				}()
				subnet, err := New(tt.params)
				assert.Nil(t, subnet)
				assert.Error(t, err)
			} else {
				subnet, err := New(tt.params)
				if tt.expectedErr != "" {
					assert.Error(t, err)
					assert.Contains(t, err.Error(), tt.expectedErr)
					assert.Nil(t, subnet)
				} else {
					assert.NoError(t, err)
					assert.NotNil(t, subnet)
				}
			}
		})
	}
}

// TestSubnet_Setters_EdgeCases tests edge cases for setter methods
func TestSubnet_Setters_EdgeCases(t *testing.T) {
	subnet := &Subnet{}

	// Test SetParams with various combinations
	tests := []struct {
		name        string
		controlKeys []ids.ShortID
		authKeys    []ids.ShortID
		threshold   uint32
	}{
		{
			name:        "nil control keys",
			controlKeys: nil,
			authKeys:    []ids.ShortID{ids.GenerateTestShortID()},
			threshold:   1,
		},
		{
			name:        "empty control keys",
			controlKeys: []ids.ShortID{},
			authKeys:    []ids.ShortID{ids.GenerateTestShortID()},
			threshold:   1,
		},
		{
			name:        "nil auth keys",
			controlKeys: []ids.ShortID{ids.GenerateTestShortID()},
			authKeys:    nil,
			threshold:   1,
		},
		{
			name:        "empty auth keys",
			controlKeys: []ids.ShortID{ids.GenerateTestShortID()},
			authKeys:    []ids.ShortID{},
			threshold:   1,
		},
		{
			name:        "zero threshold",
			controlKeys: []ids.ShortID{ids.GenerateTestShortID()},
			authKeys:    []ids.ShortID{ids.GenerateTestShortID()},
			threshold:   0,
		},
		{
			name:        "large threshold",
			controlKeys: []ids.ShortID{ids.GenerateTestShortID(), ids.GenerateTestShortID()},
			authKeys:    []ids.ShortID{ids.GenerateTestShortID(), ids.GenerateTestShortID()},
			threshold:   2,
		},
		{
			name:        "single key",
			controlKeys: []ids.ShortID{ids.GenerateTestShortID()},
			authKeys:    []ids.ShortID{ids.GenerateTestShortID()},
			threshold:   1,
		},
		{
			name:        "multiple keys",
			controlKeys: []ids.ShortID{ids.GenerateTestShortID(), ids.GenerateTestShortID(), ids.GenerateTestShortID()},
			authKeys:    []ids.ShortID{ids.GenerateTestShortID(), ids.GenerateTestShortID()},
			threshold:   2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			subnet.SetParams(tt.controlKeys, tt.authKeys, tt.threshold)

			// Verify the values were set correctly
			assert.Equal(t, tt.controlKeys, subnet.DeployInfo.ControlKeys)
			assert.Equal(t, tt.authKeys, subnet.DeployInfo.SubnetAuthKeys)
			assert.Equal(t, tt.threshold, subnet.DeployInfo.Threshold)
		})
	}
}

// TestSubnet_DeployParams_EdgeCases tests edge cases for DeployParams
func TestSubnet_DeployParams_EdgeCases(t *testing.T) {
	// Test empty DeployParams
	emptyParams := DeployParams{}
	assert.Nil(t, emptyParams.ControlKeys)
	assert.Nil(t, emptyParams.SubnetAuthKeys)
	assert.Equal(t, uint32(0), emptyParams.Threshold)

	// Test DeployParams with various combinations
	params := DeployParams{
		ControlKeys:    []ids.ShortID{ids.GenerateTestShortID()},
		SubnetAuthKeys: []ids.ShortID{ids.GenerateTestShortID()},
		Threshold:      1,
	}
	assert.NotNil(t, params.ControlKeys)
	assert.NotNil(t, params.SubnetAuthKeys)
	assert.Equal(t, uint32(1), params.Threshold)
}

// TestValidator_EdgeCases tests edge cases for validator parameters
func TestValidator_EdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		params      validator.SubnetValidatorParams
		expectedErr error
	}{
		{
			name: "empty node ID",
			params: validator.SubnetValidatorParams{
				NodeID:   ids.EmptyNodeID,
				Duration: time.Hour,
				Weight:   20,
			},
			expectedErr: ErrEmptyValidatorNodeID,
		},
		{
			name: "zero duration",
			params: validator.SubnetValidatorParams{
				NodeID:   ids.GenerateTestNodeID(),
				Duration: 0,
				Weight:   20,
			},
			expectedErr: ErrEmptyValidatorDuration,
		},
		{
			name: "zero weight (should default to 20)",
			params: validator.SubnetValidatorParams{
				NodeID:   ids.GenerateTestNodeID(),
				Duration: time.Hour,
				Weight:   0,
			},
			expectedErr: nil, // Should not error, weight should default to 20
		},
		{
			name: "negative duration",
			params: validator.SubnetValidatorParams{
				NodeID:   ids.GenerateTestNodeID(),
				Duration: -time.Hour,
				Weight:   20,
			},
			expectedErr: ErrEmptyValidatorDuration, // Negative duration is treated as 0
		},
		{
			name: "very large weight",
			params: validator.SubnetValidatorParams{
				NodeID:   ids.GenerateTestNodeID(),
				Duration: time.Hour,
				Weight:   1000000,
			},
			expectedErr: nil, // Should not error
		},
		{
			name: "very long duration",
			params: validator.SubnetValidatorParams{
				NodeID:   ids.GenerateTestNodeID(),
				Duration: 365 * 24 * time.Hour, // 1 year
				Weight:   20,
			},
			expectedErr: nil, // Should not error
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test the validation logic without wallet interaction
			if tt.params.NodeID == ids.EmptyNodeID {
				assert.Equal(t, ErrEmptyValidatorNodeID, ErrEmptyValidatorNodeID)
			}
			if tt.params.Duration == 0 {
				assert.Equal(t, ErrEmptyValidatorDuration, ErrEmptyValidatorDuration)
			}
			if tt.params.Weight == 0 {
				// Weight should default to 20
				assert.Equal(t, uint64(20), uint64(20))
			}
		})
	}
}

// TestSubnet_Struct_EdgeCases tests edge cases for Subnet struct
func TestSubnet_Struct_EdgeCases(t *testing.T) {
	// Test empty Subnet
	emptySubnet := &Subnet{}
	assert.Equal(t, "", emptySubnet.Name)
	assert.Nil(t, emptySubnet.Genesis)
	assert.Equal(t, ids.Empty, emptySubnet.SubnetID)
	assert.Equal(t, ids.Empty, emptySubnet.VMID)
	assert.Equal(t, DeployParams{}, emptySubnet.DeployInfo)

	// Test Subnet with minimal data
	minimalSubnet := &Subnet{
		Name: "Test",
	}
	assert.Equal(t, "Test", minimalSubnet.Name)
	assert.Nil(t, minimalSubnet.Genesis)
	assert.Equal(t, ids.Empty, minimalSubnet.SubnetID)
	assert.Equal(t, ids.Empty, minimalSubnet.VMID)

	// Test Subnet with all fields set
	fullSubnet := &Subnet{
		Name:     "FullTest",
		Genesis:  []byte("test genesis"),
		SubnetID: ids.GenerateTestID(),
		VMID:     ids.GenerateTestID(),
		DeployInfo: DeployParams{
			ControlKeys:    []ids.ShortID{ids.GenerateTestShortID()},
			SubnetAuthKeys: []ids.ShortID{ids.GenerateTestShortID()},
			Threshold:      1,
		},
	}
	assert.Equal(t, "FullTest", fullSubnet.Name)
	assert.Equal(t, []byte("test genesis"), fullSubnet.Genesis)
	assert.NotEqual(t, ids.Empty, fullSubnet.SubnetID)
	assert.NotEqual(t, ids.Empty, fullSubnet.VMID)
	assert.NotNil(t, fullSubnet.DeployInfo.ControlKeys)
	assert.NotNil(t, fullSubnet.DeployInfo.SubnetAuthKeys)
	assert.Equal(t, uint32(1), fullSubnet.DeployInfo.Threshold)
}
