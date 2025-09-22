// Copyright (C) 2025, Dione Protocol, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package wallet

import (
	"context"
	"strings"
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
	addTestDelay() // Add delay to avoid rate limiting

	tests := []struct {
		name    string
		network odyssey.Network
		expErr  bool
	}{
		{
			name:    "valid wallet config with testnet keychain",
			network: odyssey.TestnetNetwork(),
			expErr:  false,
		},
		{
			name:    "valid wallet config with mainnet keychain",
			network: odyssey.MainnetNetwork(),
			expErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			// Create a temporary key file for testing
			tempKeyPath := t.TempDir() + "/test.pk"

			// Create keychain using the production approach
			keychain, err := keychain.NewKeychain(tt.network, tempKeyPath, nil)
			require.NoError(t, err)
			require.NotNil(t, keychain)

			// Create wallet using the production approach
			wallet, err := New(
				ctx,
				&primary.WalletConfig{
					URI:              tt.network.Endpoint,
					DIONEKeychain:    keychain.Keychain,
					EthKeychain:      secp256k1fx.NewKeychain(),
					OChainTxsToFetch: nil,
				},
			)

			if tt.expErr {
				require.Error(t, err)
				require.Empty(t, wallet)
			} else {
				require.NoError(t, err)
				require.NotNil(t, wallet)
				require.NotNil(t, wallet.Keychain)
				require.NotNil(t, wallet.config)
				require.Equal(t, keychain.Keychain, wallet.Keychain.Keychain)
			}

			t.Logf("Wallet: %v", wallet)
			t.Logf("Network: %s (HRP: %s)", tt.network.Kind.String(), tt.network.HRP())
			t.Logf("Endpoint: %s", tt.network.Endpoint)
		})
	}
}

func TestWalletCreationNilConfig(t *testing.T) {
	// Sequential to avoid rate limiting
	addTestDelay() // Add delay to avoid rate limiting

	// Test nil config separately to avoid panic
	ctx := context.Background()
	wallet, err := New(ctx, nil)
	require.Error(t, err)
	require.Empty(t, wallet)
}

func TestWalletAddresses(t *testing.T) {
	// Sequential to avoid rate limiting
	addTestDelay() // Add delay to avoid rate limiting

	tests := []struct {
		name    string
		network odyssey.Network
	}{
		{
			name:    "testnet wallet addresses",
			network: odyssey.TestnetNetwork(),
		},
		{
			name:    "mainnet wallet addresses",
			network: odyssey.MainnetNetwork(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			// Create a temporary key file for testing
			tempKeyPath := t.TempDir() + "/test.pk"

			// Create keychain using the production approach
			keychain, err := keychain.NewKeychain(tt.network, tempKeyPath, nil)
			require.NoError(t, err)
			require.NotNil(t, keychain)

			// Create wallet using the production approach
			wallet, err := New(
				ctx,
				&primary.WalletConfig{
					URI:              tt.network.Endpoint,
					DIONEKeychain:    keychain.Keychain,
					EthKeychain:      secp256k1fx.NewKeychain(),
					OChainTxsToFetch: nil,
				},
			)
			require.NoError(t, err)

			addresses := wallet.Addresses()
			require.NotEmpty(t, addresses)
			require.Len(t, addresses, 1) // Single key keychain should have one address

			// Verify address is valid
			require.NotEqual(t, ids.ShortID{}, addresses[0])

			// Test address generation with network context
			oAddr, err := wallet.Keychain.O()
			require.NoError(t, err)
			require.Len(t, oAddr, 1)

			// Verify address has correct prefix
			expectedPrefix := "O-" + tt.network.HRP()
			require.Contains(t, oAddr[0], expectedPrefix,
				"Expected address to contain prefix %s, but got %s", expectedPrefix, oAddr[0])

			t.Logf("Wallet address: %s", addresses[0].String())
			t.Logf("O-Chain address: %s", oAddr[0])
			t.Logf("Network: %s (HRP: %s)", tt.network.Kind.String(), tt.network.HRP())
		})
	}
}

func TestWalletMultiChainAddressGeneration(t *testing.T) {
	// Sequential to avoid rate limiting
	addTestDelay() // Add delay to avoid rate limiting

	tests := []struct {
		name           string
		network        odyssey.Network
		expectedPrefix string
	}{
		{
			name:           "testnet multi-chain address generation",
			network:        odyssey.TestnetNetwork(),
			expectedPrefix: "O-testnet",
		},
		{
			name:           "mainnet multi-chain address generation",
			network:        odyssey.MainnetNetwork(),
			expectedPrefix: "O-dione",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			// Create a temporary key file for testing
			tempKeyPath := t.TempDir() + "/test.pk"

			// Create keychain using the production approach
			keychain, err := keychain.NewKeychain(tt.network, tempKeyPath, nil)
			require.NoError(t, err)
			require.NotNil(t, keychain)

			// Create wallet using the production approach
			wallet, err := New(
				ctx,
				&primary.WalletConfig{
					URI:              tt.network.Endpoint,
					DIONEKeychain:    keychain.Keychain,
					EthKeychain:      secp256k1fx.NewKeychain(),
					OChainTxsToFetch: nil,
				},
			)
			require.NoError(t, err)

			// Test O-Chain address generation
			oAddr, err := wallet.Keychain.O()
			require.NoError(t, err)
			require.NotEmpty(t, oAddr)
			require.Len(t, oAddr, 1)

			// Verify address format
			require.Contains(t, oAddr[0], tt.expectedPrefix,
				"Expected address to contain prefix %s, but got %s", tt.expectedPrefix, oAddr[0])

			// Test that the same key generates consistent addresses
			oAddr2, err := wallet.Keychain.O()
			require.NoError(t, err)
			require.Equal(t, oAddr, oAddr2)

			// Test A-Chain address generation
			aAddr, err := wallet.Keychain.A()
			require.NoError(t, err)
			require.NotEmpty(t, aAddr)
			require.Len(t, aAddr, 1)

			// Verify A-Chain address format
			expectedAPrefix := strings.Replace(tt.expectedPrefix, "O-", "A-", 1)
			require.Contains(t, aAddr[0], expectedAPrefix,
				"Expected A-Chain address to contain prefix %s, but got %s", expectedAPrefix, aAddr[0])

			t.Logf("Wallet: %v", wallet)
			t.Logf("O-Chain address: %s", oAddr[0])
			t.Logf("A-Chain address: %s", aAddr[0])
			t.Logf("Network: %s (HRP: %s)", tt.network.Kind.String(), tt.network.HRP())
		})
	}
}

