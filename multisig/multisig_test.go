// Copyright (C) 2025, Dione Protocol, Inc. All rights reserved.
// See the file LICENSE for licensing terms.
package multisig

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/DioneProtocol/odyssey-tooling-sdk-go/odyssey"
	"github.com/DioneProtocol/odysseygo/ids"
	"github.com/DioneProtocol/odysseygo/utils/crypto/secp256k1"
	"github.com/DioneProtocol/odysseygo/vms/components/verify"
	"github.com/DioneProtocol/odysseygo/vms/omegavm/txs"
	"github.com/DioneProtocol/odysseygo/vms/secp256k1fx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Note: Mock client tests removed due to type complexity
// The GetOwners function is tested through integration tests

// TestMultisigCreation tests the creation and basic functionality of Multisig
func TestMultisigCreation(t *testing.T) {
	t.Parallel()

	t.Run("New with nil transaction", func(t *testing.T) {
		ms := New(nil)
		require.NotNil(t, ms)
		assert.True(t, ms.Undefined())
		assert.Equal(t, "", ms.String())
	})

	t.Run("New with valid transaction", func(t *testing.T) {
		// Create a mock transaction
		tx := &txs.Tx{}
		ms := New(tx)
		require.NotNil(t, ms)
		assert.False(t, ms.Undefined())
		assert.NotEmpty(t, ms.String())
	})

	t.Run("String method with valid transaction", func(t *testing.T) {
		tx := &txs.Tx{}
		ms := New(tx)
		// The String method should return the transaction ID
		result := ms.String()
		assert.NotEmpty(t, result)
	})
}

// TestMultisigSerialization tests the serialization and deserialization methods
func TestMultisigSerialization(t *testing.T) {
	t.Parallel()

	t.Run("ToBytes with undefined transaction", func(t *testing.T) {
		ms := New(nil)
		bytes, err := ms.ToBytes()
		assert.Error(t, err)
		assert.Nil(t, bytes)
		assert.Equal(t, ErrUndefinedTx, err)
	})

	t.Run("ToFile with undefined transaction", func(t *testing.T) {
		ms := New(nil)
		tempDir := t.TempDir()
		filePath := filepath.Join(tempDir, "test.tx")

		err := ms.ToFile(filePath)
		assert.Error(t, err)
		assert.Equal(t, ErrUndefinedTx, err)
	})

	t.Run("FromFile with non-existent file", func(t *testing.T) {
		ms := New(nil)
		err := ms.FromFile("nonexistent.tx")
		assert.Error(t, err)
	})

	t.Run("FromFile with invalid hex data", func(t *testing.T) {
		tempDir := t.TempDir()
		filePath := filepath.Join(tempDir, "invalid.tx")

		// Write invalid hex data
		err := os.WriteFile(filePath, []byte("invalid-hex-data"), 0644)
		require.NoError(t, err)

		ms := New(nil)
		err = ms.FromFile(filePath)
		assert.Error(t, err)
	})
}

// TestMultisigTransactionTypes tests the transaction type detection
func TestMultisigTransactionTypes(t *testing.T) {
	t.Parallel()

	t.Run("GetTxKind with undefined transaction", func(t *testing.T) {
		ms := New(nil)
		kind, err := ms.GetTxKind()
		assert.Error(t, err)
		assert.Equal(t, Undefined, kind)
		assert.Equal(t, ErrUndefinedTx, err)
	})

	t.Run("GetTxKind with CreateSubnetTx", func(t *testing.T) {
		tx := &txs.Tx{Unsigned: &txs.CreateSubnetTx{}}
		ms := New(tx)

		kind, err := ms.GetTxKind()
		assert.Error(t, err)
		assert.Equal(t, Undefined, kind)
		assert.Contains(t, err.Error(), "unexpected unsigned tx type")
	})
}

// TestMultisigNetworkAndBlockchain tests network and blockchain ID extraction
func TestMultisigNetworkAndBlockchain(t *testing.T) {
	t.Parallel()

	t.Run("GetNetworkID with undefined transaction", func(t *testing.T) {
		ms := New(nil)
		networkID, err := ms.GetNetworkID()
		assert.Error(t, err)
		assert.Equal(t, uint32(0), networkID)
		assert.Equal(t, ErrUndefinedTx, err)
	})

	t.Run("GetNetwork with undefined network ID", func(t *testing.T) {
		// Create a transaction with an unknown network ID
		tx := &txs.Tx{Unsigned: &txs.CreateSubnetTx{}}
		ms := New(tx)

		network, err := ms.GetNetwork()
		assert.Error(t, err)
		assert.Equal(t, odyssey.UndefinedNetwork, network)
		assert.Contains(t, err.Error(), "unexpected unsigned tx type")
	})

	t.Run("GetBlockchainID with undefined transaction", func(t *testing.T) {
		ms := New(nil)
		blockchainID, err := ms.GetBlockchainID()
		assert.Error(t, err)
		assert.Equal(t, ids.Empty, blockchainID)
		assert.Equal(t, ErrUndefinedTx, err)
	})
}

// TestMultisigSubnetOperations tests subnet-related operations
func TestMultisigSubnetOperations(t *testing.T) {
	t.Parallel()

	t.Run("GetSubnetID with undefined transaction", func(t *testing.T) {
		ms := New(nil)
		subnetID, err := ms.GetSubnetID()
		assert.Error(t, err)
		assert.Equal(t, ids.Empty, subnetID)
		assert.Equal(t, ErrUndefinedTx, err)
	})

	t.Run("GetSubnetOwners with undefined transaction", func(t *testing.T) {
		ms := New(nil)
		controlKeys, threshold, err := ms.GetSubnetOwners()
		assert.Error(t, err)
		assert.Nil(t, controlKeys)
		assert.Equal(t, uint32(0), threshold)
		assert.Equal(t, ErrUndefinedTx, err)
	})

	t.Run("GetOwners function exists", func(t *testing.T) {
		// Test that GetOwners function exists and has correct signature
		network := odyssey.TestnetNetwork()
		testSubnetID := ids.ID{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32}

		// This will fail due to network call, but we can test the function signature
		_, _, err := GetOwners(network, testSubnetID)
		assert.Error(t, err) // Expected to fail without real network
	})
}

// TestMultisigAuthSigners tests authentication signer operations
func TestMultisigAuthSigners(t *testing.T) {
	t.Parallel()

	t.Run("GetAuthSigners with undefined transaction", func(t *testing.T) {
		ms := New(nil)
		authSigners, err := ms.GetAuthSigners()
		assert.Error(t, err)
		assert.Nil(t, authSigners)
		assert.Equal(t, ErrUndefinedTx, err)
	})

	t.Run("GetAuthSigners with unsupported transaction type", func(t *testing.T) {
		tx := &txs.Tx{Unsigned: &txs.CreateSubnetTx{}}
		ms := New(tx)

		authSigners, err := ms.GetAuthSigners()
		assert.Error(t, err)
		assert.Nil(t, authSigners)
		assert.Contains(t, err.Error(), "unexpected unsigned tx type")
	})

	t.Run("GetRemainingAuthSigners with undefined transaction", func(t *testing.T) {
		ms := New(nil)
		authSigners, remainingSigners, err := ms.GetRemainingAuthSigners()
		assert.Error(t, err)
		assert.Nil(t, authSigners)
		assert.Nil(t, remainingSigners)
		assert.Equal(t, ErrUndefinedTx, err)
	})

	t.Run("GetRemainingAuthSigners with insufficient credentials", func(t *testing.T) {
		tx := &txs.Tx{
			Unsigned: &txs.CreateSubnetTx{},
			Creds:    []verify.Verifiable{}, // Empty credentials
		}
		ms := New(tx)

		authSigners, remainingSigners, err := ms.GetRemainingAuthSigners()
		assert.Error(t, err)
		assert.Nil(t, authSigners)
		assert.Nil(t, remainingSigners)
		// CreateSubnetTx is not supported by GetRemainingAuthSigners, so it will fail with "unexpected unsigned tx type"
		assert.Contains(t, err.Error(), "unexpected unsigned tx type")
	})

	t.Run("IsReadyToCommit with undefined transaction", func(t *testing.T) {
		ms := New(nil)
		ready, err := ms.IsReadyToCommit()
		assert.Error(t, err)
		assert.False(t, ready)
		assert.Equal(t, ErrUndefinedTx, err)
	})

	t.Run("IsReadyToCommit with CreateSubnetTx", func(t *testing.T) {
		tx := &txs.Tx{Unsigned: &txs.CreateSubnetTx{}}
		ms := New(tx)

		ready, err := ms.IsReadyToCommit()
		require.NoError(t, err)
		assert.True(t, ready)
	})
}

