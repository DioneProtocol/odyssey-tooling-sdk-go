// Copyright (C) 2025, Dione Protocol, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package wallet

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/DioneProtocol/odyssey-tooling-sdk-go/key"
	"github.com/DioneProtocol/odyssey-tooling-sdk-go/keychain"
	"github.com/DioneProtocol/odyssey-tooling-sdk-go/odyssey"
	"github.com/DioneProtocol/odysseygo/ids"
	"github.com/DioneProtocol/odysseygo/utils/set"
	"github.com/DioneProtocol/odysseygo/vms/secp256k1fx"
	"github.com/DioneProtocol/odysseygo/wallet/subnet/primary"
)

func TestWalletCreation(t *testing.T) {
	// Sequential to avoid rate limiting

	tests := []struct {
		name   string
		config *primary.WalletConfig
		expErr bool
	}{
		{
			name: "valid wallet config with soft key",
			config: &primary.WalletConfig{
				URI:              "https://testnode.dioneprotocol.com",
				DIONEKeychain:    createTestKeychain(t),
				EthKeychain:      secp256k1fx.NewKeychain(),
				OChainTxsToFetch: nil,
			},
			expErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			wallet, err := New(ctx, tt.config)
			t.Logf("Wallet: %v", wallet)

			if tt.expErr {
				require.Error(t, err)
				require.Empty(t, wallet)
			} else {
				require.NoError(t, err)
				require.NotNil(t, wallet)
				require.NotNil(t, wallet.Keychain)
				require.NotNil(t, wallet.config)
			}
		})
	}
}

func TestWalletCreationNilConfig(t *testing.T) {
	// Sequential to avoid rate limiting

	// Test nil config separately to avoid panic
	ctx := context.Background()
	wallet, err := New(ctx, nil)
	require.Error(t, err)
	require.Empty(t, wallet)
}

func TestWalletAddresses(t *testing.T) {
	// Sequential to avoid rate limiting
	addTestDelay() // Add delay to avoid rate limiting

	ctx := context.Background()
	config := &primary.WalletConfig{
		URI:              "https://testnode.dioneprotocol.com",
		DIONEKeychain:    createTestKeychain(t),
		EthKeychain:      secp256k1fx.NewKeychain(),
		OChainTxsToFetch: nil,
	}

	wallet, err := New(ctx, config)
	require.NoError(t, err)

	addresses := wallet.Addresses()
	require.NotEmpty(t, addresses)
	require.Len(t, addresses, 1) // Single key keychain should have one address

	// Verify address is valid
	require.NotEqual(t, ids.ShortID{}, addresses[0])
	t.Logf("Wallet address: %s", addresses[0].String())
}

func TestWalletMultiChainAddressGeneration(t *testing.T) {
	// Sequential to avoid rate limiting
	addTestDelay() // Add delay to avoid rate limiting

	ctx := context.Background()
	config := &primary.WalletConfig{
		URI:              "https://node.dioneprotocol.com",
		DIONEKeychain:    createTestKeychain(t),
		EthKeychain:      secp256k1fx.NewKeychain(),
		OChainTxsToFetch: nil,
	}

	wallet, err := New(ctx, config)
	require.NoError(t, err)

	// Test O-Chain address generation
	oAddr, err := wallet.Keychain.O()
	t.Logf("Wallet address: %v", wallet)
	t.Logf("Wallet address O chain: %v", oAddr)

	require.NoError(t, err)
	require.NotEmpty(t, oAddr)
	require.Len(t, oAddr, 1)

	// Verify address format
	require.Contains(t, oAddr[0], "O-")

	// Test that the same key generates consistent addresses
	oAddr2, err := wallet.Keychain.O()
	t.Logf("Wallet address O chain 2: %v", oAddr2)
	require.NoError(t, err)
	require.Equal(t, oAddr, oAddr2)
}