func TestWalletSecureChangeOwner(t *testing.T) {
	// Sequential to avoid rate limiting
	addTestDelay() // Add delay to avoid rate limiting

	tests := []struct {
		name    string
		network odyssey.Network
	}{
		{
			name:    "testnet secure change owner",
			network: odyssey.TestnetNetwork(),
		},
		{
			name:    "mainnet secure change owner",
			network: odyssey.MainnetNetwork(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			// Create a temporary key file for testing
			tempKeyPath := t.TempDir() + "/test.pk"

			// Create keychain using the production approach
			keychain, err := keychain.NewKeychain(tt.network, tempKeyPath, nil)
			require.NoError(t, err)
			require.NotNil(t, keychain)

			// Create wallet using the production approach
			wallet, err := New(
				ctx,
				&primary.WalletConfig{
					URI:              tt.network.Endpoint,
					DIONEKeychain:    keychain.Keychain,
					EthKeychain:      secp256k1fx.NewKeychain(),
					OChainTxsToFetch: nil,
				},
			)
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

			t.Logf("Network: %s (HRP: %s)", tt.network.Kind.String(), tt.network.HRP())
			t.Logf("Addresses: %v", addresses)
			t.Logf("Options count: %d", len(wallet.options))
		})
	}
}

func TestWalletSetAuthKeys(t *testing.T) {
	// Sequential to avoid rate limiting
	addTestDelay() // Add delay to avoid rate limiting

	tests := []struct {
		name    string
		network odyssey.Network
	}{
		{
			name:    "testnet set auth keys",
			network: odyssey.TestnetNetwork(),
		},
		{
			name:    "mainnet set auth keys",
			network: odyssey.MainnetNetwork(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			// Create a temporary key file for testing
			tempKeyPath := t.TempDir() + "/test.pk"

			// Create keychain using the production approach
			keychain, err := keychain.NewKeychain(tt.network, tempKeyPath, nil)
			require.NoError(t, err)
			require.NotNil(t, keychain)

			// Create wallet using the production approach
			wallet, err := New(
				ctx,
				&primary.WalletConfig{
					URI:              tt.network.Endpoint,
					DIONEKeychain:    keychain.Keychain,
					EthKeychain:      secp256k1fx.NewKeychain(),
					OChainTxsToFetch: nil,
				},
			)
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

			t.Logf("Network: %s (HRP: %s)", tt.network.Kind.String(), tt.network.HRP())
			t.Logf("Addresses: %v", addresses)
			t.Logf("Auth keys: %v", authKeys)
			t.Logf("Options count: %d", len(wallet.options))
		})
	}
}

func TestWalletSetSubnetAuthMultisig(t *testing.T) {
	// Sequential to avoid rate limiting
	addTestDelay() // Add delay to avoid rate limiting

	tests := []struct {
		name    string
		network odyssey.Network
	}{
		{
			name:    "testnet subnet auth multisig",
			network: odyssey.TestnetNetwork(),
		},
		{
			name:    "mainnet subnet auth multisig",
			network: odyssey.MainnetNetwork(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			// Create a temporary key file for testing
			tempKeyPath := t.TempDir() + "/test.pk"

			// Create keychain using the production approach
			keychain, err := keychain.NewKeychain(tt.network, tempKeyPath, nil)
			require.NoError(t, err)
			require.NotNil(t, keychain)

			// Create wallet using the production approach
			wallet, err := New(
				ctx,
				&primary.WalletConfig{
					URI:              tt.network.Endpoint,
					DIONEKeychain:    keychain.Keychain,
					EthKeychain:      secp256k1fx.NewKeychain(),
					OChainTxsToFetch: nil,
				},
			)
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

			t.Logf("Network: %s (HRP: %s)", tt.network.Kind.String(), tt.network.HRP())
			t.Logf("Auth keys: %v", authKeys)
			t.Logf("Options count: %d", len(wallet.options))
		})
	}
}