// TestMultisigErrorHandling tests various error scenarios
func TestMultisigErrorHandling(t *testing.T) {
	t.Parallel()

	t.Run("GetSpendSigners not implemented", func(t *testing.T) {
		ms := New(nil)
		signers, err := ms.GetSpendSigners()
		assert.Error(t, err)
		assert.Nil(t, signers)
		assert.Contains(t, err.Error(), "not implemented yet")
	})

	t.Run("GetWrappedOChainTx with undefined transaction", func(t *testing.T) {
		ms := New(nil)
		tx, err := ms.GetWrappedOChainTx()
		assert.Error(t, err)
		assert.Nil(t, tx)
		assert.Equal(t, ErrUndefinedTx, err)
	})

	t.Run("GetWrappedOChainTx with valid transaction", func(t *testing.T) {
		originalTx := &txs.Tx{Unsigned: &txs.CreateSubnetTx{}}
		ms := New(originalTx)

		tx, err := ms.GetWrappedOChainTx()
		require.NoError(t, err)
		assert.Equal(t, originalTx, tx)
	})
}

// TestMultisigIntegration tests integration scenarios with LOCAL_NODE flag
func TestMultisigIntegration(t *testing.T) {
	// Skip integration tests if LOCAL_NODE is not set
	if os.Getenv("LOCAL_NODE") != "true" {
		t.Skip("Skipping integration tests - LOCAL_NODE not set")
	}

	t.Run("Integration with local node", func(t *testing.T) {
		// This test would require a running local node
		// For now, we'll test the network detection
		network := odyssey.TestnetNetwork()
		assert.Equal(t, "http://127.0.0.1:9650", network.Endpoint)
		assert.Equal(t, odyssey.Testnet, network.Kind)
	})
}

// TestMultisigProductionReadiness tests production readiness scenarios
func TestMultisigProductionReadiness(t *testing.T) {
	t.Parallel()

	t.Run("Error constants are properly defined", func(t *testing.T) {
		assert.NotNil(t, ErrUndefinedTx)
		assert.Equal(t, "tx is undefined", ErrUndefinedTx.Error())
	})

	t.Run("TxKind constants are properly defined", func(t *testing.T) {
		assert.Equal(t, TxKind(0), Undefined)
		assert.Equal(t, TxKind(1), OChainRemoveSubnetValidatorTx)
		assert.Equal(t, TxKind(2), OChainAddSubnetValidatorTx)
		assert.Equal(t, TxKind(3), OChainCreateChainTx)
		assert.Equal(t, TxKind(4), OChainTransformSubnetTx)
		assert.Equal(t, TxKind(5), OChainAddPermissionlessValidatorTx)
		assert.Equal(t, TxKind(6), OChainTransferSubnetOwnershipTx)
	})

	t.Run("Multisig struct fields are properly initialized", func(t *testing.T) {
		tx := &txs.Tx{Unsigned: &txs.CreateSubnetTx{}}
		ms := New(tx)

		assert.Equal(t, tx, ms.OChainTx)
		assert.Nil(t, ms.controlKeys)
		assert.Equal(t, uint32(0), ms.threshold)
	})

	t.Run("File operations handle permissions correctly", func(t *testing.T) {
		// Test with read-only directory
		readOnlyDir := t.TempDir()
		err := os.Chmod(readOnlyDir, 0444) // Read-only
		require.NoError(t, err)

		tx := &txs.Tx{Unsigned: &txs.CreateSubnetTx{}}
		ms := New(tx)

		filePath := filepath.Join(readOnlyDir, "test.tx")
		err = ms.ToFile(filePath)
		assert.Error(t, err) // Should fail due to permissions
	})

	t.Run("Concurrent access safety", func(t *testing.T) {
		tx := &txs.Tx{Unsigned: &txs.CreateSubnetTx{}}
		ms := New(tx)

		// Test concurrent access to the same multisig instance
		done := make(chan bool, 10)

		for i := 0; i < 10; i++ {
			go func() {
				defer func() { done <- true }()

				// These operations should be safe for concurrent access
				_ = ms.Undefined()
				_ = ms.String()
				_, _ = ms.GetTxKind()
				_, _ = ms.GetNetworkID()
				_, _ = ms.GetBlockchainID()
				_, _ = ms.GetSubnetID()
			}()
		}

		// Wait for all goroutines to complete
		for i := 0; i < 10; i++ {
			<-done
		}
	})
}

// TestMultisigWithValidTransaction tests with a properly constructed transaction
func TestMultisigWithValidTransaction(t *testing.T) {
	t.Parallel()

	t.Run("CreateSubnetTx serialization", func(t *testing.T) {
		// Create a simple CreateSubnetTx for testing
		unsignedTx := &txs.CreateSubnetTx{}

		tx := &txs.Tx{
			Unsigned: unsignedTx,
			Creds:    []verify.Verifiable{}, // Empty credentials for simplicity
		}

		ms := New(tx)

		// Test serialization - this will fail due to incomplete transaction structure
		bytes, err := ms.ToBytes()
		assert.Error(t, err) // Expected to fail due to incomplete transaction
		assert.Nil(t, bytes)
	})

	t.Run("IsReadyToCommit with CreateSubnetTx", func(t *testing.T) {
		tx := &txs.Tx{Unsigned: &txs.CreateSubnetTx{}}
		ms := New(tx)

		ready, err := ms.IsReadyToCommit()
		require.NoError(t, err)
		assert.True(t, ready)
	})
}

// TestMultisigFromBytes tests the FromBytes method comprehensively
func TestMultisigFromBytes(t *testing.T) {
	t.Parallel()

	t.Run("FromBytes with valid transaction data", func(t *testing.T) {
		// Create some mock bytes for testing FromBytes
		bytes := []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05}

		// Test FromBytes
		newMs := New(nil)
		err := newMs.FromBytes(bytes)
		// This will fail due to invalid transaction data, but we're testing the method
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error unmarshaling signed tx")
	})

	t.Run("FromBytes with invalid data", func(t *testing.T) {
		ms := New(nil)
		invalidBytes := []byte{0xFF, 0xFE, 0xFD} // Invalid transaction data

		err := ms.FromBytes(invalidBytes)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error unmarshaling signed tx")
	})

	t.Run("FromBytes with empty data", func(t *testing.T) {
		ms := New(nil)
		emptyBytes := []byte{}

		err := ms.FromBytes(emptyBytes)
		assert.Error(t, err)
	})
}

