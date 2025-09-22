// Copyright (C) 2025, Dione Protocol, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package subnet

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/DioneProtocol/odyssey-tooling-sdk-go/keychain"
	"github.com/DioneProtocol/odyssey-tooling-sdk-go/odyssey"
	"github.com/DioneProtocol/odyssey-tooling-sdk-go/validator"
	"github.com/DioneProtocol/odyssey-tooling-sdk-go/wallet"
	"github.com/DioneProtocol/odysseygo/ids"
	"github.com/DioneProtocol/odysseygo/vms/secp256k1fx"
	"github.com/DioneProtocol/odysseygo/wallet/subnet/primary"
	"github.com/DioneProtocol/subnet-evm/commontype"
	"github.com/DioneProtocol/subnet-evm/core"
	"github.com/DioneProtocol/subnet-evm/params"
)

func createTestWallet(t *testing.T) wallet.Wallet {
	ctx := context.Background()
	network := odyssey.TestnetNetwork()

	// Create a temporary key file for testing
	tempKeyPath := t.TempDir() + "/test.pk"

	// Create keychain using the production approach
	keychain, err := keychain.NewKeychain(network, tempKeyPath, nil)
	require.NoError(t, err)
	require.NotNil(t, keychain)

	// Create wallet using the production approach
	testWallet, err := wallet.New(
		ctx,
		&primary.WalletConfig{
			URI:              network.Endpoint,
			DIONEKeychain:    keychain.Keychain,
			EthKeychain:      secp256k1fx.NewKeychain(),
			OChainTxsToFetch: nil,
		},
	)
	require.NoError(t, err)
	require.NotNil(t, testWallet)

	return testWallet
}

func TestSubnet_CreateSubnetTx_Success(t *testing.T) {
	addTestDelay() // Add delay to avoid rate limiting
	subnet := &Subnet{
		DeployInfo: DeployParams{
			ControlKeys: []ids.ShortID{ids.GenerateTestShortID()},
			Threshold:   1,
		},
	}

	testWallet := createTestWallet(t)
	result, err := subnet.CreateSubnetTx(testWallet)

	require.NoError(t, err)
	require.NotNil(t, result)
}

func TestSubnet_CreateSubnetTx_ValidationErrors(t *testing.T) {
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
			testWallet := createTestWallet(t)
			result, err := tt.subnet.CreateSubnetTx(testWallet)
			assert.Nil(t, result)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.expectedErr)
		})
	}
}

func TestSubnet_CreateBlockchainTx_Success(t *testing.T) {
	addTestDelay() // Add delay to avoid rate limiting
	subnet := &Subnet{
		SubnetID: ids.GenerateTestID(),
		VMID:     ids.GenerateTestID(),
		Name:     "TestChain",
		Genesis:  []byte("test genesis"),
		DeployInfo: DeployParams{
			SubnetAuthKeys: []ids.ShortID{ids.GenerateTestShortID()},
		},
	}

	testWallet := createTestWallet(t)
	result, err := subnet.CreateBlockchainTx(testWallet)

	require.NoError(t, err)
	require.NotNil(t, result)
}

func TestSubnet_CreateBlockchainTx_ValidationErrors(t *testing.T) {
	addTestDelay() // Add delay to avoid rate limiting
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
			testWallet := createTestWallet(t)
			result, err := tt.subnet.CreateBlockchainTx(testWallet)
			assert.Nil(t, result)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.expectedErr)
		})
	}
}

func TestSubnet_AddValidator_Success(t *testing.T) {
	addTestDelay() // Add delay to avoid rate limiting
	subnet := &Subnet{
		SubnetID: ids.GenerateTestID(),
		DeployInfo: DeployParams{
			SubnetAuthKeys: []ids.ShortID{ids.GenerateTestShortID()},
		},
	}

	validatorParams := validator.SubnetValidatorParams{
		NodeID:   ids.GenerateTestNodeID(),
		Duration: time.Hour,
		Weight:   20,
	}

	testWallet := createTestWallet(t)
	result, err := subnet.AddValidator(testWallet, validatorParams)

	require.NoError(t, err)
	require.NotNil(t, result)
}

func TestSubnet_AddValidator_DefaultWeight(t *testing.T) {
	addTestDelay() // Add delay to avoid rate limiting
	subnet := &Subnet{
		SubnetID: ids.GenerateTestID(),
		DeployInfo: DeployParams{
			SubnetAuthKeys: []ids.ShortID{ids.GenerateTestShortID()},
		},
	}

	validatorParams := validator.SubnetValidatorParams{
		NodeID:   ids.GenerateTestNodeID(),
		Duration: time.Hour,
		Weight:   0, // Should default to 20
	}

	testWallet := createTestWallet(t)
	result, err := subnet.AddValidator(testWallet, validatorParams)

	require.NoError(t, err)
	require.NotNil(t, result)
}

func TestSubnet_AddValidator_ValidationErrors(t *testing.T) {
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
				Duration: time.Hour,
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
				Duration: time.Hour,
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
				Duration: time.Hour,
				Weight:   20,
			},
			expectedErr: ErrEmptySubnetAuth,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testWallet := createTestWallet(t)
			result, err := tt.subnet.AddValidator(testWallet, tt.validatorParams)
			assert.Nil(t, result)
			assert.Error(t, err)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

// Integration test that creates a complete subnet workflow
func TestSubnet_CompleteWorkflow(t *testing.T) {
	addTestDelay() // Add delay to avoid rate limiting
	// Create a subnet with proper parameters
	subnetParams := &SubnetParams{
		SubnetEVM: &SubnetEVMParams{
			ChainID:     big.NewInt(123456),
			FeeConfig:   commontype.FeeConfig{GasLimit: big.NewInt(8000000)},
			Allocation:  core.GenesisAlloc{},
			Precompiles: params.Precompiles{},
		},
		Name: "TestSubnet",
	}

	subnet, err := New(subnetParams)
	require.NoError(t, err)
	require.NotNil(t, subnet)

	// Set up control keys and threshold
	controlKeys := []ids.ShortID{ids.GenerateTestShortID()}
	subnet.SetSubnetControlParams(controlKeys, 1)

	// Create wallet
	testWallet := createTestWallet(t)

	// Test CreateSubnetTx
	createSubnetTx, err := subnet.CreateSubnetTx(testWallet)
	require.NoError(t, err)
	require.NotNil(t, createSubnetTx)

	// Set subnet ID (simulating successful subnet creation)
	subnetID := ids.GenerateTestID()
	subnet.SetSubnetID(subnetID)

	// Set subnet auth keys
	subnetAuthKeys := []ids.ShortID{ids.GenerateTestShortID()}
	subnet.SetSubnetAuthKeys(subnetAuthKeys)

	// Test CreateBlockchainTx
	createBlockchainTx, err := subnet.CreateBlockchainTx(testWallet)
	require.NoError(t, err)
	require.NotNil(t, createBlockchainTx)

	// Test AddValidator
	validatorParams := validator.SubnetValidatorParams{
		NodeID:   ids.GenerateTestNodeID(),
		Duration: time.Hour,
		Weight:   20,
	}

	addValidatorTx, err := subnet.AddValidator(testWallet, validatorParams)
	require.NoError(t, err)
	require.NotNil(t, addValidatorTx)
}