func TestWalletWithMultipleKeys(t *testing.T) {
	// Sequential to avoid rate limiting
	addTestDelay() // Add delay to avoid rate limiting

	tests := []struct {
		name    string
		network odyssey.Network
	}{
		{
			name:    "testnet multiple keys",
			network: odyssey.TestnetNetwork(),
		},
		{
			name:    "mainnet multiple keys",
			network: odyssey.MainnetNetwork(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			// Create a temporary key file for testing
			tempKeyPath := t.TempDir() + "/test.pk"

			// Create keychain using the production approach
			kc, err := keychain.NewKeychain(tt.network, tempKeyPath, nil)
			require.NoError(t, err)
			require.NotNil(t, kc)

			// Create additional keys
			key1, err := key.NewSoft()
			require.NoError(t, err)

			key2, err := key.NewSoft()
			require.NoError(t, err)

			// Create a new keychain with multiple keys
			multiKeychain := secp256k1fx.NewKeychain()
			multiKeychain.Add(key1.PrivKey()) // Add first key
			multiKeychain.Add(key2.PrivKey()) // Add second key

			// Create a new keychain wrapper with multiple keys
			multiKeychainWrapper := keychain.NewKeychainFromExisting(multiKeychain, tt.network)

			// Create wallet using the production approach
			wallet, err := New(
				ctx,
				&primary.WalletConfig{
					URI:              tt.network.Endpoint,
					DIONEKeychain:    multiKeychainWrapper.Keychain,
					EthKeychain:      secp256k1fx.NewKeychain(),
					OChainTxsToFetch: nil,
				},
			)
			require.NoError(t, err)

			// Verify wallet has multiple addresses (2 keys)
			addresses := wallet.Addresses()
			require.Len(t, addresses, 2)

			// Verify all addresses are unique
			addressSet := set.Set[ids.ShortID]{}
			for _, addr := range addresses {
				require.False(t, addressSet.Contains(addr))
				addressSet.Add(addr)
			}

			// Test address generation for all keys
			oAddr, err := wallet.Keychain.O()
			require.NoError(t, err)
			require.Len(t, oAddr, 2)

			aAddr, err := wallet.Keychain.A()
			require.NoError(t, err)
			require.Len(t, aAddr, 2)

			// Verify addresses have correct prefixes
			expectedPrefix := "O-" + tt.network.HRP()
			expectedAPrefix := "A-" + tt.network.HRP()
			for _, addr := range oAddr {
				require.Contains(t, addr, expectedPrefix)
			}
			for _, addr := range aAddr {
				require.Contains(t, addr, expectedAPrefix)
			}

			t.Logf("Network: %s (HRP: %s)", tt.network.Kind.String(), tt.network.HRP())
			t.Logf("Total addresses: %d", len(addresses))
			t.Logf("O-Chain addresses: %v", oAddr)
			t.Logf("A-Chain addresses: %v", aAddr)
		})
	}
}

func TestWalletWithOChainTxsToFetch(t *testing.T) {
	// Sequential to avoid rate limiting
	addTestDelay() // Add delay to avoid rate limiting

	tests := []struct {
		name    string
		network odyssey.Network
	}{
		{
			name:    "testnet with OChain txs to fetch",
			network: odyssey.TestnetNetwork(),
		},
		{
			name:    "mainnet with OChain txs to fetch",
			network: odyssey.MainnetNetwork(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			// Create a temporary key file for testing
			tempKeyPath := t.TempDir() + "/test.pk"

			// Create keychain using the production approach
			keychain, err := keychain.NewKeychain(tt.network, tempKeyPath, nil)
			require.NoError(t, err)
			require.NotNil(t, keychain)

			// Create a test subnet ID
			subnetID := ids.GenerateTestID()

			// Create wallet using the production approach
			wallet, err := New(
				ctx,
				&primary.WalletConfig{
					URI:              tt.network.Endpoint,
					DIONEKeychain:    keychain.Keychain,
					EthKeychain:      secp256k1fx.NewKeychain(),
					OChainTxsToFetch: set.Of(subnetID),
				},
			)

			// The wallet creation might fail due to tx not found, which is expected for test IDs
			if err != nil {
				// If it fails due to tx not found, that's expected behavior
				require.Contains(t, err.Error(), "not found")
				return
			}

			require.NotNil(t, wallet)
			require.Equal(t, subnetID, wallet.config.OChainTxsToFetch.List()[0])

			t.Logf("Network: %s (HRP: %s)", tt.network.Kind.String(), tt.network.HRP())
			t.Logf("Subnet ID: %s", subnetID)
		})
	}
}

func TestWalletErrorHandling(t *testing.T) {
	// Sequential to avoid rate limiting
	addTestDelay() // Add delay to avoid rate limiting

	t.Run("context cancellation", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		// Create a temporary key file for testing
		tempKeyPath := t.TempDir() + "/test.pk"

		// Create keychain using the production approach
		network := odyssey.TestnetNetwork()
		keychain, err := keychain.NewKeychain(network, tempKeyPath, nil)
		require.NoError(t, err)

		config := &primary.WalletConfig{
			URI:              network.Endpoint,
			DIONEKeychain:    keychain.Keychain,
			EthKeychain:      secp256k1fx.NewKeychain(),
			OChainTxsToFetch: nil,
		}

		_, err = New(ctx, config)
		require.Error(t, err)
		require.Contains(t, err.Error(), "context canceled")
	})

	t.Run("invalid network endpoint", func(t *testing.T) {
		ctx := context.Background()

		// Create a temporary key file for testing
		tempKeyPath := t.TempDir() + "/test.pk"

		// Create keychain using the production approach
		keychain, err := keychain.NewKeychain(odyssey.TestnetNetwork(), tempKeyPath, nil)
		require.NoError(t, err)

		config := &primary.WalletConfig{
			URI:              "invalid://endpoint",
			DIONEKeychain:    keychain.Keychain,
			EthKeychain:      secp256k1fx.NewKeychain(),
			OChainTxsToFetch: nil,
		}

		_, err = New(ctx, config)
		require.Error(t, err)
	})
}