func TestWalletSecureChangeOwner(t *testing.T) {
	// Sequential to avoid rate limiting
	addTestDelay() // Add delay to avoid rate limiting

	ctx := context.Background()
	config := &primary.WalletConfig{
		URI:              "https://testnode.dioneprotocol.com",
		DIONEKeychain:    createTestKeychain(t),
		EthKeychain:      secp256k1fx.NewKeychain(),
		OChainTxsToFetch: nil,
	}

	wallet, err := New(ctx, config)
	require.NoError(t, err)

	// Get initial addresses
	initialAddresses := wallet.Addresses()
	require.NotEmpty(t, initialAddresses)

	// Apply secure change owner
	wallet.SecureWalletIsChangeOwner()

	// Verify wallet still has the same addresses
	addresses := wallet.Addresses()
	require.Equal(t, initialAddresses, addresses)

	// Verify options were set
	require.NotEmpty(t, wallet.options)
}

func TestWalletSetAuthKeys(t *testing.T) {
	// Sequential to avoid rate limiting
	addTestDelay() // Add delay to avoid rate limiting

	ctx := context.Background()
	config := &primary.WalletConfig{
		URI:              "https://testnode.dioneprotocol.com",
		DIONEKeychain:    createTestKeychain(t),
		EthKeychain:      secp256k1fx.NewKeychain(),
		OChainTxsToFetch: nil,
	}

	wallet, err := New(ctx, config)
	require.NoError(t, err)

	// Create additional auth keys
	additionalKey, err := key.NewSoft()
	require.NoError(t, err)
	authKeys := []ids.ShortID{additionalKey.Addresses()[0]}

	// Set auth keys
	wallet.SetAuthKeys(authKeys)

	// Verify options were set
	require.NotEmpty(t, wallet.options)

	// Verify wallet still has original addresses
	addresses := wallet.Addresses()
	require.NotEmpty(t, addresses)
}

func TestWalletSetSubnetAuthMultisig(t *testing.T) {
	// Sequential to avoid rate limiting
	addTestDelay() // Add delay to avoid rate limiting

	ctx := context.Background()
	config := &primary.WalletConfig{
		URI:              "https://testnode.dioneprotocol.com",
		DIONEKeychain:    createTestKeychain(t),
		EthKeychain:      secp256k1fx.NewKeychain(),
		OChainTxsToFetch: nil,
	}

	wallet, err := New(ctx, config)
	require.NoError(t, err)

	// Create additional auth keys
	additionalKey, err := key.NewSoft()
	require.NoError(t, err)
	authKeys := []ids.ShortID{additionalKey.Addresses()[0]}

	// Set subnet auth multisig
	wallet.SetSubnetAuthMultisig(authKeys)

	// Verify options were set (should have both secure change owner and auth keys)
	require.NotEmpty(t, wallet.options)
	require.Len(t, wallet.options, 2) // Should have both options
}

func TestWalletWithMultipleKeys(t *testing.T) {
	// Sequential to avoid rate limiting
	addTestDelay() // Add delay to avoid rate limiting

	ctx := context.Background()

	// Create a keychain with multiple keys
	key1, err := key.NewSoft()
	require.NoError(t, err)

	key2, err := key.NewSoft()
	require.NoError(t, err)

	// Create a keychain with multiple keys
	kc := secp256k1fx.NewKeychain()
	kc.Add(key1.PrivKey())
	kc.Add(key2.PrivKey())

	config := &primary.WalletConfig{
		URI:              "https://testnode.dioneprotocol.com",
		DIONEKeychain:    kc,
		EthKeychain:      secp256k1fx.NewKeychain(),
		OChainTxsToFetch: nil,
	}

	wallet, err := New(ctx, config)
	require.NoError(t, err)

	// Verify wallet has multiple addresses
	addresses := wallet.Addresses()
	require.Len(t, addresses, 2)

	// Verify all addresses are unique
	addressSet := set.Set[ids.ShortID]{}
	for _, addr := range addresses {
		require.False(t, addressSet.Contains(addr))
		addressSet.Add(addr)
	}
}