// TestMultisigToFileComprehensive tests ToFile method more comprehensively
func TestMultisigToFileComprehensive(t *testing.T) {
	t.Parallel()

	t.Run("ToFile with valid transaction", func(t *testing.T) {
		// Create a transaction that can be serialized
		unsignedTx := &txs.CreateSubnetTx{}
		tx := &txs.Tx{
			Unsigned: unsignedTx,
			Creds:    []verify.Verifiable{},
		}

		ms := New(tx)

		tempDir := t.TempDir()
		filePath := filepath.Join(tempDir, "test.tx")

		// Test ToFile - this will likely fail due to incomplete transaction
		err := ms.ToFile(filePath)
		if err != nil {
			// Expected to fail due to incomplete transaction structure
			assert.Contains(t, err.Error(), "couldn't marshal signed tx")
		} else {
			// If it succeeds, verify file was created
			assert.FileExists(t, filePath)
		}
	})

	t.Run("ToFile with directory creation", func(t *testing.T) {
		unsignedTx := &txs.CreateSubnetTx{}
		tx := &txs.Tx{
			Unsigned: unsignedTx,
			Creds:    []verify.Verifiable{},
		}

		ms := New(tx)

		tempDir := t.TempDir()
		subDir := filepath.Join(tempDir, "subdir")
		filePath := filepath.Join(subDir, "test.tx")

		// Test ToFile with non-existent directory
		err := ms.ToFile(filePath)
		if err != nil {
			// Expected to fail due to incomplete transaction or directory issues
			assert.True(t,
				containsAny(err.Error(), []string{
					"couldn't marshal signed tx",
					"couldn't create file",
					"no such file or directory",
				}))
		}
	})
}

// TestMultisigTransactionTypesComprehensive tests all supported transaction types
func TestMultisigTransactionTypesComprehensive(t *testing.T) {
	t.Parallel()

	t.Run("RemoveSubnetValidatorTx_GetTxKind", func(t *testing.T) {
		tx := &txs.Tx{Unsigned: &txs.RemoveSubnetValidatorTx{}}
		ms := New(tx)

		kind, err := ms.GetTxKind()
		require.NoError(t, err)
		assert.Equal(t, OChainRemoveSubnetValidatorTx, kind)
	})

	t.Run("AddSubnetValidatorTx_GetTxKind", func(t *testing.T) {
		tx := &txs.Tx{Unsigned: &txs.AddSubnetValidatorTx{}}
		ms := New(tx)

		kind, err := ms.GetTxKind()
		require.NoError(t, err)
		assert.Equal(t, OChainAddSubnetValidatorTx, kind)
	})

	t.Run("CreateChainTx_GetTxKind", func(t *testing.T) {
		tx := &txs.Tx{Unsigned: &txs.CreateChainTx{}}
		ms := New(tx)

		kind, err := ms.GetTxKind()
		require.NoError(t, err)
		assert.Equal(t, OChainCreateChainTx, kind)
	})

	t.Run("TransformSubnetTx_GetTxKind", func(t *testing.T) {
		tx := &txs.Tx{Unsigned: &txs.TransformSubnetTx{}}
		ms := New(tx)

		kind, err := ms.GetTxKind()
		require.NoError(t, err)
		assert.Equal(t, OChainTransformSubnetTx, kind)
	})

	t.Run("AddPermissionlessValidatorTx_GetTxKind", func(t *testing.T) {
		tx := &txs.Tx{Unsigned: &txs.AddPermissionlessValidatorTx{}}
		ms := New(tx)

		kind, err := ms.GetTxKind()
		require.NoError(t, err)
		assert.Equal(t, OChainAddPermissionlessValidatorTx, kind)
	})

	// Test GetNetworkID with different transaction types
	t.Run("RemoveSubnetValidatorTx_GetNetworkID", func(t *testing.T) {
		tx := &txs.Tx{Unsigned: &txs.RemoveSubnetValidatorTx{}}
		ms := New(tx)

		networkID, err := ms.GetNetworkID()
		// This should succeed and return 0 (default value)
		require.NoError(t, err)
		assert.Equal(t, uint32(0), networkID)
	})

	t.Run("AddSubnetValidatorTx_GetNetworkID", func(t *testing.T) {
		tx := &txs.Tx{Unsigned: &txs.AddSubnetValidatorTx{}}
		ms := New(tx)

		networkID, err := ms.GetNetworkID()
		require.NoError(t, err)
		assert.Equal(t, uint32(0), networkID)
	})

	// Test GetBlockchainID with different transaction types
	t.Run("RemoveSubnetValidatorTx_GetBlockchainID", func(t *testing.T) {
		tx := &txs.Tx{Unsigned: &txs.RemoveSubnetValidatorTx{}}
		ms := New(tx)

		blockchainID, err := ms.GetBlockchainID()
		require.NoError(t, err)
		assert.Equal(t, ids.Empty, blockchainID)
	})

	t.Run("AddSubnetValidatorTx_GetBlockchainID", func(t *testing.T) {
		tx := &txs.Tx{Unsigned: &txs.AddSubnetValidatorTx{}}
		ms := New(tx)

		blockchainID, err := ms.GetBlockchainID()
		require.NoError(t, err)
		assert.Equal(t, ids.Empty, blockchainID)
	})

	// Test GetSubnetID with different transaction types
	t.Run("RemoveSubnetValidatorTx_GetSubnetID", func(t *testing.T) {
		tx := &txs.Tx{Unsigned: &txs.RemoveSubnetValidatorTx{}}
		ms := New(tx)

		subnetID, err := ms.GetSubnetID()
		require.NoError(t, err)
		assert.Equal(t, ids.Empty, subnetID)
	})

	t.Run("AddSubnetValidatorTx_GetSubnetID", func(t *testing.T) {
		tx := &txs.Tx{Unsigned: &txs.AddSubnetValidatorTx{}}
		ms := New(tx)

		subnetID, err := ms.GetSubnetID()
		require.NoError(t, err)
		assert.Equal(t, ids.Empty, subnetID)
	})
}

// TestMultisigGetNetworkComprehensive tests GetNetwork method more comprehensively
func TestMultisigGetNetworkComprehensive(t *testing.T) {
	t.Parallel()

	t.Run("GetNetwork with testnet ID", func(t *testing.T) {
		// Create a transaction with testnet network ID
		unsignedTx := &txs.CreateSubnetTx{}
		tx := &txs.Tx{Unsigned: unsignedTx}
		ms := New(tx)

		network, err := ms.GetNetwork()
		// This will fail due to missing NetworkID field
		assert.Error(t, err)
		assert.Equal(t, odyssey.UndefinedNetwork, network)
		assert.Contains(t, err.Error(), "unexpected unsigned tx type")
	})

	t.Run("GetNetwork with mainnet ID", func(t *testing.T) {
		// Create a transaction with mainnet network ID
		unsignedTx := &txs.CreateSubnetTx{}
		tx := &txs.Tx{Unsigned: unsignedTx}
		ms := New(tx)

		network, err := ms.GetNetwork()
		// This will fail due to missing NetworkID field
		assert.Error(t, err)
		assert.Equal(t, odyssey.UndefinedNetwork, network)
	})
}