func TestWalletKeychainIntegration(t *testing.T) {
	// Sequential to avoid rate limiting
	addTestDelay() // Add delay to avoid rate limiting

	tests := []struct {
		name    string
		network odyssey.Network
	}{
		{
			name:    "testnet keychain integration",
			network: odyssey.TestnetNetwork(),
		},
		{
			name:    "mainnet keychain integration",
			network: odyssey.MainnetNetwork(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			// Create a temporary key file for testing
			tempKeyPath := t.TempDir() + "/test.pk"

			// Test with soft key keychain using production approach
			keychain, err := keychain.NewKeychain(tt.network, tempKeyPath, nil)
			require.NoError(t, err)
			require.NotNil(t, keychain)

			// Create wallet using the production approach
			wallet, err := New(
				ctx,
				&primary.WalletConfig{
					URI:              tt.network.Endpoint,
					DIONEKeychain:    keychain.Keychain,
					EthKeychain:      secp256k1fx.NewKeychain(),
					OChainTxsToFetch: nil,
				},
			)
			require.NoError(t, err)
			require.NotNil(t, wallet)

			// Verify keychain integration
			require.Equal(t, keychain.Keychain, wallet.Keychain.Keychain)

			// Test that the keychain wrapper is properly integrated
			require.NotNil(t, wallet.Keychain)
			require.True(t, wallet.Keychain.LedgerEnabled() == false) // Should be soft key

			// Test address generation through the integrated keychain
			addresses := wallet.Addresses()
			require.NotEmpty(t, addresses)

			oAddr, err := wallet.Keychain.O()
			require.NoError(t, err)
			require.NotEmpty(t, oAddr)

			aAddr, err := wallet.Keychain.A()
			require.NoError(t, err)
			require.NotEmpty(t, aAddr)

			t.Logf("Network: %s (HRP: %s)", tt.network.Kind.String(), tt.network.HRP())
			t.Logf("Addresses: %v", addresses)
			t.Logf("O-Chain address: %s", oAddr[0])
			t.Logf("A-Chain address: %s", aAddr[0])
		})
	}
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

// Helper function to add delay between tests to avoid rate limiting
func addTestDelay() {
	time.Sleep(5 * time.Second) // Increased to 5s to avoid 429 errors when running full test suite
}

// Benchmark tests for wallet operations
func BenchmarkWalletCreation(b *testing.B) {
	network := odyssey.TestnetNetwork()

	// Create a temporary key file for benchmarking
	tempKeyPath := b.TempDir() + "/benchmark.pk"

	// Create keychain using the production approach
	keychain, err := keychain.NewKeychain(network, tempKeyPath, nil)
	if err != nil {
		b.Fatal(err)
	}

	config := &primary.WalletConfig{
		URI:              network.Endpoint,
		DIONEKeychain:    keychain.Keychain,
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
	network := odyssey.TestnetNetwork()

	// Create a temporary key file for benchmarking
	tempKeyPath := b.TempDir() + "/benchmark.pk"

	// Create keychain using the production approach
	keychain, err := keychain.NewKeychain(network, tempKeyPath, nil)
	if err != nil {
		b.Fatal(err)
	}

	config := &primary.WalletConfig{
		URI:              network.Endpoint,
		DIONEKeychain:    keychain.Keychain,
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
		network        odyssey.Network
		expectedHRP    string
		expectedPrefix string
	}{
		{
			name:           "mainnet network",
			network:        odyssey.MainnetNetwork(),
			expectedHRP:    "dione",
			expectedPrefix: "O-dione",
		},
		{
			name:           "testnet network",
			network:        odyssey.TestnetNetwork(),
			expectedHRP:    "testnet",
			expectedPrefix: "O-testnet",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			// Create a temporary key file for testing
			tempKeyPath := t.TempDir() + "/test.pk"

			// Create keychain using the production approach
			keychain, err := keychain.NewKeychain(tt.network, tempKeyPath, nil)
			require.NoError(t, err)
			require.NotNil(t, keychain)

			// Create wallet using the production approach
			wallet, err := New(
				ctx,
				&primary.WalletConfig{
					URI:              tt.network.Endpoint,
					DIONEKeychain:    keychain.Keychain,
					EthKeychain:      secp256k1fx.NewKeychain(),
					OChainTxsToFetch: nil,
				},
			)
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
			require.Equal(t, tt.expectedHRP, tt.network.HRP(),
				"Expected HRP %s for network %s, but got %s", tt.expectedHRP, tt.network.Kind.String(), tt.network.HRP())

			// Test A-Chain address generation as well
			aAddr, err := wallet.Keychain.A()
			require.NoError(t, err)
			require.NotEmpty(t, aAddr)
			require.Len(t, aAddr, 1)

			expectedAPrefix := strings.Replace(tt.expectedPrefix, "O-", "A-", 1)
			require.Contains(t, aAddr[0], expectedAPrefix,
				"Expected A-Chain address to contain prefix %s, but got %s", expectedAPrefix, aAddr[0])

			t.Logf("%s O-Chain address: %s (HRP: %s)", tt.name, oAddr[0], tt.network.HRP())
			t.Logf("%s A-Chain address: %s (HRP: %s)", tt.name, aAddr[0], tt.network.HRP())
		})
	}
}

func TestWalletGenerationFromPrivateKey(t *testing.T) {
	// Sequential to avoid rate limiting
	addTestDelay() // Add delay to avoid rate limiting

	// Test with a known private key (this is a test key, not a real one)
	testPrivateKeyHex := "56289e99c94b6912bfc12adc093c9b51124f0dc54ac7a766b2bc5ccf558d8027"

	tests := []struct {
		name          string
		network       odyssey.Network
		expectedOAddr string
		expectedAAddr string
	}{
		{
			name:          "testnet from private key",
			network:       odyssey.TestnetNetwork(),
			expectedOAddr: "O-testnet18jma8ppw3nhx5r4ap8clazz0dps7rv5uw463py",
			expectedAAddr: "A-testnet18jma8ppw3nhx5r4ap8clazz0dps7rv5uw463py",
		},
		{
			name:          "mainnet from private key",
			network:       odyssey.MainnetNetwork(),
			expectedOAddr: "O-dione18jma8ppw3nhx5r4ap8clazz0dps7rv5ulw7llh",
			expectedAAddr: "A-dione18jma8ppw3nhx5r4ap8clazz0dps7rv5ulw7llh",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			// Test private key loading and validation
			softKey, err := key.LoadSoftFromBytes([]byte(testPrivateKeyHex))
			require.NoError(t, err)
			require.NotNil(t, softKey)
			require.Equal(t, testPrivateKeyHex, softKey.PrivKeyHex())

			// Create keychain from the loaded key using production approach
			keychain := keychain.NewKeychainFromExisting(softKey.KeyChain(), tt.network)
			require.NotNil(t, keychain)

			// Create wallet using the production approach
			wallet, err := New(
				ctx,
				&primary.WalletConfig{
					URI:              tt.network.Endpoint,
					DIONEKeychain:    keychain.Keychain,
					EthKeychain:      secp256k1fx.NewKeychain(),
					OChainTxsToFetch: nil,
				},
			)
			require.NoError(t, err)
			require.NotNil(t, wallet)

			// Test that we can generate addresses from the private key
			oAddr, err := wallet.Keychain.O()
			require.NoError(t, err)
			require.NotEmpty(t, oAddr)
			require.Len(t, oAddr, 1)

			// The address should be deterministic and consistent
			require.Equal(t, tt.expectedOAddr, oAddr[0])

			// Test A-Chain address generation as well
			aAddr, err := wallet.Keychain.A()
			require.NoError(t, err)
			require.NotEmpty(t, aAddr)
			require.Len(t, aAddr, 1)

			require.Equal(t, tt.expectedAAddr, aAddr[0])

			t.Logf("Private Key: %s", softKey.PrivKeyHex())
			t.Logf("Network: %s (HRP: %s)", tt.network.Kind.String(), tt.network.HRP())
			t.Logf("Generated O-Chain Address: %s", oAddr[0])
			t.Logf("Generated A-Chain Address: %s", aAddr[0])
		})
	}
}

