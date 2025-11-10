// Copyright (C) 2025, Dione Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package subnet

import (
	"encoding/json"
	"math/big"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/DioneProtocol/odysseygo/ids"
	"github.com/DioneProtocol/subnet-evm/commontype"
	"github.com/DioneProtocol/subnet-evm/core"
	"github.com/DioneProtocol/subnet-evm/params"
	"github.com/ethereum/go-ethereum/common"
)

func TestNew_WithGenesisFilePath(t *testing.T) {
	// Create a temporary genesis file
	tempFile, err := os.CreateTemp("", "test_genesis.json")
	require.NoError(t, err)
	defer os.Remove(tempFile.Name())

	genesisContent := `{"config": {"chainId": 123456}}`
	_, err = tempFile.WriteString(genesisContent)
	require.NoError(t, err)
	tempFile.Close()

	params := &SubnetParams{
		GenesisFilePath: tempFile.Name(),
		Name:            "TestSubnet",
	}

	subnet, err := New(params)
	require.NoError(t, err)
	require.NotNil(t, subnet)
	assert.Equal(t, "TestSubnet", subnet.Name)
	assert.Equal(t, []byte(genesisContent), subnet.Genesis)
}

func TestNew_WithSubnetEVMParams(t *testing.T) {
	allocation := core.GenesisAlloc{
		common.HexToAddress("0x1234567890123456789012345678901234567890"): core.GenesisAccount{
			Balance: big.NewInt(1000000000000000000),
		},
	}

	params := &SubnetParams{
		SubnetEVM: &SubnetEVMParams{
			ChainID:     big.NewInt(123456),
			FeeConfig:   commontype.FeeConfig{GasLimit: big.NewInt(8000000)},
			Allocation:  allocation,
			Precompiles: params.Precompiles{},
		},
		Name: "TestSubnet",
	}

	subnet, err := New(params)
	require.NoError(t, err)
	require.NotNil(t, subnet)
	assert.Equal(t, "TestSubnet", subnet.Name)
	assert.NotEmpty(t, subnet.Genesis)
}

func TestNew_ValidationErrors(t *testing.T) {
	tests := []struct {
		name        string
		params      *SubnetParams
		expectedErr string
	}{
		{
			name: "both genesis file and subnet evm params provided",
			params: &SubnetParams{
				GenesisFilePath: "test.json",
				SubnetEVM:       &SubnetEVMParams{},
				Name:            "TestSubnet",
			},
			expectedErr: "genesis file path cannot be non-empty if SubnetEVM params is not empty",
		},
		{
			name: "both genesis file and subnet evm params empty",
			params: &SubnetParams{
				Name: "TestSubnet",
			},
			expectedErr: "genesis file path and SubnetEVM params params cannot all be empty",
		},
		{
			name: "empty name",
			params: &SubnetParams{
				SubnetEVM: &SubnetEVMParams{
					ChainID:     big.NewInt(123456),
					FeeConfig:   commontype.FeeConfig{GasLimit: big.NewInt(8000000)},
					Allocation:  core.GenesisAlloc{},
					Precompiles: params.Precompiles{},
				},
			},
			expectedErr: "SubnetEVM name cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			subnet, err := New(tt.params)
			assert.Nil(t, subnet)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.expectedErr)
		})
	}
}

func TestNew_GenesisFileNotFound(t *testing.T) {
	params := &SubnetParams{
		GenesisFilePath: "nonexistent.json",
		Name:            "TestSubnet",
	}

	subnet, err := New(params)
	assert.Nil(t, subnet)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no such file or directory")
}

func TestNew_InvalidVMName(t *testing.T) {
	params := &SubnetParams{
		SubnetEVM: &SubnetEVMParams{
			ChainID:     big.NewInt(123456),
			FeeConfig:   commontype.FeeConfig{GasLimit: big.NewInt(8000000)},
			Allocation:  core.GenesisAlloc{},
			Precompiles: params.Precompiles{},
		},
		Name: string(make([]byte, 33)), // 33 bytes, exceeds 32 byte limit
	}

	subnet, err := New(params)
	assert.Nil(t, subnet)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "VM name must be <= 32 bytes")
}