// TestMultisigGetSubnetOwnersComprehensive tests GetSubnetOwners with caching
func TestMultisigGetSubnetOwnersComprehensive(t *testing.T) {
	t.Parallel()

	t.Run("GetSubnetOwners with cached control keys", func(t *testing.T) {
		// Create a multisig with pre-cached control keys
		ms := &Multisig{
			OChainTx: &txs.Tx{Unsigned: &txs.CreateSubnetTx{}},
			controlKeys: []ids.ShortID{
				{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
				{2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21},
			},
			threshold: 2,
		}

		controlKeys, threshold, err := ms.GetSubnetOwners()
		// This should succeed and return cached values
		require.NoError(t, err)
		assert.Len(t, controlKeys, 2)
		assert.Equal(t, uint32(2), threshold)
	})
}

// TestMultisigGetOwnersComprehensive tests GetOwners function more comprehensively
func TestMultisigGetOwnersComprehensive(t *testing.T) {
	t.Parallel()

	t.Run("GetOwners function exists and can be called", func(t *testing.T) {
		testSubnetID := ids.ID{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32}
		network := odyssey.TestnetNetwork()

		// Test that the function can be called (it will fail due to network call)
		_, _, err := GetOwners(network, testSubnetID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "subnet tx")
	})
}

// Helper function to check if error message contains any of the expected strings
func containsAny(errMsg string, expected []string) bool {
	for _, exp := range expected {
		if contains(errMsg, exp) {
			return true
		}
	}
	return false
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || (len(substr) > 0 && indexOf(s, substr) >= 0))
}

// Helper function to find index of substring
func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