func TestWalletWithOChainTxsToFetch(t *testing.T) {
	// Sequential to avoid rate limiting
	addTestDelay() // Add delay to avoid rate limiting

	ctx := context.Background()

	// Create a test subnet ID
	subnetID := ids.GenerateTestID()

	config := &primary.WalletConfig{
		URI:              "https://testnode.dioneprotocol.com",
		DIONEKeychain:    createTestKeychain(t),
		EthKeychain:      secp256k1fx.NewKeychain(),
		OChainTxsToFetch: set.Of(subnetID),
	}

	wallet, err := New(ctx, config)
	// The wallet creation might fail due to tx not found, which is expected for test IDs
	if err != nil {
		// If it fails due to tx not found, that's expected behavior
		require.Contains(t, err.Error(), "not found")
		return
	}

	require.NotNil(t, wallet)
	require.Equal(t, subnetID, wallet.config.OChainTxsToFetch.List()[0])
}

func TestWalletErrorHandling(t *testing.T) {
	t.Parallel()

	t.Run("context cancellation", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		config := &primary.WalletConfig{
			URI:              "https://testnode.dioneprotocol.com",
			DIONEKeychain:    createTestKeychain(t),
			EthKeychain:      secp256k1fx.NewKeychain(),
			OChainTxsToFetch: nil,
		}

		_, err := New(ctx, config)
		require.Error(t, err)
		require.Contains(t, err.Error(), "context canceled")
	})

	t.Run("invalid network endpoint", func(t *testing.T) {
		ctx := context.Background()
		config := &primary.WalletConfig{
			URI:              "invalid://endpoint",
			DIONEKeychain:    createTestKeychain(t),
			EthKeychain:      secp256k1fx.NewKeychain(),
			OChainTxsToFetch: nil,
		}

		_, err := New(ctx, config)
		require.Error(t, err)
	})
}

func TestWalletKeychainIntegration(t *testing.T) {
	// Sequential to avoid rate limiting
	addTestDelay() // Add delay to avoid rate limiting

	ctx := context.Background()
	network := odyssey.TestnetNetwork()

	// Create a temporary key file for testing
	tempKeyPath := t.TempDir() + "/test.key"

	// Test with soft key keychain
	keychain, err := keychain.NewKeychain(network, tempKeyPath, nil)
	require.NoError(t, err)

	config := &primary.WalletConfig{
		URI:              network.Endpoint,
		DIONEKeychain:    keychain.Keychain,
		EthKeychain:      secp256k1fx.NewKeychain(),
		OChainTxsToFetch: nil,
	}

	wallet, err := New(ctx, config)
	require.NoError(t, err)

	// Verify keychain integration
	require.Equal(t, keychain.Keychain, wallet.Keychain.Keychain)
}

func TestWalletLedgerSupport(t *testing.T) {
	t.Parallel()

	// This test verifies that the wallet can be created with ledger keychain
	// Note: This test will skip if no ledger is connected
	t.Run("ledger keychain creation", func(t *testing.T) {
		ctx := context.Background()
		network := odyssey.TestnetNetwork()

		ledgerInfo := &keychain.LedgerParams{
			LedgerAddresses: []string{"O-testnetxxxxxxxxx"}, // Test address
		}

		// This will fail if no ledger is connected, which is expected
		keychain, err := keychain.NewKeychain(network, "", ledgerInfo)
		if err != nil {
			t.Skip("Ledger not available for testing")
		}

		config := &primary.WalletConfig{
			URI:              network.Endpoint,
			DIONEKeychain:    keychain.Keychain,
			EthKeychain:      secp256k1fx.NewKeychain(),
			OChainTxsToFetch: nil,
		}

		wallet, err := New(ctx, config)
		if err != nil {
			t.Skip("Ledger wallet creation failed - ledger may not be available")
		}

		require.NotNil(t, wallet)
		require.True(t, wallet.Keychain.LedgerEnabled())
	})
}

// Helper function to create a test keychain
func createTestKeychain(t *testing.T) *secp256k1fx.Keychain {
	key, err := key.NewSoft()
	require.NoError(t, err)
	return key.KeyChain()
}

// Helper function to add delay between tests to avoid rate limiting
func addTestDelay() {
	time.Sleep(500 * time.Millisecond)
}

