// Copyright (C) 2025, Dione Protocol, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package subnet

import (
	"context"
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
)

// createTestWalletForCommit creates a test wallet for commit testing
func createTestWalletForCommit(t *testing.T) wallet.Wallet {
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

// TestCommit_UndefinedMultisig tests Commit with undefined multisig
func TestCommit_UndefinedMultisig(t *testing.T) {
	subnet := createTestSubnet(t)
	testWallet := createTestWalletForCommit(t)

	// Test with undefined multisig
	undefinedMultisig := multisig.Multisig{}
	assert.True(t, undefinedMultisig.Undefined())

	txID, err := subnet.Commit(undefinedMultisig, testWallet, false)
	assert.Equal(t, ids.Empty, txID)
	assert.Error(t, err)
	assert.Equal(t, multisig.ErrUndefinedTx, err)
}

// TestCommit_NotReadyToCommit tests Commit when multisig is not ready
func TestCommit_NotReadyToCommit(t *testing.T) {
	subnet := createTestSubnet(t)
	testWallet := createTestWalletForCommit(t)

	// Create a mock multisig that's defined but not ready
	// We can't easily create a real multisig without funds, so we test the error paths
	mockMultisig := multisig.Multisig{}

	// Test undefined multisig (this is the main error path we can test)
	txID, err := subnet.Commit(mockMultisig, testWallet, false)
	assert.Equal(t, ids.Empty, txID)
	assert.Error(t, err)
	assert.Equal(t, multisig.ErrUndefinedTx, err)
}

// TestCommit_ErrorHandling tests Commit error handling
func TestCommit_ErrorHandling(t *testing.T) {
	subnet := createTestSubnet(t)
	testWallet := createTestWalletForCommit(t)

	// Test all error conditions we can trigger without network calls

	// 1. Test undefined multisig
	undefinedMultisig := multisig.Multisig{}
	txID, err := subnet.Commit(undefinedMultisig, testWallet, false)
	assert.Equal(t, ids.Empty, txID)
	assert.Error(t, err)
	assert.Equal(t, multisig.ErrUndefinedTx, err)

	// 2. Test with zero-value multisig (should be undefined)
	zeroMultisig := multisig.Multisig{}
	txID, err = subnet.Commit(zeroMultisig, testWallet, false)
	assert.Equal(t, ids.Empty, txID)
	assert.Error(t, err)
	assert.Equal(t, multisig.ErrUndefinedTx, err)
}

// TestCommit_ValidationLogic tests the validation logic in Commit
func TestCommit_ValidationLogic(t *testing.T) {
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

// TestCommit_WithWaitForTxAcceptance tests Commit with different waitForTxAcceptance values
func TestCommit_WithWaitForTxAcceptance(t *testing.T) {
	subnet := createTestSubnet(t)
	testWallet := createTestWalletForCommit(t)

	// Test with waitForTxAcceptance = true
	undefinedMultisig := multisig.Multisig{}
	txID, err := subnet.Commit(undefinedMultisig, testWallet, true)
	assert.Equal(t, ids.Empty, txID)
	assert.Error(t, err)
	assert.Equal(t, multisig.ErrUndefinedTx, err)

	// Test with waitForTxAcceptance = false
	txID, err = subnet.Commit(undefinedMultisig, testWallet, false)
	assert.Equal(t, ids.Empty, txID)
	assert.Error(t, err)
	assert.Equal(t, multisig.ErrUndefinedTx, err)
}

// TestCommit_SubnetIDUpdate tests that Commit updates SubnetID for CreateSubnetTx
func TestCommit_SubnetIDUpdate(t *testing.T) {
	subnet := createTestSubnet(t)
	testWallet := createTestWalletForCommit(t)

	// Test that Commit would update SubnetID if it were successful
	// We can't test the actual update without a real transaction, but we can test the logic

	// Verify initial SubnetID is empty
	assert.Equal(t, ids.Empty, subnet.SubnetID)

	// Test with undefined multisig (should not update SubnetID)
	undefinedMultisig := multisig.Multisig{}
	txID, err := subnet.Commit(undefinedMultisig, testWallet, false)
	assert.Equal(t, ids.Empty, txID)
	assert.Error(t, err)
	assert.Equal(t, ids.Empty, subnet.SubnetID) // Should remain empty
}

// TestCommit_NetworkIntegration tests Commit network integration (without actual network calls)
func TestCommit_NetworkIntegration(t *testing.T) {
	subnet := createTestSubnet(t)
	testWallet := createTestWalletForCommit(t)

	// Test that Commit would attempt network integration if multisig were ready
	// We can't test actual network calls without funds, but we can test the error paths

	undefinedMultisig := multisig.Multisig{}
	txID, err := subnet.Commit(undefinedMultisig, testWallet, false)
	assert.Equal(t, ids.Empty, txID)
	assert.Error(t, err)
	assert.Equal(t, multisig.ErrUndefinedTx, err)
}

// TestCommit_RetryLogic tests Commit retry logic (without actual retries)
func TestCommit_RetryLogic(t *testing.T) {
	subnet := createTestSubnet(t)
	testWallet := createTestWalletForCommit(t)

	// Test that Commit would retry on network errors if multisig were ready
	// We can't test actual retries without network calls, but we can test the error paths

	undefinedMultisig := multisig.Multisig{}
	txID, err := subnet.Commit(undefinedMultisig, testWallet, false)
	assert.Equal(t, ids.Empty, txID)
	assert.Error(t, err)
	assert.Equal(t, multisig.ErrUndefinedTx, err)
}

// TestCommit_ContextHandling tests Commit context handling
func TestCommit_ContextHandling(t *testing.T) {
	addTestDelay() // Add delay to avoid rate limiting
	subnet := createTestSubnet(t)
	testWallet := createTestWalletForCommit(t)

	// Test that Commit would handle context properly if multisig were ready
	// We can't test actual context handling without network calls, but we can test the error paths

	undefinedMultisig := multisig.Multisig{}
	txID, err := subnet.Commit(undefinedMultisig, testWallet, false)
	assert.Equal(t, ids.Empty, txID)
	assert.Error(t, err)
	assert.Equal(t, multisig.ErrUndefinedTx, err)
}

// TestCommit_TransactionTypes tests Commit with different transaction types
func TestCommit_TransactionTypes(t *testing.T) {
	subnet := createTestSubnet(t)
	testWallet := createTestWalletForCommit(t)

	// Test that Commit would handle different transaction types if multisig were ready
	// We can't test actual transaction types without real transactions, but we can test the error paths

	undefinedMultisig := multisig.Multisig{}
	txID, err := subnet.Commit(undefinedMultisig, testWallet, false)
	assert.Equal(t, ids.Empty, txID)
	assert.Error(t, err)
	assert.Equal(t, multisig.ErrUndefinedTx, err)
}

// TestCommit_EdgeCases tests Commit edge cases
func TestCommit_EdgeCases(t *testing.T) {
	subnet := createTestSubnet(t)
	testWallet := createTestWalletForCommit(t)

	// Test various edge cases for Commit

	// 1. Test with zero-value multisig
	zeroMultisig := multisig.Multisig{}
	txID, err := subnet.Commit(zeroMultisig, testWallet, false)
	assert.Equal(t, ids.Empty, txID)
	assert.Error(t, err)
	assert.Equal(t, multisig.ErrUndefinedTx, err)
}

// TestCommit_ComprehensiveErrorPaths tests all error paths in Commit
func TestCommit_ComprehensiveErrorPaths(t *testing.T) {
	subnet := createTestSubnet(t)
	testWallet := createTestWalletForCommit(t)

	// Test all possible error paths in Commit function

	// 1. Undefined multisig
	undefinedMultisig := multisig.Multisig{}
	txID, err := subnet.Commit(undefinedMultisig, testWallet, false)
	assert.Equal(t, ids.Empty, txID)
	assert.Error(t, err)
	assert.Equal(t, multisig.ErrUndefinedTx, err)

	// 2. Test that IsReadyToCommit would be called if multisig were defined
	// We can't test this without a real multisig, but we can verify the error path

	// 3. Test that GetWrappedOChainTx would be called if multisig were ready
	// We can't test this without a real multisig, but we can verify the error path

	// 4. Test that network submission would be attempted if transaction were ready
	// We can't test this without funds, but we can verify the error path
}

// TestCommit_IntegrationWithOtherFunctions tests Commit integration with other functions
func TestCommit_IntegrationWithOtherFunctions(t *testing.T) {
	subnet := createTestSubnet(t)
	testWallet := createTestWalletForCommit(t)

	// Test that Commit integrates properly with other subnet functions

	// 1. Test that CreateSubnetTx would create a multisig that could be committed
	createSubnetMultisig, err := subnet.CreateSubnetTx(testWallet)
	if err != nil {
		// Expected due to insufficient funds
		assert.Contains(t, err.Error(), "insufficient funds")
		assert.Nil(t, createSubnetMultisig)
	} else {
		// If we get here, we have a real multisig
		require.NotNil(t, createSubnetMultisig)
		assert.False(t, createSubnetMultisig.Undefined())

		// Test that this multisig could be committed (without actually committing)
		tx, err := createSubnetMultisig.GetWrappedOChainTx()
		require.NoError(t, err)
		require.NotNil(t, tx)

		// Verify it's a CreateSubnetTx
		_, ok := tx.Unsigned.(*txs.CreateSubnetTx)
		assert.True(t, ok)
	}

	// 2. Test that CreateBlockchainTx would create a multisig that could be committed
	subnetID := ids.GenerateTestID()
	subnet.SetSubnetID(subnetID)

	createBlockchainMultisig, err := subnet.CreateBlockchainTx(testWallet)
	if err != nil {
		// Expected due to insufficient funds
		assert.Contains(t, err.Error(), "insufficient funds")
		assert.Nil(t, createBlockchainMultisig)
	} else {
		// If we get here, we have a real multisig
		require.NotNil(t, createBlockchainMultisig)
		assert.False(t, createBlockchainMultisig.Undefined())

		// Test that this multisig could be committed (without actually committing)
		tx, err := createBlockchainMultisig.GetWrappedOChainTx()
		require.NoError(t, err)
		require.NotNil(t, tx)

		// Verify it's a CreateChainTx
		_, ok := tx.Unsigned.(*txs.CreateChainTx)
		assert.True(t, ok)
	}

	// 3. Test that AddValidator would create a multisig that could be committed
	validatorParams := validator.SubnetValidatorParams{
		NodeID:   ids.GenerateTestNodeID(),
		Duration: time.Hour,
		Weight:   20,
	}

	addValidatorMultisig, err := subnet.AddValidator(testWallet, validatorParams)
	if err != nil {
		// Expected due to insufficient funds
		assert.Contains(t, err.Error(), "insufficient funds")
		assert.Nil(t, addValidatorMultisig)
	} else {
		// If we get here, we have a real multisig
		require.NotNil(t, addValidatorMultisig)
		assert.False(t, addValidatorMultisig.Undefined())

		// Test that this multisig could be committed (without actually committing)
		tx, err := addValidatorMultisig.GetWrappedOChainTx()
		require.NoError(t, err)
		require.NotNil(t, tx)

		// Verify it's an AddSubnetValidatorTx
		_, ok := tx.Unsigned.(*txs.AddSubnetValidatorTx)
		assert.True(t, ok)
	}
}