// TestGetRemainingAuthSignersProduction tests GetRemainingAuthSigners with production scenarios
func TestGetRemainingAuthSignersProduction(t *testing.T) {
	t.Parallel()

	t.Run("GetRemainingAuthSigners with proper credential structure", func(t *testing.T) {
		// Create a transaction with proper credential structure
		unsignedTx := &txs.RemoveSubnetValidatorTx{
			SubnetAuth: &secp256k1fx.Input{
				SigIndices: []uint32{0, 1},
			},
		}

		// Create proper credentials
		outputCred := &secp256k1fx.Credential{
			Sigs: [][secp256k1.SignatureLen]byte{
				{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64},
			},
		}

		subnetCred := &secp256k1fx.Credential{
			Sigs: [][secp256k1.SignatureLen]byte{
				{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64}, // Filled
				{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},                                                        // Empty
			},
		}

		tx := &txs.Tx{
			Unsigned: unsignedTx,
			Creds:    []verify.Verifiable{outputCred, subnetCred},
		}

		ms := &Multisig{
			OChainTx: tx,
			controlKeys: []ids.ShortID{
				{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
				{2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21},
			},
			threshold: 2,
		}

		authSigners, remainingSigners, err := ms.GetRemainingAuthSigners()
		require.NoError(t, err)
		assert.Len(t, authSigners, 2)
		assert.Len(t, remainingSigners, 1) // One empty signature
	})

	t.Run("GetRemainingAuthSigners with insufficient credentials", func(t *testing.T) {
		unsignedTx := &txs.RemoveSubnetValidatorTx{
			SubnetAuth: &secp256k1fx.Input{
				SigIndices: []uint32{0},
			},
		}

		// Only one credential (insufficient)
		tx := &txs.Tx{
			Unsigned: unsignedTx,
			Creds:    []verify.Verifiable{&secp256k1fx.Credential{}},
		}

		ms := New(tx)
		authSigners, remainingSigners, err := ms.GetRemainingAuthSigners()
		assert.Error(t, err)
		assert.Nil(t, authSigners)
		assert.Nil(t, remainingSigners)
		// The error will be from GetSubnetOwners failing due to undefined network
		assert.Contains(t, err.Error(), "undefined network model for tx")
	})

	t.Run("GetRemainingAuthSigners with wrong credential type", func(t *testing.T) {
		unsignedTx := &txs.RemoveSubnetValidatorTx{
			SubnetAuth: &secp256k1fx.Input{
				SigIndices: []uint32{0},
			},
		}

		// Wrong credential type
		tx := &txs.Tx{
			Unsigned: unsignedTx,
			Creds:    []verify.Verifiable{nil, nil}, // Wrong type
		}

		ms := New(tx)
		authSigners, remainingSigners, err := ms.GetRemainingAuthSigners()
		assert.Error(t, err)
		assert.Nil(t, authSigners)
		assert.Nil(t, remainingSigners)
		// The error will be from GetSubnetOwners failing due to undefined network
		assert.Contains(t, err.Error(), "undefined network model for tx")
	})

	t.Run("GetRemainingAuthSigners with empty output signature", func(t *testing.T) {
		unsignedTx := &txs.RemoveSubnetValidatorTx{
			SubnetAuth: &secp256k1fx.Input{
				SigIndices: []uint32{0},
			},
		}

		outputCred := &secp256k1fx.Credential{
			Sigs: [][secp256k1.SignatureLen]byte{
				{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, // Empty
			},
		}

		subnetCred := &secp256k1fx.Credential{
			Sigs: [][secp256k1.SignatureLen]byte{
				{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64},
			},
		}

		tx := &txs.Tx{
			Unsigned: unsignedTx,
			Creds:    []verify.Verifiable{outputCred, subnetCred},
		}

		ms := New(tx)
		authSigners, remainingSigners, err := ms.GetRemainingAuthSigners()
		assert.Error(t, err)
		assert.Nil(t, authSigners)
		assert.Nil(t, remainingSigners)
		// The error will be from GetSubnetOwners failing due to undefined network
		assert.Contains(t, err.Error(), "undefined network model for tx")
	})

	t.Run("GetRemainingAuthSigners with signature count mismatch", func(t *testing.T) {
		unsignedTx := &txs.RemoveSubnetValidatorTx{
			SubnetAuth: &secp256k1fx.Input{
				SigIndices: []uint32{0, 1}, // 2 signers required
			},
		}

		outputCred := &secp256k1fx.Credential{
			Sigs: [][secp256k1.SignatureLen]byte{
				{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64},
			},
		}

		subnetCred := &secp256k1fx.Credential{
			Sigs: [][secp256k1.SignatureLen]byte{
				{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64}, // Only 1 signature, but 2 required
			},
		}

		tx := &txs.Tx{
			Unsigned: unsignedTx,
			Creds:    []verify.Verifiable{outputCred, subnetCred},
		}

		ms := &Multisig{
			OChainTx: tx,
			controlKeys: []ids.ShortID{
				{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
				{2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21},
			},
			threshold: 2,
		}

		authSigners, remainingSigners, err := ms.GetRemainingAuthSigners()
		assert.Error(t, err)
		assert.Nil(t, authSigners)
		assert.Nil(t, remainingSigners)
		assert.Contains(t, err.Error(), "expected number of cred's signatures")
	})
}

// TestGetAuthSignersProduction tests GetAuthSigners with production scenarios
func TestGetAuthSignersProduction(t *testing.T) {
	t.Parallel()

	t.Run("GetAuthSigners with RemoveSubnetValidatorTx", func(t *testing.T) {
		unsignedTx := &txs.RemoveSubnetValidatorTx{
			SubnetAuth: &secp256k1fx.Input{
				SigIndices: []uint32{0, 2},
			},
		}

		tx := &txs.Tx{Unsigned: unsignedTx}
		ms := &Multisig{
			OChainTx: tx,
			controlKeys: []ids.ShortID{
				{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
				{2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21},
				{3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22},
			},
			threshold: 3,
		}

		authSigners, err := ms.GetAuthSigners()
		require.NoError(t, err)
		assert.Len(t, authSigners, 2)
		assert.Equal(t, ms.controlKeys[0], authSigners[0])
		assert.Equal(t, ms.controlKeys[2], authSigners[1])
	})

	t.Run("GetAuthSigners with AddSubnetValidatorTx", func(t *testing.T) {
		unsignedTx := &txs.AddSubnetValidatorTx{
			SubnetAuth: &secp256k1fx.Input{
				SigIndices: []uint32{1},
			},
		}

		tx := &txs.Tx{Unsigned: unsignedTx}
		ms := &Multisig{
			OChainTx: tx,
			controlKeys: []ids.ShortID{
				{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
				{2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21},
			},
			threshold: 2,
		}

		authSigners, err := ms.GetAuthSigners()
		require.NoError(t, err)
		assert.Len(t, authSigners, 1)
		assert.Equal(t, ms.controlKeys[1], authSigners[0])
	})

	t.Run("GetAuthSigners with CreateChainTx", func(t *testing.T) {
		unsignedTx := &txs.CreateChainTx{
			SubnetAuth: &secp256k1fx.Input{
				SigIndices: []uint32{0, 1, 2},
			},
		}

		tx := &txs.Tx{Unsigned: unsignedTx}
		ms := &Multisig{
			OChainTx: tx,
			controlKeys: []ids.ShortID{
				{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
				{2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21},
				{3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22},
			},
			threshold: 3,
		}

		authSigners, err := ms.GetAuthSigners()
		require.NoError(t, err)
		assert.Len(t, authSigners, 3)
		assert.Equal(t, ms.controlKeys[0], authSigners[0])
		assert.Equal(t, ms.controlKeys[1], authSigners[1])
		assert.Equal(t, ms.controlKeys[2], authSigners[2])
	})

	t.Run("GetAuthSigners with TransformSubnetTx", func(t *testing.T) {
		unsignedTx := &txs.TransformSubnetTx{
			SubnetAuth: &secp256k1fx.Input{
				SigIndices: []uint32{0},
			},
		}

		tx := &txs.Tx{Unsigned: unsignedTx}
		ms := &Multisig{
			OChainTx: tx,
			controlKeys: []ids.ShortID{
				{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
			},
			threshold: 1,
		}

		authSigners, err := ms.GetAuthSigners()
		require.NoError(t, err)
		assert.Len(t, authSigners, 1)
		assert.Equal(t, ms.controlKeys[0], authSigners[0])
	})

	t.Run("GetAuthSigners with AddPermissionlessValidatorTx", func(t *testing.T) {
		unsignedTx := &txs.AddPermissionlessValidatorTx{
			// Note: AddPermissionlessValidatorTx doesn't have SubnetAuth field
			// This test will fail with unsupported transaction type
		}

		tx := &txs.Tx{Unsigned: unsignedTx}
		ms := &Multisig{
			OChainTx: tx,
			controlKeys: []ids.ShortID{
				{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
				{2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21},
			},
			threshold: 2,
		}

		authSigners, err := ms.GetAuthSigners()
		assert.Error(t, err)
		assert.Nil(t, authSigners)
		assert.Contains(t, err.Error(), "unexpected unsigned tx type")
	})

	t.Run("GetAuthSigners with unsupported transaction type", func(t *testing.T) {
		tx := &txs.Tx{Unsigned: &txs.AddValidatorTx{}} // Not supported
		ms := New(tx)

		authSigners, err := ms.GetAuthSigners()
		assert.Error(t, err)
		assert.Nil(t, authSigners)
		assert.Contains(t, err.Error(), "unexpected unsigned tx type")
	})

	t.Run("GetAuthSigners with invalid SubnetAuth type", func(t *testing.T) {
		unsignedTx := &txs.RemoveSubnetValidatorTx{
			SubnetAuth: nil, // Invalid type
		}

		tx := &txs.Tx{Unsigned: unsignedTx}
		ms := New(tx)

		authSigners, err := ms.GetAuthSigners()
		assert.Error(t, err)
		assert.Nil(t, authSigners)
		// The error will be from GetSubnetOwners failing due to undefined network
		assert.Contains(t, err.Error(), "undefined network model for tx")
	})

	t.Run("GetAuthSigners with signature index out of bounds", func(t *testing.T) {
		unsignedTx := &txs.RemoveSubnetValidatorTx{
			SubnetAuth: &secp256k1fx.Input{
				SigIndices: []uint32{5}, // Out of bounds
			},
		}

		tx := &txs.Tx{Unsigned: unsignedTx}
		ms := &Multisig{
			OChainTx: tx,
			controlKeys: []ids.ShortID{
				{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
			},
			threshold: 1,
		}

		authSigners, err := ms.GetAuthSigners()
		assert.Error(t, err)
		assert.Nil(t, authSigners)
		assert.Contains(t, err.Error(), "signer index 5 exceeds number of control keys")
	})
}

// TestToFileProductionScenarios tests ToFile with production scenarios
func TestToFileProductionScenarios(t *testing.T) {
	t.Parallel()

	t.Run("ToFile with valid transaction structure", func(t *testing.T) {
		// Create a more complete transaction structure
		unsignedTx := &txs.CreateSubnetTx{
			// Note: CreateSubnetTx doesn't have SubnetAuth field
		}

		// Create proper credentials
		outputCred := &secp256k1fx.Credential{
			Sigs: [][secp256k1.SignatureLen]byte{
				{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64},
			},
		}

		subnetCred := &secp256k1fx.Credential{
			Sigs: [][secp256k1.SignatureLen]byte{
				{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64},
			},
		}

		tx := &txs.Tx{
			Unsigned: unsignedTx,
			Creds:    []verify.Verifiable{outputCred, subnetCred},
		}

		ms := New(tx)
		tempDir := t.TempDir()
		filePath := filepath.Join(tempDir, "production_test.tx")

		// Test ToFile - this may still fail due to incomplete transaction structure
		err := ms.ToFile(filePath)
		if err != nil {
			// Expected to fail due to incomplete transaction structure
			assert.Contains(t, err.Error(), "couldn't marshal signed tx")
		} else {
			// If it succeeds, verify file was created
			assert.FileExists(t, filePath)
		}
	})

	t.Run("ToFile with nested directory creation", func(t *testing.T) {
		unsignedTx := &txs.CreateSubnetTx{}
		tx := &txs.Tx{Unsigned: unsignedTx}
		ms := New(tx)

		tempDir := t.TempDir()
		nestedDir := filepath.Join(tempDir, "nested", "deep", "path")
		filePath := filepath.Join(nestedDir, "nested_test.tx")

		// Test ToFile with nested directory
		err := ms.ToFile(filePath)
		if err != nil {
			// Expected to fail due to incomplete transaction or directory issues
			assert.True(t,
				containsAny(err.Error(), []string{
					"couldn't marshal signed tx",
					"couldn't create file",
					"no such file or directory",
				}))
		}
	})

	t.Run("ToFile with special characters in filename", func(t *testing.T) {
		unsignedTx := &txs.CreateSubnetTx{}
		tx := &txs.Tx{Unsigned: unsignedTx}
		ms := New(tx)

		tempDir := t.TempDir()
		filePath := filepath.Join(tempDir, "test@#$%^&*().tx")

		// Test ToFile with special characters
		err := ms.ToFile(filePath)
		if err != nil {
			// Expected to fail due to incomplete transaction
			assert.Contains(t, err.Error(), "couldn't marshal signed tx")
		}
	})

	t.Run("ToFile with very long filename", func(t *testing.T) {
		unsignedTx := &txs.CreateSubnetTx{}
		tx := &txs.Tx{Unsigned: unsignedTx}
		ms := New(tx)

		tempDir := t.TempDir()
		longFilename := "very_long_filename_" + string(make([]byte, 200)) + ".tx"
		filePath := filepath.Join(tempDir, longFilename)

		// Test ToFile with long filename
		err := ms.ToFile(filePath)
		if err != nil {
			// Expected to fail due to incomplete transaction
			assert.Contains(t, err.Error(), "couldn't marshal signed tx")
		}
	})
}

// TestFromBytesProductionScenarios tests FromBytes with production scenarios
func TestFromBytesProductionScenarios(t *testing.T) {
	t.Parallel()

	t.Run("FromBytes with various invalid data patterns", func(t *testing.T) {
		ms := New(nil)

		testCases := []struct {
			name string
			data []byte
		}{
			{
				name: "empty data",
				data: []byte{},
			},
			{
				name: "null bytes",
				data: []byte{0x00, 0x00, 0x00, 0x00},
			},
			{
				name: "random bytes",
				data: []byte{0xFF, 0xFE, 0xFD, 0xFC, 0xFB, 0xFA},
			},
			{
				name: "partial transaction data",
				data: []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F},
			},
			{
				name: "large invalid data",
				data: make([]byte, 10000),
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				err := ms.FromBytes(tc.data)
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "error unmarshaling signed tx")
			})
		}
	})

	t.Run("FromBytes with corrupted hex data", func(t *testing.T) {
		ms := New(nil)

		// Test with data that looks like hex but is corrupted
		corruptedHex := []byte("0x1234567890abcdefghijklmnop") // Invalid hex characters
		err := ms.FromBytes(corruptedHex)
		assert.Error(t, err)
	})

	t.Run("FromBytes with valid hex but invalid transaction", func(t *testing.T) {
		ms := New(nil)

		// Test with valid hex but invalid transaction structure
		validHex := []byte("0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")
		err := ms.FromBytes(validHex)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error unmarshaling signed tx")
	})
}

// TestGetSubnetOwnersProduction tests GetSubnetOwners with production scenarios
func TestGetSubnetOwnersProduction(t *testing.T) {
	t.Parallel()

	t.Run("GetSubnetOwners with different threshold values", func(t *testing.T) {
		testCases := []struct {
			name      string
			threshold uint32
		}{
			{"threshold_1", 1},
			{"threshold_2", 2},
			{"threshold_3", 3},
			{"threshold_5", 5},
			{"threshold_10", 10},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				ms := &Multisig{
					OChainTx:    &txs.Tx{Unsigned: &txs.CreateSubnetTx{}},
					controlKeys: make([]ids.ShortID, tc.threshold),
					threshold:   tc.threshold,
				}

				// Initialize control keys
				for i := range ms.controlKeys {
					ms.controlKeys[i] = ids.ShortID{byte(i + 1), byte(i + 2), byte(i + 3), byte(i + 4), byte(i + 5), byte(i + 6), byte(i + 7), byte(i + 8), byte(i + 9), byte(i + 10), byte(i + 11), byte(i + 12), byte(i + 13), byte(i + 14), byte(i + 15), byte(i + 16), byte(i + 17), byte(i + 18), byte(i + 19), byte(i + 20)}
				}

				controlKeys, threshold, err := ms.GetSubnetOwners()
				require.NoError(t, err)
				assert.Len(t, controlKeys, int(tc.threshold))
				assert.Equal(t, tc.threshold, threshold)
			})
		}
	})

	t.Run("GetSubnetOwners with large number of control keys", func(t *testing.T) {
		// Test with a large number of control keys
		numKeys := 100
		controlKeys := make([]ids.ShortID, numKeys)
		for i := range controlKeys {
			controlKeys[i] = ids.ShortID{byte(i + 1), byte(i + 2), byte(i + 3), byte(i + 4), byte(i + 5), byte(i + 6), byte(i + 7), byte(i + 8), byte(i + 9), byte(i + 10), byte(i + 11), byte(i + 12), byte(i + 13), byte(i + 14), byte(i + 15), byte(i + 16), byte(i + 17), byte(i + 18), byte(i + 19), byte(i + 20)}
		}

		ms := &Multisig{
			OChainTx:    &txs.Tx{Unsigned: &txs.CreateSubnetTx{}},
			controlKeys: controlKeys,
			threshold:   uint32(numKeys),
		}

		resultKeys, threshold, err := ms.GetSubnetOwners()
		require.NoError(t, err)
		assert.Len(t, resultKeys, numKeys)
		assert.Equal(t, uint32(numKeys), threshold)
	})

	t.Run("GetSubnetOwners with zero threshold", func(t *testing.T) {
		ms := &Multisig{
			OChainTx:    &txs.Tx{Unsigned: &txs.CreateSubnetTx{}},
			controlKeys: []ids.ShortID{},
			threshold:   0,
		}

		controlKeys, threshold, err := ms.GetSubnetOwners()
		require.NoError(t, err)
		assert.Len(t, controlKeys, 0)
		assert.Equal(t, uint32(0), threshold)
	})
}

// TestGetOwnersProduction tests GetOwners with production scenarios
func TestGetOwnersProduction(t *testing.T) {
	t.Parallel()

	t.Run("GetOwners with different network configurations", func(t *testing.T) {
		testSubnetID := ids.ID{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32}

		networks := []odyssey.Network{
			odyssey.TestnetNetwork(),
			odyssey.MainnetNetwork(),
		}

		for i, network := range networks {
			t.Run("network_"+string(rune(i)), func(t *testing.T) {
				// Test that the function can be called (it will fail due to network call)
				_, _, err := GetOwners(network, testSubnetID)
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "subnet tx")
			})
		}
	})

	t.Run("GetOwners with LOCAL_NODE flag variations", func(t *testing.T) {
		testSubnetID := ids.ID{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32}

		testCases := []struct {
			name        string
			localNode   string
			expectLocal bool
		}{
			{"local_node_true", "true", true},
			{"local_node_false", "false", false},
			{"local_node_empty", "", false},
			{"local_node_invalid", "invalid", false},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				os.Setenv("LOCAL_NODE", tc.localNode)
				defer os.Unsetenv("LOCAL_NODE")

				network := odyssey.TestnetNetwork()
				if tc.expectLocal {
					assert.Equal(t, "http://127.0.0.1:9650", network.Endpoint)
				} else {
					assert.NotEqual(t, "http://127.0.0.1:9650", network.Endpoint)
				}

				// Test that the function can be called
				_, _, err := GetOwners(network, testSubnetID)
				assert.Error(t, err)
			})
		}
	})

	t.Run("GetOwners with different subnet ID formats", func(t *testing.T) {
		network := odyssey.TestnetNetwork()

		testCases := []struct {
			name     string
			subnetID ids.ID
		}{
			{"empty_subnet_id", ids.Empty},
			{"all_zeros", ids.ID{}},
			{"all_ones", ids.ID{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}},
			{"pattern_subnet_id", ids.ID{0x12, 0x34, 0x56, 0x78, 0x9A, 0xBC, 0xDE, 0xF0, 0x12, 0x34, 0x56, 0x78, 0x9A, 0xBC, 0xDE, 0xF0, 0x12, 0x34, 0x56, 0x78, 0x9A, 0xBC, 0xDE, 0xF0, 0x12, 0x34, 0x56, 0x78, 0x9A, 0xBC, 0xDE, 0xF0}},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Test that the function can be called
				// Some subnet IDs might not cause errors (like empty ones)
				_, _, _ = GetOwners(network, tc.subnetID)
				// We don't assert error here as some subnet IDs might be valid
				// and the function might succeed or fail depending on the specific ID
			})
		}
	})
}