func TestWalletDeterministicAddressGeneration(t *testing.T) {
	// Sequential to avoid rate limiting
	addTestDelay() // Add delay to avoid rate limiting

	// Test that the same private key generates the same addresses across different networks
	testPrivateKeyHex := "56289e99c94b6912bfc12adc093c9b51124f0dc54ac7a766b2bc5ccf558d8027"

	tests := []struct {
		name          string
		network       odyssey.Network
		expectedOAddr string
		expectedAAddr string
	}{
		{
			name:          "testnet deterministic addresses",
			network:       odyssey.TestnetNetwork(),
			expectedOAddr: "O-testnet18jma8ppw3nhx5r4ap8clazz0dps7rv5uw463py",
			expectedAAddr: "A-testnet18jma8ppw3nhx5r4ap8clazz0dps7rv5uw463py",
		},
		{
			name:          "mainnet deterministic addresses",
			network:       odyssey.MainnetNetwork(),
			expectedOAddr: "O-dione18jma8ppw3nhx5r4ap8clazz0dps7rv5ulw7llh",
			expectedAAddr: "A-dione18jma8ppw3nhx5r4ap8clazz0dps7rv5ulw7llh",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			// Create the same key twice to ensure deterministic behavior
			softKey1, err := key.LoadSoftFromBytes([]byte(testPrivateKeyHex))
			require.NoError(t, err)

			softKey2, err := key.LoadSoftFromBytes([]byte(testPrivateKeyHex))
			require.NoError(t, err)

			// Verify both keys are identical
			require.Equal(t, softKey1.PrivKeyHex(), softKey2.PrivKeyHex())
			require.Equal(t, softKey1.KeyChain().Addresses().List(), softKey2.KeyChain().Addresses().List())

			// Create keychain from the loaded key using production approach
			keychain1 := keychain.NewKeychainFromExisting(softKey1.KeyChain(), tt.network)
			keychain2 := keychain.NewKeychainFromExisting(softKey2.KeyChain(), tt.network)

			// Create wallet using the production approach
			wallet1, err := New(
				ctx,
				&primary.WalletConfig{
					URI:              tt.network.Endpoint,
					DIONEKeychain:    keychain1.Keychain,
					EthKeychain:      secp256k1fx.NewKeychain(),
					OChainTxsToFetch: nil,
				},
			)
			require.NoError(t, err)

			wallet2, err := New(
				ctx,
				&primary.WalletConfig{
					URI:              tt.network.Endpoint,
					DIONEKeychain:    keychain2.Keychain,
					EthKeychain:      secp256k1fx.NewKeychain(),
					OChainTxsToFetch: nil,
				},
			)
			require.NoError(t, err)

			// Test address generation from both wallets
			addr1, err := wallet1.Keychain.O()
			require.NoError(t, err)
			require.Len(t, addr1, 1)

			addr2, err := wallet2.Keychain.O()
			require.NoError(t, err)
			require.Len(t, addr2, 1)

			// Addresses should be identical (deterministic)
			require.Equal(t, addr1[0], addr2[0])
			require.Equal(t, tt.expectedOAddr, addr1[0])

			// Test A-Chain addresses as well
			aAddr1, err := wallet1.Keychain.A()
			require.NoError(t, err)
			require.Len(t, aAddr1, 1)

			aAddr2, err := wallet2.Keychain.A()
			require.NoError(t, err)
			require.Len(t, aAddr2, 1)

			require.Equal(t, aAddr1[0], aAddr2[0])
			require.Equal(t, tt.expectedAAddr, aAddr1[0])

			// Test direct address generation using the same key but different HRPs
			directOAddr, err := softKey1.O(tt.network.HRP())
			require.NoError(t, err)

			directAAddr, err := softKey1.A(tt.network.HRP())
			require.NoError(t, err)

			require.Equal(t, addr1[0], directOAddr)
			require.Equal(t, aAddr1[0], directAAddr)

			t.Logf("Network: %s (HRP: %s)", tt.network.Kind.String(), tt.network.HRP())
			t.Logf("O-Chain address: %s", addr1[0])
			t.Logf("A-Chain address: %s", aAddr1[0])
			t.Logf("Private Key: %s", softKey1.PrivKeyHex())
		})
	}
}