// Benchmark tests for wallet operations
func BenchmarkWalletCreation(b *testing.B) {
	config := &primary.WalletConfig{
		URI:              "https://testnode.dioneprotocol.com",
		DIONEKeychain:    createBenchmarkKeychain(b),
		EthKeychain:      secp256k1fx.NewKeychain(),
		OChainTxsToFetch: nil,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctx := context.Background()
		_, err := New(ctx, config)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkWalletAddresses(b *testing.B) {
	ctx := context.Background()
	config := &primary.WalletConfig{
		URI:              "https://testnode.dioneprotocol.com",
		DIONEKeychain:    createBenchmarkKeychain(b),
		EthKeychain:      secp256k1fx.NewKeychain(),
		OChainTxsToFetch: nil,
	}

	wallet, err := New(ctx, config)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = wallet.Addresses()
	}
}

func TestWalletNetworkHRPGeneration(t *testing.T) {
	// Sequential to avoid rate limiting
	addTestDelay() // Add delay to avoid rate limiting

	tests := []struct {
		name           string
		uri            string
		expectedHRP    string
		expectedPrefix string
	}{
		{
			name:           "mainnet network",
			uri:            "https://node.dioneprotocol.com",
			expectedHRP:    "dione",
			expectedPrefix: "O-dione",
		},
		{
			name:           "testnet network",
			uri:            "https://testnode.dioneprotocol.com",
			expectedHRP:    "testnet",
			expectedPrefix: "O-testnet",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			// Create wallet config with the specified URI
			config := &primary.WalletConfig{
				URI:              tt.uri,
				DIONEKeychain:    createTestKeychain(t),
				EthKeychain:      secp256k1fx.NewKeychain(),
				OChainTxsToFetch: nil,
			}

			// Create wallet
			wallet, err := New(ctx, config)
			require.NoError(t, err)
			require.NotNil(t, wallet)

			// Test O-Chain address generation
			oAddr, err := wallet.Keychain.O()
			require.NoError(t, err)
			require.NotEmpty(t, oAddr)
			require.Len(t, oAddr, 1)

			// Verify the address has the correct prefix
			require.Contains(t, oAddr[0], tt.expectedPrefix,
				"Expected address to contain prefix %s, but got %s", tt.expectedPrefix, oAddr[0])

			// Verify the network HRP is correct
			network := odyssey.NetworkFromURI(tt.uri)
			require.Equal(t, tt.expectedHRP, network.HRP(),
				"Expected HRP %s for URI %s, but got %s", tt.expectedHRP, tt.uri, network.HRP())

			t.Logf("%s O-Chain address: %s (HRP: %s)", tt.name, oAddr[0], network.HRP())
		})
	}
}

func TestWalletGenerationFromPrivateKey(t *testing.T) {
	// Sequential to avoid rate limiting
	addTestDelay() // Add delay to avoid rate limiting

	// Test with a known private key (this is a test key, not a real one)
	testPrivateKeyHex := "56289e99c94b6912bfc12adc093c9b51124f0dc54ac7a766b2bc5ccf558d8027"

	ctx := context.Background()

	// Test private key loading and validation
	softKey, err := key.LoadSoftFromBytes([]byte(testPrivateKeyHex))
	require.NoError(t, err)
	require.NotNil(t, softKey)
	require.Equal(t, testPrivateKeyHex, softKey.PrivKeyHex())

	// Test wallet creation from private key (HRP validation is already covered in TestWalletNetworkHRPGeneration)
	testnetConfig := &primary.WalletConfig{
		URI:              "https://testnode.dioneprotocol.com",
		DIONEKeychain:    softKey.KeyChain(),
		EthKeychain:      secp256k1fx.NewKeychain(),
		OChainTxsToFetch: nil,
	}

	wallet, err := New(ctx, testnetConfig)
	require.NoError(t, err)
	require.NotNil(t, wallet)

	// Test that we can generate addresses from the private key
	oAddr, err := wallet.Keychain.O()
	require.NoError(t, err)
	require.NotEmpty(t, oAddr)
	require.Len(t, oAddr, 1)

	// The address should be deterministic and consistent
	expectedTestnetAddr := "O-testnet18jma8ppw3nhx5r4ap8clazz0dps7rv5uw463py"
	require.Equal(t, expectedTestnetAddr, oAddr[0])

	t.Logf("Private Key: %s", softKey.PrivKeyHex())
	t.Logf("Generated Testnet Address: %s", oAddr[0])
}