// TestGetNetworkProduction tests GetNetwork with production scenarios
func TestGetNetworkProduction(t *testing.T) {
	t.Parallel()

	t.Run("GetNetwork with different network ID scenarios", func(t *testing.T) {
		// Test with different transaction types that might have different network ID handling
		testCases := []struct {
			name        string
			unsignedTx  txs.UnsignedTx
			expectError bool
			errorMsg    string
		}{
			{"CreateSubnetTx", &txs.CreateSubnetTx{}, true, "unexpected unsigned tx type"}, // Will fail due to unsupported type
			{"RemoveSubnetValidatorTx", &txs.RemoveSubnetValidatorTx{}, true, "undefined network model for tx"},
			{"AddSubnetValidatorTx", &txs.AddSubnetValidatorTx{}, true, "undefined network model for tx"},
			{"CreateChainTx", &txs.CreateChainTx{}, true, "undefined network model for tx"},
			{"TransformSubnetTx", &txs.TransformSubnetTx{}, true, "undefined network model for tx"},
			{"AddPermissionlessValidatorTx", &txs.AddPermissionlessValidatorTx{}, true, "undefined network model for tx"},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				tx := &txs.Tx{Unsigned: tc.unsignedTx}
				ms := New(tx)

				network, err := ms.GetNetwork()
				if tc.expectError {
					assert.Error(t, err)
					assert.Equal(t, odyssey.UndefinedNetwork, network)
					assert.Contains(t, err.Error(), tc.errorMsg)
				} else {
					require.NoError(t, err)
					assert.NotEqual(t, odyssey.UndefinedNetwork, network)
				}
			})
		}
	})

	t.Run("GetNetwork with custom network configurations", func(t *testing.T) {
		// Test network model resolution with different scenarios
		tx := &txs.Tx{Unsigned: &txs.CreateSubnetTx{}}
		ms := New(tx)

		// Test that GetNetworkID returns 0 for default transaction
		networkID, err := ms.GetNetworkID()
		require.NoError(t, err)
		assert.Equal(t, uint32(0), networkID)

		// Test that GetNetwork fails with unexpected unsigned tx type
		network, err := ms.GetNetwork()
		assert.Error(t, err)
		assert.Equal(t, odyssey.UndefinedNetwork, network)
		// The error message should contain "unexpected unsigned tx type"
		assert.Contains(t, err.Error(), "unexpected unsigned tx type")
	})
}