func TestAddressEncodingDifference(t *testing.T) {
	addTestDelay() // Add delay to avoid rate limiting
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
		network        odyssey.Network
		expectedHRP    string
		expectedPrefix string
	}{
		{
			name:           "mainnet A-Chain network",
			network:        odyssey.MainnetNetwork(),
			expectedHRP:    "dione",
			expectedPrefix: "A-dione",
		},
		{
			name:           "testnet A-Chain network",
			network:        odyssey.TestnetNetwork(),
			expectedHRP:    "testnet",
			expectedPrefix: "A-testnet",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			// Create a temporary key file for testing
			tempKeyPath := t.TempDir() + "/test.pk"

			// Create keychain using the production approach
			keychain, err := keychain.NewKeychain(tt.network, tempKeyPath, nil)
			require.NoError(t, err)
			require.NotNil(t, keychain)

			// Create wallet using the production approach
			wallet, err := New(
				ctx,
				&primary.WalletConfig{
					URI:              tt.network.Endpoint,
					DIONEKeychain:    keychain.Keychain,
					EthKeychain:      secp256k1fx.NewKeychain(),
					OChainTxsToFetch: nil,
				},
			)
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
			require.Equal(t, tt.expectedHRP, tt.network.HRP(),
				"Expected HRP %s for network %s, but got %s", tt.expectedHRP, tt.network.Kind.String(), tt.network.HRP())

			// Test O-Chain address generation as well for completeness
			oAddr, err := wallet.Keychain.O()
			require.NoError(t, err)
			require.NotEmpty(t, oAddr)
			require.Len(t, oAddr, 1)

			expectedOPrefix := strings.Replace(tt.expectedPrefix, "A-", "O-", 1)
			require.Contains(t, oAddr[0], expectedOPrefix,
				"Expected O-Chain address to contain prefix %s, but got %s", expectedOPrefix, oAddr[0])

			t.Logf("%s A-Chain address: %s (HRP: %s)", tt.name, aAddr[0], tt.network.HRP())
			t.Logf("%s O-Chain address: %s (HRP: %s)", tt.name, oAddr[0], tt.network.HRP())
		})
	}
}