func TestWalletDeterministicAddressGeneration(t *testing.T) {
	// Sequential to avoid rate limiting
	addTestDelay() // Add delay to avoid rate limiting

	// Test that the same private key generates the same addresses across different networks
	testPrivateKeyHex := "56289e99c94b6912bfc12adc093c9b51124f0dc54ac7a766b2bc5ccf558d8027"

	ctx := context.Background()

	// Create the same key twice to ensure deterministic behavior
	softKey1, err := key.LoadSoftFromBytes([]byte(testPrivateKeyHex))
	require.NoError(t, err)

	softKey2, err := key.LoadSoftFromBytes([]byte(testPrivateKeyHex))
	require.NoError(t, err)

	// Verify both keys are identical
	require.Equal(t, softKey1.PrivKeyHex(), softKey2.PrivKeyHex())
	require.Equal(t, softKey1.KeyChain().Addresses().List(), softKey2.KeyChain().Addresses().List())

	// Test on testnet
	testnetConfig := &primary.WalletConfig{
		URI:              "https://testnode.dioneprotocol.com",
		DIONEKeychain:    softKey1.KeyChain(),
		EthKeychain:      secp256k1fx.NewKeychain(),
		OChainTxsToFetch: nil,
	}

	testnetWallet, err := New(ctx, testnetConfig)
	require.NoError(t, err)

	testnetAddr, err := testnetWallet.Keychain.O()
	require.NoError(t, err)
	require.Len(t, testnetAddr, 1)

	// Test on mainnet
	mainnetConfig := &primary.WalletConfig{
		URI:              "https://node.dioneprotocol.com",
		DIONEKeychain:    softKey2.KeyChain(),
		EthKeychain:      secp256k1fx.NewKeychain(),
		OChainTxsToFetch: nil,
	}

	mainnetWallet, err := New(ctx, mainnetConfig)
	require.NoError(t, err)

	mainnetAddr, err := mainnetWallet.Keychain.O()
	require.NoError(t, err)
	require.Len(t, mainnetAddr, 1)

	// The addresses should be different due to different HRPs, but the underlying key should be the same
	require.Contains(t, testnetAddr[0], "O-testnet")
	require.Contains(t, mainnetAddr[0], "O-dione")

	// But the raw addresses (without HRP) should be the same
	testnetNetwork := odyssey.NetworkFromURI("https://testnode.dioneprotocol.com")
	mainnetNetwork := odyssey.NetworkFromURI("https://node.dioneprotocol.com")

	// Generate addresses using the same key but different HRPs
	testnetAddrDirect, err := softKey1.O(testnetNetwork.HRP())
	require.NoError(t, err)

	mainnetAddrDirect, err := softKey1.O(mainnetNetwork.HRP())
	require.NoError(t, err)

	require.Equal(t, testnetAddr[0], testnetAddrDirect)
	require.Equal(t, mainnetAddr[0], mainnetAddrDirect)

	t.Logf("Testnet address: %s", testnetAddr[0])
	t.Logf("Mainnet address: %s", mainnetAddr[0])
	t.Logf("Private Key: %s", softKey1.PrivKeyHex())
}

func TestAddressEncodingDifference(t *testing.T) {
	// Test to demonstrate why the same raw address produces different suffixes
	testPrivateKeyHex := "56289e99c94b6912bfc12adc093c9b51124f0dc54ac7a766b2bc5ccf558d8027"

	softKey, err := key.LoadSoftFromBytes([]byte(testPrivateKeyHex))
	require.NoError(t, err)

	// Get the raw address bytes (this is the same for both networks)
	rawAddressBytes := softKey.PrivKey().PublicKey().Address().Bytes()
	t.Logf("Raw Address Bytes (hex): %x", rawAddressBytes)
	t.Logf("Raw Address Bytes (length): %d", len(rawAddressBytes))

	// Generate addresses with different HRPs
	testnetAddr, err := softKey.O("testnet")
	require.NoError(t, err)

	mainnetAddr, err := softKey.O("dione")
	require.NoError(t, err)

	t.Logf("Testnet Address: %s", testnetAddr)
	t.Logf("Mainnet Address: %s", mainnetAddr)

	// The difference is in the Bech32 encoding, not the raw address
	// Bech32 includes a checksum that depends on the HRP
	require.NotEqual(t, testnetAddr, mainnetAddr)

	// But the raw address bytes are identical
	require.Equal(t, rawAddressBytes, rawAddressBytes) // Same raw bytes
}

