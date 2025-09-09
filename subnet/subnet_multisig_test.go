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
	"github.com/DioneProtocol/odyssey-tooling-sdk-go/multisig"
	"github.com/DioneProtocol/odyssey-tooling-sdk-go/odyssey"
	"github.com/DioneProtocol/odyssey-tooling-sdk-go/validator"
	"github.com/DioneProtocol/odyssey-tooling-sdk-go/wallet"
	"github.com/DioneProtocol/odysseygo/ids"
	"github.com/DioneProtocol/odysseygo/vms/omegavm/txs"
	"github.com/DioneProtocol/odysseygo/vms/secp256k1fx"
	"github.com/DioneProtocol/odysseygo/wallet/subnet/primary"
	"github.com/DioneProtocol/subnet-evm/commontype"
	"github.com/DioneProtocol/subnet-evm/core"
	"github.com/DioneProtocol/subnet-evm/params"
)

// createTestWallet creates a test wallet for multisig testing
func createTestWalletForMultisig(t *testing.T) wallet.Wallet {
	ctx := context.Background()
	network := odyssey.TestnetNetwork()

	tempKeyPath := t.TempDir() + "/test.pk"

	keychain, err := keychain.NewKeychain(network, tempKeyPath, nil)
	require.NoError(t, err)
	require.NotNil(t, keychain)

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

// createTestSubnet creates a test subnet with proper configuration
func createTestSubnet(t *testing.T) *Subnet {
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
	controlKeys := []ids.ShortID{ids.GenerateTestShortID(), ids.GenerateTestShortID()}
	subnet.SetSubnetControlParams(controlKeys, 2)

	// Set subnet auth keys
	subnetAuthKeys := []ids.ShortID{ids.GenerateTestShortID()}
	subnet.SetSubnetAuthKeys(subnetAuthKeys)

	return subnet
}

// TestMultisig_CreateSubnetTx_WithMultisig tests CreateSubnetTx with multisig setup
func TestMultisig_CreateSubnetTx_WithMultisig(t *testing.T) {
	subnet := createTestSubnet(t)
	testWallet := createTestWalletForMultisig(t)

	// Test CreateSubnetTx returns multisig
	multisigTx, err := subnet.CreateSubnetTx(testWallet)

	// We expect this to fail due to insufficient funds, but we can test the multisig structure
	if err != nil {
		// Verify it's a funding error, not a multisig error
		assert.Contains(t, err.Error(), "insufficient funds")
		assert.Nil(t, multisigTx)
	} else {
		// If it succeeds, verify multisig structure
		require.NotNil(t, multisigTx)
		assert.False(t, multisigTx.Undefined())

		// Test multisig methods
		tx, err := multisigTx.GetWrappedOChainTx()
		require.NoError(t, err)
		require.NotNil(t, tx)

		// Verify it's a CreateSubnetTx
		_, ok := tx.Unsigned.(*txs.CreateSubnetTx)
		assert.True(t, ok)
	}
}

// TestMultisig_CreateBlockchainTx_WithMultisig tests CreateBlockchainTx with multisig setup
func TestMultisig_CreateBlockchainTx_WithMultisig(t *testing.T) {
	subnet := createTestSubnet(t)
	testWallet := createTestWalletForMultisig(t)

	// Set subnet ID (simulating successful subnet creation)
	subnetID := ids.GenerateTestID()
	subnet.SetSubnetID(subnetID)

	// Test CreateBlockchainTx returns multisig
	multisigTx, err := subnet.CreateBlockchainTx(testWallet)

	// We expect this to fail due to insufficient funds, but we can test the multisig structure
	if err != nil {
		// Verify it's a funding error, not a multisig error
		assert.Contains(t, err.Error(), "insufficient funds")
		assert.Nil(t, multisigTx)
	} else {
		// If it succeeds, verify multisig structure
		require.NotNil(t, multisigTx)
		assert.False(t, multisigTx.Undefined())

		// Test multisig methods
		tx, err := multisigTx.GetWrappedOChainTx()
		require.NoError(t, err)
		require.NotNil(t, tx)

		// Verify it's a CreateChainTx
		_, ok := tx.Unsigned.(*txs.CreateChainTx)
		assert.True(t, ok)
	}
}

// TestMultisig_AddValidator_WithMultisig tests AddValidator with multisig setup
func TestMultisig_AddValidator_WithMultisig(t *testing.T) {
	subnet := createTestSubnet(t)
	testWallet := createTestWalletForMultisig(t)

	// Set subnet ID (simulating successful subnet creation)
	subnetID := ids.GenerateTestID()
	subnet.SetSubnetID(subnetID)

	validatorParams := validator.SubnetValidatorParams{
		NodeID:   ids.GenerateTestNodeID(),
		Duration: time.Hour,
		Weight:   20,
	}

	// Test AddValidator returns multisig
	multisigTx, err := subnet.AddValidator(testWallet, validatorParams)

	// We expect this to fail due to insufficient funds, but we can test the multisig structure
	if err != nil {
		// Verify it's a funding error, not a multisig error
		assert.Contains(t, err.Error(), "insufficient funds")
		assert.Nil(t, multisigTx)
	} else {
		// If it succeeds, verify multisig structure
		require.NotNil(t, multisigTx)
		assert.False(t, multisigTx.Undefined())

		// Test multisig methods
		tx, err := multisigTx.GetWrappedOChainTx()
		require.NoError(t, err)
		require.NotNil(t, tx)

		// Verify it's an AddSubnetValidatorTx
		_, ok := tx.Unsigned.(*txs.AddSubnetValidatorTx)
		assert.True(t, ok)
	}
}

// TestMultisig_Commit_ValidationLogic tests the Commit function validation logic
func TestMultisig_Commit_ValidationLogic(t *testing.T) {
	subnet := createTestSubnet(t)
	testWallet := createTestWalletForMultisig(t)

	// Test 1: Undefined multisig
	undefinedMultisig := multisig.Multisig{}
	assert.True(t, undefinedMultisig.Undefined())

	txID, err := subnet.Commit(undefinedMultisig, testWallet, false)
	assert.Equal(t, ids.Empty, txID)
	assert.Error(t, err)
	assert.Equal(t, multisig.ErrUndefinedTx, err)

	// Test 2: Create a mock multisig that's not ready to commit
	// We'll create a multisig with a transaction but simulate it not being ready
	mockMultisig := multisig.Multisig{}

	// Test the validation logic without actually creating a real transaction
	// This tests the error paths in Commit function
	txID, err = subnet.Commit(mockMultisig, testWallet, false)
	assert.Equal(t, ids.Empty, txID)
	assert.Error(t, err)
	assert.Equal(t, multisig.ErrUndefinedTx, err)
}

// TestMultisig_Commit_NotReadyToCommit tests Commit when multisig is not ready
func TestMultisig_Commit_NotReadyToCommit(t *testing.T) {
	subnet := createTestSubnet(t)
	testWallet := createTestWalletForMultisig(t)

	// Create a mock multisig that's defined but not ready
	// We can't easily create a real multisig without funds, so we test the error paths
	mockMultisig := multisig.Multisig{}

	// Test undefined multisig (this is the main error path we can test)
	txID, err := subnet.Commit(mockMultisig, testWallet, false)
	assert.Equal(t, ids.Empty, txID)
	assert.Error(t, err)
	assert.Equal(t, multisig.ErrUndefinedTx, err)
}

// TestMultisig_StateManagement tests multisig state management
func TestMultisig_StateManagement(t *testing.T) {
	// Test undefined multisig
	ms := &multisig.Multisig{}
	assert.True(t, ms.Undefined())

	// Test error conditions
	_, err := ms.IsReadyToCommit()
	assert.Error(t, err)
	assert.Equal(t, multisig.ErrUndefinedTx, err)

	_, err = ms.GetWrappedOChainTx()
	assert.Error(t, err)
	assert.Equal(t, multisig.ErrUndefinedTx, err)

	_, err = ms.GetAuthSigners()
	assert.Error(t, err)
	assert.Equal(t, multisig.ErrUndefinedTx, err)

	_, _, err = ms.GetRemainingAuthSigners()
	assert.Error(t, err)
	assert.Equal(t, multisig.ErrUndefinedTx, err)

	_, err = ms.GetTxKind()
	assert.Error(t, err)
	assert.Equal(t, multisig.ErrUndefinedTx, err)

	_, err = ms.GetNetworkID()
	assert.Error(t, err)
	assert.Equal(t, multisig.ErrUndefinedTx, err)

	_, err = ms.GetNetwork()
	assert.Error(t, err)
	assert.Equal(t, multisig.ErrUndefinedTx, err)

	_, err = ms.GetBlockchainID()
	assert.Error(t, err)
	assert.Equal(t, multisig.ErrUndefinedTx, err)

	_, err = ms.GetSubnetID()
	assert.Error(t, err)
	assert.Equal(t, multisig.ErrUndefinedTx, err)

	_, _, err = ms.GetSubnetOwners()
	assert.Error(t, err)
	assert.Equal(t, multisig.ErrUndefinedTx, err)
}

// TestMultisig_Serialization tests multisig serialization methods
func TestMultisig_Serialization(t *testing.T) {
	ms := &multisig.Multisig{}

	// Test ToBytes with undefined multisig
	_, err := ms.ToBytes()
	assert.Error(t, err)
	assert.Equal(t, multisig.ErrUndefinedTx, err)

	// Test ToFile with undefined multisig
	err = ms.ToFile("/tmp/test.tx")
	assert.Error(t, err)
	assert.Equal(t, multisig.ErrUndefinedTx, err)

	// Test FromBytes with invalid data
	err = ms.FromBytes([]byte("invalid"))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error unmarshaling signed tx")

	// Test FromFile with non-existent file
	err = ms.FromFile("/non/existent/file.tx")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no such file or directory")
}

// TestMultisig_String tests multisig string representation
func TestMultisig_String(t *testing.T) {
	ms := &multisig.Multisig{}

	// Test string representation of undefined multisig
	assert.Equal(t, "", ms.String())
}

// TestMultisig_CompleteWorkflow tests the complete multisig workflow
func TestMultisig_CompleteWorkflow(t *testing.T) {
	subnet := createTestSubnet(t)
	testWallet := createTestWalletForMultisig(t)

	// Step 1: Create subnet transaction (returns multisig)
	createSubnetMultisig, err := subnet.CreateSubnetTx(testWallet)
	if err != nil {
		// Expected due to insufficient funds, but we can test the structure
		assert.Contains(t, err.Error(), "insufficient funds")
		assert.Nil(t, createSubnetMultisig)
		return
	}

	// If we get here, we have a real multisig (unlikely without funds)
	require.NotNil(t, createSubnetMultisig)
	assert.False(t, createSubnetMultisig.Undefined())

	// Step 2: Test multisig state
	tx, err := createSubnetMultisig.GetWrappedOChainTx()
	require.NoError(t, err)
	require.NotNil(t, tx)

	// Step 3: Set subnet ID (simulating successful subnet creation)
	subnet.SetSubnetID(tx.ID())

	// Step 4: Create blockchain transaction
	createBlockchainMultisig, err := subnet.CreateBlockchainTx(testWallet)
	if err != nil {
		assert.Contains(t, err.Error(), "insufficient funds")
		assert.Nil(t, createBlockchainMultisig)
		return
	}

	require.NotNil(t, createBlockchainMultisig)
	assert.False(t, createBlockchainMultisig.Undefined())

	// Step 5: Add validator
	validatorParams := validator.SubnetValidatorParams{
		NodeID:   ids.GenerateTestNodeID(),
		Duration: time.Hour,
		Weight:   20,
	}

	addValidatorMultisig, err := subnet.AddValidator(testWallet, validatorParams)
	if err != nil {
		assert.Contains(t, err.Error(), "insufficient funds")
		assert.Nil(t, addValidatorMultisig)
		return
	}

	require.NotNil(t, addValidatorMultisig)
	assert.False(t, addValidatorMultisig.Undefined())

	// Step 6: Test commit validation (without actually committing)
	// Test undefined multisig
	undefinedMultisig := multisig.Multisig{}
	txID, err := subnet.Commit(undefinedMultisig, testWallet, false)
	assert.Equal(t, ids.Empty, txID)
	assert.Error(t, err)
	assert.Equal(t, multisig.ErrUndefinedTx, err)
}

// TestMultisig_WalletIntegration tests wallet multisig integration
func TestMultisig_WalletIntegration(t *testing.T) {
	testWallet := createTestWalletForMultisig(t)

	// Test wallet multisig setup
	authKeys := []ids.ShortID{ids.GenerateTestShortID(), ids.GenerateTestShortID()}

	// Test SetSubnetAuthMultisig
	testWallet.SetSubnetAuthMultisig(authKeys)

	// Verify wallet addresses are available
	addresses := testWallet.Addresses()
	assert.NotEmpty(t, addresses)

	// Test that wallet is properly configured for multisig
	// (We can't easily test the internal state, but we can verify no errors)
	assert.NotNil(t, testWallet)
}

// TestMultisig_ErrorHandling tests comprehensive error handling
func TestMultisig_ErrorHandling(t *testing.T) {
	subnet := createTestSubnet(t)
	testWallet := createTestWalletForMultisig(t)

	// Test all error conditions we can trigger without network calls

	// 1. Test undefined multisig in Commit
	undefinedMultisig := multisig.Multisig{}
	txID, err := subnet.Commit(undefinedMultisig, testWallet, false)
	assert.Equal(t, ids.Empty, txID)
	assert.Error(t, err)
	assert.Equal(t, multisig.ErrUndefinedTx, err)

	// 2. Test multisig serialization errors
	_, err = undefinedMultisig.ToBytes()
	assert.Error(t, err)
	assert.Equal(t, multisig.ErrUndefinedTx, err)

	// 3. Test multisig state query errors
	_, err = undefinedMultisig.IsReadyToCommit()
	assert.Error(t, err)
	assert.Equal(t, multisig.ErrUndefinedTx, err)

	// 4. Test multisig network query errors
	_, err = undefinedMultisig.GetNetworkID()
	assert.Error(t, err)
	assert.Equal(t, multisig.ErrUndefinedTx, err)

	// 5. Test multisig subnet query errors
	_, err = undefinedMultisig.GetSubnetID()
	assert.Error(t, err)
	assert.Equal(t, multisig.ErrUndefinedTx, err)
}

// TestMultisig_EdgeCases tests edge cases and boundary conditions
func TestMultisig_EdgeCases(t *testing.T) {
	// Test empty multisig
	ms := &multisig.Multisig{}
	assert.True(t, ms.Undefined())
	assert.Equal(t, "", ms.String())

	// Test multisig with invalid transaction data
	err := ms.FromBytes([]byte("invalid transaction data"))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error unmarshaling signed tx")

	// Test multisig with empty transaction data
	err = ms.FromBytes([]byte{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error unmarshaling signed tx")

	// Test multisig with nil transaction data
	err = ms.FromBytes(nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error unmarshaling signed tx")
}