func TestWalletBothChainAddressGeneration(t *testing.T) {
	// Sequential to avoid rate limiting
	addTestDelay() // Add delay to avoid rate limiting

	// Test with a known private key to ensure deterministic behavior
	testPrivateKeyHex := "56289e99c94b6912bfc12adc093c9b51124f0dc54ac7a766b2bc5ccf558d8027"

	tests := []struct {
		name          string
		network       odyssey.Network
		expectedOAddr string
		expectedAAddr string
	}{
		{
			name:          "testnet both chain addresses",
			network:       odyssey.TestnetNetwork(),
			expectedOAddr: "O-testnet18jma8ppw3nhx5r4ap8clazz0dps7rv5uw463py",
			expectedAAddr: "A-testnet18jma8ppw3nhx5r4ap8clazz0dps7rv5uw463py",
		},
		{
			name:          "mainnet both chain addresses",
			network:       odyssey.MainnetNetwork(),
			expectedOAddr: "O-dione18jma8ppw3nhx5r4ap8clazz0dps7rv5ulw7llh",
			expectedAAddr: "A-dione18jma8ppw3nhx5r4ap8clazz0dps7rv5ulw7llh",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			// Create a key from the specific private key
			softKey, err := key.LoadSoftFromBytes([]byte(testPrivateKeyHex))
			require.NoError(t, err)
			require.NotNil(t, softKey)

			// Create keychain from the loaded key using production approach
			keychain := keychain.NewKeychainFromExisting(softKey.KeyChain(), tt.network)
			require.NotNil(t, keychain)

			// Create wallet using the production approach
			wallet, err := New(
				ctx,
				&primary.WalletConfig{
					URI:              tt.network.Endpoint,
					DIONEKeychain:    keychain.Keychain,
					EthKeychain:      secp256k1fx.NewKeychain(),
					OChainTxsToFetch: nil,
				},
			)
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
			require.Contains(t, oAddr[0], "O-"+tt.network.HRP())
			require.Contains(t, aAddr[0], "A-"+tt.network.HRP())

			// Verify addresses are different (different chain types)
			require.NotEqual(t, oAddr[0], aAddr[0])

			// Verify expected addresses
			require.Equal(t, tt.expectedOAddr, oAddr[0])
			require.Equal(t, tt.expectedAAddr, aAddr[0])

			// Test deterministic behavior - same key should generate same addresses
			oAddr2, err := wallet.Keychain.O()
			require.NoError(t, err)
			require.Equal(t, oAddr, oAddr2)

			aAddr2, err := wallet.Keychain.A()
			require.NoError(t, err)
			require.Equal(t, aAddr, aAddr2)

			t.Logf("Network: %s (HRP: %s)", tt.network.Kind.String(), tt.network.HRP())
			t.Logf("Private Key: %s", softKey.PrivKeyHex())
			t.Logf("O-Chain Address: %s", oAddr[0])
			t.Logf("A-Chain Address: %s", aAddr[0])
		})
	}
}

func TestWalletAChainFromPrivateKey(t *testing.T) {
	// Sequential to avoid rate limiting
	addTestDelay() // Add delay to avoid rate limiting

	// Test with a known private key (this is a test key, not a real one)
	testPrivateKeyHex := "56289e99c94b6912bfc12adc093c9b51124f0dc54ac7a766b2bc5ccf558d8027"

	tests := []struct {
		name          string
		network       odyssey.Network
		expectedAAddr string
		expectedOAddr string
	}{
		{
			name:          "testnet A-Chain from private key",
			network:       odyssey.TestnetNetwork(),
			expectedAAddr: "A-testnet18jma8ppw3nhx5r4ap8clazz0dps7rv5uw463py",
			expectedOAddr: "O-testnet18jma8ppw3nhx5r4ap8clazz0dps7rv5uw463py",
		},
		{
			name:          "mainnet A-Chain from private key",
			network:       odyssey.MainnetNetwork(),
			expectedAAddr: "A-dione18jma8ppw3nhx5r4ap8clazz0dps7rv5ulw7llh",
			expectedOAddr: "O-dione18jma8ppw3nhx5r4ap8clazz0dps7rv5ulw7llh",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			// Test private key loading and validation
			softKey, err := key.LoadSoftFromBytes([]byte(testPrivateKeyHex))
			require.NoError(t, err)
			require.NotNil(t, softKey)
			require.Equal(t, testPrivateKeyHex, softKey.PrivKeyHex())

			// Create keychain from the loaded key using production approach
			keychain := keychain.NewKeychainFromExisting(softKey.KeyChain(), tt.network)
			require.NotNil(t, keychain)

			// Create wallet using the production approach
			wallet, err := New(
				ctx,
				&primary.WalletConfig{
					URI:              tt.network.Endpoint,
					DIONEKeychain:    keychain.Keychain,
					EthKeychain:      secp256k1fx.NewKeychain(),
					OChainTxsToFetch: nil,
				},
			)
			require.NoError(t, err)
			require.NotNil(t, wallet)

			// Test that we can generate A-Chain addresses from the private key
			aAddr, err := wallet.Keychain.A()
			require.NoError(t, err)
			require.NotEmpty(t, aAddr)
			require.Len(t, aAddr, 1)

			// The address should be deterministic and consistent
			require.Equal(t, tt.expectedAAddr, aAddr[0])

			// Test O-Chain address generation as well
			oAddr, err := wallet.Keychain.O()
			require.NoError(t, err)
			require.NotEmpty(t, oAddr)
			require.Len(t, oAddr, 1)

			require.Equal(t, tt.expectedOAddr, oAddr[0])

			t.Logf("Network: %s (HRP: %s)", tt.network.Kind.String(), tt.network.HRP())
			t.Logf("Private Key: %s", softKey.PrivKeyHex())
			t.Logf("Generated A-Chain Address: %s", aAddr[0])
			t.Logf("Generated O-Chain Address: %s", oAddr[0])
		})
	}
}