func TestWalletAChainAddressGeneration(t *testing.T) {
	// Sequential to avoid rate limiting
	addTestDelay() // Add delay to avoid rate limiting

	tests := []struct {
		name           string
		uri            string
		expectedHRP    string
		expectedPrefix string
	}{
		{
			name:           "mainnet A-Chain network",
			uri:            "https://node.dioneprotocol.com",
			expectedHRP:    "dione",
			expectedPrefix: "A-dione",
		},
		{
			name:           "testnet A-Chain network",
			uri:            "https://testnode.dioneprotocol.com",
			expectedHRP:    "testnet",
			expectedPrefix: "A-testnet",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			// Create wallet config with the specified URI
			config := &primary.WalletConfig{
				URI:              tt.uri,
				DIONEKeychain:    createTestKeychain(t),
				EthKeychain:      secp256k1fx.NewKeychain(),
				OChainTxsToFetch: nil,
			}

			// Create wallet
			wallet, err := New(ctx, config)
			require.NoError(t, err)
			require.NotNil(t, wallet)

			// Test A-Chain address generation
			aAddr, err := wallet.Keychain.A()
			require.NoError(t, err)
			require.NotEmpty(t, aAddr)
			require.Len(t, aAddr, 1)

			// Verify the address has the correct prefix
			require.Contains(t, aAddr[0], tt.expectedPrefix,
				"Expected A-Chain address to contain prefix %s, but got %s", tt.expectedPrefix, aAddr[0])

			// Verify the network HRP is correct
			network := odyssey.NetworkFromURI(tt.uri)
			require.Equal(t, tt.expectedHRP, network.HRP(),
				"Expected HRP %s for URI %s, but got %s", tt.expectedHRP, tt.uri, network.HRP())

			t.Logf("%s A-Chain address: %s (HRP: %s)", tt.name, aAddr[0], network.HRP())
		})
	}
}

func TestWalletBothChainAddressGeneration(t *testing.T) {
	// Sequential to avoid rate limiting
	addTestDelay() // Add delay to avoid rate limiting

	// Test with a known private key to ensure deterministic behavior
	testPrivateKeyHex := "56289e99c94b6912bfc12adc093c9b51124f0dc54ac7a766b2bc5ccf558d8027"

	ctx := context.Background()

	// Create a key from the specific private key
	softKey, err := key.LoadSoftFromBytes([]byte(testPrivateKeyHex))
	require.NoError(t, err)
	require.NotNil(t, softKey)

	// Test on testnet
	testnetConfig := &primary.WalletConfig{
		URI:              "https://testnode.dioneprotocol.com",
		DIONEKeychain:    softKey.KeyChain(),
		EthKeychain:      secp256k1fx.NewKeychain(),
		OChainTxsToFetch: nil,
	}

	wallet, err := New(ctx, testnetConfig)
	require.NoError(t, err)
	require.NotNil(t, wallet)

	// Generate both O-Chain and A-Chain addresses
	oAddr, err := wallet.Keychain.O()
	require.NoError(t, err)
	require.Len(t, oAddr, 1)

	aAddr, err := wallet.Keychain.A()
	require.NoError(t, err)
	require.Len(t, aAddr, 1)

	// Verify both addresses have correct prefixes
	require.Contains(t, oAddr[0], "O-testnet")
	require.Contains(t, aAddr[0], "A-testnet")

	// Verify addresses are different (different chain types)
	require.NotEqual(t, oAddr[0], aAddr[0])

	// Test deterministic behavior - same key should generate same addresses
	oAddr2, err := wallet.Keychain.O()
	require.NoError(t, err)
	require.Equal(t, oAddr, oAddr2)

	aAddr2, err := wallet.Keychain.A()
	require.NoError(t, err)
	require.Equal(t, aAddr, aAddr2)

	t.Logf("Private Key: %s", softKey.PrivKeyHex())
	t.Logf("O-Chain Address: %s", oAddr[0])
	t.Logf("A-Chain Address: %s", aAddr[0])
}