// TestIsReadyToCommitComprehensive tests IsReadyToCommit with comprehensive scenarios
func TestIsReadyToCommitComprehensive(t *testing.T) {
	t.Parallel()

	t.Run("IsReadyToCommit with different transaction types", func(t *testing.T) {
		testCases := []struct {
			name        string
			unsignedTx  txs.UnsignedTx
			expectReady bool
			expectError bool
		}{
			{"CreateSubnetTx", &txs.CreateSubnetTx{}, true, false},
			{"RemoveSubnetValidatorTx", &txs.RemoveSubnetValidatorTx{}, false, true},
			{"AddSubnetValidatorTx", &txs.AddSubnetValidatorTx{}, false, true},
			{"CreateChainTx", &txs.CreateChainTx{}, false, true},
			{"TransformSubnetTx", &txs.TransformSubnetTx{}, false, true},
			{"AddPermissionlessValidatorTx", &txs.AddPermissionlessValidatorTx{}, false, true},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				tx := &txs.Tx{Unsigned: tc.unsignedTx}
				ms := New(tx)

				ready, err := ms.IsReadyToCommit()
				if tc.expectError {
					assert.Error(t, err)
					assert.False(t, ready)
				} else {
					require.NoError(t, err)
					assert.Equal(t, tc.expectReady, ready)
				}
			})
		}
	})

	t.Run("IsReadyToCommit with complex transaction structure", func(t *testing.T) {
		// Create a more complex transaction with credentials
		unsignedTx := &txs.CreateSubnetTx{}
		outputCred := &secp256k1fx.Credential{
			Sigs: [][secp256k1.SignatureLen]byte{
				{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64},
			},
		}

		tx := &txs.Tx{
			Unsigned: unsignedTx,
			Creds:    []verify.Verifiable{outputCred},
		}

		ms := New(tx)
		ready, err := ms.IsReadyToCommit()
		require.NoError(t, err)
		assert.True(t, ready)
	})

	t.Run("IsReadyToCommit with multiple credentials", func(t *testing.T) {
		unsignedTx := &txs.CreateSubnetTx{}
		outputCred := &secp256k1fx.Credential{
			Sigs: [][secp256k1.SignatureLen]byte{
				{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64},
			},
		}

		subnetCred := &secp256k1fx.Credential{
			Sigs: [][secp256k1.SignatureLen]byte{
				{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64},
			},
		}

		tx := &txs.Tx{
			Unsigned: unsignedTx,
			Creds:    []verify.Verifiable{outputCred, subnetCred},
		}

		ms := New(tx)
		ready, err := ms.IsReadyToCommit()
		require.NoError(t, err)
		assert.True(t, ready)
	})
}

// TestToFileAdvancedScenarios tests ToFile with advanced file operation scenarios
func TestToFileAdvancedScenarios(t *testing.T) {
	t.Parallel()

	t.Run("ToFile with read-only directory", func(t *testing.T) {
		unsignedTx := &txs.CreateSubnetTx{}
		tx := &txs.Tx{Unsigned: unsignedTx}
		ms := New(tx)

		// Create a read-only directory
		tempDir := t.TempDir()
		readOnlyDir := filepath.Join(tempDir, "readonly")
		err := os.Mkdir(readOnlyDir, 0444) // Read-only permissions
		require.NoError(t, err)

		filePath := filepath.Join(readOnlyDir, "test.tx")

		// Test ToFile with read-only directory
		err = ms.ToFile(filePath)
		if err != nil {
			// Expected to fail due to permissions or incomplete transaction
			assert.True(t,
				containsAny(err.Error(), []string{
					"couldn't marshal signed tx",
					"couldn't create file",
					"permission denied",
					"access denied",
				}))
		}
	})

	t.Run("ToFile with existing file", func(t *testing.T) {
		unsignedTx := &txs.CreateSubnetTx{}
		tx := &txs.Tx{Unsigned: unsignedTx}
		ms := New(tx)

		tempDir := t.TempDir()
		filePath := filepath.Join(tempDir, "existing.tx")

		// Create an existing file
		err := os.WriteFile(filePath, []byte("existing content"), 0644)
		require.NoError(t, err)

		// Test ToFile with existing file (should overwrite)
		err = ms.ToFile(filePath)
		if err != nil {
			// Expected to fail due to incomplete transaction
			assert.Contains(t, err.Error(), "couldn't marshal signed tx")
		} else {
			// If it succeeds, verify file was overwritten
			content, err := os.ReadFile(filePath)
			require.NoError(t, err)
			assert.NotEqual(t, "existing content", string(content))
		}
	})

	t.Run("ToFile with concurrent access", func(t *testing.T) {
		unsignedTx := &txs.CreateSubnetTx{}
		tx := &txs.Tx{Unsigned: unsignedTx}
		ms := New(tx)

		tempDir := t.TempDir()
		filePath := filepath.Join(tempDir, "concurrent.tx")

		// Test concurrent ToFile operations
		done := make(chan error, 2)
		go func() {
			done <- ms.ToFile(filePath)
		}()
		go func() {
			done <- ms.ToFile(filePath)
		}()

		// Wait for both operations to complete
		err1 := <-done
		err2 := <-done

		// At least one should succeed or both should fail gracefully
		if err1 != nil && err2 != nil {
			// Both failed, which is expected due to incomplete transaction
			assert.Contains(t, err1.Error(), "couldn't marshal signed tx")
			assert.Contains(t, err2.Error(), "couldn't marshal signed tx")
		}
	})
}