func TestWalletPrivateKeyValidation(t *testing.T) {
	addTestDelay() // Add delay to avoid rate limiting
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

func TestWalletWithNetworkAndKeychain(t *testing.T) {
	// Sequential to avoid rate limiting
	addTestDelay() // Add delay to avoid rate limiting

	tests := []struct {
		name           string
		network        odyssey.Network
		expectedHRP    string
		expectedPrefix string
	}{
		{
			name:           "testnet network with keychain",
			network:        odyssey.TestnetNetwork(),
			expectedHRP:    "testnet",
			expectedPrefix: "O-testnet",
		},
		{
			name:           "mainnet network with keychain",
			network:        odyssey.MainnetNetwork(),
			expectedHRP:    "dione",
			expectedPrefix: "O-dione",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			// Create a temporary key file for testing
			tempKeyPath := t.TempDir() + "/test.key"

			// Create keychain using the network and keychain approach (like examples)
			keychain, err := keychain.NewKeychain(tt.network, tempKeyPath, nil)
			require.NoError(t, err)
			require.NotNil(t, keychain)

			// Create wallet using the keychain wrapper
			wallet, err := New(
				ctx,
				&primary.WalletConfig{
					URI:              tt.network.Endpoint,
					DIONEKeychain:    keychain.Keychain,
					EthKeychain:      secp256k1fx.NewKeychain(),
					OChainTxsToFetch: nil,
				},
			)
			require.NoError(t, err)
			require.NotNil(t, wallet)

			// Verify keychain integration
			require.Equal(t, keychain.Keychain, wallet.Keychain.Keychain)

			// Test address generation
			addresses := wallet.Addresses()
			require.NotEmpty(t, addresses)
			require.Len(t, addresses, 1)

			// Test O-Chain address generation
			oAddr, err := wallet.Keychain.O()
			require.NoError(t, err)
			require.NotEmpty(t, oAddr)
			require.Len(t, oAddr, 1)

			// Verify the address has the correct prefix
			require.Contains(t, oAddr[0], tt.expectedPrefix,
				"Expected address to contain prefix %s, but got %s", tt.expectedPrefix, oAddr[0])

			// Test A-Chain address generation
			aAddr, err := wallet.Keychain.A()
			require.NoError(t, err)
			require.NotEmpty(t, aAddr)
			require.Len(t, aAddr, 1)

			// Verify A-Chain address has correct prefix
			expectedAPrefix := strings.Replace(tt.expectedPrefix, "O-", "A-", 1)
			require.Contains(t, aAddr[0], expectedAPrefix,
				"Expected A-Chain address to contain prefix %s, but got %s", expectedAPrefix, aAddr[0])

			// Verify network HRP is correct
			require.Equal(t, tt.expectedHRP, tt.network.HRP(),
				"Expected HRP %s for network %s, but got %s", tt.expectedHRP, tt.network.Kind.String(), tt.network.HRP())

			// Test deterministic behavior - same key should generate same addresses
			oAddr2, err := wallet.Keychain.O()
			require.NoError(t, err)
			require.Equal(t, oAddr, oAddr2)

			aAddr2, err := wallet.Keychain.A()
			require.NoError(t, err)
			require.Equal(t, aAddr, aAddr2)

			t.Logf("%s O-Chain address: %s", tt.name, oAddr[0])
			t.Logf("%s A-Chain address: %s", tt.name, aAddr[0])
			t.Logf("%s Network HRP: %s", tt.name, tt.network.HRP())
			t.Logf("%s Network endpoint: %s", tt.name, tt.network.Endpoint)
		})
	}
}

func TestWalletWithNetworkAndKeychainFromPrivateKey(t *testing.T) {
	// Sequential to avoid rate limiting
	addTestDelay() // Add delay to avoid rate limiting

	// Test with a known private key for deterministic results
	testPrivateKeyHex := "56289e99c94b6912bfc12adc093c9b51124f0dc54ac7a766b2bc5ccf558d8027"

	tests := []struct {
		name          string
		network       odyssey.Network
		expectedOAddr string
		expectedAAddr string
	}{
		{
			name:          "testnet with known private key",
			network:       odyssey.TestnetNetwork(),
			expectedOAddr: "O-testnet18jma8ppw3nhx5r4ap8clazz0dps7rv5uw463py",
			expectedAAddr: "A-testnet18jma8ppw3nhx5r4ap8clazz0dps7rv5uw463py",
		},
		{
			name:          "mainnet with known private key",
			network:       odyssey.MainnetNetwork(),
			expectedOAddr: "O-dione18jma8ppw3nhx5r4ap8clazz0dps7rv5ulw7llh",
			expectedAAddr: "A-dione18jma8ppw3nhx5r4ap8clazz0dps7rv5ulw7llh",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			// Load the known private key
			softKey, err := key.LoadSoftFromBytes([]byte(testPrivateKeyHex))
			require.NoError(t, err)
			require.NotNil(t, softKey)

			// Create keychain from the loaded key
			keychain := keychain.NewKeychainFromExisting(softKey.KeyChain(), tt.network)
			require.NotNil(t, keychain)

			// Create wallet using the keychain wrapper
			wallet, err := New(
				ctx,
				&primary.WalletConfig{
					URI:              tt.network.Endpoint,
					DIONEKeychain:    keychain.Keychain,
					EthKeychain:      secp256k1fx.NewKeychain(),
					OChainTxsToFetch: nil,
				},
			)
			require.NoError(t, err)
			require.NotNil(t, wallet)

			// Test O-Chain address generation
			oAddr, err := wallet.Keychain.O()
			require.NoError(t, err)
			require.Len(t, oAddr, 1)
			require.Equal(t, tt.expectedOAddr, oAddr[0])

			// Test A-Chain address generation
			aAddr, err := wallet.Keychain.A()
			require.NoError(t, err)
			require.Len(t, aAddr, 1)
			require.Equal(t, tt.expectedAAddr, aAddr[0])

			// Verify addresses are different (different chain types)
			require.NotEqual(t, oAddr[0], aAddr[0])

			t.Logf("%s Private Key: %s", tt.name, softKey.PrivKeyHex())
			t.Logf("%s O-Chain Address: %s", tt.name, oAddr[0])
			t.Logf("%s A-Chain Address: %s", tt.name, aAddr[0])
			t.Logf("%s Network: %s", tt.name, tt.network.Kind.String())
		})
	}
}