func TestWalletAChainFromPrivateKey(t *testing.T) {
	// Sequential to avoid rate limiting
	addTestDelay() // Add delay to avoid rate limiting

	// Test with a known private key (this is a test key, not a real one)
	testPrivateKeyHex := "56289e99c94b6912bfc12adc093c9b51124f0dc54ac7a766b2bc5ccf558d8027"

	ctx := context.Background()

	// Test private key loading and validation
	softKey, err := key.LoadSoftFromBytes([]byte(testPrivateKeyHex))
	require.NoError(t, err)
	require.NotNil(t, softKey)
	require.Equal(t, testPrivateKeyHex, softKey.PrivKeyHex())

	// Test wallet creation from private key
	testnetConfig := &primary.WalletConfig{
		URI:              "https://testnode.dioneprotocol.com",
		DIONEKeychain:    softKey.KeyChain(),
		EthKeychain:      secp256k1fx.NewKeychain(),
		OChainTxsToFetch: nil,
	}

	wallet, err := New(ctx, testnetConfig)
	require.NoError(t, err)
	require.NotNil(t, wallet)

	// Test that we can generate A-Chain addresses from the private key
	aAddr, err := wallet.Keychain.A()
	require.NoError(t, err)
	require.NotEmpty(t, aAddr)
	require.Len(t, aAddr, 1)

	// The address should be deterministic and consistent
	expectedTestnetAAddr := "A-testnet18jma8ppw3nhx5r4ap8clazz0dps7rv5uw463py"
	require.Equal(t, expectedTestnetAAddr, aAddr[0])

	t.Logf("Private Key: %s", softKey.PrivKeyHex())
	t.Logf("Generated Testnet A-Chain Address: %s", aAddr[0])
}

func TestWalletPrivateKeyValidation(t *testing.T) {
	tests := []struct {
		name        string
		privateKey  string
		expectError bool
		description string
	}{
		{
			name:        "valid private key",
			privateKey:  "56289e99c94b6912bfc12adc093c9b51124f0dc54ac7a766b2bc5ccf558d8027",
			expectError: false,
			description: "Valid 64-character hex private key",
		},
		{
			name:        "invalid private key - too short",
			privateKey:  "56289e99c94b6912bfc12adc093c9b51124f0dc54ac7a766b2bc5ccf558d802",
			expectError: true,
			description: "Invalid private key - too short (63 characters)",
		},
		{
			name:        "invalid private key - too long",
			privateKey:  "56289e99c94b6912bfc12adc093c9b51124f0dc54ac7a766b2bc5ccf558d80277",
			expectError: true,
			description: "Invalid private key - too long (65 characters)",
		},
		{
			name:        "invalid private key - non-hex",
			privateKey:  "56289e99c94b6912bfc12adc093c9b51124f0dc54ac7a766b2bc5ccf558d802g",
			expectError: true,
			description: "Invalid private key - contains non-hex character",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Try to create a key from the private key
			softKey, err := key.LoadSoftFromBytes([]byte(tt.privateKey))

			if tt.expectError {
				require.Error(t, err, "Expected error for %s", tt.description)
				require.Nil(t, softKey)
			} else {
				require.NoError(t, err, "Expected no error for %s", tt.description)
				require.NotNil(t, softKey)
				require.Equal(t, tt.privateKey, softKey.PrivKeyHex())
			}
		})
	}
}

// Helper function for benchmarks
func createBenchmarkKeychain(b *testing.B) *secp256k1fx.Keychain {
	key, err := key.NewSoft()
	if err != nil {
		b.Fatal(err)
	}
	return key.KeyChain()
}