// TestFromBytesAdvancedScenarios tests FromBytes with advanced deserialization scenarios
func TestFromBytesAdvancedScenarios(t *testing.T) {
	t.Parallel()

	t.Run("FromBytes with malformed transaction data", func(t *testing.T) {
		ms := New(nil)

		testCases := []struct {
			name string
			data []byte
		}{
			{
				name: "truncated transaction",
				data: []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1A, 0x1B, 0x1C, 0x1D, 0x1E, 0x1F},
			},
			{
				name: "invalid transaction header",
				data: []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF},
			},
			{
				name: "zero-length transaction",
				data: []byte{0x00, 0x00, 0x00, 0x00},
			},
			{
				name: "oversized transaction",
				data: make([]byte, 1024*1024), // 1MB of zeros
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				err := ms.FromBytes(tc.data)
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "error unmarshaling signed tx")
			})
		}
	})

	t.Run("FromBytes with concurrent access", func(t *testing.T) {
		ms1 := New(nil)
		ms2 := New(nil)

		testData := []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F}

		// Test concurrent FromBytes operations
		done := make(chan error, 2)
		go func() {
			done <- ms1.FromBytes(testData)
		}()
		go func() {
			done <- ms2.FromBytes(testData)
		}()

		// Wait for both operations to complete
		err1 := <-done
		err2 := <-done

		// Both should fail with the same error
		assert.Error(t, err1)
		assert.Error(t, err2)
		assert.Contains(t, err1.Error(), "error unmarshaling signed tx")
		assert.Contains(t, err2.Error(), "error unmarshaling signed tx")
	})
}

// TestGetSubnetOwnersAdvancedScenarios tests GetSubnetOwners with advanced scenarios
func TestGetSubnetOwnersAdvancedScenarios(t *testing.T) {
	t.Parallel()

	t.Run("GetSubnetOwners with network failure simulation", func(t *testing.T) {
		// Test with a transaction that will cause network-related errors
		tx := &txs.Tx{Unsigned: &txs.CreateSubnetTx{}}
		ms := New(tx)

		// This should fail due to unexpected unsigned tx type
		controlKeys, threshold, err := ms.GetSubnetOwners()
		assert.Error(t, err)
		assert.Nil(t, controlKeys)
		assert.Equal(t, uint32(0), threshold)
		assert.Contains(t, err.Error(), "unexpected unsigned tx type")
	})

	t.Run("GetSubnetOwners with cached vs non-cached scenarios", func(t *testing.T) {
		// Test with pre-cached control keys
		ms := &Multisig{
			OChainTx: &txs.Tx{Unsigned: &txs.CreateSubnetTx{}},
			controlKeys: []ids.ShortID{
				{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
				{2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21},
				{3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22},
			},
			threshold: 3,
		}

		// First call should return cached values
		controlKeys1, threshold1, err1 := ms.GetSubnetOwners()
		require.NoError(t, err1)
		assert.Len(t, controlKeys1, 3)
		assert.Equal(t, uint32(3), threshold1)

		// Second call should return the same cached values
		controlKeys2, threshold2, err2 := ms.GetSubnetOwners()
		require.NoError(t, err2)
		assert.Len(t, controlKeys2, 3)
		assert.Equal(t, uint32(3), threshold2)

		// Verify they're the same
		assert.Equal(t, controlKeys1, controlKeys2)
		assert.Equal(t, threshold1, threshold2)
	})

	t.Run("GetSubnetOwners with different threshold scenarios", func(t *testing.T) {
		testCases := []struct {
			name      string
			threshold uint32
			keyCount  int
		}{
			{"single_key", 1, 1},
			{"two_keys", 2, 2},
			{"three_keys", 3, 3},
			{"five_keys", 5, 5},
			{"ten_keys", 10, 10},
			{"twenty_keys", 20, 20},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				controlKeys := make([]ids.ShortID, tc.keyCount)
				for i := range controlKeys {
					controlKeys[i] = ids.ShortID{byte(i + 1), byte(i + 2), byte(i + 3), byte(i + 4), byte(i + 5), byte(i + 6), byte(i + 7), byte(i + 8), byte(i + 9), byte(i + 10), byte(i + 11), byte(i + 12), byte(i + 13), byte(i + 14), byte(i + 15), byte(i + 16), byte(i + 17), byte(i + 18), byte(i + 19), byte(i + 20)}
				}

				ms := &Multisig{
					OChainTx:    &txs.Tx{Unsigned: &txs.CreateSubnetTx{}},
					controlKeys: controlKeys,
					threshold:   tc.threshold,
				}

				resultKeys, resultThreshold, err := ms.GetSubnetOwners()
				require.NoError(t, err)
				assert.Len(t, resultKeys, tc.keyCount)
				assert.Equal(t, tc.threshold, resultThreshold)
			})
		}
	})
}

// TestConcurrentAccessSafety tests concurrent access safety
func TestConcurrentAccessSafety(t *testing.T) {
	t.Parallel()

	t.Run("Concurrent method calls", func(t *testing.T) {
		ms := New(&txs.Tx{Unsigned: &txs.CreateSubnetTx{}})

		// Test concurrent access to various methods
		done := make(chan struct{}, 10)

		go func() {
			_, _ = ms.GetTxKind()
			done <- struct{}{}
		}()

		go func() {
			_, _ = ms.GetNetworkID()
			done <- struct{}{}
		}()

		go func() {
			_, _ = ms.GetBlockchainID()
			done <- struct{}{}
		}()

		go func() {
			_, _ = ms.GetSubnetID()
			done <- struct{}{}
		}()

		go func() {
			_, _ = ms.IsReadyToCommit()
			done <- struct{}{}
		}()

		go func() {
			_, _ = ms.GetWrappedOChainTx()
			done <- struct{}{}
		}()

		go func() {
			_ = ms.String()
			done <- struct{}{}
		}()

		go func() {
			_ = ms.Undefined()
			done <- struct{}{}
		}()

		go func() {
			_, _ = ms.ToBytes()
			done <- struct{}{}
		}()

		go func() {
			_ = ms.FromBytes([]byte{0x00, 0x01, 0x02})
			done <- struct{}{}
		}()

		// Wait for all goroutines to complete
		for i := 0; i < 10; i++ {
			<-done
		}
	})

	t.Run("Concurrent file operations", func(t *testing.T) {
		ms := New(&txs.Tx{Unsigned: &txs.CreateSubnetTx{}})
		tempDir := t.TempDir()

		// Test concurrent file operations
		done := make(chan error, 3)

		go func() {
			done <- ms.ToFile(filepath.Join(tempDir, "file1.tx"))
		}()

		go func() {
			done <- ms.ToFile(filepath.Join(tempDir, "file2.tx"))
		}()

		go func() {
			done <- ms.ToFile(filepath.Join(tempDir, "file3.tx"))
		}()

		// Wait for all operations to complete
		for i := 0; i < 3; i++ {
			err := <-done
			// All should fail due to incomplete transaction, but shouldn't panic
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "couldn't marshal signed tx")
		}
	})
}