func TestCreateEvmGenesis_Success(t *testing.T) {
	allocation := core.GenesisAlloc{
		common.HexToAddress("0x1234567890123456789012345678901234567890"): core.GenesisAccount{
			Balance: big.NewInt(1000000000000000000),
		},
	}

	params := &SubnetEVMParams{
		ChainID:     big.NewInt(123456),
		FeeConfig:   commontype.FeeConfig{GasLimit: big.NewInt(8000000)},
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

func TestCreateEvmGenesis_ValidationErrors(t *testing.T) {
	tests := []struct {
		name        string
		params      *SubnetEVMParams
		expectedErr string
	}{
		{
			name: "nil chain ID",
			params: &SubnetEVMParams{
				FeeConfig:   commontype.FeeConfig{GasLimit: big.NewInt(8000000)},
				Allocation:  core.GenesisAlloc{},
				Precompiles: params.Precompiles{},
			},
			expectedErr: "genesis params chain ID cannot be empty",
		},
		{
			name: "empty fee config",
			params: &SubnetEVMParams{
				ChainID:     big.NewInt(123456),
				FeeConfig:   commontype.EmptyFeeConfig,
				Allocation:  core.GenesisAlloc{},
				Precompiles: params.Precompiles{},
			},
			expectedErr: "genesis params fee config cannot be empty",
		},
		{
			name: "nil allocation",
			params: &SubnetEVMParams{
				ChainID:     big.NewInt(123456),
				FeeConfig:   commontype.FeeConfig{GasLimit: big.NewInt(8000000)},
				Precompiles: params.Precompiles{},
			},
			expectedErr: "genesis params allocation cannot be empty",
		},
		{
			name: "nil precompiles",
			params: &SubnetEVMParams{
				ChainID:    big.NewInt(123456),
				FeeConfig:  commontype.FeeConfig{GasLimit: big.NewInt(8000000)},
				Allocation: core.GenesisAlloc{},
			},
			expectedErr: "genesis params precompiles cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			genesis, err := createEvmGenesis(tt.params)
			assert.Nil(t, genesis)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.expectedErr)
		})
	}
}

func TestVmID(t *testing.T) {
	tests := []struct {
		name        string
		vmName      string
		expectedErr bool
	}{
		{
			name:        "valid vm name",
			vmName:      "TestVM",
			expectedErr: false,
		},
		{
			name:        "empty vm name",
			vmName:      "",
			expectedErr: false,
		},
		{
			name:        "vm name exactly 32 bytes",
			vmName:      string(make([]byte, 32)),
			expectedErr: false,
		},
		{
			name:        "vm name too long",
			vmName:      string(make([]byte, 33)),
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vmID, err := vmID(tt.vmName)
			if tt.expectedErr {
				assert.Error(t, err)
				assert.Equal(t, ids.Empty, vmID)
			} else {
				assert.NoError(t, err)
				// For empty strings and zero-filled strings, the ID might be empty, which is valid
				// We just need to ensure no error is returned
			}
		})
	}
}

func TestSubnet_SetParams(t *testing.T) {
	subnet := &Subnet{}
	controlKeys := []ids.ShortID{ids.GenerateTestShortID(), ids.GenerateTestShortID()}
	subnetAuthKeys := []ids.ShortID{ids.GenerateTestShortID()}
	threshold := uint32(2)

	subnet.SetParams(controlKeys, subnetAuthKeys, threshold)

	assert.Equal(t, controlKeys, subnet.DeployInfo.ControlKeys)
	assert.Equal(t, subnetAuthKeys, subnet.DeployInfo.SubnetAuthKeys)
	assert.Equal(t, threshold, subnet.DeployInfo.Threshold)
}

func TestSubnet_SetSubnetControlParams(t *testing.T) {
	subnet := &Subnet{}
	controlKeys := []ids.ShortID{ids.GenerateTestShortID(), ids.GenerateTestShortID()}
	threshold := uint32(2)

	subnet.SetSubnetControlParams(controlKeys, threshold)

	assert.Equal(t, controlKeys, subnet.DeployInfo.ControlKeys)
	assert.Equal(t, threshold, subnet.DeployInfo.Threshold)
}

func TestSubnet_SetSubnetAuthKeys(t *testing.T) {
	subnet := &Subnet{}
	subnetAuthKeys := []ids.ShortID{ids.GenerateTestShortID(), ids.GenerateTestShortID()}

	subnet.SetSubnetAuthKeys(subnetAuthKeys)

	assert.Equal(t, subnetAuthKeys, subnet.DeployInfo.SubnetAuthKeys)
}

func TestSubnet_SetSubnetID(t *testing.T) {
	subnet := &Subnet{}
	subnetID := ids.GenerateTestID()

	subnet.SetSubnetID(subnetID)

	assert.Equal(t, subnetID, subnet.SubnetID)
}
